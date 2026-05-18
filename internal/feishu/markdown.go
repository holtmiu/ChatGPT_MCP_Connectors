package feishu

import "strings"

func exportMarkdown(blocks []NormalizedBlock) string {
	var b strings.Builder
	for _, block := range blocks {
		writeMarkdownBlock(&b, block, 0)
	}
	return strings.TrimSpace(b.String()) + "\n"
}

func writeMarkdownBlock(b *strings.Builder, block NormalizedBlock, depth int) {
	indent := strings.Repeat("  ", depth)
	text := strings.TrimSpace(block.Text)
	switch block.Type {
	case "heading":
		level := intAttr(block.Attrs, "level", 2)
		if level < 1 || level > 9 {
			level = 2
		}
		b.WriteString(strings.Repeat("#", level))
		b.WriteString(" ")
		b.WriteString(text)
		b.WriteString("\n\n")
	case "bullet_list":
		b.WriteString(indent)
		b.WriteString("- ")
		b.WriteString(text)
		b.WriteString("\n")
	case "ordered_list":
		b.WriteString(indent)
		b.WriteString("1. ")
		b.WriteString(text)
		b.WriteString("\n")
	case "todo_list":
		b.WriteString(indent)
		b.WriteString("- [ ] ")
		b.WriteString(text)
		b.WriteString("\n")
	case "code_block":
		b.WriteString("```\n")
		b.WriteString(block.Text)
		if !strings.HasSuffix(block.Text, "\n") {
			b.WriteString("\n")
		}
		b.WriteString("```\n\n")
	case "quote":
		for _, line := range strings.Split(text, "\n") {
			b.WriteString("> ")
			b.WriteString(line)
			b.WriteString("\n")
		}
		b.WriteString("\n")
	case "divider":
		b.WriteString("---\n\n")
	case "image", "file":
		token := stringAttr(block.Attrs, "token")
		if token == "" {
			token = block.ID
		}
		b.WriteString("[")
		b.WriteString(block.Type)
		b.WriteString(": ")
		b.WriteString(token)
		b.WriteString("]\n\n")
	case "table":
		b.WriteString("```json\n")
		b.WriteString("{\"type\":\"table\",\"id\":\"")
		b.WriteString(block.ID)
		b.WriteString("\"}\n")
		b.WriteString("```\n\n")
	case "unsupported":
		b.WriteString("<!-- unsupported block")
		if block.Source != nil && block.Source.RawType != "" {
			b.WriteString(": ")
			b.WriteString(block.Source.RawType)
		}
		b.WriteString(" -->\n\n")
	default:
		if text != "" {
			b.WriteString(text)
			b.WriteString("\n\n")
		}
	}

	for _, child := range block.Children {
		writeMarkdownBlock(b, child, depth+1)
	}

	if block.Type == "bullet_list" || block.Type == "ordered_list" || block.Type == "todo_list" {
		b.WriteString("\n")
	}
}

func intAttr(attrs map[string]any, key string, fallback int) int {
	if attrs == nil {
		return fallback
	}
	switch v := attrs[key].(type) {
	case int:
		return v
	case float64:
		return int(v)
	default:
		return fallback
	}
}

func stringAttr(attrs map[string]any, key string) string {
	if attrs == nil {
		return ""
	}
	if v, ok := attrs[key].(string); ok {
		return v
	}
	return ""
}
