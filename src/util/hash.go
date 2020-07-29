package util

import (
	"encoding/base64"
	"tree-hole/config"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/blake2b"
)

func PasswordHash(password string) string {
	return base64.StdEncoding.EncodeToString(argon2.IDKey([]byte(password), []byte(config.Config.App.Salt), 1, 64*1024, 4, 32))
}

func UserIDHash(salt string, userID string) string {
	res := blake2b.Sum256([]byte(salt + userID))
	return base64.StdEncoding.EncodeToString(res[:])
}
