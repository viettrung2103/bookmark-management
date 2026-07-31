package stringutils

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"strings"
	"time"

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
	rng *rand.Rand
}

// NewKeyGenerator creates a new key generator
func NewKeyGenerator() KeyGenerator {
	return &stringGenerator{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GenerateKey generates a random string of the given length
func (r *stringGenerator) GenerateKey(length int) string {

	return randomString(r.rng, length)
}

// GenerateRandomString generates a random string of the given length
func GenerateRandomString(length int) string {

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return randomString(rng, length)

}

func randomString(rng *rand.Rand, length int) string {
	var strBuilder bytes.Buffer

	for i := 0; i < length; i++ {

		strBuilder.WriteByte(charset[rng.Intn(len(charset))])
	}

	return strBuilder.String()
}

// GenerateRedisKey generates a redis-specific code with prefix (a-h)
func (r *stringGenerator) GenerateRedisKey(length int) string {
	var strBuilder bytes.Buffer

	// 1. Get the first letter randomly from the redis prefix list (a-h)
	strBuilder.WriteByte(redisPrefixChars[r.rng.Intn(len(redisPrefixChars))])

	// 2. Concat with the remaining random characters using your existing logic
	if length > 1 {
		strBuilder.WriteString(randomString(r.rng, length-1))
	}

	return strBuilder.String()
}

func (r *stringGenerator) GenerateDBPrefix() string {
	return string(dbPrefixChars[r.rng.Intn(len(dbPrefixChars))])
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
	return strings.Contains(redisPrefixChars, code[:1])
}
func (r *stringGenerator) IsDBCode(code string) bool {
	return strings.Contains(dbPrefixChars, code[:1])

}
