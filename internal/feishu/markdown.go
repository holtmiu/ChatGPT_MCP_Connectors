package feishu

import (
	"strings"
)

func exportMarkdown(blocks []NormalizedBlock) string {
	var b strings.Builder
	for _, block := range blocks {
		writeMarkdownBlock(&b, block, 0)
	}
	out := strings.TrimSpace(b.String())
	if out == "" {
		return ""
	}
	return out + "\n"
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

func markdownToBlocks(markdown string) []NormalizedBlock {
	lines := strings.Split(strings.ReplaceAll(markdown, "\r\n", "\n"), "\n")
	blocks := make([]NormalizedBlock, 0, len(lines))
	inCode := false
	var code strings.Builder
	for _, raw := range lines {
		line := strings.TrimRight(raw, " \t")
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			if inCode {
				blocks = append(blocks, NormalizedBlock{Type: "code_block", Text: strings.TrimRight(code.String(), "\n")})
				code.Reset()
				inCode = false
			} else {
				inCode = true
			}
			continue
		}
		if inCode {
			code.WriteString(line)
			code.WriteByte('\n')
			continue
		}
		if trimmed == "" {
			continue
		}
		if trimmed == "---" || trimmed == "***" {
			blocks = append(blocks, NormalizedBlock{Type: "divider"})
			continue
		}
		if level, text := parseHeading(trimmed); level > 0 {
			blocks = append(blocks, NormalizedBlock{Type: "heading", Text: text, Attrs: map[string]any{"level": level}})
			continue
		}
		if strings.HasPrefix(trimmed, "- [ ] ") || strings.HasPrefix(trimmed, "- [x] ") || strings.HasPrefix(trimmed, "- [X] ") {
			blocks = append(blocks, NormalizedBlock{Type: "todo_list", Text: strings.TrimSpace(trimmed[6:])})
			continue
		}
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			blocks = append(blocks, NormalizedBlock{Type: "bullet_list", Text: strings.TrimSpace(trimmed[2:])})
			continue
		}
		if idx := orderedListIndex(trimmed); idx > 0 {
			blocks = append(blocks, NormalizedBlock{Type: "ordered_list", Text: strings.TrimSpace(trimmed[idx:])})
			continue
		}
		if strings.HasPrefix(trimmed, "> ") {
			blocks = append(blocks, NormalizedBlock{Type: "quote", Text: strings.TrimSpace(trimmed[2:])})
			continue
		}
		blocks = append(blocks, NormalizedBlock{Type: "paragraph", Text: trimmed})
	}
	if inCode && code.Len() > 0 {
		blocks = append(blocks, NormalizedBlock{Type: "code_block", Text: strings.TrimRight(code.String(), "\n")})
	}
	return blocks
}

func parseHeading(line string) (int, string) {
	count := 0
	for count < len(line) && line[count] == '#' {
		count++
	}
	if count == 0 || count > 9 || count >= len(line) || line[count] != ' ' {
		return 0, ""
	}
	return count, strings.TrimSpace(line[count:])
}

func orderedListIndex(line string) int {
	for i, r := range line {
		if r >= '0' && r <= '9' {
			continue
		}
		if r == '.' && i+1 < len(line) && line[i+1] == ' ' {
			return i + 2
		}
		return 0
	}
	return 0
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
