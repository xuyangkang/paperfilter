package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

type PaperFilterResponse struct {
	Match         bool   `json:"match"`
	Justification string `json:"justification"`
}

func main() {
	inputPath := flag.String("input", "", "Path to the input file containing the paper text. If not set, reads from stdin.")
	apiKey := flag.String("apikey", "", "API key for Grok (OpenAI compatible)")
	baseURL := flag.String("baseurl", "https://api.x.ai/v1", "Base URL for the API")
	model := flag.String("model", "grok-4-1-fast-reasoning", "Model to use for completion")
	flag.Parse()

	if *apiKey == "" {
		log.Fatal("Error: --apikey flag is required")
	}

	var paperText []byte
	var err error

	if *inputPath != "" {
		paperText, err = os.ReadFile(*inputPath)
		if err != nil {
			log.Fatalf("Error reading input file: %v", err)
		}
	} else {
		// Verify if stdin has data
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			paperText, err = io.ReadAll(os.Stdin)
			if err != nil {
				log.Fatalf("Error reading from stdin: %v", err)
			}
		} else {
			log.Fatal("Error: No input provided via stdin or --input flag")
		}
	}

	if len(paperText) == 0 {
		log.Fatal("Error: Input text is empty")
	}

	ctx := context.Background()

	config := openai.DefaultConfig(*apiKey)
	config.BaseURL = *baseURL
	client := openai.NewClientWithConfig(config)

	systemPrompt := `You are an AI research assistant. Your task is to analyze the provided abstract or text of an academic paper and determine if it aligns with the user's research interests.

The user's research interests includes:
- String algorithms
- Data structures
But the user is new to this field.

The user's skills include:
- Parallel computing, including OpenMP programming, MPI programming, and GPU programming
- Multi-agent system development
So user can potentially improve the existing paper by applying his skills to it.

Please evaluate the paper and respond in JSON format with two fields:
1. "match" (boolean): true if the paper is highly relevant to the research interests, false otherwise.
2. "justification" (string): a brief explanation of why the paper matches or does not match.`

	// Define JSON schema for structured output
	schema := jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"match":         {Type: jsonschema.Boolean, Description: "Whether the paper matches the research interests"},
			"justification": {Type: jsonschema.String, Description: "A brief justification for the match or non-match"},
		},
		Required: []string{"match", "justification"},
	}

	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: *model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: string(paperText),
			},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
			JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
				Name:        "PaperFilterResponse",
				Description: "Structured response indicating whether a paper matches research interests",
				Schema:      &schema,
				Strict:      true,
			},
		},
		Temperature: 0.1,
	})

	if err != nil {
		log.Fatalf("Error calling Grok API: %v", err)
	}

	if len(resp.Choices) == 0 {
		log.Fatal("Error: No choices returned from API")
	}

	content := resp.Choices[0].Message.Content

	var filterResp PaperFilterResponse
	err = json.Unmarshal([]byte(content), &filterResp)
	if err != nil {
		log.Fatalf("Error parsing JSON response: %v\nRaw response: %s", err, content)
	}

	// Output the parsed response in JSON to stdout
	out, err := json.MarshalIndent(filterResp, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling output: %v", err)
	}

	fmt.Println(string(out))
}
