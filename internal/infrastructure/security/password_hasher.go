package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// PasswordHasher Argon2id密码哈希器实现
type PasswordHasher struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

// NewPasswordHasher 创建密码哈希器
func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{
		memory:      64 * 1024, // 64 MB
		iterations:  3,         // 3 iterations
		parallelism: 2,         // 2 threads
		saltLength:  16,        // 16 bytes salt
		keyLength:   32,        // 32 bytes key
	}
}

// HashPassword 哈希密码
func (p *PasswordHasher) HashPassword(password string) (string, error) {
	// 生成随机盐
	salt, err := p.generateRandomBytes(p.saltLength)
	if err != nil {
		return "", fmt.Errorf("生成盐失败: %w", err)
	}

	// 使用Argon2id生成哈希
	hash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// 编码为base64字符串
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// 格式: $argon2id$v=19$m=65536,t=3,p=2$salt$hash
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, p.memory, p.iterations, p.parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

// VerifyPassword 验证密码
func (p *PasswordHasher) VerifyPassword(hashedPassword, password string) bool {
	// 解析编码的哈希
	salt, hash, memory, iterations, parallelism, err := p.decodeHash(hashedPassword)
	if err != nil {
		return false
	}

	// 使用相同参数重新计算哈希
	otherHash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, p.keyLength)

	// 使用constant time比较防止时序攻击
	return subtle.ConstantTimeCompare(hash, otherHash) == 1
}

// generateRandomBytes 生成随机字节
func (p *PasswordHasher) generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// decodeHash 解码哈希字符串
func (p *PasswordHasher) decodeHash(encodedHash string) (salt, hash []byte, memory uint32, iterations uint32, parallelism uint8, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, 0, 0, 0, fmt.Errorf("无效的哈希格式")
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, 0, 0, 0, err
	}
	if version != argon2.Version {
		return nil, nil, 0, 0, 0, fmt.Errorf("不兼容的argon2版本")
	}

	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return nil, nil, 0, 0, 0, err
	}

	salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, 0, 0, 0, err
	}

	hash, err = base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, 0, 0, 0, err
	}

	return salt, hash, memory, iterations, parallelism, nil
}
