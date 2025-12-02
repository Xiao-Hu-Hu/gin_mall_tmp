# 服务器部署指南

## 配置修改总结

我已经对以下配置文件进行了优化，确保它们适合服务器部署：

### 1. 后端配置修改
- ✅ **conf/config.docker.ini**: 添加了SMTP端口配置（465端口）
- ✅ **Dockerfile**: 已优化，包含正确的Go版本和依赖管理

### 2. 前端配置修改
- ✅ **frontend/.env.production**: 创建了生产环境变量文件
- ✅ **frontend/Dockerfile**: 添加了环境变量配置
- ✅ **frontend/src/utils/api.js**: 修改为支持环境变量配置
- ✅ **docker-compose.yml**: 添加了前端环境变量配置

### 3. 样式优化
- ✅ **frontend/src/style.css**: 已更新为手绘线条卡通风格

## 部署步骤

### 1. 上传项目到服务器
```bash
# 将整个项目目录上传到服务器
scp -r gin_mall_tmp/ user@your-server-ip:/path/to/project/
```

### 2. 修改服务器相关配置
在服务器上，需要修改以下配置：

#### 修改IP地址配置
编辑 `docker-compose.yml` 和 `conf/config.docker.ini`，将 `8.137.53.3` 替换为您的服务器IP：

```bash
sed -i 's/8.137.53.3/YOUR_SERVER_IP/g' docker-compose.yml
sed -i 's/8.137.53.3/YOUR_SERVER_IP/g' conf/config.docker.ini
```

#### 修改邮箱配置（可选）
如果您使用不同的邮箱服务，请修改 `conf/config.docker.ini` 中的邮箱配置：
```ini
[email]
SmtpHost = your-smtp-server.com
SmtpPort = 465
SmptEmail = your-email@domain.com
SmptPass = your-app-password
```

### 3. 启动服务
```bash
# 进入项目目录
cd gin_mall_tmp

# 使用docker-compose启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### 4. 验证部署

#### 检查服务状态
- 前端服务：访问 `http://YOUR_SERVER_IP`
- 后端API：访问 `http://YOUR_SERVER_IP:3000/api/v1/health`
- MySQL数据库：端口 3308
- Redis缓存：端口 6379

#### 测试功能
1. 注册功能：测试邮箱验证码注册
2. 登录功能：验证用户登录
3. 商品浏览：检查商品列表和详情页
4. 购物车：测试添加商品到购物车

## 重要注意事项

### 安全配置
1. **修改默认密码**：
   - MySQL root密码：123456 → 改为强密码
   - Redis密码：123456 → 改为强密码

2. **防火墙配置**：
   - 开放端口：80（前端）、3000（后端）、3308（MySQL）、6379（Redis）
   - 建议使用防火墙限制访问IP

### 性能优化
1. **数据库优化**：
   - 根据服务器配置调整MySQL内存设置
   - 考虑使用数据库连接池

2. **缓存策略**：
   - Redis配置持久化
   - 设置合适的缓存过期时间

### 监控和维护
1. **日志管理**：
   - 配置日志轮转
   - 监控错误日志

2. **备份策略**：
   - 定期备份数据库
   - 备份重要配置文件

## 故障排除

### 常见问题
1. **端口冲突**：检查端口是否被其他服务占用
2. **数据库连接失败**：验证MySQL服务状态和密码
3. **Redis连接失败**：检查Redis密码和网络配置
4. **邮箱发送失败**：验证SMTP配置和防火墙

### 日志查看
```bash
# 查看所有服务日志
docker-compose logs

# 查看特定服务日志
docker-compose logs app
docker-compose logs frontend
docker-compose logs mysql
docker-compose logs redis
```

## 更新部署

当需要更新代码时：
```bash
# 停止服务
docker-compose down

# 拉取最新代码
# 重新构建镜像
docker-compose build --no-cache

# 启动服务
docker-compose up -d
```

---

**注意**：部署前请确保服务器已安装Docker和Docker Compose，并且有足够的资源运行所有服务。