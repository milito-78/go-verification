package go_verification

import (
	"context"
	"testing"
	"time"
)

func TestRedisCodeRepository(t *testing.T) {
	// Replace these values with your actual Redis configuration
	redisConfig := RedisConfig{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Prefix:   "test",
	}

	ctx := context.TODO()
	repo := NewRedisCodeRepository(ctx, redisConfig)

	username := "testuser"
	code := "123456"
	scope := "test_scope"
	expiresTime := 10 * time.Minute

	// Test SaveCode
	verification, err := repo.SaveCode(username, code, scope, expiresTime)
	if err != nil {
		t.Fatalf("SaveCode error: %v", err)
	}
	if verification.ExpiredTime != Duration(expiresTime) {
		t.Fatalf("ExpiredTime is not equals to input value")
	}
	// Test GetCode
	savedVerification, err := repo.GetCode(username, scope)
	if err != nil {
		t.Fatalf("GetCode error: %v", err)
	}

	if savedVerification.Code != code && savedVerification.Code != verification.Code {
		t.Errorf("Expected code to be %s, got %s", code, savedVerification.Code)
	}

	// Test DeleteCode
	if deleted := repo.DeleteCode(username, scope); !deleted {
		t.Error("DeleteCode failed to delete the code")
	}

	// Verify that the code is deleted
	_, err = repo.GetCode(username, scope)
	if err == nil {
		t.Error("GetCode expected to return an error after deletion")
	}
}

func TestRedisCodeRepository_DeleteAllCodes(t *testing.T) {
	// Replace these values with your actual Redis configuration
	redisConfig := RedisConfig{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Prefix:   "test",
	}

	ctx := context.TODO()
	repo := NewRedisCodeRepository(ctx, redisConfig)

	username := "testuser"
	scope1 := "test_scope1"
	scope2 := "test_scope2"
	code := "123456"
	expiresTime := 10 * time.Minute

	// Save codes with different scopes
	_, err := repo.SaveCode(username, code, scope1, expiresTime)
	if err != nil {
		t.Fatalf("SaveCode error: %v", err)
	}
	_, err = repo.SaveCode(username, code, scope2, expiresTime)
	if err != nil {
		t.Fatalf("SaveCode error: %v", err)
	}

	// Test DeleteAllCodes
	if deleted := repo.DeleteAllCodes(username); !deleted {
		t.Error("DeleteAllCodes failed to delete all codes")
	}

	// Verify that the codes are deleted
	_, err = repo.GetCode(username, scope1)
	if err == nil {
		t.Error("GetCode expected to return an error after deletion")
	}
	_, err = repo.GetCode(username, scope2)
	if err == nil {
		t.Error("GetCode expected to return an error after deletion")
	}
}
