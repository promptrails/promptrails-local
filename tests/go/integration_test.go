package integration

import (
	"context"
	"os"
	"testing"

	promptrails "github.com/promptrails/go-sdk"
)

func baseURL() string {
	if v := os.Getenv("PROMPTRAILS_LOCAL_URL"); v != "" {
		return v
	}
	return "http://localhost:8080"
}

func newClient() *promptrails.Client {
	return promptrails.NewClient("test-key", promptrails.WithBaseURL(baseURL()))
}

// --- Agents ---

func TestListAgents(t *testing.T) {
	client := newClient()
	ctx := context.Background()

	result, err := client.Agents.List(ctx, nil)
	if err != nil {
		t.Fatalf("Agents.List: %v", err)
	}
	if len(result.Data) == 0 {
		t.Fatal("expected seed agents, got 0")
	}
	t.Logf("found %d agents", len(result.Data))
}

func TestGetAgent(t *testing.T) {
	client := newClient()
	ctx := context.Background()

	agent, err := client.Agents.Get(ctx, "39wNZZu78VawB207IOPonkoP38J")
	if err != nil {
		t.Fatalf("Agents.Get: %v", err)
	}
	if agent.Name != "Simple Agent" {
		t.Fatalf("expected 'Simple Agent', got %q", agent.Name)
	}
}

func TestCreateAndDeleteAgent(t *testing.T) {
	client := newClient()
	ctx := context.Background()

	agent, err := client.Agents.Create(ctx, &promptrails.CreateAgentParams{
		Name:        "Integration Test Agent",
		Type:        "simple",
		Description: "Created by Go integration test",
	})
	if err != nil {
		t.Fatalf("Agents.Create: %v", err)
	}
	if agent.ID == "" {
		t.Fatal("expected agent ID")
	}
	if agent.Name != "Integration Test Agent" {
		t.Fatalf("expected 'Integration Test Agent', got %q", agent.Name)
	}

	err = client.Agents.Delete(ctx, agent.ID)
	if err != nil {
		t.Fatalf("Agents.Delete: %v", err)
	}
}

func TestExecuteAgent(t *testing.T) {
	client := newClient()
	ctx := context.Background()

	result, err := client.Agents.Execute(ctx, "39wNZZu78VawB207IOPonkoP38J", &promptrails.ExecuteAgentParams{
		Input: map[string]any{"topic": "integration testing"},
	})
	if err != nil {
		t.Fatalf("Agents.Execute: %v", err)
	}
	if result.Status != "completed" {
		t.Fatalf("expected status 'completed', got %q", result.Status)
	}
	t.Logf("execution completed, cost=%.6f", result.Cost)
}

// --- Prompts ---

func TestListPrompts(t *testing.T) {
	client := newClient()
	ctx := context.Background()

	result, err := client.Prompts.List(ctx, nil)
	if err != nil {
		t.Fatalf("Prompts.List: %v", err)
	}
	if len(result.Data) == 0 {
		t.Fatal("expected seed prompts, got 0")
	}
	t.Logf("found %d prompts", len(result.Data))
}

func TestCreateAndDeletePrompt(t *testing.T) {
	client := newClient()
	ctx := context.Background()

	prompt, err := client.Prompts.Create(ctx, &promptrails.CreatePromptParams{
		Name:        "Test Prompt",
		Description: "Created by Go integration test",
	})
	if err != nil {
		t.Fatalf("Prompts.Create: %v", err)
	}
	if prompt.ID == "" {
		t.Fatal("expected prompt ID")
	}

	err = client.Prompts.Delete(ctx, prompt.ID)
	if err != nil {
		t.Fatalf("Prompts.Delete: %v", err)
	}
}

// --- Executions ---

func TestListExecutions(t *testing.T) {
	client := newClient()
	ctx := context.Background()

	_, err := client.Agents.Execute(ctx, "39wNZZu78VawB207IOPonkoP38J", &promptrails.ExecuteAgentParams{
		Input: map[string]any{"topic": "test"},
	})
	if err != nil {
		t.Fatalf("Agents.Execute: %v", err)
	}

	result, err := client.Executions.List(ctx, nil)
	if err != nil {
		t.Fatalf("Executions.List: %v", err)
	}
	if len(result.Data) == 0 {
		t.Fatal("expected at least one execution")
	}
	t.Logf("found %d executions", len(result.Data))
}

// --- Chat ---

func TestChatCreateSession(t *testing.T) {
	client := newClient()
	ctx := context.Background()

	session, err := client.Chat.CreateSession(ctx, &promptrails.CreateSessionParams{
		AgentID: "3A1tXOt9iovkA7LEusDSjcKbJQM",
		Title:   "Go Integration Test",
	})
	if err != nil {
		t.Fatalf("Chat.CreateSession: %v", err)
	}
	if session.ID == "" {
		t.Fatal("expected session ID")
	}

	// List sessions
	sessions, err := client.Chat.ListSessions(ctx, nil)
	if err != nil {
		t.Fatalf("Chat.ListSessions: %v", err)
	}
	if len(sessions.Data) == 0 {
		t.Fatal("expected at least one session")
	}

	err = client.Chat.DeleteSession(ctx, session.ID)
	if err != nil {
		t.Fatalf("Chat.DeleteSession: %v", err)
	}
}

// --- Credentials ---

func TestListCredentials(t *testing.T) {
	client := newClient()
	ctx := context.Background()

	result, err := client.Credentials.List(ctx, nil)
	if err != nil {
		t.Fatalf("Credentials.List: %v", err)
	}
	if len(result.Data) == 0 {
		t.Fatal("expected seed credentials, got 0")
	}
	t.Logf("found %d credentials", len(result.Data))
}

// --- Traces ---

func TestListTraces(t *testing.T) {
	client := newClient()
	ctx := context.Background()

	_, _ = client.Agents.Execute(ctx, "39wNZZu78VawB207IOPonkoP38J", &promptrails.ExecuteAgentParams{
		Input: map[string]any{"topic": "traces"},
	})

	result, err := client.Traces.List(ctx, nil)
	if err != nil {
		t.Fatalf("Traces.List: %v", err)
	}
	if len(result.Data) == 0 {
		t.Fatal("expected traces from execution")
	}
	t.Logf("found %d traces", len(result.Data))
}

// --- Data Sources ---

func TestListDataSources(t *testing.T) {
	client := newClient()
	ctx := context.Background()

	result, err := client.DataSources.List(ctx, nil)
	if err != nil {
		t.Fatalf("DataSources.List: %v", err)
	}
	if len(result.Data) == 0 {
		t.Fatal("expected seed data sources, got 0")
	}
	t.Logf("found %d data sources", len(result.Data))
}

// --- Scores ---

func TestListScores(t *testing.T) {
	client := newClient()
	ctx := context.Background()

	// Just verify the endpoint works
	_, err := client.Scores.List(ctx, nil)
	if err != nil {
		t.Fatalf("Scores.List: %v", err)
	}
}
