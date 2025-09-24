#!/bin/bash

# Build script for VaultHub with embedded web assets
set -e

echo "🔨 VaultHub Build Script"
echo "========================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if pnpm is installed
if ! command -v pnpm &> /dev/null; then
    echo -e "${RED}❌ pnpm is not installed. Please install it first.${NC}"
    echo "   Run: npm install -g pnpm"
    exit 1
fi

# Check if go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed. Please install it first.${NC}"
    exit 1
fi

# Clean previous builds
echo -e "${YELLOW}🧹 Cleaning previous builds...${NC}"
rm -rf apps/web/dist
rm -rf internal/embed/dist
rm -f vault-hub-server vault-hub-cli vault-hub-cron

# Build web application
echo -e "${YELLOW}📦 Building web application...${NC}"
cd apps/web
if [ ! -d "node_modules" ]; then
    echo "   Installing web dependencies..."
    pnpm install --silent
fi
pnpm build
cd ../..

# Copy dist to embed directory
echo -e "${YELLOW}📂 Copying dist to embed directory...${NC}"
cp -r apps/web/dist internal/embed/

# Build Go binaries
echo -e "${YELLOW}🚀 Building Go binaries...${NC}"

echo "   Building server with embedded web assets..."
go build -o vault-hub-server ./apps/server

echo "   Building CLI..."
go build -o vault-hub-cli ./apps/cli

echo "   Building cron..."
go build -o vault-hub-cron ./apps/cron

# Check binary sizes
echo -e "${GREEN}✅ Build completed successfully!${NC}"
echo ""
echo "Binary sizes:"
ls -lh vault-hub-* | awk '{print "  " $9 ": " $5}'
echo ""
echo "Embedded web assets size:"
du -sh internal/embed/dist | awk '{print "  " $2 ": " $1}'

echo ""
echo -e "${GREEN}🎉 All binaries built with embedded web assets!${NC}"
echo ""
echo "To run the server:"
echo "  ./vault-hub-server"
echo ""
echo "Note: Make sure to set required environment variables:"
echo "  - JWT_SECRET"
echo "  - ENCRYPTION_KEY"
echo "  - DATABASE_URL (optional, defaults to 'data.db')"
echo "  - APP_PORT (optional, defaults to '3000')"