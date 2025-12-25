# 订单微服务（OrderService）
基于 CloudWeGo Kitex 构建的分布式订单微服务，提供订单的创建、查询、分页查询、更新全生命周期管理，适配 Thrift 协议，底层基于 MongoDB 存储。

## 一、项目介绍
### 核心功能
| 接口          | 功能说明                     |
|---------------|------------------------------|
| Create        | 创建订单（生成唯一订单项ID） |
| QueryInfo     | 根据订单ID查询订单详情       |
| QueryOrderId  | 分页查询指定条件的订单ID列表 |
| Update        | 更新订单状态/扩展字段        |

### 技术栈
- **框架**：CloudWeGo Kitex（高性能 RPC 框架）
- **协议**：Thrift
- **存储**：MongoDB（文档型数据库）
- **ID生成**：雪花算法（Snowflake）
- **部署**：Docker（容器化部署）
- **语言**：Go 1.22.2

## 二、环境要求
| 依赖         | 版本要求       | 备注                     |
|--------------|----------------|--------------------------|
| Go           | ≥ 1.22.0       | 推荐1.22.2               |
| MongoDB      | ≥ 6.0          | 需配置用户名/密码（root）|
| Kitex        | latest         | `go install github.com/cloudwego/kitex/tool/cmd/kitex@latest` |
| Thriftgo     | latest         | `go install github.com/cloudwego/thriftgo@latest` |
| Docker       | ≥ 20.10        | 可选，容器化部署用       |

## 三、快速开始（本地部署）
### 1. 克隆项目
```bash
git clone <你的项目仓库地址>
cd cloudwego_instance/rpc/order
```

### 2. 配置环境变量
```bash
# 必选：MongoDB 连接地址（替换为你的Mongo地址）
export MONGODB_URI="mongodb://root:root123456@localhost:27017/order_db?authSource=admin"
```

### 3. 编译&启动
```bash
# 编译项目
make build

# 启动服务
./output/bootstrap.sh
```

### 4. 验证启动成功
查看日志输出，出现以下内容即为成功：
```
2025/12/23 15:33:36.772192 mongo.go:45: [Info] Pinged your deployment. You successfully connected to MongoDB!
2025/12/23 15:33:36.772898 server.go:79: [Info] KITEX: server listen at addr=[::]:8002
```

## 四、核心接口说明
### 1. Create（创建订单）
#### 请求参数（OrderService.CreateRequest）
| 字段         | 类型        | 说明                     |
|--------------|-------------|--------------------------|
| Type         | int32       | 订单类型                 |
| Status       | int32       | 订单状态（1-待支付，2-已支付） |
| ReqUserId    | int64       | 下单用户ID               |
| RespUserId   | int64       | 接单用户ID（可选）|
| Items        | []OrderItem | 订单项列表               |
| Ext          | map[string]string | 扩展字段         |

#### 返回结果
| 字段         | 类型        | 说明                     |
|--------------|-------------|--------------------------|
| BaseResp     | BaseResponse | 通用响应（code/msg）|
| OrderId      | string      | 生成的订单ID（Mongo ObjectID） |

### 2. QueryInfo（查询订单详情）
#### 请求参数
| 字段         | 类型        | 说明                     |
|--------------|-------------|--------------------------|
| Id           | string      | 订单ID                   |

#### 返回结果
| 字段         | 类型        | 说明                     |
|--------------|-------------|--------------------------|
| BaseResp     | BaseResponse | 通用响应                 |
| Order        | Order       | 完整订单信息             |

### 3. QueryOrderId（分页查询订单ID）
#### 请求参数
| 字段         | 类型        | 说明                     |
|--------------|-------------|--------------------------|
| Type         | QueryOrderIdType | 查询类型（REQ_USER/RESP_USER/EXT_KEY） |
| UserId       | int64       | 用户ID（REQ/RESP_USER类型必填） |
| ExtKey       | string      | 扩展字段Key（EXT_KEY类型必填） |
| ExtVal       | string      | 扩展字段Value（EXT_KEY类型必填） |
| Page         | int32       | 页码（默认1）|
| PageSize     | int32       | 页大小（默认20，最大100） |

