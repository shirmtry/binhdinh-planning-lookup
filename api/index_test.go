// Smoke tests httptest cho Handler (/, /about, /api/search).
// HTK – 0396241674 – cập nhật 09/11/2026
package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func goRequest(path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	w := httptest.NewRecorder()
	Handler(w, req)
	return w
}

func TestHome(t *testing.T) {
	w := goRequest("/")
	if w.Code != 200 {
		t.Fatalf("/ muốn 200, nhận %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "0396241674") {
		t.Fatal("/ thiếu hotline HTK 0396241674")
	}
}

func TestAbout(t *testing.T) {
	w := goRequest("/about")
	if w.Code != 200 || !strings.Contains(w.Body.String(), "Cập nhật: 09/11/2026") {
		t.Fatalf("/about lỗi: code=%d", w.Code)
	}
}

func TestNotFound(t *testing.T) {
	w := goRequest("/khong-ton-tai")
	if w.Code != 404 || !strings.Contains(w.Body.String(), "0396241674") {
		t.Fatalf("404 lỗi: code=%d", w.Code)
	}
}

func TestSearchJSON(t *testing.T) {
	w := goRequest("/api/search?q=quy%20nhon")
	if w.Code != 200 {
		t.Fatalf("/api/search muốn 200, nhận %d", w.Code)
	}
	var body struct {
		OK      bool     `json:"ok"`
		Total   int      `json:"total"`
		Records []BanGhi `json:"records"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("JSON không hợp lệ: %v", err)
	}
	if !body.OK || body.Total == 0 || len(body.Records) != body.Total {
		t.Fatalf("kết quả tìm 'quy nhon' không đúng: total=%d", body.Total)
	}
}

func TestSearchTreo(t *testing.T) {
	w := goRequest("/api/search?treo=true")
	var body struct {
		Total   int      `json:"total"`
		SoTreo  int      `json:"soQuyHoachTreo"`
		Records []BanGhi `json:"records"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Total == 0 || body.Total != body.SoTreo {
		t.Fatalf("lọc treo không đúng: total=%d treo=%d", body.Total, body.SoTreo)
	}
	for _, bg := range body.Records {
		if !bg.QuyHoachTreo || bg.CanhBao == "" {
			t.Fatalf("bản ghi %d thiếu cảnh báo treo", bg.ID)
		}
	}
}

func TestRecordsCount(t *testing.T) {
	w := goRequest("/api/records")
	var body struct {
		Total int `json:"total"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Total < 30 {
		t.Fatalf("cần >=30 bản ghi mock, có %d", body.Total)
	}
}

// Giả lập rewrite của Vercel: /api/index/api/search -> /api/search
func TestVercelRewrite(t *testing.T) {
	w := goRequest("/api/index/api/search?q=An%20Nh%C6%A1n")
	if w.Code != 200 {
		t.Fatalf("rewrite Vercel lỗi: %d", w.Code)
	}
}