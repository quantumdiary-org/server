package api_types

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)



type NSMobileAPIClient struct {
	timeout   time.Duration
	retryMax  int
	retryWait time.Duration
}


func (c *NSMobileAPIClient) Login(ctx context.Context, username, password string, schoolID int, instanceURL string, loginData map[string]interface{}) (string, error) {
	client := &http.Client{Timeout: c.timeout}

	
	deviceCodeURL := fmt.Sprintf("%s/connect/deviceauthorization", instanceURL)

	deviceCodeData := url.Values{}
	deviceCodeData.Set("client_id", "parent-mobile")
	deviceCodeData.Set("scope", "mobile-api")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, deviceCodeURL, strings.NewReader(deviceCodeData.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create device code request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "NetSchoolApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get device code: %w", err)
	}

	deviceCodeBody, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", fmt.Errorf("failed to read device code response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get device code, status: %d, body: %s", resp.StatusCode, string(deviceCodeBody))
	}

	var deviceCodeResp struct {
		DeviceCode      string `json:"device_code"`
		UserCode        string `json:"user_code"`
		VerificationURI string `json:"verification_uri"`
		ExpiresIn       int    `json:"expires_in"`
		Interval        int    `json:"interval"`
	}

	if err := json.Unmarshal(deviceCodeBody, &deviceCodeResp); err != nil {
		return "", fmt.Errorf("failed to parse device code response: %w", err)
	}

	
	tokenURL := fmt.Sprintf("%s/connect/token", instanceURL)

	
	maxAttempts := deviceCodeResp.ExpiresIn / deviceCodeResp.Interval
	if maxAttempts == 0 {
		maxAttempts = 100 
	}
	interval := time.Duration(deviceCodeResp.Interval) * time.Second
	if interval == 0 {
		interval = 5 * time.Second 
	}

	for attempt := 0; attempt < maxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(interval):
			tokenData := url.Values{}
			tokenData.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
			tokenData.Set("device_code", deviceCodeResp.DeviceCode)
			tokenData.Set("client_id", "parent-mobile")
			tokenData.Set("client_secret", "04064338-13df-4747-8dea-69849f9ecdf0")

			tokenReq, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(tokenData.Encode()))
			if err != nil {
				return "", fmt.Errorf("failed to create token request: %w", err)
			}

			tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			tokenReq.Header.Set("User-Agent", "NetSchoolApp/1.0")

			tokenResp, err := client.Do(tokenReq)
			if err != nil {
				return "", fmt.Errorf("failed to get token: %w", err)
			}

			tokenBody, err := io.ReadAll(tokenResp.Body)
			tokenResp.Body.Close()
			if err != nil {
				return "", fmt.Errorf("failed to read token response: %w", err)
			}

			if tokenResp.StatusCode == http.StatusOK {
				
				var tokenResult struct {
					AccessToken  string `json:"access_token"`
					TokenType    string `json:"token_type"`
					ExpiresIn    int    `json:"expires_in"`
					RefreshToken string `json:"refresh_token"`
				}

				if err := json.Unmarshal(tokenBody, &tokenResult); err != nil {
					return "", fmt.Errorf("failed to parse token response: %w", err)
				}

				return tokenResult.AccessToken, nil
			} else if tokenResp.StatusCode == http.StatusBadRequest {
				
				var errorResult struct {
					Error            string `json:"error"`
					ErrorDescription string `json:"error_description"`
				}

				if err := json.Unmarshal(tokenBody, &errorResult); err != nil {
					return "", fmt.Errorf("failed to parse error response: %w", err)
				}

				if errorResult.Error == "authorization_pending" {
					
					continue
				} else if errorResult.Error == "slow_down" {
					
					interval += 5 * time.Second
					continue
				} else if errorResult.Error == "expired_token" {
					return "", fmt.Errorf("device code expired: %s", errorResult.ErrorDescription)
				} else {
					return "", fmt.Errorf("token request failed: %s - %s", errorResult.Error, errorResult.ErrorDescription)
				}
			} else {
				return "", fmt.Errorf("token request failed with status %d: %s", tokenResp.StatusCode, string(tokenBody))
			}
		}
	}

	return "", fmt.Errorf("max attempts exceeded for token polling")
}


