#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🧹 Cleaning up test data...${NC}"

# Remove SQLite database to start fresh
if [ -f "app.db" ]; then
    rm app.db
    echo -e "${GREEN}✅ Removed existing database${NC}"
else
    echo -e "${YELLOW}ℹ️  No existing database found${NC}"
fi

# Remove temporary files
rm -f /tmp/response.json

echo -e "${GREEN}🎉 Cleanup complete!${NC}"
echo -e "${YELLOW}💡 You can now run ./test_apis.sh for fresh testing${NC}"
