# 外部定时器

GitHub Actions 自带的 `schedule` 是 **best-effort**：2026 年实测常见延迟 1.5–4.5 小时，低活跃仓库的每小时任务实际可能每天只跑 4–8 次，丢跑不会留下任何记录。因此本项目的定时由外部调度器通过 `workflow_dispatch` 触发。

## 架构

| 工作流 | 触发方式 | 频率 | 产物 |
| --- | --- | --- | --- |
| `live.yml` | 外部定时器 → `workflow_dispatch`，GitHub `schedule` 在 `:13/:43` 兜底 | 每 10 分钟 | 只发布 `latest.json` 到 GitHub Pages，不写 git |
| `crawl.yml` | 外部定时器每小时 `:23`，GitHub `schedule` 在 `:47` 兜底 | 每小时 | 写入小时快照、日报、README |

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

生成后只复制一次，直接粘贴到外部定时器的鉴权配置里。

## 2. cron-job.org

[cron-job.org](https://cron-job.org) 免费版支持自定义请求头和失败邮件通知，适合这个场景。

### Job A：实时数据（每 10 分钟）

- Title：`hot-search-crawler live`
- URL：

  ```text
  https://api.github.com/repos/solarhell/hot-search-crawler/actions/workflows/live.yml/dispatches
  ```

- Schedule：Every 10 minutes
- Request method：`POST`
- Headers：

  ```text
  Authorization: Bearer <你的 fine-grained PAT>
  Accept: application/vnd.github+json
  X-GitHub-Api-Version: 2022-11-28
  User-Agent: hot-search-crawler-scheduler
  Content-Type: application/json
  ```

- Request body：

  ```json
  {"ref":"main"}
  ```

- 建议开启失败通知（cron-job.org 会给账号邮箱发信）。

### Job B：小时归档（每小时 :23）

- Title：`hot-search-crawler archive`
- URL：

  ```text
  https://api.github.com/repos/solarhell/hot-search-crawler/actions/workflows/crawl.yml/dispatches
  ```

- Schedule：`23 * * * *`（cron-job.org 使用 UTC；北京时间是整点偏移，所以本地也是每小时的 :23）
- Headers：与 Job A 相同
- Request body：

  ```json
  {"ref":"main"}
  ```

GitHub 自己的 `schedule` 在 `:47` 作为兜底；两边谁先跑谁归档，后跑的一方会跳过。

## 3. 验证

先手动触发一次，确认 token 和工作流都正常：

```bash
export GITHUB_TOKEN='<fine-grained PAT>'   # 只在当前 shell 使用，不要写进文件

./scripts/dispatch.sh live.yml
./scripts/dispatch.sh crawl.yml

gh run list --workflow=live.yml --limit 3
gh run list --workflow=crawl.yml --limit 3
```

等 live 工作流完成后检查 Pages：

```bash
curl -s https://solarhell.github.io/hot-search-crawler/latest.json | jq '{generated_at, boards: [.boards[] | {platform, ok, n: (.items|length)}]}'
```

归档工作流完成后，仓库会多一个 `chore(data): snapshot ...` 提交。

## 4. 替代方案

### 自己的服务器 / systemd timer

如果不想把 token 放在第三方服务，可以在自己的常在线服务器上跑一个 timer。创建 `/etc/hot-search-crawler.env`（权限 `600`）：

```bash
GITHUB_TOKEN=github_pat_xxx
```

创建 `/etc/systemd/system/hot-search-crawler-live.service`：

```ini
[Unit]
Description=Trigger hot-search-crawler live workflow

[Service]
Type=oneshot
EnvironmentFile=/etc/hot-search-crawler.env
ExecStart=/opt/hot-search-crawler/scripts/dispatch.sh live.yml
```

创建 `/etc/systemd/system/hot-search-crawler-live.timer`：

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

### Cloudflare Worker Cron

如果已经有 Cloudflare 账号，也可以用一个 Worker 定时调用同一个 GitHub API，把 PAT 放在 Worker Secret 里：

```js
export default {
  async scheduled(_event, env) {
    const res = await fetch(
      "https://api.github.com/repos/solarhell/hot-search-crawler/actions/workflows/live.yml/dispatches",
      {
        method: "POST",
        headers: {
          Authorization: `Bearer ${env.GITHUB_TOKEN}`,
          Accept: "application/vnd.github+json",
          "X-GitHub-Api-Version": "2022-11-28",
          "User-Agent": "hot-search-crawler-scheduler",
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ ref: "main" }),
      },
    );
    if (!res.ok) {
      throw new Error(`dispatch failed: ${res.status} ${await res.text()}`);
    }
  },
};
```

`wrangler.toml` 里配置 `[triggers] crons = ["*/10 * * * *"]`，再用 `wrangler secret put GITHUB_TOKEN` 写入 token。

## 5. 注意事项

- **Pages 部署频率**：GitHub Pages 通过 Actions 部署有软限制（约 10 次/小时）。10 分钟一次是 6 次/小时，加上 GitHub 兜底最多 2 次/小时，仍低于限制；不要改成 5 分钟。
- **token 最小权限**：只给这个仓库的 `Actions: Read and write`。token 泄露时立即在 GitHub 设置里 Revoke。
- **token 轮换**：fine-grained PAT 到期后外部定时器会收到 401，cron-job.org 的失败通知会提醒；更新 token 即可。
- **Pages 缓存**：`latest.json` 可能被 CDN 缓存几分钟；消费端可以加 `?t=<timestamp>` 绕过缓存。
- **监控**：`gh run list --workflow=live.yml` 应保持每 10 分钟至少一条记录；如果 30 分钟没有新运行，先看 cron-job.org 的日志，再看 GitHub Actions 是否被限流。
- **公开数据**：Pages 上的 `latest.json` 是公开的，和仓库里的归档数据一致，不包含任何账号信息。
