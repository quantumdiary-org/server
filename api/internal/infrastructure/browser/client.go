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

// BrowserAuthClientInterface определяет интерфейс для аутентификации через браузер
type BrowserAuthClientInterface interface {
	// Authenticate performs browser-based authentication and returns access token
	Authenticate(ctx context.Context, instanceURL, username, password string, schoolID int) (string, error)
	
	// IsAvailable checks if browser automation is available
	IsAvailable() bool
	
	// Close освобождает ресурсы
	Close() error
}

// ErrBrowserAutomationNotSupported ошибка, когда браузерная автоматизация не поддерживается
var ErrBrowserAutomationNotSupported = errors.New("this instance does not support Playwright browser automation")

// BrowserAuthClient реализует аутентификацию через браузер
// Если playwright не установлен, возвращает ErrBrowserAutomationNotSupported
type BrowserAuthClient struct {
	available bool
	pw        *playwright.Playwright
	browser   playwright.Browser
}

// NewBrowserAuthClient создает новый клиент для браузерной аутентификации
func NewBrowserAuthClient() (*BrowserAuthClient, error) {
	// Проверяем, доступен ли playwright
	available := checkPlaywrightAvailability()
	
	client := &BrowserAuthClient{
		available: available,
	}
	
	if available {
		// Инициализируем playwright
		pw, err := playwright.Run()
		if err != nil {
			// Если не удалось запустить playwright, помечаем как недоступный
			client.available = false
			return client, nil
		}
		
		// Запускаем браузер
		browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
			Headless: playwright.Bool(true), // Запуск в headless режиме
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

// checkPlaywrightAvailability проверяет, установлен ли playwright
func checkPlaywrightAvailability() bool {
	// Проверяем наличие бинарных файлов playwright
	cmd := exec.Command("playwright", "--version")
	if err := cmd.Run(); err != nil {
		// Если команда не выполнена, проверим установку через node
		cmd = exec.Command("node", "-e", "require('playwright')")
		if err := cmd.Run(); err != nil {
			return false
		}
	}
	
	return true
}

// Authenticate performs browser-based authentication and returns access token
func (c *BrowserAuthClient) Authenticate(ctx context.Context, instanceURL, username, password string, schoolID int) (string, error) {
	if !c.available {
		return "", ErrBrowserAutomationNotSupported
	}

	// Создаем новую страницу в браузере
	page, err := c.browser.NewPage()
	if err != nil {
		return "", fmt.Errorf("failed to create new page: %w", err)
	}
	defer page.Close()

	// Подготавливаем URL для входа
	loginURL := fmt.Sprintf("%s/login?mobile", instanceURL)

	// Переходим на страницу входа
	if _, err := page.Goto(loginURL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkIdle,
	}); err != nil {
		return "", fmt.Errorf("failed to navigate to login page: %w", err)
	}

	// Ожидаем загрузки формы входа
	if err := page.WaitForSelector("input[name='lg'], input[name='login'], input[name='username']").Error(); err != nil {
		return "", fmt.Errorf("login form not found: %w", err)
	}

	// Заполняем форму входа
	if err := page.Fill("input[name='lg'], input[name='login'], input[name='username']", username).Error(); err != nil {
		// Если не найдено поле lg, пробуем другие варианты
		if err := page.Fill("input[name='login'], input[name='username']", username).Error(); err != nil {
			return "", fmt.Errorf("failed to fill username: %w", err)
		}
	}

	if err := page.Fill("input[name='pw'], input[name='password']", password).Error(); err != nil {
		return "", fmt.Errorf("failed to fill password: %w", err)
	}

	// Выбираем школу из выпадающего списка
	if err := page.SelectOption("select[name='cl'], select[name='school'], select[name='schoolId']", playwright.PageSelectOptionOptions{
		Value: fmt.Sprintf("%d", schoolID),
	}).Error(); err != nil {
		// Если не удалось выбрать школу, продолжаем (возможно, она уже выбрана или не требуется)
	}

	// Нажимаем кнопку входа
	loginButton := page.Locator("button[type='submit'], .login-button, #login-btn")
	if err := loginButton.Click().Error(); err != nil {
		return "", fmt.Errorf("failed to click login button: %w", err)
	}

	// Ожидаем редиректа на кастомный протокол (irtech:// или другой)
	deviceCodeChan := make(chan string, 1)
	errorChan := make(chan error, 1)

	// Устанавливаем обработчик события редиректа
	page.On("framenavigated", func(frame playwright.Frame) {
		url := frame.URL()
		// Проверяем, содержит ли URL кастомный протокол с device_code
		if deviceCode := extractDeviceCodeFromURL(url); deviceCode != "" {
			deviceCodeChan <- deviceCode
		}
	})

	// Запускаем ожидание в отдельной горутине с таймаутом
	done := make(chan struct{})
	var deviceCode string
	var navErr error

	// Запускаем ожидание навигации в отдельной горутине
	go func() {
		defer close(done)
		select {
		case code := <-deviceCodeChan:
			deviceCode = code
		case err := <-errorChan:
			navErr = err
		case <-time.After(30 * time.Second): // Таймаут 30 секунд
			navErr = errors.New("timeout waiting for device code")
		case <-ctx.Done():
			navErr = ctx.Err()
		}
	}()

	// Ждем получения device_code или ошибки
	<-done

	if navErr != nil {
		return "", fmt.Errorf("failed to get device code: %w", navErr)
	}

	// Успешно получили device_code, теперь используем его для получения токена
	return c.getTokenWithDeviceCode(ctx, deviceCode, instanceURL)
}

// extractDeviceCodeFromURL извлекает device_code из URL редиректа
func extractDeviceCodeFromURL(url string) string {
	// Проверяем, содержит ли URL кастомный протокол с device_code
	// Пример: irtech://?device_code=abc123 или irtech://?pincode=def456
	if len(url) < 10 {
		return ""
	}

	// Ищем device_code в URL
	deviceCodeStart := "device_code="
	pinCodeStart := "pincode="
	
	if idx := findSubstring(url, deviceCodeStart); idx != -1 {
		startIdx := idx + len(deviceCodeStart)
		endIdx := findNextChar(url, startIdx, "&")
		if endIdx == -1 {
			endIdx = len(url)
		}
		return url[startIdx:endIdx]
	}
	
	if idx := findSubstring(url, pinCodeStart); idx != -1 {
		startIdx := idx + len(pinCodeStart)
		endIdx := findNextChar(url, startIdx, "&")
		if endIdx == -1 {
			endIdx = len(url)
		}
		return url[startIdx:endIdx]
	}
	
	return ""
}

// findSubstring находит индекс подстроки в строке
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

// findNextChar находит индекс следующего символа
func findNextChar(s string, start int, char byte) int {
	for i := start; i < len(s); i++ {
		if s[i] == char {
			return i
		}
	}
	return -1
}

// getTokenWithDeviceCode получает токен с использованием device_code через OAuth device flow
func (c *BrowserAuthClient) getTokenWithDeviceCode(ctx context.Context, deviceCode, instanceURL string) (string, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	// Подготавливаем URL для получения токена
	tokenURL := "https://identity.ir-tech.ru/connect/token"

	// Подготавливаем данные для запроса токена
	tokenData := map[string]string{
		"grant_type":    "urn:ietf:params:oauth:grant-type:device_code",
		"device_code":   deviceCode,
		"client_id":     "parent-mobile",
		"client_secret": "04064338-13df-4747-8dea-69849f9ecdf0",
	}

	// Выполняем polling для получения токена
	maxAttempts := 100 // Максимальное количество попыток
	interval := 5 * time.Second // Интервал между попытками

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
				// Успешно получили токен
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
				// Проверяем тип ошибки
				var errorResp struct {
					Error            string `json:"error"`
					ErrorDescription string `json:"error_description"`
				}

				if err := json.Unmarshal(body, &errorResp); err != nil {
					return "", fmt.Errorf("failed to parse error response: %w", err)
				}

				if errorResp.Error == "authorization_pending" {
					// Продолжаем ожидание
					continue
				} else if errorResp.Error == "slow_down" {
					// Увеличиваем интервал
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

// IsAvailable checks if browser automation is available
func (c *BrowserAuthClient) IsAvailable() bool {
	return c.available
}

// Close освобождает ресурсы playwright
func (c *BrowserAuthClient) Close() error {
	if c.browser != nil {
		c.browser.Close()
	}
	if c.pw != nil {
		c.pw.Stop()
	}
	return nil
}