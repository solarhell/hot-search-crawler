# 外部定时器

GitHub Actions 自带的 `schedule` 是 **best-effort**：2026 年实测常见延迟 1.5–4.5 小时，低活跃仓库的每小时任务实际每天可能只跑 4–8 次，丢跑不会留下任何记录。因此本项目的定时由外部调度器通过 `workflow_dispatch` 触发。

## 架构

| 工作流 | 触发方式 | 频率 | 产物 |
| --- | --- | --- | --- |
| `live.yml` | Cloudflare Worker → `workflow_dispatch`，GitHub `schedule` 在 `:13/:43` 兜底 | 每 10 分钟 | 只发布 `latest.json` 到 GitHub Pages，不写 git |
| `crawl.yml` | Cloudflare Worker 每小时 `:23`，GitHub `schedule` 在 `:47` 兜底 | 每小时 | 写入小时快照、日报、README |

两个工作流都带并发组：`live` 不会重叠部署；`crawl` 会检查当前小时的快照是否已存在，已存在则跳过，因此外部定时器和 GitHub 兜底同时触发也不会重复提交。

实时数据地址：

```text
https://solarhell.github.io/hot-search-crawler/latest.json
```

## 1. 创建 fine-grained PAT

外部定时器需要一个能触发 Actions 的 token。**不要使用账号的经典 token，也不要把 token 提交到仓库。**

GitHub → 右上角头像 → Settings → Developer settings → Personal access tokens → **Fine-grained tokens** → Generate new token：

- **Resource owner**：`solarhell`
- **Repository access**：Only select repositories → `solarhell/hot-search-crawler`
- **Permissions**：
  - `Actions` → **Read and write**
  - `Metadata` → Read-only（自动勾选）
- **Expiration**：建议 90 天，到期前轮换

生成后只复制一次，直接粘贴到 Cloudflare Worker secret 或 cron-job.org 的鉴权配置里。

## 2. Cloudflare Worker（推荐）

仓库内置了完整实现：[`deploy/cloudflare-worker`](../deploy/cloudflare-worker/README.md)。

Worker 只导出 `scheduled` handler，没有公开的 `fetch` 入口；PAT 保存在 Worker secret 中，创建后不可读。两个 cron：

| cron（UTC） | 触发 | 作用 |
| --- | --- | --- |
| `*/10 * * * *` | `live.yml` | 实时发布到 Pages |
| `23 * * * *` | `crawl.yml` | 小时归档 |

### 部署

```bash
cd deploy/cloudflare-worker
npm install
npx wrangler login
npx wrangler secret put GITHUB_TOKEN
npm run deploy
```

可选失败告警（Slack 兼容，接收 `{"text": "..."}`）：

```bash
npx wrangler secret put ALERT_WEBHOOK_URL
```

### 验证

```bash
npx wrangler tail
gh run list -R solarhell/hot-search-crawler --workflow=live.yml --limit 5
gh run list -R solarhell/hot-search-crawler --workflow=crawl.yml --limit 5
curl -s https://solarhell.github.io/hot-search-crawler/latest.json | jq '.generated_at'
```

Cron 变更最多需要 15 分钟同步到 Cloudflare 全球网络；Cloudflare 控制台的 Workers → Triggers / Past Cron Events 可以看到执行记录。

### 失败处理

- GitHub 返回 401/403/404/422 等永久错误时，Worker 调用 `controller.noRetry()`，避免无效重试；429 和 5xx 会抛错，交给 Cloudflare 的重试策略。
- 配置 `ALERT_WEBHOOK_URL` 后，失败会发送告警；重试场景下可能收到多条。
- Worker 的 `[observability]` 已开启，日志可在控制台和 `wrangler tail` 查看。

### 轮换 token

```bash
cd deploy/cloudflare-worker
npx wrangler secret put GITHUB_TOKEN
```

更新 secret 会生成新的 Worker 版本，cron 配置不需要改。

## 3. cron-job.org（备选）

如果不想部署 Worker，cron-job.org 免费版也能用。它支持自定义请求头/body，但 PAT 会保存在 job 配置里（登录后台可见），且没有自动重试。

### Job A：实时数据（每 10 分钟）

- URL：`https://api.github.com/repos/solarhell/hot-search-crawler/actions/workflows/live.yml/dispatches`
- Method：`POST`
- Headers：

  ```text
  Authorization: Bearer <你的 fine-grained PAT>
  Accept: application/vnd.github+json
  X-GitHub-Api-Version: 2022-11-28
  User-Agent: hot-search-crawler-scheduler
  Content-Type: application/json
  ```

- Body：`{"ref":"main"}`

### Job B：小时归档（每小时 :23）

- URL：`https://api.github.com/repos/solarhell/hot-search-crawler/actions/workflows/crawl.yml/dispatches`
- Headers 同上
- Body：`{"ref":"main"}`
- Schedule：`23 * * * *`（cron-job.org 使用 UTC；北京时间是整点偏移，所以本地也是 :23）

注意：cron-job.org 会忽略自定义 `User-Agent`，用它自己的 UA；30 秒超时；最近 50 次执行、响应体保留 2 天；连续失败 25 次会自动停用并邮件通知。

## 4. 自己的服务器（备选）

不想把 PAT 交给第三方服务时，可以用自己的服务器 + systemd timer，只执行一个 curl 脚本。创建 `/etc/hot-search-crawler.env`（权限 `600`）：

```bash
GITHUB_TOKEN=github_pat_xxx
```

`/etc/systemd/system/hot-search-crawler-live.service`：

```ini
[Unit]
Description=Trigger hot-search-crawler live workflow

[Service]
Type=oneshot
EnvironmentFile=/etc/hot-search-crawler.env
ExecStart=/opt/hot-search-crawler/scripts/dispatch.sh live.yml
```

`/etc/systemd/system/hot-search-crawler-live.timer`：

```ini
[Unit]
Description=Run hot-search-crawler live every 10 minutes

[Timer]
OnCalendar=*:0/10
Persistent=true

[Install]
WantedBy=timers.target
```

启用：

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now hot-search-crawler-live.timer
```

归档同理，再加一个 `OnCalendar=*-*-* *:23:00` 的 timer，执行 `dispatch.sh crawl.yml`。

## 5. 注意事项

- **Pages 部署频率**：GitHub Pages 通过 Actions 部署有软限制（约 10 次/小时）。10 分钟一次是 6 次/小时，加上 GitHub 兜底最多 2 次/小时，仍低于限制；不要改成 5 分钟。
- **token 最小权限**：只给这个仓库的 `Actions: Read and write`。token 泄露时立即在 GitHub 设置里 Revoke。
- **token 轮换**：fine-grained PAT 到期后外部定时器会收到 401；Cloudflare 会在日志/Past Cron Events 里显示失败，cron-job.org 会发失败邮件。
- **Pages 缓存**：`latest.json` 可能被 CDN 缓存几分钟；消费端可以加 `?t=<timestamp>` 绕过缓存。
- **监控**：`gh run list --workflow=live.yml` 应保持每 10 分钟至少一条记录；如果 30 分钟没有新运行，先看 Cloudflare 的 Cron Events / `wrangler tail`，再看 GitHub Actions 是否被限流。
- **公开数据**：Pages 上的 `latest.json` 是公开的，和仓库里的归档数据一致，不包含任何账号信息。
