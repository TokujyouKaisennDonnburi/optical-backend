package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

func main() {
	_ = godotenv.Load()
	model, ok := os.LookupEnv("AGENT_MODEL")
	if !ok {
		panic("'AGENT_MODEL' is not set")
	}
	apiKey, ok := os.LookupEnv("AGENT_API_KEY")
	if !ok {
		panic("'AGENT_API_KEY' is not set")
	}
	// Initialize LLM
	llm, err := openai.New(
		openai.WithBaseURL("https://openrouter.ai/api/v1"),
		openai.WithModel(model),
		openai.WithToken(apiKey),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("AIチャット ('quit' で終了)")
	fmt.Println("----------------------------------------")

	for {
		fmt.Print("--> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "quit" {
			break
		}

		_, err := llms.GenerateFromSinglePrompt(
			ctx,
			llm,
			input,
			llms.WithStreamingFunc(stdStream),
		)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Printf("\n")
	}
}

func stdStream(ctx context.Context, chunk []byte) error {
	fmt.Print(string(chunk))
	return nil
}
