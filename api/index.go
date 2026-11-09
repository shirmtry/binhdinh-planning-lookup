// Package handler phục vụ API + frontend tra cứu quy hoạch Bình Định.
// Dự án: binhdinh-planning-lookup — dân khó tra cứu quy hoạch, dễ mua đất quy hoạch treo.
// Tác giả: HTK — Hotline/Zalo: 0396241674 — Email: htk.contact@gmail.com
// Cập nhật: 09/11/2026. Dữ liệu mô phỏng – cập nhật 09/11/2026.
package handler

import (
	"io/fs"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"

	app "binhdinh-planning-lookup"

	"github.com/gin-gonic/gin"
)

const dataNote = "Dữ liệu mô phỏng – cập nhật 09/11/2026"

// BanGhi là một bản ghi quy hoạch đất (mock) tại Bình Định.
type BanGhi struct {
	ID            int     `json:"id"`
	PhuongXa      string  `json:"phuongXa"`
	QuanHuyen     string  `json:"quanHuyen"`
	Duong         string  `json:"duong"`
	LoaiDat       string  `json:"loaiDat"`
	TrangThai     string  `json:"trangThai"`
	Nam           int     `json:"nam"` // năm phê duyệt / công bố quy hoạch
	DienTichHa    float64 `json:"dienTichHa"`
	QuyHoachTreo  bool    `json:"quyHoachTreo"`
	CanhBao       string  `json:"canhBao"`
	GhiChu        string  `json:"ghiChu"`
}

