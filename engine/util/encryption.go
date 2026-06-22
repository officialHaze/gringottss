package util

import (
	"crypto/rand"
	"math/big"
	"os"
	"path"

	"github.com/goccy/go-yaml"
	"github.com/o1egl/paseto"
	"github.com/officialhaze/gringottss/api-server/logger"
)

var (
	LoadedEncryptionKeys EncryptionKeys
)

type EncryptionKeys struct {
	PWD_ENCRYPTION_KEY string `yaml:"PWD_ENCRYPTION_KEY"`
}

func LoadEncryptionKeys() error {
	encryptionKeyFile := path.Join("encryption_keys.yml")
	f, err := os.Open(encryptionKeyFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// Decode the file
	var keys EncryptionKeys
	dec := yaml.NewDecoder(f)
	dec.Decode(&keys)

	// Set the decoded data for global use
	LoadedEncryptionKeys = keys
	return nil
}

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

// Generate 32 bytes string (paseto supported)
func Generate32ByteString() string {
	allowedCharset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, 32)
	charsetLength := big.NewInt(int64(len(allowedCharset)))

	for i := 0; i < 32; i++ {
		// pick any random char from charset
		randInt, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			logger.ERROR().Println(err.Error())
			return ""
		}

		b[i] = allowedCharset[randInt.Int64()]
	}

	// Convert the bytes to hex
	return string(b)
}
