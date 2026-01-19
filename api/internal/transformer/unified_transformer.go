package transformer

import (
	"encoding/json"
	"fmt"
	"time"
)

// UnifiedTransformer handles data transformation between different API formats and Go structs
type UnifiedTransformer struct{}

// NewUnifiedTransformer creates a new instance of UnifiedTransformer
func NewUnifiedTransformer() *UnifiedTransformer {
	return &UnifiedTransformer{}
}

// TransformStudentInfo transforms student information from various API formats to standardized Go struct
func (t *UnifiedTransformer) TransformStudentInfo(apiData interface{}) (map[string]interface{}, error) {
	// Convert interface{} to map[string]interface{}
	data, ok := apiData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for student info")
	}

	// Standardize the data structure
	transformed := make(map[string]interface{})
	
	// Map common fields regardless of source
	fieldMappings := map[string]string{
		"id":         "id",
		"student_id": "id",
		"userId":     "id",
		"firstName":  "first_name",
		"first_name": "first_name",
		"lastName":   "last_name",
		"last_name":  "last_name",
		"middleName": "middle_name",
		"middle_name": "middle_name",
		"birthDate":  "birth_date",
		"birth_date": "birth_date",
		"class":      "class",
		"classId":    "class_id",
		"class_id":   "class_id",
		"schoolId":   "school_id",
		"school_id":  "school_id",
		"email":      "email",
		"phone":      "phone",
		"mobilePhone": "phone",
		"mobile_phone": "phone",
		"existsPhoto": "has_photo",
		"has_photo":   "has_photo",
	}

	for sourceField, targetField := range fieldMappings {
		if value, exists := data[sourceField]; exists {
			transformed[targetField] = value
		}
	}

	// Add timestamp
	transformed["created_at"] = time.Now().UTC()

	return transformed, nil
}

// TransformDiary transforms diary data from various API formats to standardized Go struct
func (t *UnifiedTransformer) TransformDiary(apiData interface{}) (map[string]interface{}, error) {
	data, ok := apiData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for diary")
	}

	transformed := make(map[string]interface{})
	
	// Copy basic fields with standardization
	basicFields := []string{"termName", "term_name", "className", "class_name", "weekStart", "week_start", "weekEnd", "week_end"}
	for _, field := range basicFields {
		if value, exists := data[field]; exists {
			standardField := standardizeFieldName(field)
			transformed[standardField] = value
		}
	}

	// Transform weekDays if present
	if weekDays, exists := data["weekDays"]; exists {
		if days, ok := weekDays.([]interface{}); ok {
			transformedDays := make([]interface{}, len(days))
			for i, day := range days {
				if dayMap, ok := day.(map[string]interface{}); ok {
					transformedDays[i] = t.transformDay(dayMap)
				} else {
					transformedDays[i] = day
				}
			}
			transformed["week_days"] = transformedDays
		}
	} else if weekDays, exists := data["week_days"]; exists {
		// Already in standard format
		transformed["week_days"] = weekDays
	}

	transformed["created_at"] = time.Now().UTC()
	return transformed, nil
}

// transformDay transforms a single day from various API formats
func (t *UnifiedTransformer) transformDay(dayData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})
	
	// Standardize date field
	if dateVal, exists := dayData["date"]; exists {
		transformed["date"] = dateVal
	} else if dateVal, exists := dayData["dt"]; exists {
		transformed["date"] = dateVal
	}

	// Transform lessons if present
	if lessons, exists := dayData["lessons"]; exists {
		if lessonList, ok := lessons.([]interface{}); ok {
			transformedLessons := make([]interface{}, len(lessonList))
			for i, lesson := range lessonList {
				if lessonMap, ok := lesson.(map[string]interface{}); ok {
					transformedLessons[i] = t.transformLesson(lessonMap)
				} else {
					transformedLessons[i] = lesson
				}
			}
			transformed["lessons"] = transformedLessons
		}
	}

	return transformed
}

