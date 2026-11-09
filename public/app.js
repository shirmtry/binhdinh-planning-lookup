// Tra cứu quy hoạch Bình Định – HTK 0396241674 – cập nhật 09/11/2026
// Frontend gọi API Go Gin: /api/search?q=...&loaiDat=...&trangThai=...&treo=true
"use strict";

const form = document.getElementById("search-form");
const qInput = document.getElementById("q");
const loaiDatSel = document.getElementById("loaiDat");
const trangThaiSel = document.getElementById("trangThai");
const chiTreoChk = document.getElementById("chiTreo");
const summaryEl = document.getElementById("summary");
const treoAlert = document.getElementById("treo-alert");
const resultsEl = document.getElementById("results");

function esc(s) {
  const d = document.createElement("div");
  d.textContent = String(s ?? "");
  return d.innerHTML;
}

function pillTrangThai(bg) {
  if (bg.quyHoachTreo) return '<span class="pill treo">🚨 QUY HOẠCH TREO</span>';
  if (bg.trangThai === "Đã phê duyệt" || bg.trangThai === "Đang triển khai")
    return '<span class="pill ok">' + esc(bg.trangThai) + "</span>";
  return '<span class="pill warn">' + esc(bg.trangThai) + "</span>";
}

function render(data) {
  summaryEl.hidden = false;
  summaryEl.textContent =
    "Tìm thấy " + data.total + " bản ghi · " + data.soQuyHoachTreo +
    " bản ghi quy hoạch treo · " + data.note;
  treoAlert.hidden = data.soQuyHoachTreo === 0;

  if (!data.records || data.records.length === 0) {
    resultsEl.innerHTML =
      '<div class="empty">Không tìm thấy bản ghi phù hợp. Thử từ khóa khác (VD: "Quy Nhơn", "Nguyễn Thái Học") hoặc bỏ bớt bộ lọc.<br>Cần hỗ trợ: gọi HTK <a href="tel:0396241674">0396241674</a>.</div>';
    return;
  }
  resultsEl.innerHTML = data.records.map(function (bg) {
    const cls = bg.quyHoachTreo ? "record treo" : (bg.trangThai === "Đã phê duyệt" || bg.trangThai === "Đang triển khai" ? "record" : "record pending");
    return (
      '<article class="' + cls + '">' +
      "<h3>#" + bg.id + " · " + esc(bg.phuongXa) + " — " + esc(bg.duong) + "</h3>" +
      '<p class="meta">' + esc(bg.quanHuyen) + " · Năm công bố: " + bg.nam +
      " · Diện tích: " + bg.dienTichHa.toLocaleString("vi-VN") + " ha</p>" +
      '<span class="pill type">' + esc(bg.loaiDat) + "</span>" + pillTrangThai(bg) +
      (bg.quyHoachTreo ? '<div class="warn-box">⛔ ' + esc(bg.canhBao) + "</div>" : "") +
      '<p class="note">' + esc(bg.ghiChu) + "</p>" +
      "</article>"
    );
  }).join("");
}

async function search() {
  const params = new URLSearchParams();
  if (qInput.value.trim()) params.set("q", qInput.value.trim());
  if (loaiDatSel.value) params.set("loaiDat", loaiDatSel.value);
  if (trangThaiSel.value) params.set("trangThai", trangThaiSel.value);
  if (chiTreoChk.checked) params.set("treo", "true");
  resultsEl.innerHTML = '<div class="empty">Đang tra cứu dữ liệu quy hoạch…</div>';
  try {
    const res = await fetch("/api/search?" + params.toString());
    if (!res.ok) throw new Error("HTTP " + res.status);
    render(await res.json());
  } catch (err) {
    summaryEl.hidden = true;
    treoAlert.hidden = true;
    resultsEl.innerHTML =
      '<div class="empty">Lỗi khi tra cứu (' + esc(err.message) + '). Vui lòng thử lại hoặc gọi HTK <a href="tel:0396241674">0396241674</a>.</div>';
  }
}

form.addEventListener("submit", function (e) {
  e.preventDefault();
  search();
});

document.querySelectorAll(".tag").forEach(function (btn) {
  btn.addEventListener("click", function () {
    if (btn.dataset.treo) {
      chiTreoChk.checked = true;
      qInput.value = "";
      loaiDatSel.value = "";
      trangThaiSel.value = "";
    } else {
      qInput.value = btn.dataset.q || "";
      chiTreoChk.checked = false;
    }
    search();
    document.getElementById("ket-qua").scrollIntoView({ behavior: "smooth" });
  });
});

// Tải toàn bộ bản ghi khi mở trang.
search();