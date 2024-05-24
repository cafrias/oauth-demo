package utils

import (
	"crypto/rand"
	"math/big"
)

// RandHexDecString generates a cryptographycally secure random hexdecimal string of length n
func RandHexDecString(n int) (string, error) {
	const letters = "0123456789abcdef"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}
