package main

import "fmt"

var Sbox = [256]byte{
	0x63, 0x7c, 0x77, 0x7b, 0xf2, 0x6b, 0x6f, 0xc5, 0x30, 0x01, 0x67, 0x2b, 0xfe, 0xd7, 0xab, 0x76,
	0xca, 0x82, 0xc9, 0x7d, 0xfa, 0x59, 0x47, 0xf0, 0xad, 0xd4, 0xa2, 0xaf, 0x9c, 0xa4, 0x72, 0xc0,
	0xb7, 0xfd, 0x93, 0x26, 0x36, 0x3f, 0xf7, 0xcc, 0x34, 0xa5, 0xe5, 0xf1, 0x71, 0xd8, 0x31, 0x15,
	0x04, 0xc7, 0x23, 0xc3, 0x18, 0x96, 0x05, 0x9a, 0x07, 0x12, 0x80, 0xe2, 0xeb, 0x27, 0xb2, 0x75,
	0x09, 0x83, 0x2c, 0x1a, 0x1b, 0x6e, 0x5a, 0xa0, 0x52, 0x3b, 0xd6, 0xb3, 0x29, 0xe3, 0x2f, 0x84,
	0x53, 0xd1, 0x00, 0xed, 0x20, 0xfc, 0xb1, 0x5b, 0x6a, 0xcb, 0xbe, 0x39, 0x4a, 0x4c, 0x58, 0xcf,
	0xd0, 0xef, 0xaa, 0xfb, 0x43, 0x4d, 0x33, 0x85, 0x45, 0xf9, 0x02, 0x7f, 0x50, 0x3c, 0x9f, 0xa8,
	0x51, 0xa3, 0x40, 0x8f, 0x92, 0x9d, 0x38, 0xf5, 0xbc, 0xb6, 0xda, 0x21, 0x10, 0xff, 0xf3, 0xd2,
	0xcd, 0x0c, 0x13, 0xec, 0x5f, 0x97, 0x44, 0x17, 0xc4, 0xa7, 0x7e, 0x3d, 0x64, 0x5d, 0x19, 0x73,
	0x60, 0x81, 0x4f, 0xdc, 0x22, 0x2a, 0x90, 0x88, 0x46, 0xee, 0xb8, 0x14, 0xde, 0x5e, 0x0b, 0xdb,
	0xe0, 0x32, 0x3a, 0x0a, 0x49, 0x06, 0x24, 0x5c, 0xc2, 0xd3, 0xac, 0x62, 0x91, 0x95, 0xe4, 0x79,
	0xe7, 0xc8, 0x37, 0x6d, 0x8d, 0xd5, 0x4e, 0xa9, 0x6c, 0x56, 0xf4, 0xea, 0x65, 0x7a, 0xae, 0x08,
	0xba, 0x78, 0x25, 0x2e, 0x1c, 0xa6, 0xb4, 0xc6, 0xe8, 0xdd, 0x74, 0x1f, 0x4b, 0xbd, 0x8b, 0x8a,
	0x70, 0x3e, 0xb5, 0x66, 0x48, 0x03, 0xf6, 0x0e, 0x61, 0x35, 0x57, 0xb9, 0x86, 0xc1, 0x1d, 0x9e,
	0xe1, 0xf8, 0x98, 0x11, 0x69, 0xd9, 0x8e, 0x94, 0x9b, 0x1e, 0x87, 0xe9, 0xce, 0x55, 0x28, 0xdf,
	0x8c, 0xa1, 0x89, 0x0d, 0xbf, 0xe6, 0x42, 0x68, 0x41, 0x99, 0x2d, 0x0f, 0xb0, 0x54, 0xbb, 0x16,
}
var Rcon = [11]byte{
	0x00, 0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80, 0x1b, 0x36,
}

func charToWord(str string) []uint {
	b := []byte(str)
	word := make([]uint, len(b)/4)
	for i := 0; i < len(word); i++ {
		word[i] = uint(b[4*i])<<24 | uint(b[4*i+1])<<16 |
			uint(b[4*i+2])<<8 | uint(b[4*i+3])
	}
	return word
}
func subWord(word uint) uint {
	return uint(Sbox[byte(word>>24)])<<24 | uint(Sbox[byte(word>>16)])<<16 |
		uint(Sbox[byte(word>>8)])<<8 | uint(Sbox[byte(word)])
}
func rotWord(word uint) uint {
	return word>>24 | word<<8
}
func keyExpasion(keyWord []uint, Nk int, Nr int) []uint {
	totalWord := 4*Nr + 4
	w := make([]uint, totalWord)
	for i := 0; i < Nk; i++ {
		w[i] = keyWord[i]
	}
	for i := Nk; i < totalWord; i++ {
		temp := w[i-1]
		if i%Nk == 0 {
			temp = subWord(rotWord(temp)) ^ (uint(Rcon[i/Nk]) << 24)
		} else if Nk > 6 && i%Nk == 4 {
			temp = subWord(temp)
		}
		w[i] = w[i-Nk] ^ temp
	}
	return w
}
func shiftRow(state []uint) []uint {
	c0, c1, c2, c3 := state[0], state[1], state[2], state[3]
	state[0] = (c0 & 0xFF000000) | (c1 & 0x00FF0000) | (c2 & 0x0000FF00) | (c3 & 0x000000FF)
	state[1] = (c1 & 0xFF000000) | (c2 & 0x00FF0000) | (c3 & 0x0000FF00) | (c0 & 0x000000FF)
	state[2] = (c2 & 0xFF000000) | (c3 & 0x00FF0000) | (c0 & 0x0000FF00) | (c1 & 0x000000FF)
	state[3] = (c3 & 0xFF000000) | (c0 & 0x00FF0000) | (c1 & 0x0000FF00) | (c2 & 0x000000FF)
	fmt.Printf("ShiftRow:\t%08X\n", state)
	return state
}
func x02(b byte) byte {
	if b&0x80 == 0x80 {
		return (b << 1) ^ 0x1B
	} else {
		return b << 1
	}
}
func x03(b byte) byte { return b ^ x02(b) }

