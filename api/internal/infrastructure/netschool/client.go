package netschool

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"netschool-proxy/api/api/internal/pkg/encoding"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
	retryMax   int
	retryWait  time.Duration
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: timeout},
		timeout:    timeout,
		retryMax:   3, 
		retryWait:  time.Second, 
	}
}

func (c *Client) Login(ctx context.Context, username, password string, schoolID int, loginData map[string]interface{}) (string, error) {
	
	_, err := c.httpClient.Get(fmt.Sprintf("%s/webapi/logindata", c.baseURL))
	if err != nil {
		return "", fmt.Errorf("failed to get login data cookies: %w", err)
	}

	
	salt := loginData["salt"].(string)
	pw, pw2, err := hashPassword(password, salt)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	
	formData := url.Values{}
	formData.Add("loginType", "1") 
	formData.Add("scid", fmt.Sprintf("%d", schoolID))
	formData.Add("un", username)
	formData.Add("pw", pw)
	formData.Add("pw2", pw2)

	
	for key, value := range loginData {
		if key != "salt" {
			formData.Add(key, fmt.Sprintf("%v", value))
		}
	}

	
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/webapi/login", c.baseURL),
		bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "NetCityApp/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse login response: %w", err)
	}

	if accessToken, ok := result["at"].(string); ok && accessToken != "" {
		return accessToken, nil
	}

	return "", fmt.Errorf("login failed: %s", result["message"])
}

func hashPassword(password, salt string) (string, string, error) {
	
	win1251Pass, err := encoding.EncodeWindows1251(password)
	if err != nil {
		return "", "", err
	}

	
	md5Hash := md5.Sum([]byte(win1251Pass))
	pw2 := hex.EncodeToString(md5Hash[:])

	
	saltedInput := salt + pw2
	finalHash := md5.Sum([]byte(saltedInput))
	pw := hex.EncodeToString(finalHash[:])[:len(password)]

	return pw, pw2, nil
}


func (c *Client) GetLoginData(ctx context.Context) (map[string]interface{}, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/webapi/logindata", c.baseURL), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get logindata failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse logindata response: %w", err)
	}

	return result, nil
}


func (c *Client) WithRetrySettings(maxRetries int, waitTime time.Duration) *Client {
	c.retryMax = maxRetries
	c.retryWait = waitTime
	return c
}


func (c *Client) DoRequestWithRetry(ctx context.Context, req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i <= c.retryMax; i++ {
		resp, err = c.httpClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			return resp, nil
		}

		if resp != nil {
			resp.Body.Close()
		}

		
		if i == c.retryMax {
			break
		}

		
		select {
		case <-time.After(c.retryWait):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return nil, fmt.Errorf("request failed after %d retries: %w", c.retryMax, err)
}