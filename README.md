# Go_shortenURL

Dịch vụ rút gọn URL viết bằng Go. Module: `Go_shortenURL` (Go 1.25.x).

## Trạng thái repo

Các thư mục template đã có file `.gitkeep` (Git không track folder rỗng). Thêm `main.go`, package Go, v.v. và có thể xóa `.gitkeep` khi thư mục đã có file thật.

Cấu trúc bên dưới khớp repo; `pkg/shortener/` thêm cho đúng workflow (hash slug) trong README.

## Cấu trúc thư mục (đề xuất)

```text
Go_shortenURL/
├── cmd/
│   └── server/                 # Binary chạy HTTP server (.gitkeep → thêm main.go)
├── internal/
│   ├── handler/
│   ├── service/
│   ├── storage/
│   └── config/
├── pkg/
│   └── shortener/              # Hash / tạo slug ngắn (workflow)
├── api/                        # OpenAPI / contract (tuỳ chọn)
├── scripts/                    # migrate, seed (tuỳ chọn)
├── go.mod
├── go.sum
└── README.md
```

### Ý nghĩa từng phần

| Thư mục / file     | Vai trò                                                                                                                                |
| ------------------ | -------------------------------------------------------------------------------------------------------------------------------------- |
| `cmd/server`       | Mỗi thư mục con của `cmd/` là một chương trình build được (`go build ./cmd/server`). `main.go` chỉ nên khởi tạo và gọi vào `internal`. |
| `internal/handler` | `http.Handler` / chi router: parse body, status code, redirect 302/301.                                                                |
| `internal/service` | Không phụ thuộc trực tiếp `net/http`: dễ test (tạo slug, trùng slug, giới hạn độ dài URL).                                             |
| `internal/storage` | Interface repository + implementation (in-memory cho dev, sau nâng DB).                                                                |
| `internal/config`  | `PORT`, `BASE_URL`, DSN — tách khỏi handler.                                                                                           |

### Không bắt buộc

- **`pkg/`** — chỉ khi có thư viện muốn **export** cho project/module khác import; URL shortener đơn giản thường không cần.
- **`test/`** hoặc file `*_test.go` cạnh package — test nằm cùng package hoặc `package xxx_test` cho black-box test.

### Luồng gọi (tham khảo)

```text
HTTP → internal/handler → internal/service → internal/storage
                ↑
         internal/config (đọc khi start từ main)
```

Sau khi thêm code, cập nhật mục "Trạng thái repo" và cây thư mục cho khớp thực tế.

2. Phân tích luồng dữ liệu (Workflow)

Để dự án "sạch" và dễ test, bạn nên tổ chức theo luồng: Handler → Service → Repository.
🔹 1. internal/repository (Tầng dữ liệu)

Nơi chứa code thao tác với DB. Với URL Shortener, bạn có thể dùng Redis để redirect cực nhanh hoặc PostgreSQL để lưu trữ lâu dài.

    Hàm ví dụ: Save(shortCode, longURL), Get(shortCode).

🔹 2. internal/service (Tầng nghiệp vụ)

Nơi xử lý "não bộ" của app.

    Nhận URL dài từ Handler.

    Gọi hàm băm (hashing) từ pkg/shortener để tạo code ngắn (ví dụ: aB3x9).

    Gọi Repository để lưu vào DB.

🔹 3. internal/handler (Tầng giao tiếp)

Nơi tiếp nhận HTTP request.

    POST /shorten: Nhận JSON, gọi Service, trả về URL ngắn.

    GET /{code}: Nhận code, gọi Service để lấy URL gốc, sau đó dùng http.Redirect.
