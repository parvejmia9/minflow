#!/bin/bash

# MinFlow Setup Script
# This script helps set up the entire application

set -e

echo "🚀 MinFlow Expense Tracker Setup"
echo "=================================="
echo ""

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check prerequisites
echo "📋 Checking prerequisites..."

# Check Go version
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed${NC}"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo -e "${GREEN}✓${NC} Go: $GO_VERSION"

# Check Node.js
if ! command -v node &> /dev/null; then
    echo -e "${RED}❌ Node.js is not installed${NC}"
    exit 1
fi

NODE_VERSION=$(node --version)
echo -e "${GREEN}✓${NC} Node.js: $NODE_VERSION"

# Check PostgreSQL
if ! command -v psql &> /dev/null; then
    echo -e "${YELLOW}⚠${NC} PostgreSQL client not found in PATH"
    echo "   Make sure PostgreSQL is installed and accessible"
else
    PG_VERSION=$(psql --version | awk '{print $3}')
    echo -e "${GREEN}✓${NC} PostgreSQL: $PG_VERSION"
fi

echo ""
echo "🗄️  Setting up database..."

# Ask for database details
read -p "Enter PostgreSQL username [postgres]: " DB_USER
DB_USER=${DB_USER:-postgres}

read -p "Enter PostgreSQL password [postgres]: " DB_PASSWORD
DB_PASSWORD=${DB_PASSWORD:-postgres}

read -p "Enter database name [minflow]: " DB_NAME
DB_NAME=${DB_NAME:-minflow}

# Create database if it doesn't exist
echo "Creating database '$DB_NAME'..."
PGPASSWORD=$DB_PASSWORD createdb -U $DB_USER $DB_NAME 2>/dev/null && echo -e "${GREEN}✓${NC} Database created" || echo -e "${YELLOW}⚠${NC} Database may already exist"

echo ""
echo "⚙️  Configuring backend..."

# Update server .env
cd server
if [ ! -f .env ]; then
    cp .env.example .env
fi

# Update .env with user inputs
sed -i "s/DB_USER=.*/DB_USER=$DB_USER/" .env
sed -i "s/DB_PASSWORD=.*/DB_PASSWORD=$DB_PASSWORD/" .env
sed -i "s/DB_NAME=.*/DB_NAME=$DB_NAME/" .env

# Generate random JWT secret
JWT_SECRET=$(openssl rand -base64 32)
sed -i "s/JWT_SECRET=.*/JWT_SECRET=$JWT_SECRET/" .env

echo -e "${GREEN}✓${NC} Backend configuration updated"

# Install Go dependencies
echo "Installing Go dependencies..."
go mod download
echo -e "${GREEN}✓${NC} Go dependencies installed"

# Build backend
echo "Building backend..."
go build -o bin/server ./cmd/main.go
echo -e "${GREEN}✓${NC} Backend built successfully"

cd ..

echo ""
echo "🎨 Configuring frontend..."

cd client

# Frontend already has .env.local created
echo -e "${GREEN}✓${NC} Frontend configuration ready"

# Install Node dependencies
echo "Installing Node dependencies..."
npm install --silent
echo -e "${GREEN}✓${NC} Node dependencies installed"

cd ..

echo ""
echo -e "${GREEN}✅ Setup complete!${NC}"
echo ""
echo "📚 Next steps:"
echo ""
echo "1. Start the backend server:"
echo -e "   ${YELLOW}cd server/cmd && go run main.go${NC}"
echo ""
echo "2. In a new terminal, start the frontend:"
echo -e "   ${YELLOW}cd client && npm run dev${NC}"
echo ""
echo "3. Open your browser to: http://localhost:3000"
echo ""
echo "4. Create your first account and start tracking expenses!"
echo ""
echo "💡 Tips:"
echo "   - The database tables will be created automatically on first run"
echo "   - To create an admin account, sign up then run:"
echo "     UPDATE users SET is_admin = true WHERE email = 'your@email.com';"
echo ""
echo "📖 For more information, see README.md"
