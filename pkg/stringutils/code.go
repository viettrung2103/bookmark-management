package stringutils

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"math/big"

	"strings"

	"github.com/deatil/go-encoding/base62"
)

const (
	charset          = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	redisPrefixChars = "abcdefgh"
	dbPrefixChars    = "ijklmnopqrstuvwxyz"
)

// KeyGenerator is the interface for key generator
//
//go:generate mockery --name=KeyGenerator --filename=keygenerator.go
type KeyGenerator interface {
	GenerateKey(length int) string
	GenerateRedisKey(length int) string
	GenerateDBPrefix() string
	GenerateBase62Code(codeInt uint64) string
	IsRedisCode(code string) bool
	IsDBCode(code string) bool
}

type stringGenerator struct {
}

// NewKeyGenerator creates a new key generator
func NewKeyGenerator() KeyGenerator {
	return &stringGenerator{}
}

// GenerateKey generates a random string of the given length
func (r *stringGenerator) GenerateKey(length int) string {

	return randomString(length)
}

// GenerateRandomString generates a random string of the given length
func GenerateRandomString(length int) string {

	return randomString(length)

}

func randomString(length int) string {
	b := make([]byte, length)
	max := big.NewInt(int64(len(charset)))
	for i := 0; i < length; i++ {

		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic("crypto/rand failed: " + err.Error())
		}
		b[i] = charset[n.Int64()]
	}

	return string(b)
}

func generateRandomChar(allowedChars string) byte {
	max := big.NewInt(int64(len(allowedChars)))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic("crypto/rand failed: " + err.Error())
	}
	return allowedChars[n.Int64()]
}

// GenerateRedisKey generates a redis-specific code with prefix (a-h)
func (r *stringGenerator) GenerateRedisKey(length int) string {
	var strBuilder bytes.Buffer

	// 1. Get the first letter randomly from the redis prefix list (a-h)
	strBuilder.WriteByte(generateRandomChar(redisPrefixChars))
	// 2. Concat with the remaining random characters using your existing logic
	if length > 1 {
		strBuilder.WriteString(randomString(length - 1))
	}

	return strBuilder.String()
}

func (r *stringGenerator) GenerateDBPrefix() string {
	return string(generateRandomChar(dbPrefixChars))
}

func (r *stringGenerator) GenerateBase62Code(codeInt uint64) string {
	// Chuyển uint64 thành []byte trước
	idBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(idBytes, codeInt)

	// Encode sang chuỗi Base62
	base62Str := base62.StdEncoding.EncodeToString(idBytes)

	// 4. Ghép prefix do bạn tự định nghĩa (ví dụ ký tự ngẫu nhiên từ a-h hoặc i-z) với base62Str
	finalCode := r.GenerateDBPrefix() + base62Str
	return finalCode
}

func (r *stringGenerator) IsRedisCode(code string) bool {
	if len(code) == 0 {
		return false
	}
	return strings.Contains(redisPrefixChars, code[:1])
}
func (r *stringGenerator) IsDBCode(code string) bool {
	if len(code) == 0 {
		return false
	}
	return strings.Contains(dbPrefixChars, code[:1])

}