#### 返回结果
| 字段         | 类型        | 说明                     |
|--------------|-------------|--------------------------|
| BaseResp     | BaseResponse | 通用响应                 |
| OrderId      | []string    | 订单ID列表               |
| Total        | int32       | 总订单数                 |
| Page         | int32       | 当前页码                 |
| PageSize     | int32       | 当前页大小               |

### 4. Update（更新订单）
#### 请求参数
| 字段         | 类型        | 说明                     |
|--------------|-------------|--------------------------|
| Id           | string      | 订单ID（必填）|
| Status       | int32       | 订单状态（可选）|
| Ext          | map[string]string | 扩展字段（覆盖式更新）|

#### 返回结果
| 字段         | 类型        | 说明                     |
|--------------|-------------|--------------------------|
| BaseResp     | BaseResponse | 通用响应                 |

## 五、Docker 部署
### 1. 构建镜像
```bash
# 在项目根目录执行
docker build -t order-service:v1 .
```

### 2. 启动容器
```bash
docker run -d \
  --name order-service \
  -p 8002:8002 \
  -e MONGODB_URI="mongodb://root:root123456@<宿主机IP>:27017/order_db?authSource=admin" \
  --network host \ # 可选：直接使用宿主机网络，无需端口映射
  order-service:v1
```

### 3. 验证容器运行
```bash
# 查看容器日志
docker logs order-service

# 验证接口（示例：查询订单详情）
bash ./scripts_kit/test/query_info.sh
```

## 六、目录结构
```
order/
├── kitex_gen/           # Kitex 生成的Thrift代码（自动生成）
├── pkg/                 # 通用工具包
│   ├── mongo/           # MongoDB 客户端初始化
│   ├── snowflake/       # 雪花算法ID生成
│   └── trans/           # 数据转换（DTO→Model）
├── scripts_kit/         # 脚本目录
│   ├── test/            # 接口测试脚本
│   └── enter_container/ # 进入Mongo容器脚本
├── Dockerfile           # Docker构建文件
├── build.sh             # 编译脚本
├── Makefile             # 构建规则
├── main.go              # 服务入口
└── output/              # 编译产物（自动生成）
    ├── bootstrap.sh     # 启动脚本
    └── bin/             # 二进制文件
```

## 七、常见问题
### 1. Mongo 报 `client is disconnected`
- 原因：`init()` 函数中过早执行 `defer client.Disconnect()`，导致连接刚建立就断开。
- 解决方案：移除 `init()` 中的 `defer Disconnect`，在 `main` 函数退出时手动断开。

### 2. 雪花算法生成的订单项ID重复
- 原因：`snowflake.Epoch` 配置在 `NewNode` 之后，导致Epoch不生效。
- 解决方案：先设置 `snowflake.Epoch`，再创建 `snowflake.Node`。

### 3. QueryOrderId 查不到数据
- 原因：Mongo 字段嵌套层级错误（`requserid` 嵌套在 `order` 子对象中）。
- 解决方案：查询条件改为 `order.requserid`/`order.respuserid`。

### 4. Update 接口更新无效
- 原因：更新字段路径错误（未匹配 Mongo 嵌套字段+无下划线格式）。
- 解决方案：更新路径改为 `order.status`/`order.updatedat`/`order.ext`。

## 八、注意事项
1. **时区**：服务默认使用 `Asia/Shanghai` 时区，保证订单时间戳与本地时间一致。
2. **权限**：Docker 容器使用非 root 用户运行，提高安全性。
3. **Mongo 网络**：容器内访问宿主机 Mongo 时，`MONGODB_URI` 中的地址需填写宿主机实际IP（不要用 `localhost`）。
4. **连接池**：Mongo 客户端已配置连接池（最大50，最小10），避免高频调用导致连接耗尽。
5. **参数校验**：所有接口均做参数非空/格式校验，非法参数直接返回 `Code_INVALID_PARAM`。

## 九、扩展说明
- 后续可集成 Eino 实现 AI 智能订单查询（如“查询100001用户最近的已支付订单”）。
- 可新增「取消订单」接口，联动库存服务回滚商品库存。
- 可接入链路追踪（如 OpenTelemetry），监控服务调用链路。