package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	TenantID     string `toml:"tenant-id"`
	ClientId     string `toml:"client-id"`
	ClientSecret string `toml:"client-secret"`
	Scope        string `toml:"scope"`
	GroupName    string `toml:"group-name"`
	Domain       string `toml:"domain"`
}

type MicrosoftGraphResponse struct {
	Context string `json:"@odata.context"`
	Groups  []struct {
		Name string `json:"displayName"`
		Type string `json:"@odata.type"`
	} `json:"value"`
}

// load config file
func LoadConfig(configPath string) (*Config, error) {
	if _, err := os.Stat(configPath); err != nil {
		return nil, fmt.Errorf("Unable locate config file. Error: %s", err)
	}
	var config Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		// log.Fatal(log_prefix, "Unable to load config file. Error: ", err)
		return nil, fmt.Errorf("Unable to load config file. Error: %s", err)
	}
	// Set default values where appropriate
	if config.TenantID == "" {
		config.TenantID = "common"
	}
	return &config, nil
}

func ValidateCredentials(configPath string, username string, password string) int {
	var log_prefix = fmt.Sprintf("[%s] ", username)

	// Load config file
	config, err := LoadConfig(configPath)
	if err != nil {
		log.Println(log_prefix, strings.ReplaceAll(err.Error(), "\n", ". "))
		return 4 // PAM_SYSTEM_ERR
	}

	// Generate the OAuth2 config
	oauth2Config := oauth2.Config{
		ClientID:     config.ClientId,
		ClientSecret: config.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.microsoftonline.com/" + config.TenantID + "/oauth2/v2.0/authorize",
			TokenURL: "https://login.microsoftonline.com/" + config.TenantID + "/oauth2/v2.0/token",
		},
		RedirectURL: "urn:ietf:wg:oauth:2.0:oob", // this is the "no redirect" URL
		Scopes:      []string{config.Scope},
	}

	// If there is no suffix then use the default domain
	if !strings.Contains(username, "@") {
		username = username + "@" + config.Domain
	}

	// Retrieve an OAuth token from AzureAD
	// The "password" grant type should only be used "when there is a high degree
	// of trust between the resource owner and the client (e.g., the client is
	// part of the device operating system or a highly privileged application),
	// and when other authorization grant types are not available."
	// See https://tools.ietf.org/html/rfc6749#section-4.3 for more info.
	oauthToken, err := oauth2Config.PasswordCredentialsToken(
		context.Background(),
		username,
		password,
	)
	// Note that we do not perform further validity checks as we are not using
	// this token directly but instead using it to make a further request against
	// the Microsoft Graph API that will fail if the token is invalid.
	if err != nil {
		log.Println(log_prefix, strings.ReplaceAll(err.Error(), "\n", ". "))
		return 7 // PAM_AUTH_ERR
	}

	// Use the access token to retrieve group memberships for the user in question
	// We compare these against the specified group name to determine whether
	// authentication is successful.
	aadGroupNames, err := RetrieveAADGroupMemberships(oauthToken.AccessToken)
	if err != nil {
		log.Println(log_prefix, "AzureAD groups could not be loaded for this user")
		return 8 // PAM_CRED_INSUFFICIENT
	}
	for _, aadGroupName := range aadGroupNames {
		if aadGroupName == config.GroupName {
			log.Println(log_prefix, "Authentication succeeded")
			return 0 // PAM_SUCCESS
		}
	}

	log.Println(log_prefix, "Authentication was successful but authorization failed")
	return 7 // PAM_PERM_DENIED
}

// RetrieveAADGroupMemberships returns a []string containing the names
// of Azure AD groups that this user belongs to, using the provided
// bearer token.
func RetrieveAADGroupMemberships(bearerToken string) ([]string, error) {
	groupNames := []string{}

	// AzureAD access tokens are *NOT* verifiable JWTs and can only be validated by Microsoft Graph
	// See https://stackoverflow.com/questions/60778634/failing-signature-validation-of-jwt-tokens-from-azure-ad
	// and https://github.com/AzureAD/azure-activedirectory-identitymodel-extensions-for-dotnet/issues/609#issuecomment-529537264
	parsedToken, err := jwt.ParseInsecure([]byte(bearerToken))
	if err != nil {
		return groupNames, err
	}

	// Instead of verifying the token via its signature, we verify it by its capabilities
	// Namely, we extract the userId and use this, to to make a call to the Microsoft Graph API
	userId := parsedToken.PrivateClaims()["oid"].(string)

	// Create a new request using http with correct authorization header
	req, err := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/users/"+userId+"/memberOf", nil)
	req.Header.Add("Authorization", "Bearer "+bearerToken)

	// Use http Client to send the request, closing when finished
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return groupNames, err
	}
	defer resp.Body.Close()

	// Read response and unmarshal JSON into a struct
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return groupNames, err
	}
	var response MicrosoftGraphResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return groupNames, err
	}

	// Look through the struct for Microsoft Graph groups
	for _, group := range response.Groups {
		if group.Type == "#microsoft.graph.group" {
			groupNames = append(groupNames, group.Name)
		}
	}
	return groupNames, nil
}
