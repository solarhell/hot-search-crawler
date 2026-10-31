# 热搜爬虫

抓取**微博、抖音、小红书**热榜，输出统一结构的 JSON 快照。配合仓库自带的 GitHub Actions，每小时自动抓取并提交数据，仓库本身就是一份可检索的热搜历史。

[![crawl](https://github.com/solarhell/hot-search-crawler/actions/workflows/crawl.yml/badge.svg)](https://github.com/solarhell/hot-search-crawler/actions/workflows/crawl.yml)

<!-- BEGIN HOTSEARCH -->
> 数据更新于 **2026-10-31 14:56**（Asia/Shanghai），展示今日榜单前 10。

### 微博

1. [吾恩去世](https://s.weibo.com/weibo?q=%E5%90%BE%E6%81%A9%E5%8E%BB%E4%B8%96) — 132.8万 · 新
2. [卖650万捐赠房产要缴450万的税](https://s.weibo.com/weibo?q=%23%E5%8D%96650%E4%B8%87%E6%8D%90%E8%B5%A0%E6%88%BF%E4%BA%A7%E8%A6%81%E7%BC%B4450%E4%B8%87%E7%9A%84%E7%A8%8E%23) — 122.2万
3. [黄金](https://s.weibo.com/weibo?q=%E9%BB%84%E9%87%91) — 95.5万 · 热
4. [昂跑未与姆巴佩解约或两头皆输](https://s.weibo.com/weibo?q=%23%E6%98%82%E8%B7%91%E6%9C%AA%E4%B8%8E%E5%A7%86%E5%B7%B4%E4%BD%A9%E8%A7%A3%E7%BA%A6%E6%88%96%E4%B8%A4%E5%A4%B4%E7%9A%86%E8%BE%93%23) — 88.7万
5. [开创新时代管党治党兴党强党新局面](https://s.weibo.com/weibo?q=%23%E5%BC%80%E5%88%9B%E6%96%B0%E6%97%B6%E4%BB%A3%E7%AE%A1%E5%85%9A%E6%B2%BB%E5%85%9A%E5%85%B4%E5%85%9A%E5%BC%BA%E5%85%9A%E6%96%B0%E5%B1%80%E9%9D%A2%23) — 71万
6. [六网协同夯实中国式现代化底座](https://s.weibo.com/weibo?q=%23%E5%85%AD%E7%BD%91%E5%8D%8F%E5%90%8C%E5%A4%AF%E5%AE%9E%E4%B8%AD%E5%9B%BD%E5%BC%8F%E7%8E%B0%E4%BB%A3%E5%8C%96%E5%BA%95%E5%BA%A7%23) — 69.6万
7. [微博视界大会红毯](https://s.weibo.com/weibo?q=%23%E5%BE%AE%E5%8D%9A%E8%A7%86%E7%95%8C%E5%A4%A7%E4%BC%9A%E7%BA%A2%E6%AF%AF%23) — 49.9万 · 新
8. [披荆斩棘成团名单](https://s.weibo.com/weibo?q=%E6%8A%AB%E8%8D%86%E6%96%A9%E6%A3%98%E6%88%90%E5%9B%A2%E5%90%8D%E5%8D%95) — 37.5万 · 沸
9. [今年全国实现生娃基本不花钱](https://s.weibo.com/weibo?q=%23%E4%BB%8A%E5%B9%B4%E5%85%A8%E5%9B%BD%E5%AE%9E%E7%8E%B0%E7%94%9F%E5%A8%83%E5%9F%BA%E6%9C%AC%E4%B8%8D%E8%8A%B1%E9%92%B1%23) — 32.6万
10. [小沈阳披荆斩棘冠军](https://s.weibo.com/weibo?q=%23%E5%B0%8F%E6%B2%88%E9%98%B3%E6%8A%AB%E8%8D%86%E6%96%A9%E6%A3%98%E5%86%A0%E5%86%9B%23) — 32.4万 · 热

### 抖音

1. [冬天奶茶要挖着吃才满足](https://www.douyin.com/search/%E5%86%AC%E5%A4%A9%E5%A5%B6%E8%8C%B6%E8%A6%81%E6%8C%96%E7%9D%80%E5%90%83%E6%89%8D%E6%BB%A1%E8%B6%B3) — 1146.5万
2. [中国稀土筹划重大资产重组](https://www.douyin.com/search/%E4%B8%AD%E5%9B%BD%E7%A8%80%E5%9C%9F%E7%AD%B9%E5%88%92%E9%87%8D%E5%A4%A7%E8%B5%84%E4%BA%A7%E9%87%8D%E7%BB%84) — 1129.7万
3. [中共中央举行新闻发布会](https://www.douyin.com/search/%E4%B8%AD%E5%85%B1%E4%B8%AD%E5%A4%AE%E4%B8%BE%E8%A1%8C%E6%96%B0%E9%97%BB%E5%8F%91%E5%B8%83%E4%BC%9A) — 1128.3万
4. [以色列选举最终结果公布](https://www.douyin.com/search/%E4%BB%A5%E8%89%B2%E5%88%97%E9%80%89%E4%B8%BE%E6%9C%80%E7%BB%88%E7%BB%93%E6%9E%9C%E5%85%AC%E5%B8%83) — 1115.2万
5. [攻击-21无人机在四川舰成功起降](https://www.douyin.com/search/%E6%94%BB%E5%87%BB-21%E6%97%A0%E4%BA%BA%E6%9C%BA%E5%9C%A8%E5%9B%9B%E5%B7%9D%E8%88%B0%E6%88%90%E5%8A%9F%E8%B5%B7%E9%99%8D) — 1113.5万
6. [进博会展览展示面积首破44万平](https://www.douyin.com/search/%E8%BF%9B%E5%8D%9A%E4%BC%9A%E5%B1%95%E8%A7%88%E5%B1%95%E7%A4%BA%E9%9D%A2%E7%A7%AF%E9%A6%96%E7%A0%B444%E4%B8%87%E5%B9%B3) — 1110.4万
7. [不怕远征难](https://www.douyin.com/search/%E4%B8%8D%E6%80%95%E8%BF%9C%E5%BE%81%E9%9A%BE) — 1110万
8. [比冬天先来的是一顿热火锅](https://www.douyin.com/search/%E6%AF%94%E5%86%AC%E5%A4%A9%E5%85%88%E6%9D%A5%E7%9A%84%E6%98%AF%E4%B8%80%E9%A1%BF%E7%83%AD%E7%81%AB%E9%94%85) — 1109.9万
9. [二十届五中全会公报一图速览](https://www.douyin.com/search/%E4%BA%8C%E5%8D%81%E5%B1%8A%E4%BA%94%E4%B8%AD%E5%85%A8%E4%BC%9A%E5%85%AC%E6%8A%A5%E4%B8%80%E5%9B%BE%E9%80%9F%E8%A7%88) — 1102.3万
10. [落日与风能发电机同框太治愈](https://www.douyin.com/search/%E8%90%BD%E6%97%A5%E4%B8%8E%E9%A3%8E%E8%83%BD%E5%8F%91%E7%94%B5%E6%9C%BA%E5%90%8C%E6%A1%86%E5%A4%AA%E6%B2%BB%E6%84%88) — 1053.7万

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

[完整快照](data/latest.json) · [今日归档](archives/2026-10-31.md)
<!-- END HOTSEARCH -->

## 特性

- **三个平台，一种结构**：微博、抖音、小红书统一为 `rank / title / url / hot / tag`，平台特有字段放进 `extra`。
- **过滤广告**：微博接口里混在 `realtime` 中的推广位（`is_ad=1`）会被剔除，不占榜单名次。
- **限流友好**：识别小红书 `300013` 限流码，不连续重试；普通失败按 2s → 4s → 8s 指数退避。
- **零数据库依赖**：默认写文件，可直接用 GitHub Pages、对象存储或任意静态服务消费。
- **可持续归档**：小时快照 + 按天重建的 `summary.json` + 每日 Markdown 归档，重复运行或补抓都不会污染统计。
- **可选翻译**：内置 OpenAI 兼容接口的批量翻译，三个国内平台不需要，留给后续海外数据源。
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
    09.json          # 09:07 的一次完整抓取（含各平台成功/失败状态）
    10.json
    summary.json     # 当天聚合：最高热度、最好名次、在榜次数、首次/最后出现
  latest.json        # 最近一次快照，适合直接给前端消费
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
| `--dry-run` | `false` | 只输出快照 JSON，不写文件 |
| `--timeout` | `25s` | 单次请求超时 |
| `--retries` | `2` | 普通失败的重试次数（限流不重试） |
| `--interval` | `30m` | 仅 `daemon`：抓取间隔 |
| `--debug` | `false` | 输出调试日志 |

`crawl` 在至少一个平台成功时返回 0；全部失败返回 1，但失败快照仍会写入，便于排查。`daemon` 不会因单轮失败退出。

## GitHub Actions 持续归档

仓库内置 [`.github/workflows/crawl.yml`](.github/workflows/crawl.yml)：

- 每小时第 7 分钟运行一次，也支持手动触发。
- 抓取后把 `data/`、`archives/` 和 README 的变更提交回仓库。
- 使用 `GITHUB_TOKEN` 推送，不需要额外配置 secret。

Fork 或新建仓库后，在 **Settings → Actions → General** 确认 Workflow permissions 为 **Read and write permissions**，然后启用 Actions 即可。

如果 GitHub Runner 访问国内接口不稳定，可以在 workflow 中给抓取步骤设置 `HTTPS_PROXY`，surf 会读取标准代理环境变量：

```yaml
- name: Crawl
  env:
    HTTPS_PROXY: ${{ secrets.HTTPS_PROXY }}
  run: go run ./cmd/hot-search-crawler crawl --update-readme
```

## 可选翻译

当前三个平台都是中文，不需要翻译。翻译模块是为新增海外数据源预留的，接口为 OpenAI 兼容的 Chat Completions：

```bash
export TRANSLATE_BASE_URL=https://api.openai.com/v1   # 可换成任意兼容服务
export TRANSLATE_API_KEY=sk-...
export TRANSLATE_MODEL=gpt-4.1-mini

go run ./cmd/hot-search-crawler crawl --translate
```

数据源通过 `Translatable()` 声明是否需要翻译，翻译失败时自动降级为只保留原文，不影响抓取结果。

## 新增平台

实现 `internal/source.Source`：

```go
type Source interface {
    ID() string
    Name() string
    Translatable() bool
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
