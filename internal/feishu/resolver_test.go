package feishu

import "testing"

func TestResolveDocxURL(t *testing.T) {
	r := NewResolver("feishu")
	got, err := r.Resolve("https://example.feishu.cn/docx/AbCdEf123456?from=from_copylink")
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	if got.Provider != ProviderFeishu || got.ResourceType != ResourceDocx || got.Token != "AbCdEf123456" {
		t.Fatalf("unexpected identity: %+v", got)
	}
}

func TestResolveWikiURL(t *testing.T) {
	r := NewResolver("lark")
	got, err := r.Resolve("https://example.larksuite.com/wiki/WikiToken123")
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	if got.Provider != ProviderLark || got.ResourceType != ResourceWiki || got.Token != "WikiToken123" {
		t.Fatalf("unexpected identity: %+v", got)
	}
}

func TestResolveBareToken(t *testing.T) {
	r := NewResolver("feishu")
	got, err := r.Resolve("DocToken123")
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	if got.ResourceType != ResourceDocx || got.Token != "DocToken123" {
		t.Fatalf("unexpected identity: %+v", got)
	}
}
