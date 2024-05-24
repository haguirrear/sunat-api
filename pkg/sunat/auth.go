package sunat

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AuthParams struct {
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
}

type AuthResponseBody struct {
	AccessToken string `json:"access_token"`
}

const defaultTimeout = time.Second

// body = {
//     "scope": "https://api-cpe.sunat.gob.pe",
//     "grant_type": "password",
//     "client_id": self.client_id,
//     "client_secret": self.client_secret,
//     "username": self.username,
//     "password": self.password,
// }

func GetToken(baseURL string, params AuthParams) (token string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	authURL := fmt.Sprintf("%s/v1/clientessol/%s/oauth2/token/", baseURL, params.ClientID)
	form := url.Values{}
	form.Set("scope", "https://api-cpe.sunat.gob.pe")
	form.Set("grant_type", "password")
	form.Set("client_id", params.ClientID)
	form.Set("client_secret", params.ClientSecret)
	form.Set("username", params.Username)
	form.Set("password", params.Password)

	encoded := strings.NewReader(form.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, authURL, encoded)
	if err != nil {
		return "", fmt.Errorf("error building auth request: %w", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error in auth response: %w", err)
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		return "", fmt.Errorf("error reading body of auth request with response %s: %w", res.Status, err)
	}

	if res.StatusCode >= 400 {
		return "", fmt.Errorf("error authorizing with SUNAT: %s | %s", res.Status, string(body))
	}

	var parsed AuthResponseBody
	if err = json.Unmarshal(body, &parsed); err != nil {
		return "", fmt.Errorf("error deserializing auth response body into json: %w", err)
	}

	return parsed.AccessToken, nil
}
