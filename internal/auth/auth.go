package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/patyukin/go-chat/internal/config"
	"github.com/patyukin/go-chat/internal/model"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

type Client struct {
	client  *http.Client
	address string
	token   string
}

func (c *Client) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func New(cfg *config.Config) *Client {
	return &Client{
		client:  &http.Client{},
		address: cfg.AuthServiceAddress,
		token:   cfg.ChatServiceToken,
	}
}

func (c *Client) SignUp(ctx context.Context, in model.SignUpRequest) (model.SignUpResponse, error) {
	data, err := json.Marshal(in)
	var result model.SignUpResponse
	if err != nil {
		return result, fmt.Errorf("failed to marshal in auth.SignUp: %w", err)
	}

	log.Debug().Msgf("sign up request: %s", string(data))
	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/v3/sign-up", c.address),
		bytes.NewReader(data),
	)
	if err != nil {
		return result, fmt.Errorf("failed to create request in auth.SignUp: %w", err)
	}

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", fmt.Sprintf("chat:%s", c.token))

	resp, err := c.client.Do(r)
	if err != nil {
		return result, fmt.Errorf("failed to do request in auth.SignUp: %w", err)
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Error().Msgf("failed to close response body in auth.SignUp, err %v", err)
		}
	}()

	log.Info().Msgf("status code: %d", resp.StatusCode)

	if resp.StatusCode >= http.StatusBadRequest {
		return result, fmt.Errorf("failed in auth.SignUp: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("error reading response body: %w", err)
	}

	log.Info().Msgf("body: %s", string(body))

	if err = json.Unmarshal(body, &result); err != nil {
		return result, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	log.Info().Msgf("user uuid: %v", result)

	return result, nil
}

func (c *Client) SignIn(ctx context.Context, in model.SignInRequest) (model.SignInResponse, error) {
	log.Debug().Msgf("sign in request: %s", in)
	data, err := json.Marshal(in)
	var tokens model.SignInResponse
	if err != nil {
		return tokens, fmt.Errorf("failed to marshal in auth.SignIn: %w", err)
	}

	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/v2/sign-in", c.address),
		bytes.NewReader(data),
	)
	if err != nil {
		return tokens, fmt.Errorf("failed to create request in auth.SignIn: %w", err)
	}

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", fmt.Sprintf("chat:%s", c.token))

	resp, err := c.client.Do(r)
	if err != nil {
		return tokens, fmt.Errorf("failed to do request in auth.SignIn: %w", err)
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Error().Msgf("failed to close response body in auth.SignIn, err %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tokens, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return tokens, fmt.Errorf("failed in auth.SignIn: %w", err)
	}

	if err = json.Unmarshal(body, &tokens); err != nil {
		return tokens, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	return tokens, nil
}

func (c *Client) ValidateToken(ctx context.Context, token string) (model.ValidateTokenResponse, error) {
	log.Debug().Msgf("ValidateToken request: %s", token)

	validateTokenData := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	data, err := json.Marshal(validateTokenData)
	if err != nil {
		return model.ValidateTokenResponse{}, fmt.Errorf("failed to marshal in auth.ValidateToken: %w", err)
	}

	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/v1/validate-token", c.address),
		bytes.NewReader(data),
	)
	if err != nil {
		return model.ValidateTokenResponse{}, fmt.Errorf("failed to create request in auth.ValidateToken: %w", err)
	}

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", fmt.Sprintf("chat:%s", c.token))

	resp, err := c.client.Do(r)
	if err != nil {
		return model.ValidateTokenResponse{}, fmt.Errorf("failed to do request in auth.ValidateToken: %w", err)
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Error().Msgf("failed to close response body in auth.ValidateToken, err %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.ValidateTokenResponse{}, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return model.ValidateTokenResponse{}, fmt.Errorf("failed in auth.ValidateToken: %w", err)
	}

	var userUUID model.ValidateTokenResponse
	if err = json.Unmarshal(body, &userUUID); err != nil {
		return model.ValidateTokenResponse{}, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	return userUUID, nil
}
