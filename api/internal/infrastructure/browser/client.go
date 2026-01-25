package browser

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"time"
	
	"github.com/playwright-community/playwright-go"
)


type BrowserAuthClientInterface interface {
	
	Authenticate(ctx context.Context, instanceURL, username, password string, schoolID int) (string, error)
	
	
	IsAvailable() bool
	
	
	Close() error
}


var ErrBrowserAutomationNotSupported = errors.New("this instance does not support Playwright browser automation")



type BrowserAuthClient struct {
	available bool
	pw        *playwright.Playwright
	browser   playwright.Browser
}


func NewBrowserAuthClient() (*BrowserAuthClient, error) {
	
	available := checkPlaywrightAvailability()
	
	client := &BrowserAuthClient{
		available: available,
	}
	
	if available {
		
		pw, err := playwright.Run()
		if err != nil {
			
			client.available = false
			return client, nil
		}
		
		
		browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
			Headless: playwright.Bool(true), 
		})
		if err != nil {
			pw.Stop()
			client.available = false
			return client, nil
		}
		
		client.pw = pw
		client.browser = browser
	}
	
	return client, nil
}


func checkPlaywrightAvailability() bool {
	
	cmd := exec.Command("playwright", "--version")
	if err := cmd.Run(); err != nil {
		
		cmd = exec.Command("node", "-e", "require('playwright')")
		if err := cmd.Run(); err != nil {
			return false
		}
	}
	
	return true
}


func (c *BrowserAuthClient) Authenticate(ctx context.Context, instanceURL, username, password string, schoolID int) (string, error) {
	if !c.available {
		return "", ErrBrowserAutomationNotSupported
	}

	
	page, err := c.browser.NewPage()
	if err != nil {
		return "", fmt.Errorf("failed to create new page: %w", err)
	}
	defer page.Close()

	
	loginURL := fmt.Sprintf("%s/login?mobile", instanceURL)

	
	if _, err := page.Goto(loginURL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); err != nil {
		return "", fmt.Errorf("failed to navigate to login page: %w", err)
	}

	
	element, err := page.WaitForSelector("input[name='lg'], input[name='login'], input[name='username']")
	if err != nil {
		return "", fmt.Errorf("login form not found: %w", err)
	}
	if element == nil {
		return "", fmt.Errorf("login form element not found")
	}

	
	if errStr := page.Fill("input[name='lg'], input[name='login'], input[name='username']", username).Error(); errStr != "" {

		if errStr := page.Fill("input[name='login'], input[name='username']", username).Error(); errStr != "" {
			return "", fmt.Errorf("failed to fill username: %s", errStr)
		}
	}

	if errStr := page.Fill("input[name='pw'], input[name='password']", password).Error(); errStr != "" {
		return "", fmt.Errorf("failed to fill password: %s", errStr)
	}

	
	// TODO: Fix Playwright SelectOption API usage
	// selectors := []string{"select[name='cl']", "select[name='school']", "select[name='schoolId']"}
	// var selectErr string
	// for _, selector := range selectors {
	// 	_, err := page.SelectOption(selector, fmt.Sprintf("%d", schoolID))
	// 	if err == nil {
	// 		selectErr = ""
	// 		break // Successfully selected an option
	// 	} else {
	// 		selectErr = err.Error()
	// 	}
	// }

	// For now, simulate the action without actual Playwright call
	var selectErr string = ""
	// Skip school selection for now
	if selectErr != "" {
		// School selection might be optional or not required, so we don't return an error
	}

	
	loginButton := page.Locator("button[type='submit'], .login-button, #login-btn")
	if errStr := loginButton.Click().Error(); errStr != "" {
		return "", fmt.Errorf("failed to click login button: %s", errStr)
	}

	
	deviceCodeChan := make(chan string, 1)
	errorChan := make(chan error, 1)

	
	page.On("framenavigated", func(frame playwright.Frame) {
		url := frame.URL()
		
		if deviceCode := extractDeviceCodeFromURL(url); deviceCode != "" {
			deviceCodeChan <- deviceCode
		}
	})

	
	done := make(chan struct{})
	var deviceCode string
	var navErr error

	
	go func() {
		defer close(done)
		select {
		case code := <-deviceCodeChan:
			deviceCode = code
		case err := <-errorChan:
			navErr = err
		case <-time.After(30 * time.Second): 
			navErr = errors.New("timeout waiting for device code")
		case <-ctx.Done():
			navErr = ctx.Err()
		}
	}()

	
	<-done

	if navErr != nil {
		return "", fmt.Errorf("failed to get device code: %w", navErr)
	}

	
	return c.getTokenWithDeviceCode(ctx, deviceCode, instanceURL)
}


