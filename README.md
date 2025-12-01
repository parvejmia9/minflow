# MinFlow - Expense Tracker Application

A full-stack expense tracking SaaS application built with Go (Fiber + GORM) backend and Next.js frontend.

## Features

- **Authentication**: JWT-based authentication with bcrypt password hashing
- **User Management**: Admin panel to manage users
- **Expense Tracking**: Add, view, and delete expenses with categories
- **Analytics**: View spending patterns with date range filters and charts
- **Categories**: Create and manage expense categories
- **Role-Based Access**: Admin and regular user roles

## Tech Stack

### Backend
- **Go 1.24.0**: Programming language
- **Fiber v2.52.10**: Web framework
- **GORM v1.31.1**: ORM
- **PostgreSQL**: Database
- **JWT**: Authentication
- **Bcrypt**: Password hashing

### Frontend
- **Next.js 16.0.6**: React framework
- **TypeScript**: Type safety
- **Tailwind CSS**: Styling
- **Zustand**: State management
- **Recharts**: Data visualization
- **Axios**: HTTP client

## Architecture

The backend follows Clean Architecture principles with:
- **Models**: Database entities
- **Services**: Business logic with dependency injection
- **Handlers**: HTTP request handlers
- **Routes**: API route definitions
- **Middleware**: Authentication and authorization

## Setup Instructions

### Prerequisites
- Go 1.24.0 or higher
- Node.js 20.9.0 or higher
- PostgreSQL

### Database Setup

1. Create PostgreSQL database:
```bash
createdb minflow
```

2. Update database credentials in `server/.env`

### Backend Setup

1. Navigate to server directory:
```bash
cd server
```

2. Copy environment file:
```bash
cp .env.example .env
```

3. Edit `.env` and update the following:
   - `DB_PASSWORD`: Your PostgreSQL password
   - `JWT_SECRET`: A strong secret key for JWT signing
   - Other database settings if needed

4. Install dependencies (already done via go.mod):
```bash
go mod download
```

5. Run the server:
```bash
cd cmd
go run main.go
```

The server will start on `http://localhost:8080` and automatically create database tables on first run.

### Frontend Setup

1. Navigate to client directory:
```bash
cd client
```

2. The `.env.local` file is already created with default values

3. Install dependencies:
```bash
npm install
```

4. Run the development server:
```bash
npm run dev
```

The frontend will start on `http://localhost:3000`

## API Endpoints

### Authentication
- `POST /api/auth/signup` - Register new user
- `POST /api/auth/login` - Login user

### Expenses
- `GET /api/expenses` - Get all expenses for logged-in user
- `GET /api/expenses/:id` - Get single expense
- `POST /api/expenses` - Create new expense
- `DELETE /api/expenses/:id` - Delete expense
- `POST /api/expenses/date-range` - Get expenses by date range
- `POST /api/expenses/analytics` - Get analytics data

### Categories
- `GET /api/categories` - Get all categories for logged-in user
- `GET /api/categories/:id` - Get single category
- `POST /api/categories` - Create new category
- `PUT /api/categories/:id` - Update category
- `DELETE /api/categories/:id` - Delete category

### Users (Admin Only)
- `GET /api/users` - Get all users
- `GET /api/users/:id` - Get single user
- `DELETE /api/users/:id` - Delete user

## Usage

1. **Sign Up**: Create a new account at `/auth/signup`
2. **Login**: Login at `/auth/login`
3. **Dashboard**: View the main dashboard with quick actions
4. **Add Expense**: Record expenses with name, category, unit, per unit cost
5. **View Expenses**: See all your expenses in a table view
6. **Analytics**: View spending patterns with charts and date filters
7. **Admin Panel**: Manage users (admin accounts only)

## Default Admin Account

Create an admin account by signing up, then manually update the `is_admin` field in the database:

```sql
UPDATE users SET is_admin = true WHERE email = 'admin@example.com';
```

## Project Structure

```
server/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── models/              # Database models
│   ├── services/            # Business logic
│   ├── handlers/            # HTTP handlers
│   ├── routes/              # Route definitions
│   ├── middleware/          # Middleware functions
│   └── db/                  # Database connection
├── go.mod
└── .env

client/
├── app/
│   ├── auth/               # Authentication pages
│   ├── dashboard/          # Dashboard page
│   ├── expenses/           # Expense pages
│   ├── analytics/          # Analytics page
│   ├── admin/              # Admin panel
│   └── layout.tsx          # Root layout
├── components/             # React components (if any)
├── lib/                    # Utility functions
├── store/                  # State management
├── types/                  # TypeScript types
└── .env.local
```

## Development Notes

- JWT tokens expire after 7 days
- Passwords are hashed using bcrypt with cost 10
- Total expense is automatically calculated: `total = unit * per_unit_cost`
- Categories are user-specific
- Admin users cannot be deleted from the admin panel
- Date range analytics cannot select future dates or dates before the first expense

## License

MIT
