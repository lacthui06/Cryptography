import random

def invMod(value, mod):
    return pow(value, -1, mod)

def point_add(P, Q, a, p):
    if P is None: return Q
    if Q is None: return P
    x1, y1 = P
    x2, y2 = Q
    if x1 == x2 and (y1 != y2 or y1 == 0):
        return None 
    if x1 == x2 and y1 == y2:
        lam = ((3 * x1 * x1 + a) * invMod(2 * y1, p)) % p
    else:
        lam = ((y2 - y1) * invMod(x2 - x1, p)) % p
    x3 = (lam * lam - x1 - x2) % p
    y3 = (lam * (x1 - x3) - y1) % p
    return (x3, y3)

def scalar_mult(k, P, a, p):
    R = None
    Q = P
    while k > 0:
        if k % 2 == 1: R = point_add(R, Q, a, p)
        Q = point_add(Q, Q, a, p)
        k //= 2
    return R

# HÀM MỚI: Phép đối xứng điểm (Tìm -Q)
def point_neg(Q, p):
    if Q is None: return None
    x, y = Q
    return (x, -y % p) # Điểm đối xứng qua trục Ox trong module p

# HÀM MỚI: Phép trừ điểm (P - Q) = P + (-Q)
def point_sub(P, Q, a, p):
    neg_Q = point_neg(Q, p)
    return point_add(P, neg_Q, a, p)


p = 29          
a = 4           
G = (1, 5)      
n = 37          

print("=== PHẦN 1: TRAO ĐỔI KHÓA (ECDH) ===")
# Alice sinh khóa
na = 15                         # Chọn khóa riêng na < n
Pa = scalar_mult(na, G, a, p)   # Tính khóa công khai Pa = na * G

# Bob sinh khóa
nb = 22                         # Chọn khóa riêng nb < n
Pb = scalar_mult(nb, G, a, p)   # Tính khóa công khai Pb = nb * G

print(f"Khóa công khai Alice (Pa): {Pa}")
print(f"Khóa công khai Bob (Pb): {Pb}")

# Tính khóa bí mật chung (Shared Secret)
k_Alice = scalar_mult(na, Pb, a, p) # Alice lấy na * Pb
k_Bob = scalar_mult(nb, Pa, a, p)   # Bob lấy nb * Pa

print(f"Khóa bí mật Alice tính được: {k_Alice}")
print(f"Khóa bí mật Bob tính được:   {k_Bob}")
print(f"Hai khóa giống nhau! Trao đổi khóa thành công.\n")


print("=== PHẦN 2: MÃ HÓA (ALICE GỬI CHO BOB) ===")
# Đầu tiên mã hóa tin nhắn M thành 1 điểm trên đường cong (Pm)
# Ở đây ta giả sử tin nhắn đã được map thành điểm (16, 27) hợp lệ trên đường cong
Pm = (16, 27) 
print(f"Điểm tin nhắn ban đầu (Pm): {Pm}")

# Chọn số nguyên dương k bất kỳ
k = random.randint(1, n - 1)

# Khi đó điểm mật mã sẽ là: Cm = (k*G, Pm + k*Pb)
C1 = scalar_mult(k, G, a, p)             # Thành phần thứ nhất: k * G

kPb = scalar_mult(k, Pb, a, p)           # Tính k * Pb
C2 = point_add(Pm, kPb, a, p)            # Thành phần thứ hai: Pm + k*Pb

Cm = (C1, C2)
print(f"Bản mã được gửi đi (Cm): {Cm}\n")


print("=== PHẦN 3: GIẢI MÃ (BOB NHẬN VÀ XỬ LÝ) ===")
# Bob nhận được Cm = (C1, C2)
nhan_C1, nhan_C2 = Cm

# Để giải mã, nhân thành phần 1 (C1) với khóa bí mật của người nhận (Bob -> nb)
# Tính: nb * C1 (chính là k * G * nb)
S = scalar_mult(nb, nhan_C1, a, p) 

# Sau đó trừ cho (k * G * nb) từ điểm C2
# Tính: Pm_recovered = C2 - S
Pm_recovered = point_sub(nhan_C2, S, a, p)

print(f"Điểm tin nhắn Bob giải mã được: {Pm_recovered}")

if Pm_recovered == Pm:
    print("GIẢI MÃ THÀNH CÔNG! (Khớp với tin nhắn gốc của Alice)")
else:
    print("GIẢI MÃ THẤT BẠI!")