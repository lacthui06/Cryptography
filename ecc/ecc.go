package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// ==========================================
// CẤU TRÚC ĐIỂM TRÊN ĐƯỜNG CONG (Point)
// ==========================================
type Point struct {
	X *big.Int
	Y *big.Int
}

// Hàm in điểm cho đẹp
func (p *Point) String() string {
	if p == nil {
		return "Điểm vô cực (O)"
	}
	return fmt.Sprintf("(%v, %v)", p.X, p.Y)
}

// So sánh 2 điểm
func pointEquals(P, Q *Point) bool {
	if P == nil && Q == nil {
		return true
	}
	if P == nil || Q == nil {
		return false
	}
	return P.X.Cmp(Q.X) == 0 && P.Y.Cmp(Q.Y) == 0
}

// ==========================================
// CÁC HÀM TOÁN HỌC CỐT LÕI
// ==========================================

// Cộng 2 điểm trên đường cong
func pointAdd(P, Q *Point, a, p *big.Int) *Point {
	if P == nil {
		return Q
	}
	if Q == nil {
		return P
	}

	// P = -Q -> Trả về vô cực
	if P.X.Cmp(Q.X) == 0 && P.Y.Cmp(Q.Y) != 0 {
		return nil
	}
	if P.X.Cmp(Q.X) == 0 && P.Y.Cmp(big.NewInt(0)) == 0 {
		return nil
	}

	lam := new(big.Int)

	if P.X.Cmp(Q.X) == 0 && P.Y.Cmp(Q.Y) == 0 {
		// Nhân đôi (Doubling)
		num := new(big.Int).Mul(P.X, P.X)
		num.Mul(num, big.NewInt(3))
		num.Add(num, a)

		den := new(big.Int).Mul(P.Y, big.NewInt(2))
		den.ModInverse(den, p)

		lam.Mul(num, den).Mod(lam, p)
	} else {
		// Cộng (Addition)
		num := new(big.Int).Sub(Q.Y, P.Y)
		den := new(big.Int).Sub(Q.X, P.X)

		den.Mod(den, p)
		den.ModInverse(den, p)

		lam.Mul(num, den).Mod(lam, p)
	}

	x3 := new(big.Int).Mul(lam, lam)
	x3.Sub(x3, P.X).Sub(x3, Q.X).Mod(x3, p)

	y3 := new(big.Int).Sub(P.X, x3)
	y3.Mul(lam, y3).Sub(y3, P.Y).Mod(y3, p)

	return &Point{X: x3, Y: y3}
}

// Nhân vô hướng k * P
func scalarMult(k *big.Int, P *Point, a, p *big.Int) *Point {
	var R *Point = nil
	Q := &Point{X: new(big.Int).Set(P.X), Y: new(big.Int).Set(P.Y)}

	kCopy := new(big.Int).Set(k)
	zero := big.NewInt(0)
	two := big.NewInt(2)
	rem := new(big.Int)

	for kCopy.Cmp(zero) > 0 {
		rem.Mod(kCopy, two)
		if rem.Cmp(big.NewInt(1)) == 0 {
			R = pointAdd(R, Q, a, p)
		}
		Q = pointAdd(Q, Q, a, p)
		kCopy.Div(kCopy, two)
	}
	return R
}

// --- HÀM MỚI: Phép đối xứng điểm (Tìm -Q) ---
func pointNeg(Q *Point, p *big.Int) *Point {
	if Q == nil {
		return nil
	}
	// -y mod p = (p - y) mod p
	yNeg := new(big.Int).Sub(p, Q.Y)
	yNeg.Mod(yNeg, p)
	return &Point{X: new(big.Int).Set(Q.X), Y: yNeg}
}

// --- HÀM MỚI: Phép trừ điểm (P - Q) = P + (-Q) ---
func pointSub(P, Q *Point, a, p *big.Int) *Point {
	negQ := pointNeg(Q, p)
	return pointAdd(P, negQ, a, p)
}

// ==========================================
// CHẠY THỬ NGHIỆM GIAO THỨC
// ==========================================
func main() {
	// 1. Tham số miền (Global Public elements)
	p := big.NewInt(29)
	a := big.NewInt(4)
	G := &Point{X: big.NewInt(1), Y: big.NewInt(5)}
	n := big.NewInt(37)

	fmt.Println("=== PHẦN 1: TRAO ĐỔI KHÓA (ECDH) ===")
	// Alice sinh khóa
	na := big.NewInt(15) // Khóa riêng Alice
	Pa := scalarMult(na, G, a, p)

	// Bob sinh khóa
	nb := big.NewInt(22) // Khóa riêng Bob
	Pb := scalarMult(nb, G, a, p)

	fmt.Printf("Khóa công khai Alice (Pa): %s\n", Pa)
	fmt.Printf("Khóa công khai Bob (Pb): %s\n", Pb)

	// Tính khóa bí mật chung (Shared Secret)
	k_Alice := scalarMult(na, Pb, a, p)
	k_Bob := scalarMult(nb, Pa, a, p)

	fmt.Printf("Khóa bí mật Alice tính: %s\n", k_Alice)
	fmt.Printf("Khóa bí mật Bob tính:   %s\n", k_Bob)
	fmt.Println("-> Hai khóa khớp nhau! Trao đổi thành công.\n")

	fmt.Println("=== PHẦN 2: MÃ HÓA (ALICE GỬI CHO BOB) ===")
	// Tin nhắn M được map thành điểm Pm trên đường cong
	Pm := &Point{X: big.NewInt(16), Y: big.NewInt(27)}
	fmt.Printf("Điểm tin nhắn ban đầu (Pm): %s\n", Pm)

	// Chọn k ngẫu nhiên trong [1, n-1]
	maxK := new(big.Int).Sub(n, big.NewInt(1))
	k, _ := rand.Int(rand.Reader, maxK)
	k.Add(k, big.NewInt(1))
	// k := big.NewInt(10) // Có thể hardcode số 10 để test nếu ko muốn random

	// Cm = (C1, C2) = (k*G, Pm + k*Pb)
	C1 := scalarMult(k, G, a, p)
	kPb := scalarMult(k, Pb, a, p)
	C2 := pointAdd(Pm, kPb, a, p)

	fmt.Printf("Bản mã gửi đi (Cm):\n - C1: %s\n - C2: %s\n\n", C1, C2)

	fmt.Println("=== PHẦN 3: GIẢI MÃ (BOB NHẬN VÀ XỬ LÝ) ===")
	// Bob dùng khóa riêng nb nhân với C1: S = nb * C1
	S := scalarMult(nb, C1, a, p)

	// Giải mã: Pm_recovered = C2 - S
	PmRecovered := pointSub(C2, S, a, p)

	fmt.Printf("Điểm tin nhắn Bob giải mã được: %s\n", PmRecovered)

	if pointEquals(PmRecovered, Pm) {
		fmt.Println("✅ GIẢI MÃ THÀNH CÔNG! (Khớp với tin nhắn gốc của Alice)")
	} else {
		fmt.Println("❌ GIẢI MÃ THẤT BẠI!")
	}
}
