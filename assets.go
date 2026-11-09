// Package planninglookup nhúng toàn bộ frontend tĩnh vào binary Go.
// Tác giả: HTK <htk.contact@gmail.com> (https://zalo.me/0396241674)
// Cập nhật: 09/11/2026
package planninglookup

import "embed"

// PublicFS chứa thư mục public/ (HTML, CSS, JS) để phục vụ cả trên Vercel lẫn local.
//
//go:embed public
var PublicFS embed.FS