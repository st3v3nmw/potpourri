package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Credentials1 interface {
	CredentialsType() string

	String() string
}

var credentialTypesRegistry = map[string]func() Credentials1{
	"username-password": func() Credentials1 { return &UsernamePasswordCredz{} },
	"token":             func() Credentials1 { return &TokenCredz{} },
}

type UsernamePasswordCredz struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (c UsernamePasswordCredz) CredentialsType() string {
	return "username-password"
}

func (c UsernamePasswordCredz) String() string {
	return fmt.Sprintf("Username:%s Password:%s", c.Username, c.Password)
}

type TokenCredz struct {
	Token string `json:"token" validate:"required"`
}

func (c TokenCredz) CredentialsType() string {
	return "token"
}

func (c TokenCredz) String() string {
	return fmt.Sprintf("Token:%s", c.Token)
}

type Config1 struct {
	CredentialsType string       `json:"credentials-type" validate:"required"`
	CredentialsBody Credentials1 `json:"credentials-body" validate:"required"`
}

func (c *Config1) UnmarshalJSON(data []byte) error {
	var raw struct {
		CredentialsType string          `json:"credentials-type"`
		CredentialsBody json.RawMessage `json:"credentials-body"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	constructor, exists := credentialTypesRegistry[raw.CredentialsType]
	if !exists {
		return fmt.Errorf("unknown credentials type: %s", raw.CredentialsType)
	}

	c.CredentialsType = raw.CredentialsType

	c.CredentialsBody = constructor()
	err := json.Unmarshal(raw.CredentialsBody, c.CredentialsBody)
	if err != nil {
		return err
	}

	return c.Validate()
}

func (c *Config1) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}
