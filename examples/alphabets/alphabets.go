package main

import (
	"context"
	"fmt"
	go_verification "github.com/milito-78/go-verification"
	"log"
	"time"
)

func main() {
	ctx := context.Background()
	verification, _ := go_verification.NewVerificationCodeHandler(
		go_verification.NewAlphabetGenerator(6, false, false), //Code Generator
		go_verification.NewRedisCodeRepository(ctx, go_verification.RedisConfig{
			Prefix: "verification",
			Addr:   "localhost:6379",
			DB:     0,
		}), // Code Repository
		&go_verification.Config{
			ExpiredAfterSec: 180 * time.Second, // ExpiredAfterSec min is 2minutes, ExpiredAfterSec max is 10minutes
		}, // Options
	)

	//Generate code for `user_test` and for `forget-password` scope.
	//This code doesn't work in other scopes.
	//You can set a scope for global and use it anywhere
	code, err := verification.GenerateCode("user_test", "2step-verification")
	if err != nil {
		log.Printf("Error during generate code : %s", err)
	}
	fmt.Printf("Code is %s for scope %s and will be expired after %d", code.Code, code.Scope, code.ExpireAfter)

}
