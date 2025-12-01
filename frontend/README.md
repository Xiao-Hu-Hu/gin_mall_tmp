# 前端项目说明

## 技术栈

- Vue 3
- Vite
- Vue Router
- Pinia
- Axios

## 开发

```bash
# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build
```

## 项目结构

```
frontend/
├── src/
│   ├── views/          # 页面组件
│   ├── stores/        # 状态管理
│   ├── router/        # 路由配置
│   ├── utils/         # 工具函数
│   ├── App.vue        # 根组件
│   ├── main.js        # 入口文件
│   └── style.css      # 全局样式
├── index.html
├── vite.config.js
└── package.json
```

## 功能模块

1. **用户认证**
   - 用户注册
   - 用户登录
   - 邮箱绑定验证

2. **商品管理**
   - 商品列表展示
   - 商品详情
   - 商品搜索
   - 商品上架

3. **购物车**
   - 添加商品到购物车
   - 修改商品数量
   - 删除商品

4. **订单管理**
   - 创建订单
   - 订单列表
   - 订单支付
   - 地址管理

5. **个人中心**
   - 个人信息管理
   - 头像上传
   - 余额查询

