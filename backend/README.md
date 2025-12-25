# Health Tracker Backend

Backend API untuk aplikasi Health Tracker, dibangun dengan Go + Gin + SQLite.

## Prerequisites

- Go 1.21 atau lebih baru
- Git

## Quick Start

```bash
# Masuk ke folder backend
cd backend

# Download dependencies
go mod tidy

# Jalankan server
go run main.go
```

Server akan berjalan di `http://localhost:8080`

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register user baru
- `POST /api/auth/login` - Login dan dapatkan token
- `GET /api/auth/me` - Get profil user (protected)
- `PUT /api/auth/profile` - Update profil (protected)

### Health Data
- `POST /api/health` - Submit data kesehatan
- `GET /api/health` - Get semua data kesehatan
- `GET /api/health/latest` - Get data terbaru
- `GET /api/health/dashboard` - Get dashboard summary
- `GET /api/health/graph/:period` - Get data grafik (week/month/year)

### Symptoms
- `GET /api/symptoms/list` - Get daftar gejala
- `POST /api/symptoms` - Log gejala
- `POST /api/symptoms/batch` - Log multiple gejala
- `GET /api/symptoms/history` - Get riwayat gejala
- `GET /api/symptoms/stats` - Get statistik gejala

### Family
- `POST /api/family/invite` - Undang anggota keluarga
- `GET /api/family/members` - Get daftar anggota keluarga
- `GET /api/family/requests` - Get permintaan tertunda
- `PUT /api/family/approve/:id` - Setujui permintaan
- `PUT /api/family/reject/:id` - Tolak permintaan
- `GET /api/family/:id/health` - Lihat kesehatan anggota
- `DELETE /api/family/:id` - Hapus anggota

### Recommendations
- `GET /api/recommendations/food` - Rekomendasi makanan
- `GET /api/recommendations/exercise` - Rekomendasi olahraga
- `GET /api/recommendations/emotional` - Rekomendasi aktivitas emosional

## Environment Variables

Buat file `.env` di folder backend:

```env
PORT=8080
GIN_MODE=debug
JWT_SECRET=your-secret-key
JWT_EXPIRY_HOURS=24
DATABASE_PATH=./health_tracker.db
```

## Project Structure

```
backend/
├── main.go              # Entry point
├── go.mod               # Dependencies
├── .env                 # Environment config
├── config/              # Configuration
├── database/            # Database setup
├── models/              # Data models
├── handlers/            # API handlers
├── middleware/          # Auth & CORS
├── routes/              # Route definitions
└── utils/               # Helpers
```
