package go_verification

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = Duration(time.Duration(value))
		return nil
	case string:
		tmp, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(tmp)
		return nil
	default:
		return errors.New("invalid duration")
	}
}

type Config struct {
	ExpiredAfterSec time.Duration
}

type VerificationCode struct {
	ExpireAfter int
	ExpiredTime Duration
	ExpiredAt   time.Time
	Username    string
	Scope       string
	Code        string
}

type VerificationCodeHandler struct {
	repository CodeRepositoryInterface
	generator  CodeGenerator
	config     *Config
}

func NewVerificationCodeHandler(generator CodeGenerator, repository CodeRepositoryInterface, options *Config) (*VerificationCodeHandler, error) {
	checkConfig(options)

	return &VerificationCodeHandler{
		repository: repository,
		generator:  generator,
		config:     options,
	}, nil
}

func (v *VerificationCodeHandler) GenerateCode(username, scope string) (*VerificationCode, error) {
	verify, err := v.repository.GetCode(username, scope)
	if err == nil {
		return verify, nil
	}

	code := v.generator.Generate()
	verify, err = v.repository.SaveCode(username, code, scope, v.config.ExpiredAfterSec)
	if err != nil {
		return nil, err
	}
	return verify, nil
}

func (v *VerificationCodeHandler) GetCode(username, scope string) (*VerificationCode, error) {
	verify, err := v.repository.GetCode(username, scope)
	fmt.Println(verify)
	if err != nil {
		return nil, err
	}
	return verify, nil
}

func (v *VerificationCodeHandler) CheckCode(username, code, scope string) (bool, error) {
	verify, err := v.repository.GetCode(username, scope)
	if err != nil {
		return false, err
	}
	if verify.Code == code && verify.ExpiredAt.After(time.Now()) {
		return true, nil
	}
	return false, errors.New("code expired")
}

func (v *VerificationCodeHandler) DeleteCode(username, scope string) bool {
	return v.repository.DeleteCode(username, scope)
}

func (v *VerificationCodeHandler) RegenerateCode(username, scope string, resetExpireTime bool) (*VerificationCode, error) {
	verify, err := v.repository.GetCode(username, scope)
	if err != nil {
		return nil, err
	}

	if resetExpireTime {
		v.DeleteCode(username, scope)
		saveCode, err := v.GenerateCode(username, scope)
		if err != nil {
			return nil, err
		}
		return saveCode, nil
	}

	timeExpired := verify.ExpireAfter
	if verify.ExpiredAt.After(time.Now()) {
		timeExpired = int(verify.ExpiredAt.Sub(time.Now()).Seconds())
	}
	code := v.generator.Generate()
	v.repository.DeleteCode(username, scope)
	saveCode, err := v.repository.SaveCode(username, code, scope, time.Duration(timeExpired)*time.Second)
	if err != nil {
		return nil, err
	}
	return saveCode, nil
}

func checkConfig(config *Config) {
	if config.ExpiredAfterSec.Seconds() < time.Minute.Seconds() {
		config.ExpiredAfterSec = 2 * time.Minute
	} else if config.ExpiredAfterSec.Seconds() > 10*time.Minute.Seconds() {
		config.ExpiredAfterSec = 10 * time.Minute
	}
}
