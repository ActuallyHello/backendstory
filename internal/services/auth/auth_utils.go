package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type PasswordConfig struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
	saltLen uint32
}

var DefaultPasswordConfig = &PasswordConfig{
	time:    1,
	memory:  64 * 1024,
	threads: 4,
	keyLen:  32,
	saltLen: 16,
}

// GenerateHash создает хэш пароля используя Argon2id
func GenerateHash(password string) (string, error) {
	config := DefaultPasswordConfig

	// Генерируем случайную соль
	salt := make([]byte, config.saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Создаем хэш используя Argon2id
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		config.time,
		config.memory,
		config.threads,
		config.keyLen,
	)

	// Кодируем в строку формата: argon2id$time$memory$threads$keyLen$salt$hash
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf(
		"argon2id$%d$%d$%d$%d$%s$%s",
		config.time,
		config.memory,
		config.threads,
		config.keyLen,
		b64Salt,
		b64Hash,
	)

	return encoded, nil
}

// VerifyPassword проверяет пароль против хэша
func VerifyPassword(password, encodedHash string) (bool, error) {
	// Парсим encoded hash
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 7 {
		return false, fmt.Errorf("invalid hash format")
	}
	if parts[0] != "argon2id" {
		return false, fmt.Errorf("unsupported hash algorithm")
	}

	var config PasswordConfig
	if _, err := fmt.Sscanf(parts[1], "%d", &config.time); err != nil {
		return false, err
	}
	if _, err := fmt.Sscanf(parts[2], "%d", &config.memory); err != nil {
		return false, err
	}
	if _, err := fmt.Sscanf(parts[3], "%d", &config.threads); err != nil {
		return false, err
	}
	if _, err := fmt.Sscanf(parts[4], "%d", &config.keyLen); err != nil {
		return false, err
	}

	// Декодируем соль и хэш
	salt, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[6])
	if err != nil {
		return false, err
	}

	// Вычисляем хэш для предоставленного пароля
	actualHash := argon2.IDKey(
		[]byte(password),
		salt,
		config.time,
		config.memory,
		config.threads,
		config.keyLen,
	)

	// Сравниваем хэши безопасно (constant-time comparison)
	if subtle.ConstantTimeCompare(actualHash, expectedHash) == 1 {
		return true, nil
	}

	return false, nil
}

// IsPasswordHashed проверяет, является ли строка хэшированным паролем
func IsPasswordHashed(password string) bool {
	return strings.HasPrefix(password, "argon2id$")
}
