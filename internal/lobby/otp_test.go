package lobby_test

import (
	"testing"

	"github.com/xGihyun/itso-quiz-bee/internal/lobby"
)

const LENGTH = 6

func TestGenerateOTP(t *testing.T) {
	otp, err := lobby.GenerateOTP(LENGTH)
	if err != nil {
		t.Fatal(err)
	}

	if LENGTH != len(otp) {
		t.Fatalf("Expected: %d, Got: %d - OTP: %s", LENGTH, len(otp), otp)
	}
}
