package fake

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/promptrails/promptrails-local/internal/model"
)

func TestGenerateExecutionOutput(t *testing.T) {
	out := GenerateExecutionOutput("summarizer", map[string]any{"text": "hi"})
	if out["type"] != "text" {
		t.Errorf("type = %v, want text", out["type"])
	}
	content, _ := out["content"].(string)
	if !strings.Contains(content, "summarizer") {
		t.Errorf("content should mention the agent name: %q", content)
	}
	if !strings.Contains(content, "text") {
		t.Errorf("content should list the input key: %q", content)
	}
}

func TestGeneratePromptPreview(t *testing.T) {
	prompt := model.Prompt{
		Name: "greeting",
		CurrentVersion: &model.PromptVersion{
			SystemPrompt: "You are a helpful assistant.",
			UserPrompt:   "Say hello to {{ name }}.",
		},
	}
	input := map[string]any{"name": "Ada"}

	resp := GeneratePromptPreview(prompt, input)

	// Preview is a dry-run render — the content-only templates come back verbatim.
	if resp.SystemPrompt != "You are a helpful assistant." {
		t.Errorf("SystemPrompt = %q", resp.SystemPrompt)
	}
	if resp.UserPrompt != "Say hello to {{ name }}." {
		t.Errorf("UserPrompt = %q", resp.UserPrompt)
	}
	if resp.Input["name"] != "Ada" {
		t.Errorf("Input not echoed: %v", resp.Input)
	}
}

func TestGeneratePromptPreview_NoCurrentVersion(t *testing.T) {
	resp := GeneratePromptPreview(model.Prompt{Name: "empty"}, nil)
	if resp.SystemPrompt != "" || resp.UserPrompt != "" {
		t.Errorf("expected empty render with no current version, got %+v", resp)
	}
}

func TestGenerateChatResponse_Truncates(t *testing.T) {
	long := strings.Repeat("x", 200)
	resp := GenerateChatResponse("agent-1", long)
	if !strings.Contains(resp, "agent-1") {
		t.Errorf("reply should mention agent id: %q", resp)
	}
	if !strings.Contains(resp, "...") {
		t.Errorf("long user content should be truncated with ellipsis: %q", resp)
	}
}

func TestTraceAndSpanIDs(t *testing.T) {
	tid := GenerateTraceID()
	if len(tid) != 32 {
		t.Errorf("trace id len = %d, want 32", len(tid))
	}
	if _, err := hex.DecodeString(tid); err != nil {
		t.Errorf("trace id not valid hex: %v", err)
	}
	sid := GenerateSpanID()
	if len(sid) != 16 {
		t.Errorf("span id len = %d, want 16", len(sid))
	}
	if _, err := hex.DecodeString(sid); err != nil {
		t.Errorf("span id not valid hex: %v", err)
	}
	if tid == GenerateTraceID() {
		t.Error("trace ids should be unique across calls")
	}
}

func TestCreateExecutionTrace(t *testing.T) {
	dur := int64(500)
	agentID := "agent-1"
	exec := model.Execution{
		ID:          "exec-1",
		WorkspaceID: "ws-1",
		TraceID:     "trace-abc",
		AgentID:     &agentID,
		DurationMS:  &dur,
	}

	traces := CreateExecutionTrace(exec, "my-agent")
	if len(traces) != 2 {
		t.Fatalf("expected root + child trace, got %d", len(traces))
	}
	root, child := traces[0], traces[1]

	if root.Kind != "agent" || child.Kind != "llm" {
		t.Errorf("kinds = (%q, %q), want (agent, llm)", root.Kind, child.Kind)
	}
	if child.ParentSpanID != root.SpanID {
		t.Errorf("child ParentSpanID = %q, want root SpanID %q", child.ParentSpanID, root.SpanID)
	}
	if root.TraceID != "trace-abc" || child.TraceID != "trace-abc" {
		t.Errorf("trace ids = (%q, %q), want both trace-abc", root.TraceID, child.TraceID)
	}
	if root.SpanID == child.SpanID {
		t.Error("root and child span ids must differ")
	}
	// Child duration is shortened relative to the root.
	if child.DurationMS == nil || root.DurationMS == nil || *child.DurationMS > *root.DurationMS {
		t.Errorf("child duration (%v) should be <= root duration (%v)", child.DurationMS, root.DurationMS)
	}
}
