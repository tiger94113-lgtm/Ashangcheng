#!/bin/bash

# Asset Mall Backend 部署脚本
# 支持 Render、Railway、Fly.io

echo "🚀 Asset Mall Backend 部署工具"
echo "=============================="

# 检查命令
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# 菜单
show_menu() {
    echo ""
    echo "选择部署平台:"
    echo "1) Render (推荐，简单)"
    echo "2) Railway (推荐，稳定)"
    echo "3) Fly.io (专业，稳定)"
    echo "4) 本地测试"
    echo "5) 退出"
    echo ""
}

# 部署到 Render
deploy_render() {
    echo "📦 部署到 Render..."
    
    if ! command_exists git; then
        echo "❌ 请先安装 Git"
        exit 1
    fi
    
    # 检查是否有 git 仓库
    if [ ! -d .git ]; then
        echo "📝 初始化 Git 仓库..."
        git init
        git add .
        git commit -m "Initial commit"
    fi
    
    echo ""
    echo "请按以下步骤操作:"
    echo "1. 在 GitHub 创建新仓库"
    echo "2. 运行以下命令推送代码:"
    echo ""
    echo "   git remote add origin https://github.com/你的用户名/仓库名.git"
    echo "   git branch -M main"
    echo "   git push -u origin main"
    echo ""
    echo "3. 访问 https://dashboard.render.com"
    echo "4. 点击 'New +' → 'Web Service'"
    echo "5. 选择你的 GitHub 仓库"
    echo "6. 点击 'Create Web Service'"
    echo ""
    echo "✅ Render 会自动检测配置并部署"
}

# 部署到 Railway
deploy_railway() {
    echo "📦 部署到 Railway..."
    
    if ! command_exists npm; then
        echo "⚠️  建议安装 Railway CLI: npm install -g @railway/cli"
    fi
    
    if command_exists railway; then
        echo "🚀 使用 Railway CLI 部署..."
        railway login
        railway init
        railway up
    else
        echo ""
        echo "请按以下步骤操作:"
        echo "1. 访问 https://railway.app"
        echo "2. 用 GitHub 登录"
        echo "3. 点击 'New Project'"
        echo "4. 选择 'Deploy from GitHub repo'"
        echo "5. 选择你的仓库"
        echo ""
        echo "✅ Railway 会自动部署"
    fi
}

# 部署到 Fly.io
deploy_fly() {
    echo "📦 部署到 Fly.io..."
    
    if ! command_exists fly; then
        echo "📥 安装 Fly CLI..."
        curl -L https://fly.io/install.sh | sh
        export PATH="$HOME/.fly/bin:$PATH"
    fi
    
    echo "🔐 登录 Fly.io..."
    fly auth login
    
    echo "🚀 部署应用..."
    fly launch --name asset-mall-backend --region hkg --no-deploy
    fly deploy
    
    echo ""
    echo "✅ 部署完成!"
    echo "访问: https://asset-mall-backend.fly.dev"
}

# 本地测试
local_test() {
    echo "🧪 本地测试..."
    
    if ! command_exists go; then
        echo "❌ 请先安装 Go: https://golang.org/dl/"
        exit 1
    fi
    
    echo "📥 下载依赖..."
    go mod download
    
    echo "🚀 启动服务..."
    echo "服务将在 http://localhost:8080 运行"
    echo "按 Ctrl+C 停止"
    echo ""
    
    go run main.go
}

# 主程序
main() {
    show_menu
    read -p "请选择 (1-5): " choice
    
    case $choice in
        1)
            deploy_render
            ;;
        2)
            deploy_railway
            ;;
        3)
            deploy_fly
            ;;
        4)
            local_test
            ;;
        5)
            echo "👋 再见!"
            exit 0
            ;;
        *)
            echo "❌ 无效选择"
            exit 1
            ;;
    esac
}

main
