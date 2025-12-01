# Quick Start Guide

## Automated Setup (Recommended)

Run the setup script:
```bash
./setup.sh
```

This will:
- Check prerequisites
- Create the database
- Configure environment variables
- Install dependencies
- Build the backend

## Manual Setup

### 1. Database Setup
```bash
createdb minflow
```

### 2. Backend Setup
```bash
cd server
cp .env.example .env
# Edit .env and update DB_PASSWORD and JWT_SECRET
go mod download
cd cmd
go run main.go
```

### 3. Frontend Setup
```bash
cd client
npm install
npm run dev
```

## Access the Application

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080/api

## First Steps

1. **Sign Up**: Create your account at `/auth/signup`
2. **Login**: Access your dashboard
3. **Add Expense**: Click "Add Expense" to record your first expense
4. **View Analytics**: Click "View Analytics" to see spending patterns

## Creating an Admin Account

After signing up, connect to PostgreSQL:
```bash
psql -U postgres -d minflow
```

Then run:
```sql
UPDATE users SET is_admin = true WHERE email = 'your@email.com';
```

## Troubleshooting

### Backend won't start
- Check PostgreSQL is running: `sudo systemctl status postgresql`
- Verify database credentials in `server/.env`
- Check port 8080 is not in use: `lsof -i :8080`

### Frontend won't start
- Clear Next.js cache: `rm -rf client/.next`
- Reinstall dependencies: `cd client && rm -rf node_modules package-lock.json && npm install`
- Check port 3000 is not in use: `lsof -i :3000`

### Cannot login/signup
- Check backend is running and accessible
- Verify CORS settings in backend
- Check browser console for errors
- Ensure `NEXT_PUBLIC_API_URL` is set correctly in `client/.env.local`

## API Testing

Test the backend with curl:

```bash
# Health check (add this endpoint if needed)
curl http://localhost:8080/api

# Signup
curl -X POST http://localhost:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

## Development Tips

- Backend changes require restart (or use `air` for hot reload)
- Frontend has hot reload built-in
- Check backend logs for API errors
- Use browser DevTools Network tab to debug API calls
- JWT token is stored in localStorage with key `minflow_token`

## Production Deployment

Before deploying to production:

1. Update `JWT_SECRET` to a strong random value
2. Set `DB_PASSWORD` to a secure password
3. Update `ALLOWED_ORIGINS` to your production domain
4. Enable PostgreSQL SSL: `DB_SSLMODE=require`
5. Build frontend: `cd client && npm run build`
6. Use a process manager for backend (systemd, PM2, etc.)
7. Set up HTTPS with a reverse proxy (nginx, Caddy)

## Need Help?

Check the full README.md for detailed information about:
- Architecture and project structure
- API endpoints documentation
- Feature descriptions
- Technology stack details
