# Docker 部署详细指南

## 📋 前置条件

- ✅ 服务器已安装 Docker（版本 20.10+）
- ✅ 服务器已安装 Docker Compose（版本 1.29+）
- ✅ 服务器已开放端口：80, 3000, 3308, 6379
- ✅ 服务器IP：8.137.53.3

## 🚀 部署步骤

### 第一步：上传项目到服务器

#### 方法一：使用 Git（推荐）

```bash
# 在服务器上执行
cd /opt
git clone YOUR_REPO_URL gin_mall_tmp
cd gin_mall_tmp
```

#### 方法二：使用 SCP 上传

```bash
# 在本地电脑执行
scp -r /path/to/gin_mall_tmp root@8.137.53.3:/opt/
```

#### 方法三：使用压缩包上传

```bash
# 在本地压缩项目
tar -czf gin_mall_tmp.tar.gz gin_mall_tmp/

# 上传到服务器
scp gin_mall_tmp.tar.gz root@8.137.53.3:/opt/

# 在服务器上解压
ssh root@8.137.53.3
cd /opt
tar -xzf gin_mall_tmp.tar.gz
cd gin_mall_tmp
```

---

### 第二步：检查项目文件

进入项目目录，确认以下文件存在：

```bash
cd /opt/gin_mall_tmp

# 检查关键文件
ls -la docker-compose.yml
ls -la Dockerfile
ls -la frontend/Dockerfile
ls -la conf/config.docker.ini
```

**应该看到：**
- ✅ `docker-compose.yml` - Docker Compose 配置文件
- ✅ `Dockerfile` - 后端 Dockerfile
- ✅ `frontend/Dockerfile` - 前端 Dockerfile
- ✅ `frontend/nginx.conf` - Nginx 配置文件
- ✅ `conf/config.docker.ini` - Docker 环境配置文件

---

### 第三步：检查 Docker 和 Docker Compose

```bash
# 检查 Docker 版本
docker --version
# 应该显示：Docker version 20.10.x 或更高

# 检查 Docker Compose 版本
docker-compose --version
# 应该显示：docker-compose version 1.29.x 或更高

# 如果未安装 Docker Compose，执行：
sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

---

### 第四步：检查端口占用

```bash
# 检查端口是否被占用
netstat -tulpn | grep -E ':(80|3000|3308|6379)'

# 或者使用 ss 命令
ss -tulpn | grep -E ':(80|3000|3308|6379)'

# 如果端口被占用，需要停止占用端口的服务或修改 docker-compose.yml 中的端口映射
```

---

### 第五步：配置防火墙（如果需要）

#### Ubuntu/Debian 系统：

```bash
# 检查防火墙状态
sudo ufw status

# 如果防火墙开启，开放端口
sudo ufw allow 80/tcp
sudo ufw allow 3000/tcp
sudo ufw allow 3308/tcp
sudo ufw allow 6379/tcp

# 重新加载防火墙
sudo ufw reload
```

#### CentOS/RHEL 系统：

```bash
# 检查防火墙状态
sudo firewall-cmd --state

# 如果防火墙开启，开放端口
sudo firewall-cmd --permanent --add-port=80/tcp
sudo firewall-cmd --permanent --add-port=3000/tcp
sudo firewall-cmd --permanent --add-port=3308/tcp
sudo firewall-cmd --permanent --add-port=6379/tcp

# 重新加载防火墙
sudo firewall-cmd --reload
```

---

### 第六步：构建和启动服务

```bash
# 确保在项目根目录
cd /opt/gin_mall_tmp

# 构建并启动所有服务（首次运行会下载镜像，需要一些时间）
docker-compose up -d --build

# 查看构建和启动日志
docker-compose logs -f
```

**预期输出：**
- ✅ 开始构建后端镜像
- ✅ 开始构建前端镜像
- ✅ 启动 MySQL 容器
- ✅ 启动 Redis 容器
- ✅ 启动后端应用容器
- ✅ 启动前端 Nginx 容器

**首次构建可能需要 5-10 分钟**，请耐心等待。

---

### 第七步：检查服务状态

```bash
# 查看所有容器状态
docker-compose ps

# 应该看到所有服务都是 "Up" 状态：
# - app (后端)
# - frontend (前端)
# - mysql (数据库)
# - redis (缓存)
```

**如果某个容器状态不是 "Up"，查看日志：**

```bash
# 查看特定服务的日志
docker-compose logs app        # 后端日志
docker-compose logs frontend   # 前端日志
docker-compose logs mysql      # 数据库日志
docker-compose logs redis      # Redis日志

# 查看所有服务的日志
docker-compose logs
```

---

### 第八步：验证部署

#### 1. 检查容器是否运行

```bash
docker ps

# 应该看到 4 个容器在运行
```

#### 2. 检查前端访问

在浏览器中访问：**http://8.137.53.3**

应该能看到商城首页。

#### 3. 检查后端 API

在浏览器中访问：**http://8.137.53.3:3000/api/v1/ping**

应该返回：`"success"`

#### 4. 检查数据库连接

```bash
# 进入 MySQL 容器
docker-compose exec mysql mysql -uroot -p123456

# 在 MySQL 中执行
USE mall_db_tmp;
SHOW TABLES;
EXIT;
```

#### 5. 检查 Redis 连接

```bash
# 进入 Redis 容器
docker-compose exec redis redis-cli -a 123456

# 在 Redis 中执行
PING
# 应该返回：PONG
EXIT
```

---

## 🔧 常用操作命令

### 查看日志

```bash
# 实时查看所有日志
docker-compose logs -f

# 查看特定服务的日志
docker-compose logs -f app
docker-compose logs -f frontend
docker-compose logs -f mysql
docker-compose logs -f redis

