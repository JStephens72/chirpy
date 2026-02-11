package auth_test

import (
	"testing"
	"time"

	"github.com/JStephens72/chirpy/internal/auth"
	"github.com/google/uuid"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "correct horse battery staple"

	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if hash == "" {
		t.Fatalf("expected non-empty hash")
	}

	// Should match
	match, err := auth.CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned error: %v", err)
	}
	if !match {
		t.Fatalf("expected password to match hash")
	}

	// Should NOT match
	match, err = auth.CheckPasswordHash("wrong password", hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned error on wrong password: %v", err)
	}
	if match {
		t.Fatalf("expected password mismatch")
	}
}

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "supersecretkey"
	expires := time.Minute * 5

	token, err := auth.MakeJWT(userID, secret, expires)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}
	if token == "" {
		t.Fatalf("expected non-empty token")
	}

	parsedID, err := auth.ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT returned error: %v", err)
	}

	if parsedID != userID {
		t.Fatalf("expected userID %v, got %v", userID, parsedID)
	}
}

func TestValidateJWT_InvalidSecret(t *testing.T) {
	userID := uuid.New()
	secret := "correctsecret"
	wrongSecret := "wrongsecret"

	token, err := auth.MakeJWT(userID, secret, time.Minute)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	_, err = auth.ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Fatalf("expected error when validating with wrong secret")
	}
}

func TestValidateJWT_Expired(t *testing.T) {
	userID := uuid.New()
	secret := "secret"

	// Token expired 1 second ago
	token, err := auth.MakeJWT(userID, secret, -1*time.Second)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	_, err = auth.ValidateJWT(token, secret)
	if err == nil {
		t.Fatalf("expected error validating expired token")
	}
}