func (c *NSMobileAPIClient) GetLoginData(ctx context.Context, instanceURL string) (map[string]interface{}, error) {
	
	loginData := map[string]interface{}{
		"auth_method": "oauth_device_flow",
		"client_id":   "parent-mobile",
		"scopes":      []string{"mobile-api"},
		"instance_url": instanceURL,
	}

	return loginData, nil
}


func (c *NSMobileAPIClient) GetStudentInfo(ctx context.Context, userID, instanceURL string) (interface{}, error) {
	client := &http.Client{Timeout: c.timeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/students", instanceURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+userID)
	req.Header.Set("User-Agent", "NetSchoolApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get student info, status: %d, body: %s", resp.StatusCode, string(body))
	}

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}


func (c *NSMobileAPIClient) GetGrades(ctx context.Context, userID, studentID, instanceURL string) (interface{}, error) {
	client := &http.Client{Timeout: c.timeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/grades", instanceURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	
	q := req.URL.Query()
	q.Set("studentId", studentID)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+userID)
	req.Header.Set("User-Agent", "NetSchoolApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get grades, status: %d, body: %s", resp.StatusCode, string(body))
	}

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}


func (c *NSMobileAPIClient) GetSchedule(ctx context.Context, userID, instanceURL string, weekStart time.Time) (interface{}, error) {
	client := &http.Client{Timeout: c.timeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/classmeetings", instanceURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	
	q := req.URL.Query()
	q.Set("weekStart", weekStart.Format("2006-01-02"))
	q.Set("weekEnd", weekStart.AddDate(0, 0, 6).Format("2006-01-02"))
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+userID)
	req.Header.Set("User-Agent", "NetSchoolApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get schedule, status: %d, body: %s", resp.StatusCode, string(body))
	}

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}


func (c *NSMobileAPIClient) GetSchoolInfo(ctx context.Context, userID, instanceURL string) (interface{}, error) {
	client := &http.Client{Timeout: c.timeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/education", instanceURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+userID)
	req.Header.Set("User-Agent", "NetSchoolApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get school info, status: %d, body: %s", resp.StatusCode, string(body))
	}

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}


func (c *NSMobileAPIClient) GetClasses(ctx context.Context, userID, instanceURL string) (interface{}, error) {
	client := &http.Client{Timeout: c.timeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/students/class", instanceURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+userID)
	req.Header.Set("User-Agent", "NetSchoolApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get classes, status: %d, body: %s", resp.StatusCode, string(body))
	}

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}


func (c *NSMobileAPIClient) GetDiary(ctx context.Context, userID, studentID, instanceURL string, start, end time.Time) (interface{}, error) {
	client := &http.Client{Timeout: c.timeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/diary", instanceURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	
	q := req.URL.Query()
	q.Set("studentId", studentID)
	q.Set("startDate", start.Format("2006-01-02"))
	q.Set("endDate", end.Format("2006-01-02"))
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+userID)
	req.Header.Set("User-Agent", "NetSchoolApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get diary, status: %d, body: %s", resp.StatusCode, string(body))
	}

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}


func (c *NSMobileAPIClient) GetAssignment(ctx context.Context, userID, studentID, assignmentID, instanceURL string) (interface{}, error) {
	client := &http.Client{Timeout: c.timeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/assignments/%s", instanceURL, assignmentID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	
	q := req.URL.Query()
	q.Set("studentId", studentID)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+userID)
	req.Header.Set("User-Agent", "NetSchoolApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get assignment, status: %d, body: %s", resp.StatusCode, string(body))
	}

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}


func (c *NSMobileAPIClient) GetAssignmentTypes(ctx context.Context, userID, instanceURL string) (interface{}, error) {
	client := &http.Client{Timeout: c.timeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/assignmentTypes", instanceURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+userID)
	req.Header.Set("User-Agent", "NetSchoolApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get assignment types, status: %d, body: %s", resp.StatusCode, string(body))
	}

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}


func (c *NSMobileAPIClient) GetDownloadFile(ctx context.Context, userID, studentID, assignmentID, fileID, instanceURL string) (interface{}, error) {
	client := &http.Client{Timeout: c.timeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/attachments/%s", instanceURL, fileID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	
	q := req.URL.Query()
	q.Set("studentId", studentID)
	q.Set("assignmentId", assignmentID)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+userID)
	req.Header.Set("User-Agent", "NetSchoolApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to download file, status: %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return body, nil
}


func (c *NSMobileAPIClient) GetReportFile(ctx context.Context, userID, instanceURL, reportURL string, filters map[string]interface{}, yearID int, timeout int, transport *int) (interface{}, error) {
	client := &http.Client{Timeout: c.timeout}

	
	fullURL := fmt.Sprintf("%s/%s", instanceURL, reportURL)

	
	requestBody, err := json.Marshal(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal filters: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+userID)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "NetSchoolApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get report file, status: %d, body: %s", resp.StatusCode, string(body))
	}

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}


