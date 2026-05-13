import hashlib #đây là thuật toán chữ kí số nên dùng hàm hash demo
import random

def invMod(value, mod):
    return pow(value, -1, mod)

def add2Point(P, Q, a, p):
    if P is None: return Q
    if Q is None: return P
    
    x1, y1 = P
    x2, y2 = Q
    
    if x1 == x2 and (y1 != y2 or y1 == 0):
        return None # Điểm vô cực (Point at Infinity)
        
    if x1 == x2 and y1 == y2:
        lam = ((3 * x1 * x1 + a) * invMod(2 * y1, p)) % p
    else:
        lam = ((y2 - y1) * invMod(x2 - x1, p)) % p
        
    x3 = (lam * lam - x1 - x2) % p
    y3 = (lam * (x1 - x3) - y1) % p
    return (x3, y3)

def scalarMult(k, P, a, p):
    R = None
    Q = P
    while k > 0:
        if k % 2 == 1:
            R = add2Point(R, Q, a, p)
        Q = add2Point(Q, Q, a, p)
        k //= 2
    return R


# THUẬT TOÁN TẠO CHỮ KÝ
def ecdsaSign(message, d, G, n, p, a):
    hashytes = hashlib.sha256(message.encode('utf-8')).digest()
    e = int.from_bytes(hashytes, byteorder='big')
    
    while True:
        k = random.SystemRandom().randint(1, n - 1)
        P = scalarMult(k, G, a, p)
        if P is None: continue
            
        r = P[0] % n
        if r == 0: continue
            
        t = invMod(k, n)
        s = (t * (e + d * r)) % n
        if s == 0: continue
            
        return (r, s)


# THUẬT TOÁN XÁC THỰC CHỮ KÝ
def ecdsaVerify(message, signature, Q, G, n, p, a):
    """
    signature: Cặp (r, s)
    Q: Khóa công khai của người ký (Public Key)
    """
    r, s = signature
    
    # BƯỚC 1: Xác minh r và s nằm trong phạm vi [1, n - 1]
    if not (1 <= r <= n - 1) or not (1 <= s <= n - 1):
        print("Xác thực thất bại: r hoặc s không nằm trong khoảng hợp lệ.")
        return False
        
    # BƯỚC 2: Tính Hash(m) và chuyển thành số nguyên e
    hashytes = hashlib.sha256(message.encode('utf-8')).digest()
    e = int.from_bytes(hashytes, byteorder='big')
    
    # BƯỚC 3: Tính toán w = s^-1 mod n
    w = invMod(s, n)
    
    # BƯỚC 4: Tính toán u1 = ew (mod n) và u2 = rw (mod n)
    u1 = (e * w) % n
    u2 = (r * w) % n
    
    # BƯỚC 5: Tính toán điểm X(x1, y1) = u1*G + u2*Q
    u1 = scalarMult(u1, G, a, p)
    u2 = scalarMult(u2, Q, a, p)
    X = add2Point(u1, u2, a, p)
    
    # BƯỚC 6: Nếu X = 0 (Điểm vô cực), từ chối chữ ký.
    if X is None:
        print("Xác thực thất bại: Điểm X là điểm vô cực.")
        return False
        
    # Nếu X != 0, tính v = x1 mod n
    x1 = X[0]
    v = x1 % n
    
    # BƯỚC 7: Chấp nhận chữ ký khi và chỉ khi v = r
    if v == r:
        print("Xác thực thành công! (v == r)")
        return True
    else:
        print(f"Xác thực thất bại! (v = {v} khác r = {r})")
        return False


'''CHẠY THỬ NGHIỆM ĐẦY ĐỦ (KÝ & XÁC THỰC)
Đường cong y^2 = x^3 + ax + b mod 29'''
p = 29
a = 4 #b = 20
G = (1, 5)
n = 37 # must prime
# Tạo cặp khóa
d = 7           # Khóa riêng tư (Alice giữ bí mật)
# Khóa công khai Q = d * G (Alice công bố cho Bob)
Q = scalarMult(d, G, a, p)


m = "TDHGTVTTPHCM-UTH"
print(f"--- KÝ TIN NHẮN ---")
signature = ecdsaSign(m, d, G, n, p, a)
print(f"Chữ ký (r, s): {signature}\n")

print(f"--- XÁC THỰC ---")
isalid = ecdsaVerify(m, signature, Q, G, n, p, a)


# Thử nghiệm giả mạo (Sửa đổi nội dung tin nhắn)
print(f"\n--- THỬ NGHIỆM GIẢ MẠO TIN NHẮN ---")
fakeessage = "Hello Bob, send me $1000. - Alice."
ecdsaVerify(fakeessage, signature, Q, G, n, p, a)