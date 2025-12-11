package auth

import "os"

func GetClientId() string {
	clientId, ok := os.LookupEnv("GITHUB_CLIENT_ID")
	if !ok {
		panic("'GITHUB_CLIENT_ID' is not  set")
	}
	return clientId
}

func GetClientSecret() string {
	secret, ok := os.LookupEnv("GITHUB_CLIENT_SECRET")
	if !ok {
		panic("'GITHUB_CLIENT_SECRET' is not  set")
	}
	return secret
}

func GetGithubAppId() string {
	appId, ok := os.LookupEnv("GITHUB_APP_ID")
	if !ok {
		panic("'GITHUB_APP_ID' is not set")
	}
	return appId
}

func GetGithubOauthRedirectURI() string {
	redirectUri, ok := os.LookupEnv("GITHUB_OAUTH_REDIRECT_URI")
	if !ok {
		panic("'GITHUB_OAUTH_REDIRECT_URI' is not set")
	}
	return redirectUri
}
