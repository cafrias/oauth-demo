package security

import "encoding/base64"

var defaultCSRFTokenLength uint32 = 32

func GenerateCSRFToken() (string, error) {
	hash, err := generateRandomBytes(defaultCSRFTokenLength)
	if err != nil {
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(hash), nil
}
