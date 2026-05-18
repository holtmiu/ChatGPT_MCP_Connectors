package feishu

import "testing"

func TestBuildAppendBlocksRequest(t *testing.T) {
	body, ids, err := buildAppendBlocksRequest(AppendRequest{Markdown: "# Title\n\nBody"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	children, ok := body["children"].([]any)
	if !ok || len(children) != 2 || len(ids) != 2 {
		t.Fatalf("unexpected body=%+v ids=%+v", body, ids)
	}
}

func TestBuildCreateDocumentRequest(t *testing.T) {
	body, err := buildCreateDocumentRequest(CreateDocumentRequest{Title: "Spec", FolderToken: "folder"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if body["title"] != "Spec" || body["folder_token"] != "folder" {
		t.Fatalf("unexpected body: %+v", body)
	}
}
