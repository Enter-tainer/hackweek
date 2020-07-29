package util

import (
	"encoding/base64"
	"tree-hole/config"

	"golang.org/x/crypto/argon2"
)

func PasswordHash(password string) string {
	return base64.StdEncoding.EncodeToString(argon2.IDKey([]byte(password), []byte(config.Config.App.Salt), 1, 64*1024, 4, 32))
}