// transformLesson transforms a single lesson from various API formats
func (t *UnifiedTransformer) transformLesson(lessonData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})
	
	// Standardize lesson fields
	fieldMappings := map[string]string{
		"id":          "id",
		"lessonId":    "id",
		"lesson_id":   "id",
		"number":      "number",
		"lessonNumber": "number",
		"lesson_number": "number",
		"subject":     "subject",
		"subjectName": "subject",
		"subject_name": "subject",
		"teacher":     "teacher",
		"teacherName": "teacher",
		"teacher_name": "teacher",
		"room":        "room",
		"roomId":      "room",
		"room_id":     "room",
		"startTime":   "start_time",
		"start_time":  "start_time",
		"start":       "start_time",
		"endTime":     "end_time",
		"end_time":    "end_time",
		"end":         "end_time",
		"date":        "date",
		"homeWork":    "homework",
		"homework":    "homework",
		"hw":          "homework",
		"assignments": "assignments",
		"tasks":       "assignments",
	}

	for sourceField, targetField := range fieldMappings {
		if value, exists := lessonData[sourceField]; exists {
			transformed[targetField] = value
		}
	}

	return transformed
}

// transformAssignment transforms a single assignment from various API formats
func (t *UnifiedTransformer) transformAssignment(assignmentData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})
	
	// Standardize assignment fields
	fieldMappings := map[string]string{
		"id":          "id",
		"assignmentId": "id",
		"assignment_id": "id",
		"text":        "text",
		"content":     "text",
		"desc":        "text",
		"description": "text",
		"date":        "due_date",
		"dueDate":     "due_date",
		"due_date":    "due_date",
		"mark":        "mark",
		"grade":       "mark",
		"score":       "mark",
		"typeId":      "type_id",
		"type_id":     "type_id",
		"type":        "type",
		"comment":     "comment",
		"comments":    "comment",
		"lessonId":    "lesson_id",
		"lesson_id":   "lesson_id",
		"attachments": "attachments",
		"files":       "attachments",
		"dot":         "overdue",
		"overdue":     "overdue",
		"isDeleted":   "deleted",
		"deleted":     "deleted",
		"weight":      "weight",
		"value":       "value",
	}

	for sourceField, targetField := range fieldMappings {
		if value, exists := assignmentData[sourceField]; exists {
			transformed[targetField] = value
		}
	}

	return transformed
}

// TransformGrades transforms grades data from various API formats to standardized Go struct
func (t *UnifiedTransformer) TransformGrades(apiData interface{}) (map[string]interface{}, error) {
	// For now, just return the data with standardization
	// In a real implementation, this would parse the HTML and extract structured data
	data, ok := apiData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for grades")
	}

	transformed := make(map[string]interface{})
	
	// Copy and standardize fields
	for k, v := range data {
		transformed[standardizeFieldName(k)] = v
	}

	transformed["created_at"] = time.Now().UTC()
	return transformed, nil
}

// TransformSchedule transforms schedule data from various API formats to standardized Go struct
func (t *UnifiedTransformer) TransformSchedule(apiData interface{}) (map[string]interface{}, error) {
	data, ok := apiData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for schedule")
	}

	transformed := make(map[string]interface{})
	
	// Standardize schedule fields
	for k, v := range data {
		transformed[standardizeFieldName(k)] = v
	}

	// Transform days if present
	if days, exists := data["days"]; exists {
		if dayList, ok := days.([]interface{}); ok {
			transformedDays := make([]interface{}, len(dayList))
			for i, day := range dayList {
				if dayMap, ok := day.(map[string]interface{}); ok {
					transformedDays[i] = t.transformDay(dayMap)
				} else {
					transformedDays[i] = day
				}
			}
			transformed["days"] = transformedDays
		}
	}

	transformed["created_at"] = time.Now().UTC()
	return transformed, nil
}

