package api_types

import (
	"context"
	"fmt"
	"time"
)


type DevMockAPIClient struct {
	Timeout time.Duration
}


func (c *DevMockAPIClient) Login(ctx context.Context, username, password string, schoolID int, instanceURL string, loginData map[string]interface{}) (string, error) {
	
	if username == "errorcode" && password == "errorcode" {
		return "", ErrAuthenticationFailed
	}

	if username == "nsfail" && password == "nsfail" {
		
		return "mock_token_for_nsfail", nil
	}

	
	return fmt.Sprintf("mock_token_for_%s", username), nil
}


func (c *DevMockAPIClient) GetLoginData(ctx context.Context, instanceURL string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"salt": "mock_salt_value",
		"lt":   "mock_lt_value",
		"ver":  "mock_ver_value",
	}, nil
}


func (c *DevMockAPIClient) GetStudentInfo(ctx context.Context, userID, instanceURL string) (interface{}, error) {
	studentInfo := map[string]interface{}{
		"id":         fmt.Sprintf("student_%s", userID),
		"first_name": "Тест",
		"last_name":  "Студент",
		"middle_name": "Тестович",
		"birth_date": "2000-01-01",
		"class":      "9А",
		"school_id":  1,
	}

	return studentInfo, nil
}


func (c *DevMockAPIClient) GetGrades(ctx context.Context, userID, studentID, instanceURL string) (interface{}, error) {
	grades := []interface{}{
		map[string]interface{}{
			"id":          fmt.Sprintf("grade_%s_1", userID),
			"student_id":  studentID,
			"subject_id":  "math",
			"value":       "5",
			"date":        "2023-09-15",
			"description": "Контрольная работа",
			"teacher_id":  "teacher_1",
			"weight":      10,
		},
		map[string]interface{}{
			"id":          fmt.Sprintf("grade_%s_2", userID),
			"student_id":  studentID,
			"subject_id":  "math",
			"value":       "4",
			"date":        "2023-09-10",
			"description": "Самостоятельная работа",
			"teacher_id":  "teacher_1",
			"weight":      5,
		},
	}

	return grades, nil
}


func (c *DevMockAPIClient) GetSchedule(ctx context.Context, userID, instanceURL string, weekStart time.Time) (interface{}, error) {
	
	days := make([]interface{}, 7)
	for i := 0; i < 7; i++ {
		currentDay := weekStart.AddDate(0, 0, i)

		var lessons []interface{}
		if i < 5 { 
			lessons = []interface{}{
				map[string]interface{}{
					"id":       fmt.Sprintf("lesson_%s_%d_1", userID, i),
					"number":   1,
					"subject":  "Математика",
					"teacher":  "Иванова А.А.",
					"room":     "301",
					"start":    currentDay.Add(8 * time.Hour).Add(30 * time.Minute),
					"end":      currentDay.Add(9 * time.Hour).Add(15 * time.Minute),
					"date":     currentDay,
				},
				map[string]interface{}{
					"id":       fmt.Sprintf("lesson_%s_%d_2", userID, i),
					"number":   2,
					"subject":  "Русский язык",
					"teacher":  "Петрова Б.Б.",
					"room":     "302",
					"start":    currentDay.Add(9 * time.Hour).Add(30 * time.Minute),
					"end":      currentDay.Add(10 * time.Hour).Add(15 * time.Minute),
					"date":     currentDay,
				},
			}
		}

		days[i] = map[string]interface{}{
			"date":    currentDay,
			"lessons": lessons,
		}
	}

	schedule := map[string]interface{}{
		"week_start": weekStart,
		"week_end":   weekStart.AddDate(0, 0, 6),
		"days":       days,
	}

	return schedule, nil
}


func (c *DevMockAPIClient) GetSchoolInfo(ctx context.Context, userID, instanceURL string) (interface{}, error) {
	schoolInfo := map[string]interface{}{
		"id":           1,
		"name":         "Тестовая школа №1",
		"address":      "ул. Тестовая, д. 1",
		"phone":        "+7 (123) 456-78-90",
		"email":        "test-school@example.com",
		"principal":    "Тестов Тест Тестович",
		"foundation_year": 1990,
		"website":      "https://test-school.edu.ru",
	}

	return schoolInfo, nil
}


func (c *DevMockAPIClient) GetClasses(ctx context.Context, userID, instanceURL string) (interface{}, error) {
	classes := []interface{}{
		map[string]interface{}{
			"id":              "class_1",
			"name":            "1А",
			"grade":           1,
			"letter":          "А",
			"teacher":         "Тестова Т.Т.",
			"students_count":  25,
		},
		map[string]interface{}{
			"id":              "class_9b",
			"name":            "9Б",
			"grade":           9,
			"letter":          "Б",
			"teacher":         "Тестов Т.Т.",
			"students_count":  22,
		},
	}

	return classes, nil
}


