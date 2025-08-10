package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userId := uuid.New()
	tokenSecret := "hello"
	expiresIn := time.Minute
	tokenString, err := MakeJWT(userId, tokenSecret, expiresIn)
	if err != nil {
		t.Fatal("error making JWT")
	}

	newUserId, err := ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Fatal("invalid tokenString or tokenSecret")
	}

	if newUserId != userId {
		t.Error("userId is not equal to returned userID")
	}
}

func TextExpiredJWT(t *testing.T) {
	userId := uuid.New()
	tokenSecret := "hello"
	expiresIn := time.Second
	tokenString, err := MakeJWT(userId, tokenSecret, expiresIn)
	if err != nil {
		t.Fatal("error making JWT")
	}
	time.Sleep(2 * expiresIn)

	_, err = ValidateJWT(tokenString, tokenSecret)
	if err == nil {
		t.Fatal("valid tokens, time did not expire")
	}
}

func TestInValidateJWT(t *testing.T) {
	userId := uuid.New()
	knownSecret := "correctSecret"
	expiresIn := time.Minute
	tokenString, err := MakeJWT(userId, knownSecret, expiresIn)
	if err != nil {
		t.Fatal("error making JWT")
	}

	unknownSecret := "WrongSecret"
	_, err = ValidateJWT(tokenString, unknownSecret)
	if err == nil {
		t.Fatal("validated wrong secret")
	}
}