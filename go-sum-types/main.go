package main

import (
	"encoding/json"
	"fmt"
)

type Credentials interface {
	CredentialsType() string

	UnmarshalJSON(data []byte) error
}

var credentialTypes = map[string]func() Credentials{
	"username-password": func() Credentials { return &UsernamePasswordCredz{} },
	"token":             func() Credentials { return &TokenCredz{} },
}

type UsernamePasswordCredz struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c UsernamePasswordCredz) CredentialsType() string {
	return "username-password"
}

func (c *UsernamePasswordCredz) UnmarshalJSON(data []byte) error {
	type Alias UsernamePasswordCredz
	var alias Alias

	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	if alias.Username == "" || alias.Password == "" {
		return fmt.Errorf("username and password are required")
	}

	*c = UsernamePasswordCredz(alias)
	return nil
}

type TokenCredz struct {
	Token string `json:"token"`
}

func (c TokenCredz) CredentialsType() string {
	return "token"
}

func (c *TokenCredz) UnmarshalJSON(data []byte) error {
	type Alias TokenCredz
	var alias Alias

	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	if alias.Token == "" {
		return fmt.Errorf("token is required")
	}

	*c = TokenCredz(alias)
	return nil
}

type Config struct {
	CredentialsType string      `json:"credentials-type"`
	CredentialsBody Credentials `json:"-"`
}

func (c *Config) UnmarshalJSON(data []byte) error {
	type ConfigAlias Config
	var raw struct {
		*ConfigAlias

		RawCredzBody json.RawMessage `json:"credentials-body"`
	}

	raw.ConfigAlias = (*ConfigAlias)(c)
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	constructor, exists := credentialTypes[raw.CredentialsType]
	if !exists {
		return fmt.Errorf("unknown credentials type: %s", raw.CredentialsType)
	}

	raw.CredentialsBody = constructor()
	err := json.Unmarshal(raw.RawCredzBody, raw.CredentialsBody)
	if err != nil {
		return err
	}

	*c = Config(*raw.ConfigAlias)
	return nil
}

func main() {
	fmt.Println("USERNAME-PASSWORD")
	userPassJSON := `{
		"credentials-type": "username-password",
		"credentials-body": {
			"username": "john",
			"password": "secret123"
		}
	}`

	var config1 Config
	if err := json.Unmarshal([]byte(userPassJSON), &config1); err != nil {
		fmt.Printf("Error unmarshaling username-password: %v\n", err)
	} else {
		fmt.Printf("Type: %s\n", config1.CredentialsBody.CredentialsType())

		if upCreds, ok := config1.CredentialsBody.(*UsernamePasswordCredz); ok {
			fmt.Printf("Username: %s\n", upCreds.Username)
			fmt.Printf("Password: %s\n", upCreds.Password)
		}

	}

	fmt.Println("\nTOKEN")
	tokenJSON := `{
		"credentials-type": "token",
		"credentials-body": {
			"token": "abc123xyz"
		}
	}`

	var config2 Config
	if err := json.Unmarshal([]byte(tokenJSON), &config2); err != nil {
		fmt.Printf("Error unmarshaling token: %v\n", err)
	} else {
		fmt.Printf("Type: %s\n", config2.CredentialsBody.CredentialsType())

		if tokenCreds, ok := config2.CredentialsBody.(*TokenCredz); ok {
			fmt.Printf("Token: %s\n", tokenCreds.Token)
		}
	}
}
