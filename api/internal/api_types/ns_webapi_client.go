package api_types

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
	"strconv"
	"time"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)


type NSWebAPIClient struct {
	timeout   time.Duration
	retryMax  int
	retryWait time.Duration
}


func (c *NSWebAPIClient) Login(ctx context.Context, username, password string, schoolID int, instanceURL string, loginData map[string]interface{}) (string, error) {
	
	client := &http.Client{Timeout: c.timeout}
	
	_, err := client.Get(fmt.Sprintf("%s/webapi/logindata", instanceURL))
	if err != nil {
		return "", fmt.Errorf("failed to get logindata cookies: %w", err)
	}

	
	loginMeta := make(map[string]interface{})
	for k, v := range loginData {
		loginMeta[k] = v
	}

	
	salt, ok := loginMeta["salt"].(string)
	if !ok {
		return "", fmt.Errorf("salt not found in login data")
	}

	
	delete(loginMeta, "salt")

	
	encoder := charmap.Windows1251.NewEncoder()
	win1251Pass, _, err := transform.String(encoder, password)
	if err != nil {
		return "", fmt.Errorf("failed to encode password: %w", err)
	}

	encodedPassword := md5.Sum([]byte(win1251Pass))
	pw2Hash := md5.Sum(append([]byte(salt), encodedPassword[:]...))
	pw2 := hex.EncodeToString(pw2Hash[:])
	pw := pw2[:len(password)]

	
	loginDataFinal := map[string]string{
		"loginType": "1", 
		"scid":      strconv.Itoa(schoolID), 
		"un":        username,
		"pw":        pw,
		"pw2":       pw2,
	}

	
	for k, v := range loginMeta {
		loginDataFinal[k] = fmt.Sprintf("%v", v)
	}

	
	formData := url.Values{}
	for k, v := range loginDataFinal {
		formData.Set(k, v)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/webapi/login", instanceURL), bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create login request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "NetCityApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute login request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read login response: %w", err)
	}

	var authResult map[string]interface{}
	if err := json.Unmarshal(body, &authResult); err != nil {
		return "", fmt.Errorf("failed to parse login response: %w", err)
	}

	
	at, exists := authResult["at"]
	if !exists {
		message, msgExists := authResult["message"]
		if !msgExists {
			message = "Unknown authentication error"
		}
		return "", fmt.Errorf("authentication failed: %v", message)
	}

	accessToken, ok := at.(string)
	if !ok {
		return "", fmt.Errorf("access token is not a string")
	}

	return accessToken, nil
}


