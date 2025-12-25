# 🏥 Health Tracker - Live for Health

Aplikasi web untuk tracking kesehatan dengan fitur rekomendasi makanan, tracking aktivitas, dan manajemen keluarga.

---

## 📋 Prasyarat (Harus Diinstall)

Pastikan komputer Anda sudah menginstall:

| Software | Versi Minimum | Link Download |
|----------|---------------|---------------|
| **Node.js** | v18+ | [nodejs.org](https://nodejs.org) |
| **Go** | v1.21+ | [go.dev/dl](https://go.dev/dl/) |
| **Git** | Terbaru | [git-scm.com](https://git-scm.com/) |
| **Cloudflared** | Opsional | [Cloudflare Tunnel](https://developers.cloudflare.com/cloudflare-one/connections/connect-apps/install-and-setup/installation/) |

---

## 🚀 Cara Menjalankan

### 1. Clone/Download Project

```powershell
# Jika via Git
git clone <URL_REPOSITORY>
cd my-react-app
```

Atau extract file ZIP ke folder yang diinginkan.

---

### 2. Setup Frontend

```powershell
# Masuk ke folder project
cd my-react-app

# Install dependencies
npm install

# Jalankan development server
npm run dev
```

Frontend akan berjalan di: **http://localhost:5173**

---

### 3. Setup Backend

Buka terminal baru (jangan tutup terminal frontend):

```powershell
# Masuk ke folder backend
cd my-react-app/backend

# Copy file environment
copy .env.example .env

# Download Go dependencies
go mod download

# Jalankan backend server
go run main.go
```

Backend akan berjalan di: **http://localhost:8080**

---

## 🌐 Akses dari Luar (Opsional)

Untuk membuat web bisa diakses publik via internet:

```powershell
# Terminal 1: Tunnel untuk Backend
cloudflared tunnel --url http://localhost:8080

# Terminal 2: Tunnel untuk Frontend
cloudflared tunnel --url http://localhost:5173
```

URL publik akan muncul di terminal dengan format:
```
https://xxxxx-xxxxx.trycloudflare.com
```

---

## 📁 Struktur Project

```
my-react-app/
├── src/                    # Frontend React
│   ├── components/         # Komponen UI
│   ├── pages/              # Halaman aplikasi
│   ├── contexts/           # React Context
│   └── config/             # Konfigurasi API
├── backend/                # Backend Go
│   ├── handlers/           # API Handlers
│   ├── models/             # Database Models
│   ├── middleware/         # Auth Middleware
│   ├── routes/             # Route definitions
│   └── database/           # Database config
├── package.json            # Frontend dependencies
└── vite.config.js          # Vite configuration
```

---

## 🔧 Troubleshooting

### Error: "npm not found"
→ Install Node.js dari [nodejs.org](https://nodejs.org)

### Error: "go not found"
→ Install Go dari [go.dev](https://go.dev/dl/)

### Error: Port sudah digunakan
```powershell
# Cek port 8080
netstat -ano | findstr :8080

# Kill process (ganti PID dengan ID dari command di atas)
taskkill /PID <PID> /F
```

### Database error
→ Hapus file `backend/health_tracker.db` dan jalankan ulang backend

---

## 📞 Kontak

Jika ada pertanyaan, hubungi developer.

---

**Happy Coding! 🎉**
