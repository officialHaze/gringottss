package util

import (
	"github.com/o1egl/paseto"
)

func PasetoEncrypt(jsonToken paseto.JSONToken, key, footer string) (string, error) {
	// Encrypt data
	token, err := paseto.NewV2().Encrypt([]byte(key), jsonToken, footer)
	if err != nil {
		return "", err
	}

	return token, nil
}

func PasetoDecrypt(token, key string) (paseto.JSONToken, string, error) {
	// Decrypt data
	var jsonToken paseto.JSONToken
	var footer string
	if err := paseto.NewV2().Decrypt(token, []byte(key), &jsonToken, &footer); err != nil {
		return paseto.JSONToken{}, "", err
	}

	return jsonToken, footer, nil
}