// TransformAssignmentTypes transforms assignment types from various API formats to standardized Go struct
func (t *UnifiedTransformer) TransformAssignmentTypes(apiData interface{}) ([]map[string]interface{}, error) {
	types, ok := apiData.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for assignment types")
	}

	transformed := make([]map[string]interface{}, len(types))
	for i, item := range types {
		if itemMap, ok := item.(map[string]interface{}); ok {
			transformed[i] = make(map[string]interface{})
			
			// Standardize fields
			fieldMappings := map[string]string{
				"id":       "id",
				"typeId":   "id",
				"type_id":  "id",
				"name":     "name",
				"title":    "name",
				"abbr":     "abbreviation",
				"short":    "abbreviation",
				"order":    "order",
				"priority": "order",
			}

			for sourceField, targetField := range fieldMappings {
				if value, exists := itemMap[sourceField]; exists {
					transformed[i][targetField] = value
				}
			}
		} else {
			return nil, fmt.Errorf("invalid item format in assignment types")
		}
	}

	return transformed, nil
}

// TransformSchoolInfo transforms school information from various API formats to standardized Go struct
func (t *UnifiedTransformer) TransformSchoolInfo(apiData interface{}) (map[string]interface{}, error) {
	data, ok := apiData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for school info")
	}

	transformed := make(map[string]interface{})
	
	// Standardize school fields
	fieldMappings := map[string]string{
		"id":             "id",
		"schoolId":       "id",
		"school_id":      "id",
		"name":           "name",
		"schoolName":     "name",
		"school_name":    "name",
		"fullName":       "full_name",
		"full_name":      "full_name",
		"address":        "address",
		"addr":           "address",
		"phone":          "phone",
		"telephone":      "phone",
		"email":          "email",
		"principal":      "principal",
		"director":       "principal",
		"foundationYear": "foundation_year",
		"foundation_year": "foundation_year",
		"website":        "website",
		"url":            "website",
		"inn":            "inn",
		"ogrn":           "ogrn",
	}

	for sourceField, targetField := range fieldMappings {
		if value, exists := data[sourceField]; exists {
			transformed[targetField] = value
		}
	}

	transformed["created_at"] = time.Now().UTC()
	return transformed, nil
}

// TransformContext transforms context information from various API formats to standardized Go struct
func (t *UnifiedTransformer) TransformContext(apiData interface{}) (map[string]interface{}, error) {
	data, ok := apiData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for context")
	}

	transformed := make(map[string]interface{})
	
	// Copy and standardize fields
	for k, v := range data {
		transformed[standardizeFieldName(k)] = v
	}

	transformed["created_at"] = time.Now().UTC()
	return transformed, nil
}

// TransformJournal transforms journal data from various API formats to standardized Go struct
func (t *UnifiedTransformer) TransformJournal(apiData interface{}) (map[string]interface{}, error) {
	// For now, just return the data with standardization
	// In a real implementation, this would parse the HTML and extract structured data
	data, ok := apiData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for journal")
	}

	transformed := make(map[string]interface{})
	
	// Copy and standardize fields
	for k, v := range data {
		transformed[standardizeFieldName(k)] = v
	}

	transformed["created_at"] = time.Now().UTC()
	return transformed, nil
}

// TransformInfo transforms user info from various API formats to standardized Go struct
func (t *UnifiedTransformer) TransformInfo(apiData interface{}) (map[string]interface{}, error) {
	data, ok := apiData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for info")
	}

	transformed := make(map[string]interface{})
	
	// Standardize user info fields
	fieldMappings := map[string]string{
		"email":        "email",
		"mobilePhone":  "phone",
		"mobile_phone": "phone",
		"phone":        "phone",
		"firstName":    "first_name",
		"first_name":   "first_name",
		"lastName":     "last_name",
		"last_name":    "last_name",
		"middleName":   "middle_name",
		"middle_name":  "middle_name",
		"birthDate":    "birth_date",
		"birth_date":   "birth_date",
		"existsPhoto":  "has_photo",
		"has_photo":    "has_photo",
		"userId":       "user_id",
		"user_id":      "user_id",
		"userName":     "username",
		"username":     "username",
	}

	for sourceField, targetField := range fieldMappings {
		if value, exists := data[sourceField]; exists {
			transformed[targetField] = value
		}
	}

	transformed["created_at"] = time.Now().UTC()
	return transformed, nil
}