func (c *DevMockAPIClient) CheckHealth(ctx context.Context, instanceURL string) (bool, error) {
	return true, nil
}


func (c *DevMockAPIClient) CheckIntPing(ctx context.Context, instanceURL string) (bool, time.Duration, error) {
	start := time.Now()

	
	
	
	duration := time.Since(start)
	return true, duration, nil
}


func (c *DevMockAPIClient) GetDiary(ctx context.Context, userID, studentID, instanceURL string, start, end time.Time) (interface{}, error) {
	
	days := make([]interface{}, 0)
	current := start
	for current.Before(end) || current.Equal(end) {
		
		if current.Weekday() != time.Saturday && current.Weekday() != time.Sunday {
			lessons := []interface{}{
				map[string]interface{}{
					"id":       fmt.Sprintf("lesson_%s_%s", userID, current.Format("2006-01-02")),
					"number":   1,
					"subject":  "Математика",
					"teacher":  "Иванова А.А.",
					"room":     "301",
					"start":    current.Add(8 * time.Hour).Add(30 * time.Minute),
					"end":      current.Add(9 * time.Hour).Add(15 * time.Minute),
					"date":     current,
					"homework": "Выполнить задачи 1-5",
					"assignments": []interface{}{
						map[string]interface{}{
							"id":       fmt.Sprintf("assignment_%s_%s", userID, current.Format("2006-01-02")),
							"text":     "Домашнее задание",
							"date":     current,
							"mark":     5,
							"typeId":   1,
							"comment":  "Хорошая работа!",
						},
					},
				},
				map[string]interface{}{
					"id":       fmt.Sprintf("lesson2_%s_%s", userID, current.Format("2006-01-02")),
					"number":   2,
					"subject":  "Русский язык",
					"teacher":  "Петрова Б.Б.",
					"room":     "302",
					"start":    current.Add(9 * time.Hour).Add(30 * time.Minute),
					"end":      current.Add(10 * time.Hour).Add(15 * time.Minute),
					"date":     current,
					"homework": "Прочитать главу 3",
					"assignments": []interface{}{
						map[string]interface{}{
							"id":       fmt.Sprintf("assignment2_%s_%s", userID, current.Format("2006-01-02")),
							"text":     "Сочинение",
							"date":     current.AddDate(0, 0, 1),
							"mark":     4,
							"typeId":   2,
							"comment":  "Нужно доработать",
						},
					},
				},
			}

			days = append(days, map[string]interface{}{
				"date":    current,
				"lessons": lessons,
			})
		}

		current = current.AddDate(0, 0, 1)
	}

	diary := map[string]interface{}{
		"weekDays":   days,
		"termName":   "1 четверть",
		"className":  "9А",
		"weekStart":  start,
		"weekEnd":    end,
	}

	return diary, nil
}


func (c *DevMockAPIClient) GetAssignment(ctx context.Context, userID, studentID, assignmentID, instanceURL string) (interface{}, error) {
	assignment := map[string]interface{}{
		"id":          assignmentID,
		"date":        time.Now().AddDate(0, 0, -1),
		"text":        "Пример задания",
		"weight":      10,
		"subject":     "Математика",
		"teacher":     "Иванова А.А.",
		"isDeleted":   false,
		"description": "Описание примерного задания",
	}

	return assignment, nil
}


func (c *DevMockAPIClient) GetAssignmentTypes(ctx context.Context, userID, instanceURL string) (interface{}, error) {
	types := []interface{}{
		map[string]interface{}{
			"id":    1,
			"name":  "Домашнее задание",
			"abbr":  "ДЗ",
			"order": 1,
		},
		map[string]interface{}{
			"id":    2,
			"name":  "Контрольная работа",
			"abbr":  "КР",
			"order": 2,
		},
		map[string]interface{}{
			"id":    3,
			"name":  "Самостоятельная работа",
			"abbr":  "СР",
			"order": 3,
		},
	}

	return types, nil
}


func (c *DevMockAPIClient) GetDownloadFile(ctx context.Context, userID, studentID, assignmentID, fileID, instanceURL string) (interface{}, error) {
	
	mockFileContent := []byte("This is a mock file content for testing purposes.")
	return mockFileContent, nil
}


func (c *DevMockAPIClient) GetReportFile(ctx context.Context, userID, instanceURL, reportURL string, filters map[string]interface{}, yearID int, timeout int, transport *int) (interface{}, error) {
	report := map[string]interface{}{
		"status":  "success",
		"html":    "<html><body><h1>Mock Report Content</h1><p>This is a mock report for testing purposes.</p></body></html>",
		"filters": filters,
		"yearId":  yearID,
	}

	return report, nil
}


