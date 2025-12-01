# 个人商城项目 - 前端集成说明

## 概述

已为您的个人商城项目创建了完整的前端代码，使用 Vue 3 + Vite 构建，采用简约现代的设计风格。

## 已实现的功能

### 1. 用户认证模块
- ✅ 用户注册（包含16位密钥设置）
- ✅ 用户登录
- ✅ 邮箱绑定（发送验证邮件）
- ✅ 邮箱验证页面

### 2. 商品管理模块
- ✅ 商品列表展示
- ✅ 商品详情页
- ✅ 商品搜索
- ✅ 商品上架（支持多图片上传）

### 3. 购物车模块
- ✅ 添加商品到购物车
- ✅ 购物车列表
- ✅ 修改商品数量
- ✅ 删除商品
- ✅ 购物车结算

### 4. 订单管理模块
- ✅ 创建订单
- ✅ 订单列表
- ✅ 订单支付（需要16位密钥）
- ✅ 收货地址管理

### 5. 个人中心模块
- ✅ 个人信息展示
- ✅ 昵称修改
- ✅ 头像上传
- ✅ 余额查询（需要16位密钥）

## 项目结构

```
frontend/
├── src/
│   ├── views/              # 页面组件
│   │   ├── Home.vue        # 首页
│   │   ├── Login.vue       # 登录页
│   │   ├── Register.vue    # 注册页
│   │   ├── Products.vue    # 商品列表
│   │   ├── ProductDetail.vue  # 商品详情
│   │   ├── Cart.vue        # 购物车
│   │   ├── Orders.vue      # 订单管理
│   │   ├── Sell.vue        # 上架商品
│   │   ├── Profile.vue     # 个人中心
│   │   └── ValidEmail.vue  # 邮箱验证
│   ├── stores/             # 状态管理
│   │   └── user.js         # 用户状态
│   ├── router/             # 路由配置
│   │   └── index.js
│   ├── utils/              # 工具函数
│   │   └── api.js          # API封装
│   ├── App.vue             # 根组件
│   ├── main.js             # 入口文件
│   └── style.css           # 全局样式
├── Dockerfile              # Docker构建文件
├── nginx.conf              # Nginx配置
├── package.json
└── vite.config.js
```

## API对接说明

前端已完全对接后端API，主要接口包括：

- `POST /api/v1/user/register` - 用户注册
- `POST /api/v1/user/login` - 用户登录
- `PUT /api/v1/user` - 更新用户信息
- `POST /api/v1/avatar` - 上传头像
- `POST /api/v1/user/sending-email` - 发送验证邮件
- `POST /api/v1/user/valid-email` - 验证邮箱
- `POST /api/v1/money` - 查询余额
- `GET /api/v1/products` - 获取商品列表
- `POST /api/v1/products` - 搜索商品
- `GET /api/v1/products/:id` - 获取商品详情
- `POST /api/v1/product` - 创建商品
- `POST /api/v1/carts` - 创建购物车
- `GET /api/v1/carts` - 获取购物车列表
- `PUT /api/v1/carts/:id` - 更新购物车
- `DELETE /api/v1/carts/:id` - 删除购物车
- `POST /api/v1/orders` - 创建订单
- `GET /api/v1/orders` - 获取订单列表
- `POST /api/v1/paydown` - 订单支付
- `POST /api/v1/addresses` - 创建地址
- `GET /api/v1/addresses` - 获取地址列表

## 设计特点

1. **简约风格**：采用简洁的卡片式设计，清晰的层次结构
2. **响应式布局**：适配不同屏幕尺寸
3. **用户体验**：友好的错误提示和加载状态
4. **安全性**：Token认证，自动处理登录状态

## 开发环境运行

```bash
cd frontend
npm install
npm run dev
```

访问 http://localhost:8000

## 生产环境构建

```bash
cd frontend
npm run build
```

构建产物在 `frontend/dist` 目录

## Docker部署

前端已配置Docker支持，使用Nginx作为Web服务器。详见 `DEPLOY.md`

## 注意事项

1. **密钥管理**：用户注册时需要16位密钥，用于加密账户余额
2. **图片路径**：后端返回的图片路径已包含完整URL，前端直接使用
3. **API代理**：开发环境通过Vite代理到后端，生产环境通过Nginx代理
4. **Token存储**：Token存储在localStorage中，页面刷新后自动恢复登录状态

## 后续优化建议

1. 添加商品分类筛选
2. 优化图片加载（懒加载、压缩）
3. 添加订单状态跟踪
4. 实现商品收藏功能
5. 添加用户评价功能
6. 优化移动端体验

