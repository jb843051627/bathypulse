# bathypulse

`bathypulse` 是深海微震观测站的波形事件编排服务。它把台站上传的采样窗口
落到 SQLite 文件数据库，聚合事件、告警和维护状态，并提供一组用于观测值班
的 HTTP 接口。

## 运行

```bash
GOTOOLCHAIN=local go build ./...
GOTOOLCHAIN=local go test ./...
GOTOOLCHAIN=local go run . --smoke-test
```

默认数据库路径是 `data/bathypulse.db`，可用 `BATHYPULSE_DB` 覆盖。服务监听
`BATHYPULSE_ADDR`，默认 `127.0.0.1:8097`。
