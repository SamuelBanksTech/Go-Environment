## Environment variable loader

Why use this over `os.Getenv("FOO")`? Read below to find out!

This package will load an environment variable file (usually a .env file placed in the root) and have these variables globally available using `env.Get("FOO")`.
You can optionally allow 'real' environment variables to override the variables in the .env file allowing for some nice defaults to be set in the case of a missing environment variable.

Whats more, is unlike other env loaders, the variables contained in your .env file do not actually get set in the real environment and are only available to the running Go application making them more secure to malicious code that may be running on the same server. 

Using `env.Get("XXX")` can completly replace `os.Getenv("XXX")` if there is no variable in the file then next check will be the 'real' environment.

#### Usage

`go get github.com/SamuelBanksTech/Go-Environment`


#### Basic Example
```go
package main

import (
	"context"
	"fmt"
	"os"
	"github.com/SamuelBanksTech/Go-Environment"
)

func main() {

	//Init env vars
	envinit := env.Environment{
		EnvPath:             ".env",
		EnableOsEnvOverride: true,
		HideOutput:          false,
	}

	err := envinit.LoadEnv()
	if err != nil {
		log.Fatalln(err)
	}
	
	fmt.Println(env.Get("FOO"))
	
}
```
#### Example .env

```dotenv
FOO=bar
BAR=donk
#IGNOREDVAR=foobar
```