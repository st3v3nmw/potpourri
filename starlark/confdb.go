package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type PasswordPolicy struct {
	DefaultMinimumLength int                   `json:"defaultMinimumLength"`
	RequireSpecialChars  bool                  `json:"requireSpecialChars"`
	RoleOverrides        map[string]RoleConfig `json:"roleOverrides"`
}

type RoleConfig struct {
	MinimumLength int `json:"minimumLength"`
}

func CheckPasswordPolicyConfdb(password, userRole string) ValidationResult {
	// Load policy from JSON
	data, err := os.ReadFile("password_policy.json")
	if err != nil {
		log.Fatalf("Error reading policy file")
	}

	var config struct {
		PasswordPolicy PasswordPolicy `json:"passwordPolicy"`
	}
	json.Unmarshal(data, &config)

	policy := config.PasswordPolicy
	violations := []string{}

	// Determine minimum length
	minLength := policy.DefaultMinimumLength
	if override, exists := policy.RoleOverrides[userRole]; exists {
		minLength = override.MinimumLength
	}

	// Check length
	if len(password) < minLength {
		violations = append(violations, fmt.Sprintf("Password is too short (minimum %d)", minLength))
	}

	// Check special characters
	if policy.RequireSpecialChars {
		hasSpecial := false
		for _, char := range password {
			if strings.ContainsRune("!@#$%^&*", char) {
				hasSpecial = true
				break
			}
		}

		if !hasSpecial {
			violations = append(violations, "Password must contain special characters")
		}
	}

	return ValidationResult{
		Compliant:  len(violations) == 0,
		Violations: violations,
	}
}
