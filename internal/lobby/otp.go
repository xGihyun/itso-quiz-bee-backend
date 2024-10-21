package lobby

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
)

func GenerateOTP(length uint32) (string, error) {
	num, err := rand.Int(
		rand.Reader,
		big.NewInt(int64(math.Pow(10, float64(length)))),
	)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%0*d", length, num), nil
}
