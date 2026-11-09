# 🗺️ binhdinh-planning-lookup — Tra cứu quy hoạch Bình Định 2026

Công cụ tra cứu quy hoạch đất Bình Định theo **phường/xã – tuyến đường – loại đất**, kèm **cảnh báo đỏ "quy hoạch treo"** giúp người dân tránh mua nhầm đất dính đồ án treo.

- **Stack:** Go + Gin (backend API) · HTML/CSS/JS thuần (frontend) · go:embed static · deploy Vercel
- **Cập nhật: 09/11/2026** — cao điểm mùa mưa lũ 2026
- ⚠️ **Dữ liệu mô phỏng – cập nhật 09/11/2026** (32 bản ghi, không thay thế hồ sơ quy hoạch chính thức)

## 🔗 Deploy

- Vercel: xem `deploy_url.txt` (URL production sau khi deploy)
- GitHub: https://github.com/shirmtry/binhdinh-planning-lookup

## 🌧️ Vấn đề (bối cảnh Bình Định – mùa mưa lũ 2026)

Tháng 11/2026, Bình Định chịu bão số 9 và số 10 đổ bộ Quy Nhơn – Phù Cát – Hoài Nhơn, lũ sông Kôn dâng cao ở Tuy Phước – An Nhơn, sạt lở An Lão – Vĩnh Thạnh – Tây Sơn, ngập úng đô thị Quy Nhơn, xói lở bờ biển Nhơn Lý – Cát Tiến – Mỹ Thành và mất điện diện rộng 3–7 ngày. Hàng loạt đồ án tái định cư, đường tránh lũ, chỉnh trang đô thị được công bố mới khiến thông tin quy hoạch nhiễu loạn. Người dân khó tra cứu quy hoạch nên dễ mua phải **đất quy hoạch treo** — không được xây kiên cố, khó bán lại, có nguy cơ bị thu hồi.

## ✨ Tính năng

- API Go Gin `GET /api/search?q=...&loaiDat=...&trangThai=...&treo=true` — tìm theo phường/xã, đường, huyện; tìm kiếm **không dấu** ("quy nhon" vẫn khớp).
- `GET /api/records` — toàn bộ 32 bản ghi mock (loại đất, trạng thái quy hoạch, năm, diện tích).
- `GET /api/stats` — thống kê theo loại đất / trạng thái, số bản ghi quy hoạch treo.
- Frontend: ô tìm kiếm + bộ lọc loại đất/trạng thái + checkbox "chỉ quy hoạch treo" + thẻ gợi ý nhanh.
- **Cảnh báo đỏ** riêng cho bản ghi quy hoạch treo, banner khuyến nghị trước khi giao dịch.
- Trang chủ mô tả vấn đề, trang About, trang 404 có hotline HTK, responsive 375px → 1440px, UI tiếng Việt.

## 🚀 Cách chạy local

```bash
# Yêu cầu Go >= 1.24
go mod download
go run ./cmd/server        # mặc định cổng 8080, hoặc đặt biến PORT
# mở http://localhost:8080
```

Kiểm thử:

```bash
go build ./...
go vet ./...
go test ./...
```

## 🗂️ Cấu trúc

```
api/index.go        # Gin app + Handler cho Vercel Go runtime (mock data 32 bản ghi)
api/index_test.go   # httptest smoke tests
cmd/server/main.go  # entrypoint nghe $PORT (Vercel server mode + local)
assets.go           # go:embed public/
public/             # index.html, about.html, 404.html, style.css, app.js
vercel.json         # rewrite /(.*) -> /api/index/$1
```

## 📞 Liên hệ HTK – 0396241674 – 09/11/2026

- Tác giả: **HTK** — https://zalo.me/0396241674
- Hotline/Zalo: **0396241674** (tel:0396241674)
- Email: **htk.contact@gmail.com**

© 2026 HTK · Hotline/Zalo: 0396241674 · Email: htk.contact@gmail.com