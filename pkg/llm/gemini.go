package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiClient struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewGeminiClient(ctx context.Context, apiKey string) (*GeminiClient, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %v", err)
	}

	model := client.GenerativeModel("gemini-2.5-pro") // Use latest capable model
	model.SetTemperature(0.2) // Low temp for more deterministic code

	return &GeminiClient{
		client: client,
		model:  model,
	}, nil
}

func (g *GeminiClient) Close() {
	g.client.Close()
}

func (g *GeminiClient) GenerateDockerfile(ctx context.Context, contextInfo string) (string, error) {
	prompt := fmt.Sprintf(`You are an expert DevOps engineer. 
I want you to write a best-practice Dockerfile for the following project structure.
Output ONLY the raw Dockerfile content. Do not include markdown formatting like `+"```docker"+` or explanations.

Project context/files:
%s`, contextInfo)

	return g.generate(ctx, prompt)
}

func (g *GeminiClient) FixDockerfile(ctx context.Context, currentDockerfile, buildError string) (string, error) {
	prompt := fmt.Sprintf(`You are an expert DevOps engineer.
The following Dockerfile failed to build.

Current Dockerfile:
%s

Build Error:
%s

Analyze the error and fix the Dockerfile.
Output ONLY the raw fixed Dockerfile content. Do not include markdown formatting or explanations.`, currentDockerfile, buildError)

	return g.generate(ctx, prompt)
}

func (g *GeminiClient) generate(ctx context.Context, prompt string) (string, error) {
	resp, err := g.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %v", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content generated")
	}

	part := resp.Candidates[0].Content.Parts[0]
	text, ok := part.(genai.Text)
	if !ok {
		return "", fmt.Errorf("expected text part from model")
	}

	// Clean up potential markdown formatting that the model might stubbornly add
	cleaned := string(text)
	cleaned = strings.TrimSpace(cleaned)
	if strings.HasPrefix(cleaned, "```docker") {
		cleaned = strings.TrimPrefix(cleaned, "```docker")
	} else if strings.HasPrefix(cleaned, "```") {
		cleaned = strings.TrimPrefix(cleaned, "```")
	}
	if strings.HasSuffix(cleaned, "```") {
		cleaned = strings.TrimSuffix(cleaned, "```")
	}
	cleaned = strings.TrimSpace(cleaned)

	return cleaned, nil
}
