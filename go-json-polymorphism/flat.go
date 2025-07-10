package main

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

type CredentialsType string

const (
	CredentialsTypeUsernamePassword CredentialsType = "username-password"
	CredentialsTypeToken            CredentialsType = "token"
)

type Credentials2 struct {
	Type CredentialsType `json:"type" validate:"required"`

	Username string `json:"username,omitempty" validate:"required_if=Type username-password"`
	Password string `json:"password,omitempty" validate:"required_if=Type username-password"`

	Token string `json:"token,omitempty" validate:"required_if=Type token"`
}

type Config2 struct {
	Credentials Credentials2 `json:"credentials" validate:"required"`
}

func (c *Config2) UnmarshalJSON(data []byte) error {
	type Alias Config2
	if err := json.Unmarshal(data, (*Alias)(c)); err != nil {
		return err
	}

	return c.Validate()
}

func (c *Config2) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}
