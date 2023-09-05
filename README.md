# Go Verification Code Handler
This is a package for handling verification codes.

If you're tired of writing verification code scenarios for each project like me :joy:,
you can use this package and avoid the boring task of copying and pasting your codes!

This package will help you generate verification codes for users with specific scopes.

#### Wait!! what's the scope?!
A `scope` is a parameter that specifies the particular action for which this code is valid, such as `forget password` or any other required actions.<br/>
If you're confused about what I mean, you can check the [examples folder](examples) and the [examples](#example-section) in this README file.

## Installation
Ok! Stop talking and let's install the package.<br/>
Use this command to get package :
```cmd
go get github.com/milito-78/go-verification-code
```
#### *HINT*
This package uses Redis to store codes, but you have the flexibility to change this by implementing a new repository driver and passing it to the verification code handler struct.

<h2 id="#example-section"> Examples </h2>
You can check the examples folder. There are examples of how it works. But let me show you some examples below.


To start, first, you must create a new instance of `VerificationCodeHandler` struct:
```go
package main

import (
	"context"
	go_verification "github.com/milito-78/go-verification"
	"log"
	"time"
)

func main()  {
	ctx := context.Background()
	verification, _ := go_verification.NewVerificationCodeHandler(
		go_verification.NewRegexGenerator(`N-\d{5}`), //Regex Code Generator
		go_verification.NewRedisCodeRepository(ctx, go_verification.RedisConfig{
			Prefix: "verification",
			Addr:   "localhost:6379",
			DB:     0,
		}), // Redis Code Repository
		&go_verification.Config{
			ExpiredAfterSec: 180 * time.Second, // ExpiredAfterSec min is 2minutes, ExpiredAfterSec max is 10minutes
		}, // Options
	)
	
	code, err := verification.GenerateCode("user_test", "forget-password")
	if err != nil {
		log.Printf("Error during generate code : %s", err)
	}
	log.Printf("Code is %s for scope %s and will be expired after %d", code.Code, code.Scope, code.ExpireAfter)
}
```
`GenerateCode` will be used to generate a code.  If a code exists for the user and scope, it will return the existing code. To reset the existing code and generate a new one, you can use the `RegenerateCode` method.

In this example, we create a handler using `RedisCodeRepository` & `RegexGenerator`, and then generate a code for the user `user_test` within the `forget-password` scope.<br/>
If you want to check a code, you can use the `CheckCode` method.

```go
    //...
    valid, err := verification.CheckCode("user_test","12345","forget-password")
    if err != nil || !valid {
       log.Printf("code is invalid")
    }else{
        //Code is valid
    }
	
```

If you want to get code use `GetCode` method.

```go
    //...
    code, err := verification.GetCode("user_test","forget-password")
    if err != nil{
        log.Printf("code not exists")
    }else{
        //Use current code for user & scope
    }
	
```

If you need to regenerate a code for a user without resetting the expiration time or resetting an existing code, you can use the `RegenerateCode` method.
```go
    //...
    // Use `false` preserve the expiration time and true to reset expiration time
    code, err := verification.RegenerateCode("user_test","forget-password",false)
    if err != nil{
        log.Printf("error during regenerate code")
    }else{
        // New code for user & scope
    }
	
```

There are 4 types for generating codes :

| Types               | Struct            | Options                                                                                                                                                                                                                                                | Output |
|---------------------|-------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------|
| Numbers             | NumberGenerator   | `length:` An integer argument that specifies the length of the generated code.<br/> `notZeroAtStart:` A bool arg that handle generated code starts with 0 or not.                                                                                      | 09283  |
| Alphabets           | AlphabetGenerator | `length:` An integer argument that specifies the length of the generated code.<br/> `allCapital:` A bool arg that handle all of generated code is capital letter.<br/> `allNonCapital:` A bool arg that handle all of generated code is small letter.  | seAsaz |
| Alphabets & Numbers | WordGenerator     | `length:` An integer argument that specifies the length of the generated code.                                                                                                                                                                         | s2W09v |
| Regex               | RegexGenerator    | `regex:` A string argument that specifies a regex pattern for generating the code.                                                                                                                                                                     | de2ds4 |


## License

The Milito Go Verification package is an open-sourced package licensed under the [MIT license](https://opensource.org/licenses/MIT).
