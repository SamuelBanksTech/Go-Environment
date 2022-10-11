// Copyright 2022 SamuelBanksTech. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package env

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var keys = make(map[string]string)

// Environment struct
type Environment struct {
	// EnvPath is the path to your environment variable file, usually .env in the root but can be whatever you like
	EnvPath string
	// EnableOsEnvOverride allows any matching os level environment variables to override any contained in the environment file
	EnableOsEnvOverride bool
	// HideOutput simply shows or hides console output of parsed environment variables
	HideOutput bool
}

// basePath takes the path from the currently running executable,
// this way ensures that env files will be picked up if using go run or a compiled binary
func (e *Environment) basePath() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}

// LoadEnv loads the specified environment file and splits them into key value pairs
// these are then stored in a globally accessible map via a getter method
// if Environment.EnableOsOverride is set to true any matching keys that exist in the actual
// environment will override the value in the file
func (e *Environment) LoadEnv() error {

	if e.EnvPath == "" {
		e.EnvPath = ".env"
	}
	_, err := os.Stat(e.EnvPath)
	if err != nil {
		_, err = os.Stat(e.basePath() + "/" + e.EnvPath)
		if err == nil {
			e.EnvPath = e.basePath() + "/" + e.EnvPath
		} else {
			return err
		}
	}

	file, err := os.Open(e.EnvPath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if line != "" && string(line[0]) != "#" {
			envVar := strings.SplitN(line, "=", 2)

			if e.EnableOsEnvOverride {
				if os.Getenv(strings.TrimSpace(envVar[0])) != "" {
					keys[strings.TrimSpace(envVar[0])] = os.Getenv(strings.TrimSpace(envVar[0]))
				} else {
					keys[strings.TrimSpace(envVar[0])] = strings.TrimSpace(envVar[1])
				}
			} else {
				keys[strings.TrimSpace(envVar[0])] = strings.TrimSpace(envVar[1])
			}
		}
	}

	err = scanner.Err()
	if err != nil {
		return err
	}

	if !e.HideOutput {
		fmt.Println("ENV LOADED")
		fmt.Println("----------")
		for key, value := range keys {
			fmt.Println(key + "  :  " + value)
		}
		fmt.Println("----------")
	}

	return nil
}

// Get simply returns the matching key value or falling back to os env
func Get(key string) string {
	if keyVar, ok := keys[key]; ok {
		return keyVar
	}

	return os.Getenv(key)
}
