# Go Shortener Service

Dịch vụ rút gọn URL viết bằng Go. Module: `Go_shortenURL` (Go 1.25.x).

## Cấu trúc thư mục

```text
Go_shortenURL/
├── cmd/
│   └── api/                # Điểm khởi chạy ứng dụng
│       └── main.go         # Khởi tạo server, wire dependencies
├── internal/               # Logic cốt lõi (không cho bên ngoài import)
│   ├── handler/            # HTTP handlers — nhận request, trả response
│   ├── service/            # Business logic — tạo mã hash, kiểm tra URL
│   ├── repository/         # Tương tác với Database (Redis, Postgres)
│   └── model/              # Định nghĩa struct (URL, User, ...)
├── pkg/
│   └── shortener/          # Thuật toán băm (hashing) hoặc sinh ID ngẫu nhiên
├── configs/
│   └── config.example.env  # Biến môi trường mẫu (PORT, DB_URL, REDIS_ADDR)
├── templates/              # HTML templates (Go html/template)
├── go.mod
├── go.sum
└── README.md
```

## Vai trò từng tầng

| Thư mục             | Vai trò                                                                                              |
| ------------------- | ---------------------------------------------------------------------------------------------------- |
| `cmd/api`           | Binary duy nhất — đọc config, khởi tạo deps, gọi `http.ListenAndServe`.                             |
| `internal/handler`  | Parse HTTP request, gọi Service, trả JSON hoặc redirect 301/302.                                    |
| `internal/service`  | Không phụ thuộc `net/http` — dễ unit test. Tạo short code, validate URL.                            |
| `internal/repository` | Interface + implementation (Postgres lưu trữ lâu dài, Redis cache redirect nhanh).                |
| `internal/model`    | Plain struct dùng chung giữa các tầng. Không chứa business logic.                                   |
| `pkg/shortener`     | Có thể import từ ngoài module. Chứa thuật toán sinh slug (Base62, nanoid, ...).                      |
| `configs/`          | `.env` hoặc YAML config. Commit file `.example.env`, **không commit** file thật.                     |
| `templates/`        | File `.html` dùng với `html/template` cho trang redirect hoặc landing page.                          |

## Luồng gọi

```text
HTTP Request
    │
    ▼
handler  ──►  service  ──►  repository (Postgres / Redis)
                │
                ▼
           pkg/shortener (sinh short code)
```

## Chạy nhanh

```bash
cp configs/config.example.env .env
go run ./cmd/api
# → Server starting on :8080
```
