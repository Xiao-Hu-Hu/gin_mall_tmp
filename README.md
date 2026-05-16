# Gin Mall

基于 Go + Gin 的电商后端 API，涵盖用户认证、商品管理、购物车、订单、异步支付、秒杀、排行榜，以及基于 RAG 的 AI 智能客服。

## 功能特性

- **用户系统** — 注册登录（JWT）、头像上传、邮箱验证、AES 加密余额
- **商品管理** — 商品 CRUD、分类筛选、多图上传、关键词搜索
- **购物车** — 增删改查
- **订单系统** — 创建/查询/删除订单，自动生成订单号
- **异步支付** — 基于 RabbitMQ 的三阶段异步流程：支付扣款 → 库存扣减 → 订单取消
- **秒杀系统** — Redis Lua 脚本原子扣库存 + 异步下单队列
- **排行榜** — 四种 Redis Sorted Set 排行榜（热度、浏览、购买、收藏）
- **收藏 & 地址** — 收藏夹、收货地址 CRUD
- **AI 智能客服** — DeepSeek ReAct Agent，支持商品推荐、RAG 知识库问答、代下单

## 技术栈

| 组件 | 技术 |
|------|------|
| 语言 | Go 1.24 |
| Web 框架 | Gin |
| ORM | GORM（支持读写分离） |
| 数据库 | MySQL |
| 缓存 | Redis |
| 消息队列 | RabbitMQ |
| 向量数据库 | Milvus |
| AI 模型 | DeepSeek（via CloudWeGo Eino） |
| Embedding | SiliconFlow（Qwen3-Embedding-0.6B） |
| 前端 | Vue 3 + Vite + Element Plus |

## 项目结构

```
gin_mall_tmp/
├── api/v1/            # HTTP Handler 层
├── cache/             # Redis 操作（排行榜、秒杀库存、AI 会话）
├── cmd/               # 程序入口 main.go
├── conf/              # 配置加载（INI）
├── dao/               # 数据访问层（GORM 查询）
├── frontend/          # Vue 前端
├── knowledge/         # RAG 知识库文档（退换货、物流、支付政策等）
├── middleware/         # JWT 认证、CORS
├── model/             # 数据模型
├── mq/                # RabbitMQ 连接、生产者、消费者
├── pkg/               # 工具包（JWT、AES 加密、日志）
├── routes/            # 路由定义
├── serializer/        # 响应序列化
├── service/           # 业务逻辑层
│   └── ai/            # AI 客服子系统（Agent、RAG、向量存储）
└── static/            # 静态资源（头像、商品图片）
```

## 快速开始

### 环境要求

- Go 1.24+
- MySQL
- Redis
- RabbitMQ
- Milvus（可选，AI 知识库功能需要）

### 安装与运行

```bash
# 克隆项目
git clone https://github.com/yourname/gin_mall.git
cd gin_mall

# 安装依赖
go mod tidy

# 修改数据库配置
vim conf/config.ini

# 启动（本地）
go run cmd/main.go
```

默认监听 `:8080`，Docker 模式下为 `:3000`。

### 配置文件

编辑 `conf/config.ini`，配置 MySQL、Redis、RabbitMQ 连接信息。

AI 功能需要在 `.env` 或环境变量中配置：

```env
DEEPSEEK_API_KEY=your_key          # 必填
DEEPSEEK_MODEL=deepseek-chat       # 可选
DEEPSEEK_BASE_URL=https://api.deepseek.com  # 可选
SILICONFLOW_API_KEY=your_key       # Embedding 服务
MILVUS_ADDR=localhost:19530        # 向量数据库
```

### 启动 Milvus（可选）

```bash
docker-compose -f docker-compose.milvus.yml up -d
```

启动后调用 `POST /api/v1/ai/knowledge/sync` 同步知识库到向量数据库。

## API 概览

### 公开接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/user/register` | 用户注册 |
| POST | `/api/v1/user/login` | 用户登录 |
| GET | `/api/v1/products` | 商品列表 |
| GET | `/api/v1/products/:id` | 商品详情 |
| GET | `/api/v1/categories` | 分类列表 |
| GET | `/api/v1/carousels` | 轮播图 |
| GET | `/api/v1/products/hot` | 热度排行榜 |
| GET | `/api/v1/products/view-rank` | 浏览排行榜 |
| GET | `/api/v1/products/purchase-rank` | 购买排行榜 |
| GET | `/api/v1/products/favorite-rank` | 收藏排行榜 |
| GET | `/api/v1/seckill/:id` | 秒杀活动详情 |

### 需认证接口（Header: `token`）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/product` | 创建商品 |
| POST | `/api/v1/products` | 搜索商品 |
| CRUD | `/api/v1/carts` | 购物车 |
| CRUD | `/api/v1/orders` | 订单 |
| CRUD | `/api/v1/addresses` | 收货地址 |
| CRUD | `/api/v1/favorites` | 收藏 |
| POST | `/api/v1/paydown` | 发起支付 |
| POST | `/api/v1/seckill/order` | 秒杀下单 |
| POST | `/api/v1/ai/customer-service/chat` | AI 客服对话 |
| GET | `/api/v1/ai/sessions` | AI 会话列表 |
| POST | `/api/v1/ai/knowledge/sync` | 同步知识库 |

完整接口定义见 `routes/router.go`。

## 架构设计要点

- **异步支付**：支付请求发布到 RabbitMQ，由独立消费者在事务中完成扣款、加款、库存扣减，保证数据一致性
- **秒杀高并发**：Redis Lua 脚本原子执行库存校验 + 用户去重 + 库存扣减，异步写入订单
- **RAG 知识库**：政策文档经分块、Embedding 后存入 Milvus，AI 客服通过向量检索增强回答准确性
- **AI 工具调用**：ReAct Agent 可调用 4 个工具 — 搜索商品、查询地址、创建订单、检索知识库
