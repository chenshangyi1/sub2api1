# Claude 本地计费规则与验证（2026-10-11）

用户最终决定采用上游贾维斯的模型计价。已移除本站 Fable 额外 ×1.2 和 Max ×3，默认价卡直接使用贾维斯的单价。自定义分组/渠道价卡仍优先，且不额外叠加这两项售价倍率。

本次只修改本地文件并离线验证，没有部署、修改生产数据库或产生真实扣费。分组/用户倍率继续读取现有配置，老狗 Grok 本次未改。

## 生效规则

模型标准费用 = 输入费用 + 输出费用 + 缓存写入费用 + 缓存读取费用；实际扣费 = 模型标准费用 ×现有分组/用户倍率。渠道已有分时倍率、显式自定义价卡继续按配置生效。

Fable 普通、High、xhigh、Max 使用同一 token 单价。Max 可能产生更多 token，按实际用量计算，不再由本站把单价乘三。默认 Fable 输入/输出为 $12/$60/M；若分组倍率 1.1，实付单价为 $13.2/$66/M。自定义输入价 $2/M、倍率 1.1 时，所有思考档位均为 $2.2/M。

预扣、Adaptive HOLD、流式补扣和最终结算均使用该价卡。模型广场分组售价使用贾维斯价卡；原有参考列独立展示，不用于结算。

## 价卡与证据

单位 USD / 1M tokens。表中为本地采用的模型价卡，未乘分组/用户倍率。

| 模型 | 输入 | 输出 | 5m 写入 | 1h 写入 | 读取 | 证据 |
| --- | ---: | ---: | ---: | ---: | ---: | --- |
| Opus 5.5 | 4 | 20 | 5 | 8 | 0.20 | 贾维斯实际账单核实全部组件 |
| Sonnet 5 | 2 | 10 | 2.5 | 4 | 0.20 | 贾维斯实际账单核实输入、输出、5m 写入和读取；1h 按既有比例 |
| Haiku 4.5 | 1 | 5 | 1.25 | 2 | 0.10 | 贾维斯实际账单核实输入和输出；读取公示价，写入按既有比例 |
| Fable 5 | 12 | 60 | 15 | 24 | 1.20 | 贾维斯公示输入、输出、读取；写入按既有比例 |
| Fable 5.1 | 12 | 60 | 15 | 24 | 0.30 | 贾维斯公示输入、输出、读取；写入按既有比例 |
| Sonnet 5.5 | 2 | 10 | 2.5 | 4 | 0.20 | 当前目录，与 Sonnet 5 一致；无该精确版本的贾维斯账单 |
| Opus 5 / 4.8 | 5 | 25 | 6.25 | 10 | 0.50 | 贾维斯公示价；写入按既有比例 |
| Sonnet 4.6 | 3 | 15 | 3.75 | 6 | 0.30 | 贾维斯公示价；写入按既有比例 |

贾维斯公示上述价卡采用全部上下文同价，本地不继承旧模型的 200K 加价阶梯。该账号近 30 天扫描 453 条 Opus 记录，没有输入上下文超过 200K 的样本，因此这部分不能声称已通过实际账单验证。Opus Fast ×2 保留既有策略，实际账单样本只验证 Standard。近 30 天没有 Fable 的实际账单；其写入缓存价格仍是按既有比例推算，不能称为已经核实的贾维斯扣费规则。

原始证据：[使用记录](https://api.aijws.com/usage)、[贾维斯模型价目](https://web.aijws.com/models)，本地归档在 upstream-standard-pricing-audit-20261011/usage-proofs.json 和 jws-claude-rule-proofs.json。

## 本次验证

先更新测试期望，旧溢价实现明确出现 Max、自定义价卡、预扣和兼容网关断言失败，再移除溢价实现。

- 计费相关回归含子测试 496 项通过；明确跳过两条既有 Grok 旧失败：TestCalculateCost_Grok46UsesOfficialPricesThenUserMultiplier、TestGetModelPricing_Grok46DynamicCatalogIsNormalizedToOfficialCard。
- Fable 普通/Max/xhigh、全部缓存 TTL、自定义价卡并发隔离、免费倍率、预扣/补扣、兼容网关、价格展示、渠道区间与分时、Opus 及实际账单回放通过。
- 重点计费路径 -race 通过。
- GOEXPERIMENT=jsonv2 go build -o /tmp/xinmc-jws-final-server ./cmd/server 成功。
- git diff --check 通过。

全仓库测试在此前同一会话执行，仍被既有重复声明、过期测试字段/签名和其他模块失败阻断。本次回归使用 go list 的实际生产 GoFiles 加相关测试文件；不声称全仓库通过。完整已观察失败列表在 /tmp/xinmc-jws-full-unit.log，包括 service、handler、setup 编译失败以及 config、platform、repository、router 和 migrations 既有测试失败。

复跑本次离线验证：python3 /tmp/xinmc-jws-final-tests.py regression；python3 /tmp/xinmc-jws-final-tests.py race。日志为 /tmp/xinmc-jws-final-regression.jsonl 和 /tmp/xinmc-jws-final-race.jsonl。