// duLieu: 32 bản ghi mô phỏng quy hoạch đất Bình Định (không thay thế hồ sơ chính thức).
var duLieu = []BanGhi{
	{1, "P. Lê Lợi", "TP. Quy Nhơn", "Nguyễn Thái Học", "Đất thương mại – dịch vụ", "Quy hoạch treo", 2016, 4.2, true, "Treo hơn 10 năm chưa triển khai – KHÔNG nên mua, khó cấp phép xây dựng.", "Khu chợ cũ đề xuất thành trung tâm thương mại nhưng chưa bố trí vốn."},
	{2, "P. Trần Phú", "TP. Quy Nhơn", "Lê Hồng Phong", "Đất ở đô thị", "Đã phê duyệt", 2024, 6.8, false, "", "Khu dân cư chỉnh trang, đã có quyết định 1/500."},
	{3, "P. Ngô Mây", "TP. Quy Nhơn", "Tây Sơn", "Đất giao thông", "Đang triển khai", 2026, 2.5, false, "", "Mở rộng đường tránh ngập úng đô thị sau mưa lũ 2026."},
	{4, "P. Nguyễn Văn Cừ", "TP. Quy Nhơn", "An Dương Vương", "Đất giáo dục", "Quy hoạch treo", 2018, 3.1, true, "Đất dính quy hoạch trường học treo – giao dịch rủi ro cao, không được xây kiên cố.", "Dự án trường liên cấp chưa giải phóng mặt bằng."},
	{5, "P. Quang Trung", "TP. Quy Nhơn", "Tăng Bạt Hổ", "Đất cây xanh", "Đã phê duyệt", 2023, 1.9, false, "", "Công viên cây xanh chống ngập khu trung tâm."},
	{6, "P. Thị Nại", "TP. Quy Nhơn", "Bạch Đằng", "Đất thương mại – dịch vụ", "Đang lấy ý kiến", 2026, 5.4, false, "", "Khu dịch vụ ven đầm Thị Nại, đang lấy ý kiến cộng đồng."},
	{7, "P. Đống Đa", "TP. Quy Nhơn", "Đống Đa", "Đất ở đô thị", "Quy hoạch treo", 2014, 7.6, true, "Quy hoạch khu tái định cư treo từ 2014 – cảnh báo đỏ, tránh mua bán sang tay.", "Chưa có nguồn vốn đầu tư hạ tầng."},
	{8, "P. Hải Cảng", "TP. Quy Nhơn", "Trần Hưng Đạo", "Đất du lịch – nghỉ dưỡng", "Đã phê duyệt", 2026, 9.3, false, "", "Tuyến phố du lịch biển, ưu tiên phục hồi sau bão số 9/2026."},
	{9, "P. Nhơn Bình", "TP. Quy Nhơn", "Võ Nguyên Giáp", "Đất khu công nghiệp", "Đang triển khai", 2024, 15.2, false, "", "Mở rộng cụm công nghiệp phía bắc sông Hà Thanh."},
	{10, "P. Nhơn Phú", "TP. Quy Nhơn", "Hùng Vương", "Đất ở đô thị", "Đã phê duyệt", 2026, 8.7, false, "", "Khu dân cư mới phục vụ giãn dân vùng ngập lụt."},
	{11, "P. Ghềnh Ráng", "TP. Quy Nhơn", "Chế Lan Viên", "Đất du lịch – nghỉ dưỡng", "Điều chỉnh cục bộ", 2026, 3.6, false, "", "Điều chỉnh mật độ xây dựng khu sườn núi Ghềnh Ráng."},
	{12, "P. Bùi Thị Xuân", "TP. Quy Nhơn", "Quốc lộ 1D", "Đất giao thông", "Quy hoạch treo", 2019, 4.8, true, "Đất nằm trong hành lang mở rộng QL1D bị treo – không cấp phép xây mới.", "Chờ vốn trung ương, chưa xác định thời điểm thu hồi."},
	{13, "X. Nhơn Lý", "TP. Quy Nhơn", "Đường vào Eo Gió", "Đất du lịch – nghỉ dưỡng", "Đang lấy ý kiến", 2026, 6.2, false, "", "Khu dịch vụ du lịch Eo Gió, lưu ý xói lở bờ biển sau bão 2026."},
	{14, "X. Nhơn Hải", "TP. Quy Nhơn", "Đường ven biển Nhơn Hải", "Đất ở nông thôn", "Đã phê duyệt", 2024, 2.4, false, "", "Khu dân cư làng chài, ưu tiên hạ tầng cho hộ nuôi tôm hùm."},
	{15, "X. Phước Mỹ", "TP. Quy Nhơn", "Tỉnh lộ 639", "Đất khu công nghiệp", "Quy hoạch treo", 2017, 12.5, true, "Cụm công nghiệp treo nhiều năm – đất chuyển nhượng giá ảo, cảnh báo đỏ.", "Nhà đầu tư bỏ dở, tỉnh đang rà soát thu hồi."},
	{16, "X. Cát Tiến", "H. Phù Cát", "Đường ven biển Cát Tiến", "Đất du lịch – nghỉ dưỡng", "Đang triển khai", 2026, 10.1, false, "", "Khu du lịch Phương Mai, kiểm soát xói lở bờ biển Mỹ Thành."},
	{17, "X. Cát Minh", "H. Phù Cát", "Đường vào làng muối", "Đất nông nghiệp", "Đã phê duyệt", 2023, 18.6, false, "", "Vùng muối truyền thống, giữ quy hoạch sản xuất sau mưa lũ 2026."},
	{18, "TT. Ngô Mây", "H. Phù Cát", "Quốc lộ 1A", "Đất thương mại – dịch vụ", "Quy hoạch treo", 2015, 3.9, true, "Khu chợ đầu mối treo – người dân mua đất dễ mất trắng nếu bị thu hồi.", "Dự án đổi chủ nhiều lần, chưa khởi công."},
	{19, "P. Bình Định", "TX. An Nhơn", "Đường tránh phía Nam", "Đất giao thông", "Đang triển khai", 2026, 5.5, false, "", "Đường tránh lũ cho vùng hạ lưu sông Kôn."},
	{20, "P. Đập Đá", "TX. An Nhơn", "Tỉnh lộ 636", "Đất ở đô thị", "Quy hoạch treo", 2013, 6.4, true, "Khu dân cư Đập Đá treo hơn 10 năm – tuyệt đối thận trọng khi giao dịch.", "Vướng giải phóng mặt bằng, chưa tái định cư."},
	{21, "P. Nhơn Thành", "TX. An Nhơn", "Quốc lộ 19", "Đất khu công nghiệp", "Đã phê duyệt", 2024, 22.0, false, "", "Khu công nghiệp Nhơn Hòa mở rộng giai đoạn 2."},
	{22, "TT. Phú Phong", "H. Tây Sơn", "Nguyễn Huệ", "Đất thương mại – dịch vụ", "Đã phê duyệt", 2026, 4.4, false, "", "Khu dịch vụ trung tâm thị trấn, chỉnh trang sau lũ sông Kôn."},
	{23, "X. Bình Nghi", "H. Tây Sơn", "Đường liên xã Bình Nghi", "Đất nông nghiệp", "Điều chỉnh cục bộ", 2026, 9.8, false, "", "Chuyển đổi cơ cấu cây trồng vùng bị bồi lấp sau lũ 2026."},
	{24, "X. Tây Thuận", "H. Tây Sơn", "Đường kênh chính Phú Phong", "Đất công cộng", "Đang triển khai", 2026, 1.6, false, "", "Hành lang thoát lũ kênh chính, cấm xây dựng lấn chiếm."},
	{25, "TT. Tăng Bạt Hổ", "H. Hoài Nhơn", "Quang Trung", "Đất ở đô thị", "Đã phê duyệt", 2024, 7.2, false, "", "Khu dân cư phía nam thị trấn, hạ tầng chống ngập."},
	{26, "P. Bồng Sơn", "TX. Hoài Nhơn", "Đường ven sông Lại Giang", "Đất cây xanh", "Đang lấy ý kiến", 2026, 3.3, false, "", "Công viên ven sông kết hợp hành lang thoát lũ."},
	{27, "X. Hoài Hải", "H. Hoài Nhơn", "Đường ven biển Hoài Hải", "Đất du lịch – nghỉ dưỡng", "Quy hoạch treo", 2020, 5.7, true, "Khu resort ven biển treo từ 2020 – nguy cơ mất đất khi thu hồi, cảnh báo đỏ.", "Chủ đầu tư chưa chứng minh năng lực tài chính."},
	{28, "TT. An Lão", "H. An Lão", "Đường liên thôn An Toàn", "Đất ở nông thôn", "Điều chỉnh cục bộ", 2026, 2.2, false, "", "Bố trí tái định cư cho hộ nguy cơ sạt lở sau mưa lớn 11/2026."},
	{29, "TT. Vĩnh Thạnh", "H. Vĩnh Thạnh", "Hồ Chí Minh nhánh Đông", "Đất giao thông", "Đang triển khai", 2026, 3.8, false, "", "Nắn tuyến tránh điểm sạt lở thôn Định Bình."},
	{30, "X. Vĩnh Sơn", "H. Vĩnh Thạnh", "Đường suối Kôn thượng", "Đất nông nghiệp", "Quy hoạch treo", 2012, 11.3, true, "Khu nông nghiệp công nghệ cao treo – không nên đầu tư xây dựng kiên cố.", "Chưa giải phóng mặt bằng, người dân tiếp tục canh tác."},
	{31, "TT. Vân Canh", "H. Vân Canh", "Đường trung tâm Vân Canh", "Đất công cộng", "Đã phê duyệt", 2026, 1.4, false, "", "Khu hành chính – văn hóa, phục vụ đồng bào vùng núi."},
	{32, "X. Canh Thuận", "H. Vân Canh", "Tỉnh lộ 638", "Đất ở nông thôn", "Đang lấy ý kiến", 2026, 4.9, false, "", "Điểm dân cư mới tránh lũ quét lòng suối."},
}

