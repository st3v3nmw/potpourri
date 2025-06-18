package main

import (
	"log"

	"github.com/canonical/starlark/starlark"
)

func CheckPasswordPolicyStarlark(password, userRole string) ValidationResult {
	// Declare context to pass to Starlark
	user := &starlark.Dict{}
	user.SetKey(starlark.String("role"), starlark.String(userRole))

	context := starlark.StringDict{
		"password": starlark.String(password),
		"user":     user,
	}

	// Execute Starlark program
	thread := &starlark.Thread{Name: "Check Password Policy"}
	globals, err := starlark.ExecFile(thread, "password_policy.star", nil, context)
	if err != nil {
		log.Fatalf("Error executing program: %v\n", err)
	}

	// Check the policy
	checkPolicy := globals["check_policy"]
	v, err := starlark.Call(thread, checkPolicy, nil, nil)
	if err != nil {
		log.Fatalf("Error checking policy: %v\n", err)
	}

	// Extract returned values
	dict, ok := v.(*starlark.Dict)
	if !ok {
		log.Fatal("Expected dict return value")
	}

	compliantVal, found, err := dict.Get(starlark.String("compliant"))
	if !found || err != nil {
		log.Fatal("Error extracting .compliant value")
	}
	compliant := bool(compliantVal.(starlark.Bool))

	violationsVal, found, err := dict.Get(starlark.String("violations"))
	if !found || err != nil {
		log.Fatal("Error extracting .violations value")
	}

	violationsList := violationsVal.(*starlark.List)
	violations := make([]string, violationsList.Len())
	for i := 0; i < violationsList.Len(); i++ {
		violations[i] = string(violationsList.Index(i).(starlark.String))
	}

	return ValidationResult{
		Compliant:  compliant,
		Violations: violations,
	}
}
