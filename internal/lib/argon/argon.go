package argon

import (
	"bytes"
	"crypto/rand"
	"errors"
	"golang.org/x/crypto/argon2"
)

type HashSalt struct {
	Hash, Salt []byte
}

// Структура конфига argon2id
type Argon struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
	saltLen uint32
}

// Конструктор для инициализации структуры
func NewArgonHash(time, saltLen, memory, keyLen uint32, threads uint8) *Argon {
	return &Argon{
		time:    time,
		memory:  memory,
		threads: threads,
		keyLen:  keyLen,
		saltLen: saltLen,
	}
}

// Генерация соли
func randomSecret(length uint32) ([]byte, error) {
	secret := make([]byte, length)

	_, err := rand.Read(secret)
	if err != nil {
		return nil, err
	}
	return secret, nil
}

// Метод хэширования GenerateHash структуры Argon:
func (a *Argon) GenerateHash(password, salt []byte) (*HashSalt, error) {
	var err error
	if len(salt) == 0 {
		salt, err = randomSecret(a.saltLen)
	}
	if err != nil {
		return nil, err
	}
	hash := argon2.IDKey(password, salt, a.time, a.memory, a.threads, a.keyLen)
	return &HashSalt{Hash: hash, Salt: salt}, nil
}

// Сравнение паролей
func (a *Argon) Compare(hash, salt, password []byte) error {
	hashSalt, err := a.GenerateHash(password, salt)
	if err != nil {
		return err
	}
	if !bytes.Equal(hash, hashSalt.Hash) {
		return errors.New("hash does not match")
	}
	return nil
}
