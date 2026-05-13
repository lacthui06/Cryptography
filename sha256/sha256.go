package main

import (
	"fmt"
)

var K = []uint32{
	0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5, 0x3956c25b, 0x59f111f1, 0x923f82a4, 0xab1c5ed5,
	0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3, 0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174,
	0xe49b69c1, 0xefbe4786, 0x0fc19dc6, 0x240ca1cc, 0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da,
	0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7, 0xc6e00bf3, 0xd5a79147, 0x06ca6351, 0x14292967,
	0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13, 0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85,
	0xa2bfe8a1, 0xa81a664b, 0xc24b8b70, 0xc76c51a3, 0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070,
	0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5, 0x391c0cb3, 0x4ed8aa4a, 0x5b9cca4f, 0x682e6ff3,
	0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208, 0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2,
}

var HInit = []uint32{
	0x6a09e667, 0xbb67ae85, 0x3c6ef372, 0xa54ff53a,
	0x510e527f, 0x9b05688c, 0x1f83d9ab, 0x5be0cd19,
}

func ROTR(x uint32, n int) uint32 { return (x >> n) | (x << (32 - n)) }
func SHR(x uint32, n int) uint32  { return x >> n }

func Ch(x, y, z uint32) uint32  { return (x & y) ^ (^x & z) }
func Maj(x, y, z uint32) uint32 { return (x & y) ^ (x & z) ^ (y & z) }

func Sigma0_Lon(x uint32) uint32 { return ROTR(x, 2) ^ ROTR(x, 13) ^ ROTR(x, 22) }
func Sigma1_Lon(x uint32) uint32 { return ROTR(x, 6) ^ ROTR(x, 11) ^ ROTR(x, 25) }

func sigma0_Nho(x uint32) uint32 { return ROTR(x, 7) ^ ROTR(x, 18) ^ SHR(x, 3) }
func sigma1_Nho(x uint32) uint32 { return ROTR(x, 17) ^ ROTR(x, 19) ^ SHR(x, 10) }

// 3. HÀM XỬ LÝ CHÍNH
func sha256(message []byte) string {
	// --- BƯỚC 1: TỰ TÍNH TOÁN PADDING & CẤP PHÁT BỘ NHỚ ---
	L := len(message)
	origLenBits := uint64(L) * 8

	paddedLen := L + 1
	for paddedLen%64 != 56 {
		paddedLen++
	}
	paddedLen += 8

	paddedMessage := make([]byte, paddedLen)

	// Chép message gốc
	for i := 0; i < L; i++ {
		paddedMessage[i] = message[i]
	}

	// Đệm bit 1 (0x80)
	paddedMessage[L] = 0x80

	paddedMessage[paddedLen-8] = byte(origLenBits >> 56)
	paddedMessage[paddedLen-7] = byte(origLenBits >> 48)
	paddedMessage[paddedLen-6] = byte(origLenBits >> 40)
	paddedMessage[paddedLen-5] = byte(origLenBits >> 32)
	paddedMessage[paddedLen-4] = byte(origLenBits >> 24)
	paddedMessage[paddedLen-3] = byte(origLenBits >> 16)
	paddedMessage[paddedLen-2] = byte(origLenBits >> 8)
	paddedMessage[paddedLen-1] = byte(origLenBits)

	// --- BƯỚC 2: KHỞI TẠO BIẾN TRẠNG THÁI ---
	H := make([]uint32, 8)
	for i := 0; i < 8; i++ {
		H[i] = HInit[i]
	}

	// --- BƯỚC 3 & 4: XỬ LÝ TỪNG BLOCK ---
	for i := 0; i < len(paddedMessage); i += 64 {
		// Tạo Message Schedule W
		var W [64]uint32

		for t := 0; t < 16; t++ {
			idx := i + (t * 4) // Lấy vị trí trực tiếp trên paddedMessage
			// Tự rã 4 bytes thành 1 khối uint32 (BigEndian)
			W[t] = (uint32(paddedMessage[idx]) << 24) |
				(uint32(paddedMessage[idx+1]) << 16) |
				(uint32(paddedMessage[idx+2]) << 8) |
				uint32(paddedMessage[idx+3])
		}

		for t := 16; t < 64; t++ {
			W[t] = sigma1_Nho(W[t-2]) + W[t-7] + sigma0_Nho(W[t-15]) + W[t-16]
		}

		a, b, c, d, e, f, g, h := H[0], H[1], H[2], H[3], H[4], H[5], H[6], H[7]

		// 64 vòng lặp nén
		for t := 0; t < 64; t++ {
			T1 := h + Sigma1_Lon(e) + Ch(e, f, g) + K[t] + W[t]
			T2 := Sigma0_Lon(a) + Maj(a, b, c)
			fmt.Printf("Vòng %02d | Maj:%x Sigma0_Lon:%x T1:%x T2:%08x \n", t+1, Maj(a, b, c), Sigma0_Lon(a), T1, T2)
			h = g
			g = f
			f = e
			e = d + T1
			d = c
			c = b
			b = a
			a = T1 + T2
			fmt.Printf("%08x %08x %08x %08x %08x %08x %08x %08x\n\n",
				a, b, c, d, e, f, g, h)
		}

		// Cập nhật giá trị băm sau mỗi block
		H[0] += a
		H[1] += b
		H[2] += c
		H[3] += d
		H[4] += e
		H[5] += f
		H[6] += g
		H[7] += h
	}

	// --- BƯỚC 5: XUẤT KẾT QUẢ DẠNG HEX ---
	return fmt.Sprintf("%08x%08x%08x%08x%08x%08x%08x%08x", H[0], H[1], H[2], H[3], H[4], H[5], H[6], H[7])
}

// 4. HÀM MAIN CHẠY THỬ
func main() {
	plaintext := "TDHGTVTTPHCM-UTH"
	fmt.Printf("Thông điệp gốc: %s\n", plaintext)

	hash := sha256([]byte(plaintext))
	fmt.Printf("Mã băm SHA-256: %s\n", hash)
}
