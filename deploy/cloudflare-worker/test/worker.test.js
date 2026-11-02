import assert from "node:assert/strict";
import test from "node:test";

import worker, {
  ARCHIVE_CRON,
  LIVE_CRON,
  dispatchWorkflow,
  workflowForCron,
} from "../src/index.js";

test("workflowForCron routes live and archive", () => {
  assert.equal(workflowForCron(LIVE_CRON), "live.yml");
  assert.equal(workflowForCron(ARCHIVE_CRON), "crawl.yml");
});

test("dispatchWorkflow posts to the workflow dispatch endpoint", async () => {
  const calls = [];
  const fakeFetch = async (url, init) => {
    calls.push({ url, init });
    return new Response(null, { status: 204 });
  };

  await dispatchWorkflow("live.yml", { GITHUB_TOKEN: "test-token" }, fakeFetch);

  assert.equal(calls.length, 1);
  assert.equal(
    calls[0].url,
    "https://api.github.com/repos/solarhell/hot-search-crawler/actions/workflows/live.yml/dispatches",
  );
  assert.equal(calls[0].init.method, "POST");
  assert.equal(calls[0].init.headers.Authorization, "Bearer test-token");
  assert.equal(calls[0].init.headers.Accept, "application/vnd.github+json");
  assert.deepEqual(JSON.parse(calls[0].init.body), { ref: "main" });
});

test("dispatchWorkflow sends an alert and throws on failure", async () => {
  const calls = [];
  const fakeFetch = async (url, init) => {
    calls.push({ url, init });
    if (url.includes("dispatches")) {
      return new Response("bad credentials", { status: 401 });
    }
    return new Response(null, { status: 200 });
  };

  await assert.rejects(
    dispatchWorkflow(
      "crawl.yml",
      { GITHUB_TOKEN: "expired", ALERT_WEBHOOK_URL: "https://hooks.example/alert" },
      fakeFetch,
    ),
    /HTTP 401/,
  );

  assert.equal(calls.length, 2);
  const alert = JSON.parse(calls[1].init.body);
  assert.match(alert.text, /crawl\.yml/);
  assert.match(alert.text, /HTTP 401/);
});

test("scheduled handler disables retry for permanent failures", async () => {
  const originalFetch = globalThis.fetch;
  globalThis.fetch = async () => new Response("forbidden", { status: 403 });

  let noRetry = false;
  const controller = {
    cron: LIVE_CRON,
    noRetry() {
      noRetry = true;
    },
  };

  try {
    await assert.rejects(worker.scheduled(controller, { GITHUB_TOKEN: "t" }), /HTTP 403/);
    assert.equal(noRetry, true);
  } finally {
    globalThis.fetch = originalFetch;
  }
});

test("scheduled handler keeps retries for 5xx", async () => {
  const originalFetch = globalThis.fetch;
  globalThis.fetch = async () => new Response("server error", { status: 500 });

  let noRetry = false;
  const controller = {
    cron: ARCHIVE_CRON,
    noRetry() {
      noRetry = true;
    },
  };

  try {
    await assert.rejects(worker.scheduled(controller, { GITHUB_TOKEN: "t" }), /HTTP 500/);
    assert.equal(noRetry, false);
  } finally {
    globalThis.fetch = originalFetch;
  }
});
