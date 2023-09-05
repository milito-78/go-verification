package go_verification

import (
	"errors"
	"math/rand"
	"testing"
	"time"
)

type MockCodeRepository struct {
	data map[string]*VerificationCode
}

func NewMockCodeRepository() *MockCodeRepository {
	return &MockCodeRepository{
		data: make(map[string]*VerificationCode),
	}
}

func (m *MockCodeRepository) DeleteAllCodes(username string) bool {
	return false
}

func (m *MockCodeRepository) SaveCode(username, code, scope string, expiresTime time.Duration) (*VerificationCode, error) {
	key := username + scope
	expiredAt := time.Now().Add(expiresTime)
	data := &VerificationCode{
		ExpireAfter: int(expiresTime.Seconds()),
		ExpiredTime: Duration(expiresTime),
		ExpiredAt:   expiredAt,
		Username:    username,
		Scope:       scope,
		Code:        code,
	}
	m.data[key] = data
	return data, nil
}

func (m *MockCodeRepository) GetCode(username, scope string) (*VerificationCode, error) {
	key := username + scope
	data, ok := m.data[key]
	if !ok {
		return nil, errors.New("does not exist")
	}
	return data, nil
}

func (m *MockCodeRepository) DeleteCode(username, scope string) bool {
	key := username + scope
	_, ok := m.data[key]
	if !ok {
		return false
	}
	delete(m.data, key)
	return true
}

type MockCodeGenerator struct {
	defCode string
	length  int
}

func (m *MockCodeGenerator) Generate() string {
	if m.defCode != "" {
		return m.defCode
	}
	rand.Seed(time.Now().UnixNano())
	chars := "0123456789"
	result := make([]byte, m.length)
	for i := 0; i < m.length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func TestVerificationCodeHandler(t *testing.T) {
	options := &Config{
		ExpiredAfterSec: 5 * time.Minute,
	}
	code := "123456"

	repository := NewMockCodeRepository()
	generator := &MockCodeGenerator{defCode: code}

	handler, err := NewVerificationCodeHandler(generator, repository, options)
	if err != nil {
		t.Fatalf("Failed to create VerificationCodeHandler: %v", err)
	}

	username := "testuser"
	scope := "testscope"

	// Test GenerateCode
	verification, err := handler.GenerateCode(username, scope)
	if err != nil {
		t.Fatalf("GenerateCode error: %v", err)
	}
	if verification.Code != code {
		t.Errorf("Expected code to be %s, got %s", code, verification.Code)
	}

	// Test GetCode
	retrievedVerification, err := handler.GetCode(username, scope)
	if err != nil {
		t.Fatalf("GetCode error: %v", err)
	}
	if retrievedVerification.Code != code {
		t.Errorf("Expected code to be %s, got %s", code, retrievedVerification.Code)
	}

	// Test CheckCode
	match, err := handler.CheckCode(username, code, scope)
	if err != nil || !match {
		t.Errorf("CheckCode failed: expected code to match and not expire")
	}

	// Test DeleteCode
	deleted := handler.DeleteCode(username, scope)
	if !deleted {
		t.Errorf("DeleteCode failed: expected code to be deleted")
	}
}

func TestVerificationCodeHandler_RegenerateCode(t *testing.T) {
	options := &Config{
		ExpiredAfterSec: 5 * time.Minute,
	}

	repository := NewMockCodeRepository()
	generator := &MockCodeGenerator{defCode: "", length: 5}

	handler, err := NewVerificationCodeHandler(generator, repository, options)
	if err != nil {
		t.Fatalf("Failed to create VerificationCodeHandler: %v", err)
	}

	username := "testuser"
	scope := "testscope"

	// Generate an initial code
	initialVerification, err := handler.GenerateCode(username, scope)
	if err != nil {
		t.Fatalf("GenerateCode error: %v", err)
	}

	// Regenerate the code without resetting expiration time
	regeneratedVerification, err := handler.RegenerateCode(username, scope, false)
	if err != nil {
		t.Fatalf("RegenerateCode error: %v", err)
	}

	// Check if the code was regenerated
	if initialVerification.Code == regeneratedVerification.Code {
		t.Error("Expected regenerated code to be different from the initial code")
	}

	if !(regeneratedVerification.ExpiredAt.Equal(initialVerification.ExpiredAt) || regeneratedVerification.ExpiredAt.Before(initialVerification.ExpiredAt)) {
		t.Errorf("Expected expiration time to be reset, got %d", regeneratedVerification.ExpireAfter)
	}
	time.Sleep(time.Second)
	// Regenerate the code with resetting expiration time
	regeneratedVerification, err = handler.RegenerateCode(username, scope, true)
	if err != nil {
		t.Fatalf("RegenerateCode error: %v", err)
	}

	// Check if the code was regenerated and the expiration time reset
	if initialVerification.Code == regeneratedVerification.Code {
		t.Error("Expected regenerated code to be different from the initial code")
	}

	if !regeneratedVerification.ExpiredAt.After(initialVerification.ExpiredAt) {
		t.Errorf("Expected expiration time to be reset, got %d", regeneratedVerification.ExpireAfter)
	}
}
