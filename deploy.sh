#!/bin/bash

# 个人商城 Docker 部署脚本
# 使用方法: bash deploy.sh

set -e  # 遇到错误立即退出

echo "=========================================="
echo "  个人商城 Docker 部署脚本"
echo "=========================================="
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查是否在项目根目录
if [ ! -f "docker-compose.yml" ]; then
    echo -e "${RED}错误: 请在项目根目录执行此脚本${NC}"
    exit 1
fi

# 检查 Docker 是否安装
echo -e "${YELLOW}[1/8] 检查 Docker 环境...${NC}"
if ! command -v docker &> /dev/null; then
    echo -e "${RED}错误: Docker 未安装，请先安装 Docker${NC}"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}错误: Docker Compose 未安装，请先安装 Docker Compose${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Docker 版本: $(docker --version)${NC}"
echo -e "${GREEN}✓ Docker Compose 版本: $(docker-compose --version)${NC}"
echo ""

# 检查端口占用
echo -e "${YELLOW}[2/8] 检查端口占用...${NC}"
PORTS=(80 3000 3308 6379)
PORT_OCCUPIED=false

for port in "${PORTS[@]}"; do
    if netstat -tuln 2>/dev/null | grep -q ":$port " || ss -tuln 2>/dev/null | grep -q ":$port "; then
        echo -e "${YELLOW}警告: 端口 $port 已被占用${NC}"
        PORT_OCCUPIED=true
    fi
done

if [ "$PORT_OCCUPIED" = true ]; then
    echo -e "${YELLOW}提示: 如果端口被占用，请先停止相关服务或修改 docker-compose.yml${NC}"
    read -p "是否继续? (y/n): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi
echo ""

# 检查必要文件
echo -e "${YELLOW}[3/8] 检查项目文件...${NC}"
REQUIRED_FILES=(
    "docker-compose.yml"
    "Dockerfile"
    "frontend/Dockerfile"
    "frontend/nginx.conf"
    "conf/config.docker.ini"
)

for file in "${REQUIRED_FILES[@]}"; do
    if [ ! -f "$file" ]; then
        echo -e "${RED}错误: 缺少必要文件: $file${NC}"
        exit 1
    fi
done
echo -e "${GREEN}✓ 所有必要文件存在${NC}"
echo ""

# 停止现有容器
echo -e "${YELLOW}[4/8] 停止现有容器...${NC}"
docker-compose down 2>/dev/null || true
echo -e "${GREEN}✓ 已停止现有容器${NC}"
echo ""

# 清理旧镜像（可选）
read -p "是否清理旧的 Docker 镜像? (y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}清理旧镜像...${NC}"
    docker system prune -f
    echo -e "${GREEN}✓ 清理完成${NC}"
    echo ""
fi

# 构建镜像
echo -e "${YELLOW}[5/8] 构建 Docker 镜像...${NC}"
echo -e "${YELLOW}这可能需要几分钟时间，请耐心等待...${NC}"
docker-compose build --no-cache

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 镜像构建成功${NC}"
else
    echo -e "${RED}✗ 镜像构建失败，请检查错误信息${NC}"
    exit 1
fi
echo ""

# 启动服务
echo -e "${YELLOW}[6/8] 启动服务...${NC}"
docker-compose up -d

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 服务启动成功${NC}"
else
    echo -e "${RED}✗ 服务启动失败，请检查错误信息${NC}"
    exit 1
fi
echo ""

# 等待服务就绪
echo -e "${YELLOW}[7/8] 等待服务就绪...${NC}"
echo -e "${YELLOW}等待 MySQL 启动（最多 60 秒）...${NC}"

for i in {1..60}; do
    if docker-compose exec -T mysql mysqladmin ping -h localhost --silent 2>/dev/null; then
        echo -e "${GREEN}✓ MySQL 已就绪${NC}"
        break
    fi
    if [ $i -eq 60 ]; then
        echo -e "${YELLOW}警告: MySQL 启动超时，但将继续检查其他服务${NC}"
    else
        sleep 1
        echo -n "."
    fi
done
echo ""

# 检查服务状态
echo -e "${YELLOW}[8/8] 检查服务状态...${NC}"
sleep 5

echo ""
echo "=========================================="
echo "  服务状态"
echo "=========================================="
docker-compose ps
echo ""

# 验证服务
echo -e "${YELLOW}验证服务...${NC}"

# 检查后端
if curl -s http://localhost:3000/api/v1/ping > /dev/null; then
    echo -e "${GREEN}✓ 后端服务正常 (http://8.137.53.3:3000)${NC}"
else
    echo -e "${YELLOW}⚠ 后端服务可能未就绪，请稍后重试${NC}"
fi

# 检查前端
if curl -s http://localhost:80 > /dev/null; then
    echo -e "${GREEN}✓ 前端服务正常 (http://8.137.53.3)${NC}"
else
    echo -e "${YELLOW}⚠ 前端服务可能未就绪，请稍后重试${NC}"
fi

echo ""
echo "=========================================="
echo -e "${GREEN}  部署完成！${NC}"
echo "=========================================="
echo ""
echo "访问地址："
echo "  - 前端: http://8.137.53.3"
echo "  - 后端: http://8.137.53.3:3000"
echo ""
echo "常用命令："
echo "  - 查看日志: docker-compose logs -f"
echo "  - 停止服务: docker-compose down"
echo "  - 重启服务: docker-compose restart"
echo "  - 查看状态: docker-compose ps"
echo ""
echo "如果遇到问题，请查看日志："
echo "  docker-compose logs app"
echo "  docker-compose logs frontend"
echo ""