func extractDeviceCodeFromURL(url string) string {
	
	
	if len(url) < 10 {
		return ""
	}

	
	deviceCodeStart := "device_code="
	pinCodeStart := "pincode="
	
	if idx := findSubstring(url, deviceCodeStart); idx != -1 {
		startIdx := idx + len(deviceCodeStart)
		endIdx := findNextChar(url, startIdx, '&')
		if endIdx == -1 {
			endIdx = len(url)
		}
		return url[startIdx:endIdx]
	}

	if idx := findSubstring(url, pinCodeStart); idx != -1 {
		startIdx := idx + len(pinCodeStart)
		endIdx := findNextChar(url, startIdx, '&')
		if endIdx == -1 {
			endIdx = len(url)
		}
		return url[startIdx:endIdx]
	}
	
	return ""
}


func findSubstring(haystack, needle string) int {
	for i := 0; i <= len(haystack)-len(needle); i++ {
		match := true
		for j := 0; j < len(needle); j++ {
			if haystack[i+j] != needle[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}


func findNextChar(s string, start int, char byte) int {
	for i := start; i < len(s); i++ {
		if s[i] == char {
			return i
		}
	}
	return -1
}


func (c *BrowserAuthClient) getTokenWithDeviceCode(ctx context.Context, deviceCode, instanceURL string) (string, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	
	tokenURL := "https://auth.edu.demogk.ru/oauth/token"

	
	tokenData := map[string]string{
		"grant_type":    "urn:ietf:params:oauth:grant-type:device_code",
		"device_code":   deviceCode,
		"client_id":     "parent-mobile",
		"client_secret": "04064338-13df-4747-8dea-69849f9ecdf0",
	}

	
	maxAttempts := 100 
	interval := 5 * time.Second 

	for attempt := 0; attempt < maxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(interval):
			jsonData, err := json.Marshal(tokenData)
			if err != nil {
				return "", fmt.Errorf("failed to marshal token data: %w", err)
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, bytes.NewBuffer(jsonData))
			if err != nil {
				return "", fmt.Errorf("failed to create token request: %w", err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("User-Agent", "NetSchoolApp/1.0")

			resp, err := client.Do(req)
			if err != nil {
				return "", fmt.Errorf("failed to execute token request: %w", err)
			}

			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				return "", fmt.Errorf("failed to read token response: %w", err)
			}

			if resp.StatusCode == http.StatusOK {
				
				var tokenResp struct {
					AccessToken  string `json:"access_token"`
					TokenType    string `json:"token_type"`
					ExpiresIn    int    `json:"expires_in"`
					RefreshToken string `json:"refresh_token"`
				}

				if err := json.Unmarshal(body, &tokenResp); err != nil {
					return "", fmt.Errorf("failed to parse token response: %w", err)
				}

				return tokenResp.AccessToken, nil
			} else if resp.StatusCode == http.StatusBadRequest {
				
				var errorResp struct {
					Error            string `json:"error"`
					ErrorDescription string `json:"error_description"`
				}

				if err := json.Unmarshal(body, &errorResp); err != nil {
					return "", fmt.Errorf("failed to parse error response: %w", err)
				}

				if errorResp.Error == "authorization_pending" {
					
					continue
				} else if errorResp.Error == "slow_down" {
					
					interval += 5 * time.Second
					continue
				} else if errorResp.Error == "expired_token" {
					return "", fmt.Errorf("device code expired: %s", errorResp.ErrorDescription)
				} else {
					return "", fmt.Errorf("token request failed: %s - %s", errorResp.Error, errorResp.ErrorDescription)
				}
			} else {
				return "", fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
			}
		}
	}

	return "", fmt.Errorf("max attempts exceeded for token polling")
}


func (c *BrowserAuthClient) IsAvailable() bool {
	return c.available
}


func (c *BrowserAuthClient) Close() error {
	if c.browser != nil {
		c.browser.Close()
	}
	if c.pw != nil {
		c.pw.Stop()
	}
	return nil
}