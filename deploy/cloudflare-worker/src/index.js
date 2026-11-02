const REPO = "solarhell/hot-search-crawler";
const API_BASE = "https://api.github.com";
const REF = "main";

export const LIVE_CRON = "*/10 * * * *";
export const ARCHIVE_CRON = "23 * * * *";

/** 根据触发本次运行的 cron 表达式选择工作流。 */
export function workflowForCron(cron) {
  return cron === ARCHIVE_CRON ? "crawl.yml" : "live.yml";
}

/**
 * 调用 GitHub 的 workflow_dispatch 接口。
 * fetchImpl 仅用于测试注入，线上使用全局 fetch。
 */
export async function dispatchWorkflow(workflow, env, fetchImpl = fetch) {
  const response = await fetchImpl(
    `${API_BASE}/repos/${REPO}/actions/workflows/${workflow}/dispatches`,
    {
      method: "POST",
      headers: {
        Authorization: `Bearer ${env.GITHUB_TOKEN}`,
        Accept: "application/vnd.github+json",
        "X-GitHub-Api-Version": "2022-11-28",
        "User-Agent": "hot-search-crawler-scheduler",
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ ref: REF }),
    },
  );

  if (response.ok) {
    console.log(`dispatched ${workflow} on ${REPO}@${REF}`);
    return;
  }

  const body = await response.text();
  const error = new Error(
    `dispatch ${workflow} failed: HTTP ${response.status}: ${truncate(body, 500)}`,
  );
  error.status = response.status;

  if (env.ALERT_WEBHOOK_URL) {
    await notify(env.ALERT_WEBHOOK_URL, `[hot-search-crawler] ${error.message}`, fetchImpl);
  }

  throw error;
}

async function notify(url, text, fetchImpl) {
  try {
    await fetchImpl(url, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      // Slack 兼容格式；Discord/飞书等可自行改 body。
      body: JSON.stringify({ text }),
    });
  } catch (error) {
    console.error(`alert webhook failed: ${error}`);
  }
}

function truncate(value, max) {
  return value.length <= max ? value : `${value.slice(0, max)}...`;
}

export default {
  async scheduled(controller, env) {
    try {
      await dispatchWorkflow(workflowForCron(controller.cron), env);
    } catch (error) {
      // 4xx 通常重试也不会成功（token 过期、权限不足、工作流不存在）。
      // 429 是限流，保留重试；其余 4xx 直接告诉 Cloudflare 不要重试。
      if (error.status >= 400 && error.status < 500 && error.status !== 429) {
        controller.noRetry?.();
      }
      throw error;
    }
  },
};
