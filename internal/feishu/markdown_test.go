package feishu

import (
	"strings"
	"testing"
)

func TestExportMarkdown(t *testing.T) {
	got := exportMarkdown([]NormalizedBlock{
		{ID: "1", Type: "heading", Text: "Title", Attrs: map[string]any{"level": 1}},
		{ID: "2", Type: "paragraph", Text: "Body"},
	})
	if !strings.Contains(got, "# Title") || !strings.Contains(got, "Body") {
		t.Fatalf("unexpected markdown output")
	}
}

func TestMarkdownToBlocks(t *testing.T) {
	got := markdownToBlocks("# Title\n\n- Item\n\nBody")
	if len(got) != 3 {
		t.Fatalf("expected 3 blocks, got %d", len(got))
	}
	if got[0].Type != "heading" || got[1].Type != "bullet_list" || got[2].Type != "paragraph" {
		t.Fatalf("unexpected blocks: %+v", got)
	}
}