func (c *NSWebAPIClient) GetLoginData(ctx context.Context, instanceURL string) (map[string]interface{}, error) {
	client := &http.Client{Timeout: c.timeout}
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/webapi/logindata", instanceURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}


func (c *NSWebAPIClient) GetStudentInfo(ctx context.Context, userID, instanceURL string) (interface{}, error) {
	client := &http.Client{Timeout: c.timeout}
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/webapi/student/diary/init", instanceURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	
	req.Header.Set("at", userID) 
	req.Header.Set("User-Agent", "NetCityApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}


func (c *NSWebAPIClient) GetGrades(ctx context.Context, userID, studentID, instanceURL string) (interface{}, error) {
	client := &http.Client{Timeout: c.timeout}
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/webapi/student/%s/grades", instanceURL, studentID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	
	req.Header.Set("at", userID) 
	req.Header.Set("User-Agent", "NetCityApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}


func (c *NSWebAPIClient) GetSchedule(ctx context.Context, userID, instanceURL string, weekStart time.Time) (interface{}, error) {
	client := &http.Client{Timeout: c.timeout}
	
	
	yearResp, err := c.getCurrentYear(ctx, userID, instanceURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get current year: %w", err)
	}
	
	yearID, ok := yearResp["id"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to get year ID")
	}

	
	weekFormatted := weekStart.Format("2006-01-02")
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/webapi/student/%s/schedule", instanceURL, userID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	
	q := req.URL.Query()
	q.Set("yearId", yearID)
	q.Set("week", weekFormatted)
	req.URL.RawQuery = q.Encode()

	
	req.Header.Set("at", userID) 
	req.Header.Set("User-Agent", "NetCityApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}


func (c *NSWebAPIClient) GetSchoolInfo(ctx context.Context, userID, instanceURL string) (interface{}, error) {
	client := &http.Client{Timeout: c.timeout}
	
	endpoints := []string{
		"schools/current",
		"school/info", 
		"schools/my",
	}
	
	for _, endpoint := range endpoints {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/webapi/%s", instanceURL, endpoint), nil)
		if err != nil {
			continue 
		}

		
		req.Header.Set("at", userID) 
		req.Header.Set("User-Agent", "NetCityApp/1.0")

		resp, err := client.Do(req)
		if err != nil {
			continue 
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue 
		}

		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			continue 
		}

		
		if len(result) > 0 {
			return result, nil
		}
	}

	return nil, fmt.Errorf("failed to get school info from any endpoint")
}


func (c *NSWebAPIClient) GetClasses(ctx context.Context, userID, instanceURL string) (interface{}, error) {
	client := &http.Client{Timeout: c.timeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/webapi/classes", instanceURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	
	req.Header.Set("at", userID) 
	req.Header.Set("User-Agent", "NetCityApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}


func (c *NSWebAPIClient) CheckHealth(ctx context.Context, instanceURL string) (bool, error) {
	_, err := c.GetLoginData(ctx, instanceURL)
	return err == nil, err
}


func (c *NSWebAPIClient) CheckIntPing(ctx context.Context, instanceURL string) (bool, time.Duration, error) {
	start := time.Now()
	ok, err := c.CheckHealth(ctx, instanceURL)
	duration := time.Since(start)
	return ok, duration, err
}


func (c *NSWebAPIClient) GetDiary(ctx context.Context, userID, studentID, instanceURL string, start, end time.Time) (interface{}, error) {
	client := &http.Client{Timeout: c.timeout}

	
	yearResp, err := c.getCurrentYear(ctx, userID, instanceURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get current year: %w", err)
	}

	yearID, ok := yearResp["id"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to get year ID")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/webapi/student/diary", instanceURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	
	q := req.URL.Query()
	q.Set("yearId", yearID)
	q.Set("studentId", studentID)
	q.Set("weekEnd", end.Format("2006-01-02"))
	q.Set("weekStart", start.Format("2006-01-02"))
	req.URL.RawQuery = q.Encode()

	
	req.Header.Set("at", userID)
	req.Header.Set("User-Agent", "NetCityApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}


func (c *NSWebAPIClient) GetAssignment(ctx context.Context, userID, studentID, assignmentID, instanceURL string) (interface{}, error) {
	client := &http.Client{Timeout: c.timeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/webapi/grade/assignment/%s", instanceURL, assignmentID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	
	q := req.URL.Query()
	q.Set("studentId", studentID)
	req.URL.RawQuery = q.Encode()

	
	req.Header.Set("at", userID)
	req.Header.Set("User-Agent", "NetCityApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}


func (c *NSWebAPIClient) GetAssignmentTypes(ctx context.Context, userID, instanceURL string) (interface{}, error) {
	client := &http.Client{Timeout: c.timeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/webapi/grade/assignment/types", instanceURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	
	q := req.URL.Query()
	q.Set("all", "false")
	req.URL.RawQuery = q.Encode()

	
	req.Header.Set("at", userID)
	req.Header.Set("User-Agent", "NetCityApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}


func (c *NSWebAPIClient) GetDownloadFile(ctx context.Context, userID, studentID, assignmentID, fileID, instanceURL string) (interface{}, error) {
	client := &http.Client{Timeout: c.timeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/webapi/download/attachment/%s", instanceURL, fileID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	
	q := req.URL.Query()
	q.Set("studentId", studentID)
	q.Set("assignId", assignmentID)
	req.URL.RawQuery = q.Encode()

	
	req.Header.Set("at", userID)
	req.Header.Set("User-Agent", "NetCityApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	
	return body, nil
}


func (c *NSWebAPIClient) GetReportFile(ctx context.Context, userID, instanceURL, reportURL string, filters map[string]interface{}, yearID int, timeout int, transport *int) (interface{}, error) {
	client := &http.Client{Timeout: c.timeout}

	
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/webapi/%s", instanceURL, reportURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	
	filterData := make(map[string]interface{})
	for k, v := range filters {
		filterData[k] = v
	}

	
	if yearID > 0 {
		filterData["yearId"] = yearID
	}

	
	if transport != nil {
		filterData["transport"] = *transport
	}

	
	jsonData, err := json.Marshal(filterData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal filter data: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("at", userID)
	req.Header.Set("User-Agent", "NetCityApp/1.0")
	req.Body = io.NopCloser(bytes.NewReader(jsonData))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}


func (c *NSWebAPIClient) GetJournal(ctx context.Context, userID, studentID, instanceURL string, start, end time.Time, termID, classID int, transport *int) (interface{}, error) {
	client := &http.Client{Timeout: c.timeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/webapi/reports/studenttotal", instanceURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	
	q := req.URL.Query()
	q.Set("SID", studentID)
	q.Set("PCLID_IUP", fmt.Sprintf("%d", classID))
	q.Set("TERMID", fmt.Sprintf("%d", termID))
	q.Set("period", fmt.Sprintf("%s - %s", start.Format("2006-01-02"), end.Format("2006-01-02")))
	if transport != nil {
		q.Set("transport", fmt.Sprintf("%d", *transport))
	}
	req.URL.RawQuery = q.Encode()

	
	req.Header.Set("at", userID)
	req.Header.Set("User-Agent", "NetCityApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}


func (c *NSWebAPIClient) GetInfo(ctx context.Context, userID, instanceURL string) (interface{}, error) {
	client := &http.Client{Timeout: c.timeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/webapi/mysettings", instanceURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	
	req.Header.Set("at", userID)
	req.Header.Set("User-Agent", "NetCityApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}


func (c *NSWebAPIClient) GetPhoto(ctx context.Context, userID, studentID, instanceURL string) (interface{}, error) {
	client := &http.Client{Timeout: c.timeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/webapi/Photo", instanceURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	
	q := req.URL.Query()
	q.Set("studentId", studentID)
	req.URL.RawQuery = q.Encode()

	
	req.Header.Set("at", userID)
	req.Header.Set("User-Agent", "NetCityApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	
	return body, nil
}


func (c *NSWebAPIClient) GetGradesForSubject(ctx context.Context, userID, studentID, subjectID, instanceURL string, start, end time.Time, termID, classID int, transport *int) (interface{}, error) {
	client := &http.Client{Timeout: c.timeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/webapi/student/%s/grades", instanceURL, studentID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	
	q := req.URL.Query()
	q.Set("subjectId", subjectID)
	q.Set("period", fmt.Sprintf("%s - %s", start.Format("2006-01-02"), end.Format("2006-01-02")))
	q.Set("termId", fmt.Sprintf("%d", termID))
	q.Set("classId", fmt.Sprintf("%d", classID))
	if transport != nil {
		q.Set("transport", fmt.Sprintf("%d", *transport))
	}
	req.URL.RawQuery = q.Encode()

	
	req.Header.Set("at", userID)
	req.Header.Set("User-Agent", "NetCityApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}


func (c *NSWebAPIClient) GetFullJournal(ctx context.Context, userID, studentID, instanceURL string, start, end time.Time, termID, classID int, transport *int) (interface{}, error) {
	client := &http.Client{Timeout: c.timeout}

	
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/webapi/reports/studenttotal", instanceURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	
	reportData := map[string]interface{}{
		"SID":         studentID,
		"PCLID_IUP":   classID,
		"TERMID":      termID,
		"period":      fmt.Sprintf("%s - %s", start.Format("2006-01-02"), end.Format("2006-01-02")),
		"transport":   transport,
	}

	jsonData, err := json.Marshal(reportData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal report data: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("at", userID)
	req.Header.Set("User-Agent", "NetCityApp/1.0")
	req.Body = io.NopCloser(bytes.NewReader(jsonData))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}


func (c *NSWebAPIClient) getCurrentYear(ctx context.Context, userID, instanceURL string) (map[string]interface{}, error) {
	client := &http.Client{Timeout: c.timeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/webapi/years/current", instanceURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	
	req.Header.Set("at", userID) 
	req.Header.Set("User-Agent", "NetCityApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}