// dauMap bỏ dấu tổ hợp phổ biến (không dùng x/text để nhẹ dependency).
var dauMap = map[rune]rune{
	'à': 'a', 'á': 'a', 'ạ': 'a', 'ả': 'a', 'ã': 'a',
	'â': 'a', 'ầ': 'a', 'ấ': 'a', 'ậ': 'a', 'ẩ': 'a', 'ẫ': 'a',
	'ă': 'a', 'ằ': 'a', 'ắ': 'a', 'ặ': 'a', 'ẳ': 'a', 'ẵ': 'a',
	'è': 'e', 'é': 'e', 'ẹ': 'e', 'ẻ': 'e', 'ẽ': 'e',
	'ê': 'e', 'ề': 'e', 'ế': 'e', 'ệ': 'e', 'ể': 'e', 'ễ': 'e',
	'ì': 'i', 'í': 'i', 'ị': 'i', 'ỉ': 'i', 'ĩ': 'i',
	'ò': 'o', 'ó': 'o', 'ọ': 'o', 'ỏ': 'o', 'õ': 'o',
	'ô': 'o', 'ồ': 'o', 'ố': 'o', 'ộ': 'o', 'ổ': 'o', 'ỗ': 'o',
	'ơ': 'o', 'ờ': 'o', 'ớ': 'o', 'ợ': 'o', 'ở': 'o', 'ỡ': 'o',
	'ù': 'u', 'ú': 'u', 'ụ': 'u', 'ủ': 'u', 'ũ': 'u',
	'ư': 'u', 'ừ': 'u', 'ứ': 'u', 'ự': 'u', 'ử': 'u', 'ữ': 'u',
	'ỳ': 'y', 'ý': 'y', 'ỵ': 'y', 'ỷ': 'y', 'ỹ': 'y',
}

