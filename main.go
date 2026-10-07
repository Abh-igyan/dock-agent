package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/sruja/dock-agent/pkg/agent"
)

func main() {
	dir := flag.String("dir", ".", "Directory to analyze and dockerize")
	maxRetries := flag.Int("retries", 3, "Maximum number of self-healing retries")
	flag.Parse()

	fmt.Printf("🚀 Starting DockAgent on directory: %s\n", *dir)

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Println("❌ Error: GEMINI_API_KEY environment variable is required.")
		os.Exit(1)
	}

	ctx := context.Background()
	loop := agent.NewLoop(apiKey, *dir, *maxRetries)

	err := loop.Run(ctx)
	if err != nil {
		fmt.Printf("\n❌ Agent failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✅ Agent successfully containerized the project!")
}
