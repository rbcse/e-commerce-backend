package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateRandom6DigitOTP() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	return fmt.Sprintf("%06d", n.Int64()), err
}