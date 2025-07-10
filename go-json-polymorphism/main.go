package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"
)

func main() {
	log.SetFlags(0)

	var useSum bool
	flag.BoolVar(&useSum, "sum", false, "Use sum type approach")

	var useFlat bool
	flag.BoolVar(&useFlat, "flat", false, "Use flat struct approach")

	flag.Parse()

	if !useSum && !useFlat {
		useSum = true
	}

	if useSum {
		fmt.Println("=== Using Sum Type Approach ===")
		sumJsonConfigs := []string{
			`{
				"credentials-type": "username-password",
				"credentials-body": {
					"username": "john",
					"password": "secret123"
				}
			}`,
			`{
				"credentials-type": "token",
				"credentials-body": {
					"token": "abc123xyz"
				}
			}`,
		}
		runConfig1(sumJsonConfigs)
	} else {
		fmt.Println("=== Using Flat Struct Approach ===")
		flatJsonConfigs := []string{
			`{
				"credentials": {
					"type": "username-password",
					"username": "alice",
					"password": "mypassword"
				}
			}`,
			`{
				"credentials": {
					"type": "token",
					"token": "xyz789abc"
				}
			}`,
		}
		runConfig2(flatJsonConfigs)
	}
}

func runConfig1(jsonConfigs []string) {
	for _, jsonConfig := range jsonConfigs {
		var config Config1
		err := json.Unmarshal([]byte(jsonConfig), &config)
		if err != nil {
			log.Fatalf("Error unmarshalling config: %v\n", err)
		}

		fmt.Println(strings.ToUpper(config.CredentialsType))
		fmt.Printf("Unmarshalled: %+v\n", config)

		out, err := json.Marshal(config)
		if err != nil {
			log.Fatalf("Error marshalling config: %v\n", err)
		}

		fmt.Printf("Marshalled: %s\n\n", string(out))
	}
}

func runConfig2(jsonConfigs []string) {
	for _, jsonConfig := range jsonConfigs {
		var config Config2
		err := json.Unmarshal([]byte(jsonConfig), &config)
		if err != nil {
			log.Fatalf("Error unmarshalling config: %v\n", err)
		}

		fmt.Println(strings.ToUpper(string(config.Credentials.Type)))
		fmt.Printf("Unmarshalled: %+v\n", config)

		out, err := json.Marshal(config)
		if err != nil {
			log.Fatalf("Error marshalling config: %v\n", err)
		}

		fmt.Printf("Marshalled: %s\n\n", string(out))
	}
}
