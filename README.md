# School Management SaaS

Sistem SaaS Manajemen Sekolah - Platform multi-tenant untuk pendataan sekolah, absensi RFID, nilai akademik, dan sistem BK.

## Tech Stack

### Frontend (Web Admin)
- Vue 3 + TypeScript
- Ant Design Vue
- Pinia (State Management)
- Vue Router

### Backend
- Golang + Fiber Framework
- GORM (ORM)
- PostgreSQL (Database)
- Redis (Queue & Cache)
- Firebase Cloud Messaging (Push Notifications)

## Development Setup

### Prerequisites
- Docker & Docker Compose
- Node.js 18+
- Go 1.21+

### 1. Start Database Services

```bash
# Start PostgreSQL, Redis, and Adminer
docker compose up -d

# Check status
docker compose ps

# View logs
docker compose logs -f
```

**Services:**
| Service | Port | Description |
|---------|------|-------------|
| PostgreSQL | 5432 | Main database |
| Redis | 6379 | Queue & cache |
| Adminer | 8080 | Database UI (http://localhost:8080) |

**Database Credentials:**
- Host: `localhost`
- Port: `5432`
- Database: `school_management`
- Username: `school_admin`
- Password: `school_secret_2024`

### 2. Frontend Development

```bash
cd web-admin
npm install
npm run dev
```

Frontend akan berjalan di http://localhost:5173

### 3. Backend Development

```bash
cd backend
go mod download
go run cmd/server/main.go
```

Backend API akan berjalan di http://localhost:3000

## Project Structure

```
├── web-admin/          # Vue 3 Frontend
│   ├── src/
│   │   ├── components/ # Reusable components
│   │   ├── views/      # Page components
│   │   ├── stores/     # Pinia stores
│   │   ├── services/   # API services
│   │   └── types/      # TypeScript types
│   └── ...
├── backend/            # Golang Backend
│   ├── cmd/server/     # Entry point
│   ├── internal/       # Internal packages
│   │   ├── modules/    # Domain modules
│   │   ├── middleware/ # HTTP middleware
│   │   └── shared/     # Shared utilities
│   └── ...
├── docker-compose.yml  # Development services
└── .env               # Environment variables
```

## User Roles

| Role | Description |
|------|-------------|
| Super Admin | Platform administrator, manages all tenants |
| Admin Sekolah | School administrator, manages school data |
| Guru BK | Counseling teacher, manages BK records |
| Wali Kelas | Homeroom teacher, manages class & grades |

## Useful Commands

```bash
# Stop all services
docker compose down

# Stop and remove volumes (reset data)
docker compose down -v

# Connect to PostgreSQL
docker exec -it school-management-postgres psql -U school_admin -d school_management

# Connect to Redis CLI
docker exec -it school-management-redis redis-cli
```
