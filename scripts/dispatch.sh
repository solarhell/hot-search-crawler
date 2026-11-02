#!/usr/bin/env bash
# 手动触发 workflow_dispatch，供外部定时器或本地测试使用。
#
#   GITHUB_TOKEN=... ./scripts/dispatch.sh live.yml
#   GITHUB_TOKEN=... ./scripts/dispatch.sh crawl.yml
#
# GITHUB_TOKEN 需要是 repo 级 fine-grained PAT，权限 Actions: Read and write。
# 不要把 token 写进仓库或命令行历史。
set -euo pipefail

repo="${REPO:-solarhell/hot-search-crawler}"
workflow="${1:-live.yml}"
ref="${REF:-main}"

: "${GITHUB_TOKEN:?请先设置 GITHUB_TOKEN（fine-grained PAT，需要 Actions: Read and write）}"

curl -fsS -X POST \
  -H "Authorization: Bearer ${GITHUB_TOKEN}" \
  -H "Accept: application/vnd.github+json" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  -H "User-Agent: hot-search-crawler-scheduler" \
  -H "Content-Type: application/json" \
  "https://api.github.com/repos/${repo}/actions/workflows/${workflow}/dispatches" \
  -d "{\"ref\":\"${ref}\"}"

echo "triggered ${workflow} on ${repo}@${ref}"
