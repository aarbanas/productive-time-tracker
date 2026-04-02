package appinit

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	TokenFileName          = ".productive_token"
	OrganizationIDFileName = ".productive_org_id"
)

// LoadCredentials prompts for or loads the auth token and organization ID from the user home directory.
func LoadCredentials() (token, orgID string, err error) {
	token, err = loadOrRequestCredential(TokenFileName, "X-Auth-Token")
	if err != nil {
		return "", "", err
	}
	orgID, err = loadOrRequestCredential(OrganizationIDFileName, "X-Organization-Id")
	if err != nil {
		return "", "", err
	}
	return token, orgID, nil
}

func loadOrRequestCredential(fileName string, promptMessage string) (string, error) {
	credPath := credentialPath(fileName)
	docsURL := "https://developer.productive.io/index.html#header-authorization"

	if credential, err := loadTokenFromFile(credPath); err == nil {
		fmt.Printf("✓ %s loaded from storage\n", promptMessage)
		return credential, nil
	}

	fmt.Printf("You can read here how to setup the %s: %s\n", promptMessage, docsURL)
	fmt.Printf("%s not found. Please enter %s: ", promptMessage, promptMessage)
	var credential string
	if _, err := fmt.Scanln(&credential); err != nil {
		return "", err
	}

	if credential == "" {
		return "", fmt.Errorf("%s cannot be empty", promptMessage)
	}

	if err := saveTokenToFile(credPath, credential); err != nil {
		fmt.Printf("⚠ Warning: Could not save %s: %v\n", promptMessage, err)
	}

	fmt.Printf("✓ %s saved\n", promptMessage)
	return credential, nil
}

func credentialPath(fileName string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fileName
	}
	return filepath.Join(homeDir, fileName)
}

func loadTokenFromFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func saveTokenToFile(path string, token string) error {
	return os.WriteFile(path, []byte(token), 0600)
}
