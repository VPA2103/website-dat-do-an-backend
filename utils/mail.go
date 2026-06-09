package utils

import (
	"fmt"
	"strconv"

	"github.com/vpa/quanlynhahang-backend/config"
	"gopkg.in/gomail.v2"
)

func SendMail(to, subject, body string) error {
	host := config.GetEnv("MAIL_HOST")
	portStr := config.GetEnv("MAIL_PORT")
	username := config.GetEnv("MAIL_USERNAME")
	password := config.GetEnv("MAIL_PASSWORD")
	from := config.GetEnv("MAIL_FROM")

	port, _ := strconv.Atoi(portStr)

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(
		host,
		port,
		username,
		password,
	)

	return d.DialAndSend(m)
}

type DatDoAnMailInfo struct {
	TenKhachHang string
	MaDon        string
	NgayGio      string
	DiaChi       string
	SoMonAn      int
	TamTinh      float64
	TienGiam     float64
	TongCuoi     float64
	GhiChu       string
}

type ThanhToanMailInfo struct {
    TenKhachHang string
    MaDon        string
    NgayGio      string
    TongTien     float64
    SoTienDaTra  float64
}

func SendMailSauKhiDatDoAn(email string, info DatDoAnMailInfo) error {
	body := fmt.Sprintf(`
				<!DOCTYPE html>
				<html lang="vi">
				<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"></head>
				 <body style="margin:0;padding:20px;background:#f4ede0;font-family:'Be Vietnam Pro',Arial,sans-serif;">
				<div style="max-width:520px;margin:0 auto;background:#ffffff;border-radius:8px;overflow:hidden;">

					<!-- Header -->
					<div style="background:#1a1a1a;padding:32px 32px 24px;text-align:center;">
					<div style="font-size:22px;letter-spacing:4px;color:#e8d5b0;">✦ NHÀ HÀNG ✦</div>
					<div style="font-size:11px;letter-spacing:6px;color:#8a7a5a;margin-top:4px;font-family:'Courier New',monospace;">SAIGON KITCHEN</div>
					</div>

					<!-- Title -->
					<div style="background:#f7f0e3;padding:24px 32px 16px;text-align:center;border-bottom:1px solid #e0d0b0;">
					<div style="font-size:11px;letter-spacing:5px;color:#8a7a5a;font-family:'Courier New',monospace;margin-bottom:8px;">THÔNG BÁO</div>
					<div style="font-size:22px;color:#2a1f0a;letter-spacing:1px;">Đặt món ăn thành công</div>
					<div style="width:40px;height:1px;background:#c4a55a;margin:12px auto 0;"></div>
					</div>

					<!-- Body -->
					<div style="padding:24px 32px;background:#fdfaf4;">
					<p style="font-size:14px;color:#4a3c20;line-height:1.8;margin:0 0 16px;">
						Kính gửi <strong>%s</strong>,
					</p>
					<p style="font-size:14px;color:#4a3c20;line-height:1.8;margin:0 0 24px;">
						Chúng tôi xin xác nhận bạn đã đặt đồ ăn tại quán thành công. Dưới đây là chi tiết đơn đặt của bạn:
					</p>

					<!-- Info table -->
					<table style="width:100%%;border-collapse:collapse;font-size:13px;background:#fff;border:0.5px solid #e0d0b0;border-radius:8px;overflow:hidden;margin-bottom:24px;">
						<tr style="border-bottom:0.5px solid #f0e4c8;">
						<td style="padding:10px 14px;color:#8a7a5a;width:40%%;">Mã đơn</td>
						<td style="padding:10px 14px;color:#2a1f0a;font-family:'Courier New',monospace;font-weight:600;">#%s</td>
						</tr>
						<tr style="border-bottom:0.5px solid #f0e4c8;">
						<td style="padding:10px 14px;color:#8a7a5a;">Ngày &amp; giờ</td>
						<td style="padding:10px 14px;color:#2a1f0a;">%s</td>
						</tr>
						<tr style="border-bottom:0.5px solid #f0e4c8;">
						<td style="padding:10px 14px;color:#8a7a5a;">Tạm tính</td>
						<td style="padding:10px 14px;color:#2a1f0a;">%s</td>
						</tr>
						<tr>
						<td style="padding:10px 14px;color:#8a7a5a;">Tổng cộng</td>
						<td style="padding:10px 14px;color:#2a1f0a;font-weight:600;">%s</td>
						</tr>
					</table>

					<!-- Ghi chú -->
					<div style="background:#f7f0e3;border-left:3px solid #c4a55a;padding:12px 14px;margin-bottom:24px;">
						<p style="font-size:13px;color:#5a4520;margin:0;line-height:1.6;">
						<strong>Ghi chú:</strong> %s
						</p>
					</div>

					<p style="font-size:13px;color:#6a5a3a;line-height:1.8;margin:0;">
						Nếu bạn cần thay đổi hoặc hủy đặt bàn, vui lòng liên hệ trước <strong>2 tiếng</strong> so với giờ đặt.
					</p>
					</div>

					<!-- CTA -->
					<div style="padding:16px 32px;background:#fdfaf4;border-top:0.5px solid #e0d0b0;text-align:center;">
					<div style="display:inline-block;background:#1a1a1a;color:#e8d5b0;font-size:12px;letter-spacing:3px;padding:10px 28px;font-family:'Courier New',monospace;">
						HẸN GẶP BẠN!
					</div>
					</div>

					<!-- Footer -->
					<div style="padding:16px 32px;background:#1a1a1a;text-align:center;">
					<p style="font-size:11px;color:#6a5a3a;margin:0;letter-spacing:1px;font-family:'Courier New',monospace;">
						123 Đường ABC, Q.1, TP.HCM &nbsp;|&nbsp; 028-xxxx-xxxx
					</p>
					</div>

				</div>
				</body>
				</html>`,
		info.TenKhachHang,
		fmt.Sprintf("HD%s", info.MaDon),
		info.NgayGio,
		// info.DiaChi,
		// info.SoMonAn,
		formatVND(info.TamTinh),
		formatVND(info.TongCuoi),
		info.GhiChu,
	)

	return SendMail(email, "✦ Xác nhận đặt đồ ăn thành công", body)
}


