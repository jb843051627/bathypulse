# bathypulse 容器运行说明

镜像使用固定的 `golang:1.22-bookworm` 构建环境，依赖在镜像构建阶段缓存。
容器默认启动 HTTP 服务，数据库文件位于 `/app/data/bathypulse.db`。
