# 🚀 AI Campaign Analyzer — Backend

## 📌 Overview
Backend service for AI Campaign Analyzer built with:
- Go (Golang)
- Gin Framework
- Clean Architecture
- PostgreSQL (Supabase)
- Google Gemini AI

---

## 🏗️ Architecture
Handler → Usecase → Repository → Infrastructure

---

## ⚙️ Tech Stack
- Gin (HTTP framework)
- pgx (PostgreSQL driver)
- Supabase (DB & Auth)
- Gemini AI (genai SDK)

---

## 📂 Project Structure
```
/cmd/server
/internal
  /domain
  /usecase
  /repository
  /delivery/http
  /infrastructure
/pkg/middleware
/config
```

---

## 🔑 Features
- Campaign CRUD
- Metrics calculation (CTR, CPC, CPA)
- AI Analysis (Gemini)
- Analysis history
- Pagination support
- Consistent API response (Presenter)

---

## 🔌 API Endpoints

### Campaign
- GET /campaigns
- POST /campaigns
- GET /campaigns/:id/metrics
- GET /campaigns/:id/analyze

### Analysis
- GET /analyses
- GET /analyses/:id

---

## ▶️ Run Project

### 1. Setup Environment
Create `.env`:

```
DATABASE_URL=your_db_url
GEMINI_API_KEY=your_api_key
```

### 2. Run Server
```
go run cmd/server/main.go
```

Server runs at:
```
http://localhost:8080
```

---

## 🧠 Key Concepts
- Clean Architecture separation
- AI abstraction via interface
- Presenter pattern for API consistency
- Safe metrics calculation
- Pagination & scalable queries

---

## 📌 Status
Backend is production-ready with AI integration and scalable design.
