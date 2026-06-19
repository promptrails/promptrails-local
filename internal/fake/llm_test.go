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

func TestGeneratePromptRunResponse_Base(t *testing.T) {
	resp := GeneratePromptRunResponse("greeting", nil)

	if resp.Model != "gpt-4o" {
		t.Errorf("Model = %q, want gpt-4o", resp.Model)
	}
	if resp.Cost <= 0 {
		t.Errorf("Cost = %f, want > 0", resp.Cost)
	}
	if len(resp.TraceID) != 32 {
		t.Errorf("TraceID len = %d, want 32 hex chars", len(resp.TraceID))
	}
	assertUsageConsistent(t, resp.TokenUsage)

	// Base response has no caching or reasoning breakdown.
	if _, ok := resp.TokenUsage["cached_tokens"]; ok {
		t.Error("base response should not have cached_tokens")
	}
	if _, ok := resp.TokenUsage["reasoning_tokens"]; ok {
		t.Error("base response should not have reasoning_tokens")
	}
}

func TestGeneratePromptRunResponse_PromptCaching(t *testing.T) {
	resp := GeneratePromptRunResponse("p", &model.RunPromptRequest{PromptCaching: true})
	if _, ok := resp.TokenUsage["cached_tokens"]; !ok {
		t.Error("expected cached_tokens when PromptCaching is on")
	}
	if _, ok := resp.TokenUsage["cache_creation_tokens"]; !ok {
		t.Error("expected cache_creation_tokens when PromptCaching is on")
	}
	assertUsageConsistent(t, resp.TokenUsage)
}

func TestGeneratePromptRunResponse_Reasoning(t *testing.T) {
	resp := GeneratePromptRunResponse("p", &model.RunPromptRequest{ReasoningEffort: "high"})
	reasoning, ok := resp.TokenUsage["reasoning_tokens"]
	if !ok || reasoning <= 0 {
		t.Fatalf("expected positive reasoning_tokens, got %v (present=%v)", reasoning, ok)
	}
	// Reasoning tokens are folded into completion_tokens.
	if resp.TokenUsage["completion_tokens"] < reasoning {
		t.Errorf("completion_tokens (%d) should include reasoning (%d)",
			resp.TokenUsage["completion_tokens"], reasoning)
	}
	assertUsageConsistent(t, resp.TokenUsage)
}

// total_tokens must always equal prompt_tokens + completion_tokens.
func assertUsageConsistent(t *testing.T, u map[string]int) {
	t.Helper()
	if u["total_tokens"] != u["prompt_tokens"]+u["completion_tokens"] {
		t.Errorf("total_tokens (%d) != prompt (%d) + completion (%d)",
			u["total_tokens"], u["prompt_tokens"], u["completion_tokens"])
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
