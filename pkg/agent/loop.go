package agent

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Abh-igyan/dock-agent/pkg/docker"
	"github.com/Abh-igyan/dock-agent/pkg/llm"
)

type Loop struct {
	llmClient  *llm.GeminiClient
	builder    *docker.Builder
	dir        string
	maxRetries int
}

func NewLoop(apiKey, dir string, maxRetries int) *Loop {
	return &Loop{
		dir:        dir,
		maxRetries: maxRetries,
	}
}

func (l *Loop) Run(ctx context.Context) error {
	client, err := llm.NewGeminiClient(ctx, os.Getenv("GEMINI_API_KEY"))
	if err != nil {
		return err
	}
	defer client.Close()

	l.llmClient = client
	l.builder = docker.NewBuilder(l.dir)

	fmt.Println("🔍 Analyzing project context...")
	contextInfo, err := l.gatherContext()
	if err != nil {
		return fmt.Errorf("failed to gather context: %v", err)
	}

	fmt.Println("🧠 Generating initial Dockerfile...")
	dockerfile, err := l.llmClient.GenerateDockerfile(ctx, contextInfo)
	if err != nil {
		return fmt.Errorf("failed to generate Dockerfile: %v", err)
	}

	for attempt := 1; attempt <= l.maxRetries; attempt++ {
		fmt.Printf("🔨 Attempt %d/%d: Building image...\n", attempt, l.maxRetries)
		err := l.builder.Build(dockerfile)
		if err == nil {
			fmt.Println("🎉 Build succeeded!")
			return nil
		}

		fmt.Printf("⚠️ Build failed! Error:\n%v\n", err)
		if attempt == l.maxRetries {
			return fmt.Errorf("max retries reached. Last error: %v", err)
		}

		fmt.Println("🧠 Asking LLM for a fix...")
		dockerfile, err = l.llmClient.FixDockerfile(ctx, dockerfile, err.Error())
		if err != nil {
			return fmt.Errorf("failed to get fix from LLM: %v", err)
		}
	}

	return nil
}

func (l *Loop) gatherContext() (string, error) {
	var sb strings.Builder
	
	err := filepath.Walk(l.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip hidden dirs and common heavy dirs
		if info.IsDir() && (strings.HasPrefix(info.Name(), ".") || info.Name() == "node_modules" || info.Name() == "vendor" || info.Name() == "build" || info.Name() == "dist") {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			rel, _ := filepath.Rel(l.dir, path)
			sb.WriteString(fmt.Sprintf("- %s\n", rel))
			
			// If it's a configuration file, let's include its content for better context
			if isImportantConfig(info.Name()) {
				content, err := os.ReadFile(path)
				if err == nil {
					sb.WriteString(fmt.Sprintf("--- CONTENT OF %s ---\n%s\n--- END ---\n", rel, string(content)))
				}
			}
		}
		return nil
	})

	return sb.String(), err
}

func isImportantConfig(name string) bool {
	configs := []string{"package.json", "go.mod", "requirements.txt", "pom.xml", "Cargo.toml", "Gemfile"}
	for _, c := range configs {
		if name == c {
			return true
		}
	}
	return false
}