// TransformAssignmentInfo transforms assignment info from various API formats to standardized Go struct
func (t *UnifiedTransformer) TransformAssignmentInfo(apiData interface{}) (map[string]interface{}, error) {
	data, ok := apiData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for assignment info")
	}

	transformed := make(map[string]interface{})
	
	// Standardize assignment info fields
	fieldMappings := map[string]string{
		"id":          "id",
		"assignmentId": "id",
		"assignment_id": "id",
		"date":        "date",
		"dueDate":     "date",
		"text":        "text",
		"content":     "text",
		"desc":        "text",
		"description": "text",
		"weight":      "weight",
		"subject":     "subject",
		"subjectName": "subject",
		"teacher":     "teacher",
		"teacherName": "teacher",
		"isDeleted":   "deleted",
		"deleted":     "deleted",
	}

	for sourceField, targetField := range fieldMappings {
		if value, exists := data[sourceField]; exists {
			transformed[targetField] = value
		}
	}

	transformed["created_at"] = time.Now().UTC()
	return transformed, nil
}

// ParseDate parses date string from various API formats to time.Time
func (t *UnifiedTransformer) ParseDate(dateStr string) (time.Time, error) {
	// Try different date formats commonly used by various APIs
	formats := []string{
		"2006-01-02T15:04:05.000Z", // ISO 8601 with milliseconds
		"2006-01-02T15:04:05Z",     // ISO 8601
		"2006-01-02T15:04:05",      // ISO 8601 without Z
		"2006-01-02",               // YYYY-MM-DD
		"02.01.06",                 // DD.MM.YY
		"02.01.2006",               // DD.MM.YYYY
		"02/01/2006",               // DD/MM/YYYY
		"01/02/2006",               // MM/DD/YYYY (US format)
		"Jan 2, 2006",              // Month DD, YYYY
		"January 2, 2006",          // Full month DD, YYYY
		"2006-01-02 15:04:05",      // YYYY-MM-DD HH:MM:SS
		"02.01.2006 15:04:05",     // DD.MM.YYYY HH:MM:SS
	}

	for _, format := range formats {
		if parsed, err := time.Parse(format, dateStr); err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

// FormatDate formats time.Time to string in API compatible format
func (t *UnifiedTransformer) FormatDate(date time.Time) string {
	return date.Format("2006-01-02T15:04:05Z")
}

// standardizeFieldName converts various field name formats to a standard format
func standardizeFieldName(fieldName string) string {
	// Common field name variations to standard form
	mappings := map[string]string{
		"weekStart":    "week_start",
		"weekEnd":      "week_end",
		"className":    "class_name",
		"termName":     "term_name",
		"firstName":    "first_name",
		"lastName":     "last_name",
		"middleName":   "middle_name",
		"birthDate":    "birth_date",
		"mobilePhone":  "mobile_phone",
		"schoolId":     "school_id",
		"studentId":    "student_id",
		"assignmentId": "assignment_id",
		"lessonId":     "lesson_id",
		"subjectName":  "subject_name",
		"teacherName":  "teacher_name",
		"startTime":    "start_time",
		"endTime":      "end_time",
		"homeWork":     "homework",
		"isDeleted":    "deleted",
		"existsPhoto":  "has_photo",
		"foundationYear": "foundation_year",
		"typeId":       "type_id",
		"dueDate":      "due_date",
		"lessonNumber": "lesson_number",
	}

	if standardized, exists := mappings[fieldName]; exists {
		return standardized
	}
	return fieldName
}