func chuanHoa(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	for _, r := range s {
		if v, ok := dauMap[r]; ok {
			b.WriteRune(v)
		} else if r == 'đ' {
			b.WriteRune('d')
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// timKiem lọc bản ghi theo từ khóa (phường/xã, đường, huyện) + loại đất + trạng thái.
func timKiem(q, loaiDat, trangThai string, chiTreo bool) []BanGhi {
	kw := chuanHoa(q)
	out := make([]BanGhi, 0, len(duLieu))
	for _, bg := range duLieu {
		if kw != "" {
			gop := chuanHoa(bg.PhuongXa + " " + bg.QuanHuyen + " " + bg.Duong + " " + bg.LoaiDat + " " + bg.GhiChu)
			if !strings.Contains(gop, kw) {
				continue
			}
		}
		if loaiDat != "" && bg.LoaiDat != loaiDat {
			continue
		}
		if trangThai != "" && bg.TrangThai != trangThai {
			continue
		}
		if chiTreo && !bg.QuyHoachTreo {
			continue
		}
		out = append(out, bg)
	}
	return out
}

func demTreo(ds []BanGhi) int {
	n := 0
	for _, bg := range ds {
		if bg.QuyHoachTreo {
			n++
		}
	}
	return n
}

var (
	once     sync.Once
	ginApp   *gin.Engine
	staticFS fs.FS
)

func buildApp() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	staticFS, _ = fs.Sub(app.PublicFS, "public")

	serveFile := func(c *gin.Context, name string, status int) {
		b, err := fs.ReadFile(staticFS, name)
		if err != nil {
			c.String(500, "Không đọc được tài nguyên.")
			return
		}
		ct := "application/octet-stream"
		switch path.Ext(name) {
		case ".html":
			ct = "text/html; charset=utf-8"
		case ".css":
			ct = "text/css; charset=utf-8"
		case ".js":
			ct = "application/javascript; charset=utf-8"
		}
		c.Data(status, ct, b)
	}

	api := r.Group("/api")
	api.GET("/search", func(c *gin.Context) {
		q := c.Query("q")
		loaiDat := c.Query("loaiDat")
		trangThai := c.Query("trangThai")
		chiTreo := c.Query("treo") == "true" || c.Query("treo") == "1"
		kq := timKiem(q, loaiDat, trangThai, chiTreo)
		c.JSON(200, gin.H{
			"ok": true, "note": dataNote,
			"query": gin.H{"q": q, "loaiDat": loaiDat, "trangThai": trangThai, "chiTreo": chiTreo},
			"total": len(kq), "soQuyHoachTreo": demTreo(kq), "records": kq,
			"contact": gin.H{"author": "HTK", "hotline": "0396241674", "email": "htk.contact@gmail.com"},
		})
	})
	api.GET("/records", func(c *gin.Context) {
		limit := len(duLieu)
		if v, err := strconv.Atoi(c.Query("limit")); err == nil && v > 0 && v < limit {
			limit = v
		}
		c.JSON(200, gin.H{"ok": true, "note": dataNote, "total": len(duLieu), "soQuyHoachTreo": demTreo(duLieu), "records": duLieu[:limit]})
	})
	api.GET("/stats", func(c *gin.Context) {
		theoLoai := map[string]int{}
		theoTrangThai := map[string]int{}
		for _, bg := range duLieu {
			theoLoai[bg.LoaiDat]++
			theoTrangThai[bg.TrangThai]++
		}
		c.JSON(200, gin.H{"ok": true, "note": dataNote, "total": len(duLieu), "soQuyHoachTreo": demTreo(duLieu),
			"theoLoaiDat": theoLoai, "theoTrangThai": theoTrangThai})
	})

	r.GET("/", func(c *gin.Context) { serveFile(c, "index.html", 200) })
	r.GET("/about", func(c *gin.Context) { serveFile(c, "about.html", 200) })
	r.GET("/404", func(c *gin.Context) { serveFile(c, "404.html", 404) })
	r.GET("/style.css", func(c *gin.Context) { serveFile(c, "style.css", 200) })
	r.GET("/app.js", func(c *gin.Context) { serveFile(c, "app.js", 200) })

	r.NoRoute(func(c *gin.Context) {
		p := strings.TrimPrefix(c.Request.URL.Path, "/")
		if strings.HasPrefix(p, "api/") {
			c.JSON(404, gin.H{"ok": false, "error": "Không tìm thấy API.", "hotline": "0396241674"})
			return
		}
		serveFile(c, "404.html", 404)
	})
	return r
}

// restorePath khôi phục đường dẫn gốc khi Vercel rewrite qua /api/index/...
func restorePath(r *http.Request) {
	p := r.URL.Path
	if orig := r.Header.Get("X-Original-Path"); orig != "" {
		p = orig
	} else {
		for _, prefix := range []string{"/api/index.go", "/api/index"} {
			if p == prefix {
				p = "/"
				break
			}
			if strings.HasPrefix(p, prefix+"/") {
				p = strings.TrimPrefix(p, prefix)
				break
			}
		}
	}
	r.URL.Path = p
	if q := r.URL.RawQuery; q != "" {
		r.RequestURI = p + "?" + q
	} else {
		r.RequestURI = p
	}
}

// Handler là entrypoint cho Vercel Go runtime (api/index.go).
func Handler(w http.ResponseWriter, r *http.Request) {
	restorePath(r)
	once.Do(func() { ginApp = buildApp() })
	ginApp.ServeHTTP(w, r)
}