func (c *DevMockAPIClient) GetJournal(ctx context.Context, userID, studentID, instanceURL string, start, end time.Time, termID, classID int, transport *int) (interface{}, error) {
	journal := map[string]interface{}{
		"raw": "<html><body><h1>Mock Journal Content</h1><p>This is a mock journal for testing purposes.</p></body></html>",
		"range": map[string]interface{}{
			"start": start,
			"end":   end,
		},
		"subjects": []interface{}{
			map[string]interface{}{
				"id":   1,
				"name": "Математика",
				"marks": []interface{}{
					map[string]interface{}{
						"mark":    5,
						"date":    start.AddDate(0, 0, 1),
						"termId":  termID,
					},
					map[string]interface{}{
						"mark":    4,
						"date":    start.AddDate(0, 0, 3),
						"termId":  termID,
					},
				},
				"dotList": []interface{}{
					map[string]interface{}{
						"date":   start.AddDate(0, 0, 5),
						"termId": termID,
					},
				},
				"missedList": []interface{}{
					map[string]interface{}{
						"type":   "УП",
						"date":   start.AddDate(0, 0, 7),
						"termId": termID,
					},
				},
				"periodMiddleMark": 4.5,
			},
		},
	}

	return journal, nil
}


func (c *DevMockAPIClient) GetInfo(ctx context.Context, userID, instanceURL string) (interface{}, error) {
	info := map[string]interface{}{
		"email":        "test@example.com",
		"mobilePhone":  "+7 (999) 123-45-67",
		"firstName":    "Тест",
		"lastName":     "Пользователь",
		"middleName":   "Тестович",
		"birthDate":    "2000-01-01",
		"existsPhoto":  true,
	}

	return info, nil
}


func (c *DevMockAPIClient) GetPhoto(ctx context.Context, userID, studentID, instanceURL string) (interface{}, error) {
	
	mockImageContent := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A} 
	return mockImageContent, nil
}


func (c *DevMockAPIClient) GetGradesForSubject(ctx context.Context, userID, studentID, subjectID, instanceURL string, start, end time.Time, termID, classID int, transport *int) (interface{}, error) {
	grades := []interface{}{
		map[string]interface{}{
			"id":          fmt.Sprintf("grade_%s_%s", studentID, subjectID),
			"student_id":  studentID,
			"subject_id":  subjectID,
			"value":       "5",
			"date":        start.AddDate(0, 0, 1).Format("2006-01-02"),
			"description": "Контрольная работа",
			"teacher_id":  "teacher_1",
			"weight":      10,
		},
		map[string]interface{}{
			"id":          fmt.Sprintf("grade_%s_%s_2", studentID, subjectID),
			"student_id":  studentID,
			"subject_id":  subjectID,
			"value":       "4",
			"date":        start.AddDate(0, 0, 3).Format("2006-01-02"),
			"description": "Самостоятельная работа",
			"teacher_id":  "teacher_1",
			"weight":      5,
		},
	}

	return grades, nil
}


func (c *DevMockAPIClient) GetFullJournal(ctx context.Context, userID, studentID, instanceURL string, start, end time.Time, termID, classID int, transport *int) (interface{}, error) {
	journal := map[string]interface{}{
		"raw": "<html><body><h1>Mock Full Journal Content</h1><p>This is a mock full journal for testing purposes.</p></body></html>",
		"range": map[string]interface{}{
			"start": start,
			"end":   end,
		},
		"subjects": []interface{}{
			map[string]interface{}{
				"id":   1,
				"name": "Математика",
				"marks": []interface{}{
					map[string]interface{}{
						"mark":    5,
						"date":    start.AddDate(0, 0, 1),
						"termId":  termID,
					},
					map[string]interface{}{
						"mark":    4,
						"date":    start.AddDate(0, 0, 3),
						"termId":  termID,
					},
				},
				"dotList": []interface{}{
					map[string]interface{}{
						"date":   start.AddDate(0, 0, 5),
						"termId": termID,
					},
				},
				"missedList": []interface{}{
					map[string]interface{}{
						"type":   "УП",
						"date":   start.AddDate(0, 0, 7),
						"termId": termID,
					},
				},
				"periodMiddleMark": 4.5,
			},
			map[string]interface{}{
				"id":   2,
				"name": "Русский язык",
				"marks": []interface{}{
					map[string]interface{}{
						"mark":    5,
						"date":    start.AddDate(0, 0, 2),
						"termId":  termID,
					},
					map[string]interface{}{
						"mark":    5,
						"date":    start.AddDate(0, 0, 4),
						"termId":  termID,
					},
				},
				"dotList": []interface{}{
					map[string]interface{}{
						"date":   start.AddDate(0, 0, 6),
						"termId": termID,
					},
				},
				"missedList": []interface{}{
					map[string]interface{}{
						"type":   "Б",
						"date":   start.AddDate(0, 0, 8),
						"termId": termID,
					},
				},
				"periodMiddleMark": 4.8,
			},
		},
	}

	return journal, nil
}