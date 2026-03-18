"""Integration tests for promptrails-local using the Python SDK."""

import os

import pytest
from promptrails import PromptRails

BASE_URL = os.environ.get("PROMPTRAILS_LOCAL_URL", "http://localhost:8080")

SEED_AGENT_ID = "39wNZZu78VawB207IOPonkoP38J"
SEED_CHAT_AGENT_ID = "3A1tXOt9iovkA7LEusDSjcKbJQM"


@pytest.fixture
def client():
    return PromptRails(api_key="test-key", base_url=BASE_URL)


# --- Agents ---


def test_list_agents(client):
    result = client.agents.list()
    assert len(result.data) > 0, "expected seed agents"


def test_get_agent(client):
    agent = client.agents.get(SEED_AGENT_ID)
    assert agent.name == "Simple Agent"
    assert agent.type == "simple"


def test_create_and_delete_agent(client):
    agent = client.agents.create(
        name="Python Test Agent",
        type="simple",
        description="Created by Python integration test",
    )
    assert agent.id
    assert agent.name == "Python Test Agent"

    client.agents.delete(agent.id)


def test_execute_agent(client):
    result = client.agents.execute(
        SEED_AGENT_ID,
        input={"topic": "integration testing"},
    )
    assert result.status == "completed"
    assert result.cost >= 0


def test_list_agent_versions(client):
    versions = client.agents.list_versions(SEED_AGENT_ID)
    assert len(versions) > 0, "expected seed agent versions"


# --- Prompts ---


def test_list_prompts(client):
    result = client.prompts.list()
    assert len(result.data) > 0, "expected seed prompts"


def test_create_and_delete_prompt(client):
    prompt = client.prompts.create(
        name="Python Test Prompt",
        description="Created by Python integration test",
    )
    assert prompt.id
    assert prompt.name == "Python Test Prompt"

    client.prompts.delete(prompt.id)


# --- Executions ---


def test_list_executions(client):
    client.agents.execute(SEED_AGENT_ID, input={"topic": "test"})

    result = client.executions.list()
    assert len(result.data) > 0, "expected at least one execution"


# --- Chat ---


def test_chat_create_session(client):
    session = client.chat.create_session(
        agent_id=SEED_CHAT_AGENT_ID,
        title="Python Integration Test",
    )
    assert session.id

    sessions = client.chat.list_sessions()
    assert len(sessions.data) > 0

    client.chat.delete_session(session.id)


# --- Credentials ---


def test_list_credentials(client):
    result = client.credentials.list()
    assert len(result.data) > 0, "expected seed credentials"


# --- Data Sources ---


def test_list_data_sources(client):
    result = client.data_sources.list()
    assert len(result.data) > 0, "expected seed data sources"


# --- Traces ---


def test_list_traces(client):
    client.agents.execute(SEED_AGENT_ID, input={"topic": "traces"})

    result = client.traces.list()
    assert len(result.data) > 0, "expected traces from execution"


# --- Scores ---


def test_create_and_list_scores(client):
    exec_result = client.agents.execute(SEED_AGENT_ID, input={"topic": "scores"})

    score = client.scores.create(
        trace_id=exec_result.trace_id,
        name="accuracy",
        value=0.95,
        data_type="numeric",
        source="api",
    )
    assert score.id

    scores = client.scores.list()
    assert len(scores.data) > 0


# --- MCP Tools ---


def test_list_mcp_tools(client):
    result = client.mcp_tools.list()
    assert len(result.data) > 0, "expected seed MCP tools"


# --- Memories ---


def test_memory_list(client):
    memories = client.agents.list_memories(SEED_CHAT_AGENT_ID)
    assert len(memories.data) > 0, "expected seed memories"
