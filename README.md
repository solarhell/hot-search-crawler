# 热搜爬虫

抓取**微博、抖音、小红书**热榜，输出统一结构的 JSON 快照。仓库自带两条发布链路：

- **实时**：用仓库内置的 Cloudflare Worker 每 10 分钟触发一次 GitHub Actions，把最新榜单发布到 GitHub Pages。
- **归档**：每小时把快照、日报和 README 提交回仓库，仓库本身就是一份可检索的热搜历史。

[![crawl](https://github.com/solarhell/hot-search-crawler/actions/workflows/crawl.yml/badge.svg)](https://github.com/solarhell/hot-search-crawler/actions/workflows/crawl.yml)
[![live](https://github.com/solarhell/hot-search-crawler/actions/workflows/live.yml/badge.svg)](https://github.com/solarhell/hot-search-crawler/actions/workflows/live.yml)

**实时 JSON**：<https://solarhell.github.io/hot-search-crawler/latest.json>

<!-- BEGIN HOTSEARCH -->
> 数据更新于 **2026-11-07 02:24**（Asia/Shanghai），展示今日榜单前 10。

### 微博

1. [后会友期综史一](https://s.weibo.com/weibo?q=%23%E5%90%8E%E4%BC%9A%E5%8F%8B%E6%9C%9F%E7%BB%BC%E5%8F%B2%E4%B8%80%23) — 95.3万 · 热
2. [皮肤一碰就破女孩父母称不生二胎](https://s.weibo.com/weibo?q=%23%E7%9A%AE%E8%82%A4%E4%B8%80%E7%A2%B0%E5%B0%B1%E7%A0%B4%E5%A5%B3%E5%AD%A9%E7%88%B6%E6%AF%8D%E7%A7%B0%E4%B8%8D%E7%94%9F%E4%BA%8C%E8%83%8E%23) — 68.6万
3. [6图看我国服务业发展成绩单](https://s.weibo.com/weibo?q=%236%E5%9B%BE%E7%9C%8B%E6%88%91%E5%9B%BD%E6%9C%8D%E5%8A%A1%E4%B8%9A%E5%8F%91%E5%B1%95%E6%88%90%E7%BB%A9%E5%8D%95%23) — 56.9万
4. [特斯拉将逐步停产Cybertruck](https://s.weibo.com/weibo?q=%23%E7%89%B9%E6%96%AF%E6%8B%89%E5%B0%86%E9%80%90%E6%AD%A5%E5%81%9C%E4%BA%A7Cybertruck%23) — 47.6万
5. [一图速览进博会亮点](https://s.weibo.com/weibo?q=%23%E4%B8%80%E5%9B%BE%E9%80%9F%E8%A7%88%E8%BF%9B%E5%8D%9A%E4%BC%9A%E4%BA%AE%E7%82%B9%23) — 47.5万
6. [许飞非虚构写作级别的表达](https://s.weibo.com/weibo?q=%E8%AE%B8%E9%A3%9E%E9%9D%9E%E8%99%9A%E6%9E%84%E5%86%99%E4%BD%9C%E7%BA%A7%E5%88%AB%E7%9A%84%E8%A1%A8%E8%BE%BE) — 47.4万 · 综艺
7. [曝我们来了阵容](https://s.weibo.com/weibo?q=%E6%9B%9D%E6%88%91%E4%BB%AC%E6%9D%A5%E4%BA%86%E9%98%B5%E5%AE%B9) — 41.4万 · 综艺
8. [小米澎程主驾车门被撞仍保护车主](https://s.weibo.com/weibo?q=%E5%B0%8F%E7%B1%B3%E6%BE%8E%E7%A8%8B%E4%B8%BB%E9%A9%BE%E8%BD%A6%E9%97%A8%E8%A2%AB%E6%92%9E%E4%BB%8D%E4%BF%9D%E6%8A%A4%E8%BD%A6%E4%B8%BB) — 39.2万 · 新
9. [赵丽颖周迅打招呼](https://s.weibo.com/weibo?q=%23%E8%B5%B5%E4%B8%BD%E9%A2%96%E5%91%A8%E8%BF%85%E6%89%93%E6%8B%9B%E5%91%BC%23) — 38.9万 · 盛典
10. [王楚钦赛后情绪低落](https://s.weibo.com/weibo?q=%E7%8E%8B%E6%A5%9A%E9%92%A6%E8%B5%9B%E5%90%8E%E6%83%85%E7%BB%AA%E4%BD%8E%E8%90%BD) — 38.9万

### 抖音

1. [国防部回应攻击-21在四川舰起降](https://www.douyin.com/search/%E5%9B%BD%E9%98%B2%E9%83%A8%E5%9B%9E%E5%BA%94%E6%94%BB%E5%87%BB-21%E5%9C%A8%E5%9B%9B%E5%B7%9D%E8%88%B0%E8%B5%B7%E9%99%8D) — 1156万
2. [我的不烧心旅行地藏不住了](https://www.douyin.com/search/%E6%88%91%E7%9A%84%E4%B8%8D%E7%83%A7%E5%BF%83%E6%97%85%E8%A1%8C%E5%9C%B0%E8%97%8F%E4%B8%8D%E4%BD%8F%E4%BA%86) — 1109.3万
3. [2026世界市长对话宁波](https://www.douyin.com/search/2026%E4%B8%96%E7%95%8C%E5%B8%82%E9%95%BF%E5%AF%B9%E8%AF%9D%E5%AE%81%E6%B3%A2) — 1107.8万
4. [立冬和饺子最配了](https://www.douyin.com/search/%E7%AB%8B%E5%86%AC%E5%92%8C%E9%A5%BA%E5%AD%90%E6%9C%80%E9%85%8D%E4%BA%86) — 1107.4万
5. [你怎么舍得让有趣的知识流走](https://www.douyin.com/search/%E4%BD%A0%E6%80%8E%E4%B9%88%E8%88%8D%E5%BE%97%E8%AE%A9%E6%9C%89%E8%B6%A3%E7%9A%84%E7%9F%A5%E8%AF%86%E6%B5%81%E8%B5%B0) — 1048.6万
6. [多地目击巨型火流星划过](https://www.douyin.com/search/%E5%A4%9A%E5%9C%B0%E7%9B%AE%E5%87%BB%E5%B7%A8%E5%9E%8B%E7%81%AB%E6%B5%81%E6%98%9F%E5%88%92%E8%BF%87) — 1046.7万
7. [又到穿费尔岛的季节了](https://www.douyin.com/search/%E5%8F%88%E5%88%B0%E7%A9%BF%E8%B4%B9%E5%B0%94%E5%B2%9B%E7%9A%84%E5%AD%A3%E8%8A%82%E4%BA%86) — 1030.3万
8. [好闺闺都要吃的入冬九件套](https://www.douyin.com/search/%E5%A5%BD%E9%97%BA%E9%97%BA%E9%83%BD%E8%A6%81%E5%90%83%E7%9A%84%E5%85%A5%E5%86%AC%E4%B9%9D%E4%BB%B6%E5%A5%97) — 927万
9. [降温了羽绒服先穿为敬](https://www.douyin.com/search/%E9%99%8D%E6%B8%A9%E4%BA%86%E7%BE%BD%E7%BB%92%E6%9C%8D%E5%85%88%E7%A9%BF%E4%B8%BA%E6%95%AC) — 916.6万
10. [土耳其连环枪击案致8人死亡](https://www.douyin.com/search/%E5%9C%9F%E8%80%B3%E5%85%B6%E8%BF%9E%E7%8E%AF%E6%9E%AA%E5%87%BB%E6%A1%88%E8%87%B48%E4%BA%BA%E6%AD%BB%E4%BA%A1) — 916.2万

### 小红书

1. [用万能旅行拍照姿势美美出片](https://www.xiaohongshu.com/search_result?keyword=%E7%94%A8%E4%B8%87%E8%83%BD%E6%97%85%E8%A1%8C%E6%8B%8D%E7%85%A7%E5%A7%BF%E5%8A%BF%E7%BE%8E%E7%BE%8E%E5%87%BA%E7%89%87) — 947.5万 · 热
2. [耗时三年拍下古诗词里的中国](https://www.xiaohongshu.com/search_result?keyword=%E8%80%97%E6%97%B6%E4%B8%89%E5%B9%B4%E6%8B%8D%E4%B8%8B%E5%8F%A4%E8%AF%97%E8%AF%8D%E9%87%8C%E7%9A%84%E4%B8%AD%E5%9B%BD) — 934.8万 · 热
3. [我拍到了海鸥雨](https://www.xiaohongshu.com/search_result?keyword=%E6%88%91%E6%8B%8D%E5%88%B0%E4%BA%86%E6%B5%B7%E9%B8%A5%E9%9B%A8) — 903.2万
4. [超日常美食教程速来get](https://www.xiaohongshu.com/search_result?keyword=%E8%B6%85%E6%97%A5%E5%B8%B8%E7%BE%8E%E9%A3%9F%E6%95%99%E7%A8%8B%E9%80%9F%E6%9D%A5get) — 886万
5. [定格这一刻的日照金山](https://www.xiaohongshu.com/search_result?keyword=%E5%AE%9A%E6%A0%BC%E8%BF%99%E4%B8%80%E5%88%BB%E7%9A%84%E6%97%A5%E7%85%A7%E9%87%91%E5%B1%B1) — 871.1万
6. [你可以永远相信赛里木湖的美景](https://www.xiaohongshu.com/search_result?keyword=%E4%BD%A0%E5%8F%AF%E4%BB%A5%E6%B0%B8%E8%BF%9C%E7%9B%B8%E4%BF%A1%E8%B5%9B%E9%87%8C%E6%9C%A8%E6%B9%96%E7%9A%84%E7%BE%8E%E6%99%AF) — 783.1万 · 热
7. [拼豆上也可以作画了](https://www.xiaohongshu.com/search_result?keyword=%E6%8B%BC%E8%B1%86%E4%B8%8A%E4%B9%9F%E5%8F%AF%E4%BB%A5%E4%BD%9C%E7%94%BB%E4%BA%86) — 749.8万 · 热
8. [我的家庭旅行更像是打副本](https://www.xiaohongshu.com/search_result?keyword=%E6%88%91%E7%9A%84%E5%AE%B6%E5%BA%AD%E6%97%85%E8%A1%8C%E6%9B%B4%E5%83%8F%E6%98%AF%E6%89%93%E5%89%AF%E6%9C%AC) — 737.1万
9. [原来古诗词里的河南真的存在](https://www.xiaohongshu.com/search_result?keyword=%E5%8E%9F%E6%9D%A5%E5%8F%A4%E8%AF%97%E8%AF%8D%E9%87%8C%E7%9A%84%E6%B2%B3%E5%8D%97%E7%9C%9F%E7%9A%84%E5%AD%98%E5%9C%A8) — 727.8万 · 热
10. [蒸出了奶香爆米花馒头](https://www.xiaohongshu.com/search_result?keyword=%E8%92%B8%E5%87%BA%E4%BA%86%E5%A5%B6%E9%A6%99%E7%88%86%E7%B1%B3%E8%8A%B1%E9%A6%92%E5%A4%B4) — 716.2万 · 热

[完整快照](data/latest-hourly.json) · [今日归档](archives/2026-11-07.md)
<!-- END HOTSEARCH -->

## 特性

- **三个平台，一种结构**：微博、抖音、小红书统一为 `rank / title / url / hot / tag`，平台特有字段放进 `extra`。
- **实时与归档分离**：高频运行只发布 Pages 的 `latest.json`，不写 git；每小时才落一份归档（最新归档快照为 `data/latest-hourly.json`），避免仓库历史膨胀。
- **过滤广告**：微博接口里混在 `realtime` 中的推广位（`is_ad=1`）会被剔除，不占榜单名次。
- **限流友好**：识别小红书 `300013` 限流码，不连续重试；普通失败按 2s → 4s → 8s 指数退避。
- **零数据库依赖**：默认写文件，可直接用 GitHub Pages、对象存储或任意静态服务消费。
- **可持续归档**：小时快照 + 按天重建的 `summary.json` + 每日 Markdown 归档，重复运行或补抓都不会污染统计。
- **常驻模式**：`daemon` 子命令按间隔运行，连续失败时自动退避到 15 分钟。

## 快速开始

需要 Go 1.27+。

```bash
git clone https://github.com/solarhell/hot-search-crawler.git
cd hot-search-crawler

# 抓取一次，并更新 README 里的今日热榜
go run ./cmd/hot-search-crawler crawl --update-readme

# 只看看当前数据，不落盘
go run ./cmd/hot-search-crawler crawl --dry-run

# 只生成实时 JSON（用于 Pages 发布，不写小时快照/日报）
go run ./cmd/hot-search-crawler crawl --live --data dist
```

也可以直接安装：

```bash
go install github.com/solarhell/hot-search-crawler/cmd/hot-search-crawler@latest
```

构建二进制：

```bash
go build -o bin/hot-search-crawler ./cmd/hot-search-crawler
./bin/hot-search-crawler crawl
```

## 输出结构

```text
data/
  2026-10-30/
    09.json          # 每小时一次完整抓取（含各平台成功/失败状态）
    10.json
    summary.json     # 当天聚合：最高热度、最好名次、在榜次数、首次/最后出现
  latest-hourly.json # 最近一次小时归档快照；实时数据见 GitHub Pages
archives/
  2026-10-30.md      # 当天热榜 Markdown 归档
README.md            # 标记区间内的今日热榜
```

小时快照：

```json
{
  "schema_version": 1,
  "generated_at": "2026-10-30T09:32:05+08:00",
  "timezone": "Asia/Shanghai",
  "boards": [
    {
      "platform": "weibo",
      "name": "微博",
      "ok": true,
      "fetched_at": "2026-10-30T09:32:05+08:00",
      "elapsed_ms": 126,
      "items": [
        {
          "rank": 1,
          "title": "TES战胜BLG",
          "url": "https://s.weibo.com/weibo?q=TES%E6%88%98%E8%83%9CBLG",
          "hot": 1748429,
          "tag": "新"
        }
      ]
    }
  ]
}
```

抓取失败的平台也会保留在快照里，`ok` 为 `false`，`error` 写原因，`items` 为空数组。单个平台失败不会影响其他平台。

小红书等平台的热度是 `"947.5w"` 这样的文本，`hot` 会归一化成 `9475000`，同时用 `hot_text` 保留原始值。

## 平台说明

| 平台 | 接口 | 热度字段 | 备注 |
| --- | --- | --- | --- |
| 微博 | `weibo.com/ajax/side/hotSearch` | `num` | 使用 Chrome TLS 指纹；过滤 `is_ad=1` |
| 抖音 | `aweme.snssdk.com/aweme/v1/hot/search/list/` | `hot_value` | 只取 `word_list`，不混入“实时上升热点” |
| 小红书 | `edith.xiaohongshu.com/api/sns/v1/search/hot_list` | `score`（`947.5w`） | 复用固定 iOS 设备指纹，可能被限流 |

小红书没有公开热榜接口，`internal/source/xiaohongshu.go` 中的请求头来自客户端抓包。同一设备指纹被间歇性拒绝（`300013`）是常态，程序会按限流处理并保留其他平台的数据。指纹可能随客户端升级失效，届时需要重新抓包，或用 `--sources weibo,douyin` 暂时关闭小红书。

## 命令行

```text
hot-search-crawler crawl [flags]     抓取一次
hot-search-crawler daemon [flags]    常驻定时抓取
hot-search-crawler sources           列出内置平台
hot-search-crawler version           显示版本
```

常用参数：

| 参数 | 默认值 | 说明 |
| --- | --- | --- |
| `--sources` | `weibo,douyin,xiaohongshu` | 启用的平台，逗号分隔 |
| `--top` | `20` | 每个平台保留的条数 |
| `--data` | `data` | 小时快照目录 |
| `--archives` | `archives` | 每日 Markdown 归档目录 |
| `--readme` | `README.md` | README 路径 |
| `--update-readme` | `false` | 把今日热榜写入 README |
| `--live` | `false` | 只写实时 `latest.json`，不生成小时快照/日报/README |
| `--dry-run` | `false` | 只输出快照 JSON，不写文件 |
| `--timeout` | `25s` | 单次请求超时 |
| `--retries` | `2` | 普通失败的重试次数（限流不重试） |
| `--interval` | `30m` | 仅 `daemon`：抓取间隔 |
| `--debug` | `false` | 输出调试日志 |

`crawl` 在至少一个平台成功时返回 0；全部失败返回 1，但失败快照仍会写入，便于排查。`daemon` 不会因单轮失败退出。

## 更新频率与调度

| 数据 | 频率 | 发布位置 | 触发方式 |
| --- | --- | --- | --- |
| 实时快照 | 每 10 分钟 | GitHub Pages `latest.json` | Cloudflare Worker → `live.yml`；GitHub `schedule` `:13/:43` 兜底 |
| 小时归档 | 每小时 | `data/`、`archives/`、README | Cloudflare Worker `:23` + GitHub `schedule` `:47` 兜底 |

GitHub Actions 自带的 `schedule` 是 best-effort：2026 年实测常见延迟 1.5–4.5 小时，低活跃仓库的每小时任务实际每天只跑 4–8 次，丢跑不会留下记录。因此实时链路改为 Cloudflare Worker 调用 `workflow_dispatch`，部署步骤见 **[Cloudflare Worker 文档](deploy/cloudflare-worker/README.md)**；cron-job.org 和 systemd timer 的替代方案见 [外部定时器文档](docs/external-scheduler.md)。

高频运行只发布 Pages，不写 git，原因是单次快照约 24KB，10 分钟一次会让仓库一年增长数百 MB；归档仍保持每小时一份，兼顾历史价值和仓库体积。

## GitHub Actions

仓库内置三个工作流：

- [`live.yml`](.github/workflows/live.yml)：由 Cloudflare Worker 每 10 分钟触发，抓取后只把 `latest.json` 部署到 GitHub Pages，不产生 git 提交；GitHub `schedule` 在 `:13/:43` 做低频兜底。
- [`crawl.yml`](.github/workflows/crawl.yml)：每小时归档，写入小时快照、日报和 README。Cloudflare Worker 在 `:23` 触发，GitHub 自带 `schedule` 在 `:47` 兜底；如果当前小时已有快照则自动跳过，避免重复提交。
- [`test.yml`](.github/workflows/test.yml)：只在代码变更时运行 `go vet` 和 `go test -race`，数据提交不会触发。

第一次使用需要：

1. 在 **Settings → Actions → General** 确认 Workflow permissions 为 **Read and write permissions**。
2. 在 **Settings → Pages** 把 Source 设为 **GitHub Actions**（或运行一次 `live.yml`，按提示授权）。
3. 按 [Cloudflare Worker 部署文档](deploy/cloudflare-worker/README.md) 创建 fine-grained PAT 并部署 Worker；也支持 cron-job.org 或自己的 systemd timer，见 [外部定时器文档](docs/external-scheduler.md)。

### 部署 Cloudflare Worker

```bash
cd deploy/cloudflare-worker
npm install
npx wrangler login
npx wrangler secret put GITHUB_TOKEN
npm run deploy
```

Worker 的两个 cron 会分别触发 `live.yml`（每 10 分钟）和 `crawl.yml`（每小时 :23）。完整说明和失败告警配置见 [deploy/cloudflare-worker/README.md](deploy/cloudflare-worker/README.md)。

如果 GitHub Runner 访问国内接口不稳定，可以在 workflow 中给抓取步骤设置 `HTTPS_PROXY`，surf 会读取标准代理环境变量：

```yaml
- name: Crawl
  env:
    HTTPS_PROXY: ${{ secrets.HTTPS_PROXY }}
  run: go run ./cmd/hot-search-crawler crawl --update-readme
```

## 新增平台

实现 `internal/source.Source`：

```go
type Source interface {
    ID() string
    Name() string
    Fetch(ctx context.Context, limit int) ([]Item, error)
}
```

然后在 `internal/source/registry.go` 的 `All()` 里注册。解析逻辑建议拆成纯函数并用 `testdata/` 中的真实报文做单元测试，HTTP 请求、超时、重试和落盘由框架统一处理。

## 设计取舍

- **每次从小时文件重建日报**：`summary.json` 不是增量合并，而是读取当天所有 `HH.json` 重新聚合。这样重复运行、手动补抓、删除异常快照都不会让统计漂移。
- **广告与空标题在计数前过滤**：`--top` 限制的是有效条目数，不会被推广位占掉名额。
- **限流不重试**：小红书 `300013` 的 `RetryAfter` 是 10 分钟，连续敲门只会延长封禁；普通网络错误才走指数退避。
- **固定北京时间**：文件名和日报统一用 UTC+8，避免服务器时区不同导致目录漂移。

## 免责声明

本项目仅用于学习和技术研究。热搜数据归各平台所有，请遵守目标网站的 robots.txt、服务条款和当地法律法规；不要高频请求、绕过付费能力或用于商业用途。因使用本项目产生的任何后果由使用者自行承担。

## License

[MIT](LICENSE)
