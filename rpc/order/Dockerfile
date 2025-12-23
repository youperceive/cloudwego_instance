# ========== 阶段1：构建阶段（debian系，支持动态编译） ==========
FROM golang:1.22.2 AS builder

# 1. 换debian国内源（阿里云），加速apt下载
RUN sed -i 's/deb.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list.d/debian.sources && \
    apt update && apt install -y --no-install-recommends git make bash gcc && \
    apt clean && rm -rf /var/lib/apt/lists/*

# 2. 换Go国内模块源（阿里云），加速go mod download
RUN go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct

# 3. 工作目录
WORKDIR /app

# 4. 复制项目所有文件到构建目录
COPY . .

# 5. 给build.sh加执行权限，执行编译（匹配你的make build+build.sh）
RUN chmod +x ./build.sh && make build

# ========== 阶段2：运行阶段（轻量debian，兼容动态编译） ==========
FROM debian:bookworm-slim

# 1. 换debian国内源（加速后续安装）
RUN sed -i 's/deb.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list.d/debian.sources && \
    apt update && apt install -y --no-install-recommends bash && \
    DEBIAN_FRONTEND=noninteractive apt update && \
    DEBIAN_FRONTEND=noninteractive apt install -y --no-install-recommends tzdata && \
    # 设置上海时区（匹配业务时间戳）
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo 'Asia/Shanghai' > /etc/timezone && \
    apt clean && rm -rf /var/lib/apt/lists/*

# 2. 创建非root用户（安全优化，避免root运行）
RUN groupadd -r appgroup && useradd -r -g appgroup appuser

# 3. 工作目录
WORKDIR /app

# 4. 从构建阶段复制编译产物（output目录，匹配你的项目结构）
COPY --from=builder /app/output ./output

# 6. 暴露订单服务端口（匹配你日志里的8002端口）
EXPOSE 8002

# 7. 切换非root用户（安全最佳实践）
USER appuser

# 8. 启动订单服务（匹配你的bootstrap.sh启动逻辑）
CMD ["./output/bootstrap.sh"]