func mixColumn(state []uint) []uint {
	for i := 0; i < 4; i++ {
		r0, r1, r2, r3 := byte(state[i]>>24), byte(state[i]>>16),
			byte(state[i]>>8), byte(state[i])
		n0 := x02(r0) ^ x03(r1) ^ r2 ^ r3
		n1 := r0 ^ x02(r1) ^ x03(r2) ^ r3
		n2 := r0 ^ r1 ^ x02(r2) ^ x03(r3)
		n3 := x03(r0) ^ r1 ^ r2 ^ x02(r3)
		state[i] = uint(n0)<<24 | uint(n1)<<16 | uint(n2)<<8 | uint(n3)
	}
	fmt.Printf("Mix Column:\t%08X\n\n", state)
	return state
}
func cipher(inp []uint, Nr int, w []uint) []uint {
	state := make([]uint, len(inp))
	copy(state, inp)
	//addroundkey dau tien
	for i := 0; i < 4; i++ {
		state[i] ^= w[i]
	}
	for round := 1; round < Nr; round++ {
		//subword tung byte trong 1 word
		for i := 0; i < 4; i++ {
			state[i] = subWord(state[i])
		}
		state = shiftRow(state)

		state = mixColumn(state)
		for i := 0; i < 4; i++ {
			state[i] ^= w[4*round+i]
		}
	}
	for i := 0; i < 4; i++ {
		state[i] = subWord(state[i])
	}
	state = shiftRow(state)
	for i := 0; i < 4; i++ {
		state[i] ^= w[4*Nr+i]
	}
	return state
}
func cipherCTR(inp []uint, counter []uint, Nr int, w []uint) []uint {
	resultWords := make([]uint, len(inp))

	//copy counter vao counterBlock de counter ngoai ham giu nguyen ven
	counterBlock := make([]uint, 4)
	copy(counterBlock, counter)

	// Duyệt qua từng khối 4 word (16 bytes) của văn bản
	for i := 0; i < len(inp); i += 4 {
		tempCounter := make([]uint, 4)
		copy(tempCounter, counterBlock)

		keystream := cipher(tempCounter, Nr, w)

		// Bước 2: XOR Keystream với khối Plaintext hiện tại
		for j := 0; j < 4 && i+j < len(inp); j++ {
			resultWords[i+j] = inp[i+j] ^ keystream[j]
		}

		// Bước 3: Tăng bộ đếm (Counter) lên 1 cho khối tiếp theo
		// Trong thực tế, bộ đếm thường nằm ở word cuối cùng (index 3)
		counterBlock[3]++

		// Xử lý tràn số (Carry) nếu word cuối vượt quá giới hạn của uint32
		if counterBlock[3] == 0 {
			counterBlock[2]++
		}
	}

	return resultWords
}
func main() {
	blockStr := "TDHGTVTTPHCM-UTH"
	keyStr := "THAIANHLAC221106"
	counterStr := "TDHGTVTTPHCM-UTH"

	blockWord := charToWord(blockStr)
	keyWord := charToWord(keyStr)
	counterWord := charToWord(counterStr)

	fmt.Println("Plain text string:", blockStr)
	fmt.Println("Private key string:", keyStr)
	fmt.Printf("Block (Hex): %08X\n", blockWord)
	fmt.Printf("Key (Hex): %08X\n", keyWord)
	fmt.Printf("Counter (Hex): %08X\n", counterWord)

	keyLen := len(keyWord) * 4
	var Nk, Nr int
	switch keyLen {
	case 16:
		Nk = 4
		Nr = 10 // AES-128
	case 24:
		Nk = 6
		Nr = 12 // AES-192
	case 32:
		Nk = 8
		Nr = 14 // AES-256
	default:
		panic("Khóa không hợp lệ!")
	}

	fmt.Println(" MỞ RỘNG KHÓA ")
	w := keyExpasion(keyWord, Nk, Nr)
	ciCTR := cipherCTR(blockWord, counterWord, Nr, w)
	fmt.Printf("\nCipher text: %08X %08X %08X %08X\n\n",
		ciCTR[0], ciCTR[1], ciCTR[2], ciCTR[3])

	fmt.Println("\tGIẢI MÃ")
	plt := cipherCTR(ciCTR, counterWord, Nr, w)
	fmt.Printf("\nPlaint text: %08X %08X %08X %08X",
		plt[0], plt[1], plt[2], plt[3])
}
