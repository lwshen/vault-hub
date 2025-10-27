#!/bin/bash
set -eu  # Exit on error, undefined variables

# Build script for Echo migration testing
# This script builds both Fiber and Echo versions for comparison

echo "🚀 Building VaultHub with Echo migration..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Build Echo version
echo -e "${YELLOW}Building Echo server...${NC}"
go build -o tmp/main-echo ./apps/server/main_echo.go
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Echo server built successfully${NC}"
else
    echo -e "${RED}❌ Echo server build failed${NC}"
    exit 1
fi

# Build original Fiber version for comparison
echo -e "${YELLOW}Building original Fiber server...${NC}"
go build -o tmp/main-fiber ./apps/server/main.go
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Fiber server built successfully${NC}"
else
    echo -e "${RED}❌ Fiber server build failed${NC}"
    exit 1
fi

# Run tests for Echo implementation
echo -e "${YELLOW}Running Echo tests...${NC}"
go test ./packages/api -v -run TestEcho
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Echo tests passed${NC}"
else
    echo -e "${RED}❌ Echo tests failed${NC}"
    TEST_FAILED=true
fi

# Generate documentation
echo -e "${YELLOW}Generating migration documentation...${NC}"
echo "Phase 3 Migration Status:" > MIGRATION_STATUS.md
echo "======================" >> MIGRATION_STATUS.md
echo "" >> MIGRATION_STATUS.md
echo "✅ Dependencies: Echo v4.13.4 added" >> MIGRATION_STATUS.md
echo "✅ Server Bootstrap: main_echo.go created" >> MIGRATION_STATUS.md
echo "✅ Authentication: Echo middleware implemented" >> MIGRATION_STATUS.md
echo "✅ Routing: Echo routes configured" >> MIGRATION_STATUS.md
echo "✅ Static Assets: SPA serving configured" >> MIGRATION_STATUS.md
echo "✅ Model Adapter: Compatibility layer created" >> MIGRATION_STATUS.md

if [ "${TEST_FAILED:-false}" = "true" ]; then
    echo "⚠️  Tests: Some tests failed" >> MIGRATION_STATUS.md
else
    echo "✅ Tests: All tests passed" >> MIGRATION_STATUS.md
fi

echo "" >> MIGRATION_STATUS.md
echo "Next Steps:" >> MIGRATION_STATUS.md
echo "1. Run migration tests: go test ./packages/api -v" >> MIGRATION_STATUS.md
echo "2. Start Echo server: ./tmp/main-echo" >> MIGRATION_STATUS.md
echo "3. Test with existing client" >> MIGRATION_STATUS.md
echo "4. Compare with Fiber version: ./tmp/main-fiber" >> MIGRATION_STATUS.md

echo -e "${GREEN}📄 Migration status saved to MIGRATION_STATUS.md${NC}"

echo -e "${GREEN}🎉 Phase 3 migration preparation complete!${NC}"
echo ""
echo "Available commands:"
echo "  ./tmp/main-echo     - Start Echo server"
echo "  ./tmp/main-fiber    - Start original Fiber server"
echo "  go test ./packages/api -v -run TestEcho  - Run Echo tests"
echo ""