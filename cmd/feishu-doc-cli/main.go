package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/holtmiu/ChatGPT_MCP_Connectors/internal/config"
	"github.com/holtmiu/ChatGPT_MCP_Connectors/internal/feishu"
)

func main() {
	if len(os.Args) < 3 {
		usage()
		os.Exit(2)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	service := feishu.NewService(config.Load())
	command := os.Args[1]
	input := os.Args[2]

	var result any
	var err error
	switch command {
	case "resolve":
		result, err = service.Resolve(input)
	case "metadata":
		result, err = service.GetMetadata(ctx, input)
	case "read":
		result, err = service.ReadDocument(ctx, input, feishu.ReadOptions{Format: "both"})
	case "append-dry-run":
		dryRun := true
		markdown := ""
		if len(os.Args) > 3 {
			markdown = os.Args[3]
		}
		result, err = service.AppendDocument(ctx, input, feishu.AppendRequest{Markdown: markdown, DryRun: &dryRun})
	default:
		usage()
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(result)
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: feishu-doc-cli <resolve|metadata|read|append-dry-run> <url-or-token> [markdown]")
}
