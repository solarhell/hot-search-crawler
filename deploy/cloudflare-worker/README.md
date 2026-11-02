# Cloudflare Worker 定时器

用 Cloudflare Workers Cron Triggers 驱动 GitHub Actions：

| cron（UTC） | 触发的工作流 | 作用 |
| --- | --- | --- |
| `*/10 * * * *` | `live.yml` | 每 10 分钟抓取并发布 `latest.json` 到 GitHub Pages |
| `23 * * * *` | `crawl.yml` | 每小时 :23 归档小时快照、日报和 README |

GitHub 自带的 `schedule` 不可靠（2026 年实测常见延迟 1.5–4.5 小时，低活跃仓库每天只跑 4–8 次）。Worker 只负责按时调用 `workflow_dispatch`，真正的抓取仍然在 GitHub Runner 上执行。

Worker 只导出了 `scheduled` handler，没有 `fetch` handler，因此不会暴露公开 HTTP 入口；GitHub PAT 存在 Worker secret 里，创建后不可读。

## 前置条件

- Cloudflare 账号（免费版足够：2 个 cron trigger、约 150 次调用/天，远低于免费额度）
- Node.js 18+ 和 npm
- GitHub **fine-grained PAT**：
  - Repository access：Only select repositories → `solarhell/hot-search-crawler`
  - Permissions：`Actions` → Read and write；`Metadata` → Read-only
  - 有效期建议 90 天

## 部署

```bash
cd deploy/cloudflare-worker
npm install
npx wrangler login
npx wrangler secret put GITHUB_TOKEN
npm run deploy
```

`wrangler login` 会打开浏览器完成 Cloudflare 授权。`GITHUB_TOKEN` 直接粘贴到命令行提示里，不会写进仓库或 shell 历史。

如果这台机器没有浏览器，或不想用交互式登录，可以在 Cloudflare Dashboard 创建一个 API token（模板 **Edit Cloudflare Workers**），然后用环境变量部署：

```bash
export CLOUDFLARE_API_TOKEN=...   # 只放当前 shell，不要提交
npx wrangler whoami               # 确认账号
npx wrangler secret put GITHUB_TOKEN
npm run deploy
```

如果账号下有多个 Cloudflare account，再设置 `CLOUDFLARE_ACCOUNT_ID`。

可选：配置失败告警 webhook（Slack 兼容，接收 `{"text": "..."}`）：

```bash
npx wrangler secret put ALERT_WEBHOOK_URL
```

部署成功后 Wrangler 会输出 Worker 名称和一个 `*.workers.dev` 地址。这个地址只用于标识 Worker，项目没有注册 HTTP 路由，访问它不会触发任何抓取。

Cron 变更最多需要 15 分钟同步到 Cloudflare 全球网络；第一次部署后稍等一会儿再观察。

## 验证

看 Worker 日志：

```bash
npx wrangler tail
```

看 GitHub Actions 是否被触发：

```bash
gh run list -R solarhell/hot-search-crawler --workflow=live.yml --limit 5
gh run list -R solarhell/hot-search-crawler --workflow=crawl.yml --limit 5
```

看实时数据：

```bash
curl -s https://solarhell.github.io/hot-search-crawler/latest.json | jq '.generated_at'
```

Cloudflare 控制台 → Workers & Pages → `hot-search-crawler-scheduler` → Triggers / Past Cron Events 也能看到每次执行状态。

## 行为说明

- 每次 cron 触发都会向 GitHub 发送一次 `workflow_dispatch`，body 为 `{"ref":"main"}`。
- GitHub 返回 4xx（401/403/404/422 等，429 除外）时调用 `controller.noRetry()`：token 过期、权限不足这类错误重试也不会成功，避免无效重试。
- 5xx 和 429 会抛错，交给 Cloudflare 的重试策略处理。
- 配置了 `ALERT_WEBHOOK_URL` 时，失败会向该地址 POST 一条 Slack 兼容消息；重试场景下可能收到多条。
- `[observability]` 已开启，执行日志可在 Cloudflare 控制台和 `wrangler tail` 查看。

## 本地验证

```bash
npm test
npx wrangler deploy --dry-run
```

`npm test` 用 Node 内置 test runner 覆盖 cron 映射、鉴权头、请求体、失败告警和 4xx/5xx 的重试区分。

## 轮换 token

PAT 到期或泄露时：

```bash
cd deploy/cloudflare-worker
npx wrangler secret put GITHUB_TOKEN
```

更新 secret 会自动生成新的 Worker 版本，不需要改代码或重新配置 cron。

## 安全注意

- 不要把 PAT 写进 `wrangler.toml`、`[vars]`、代码或 git。
- PAT 只授权这个仓库的 `Actions: Read and write`，不要用账号级经典 token。
- Worker 没有公开入口，只有 Cloudflare 的 cron 能触发它。
- 如果使用告警 webhook，也不要把它提交到仓库；用 `wrangler secret put` 保存。