# 查看最近 100 行日志
docker-compose logs --tail=100 app
```

### 重启服务

```bash
# 重启所有服务
docker-compose restart

# 重启特定服务
docker-compose restart app
docker-compose restart frontend
```

### 停止服务

```bash
# 停止所有服务（不删除容器）
docker-compose stop

# 停止并删除容器
docker-compose down

# 停止并删除容器和卷（⚠️ 会删除数据）
docker-compose down -v
```

### 更新代码后重新部署

```bash
# 1. 停止服务
docker-compose down

# 2. 拉取最新代码（如果使用 Git）
git pull

# 3. 重新构建并启动
docker-compose up -d --build
```

### 进入容器调试

```bash
# 进入后端容器
docker-compose exec app sh

# 进入前端容器
docker-compose exec frontend sh

# 进入 MySQL 容器
docker-compose exec mysql bash

# 进入 Redis 容器
docker-compose exec redis sh
```

---

## 🐛 常见问题排查

### 问题 1：容器无法启动

**症状：** `docker-compose ps` 显示容器状态为 "Exit" 或 "Restarting"

**排查步骤：**

```bash
# 1. 查看错误日志
docker-compose logs app

# 2. 检查端口是否被占用
netstat -tulpn | grep -E ':(80|3000|3308|6379)'

# 3. 检查磁盘空间
df -h

# 4. 检查 Docker 资源
docker system df
```

**常见原因：**
- 端口被占用 → 修改 `docker-compose.yml` 中的端口映射
- 磁盘空间不足 → 清理 Docker 镜像和容器
- 配置文件错误 → 检查 `conf/config.docker.ini`

---

### 问题 2：前端无法访问后端 API

**症状：** 前端页面可以打开，但无法加载数据

**排查步骤：**

```bash
# 1. 检查后端是否正常运行
curl http://localhost:3000/api/v1/ping

# 2. 检查前端 Nginx 配置
docker-compose exec frontend cat /etc/nginx/conf.d/default.conf

# 3. 检查网络连接
docker network ls
docker network inspect gin_mall_tmp_mall-network
```

**解决方案：**

检查 `frontend/nginx.conf` 中的代理配置是否正确：

```nginx
location /api {
    proxy_pass http://app:3000;  # 确保这里指向 app 服务
    ...
}
```

---

### 问题 3：数据库连接失败

**症状：** 后端日志显示数据库连接错误

**排查步骤：**

```bash
# 1. 检查 MySQL 容器是否运行
docker-compose ps mysql

# 2. 检查 MySQL 日志
docker-compose logs mysql

# 3. 尝试手动连接数据库
docker-compose exec mysql mysql -uroot -p123456 -e "SHOW DATABASES;"
```

**解决方案：**

```bash
# 如果数据库未初始化，等待 MySQL 完全启动（可能需要 30-60 秒）
docker-compose logs -f mysql

# 看到 "ready for connections" 表示 MySQL 已启动
```

---

### 问题 4：静态文件无法访问

**症状：** 商品图片或头像无法显示

**排查步骤：**

```bash
# 1. 检查 static 目录是否挂载
docker-compose exec app ls -la /root/static

# 2. 检查文件权限
ls -la static/

# 3. 检查后端路由配置
docker-compose exec app ls -la /root/static/imgs/product/
```

**解决方案：**

确保 `static` 目录存在且有正确的文件：

```bash
# 在项目根目录
ls -la static/imgs/product/
ls -la static/imgs/avatar/
```

---

### 问题 5：内存不足

**症状：** 容器频繁重启或构建失败

**排查步骤：**

```bash
# 检查系统资源
free -h
df -h

# 检查 Docker 资源使用
docker stats
```

**解决方案：**

```bash
# 清理 Docker 资源
docker system prune -a

# 限制容器资源使用（在 docker-compose.yml 中添加）
services:
  app:
    deploy:
      resources:
        limits:
          memory: 512M
```

---

## 📊 监控和维护

### 查看资源使用情况

```bash
# 实时查看容器资源使用
docker stats

# 查看 Docker 磁盘使用
docker system df
```

### 备份数据库

```bash
# 备份数据库
docker-compose exec mysql mysqldump -uroot -p123456 mall_db_tmp > backup_$(date +%Y%m%d_%H%M%S).sql

# 恢复数据库
docker-compose exec -T mysql mysql -uroot -p123456 mall_db_tmp < backup.sql
```

### 清理日志

```bash
# 清理 Docker 日志（谨慎使用）
docker-compose down
docker system prune -a --volumes
```

---

## ✅ 部署成功检查清单

- [ ] 所有容器状态为 "Up"
- [ ] 前端页面可以访问：http://8.137.53.3
- [ ] 后端 API 可以访问：http://8.137.53.3:3000/api/v1/ping
- [ ] 可以注册新用户
- [ ] 可以登录
- [ ] 可以查看商品列表
- [ ] 可以上架商品
- [ ] 可以添加商品到购物车
- [ ] 可以创建订单
- [ ] 可以支付订单

---

## 🆘 获取帮助

如果遇到问题，请提供以下信息：

1. **错误日志：**
   ```bash
   docker-compose logs > error.log
   ```

2. **容器状态：**
   ```bash
   docker-compose ps
   ```

3. **系统信息：**
   ```bash
   docker --version
   docker-compose --version
   uname -a
   ```

4. **网络配置：**
   ```bash
   docker network inspect gin_mall_tmp_mall-network
   ```

---

## 📝 下一步

部署成功后，您可以：

1. **配置域名**：将域名指向服务器 IP
2. **配置 HTTPS**：使用 Let's Encrypt 配置 SSL 证书
3. **设置自动备份**：配置数据库定时备份
4. **监控告警**：配置服务监控和告警

祝您部署顺利！🎉

