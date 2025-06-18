package main

import (
	"encoding/json"
	"fmt"
)

type ValidationResult struct {
	Compliant  bool     `json:"compliant"`
	Violations []string `json:"violations"`
}

func (r ValidationResult) String() string {
	out, _ := json.MarshalIndent(r, "", "  ")
	return string(out)
}

func main() {
	password := "pass123"
	userRole := "admin"

	result := CheckPasswordPolicyStarlark(password, userRole)
	fmt.Println("Starlark:", result)

	fmt.Println()

	result = CheckPasswordPolicyConfdb(password, userRole)
	fmt.Println("Confdb:", result)
}
