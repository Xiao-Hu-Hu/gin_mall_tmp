# 部署说明

## 部署前准备

### 1. 修改配置文件

在部署到阿里云服务器之前，需要修改以下配置文件中的服务器IP地址：

#### 修改 `conf/config.docker.ini`

配置文件中的IP地址已设置为：`8.137.53.3`

```ini
[email]
ValidEmail = http://8.137.53.3/#/valid/email/

[path]
Host = http://8.137.53.3
```

#### 修改 `docker-compose.yml`

环境变量中的IP地址已设置为：`8.137.53.3`

```yaml
environment:
  - HOST_URL=http://8.137.53.3:3000
```

### 2. 服务器要求

- 操作系统: Linux (推荐 Ubuntu 20.04+ 或 CentOS 7+)
- Docker: 20.10+
- Docker Compose: 1.29+
- 内存: 至少 2GB
- 磁盘: 至少 10GB 可用空间

### 3. 安装 Docker 和 Docker Compose

如果服务器上还没有安装 Docker，请执行以下命令：

```bash
# 安装 Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 安装 Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

## 部署步骤

### 1. 上传项目到服务器

将整个项目目录上传到服务器，可以使用 `scp` 或 `git clone`：

```bash
# 使用 scp 上传
scp -r /path/to/gin_mall_tmp root@8.137.53.3:/opt/

# 或使用 git
git clone YOUR_REPO_URL /opt/gin_mall_tmp
```

### 2. 进入项目目录

```bash
cd /opt/gin_mall_tmp
```

### 3. 修改配置文件

按照上面的说明修改 `conf/config.docker.ini` 和 `docker-compose.yml` 中的IP地址。

### 4. 构建和启动服务

```bash
# 构建并启动所有服务
docker-compose up -d --build

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### 5. 检查服务

- 前端访问: http://8.137.53.3
- 后端API: http://8.137.53.3:3000
- 数据库端口: 8.137.53.3:3308
- Redis端口: 8.137.53.3:6379

## 常用命令

```bash
# 停止所有服务
docker-compose down

# 重启服务
docker-compose restart

# 查看日志
docker-compose logs -f app        # 后端日志
docker-compose logs -f frontend   # 前端日志
docker-compose logs -f mysql      # 数据库日志

# 进入容器
docker-compose exec app sh        # 进入后端容器
docker-compose exec mysql mysql -uroot -p123456 mall_db_tmp  # 进入数据库

# 清理数据（谨慎使用）
docker-compose down -v
```

## 防火墙配置

确保服务器防火墙开放以下端口：

```bash
# Ubuntu/Debian
sudo ufw allow 80/tcp
sudo ufw allow 3000/tcp
sudo ufw allow 3308/tcp
sudo ufw allow 6379/tcp

# CentOS
sudo firewall-cmd --permanent --add-port=80/tcp
sudo firewall-cmd --permanent --add-port=3000/tcp
sudo firewall-cmd --permanent --add-port=3308/tcp
sudo firewall-cmd --permanent --add-port=6379/tcp
sudo firewall-cmd --reload
```

## 阿里云安全组配置

在阿里云控制台配置安全组规则，开放以下端口：

- 80 (HTTP)
- 3000 (后端API)
- 3308 (MySQL，建议仅内网访问)
- 6379 (Redis，建议仅内网访问)

## 故障排查

### 1. 服务无法启动

```bash
# 查看详细错误日志
docker-compose logs app
docker-compose logs frontend
```

### 2. 数据库连接失败

检查 `docker-compose.yml` 中的数据库配置和环境变量是否正确。

### 3. 前端无法访问后端API

检查 nginx 配置中的代理设置，确保 `proxy_pass` 指向正确的后端地址。

### 4. 静态资源无法访问

确保 `static` 目录已正确挂载到容器中。

## 数据备份

```bash
# 备份数据库
docker-compose exec mysql mysqldump -uroot -p123456 mall_db_tmp > backup.sql

# 恢复数据库
docker-compose exec -T mysql mysql -uroot -p123456 mall_db_tmp < backup.sql
```

## 更新部署

```bash
# 拉取最新代码
git pull

# 重新构建并启动
docker-compose up -d --build
```

