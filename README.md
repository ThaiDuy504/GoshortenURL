# Go Shortener Service

Dịch vụ rút gọn URL (Go + Gin). Lưu mapping short code → URL trên **Redis**. Ứng dụng vẫn kết nối **PostgreSQL** lúc khởi động (theo `docker-compose`), nên DB phải sẵn sàng khi chạy.

- **Module:** `Go_shortenURL`
- **Go:** `1.25.5` (xem `go.mod`)

## Cấu trúc thư mục

```text
Go_shortenURL/
├── cmd/api/                 # Entrypoint — Gin, routes, LoadHTMLGlob
│   └── main.go
├── internal/
│   ├── handler/             # HTTP: form shorten, redirect
│   ├── service/             # Shorten / resolve
│   ├── repository/          # Redis (URL); Postgres client khởi tạo
│   └── model/
├── pkg/shortener/           # Sinh short code (hash)
├── configs/
│   ├── config.go            # cleanenv đọc configs/config.env
│   └── config.example.env   # Mẫu biến — copy thành config.env
├── templates/               # index.html (Gin html/template)
├── docker-compose.yml       # app + postgres + redis
├── dockerfile               # multi-stage build
├── go.mod
├── go.sum
└── README.md
```

## API & route

| Method | Path | Mô tả |
|--------|------|--------|
| `GET` | `/` | Trang form rút gọn |
| `POST` | `/shorten` | Form field `url` — trả HTML kết quả |
| `GET` | `/:shortCode` | Redirect (logic redirect trong handler có thể đang test) |
| `GET` | `/health` | `{"status":"ok"}` |

Luồng: **handler → service → repository (Redis)**; short code từ **pkg/shortener**.

## Biến môi trường & config

File bắt buộc: **`configs/config.env`** (copy từ `configs/config.example.env`).

- **Cổng HTTP:** `main` đọc `APP_PORT`, mặc định `8080` (khác với field `PORT` trong struct config nếu bạn thêm dùng sau này).
- **Postgres:** `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- **Redis:** `REDIS_ADDR`, `REDIS_PASSWORD`, `REDIS_DB`

Với `docker-compose`, service `app` đã set các biến tương thích Postgres/Redis trong compose; vẫn cần file `configs/config.env` trong image (thư mục `configs/` được copy khi build).

## Chạy local

1. Bật PostgreSQL và Redis khớp với `configs/config.env` (mặc định example trỏ `localhost`).
2. Redis: nếu dùng password (như trong compose: `123`), cập nhật `REDIS_PASSWORD` trong `config.env`.

```bash
cp configs/config.example.env configs/config.env
# chỉnh DB_*, REDIS_* cho đúng máy bạn
go run ./cmd/api
```

Mặc định server lắng nghe `:8080` (hoặc `APP_PORT` nếu set).

## Docker

```bash
# Tạo configs/config.env cho build/context (có thể copy từ example rồi chỉnh cho khớp compose)
cp configs/config.example.env configs/config.env

docker compose up --build
```

- App: [http://localhost:8080](http://localhost:8080)
- Postgres: `localhost:5432` (user/pass/db như trong `docker-compose.yml`)
- Redis: `localhost:6379`, password `123` (theo `command` của service redis)

Volume: `postgres_data`, `./redis_data` (AOF Redis).

## Vai trò từng tầng

| Thư mục | Vai trò |
|---------|---------|
| `cmd/api` | Đọc config, wire repo/service/handler, đăng ký route Gin. |
| `internal/handler` | Parse form/HTML, gọi service, redirect hoặc HTML lỗi. |
| `internal/service` | Tạo short code, gọi repository. |
| `internal/repository` | `SetURL` / `GetURL` trên Redis; kết nối Postgres khi khởi tạo. |
| `pkg/shortener` | Thuật toán encode URL → short string. |
| `configs/` | `config.env` (không commit secret; dùng `config.example.env` làm mẫu). |
| `templates/` | Giao diện form / hiển thị lỗi & short URL. |

---

**Lưu ý:** Giữ `configs/config.env` và mật khẩu Redis/Postgres khỏi git; `redis_data/` thường nên nằm trong `.gitignore` nếu chứa dữ liệu local.