func SendMailSauKhiThanhToan(email string, info ThanhToanMailInfo) error {
    body := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="vi">
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"></head>
<body style="margin:0;padding:20px;background:#f4ede0;font-family:'Be Vietnam Pro',Arial,sans-serif;">
<div style="max-width:520px;margin:0 auto;background:#ffffff;border-radius:8px;overflow:hidden;">

    <div style="background:#1a1a1a;padding:32px 32px 24px;text-align:center;">
        <div style="font-size:22px;letter-spacing:4px;color:#e8d5b0;">✦ NHÀ HÀNG ✦</div>
        <div style="font-size:11px;letter-spacing:6px;color:#8a7a5a;margin-top:4px;font-family:'Courier New',monospace;">SAIGON KITCHEN</div>
    </div>

    <div style="background:#f7f0e3;padding:24px 32px 16px;text-align:center;border-bottom:1px solid #e0d0b0;">
        <div style="font-size:11px;letter-spacing:5px;color:#8a7a5a;font-family:'Courier New',monospace;margin-bottom:8px;">XÁC NHẬN</div>
        <div style="font-size:22px;color:#2a1f0a;letter-spacing:1px;">Thanh toán thành công</div>
        <div style="width:40px;height:1px;background:#c4a55a;margin:12px auto 0;"></div>
    </div>

    <div style="background:#fdfaf4;padding:20px 32px 0;text-align:center;">
        <div style="display:inline-block;background:#e8f5e9;border-radius:50%%;width:52px;height:52px;line-height:52px;font-size:24px;">✓</div>
    </div>

    <div style="padding:20px 32px 24px;background:#fdfaf4;">
        <p style="font-size:14px;color:#4a3c20;line-height:1.8;margin:0 0 16px;">
            Kính gửi <strong>%s</strong>,
        </p>
        <p style="font-size:14px;color:#4a3c20;line-height:1.8;margin:0 0 24px;">
            Chúng tôi xác nhận đã nhận được thanh toán của bạn. Cảm ơn bạn đã dùng bữa tại <strong>Saigon Kitchen</strong>!
        </p>

        <table style="width:100%%;border-collapse:collapse;font-size:13px;background:#fff;border:0.5px solid #e0d0b0;border-radius:8px;overflow:hidden;margin-bottom:24px;">
            <tr style="border-bottom:0.5px solid #f0e4c8;">
                <td style="padding:10px 14px;color:#8a7a5a;width:45%%;">Mã hóa đơn</td>
                <td style="padding:10px 14px;color:#2a1f0a;font-family:'Courier New',monospace;font-weight:600;">#HD%s</td>
            </tr>
            <tr style="border-bottom:0.5px solid #f0e4c8;">
                <td style="padding:10px 14px;color:#8a7a5a;">Thời gian</td>
                <td style="padding:10px 14px;color:#2a1f0a;">%s</td>
            </tr>
            <tr style="border-bottom:0.5px solid #f0e4c8;">
                <td style="padding:10px 14px;color:#8a7a5a;">Tổng hóa đơn</td>
                <td style="padding:10px 14px;color:#2a1f0a;">%s</td>
            </tr>
            <tr>
                <td style="padding:10px 14px;color:#8a7a5a;">Đã thanh toán</td>
                <td style="padding:10px 14px;color:#2e7d32;font-weight:700;font-size:14px;">%s</td>
            </tr>
        </table>

        <div style="background:#f7f0e3;border-left:3px solid #c4a55a;padding:12px 14px;">
            <p style="font-size:13px;color:#5a4520;margin:0;line-height:1.6;">
                Nếu có bất kỳ thắc mắc nào về hóa đơn, vui lòng liên hệ nhân viên hoặc gọi <strong>028-xxxx-xxxx</strong>.
            </p>
        </div>
    </div>

    <div style="padding:16px 32px;background:#fdfaf4;border-top:0.5px solid #e0d0b0;text-align:center;">
        <div style="display:inline-block;background:#1a1a1a;color:#e8d5b0;font-size:12px;letter-spacing:3px;padding:10px 28px;font-family:'Courier New',monospace;">
            CẢM ƠN QUÝ KHÁCH!
        </div>
    </div>

    <div style="padding:16px 32px;background:#1a1a1a;text-align:center;">
        <p style="font-size:11px;color:#6a5a3a;margin:0;letter-spacing:1px;font-family:'Courier New',monospace;">
            123 Đường ABC, Q.1, TP.HCM &nbsp;|&nbsp; 028-xxxx-xxxx
        </p>
    </div>

</div>
</body>
</html>`,
        info.TenKhachHang,
        info.MaDon,
        info.NgayGio,
        formatVND(info.TongTien),
        formatVND(info.SoTienDaTra),
    )

    return SendMail(email, "✦ Xác nhận thanh toán thành công – Saigon Kitchen", body)
}

func formatVND(amount float64) string {
	s := fmt.Sprintf("%.0f", amount)
	// Thêm dấu chấm mỗi 3 chữ số từ phải sang
	n := len(s)
	var result []byte
	for i, c := range s {
		if i > 0 && (n-i)%3 == 0 {
			result = append(result, '.')
		}
		result = append(result, byte(c))
	}
	return string(result) + " ₫"
}

type DangKyMailInfo struct {
	TenKhachHang string
	Email        string
	MaNguoiDung  uint
}

func SendMailSauKhiDangKy(email string, info DangKyMailInfo) error {
	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="vi">
		<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"></head>
		<body style="margin:0;padding:20px;background:#f4ede0;font-family:Georgia,'Times New Roman',serif;">
		<div style="max-width:520px;margin:0 auto;background:#ffffff;border-radius:8px;overflow:hidden;">

			<!-- Header -->
			<div style="background:#1a1a1a;padding:32px 32px 24px;text-align:center;">
			<div style="font-size:22px;letter-spacing:4px;color:#e8d5b0;">✦ NHÀ HÀNG ✦</div>
			<div style="font-size:11px;letter-spacing:6px;color:#8a7a5a;margin-top:4px;font-family:'Courier New',monospace;">SAIGON KITCHEN</div>
			</div>

			<!-- Title -->
			<div style="background:#f7f0e3;padding:24px 32px 16px;text-align:center;border-bottom:1px solid #e0d0b0;">
			<div style="font-size:11px;letter-spacing:5px;color:#8a7a5a;font-family:'Courier New',monospace;margin-bottom:8px;">CHÀO MỪNG</div>
			<div style="font-size:22px;color:#2a1f0a;letter-spacing:1px;">Đăng ký thành công</div>
			<div style="width:40px;height:1px;background:#c4a55a;margin:12px auto 0;"></div>
			</div>

			<!-- Body -->
			<div style="padding:24px 32px;background:#fdfaf4;">
			<p style="font-size:14px;color:#4a3c20;line-height:1.8;margin:0 0 16px;">
				Kính gửi <strong>%s</strong>,
			</p>
			<p style="font-size:14px;color:#4a3c20;line-height:1.8;margin:0 0 24px;">
				Chúc mừng bạn đã đăng ký tài khoản thành công tại <strong>Saigon Kitchen</strong>. Dưới đây là thông tin tài khoản của bạn:
			</p>

			<!-- Info table -->
			<table style="width:100%%;border-collapse:collapse;font-size:13px;background:#fff;border:0.5px solid #e0d0b0;border-radius:8px;overflow:hidden;margin-bottom:24px;">
				<tr style="border-bottom:0.5px solid #f0e4c8;">
				<td style="padding:10px 14px;color:#8a7a5a;width:40%%;">Mã tài khoản</td>
				<td style="padding:10px 14px;color:#2a1f0a;font-family:'Courier New',monospace;font-weight:600;">#%d</td>
				</tr>
				<tr style="border-bottom:0.5px solid #f0e4c8;">
				<td style="padding:10px 14px;color:#8a7a5a;">Họ tên</td>
				<td style="padding:10px 14px;color:#2a1f0a;">%s</td>
				</tr>
				<tr>
				<td style="padding:10px 14px;color:#8a7a5a;">Email</td>
				<td style="padding:10px 14px;color:#2a1f0a;">%s</td>
				</tr>
			</table>

			<!-- Ghi chú -->
			<div style="background:#f7f0e3;border-left:3px solid #c4a55a;padding:12px 14px;margin-bottom:24px;">
				<p style="font-size:13px;color:#5a4520;margin:0;line-height:1.6;">
				Vui lòng <strong>bảo mật</strong> thông tin tài khoản và không chia sẻ mật khẩu với bất kỳ ai.
				</p>
			</div>

			<p style="font-size:13px;color:#6a5a3a;line-height:1.8;margin:0;">
				Nếu bạn không thực hiện đăng ký này, vui lòng liên hệ ngay với chúng tôi để được hỗ trợ.
			</p>
			</div>

			<!-- CTA -->
			<div style="padding:16px 32px;background:#fdfaf4;border-top:0.5px solid #e0d0b0;text-align:center;">
			<div style="display:inline-block;background:#1a1a1a;color:#e8d5b0;font-size:12px;letter-spacing:3px;padding:10px 28px;font-family:'Courier New',monospace;">
				KHÁM PHÁ THỰC ĐƠN NGAY!
			</div>
			</div>

			<!-- Footer -->
			<div style="padding:16px 32px;background:#1a1a1a;text-align:center;">
			<p style="font-size:11px;color:#6a5a3a;margin:0;letter-spacing:1px;font-family:'Courier New',monospace;">
				123 Đường ABC, Q.1, TP.HCM &nbsp;|&nbsp; 028-xxxx-xxxx
			</p>
			</div>

		</div>
		</body>
		</html>`,
		info.TenKhachHang,
		info.MaNguoiDung,
		info.TenKhachHang,
		info.Email,
	)

	return SendMail(email, "✦ Chào mừng bạn đến với Saigon Kitchen", body)
}

// type RegisterOTPData struct {
// 	Code      string
// 	ExpiredAt time.Time
// 	UserData  RegisterInput
// }

// var RegisterOTPStore = map[string]RegisterOTPData{}
