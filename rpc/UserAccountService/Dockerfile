# ========== 阶段1：构建阶段（debian系，支持动态编译） ==========
FROM golang:1.23.2 AS builder

# 1. 换debian国内源（阿里云），加速apt下载
RUN sed -i 's/deb.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list.d/debian.sources && \
    apt update && apt install -y --no-install-recommends git make bash gcc && \
    apt clean && rm -rf /var/lib/apt/lists/*

# 2. 换Go国内模块源（阿里云），加速go mod download
RUN go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct

# 3. 工作目录（匹配你的项目根目录）
WORKDIR /app

# 4. 复制项目所有文件到容器（你的目录：Dockerfile、build.sh、script/、idl/等）
COPY . .

# 5. 给build.sh加执行权限，执行编译（核心：运行build.sh，而非make build）
RUN chmod +x ./build.sh && ./build.sh

# ========== 阶段2：运行阶段（轻量debian，兼容动态编译） ==========
FROM debian:bookworm-slim

# 1. 换debian国内源，加速依赖安装
RUN sed -i 's/deb.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list.d/debian.sources && \
    apt update && apt install -y --no-install-recommends bash && \
    apt clean && rm -rf /var/lib/apt/lists/*

# 2. 创建非root用户（安全优化，避免root运行）
RUN groupadd -r appgroup && useradd -r -g appgroup appuser

# 3. 工作目录
WORKDIR /app

# 4. 从构建阶段复制编译产物（关键：output目录包含bootstrap.sh+二进制）
COPY --from=builder /app/output ./output

# 5. 给启动脚本+二进制加执行权限（适配你的UserAccountService）
RUN chmod +x ./output/bootstrap.sh && \
    chmod +x ./output/bin/UserAccountService

# 6. 环境变量（适配用户账户服务，按需调整）
ENV REDIS_ADDR=redis:6379 \
    ETCD_ADDR=etcd:2379 \
    LOG_LEVEL=info

# 7. 暴露端口（UserAccountService 端口，按需调整，比如8888）
EXPOSE 8888

# 8. 切换非root用户（安全要求）
USER appuser

# 9. 容器启动：运行output/bootstrap.sh（和你要求的一致）
CMD ["./output/bootstrap.sh"]