func (c *NSMobileAPIClient) GetJournal(ctx context.Context, userID, studentID, instanceURL string, start, end time.Time, termID, classID int, transport *int) (interface{}, error) {
	client := &http.Client{Timeout: c.timeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/journal", instanceURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	
	q := req.URL.Query()
	q.Set("studentId", studentID)
	q.Set("startDate", start.Format("2006-01-02"))
	q.Set("endDate", end.Format("2006-01-02"))
	q.Set("termId", fmt.Sprintf("%d", termID))
	q.Set("classId", fmt.Sprintf("%d", classID))
	if transport != nil {
		q.Set("transport", fmt.Sprintf("%d", *transport))
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+userID)
	req.Header.Set("User-Agent", "NetSchoolApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get journal, status: %d, body: %s", resp.StatusCode, string(body))
	}

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}


func (c *NSMobileAPIClient) GetInfo(ctx context.Context, userID, instanceURL string) (interface{}, error) {
	client := &http.Client{Timeout: c.timeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/info", instanceURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+userID)
	req.Header.Set("User-Agent", "NetSchoolApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get info, status: %d, body: %s", resp.StatusCode, string(body))
	}

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}


func (c *NSMobileAPIClient) GetPhoto(ctx context.Context, userID, studentID, instanceURL string) (interface{}, error) {
	client := &http.Client{Timeout: c.timeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/photo", instanceURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	
	q := req.URL.Query()
	q.Set("studentId", studentID)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+userID)
	req.Header.Set("User-Agent", "NetSchoolApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get photo, status: %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return body, nil
}


func (c *NSMobileAPIClient) GetGradesForSubject(ctx context.Context, userID, studentID, subjectID, instanceURL string, start, end time.Time, termID, classID int, transport *int) (interface{}, error) {
	client := &http.Client{Timeout: c.timeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/grades", instanceURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	
	q := req.URL.Query()
	q.Set("studentId", studentID)
	q.Set("subjectId", subjectID)
	q.Set("startDate", start.Format("2006-01-02"))
	q.Set("endDate", end.Format("2006-01-02"))
	q.Set("termId", fmt.Sprintf("%d", termID))
	q.Set("classId", fmt.Sprintf("%d", classID))
	if transport != nil {
		q.Set("transport", fmt.Sprintf("%d", *transport))
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+userID)
	req.Header.Set("User-Agent", "NetSchoolApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get grades for subject, status: %d, body: %s", resp.StatusCode, string(body))
	}

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}


func (c *NSMobileAPIClient) GetFullJournal(ctx context.Context, userID, studentID, instanceURL string, start, end time.Time, termID, classID int, transport *int) (interface{}, error) {
	client := &http.Client{Timeout: c.timeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/full-journal", instanceURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	
	q := req.URL.Query()
	q.Set("studentId", studentID)
	q.Set("startDate", start.Format("2006-01-02"))
	q.Set("endDate", end.Format("2006-01-02"))
	q.Set("termId", fmt.Sprintf("%d", termID))
	q.Set("classId", fmt.Sprintf("%d", classID))
	if transport != nil {
		q.Set("transport", fmt.Sprintf("%d", *transport))
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+userID)
	req.Header.Set("User-Agent", "NetSchoolApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get full journal, status: %d, body: %s", resp.StatusCode, string(body))
	}

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}


func (c *NSMobileAPIClient) CheckHealth(ctx context.Context, instanceURL string) (bool, error) {
	client := &http.Client{Timeout: c.timeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/ping", instanceURL), nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}


func (c *NSMobileAPIClient) CheckIntPing(ctx context.Context, instanceURL string) (bool, time.Duration, error) {
	client := &http.Client{Timeout: c.timeout}

	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/int-ping", instanceURL), nil)
	if err != nil {
		return false, 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, 0, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	duration := time.Since(start)

	return resp.StatusCode == http.StatusOK, duration, nil
}