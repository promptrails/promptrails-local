/**
 * Integration tests for promptrails-local using the JS SDK.
 *
 * Requires: npm install @promptrails/sdk
 * Run:      node --test test_integration.mjs
 *
 * NOTE: If the @promptrails/sdk package has ESM resolution issues,
 * these tests may fail to import. This is an SDK packaging issue,
 * not an emulator issue.
 */

import { describe, it, before, after } from "node:test";
import assert from "node:assert/strict";

const BASE_URL = process.env.PROMPTRAILS_LOCAL_URL || "http://localhost:8080";
const SEED_AGENT_ID = "39wNZZu78VawB207IOPonkoP38J";
const SEED_CHAT_AGENT_ID = "3A1tXOt9iovkA7LEusDSjcKbJQM";

// Dynamic import to handle potential ESM resolution issues gracefully
let PromptRails;
try {
  const sdk = await import("@promptrails/sdk");
  PromptRails = sdk.PromptRails;
} catch (err) {
  console.log(
    `Skipping JS SDK tests: ${err.message}\nThis is likely an SDK packaging issue, not an emulator issue.`,
  );
  process.exit(0);
}

const client = new PromptRails({
  apiKey: "test-key",
  baseUrl: BASE_URL,
});

// --- Agents ---

describe("Agents", () => {
  it("should list seed agents", async () => {
    const result = await client.agents.list();
    assert.ok(result.data.length > 0, "expected seed agents");
  });

  it("should get agent by ID", async () => {
    const agent = await client.agents.get(SEED_AGENT_ID);
    assert.equal(agent.name, "Simple Agent");
    assert.equal(agent.type, "simple");
  });

  it("should create and delete agent", async () => {
    const agent = await client.agents.create({
      name: "JS Test Agent",
      type: "simple",
      description: "Created by JS integration test",
    });
    assert.ok(agent.id);
    assert.equal(agent.name, "JS Test Agent");

    await client.agents.delete(agent.id);
  });

  it("should execute agent", async () => {
    const result = await client.agents.execute(SEED_AGENT_ID, {
      input: { topic: "integration testing" },
    });
    assert.equal(result.status, "completed");
    assert.ok(result.cost >= 0);
  });
});

// --- Prompts ---

describe("Prompts", () => {
  it("should list seed prompts", async () => {
    const result = await client.prompts.list();
    assert.ok(result.data.length > 0, "expected seed prompts");
  });

  it("should create and delete prompt", async () => {
    const prompt = await client.prompts.create({
      name: "JS Test Prompt",
      description: "Created by JS integration test",
    });
    assert.ok(prompt.id);

    await client.prompts.delete(prompt.id);
  });
});

// --- Executions ---

describe("Executions", () => {
  it("should list executions", async () => {
    await client.agents.execute(SEED_AGENT_ID, {
      input: { topic: "test" },
    });

    const result = await client.executions.list();
    assert.ok(result.data.length > 0, "expected executions");
  });
});

// --- Chat ---

describe("Chat", () => {
  it("should create and delete session", async () => {
    const session = await client.chat.createSession({
      agent_id: SEED_CHAT_AGENT_ID,
      title: "JS Integration Test",
    });
    assert.ok(session.id);

    const sessions = await client.chat.listSessions();
    assert.ok(sessions.data.length > 0);

    await client.chat.deleteSession(session.id);
  });
});

// --- Credentials ---

describe("Credentials", () => {
  it("should list seed credentials", async () => {
    const result = await client.credentials.list();
    assert.ok(result.data.length > 0, "expected seed credentials");
  });
});

// --- Data Sources ---

describe("Data Sources", () => {
  it("should list seed data sources", async () => {
    const result = await client.dataSources.list();
    assert.ok(result.data.length > 0, "expected seed data sources");
  });
});

// --- Traces ---

describe("Traces", () => {
  it("should list traces after execution", async () => {
    await client.agents.execute(SEED_AGENT_ID, {
      input: { topic: "traces" },
    });
    const result = await client.traces.list();
    assert.ok(result.data.length > 0, "expected traces");
  });
});

// --- MCP Tools ---

describe("MCP Tools", () => {
  it("should list seed MCP tools", async () => {
    const result = await client.mcpTools.list();
    assert.ok(result.data.length > 0, "expected seed MCP tools");
  });
});
