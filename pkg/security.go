package pkg

import (
	"encoding/hex"
	"math/rand"

	"golang.org/x/crypto/scrypt"
)

func generateRandString(length int) string {
	var chars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+,.?/:;{}[]`~")
	clen := len(chars)
	if clen < 2 || clen > 256 {
		panic("Wrong charset length for NewLenChars()")
	}
	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			panic("Error reading random bytes: " + err.Error())
		}
		for _, rb := range r {
			c := int(rb)
			if c > maxrb {
				continue // Skip this number to avoid modulo bias.
			}
			b[i] = chars[c%clen]
			i++
			if i == length {
				return string(b)
			}
		}
	}
}

// GenerateRandomKey20 生成20位随机字符串
func GenerateRandomKey20() string {
	return generateRandString(20)
}

// GenerateRandomKey16 生成6为随机字符串
func GenerateRandomKey16() string {
	return generateRandString(6)
}

// SetPassword 根据明文密码和加盐值生成密码
func SetPassword(password string, salt string) (verify string, err error) {
	var rb []byte
	rb, err = scrypt.Key([]byte(password), []byte(salt), 16384, 8, 1, 32)
	if err != nil {
		return
	}
	verify = hex.EncodeToString(rb)
	return
}
