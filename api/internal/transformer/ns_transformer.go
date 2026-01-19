package transformer

import (
	"encoding/json"
	"fmt"
	"time"
)

// NSTransformer handles data transformation between JavaScript client format and Go structs
type NSTransformer struct{}

// NewNSTransformer creates a new instance of NSTransformer
func NewNSTransformer() *NSTransformer {
	return &NSTransformer{}
}

// TransformStudentInfo transforms student information from JS client format to Go struct
func (t *NSTransformer) TransformStudentInfo(jsData interface{}) (map[string]interface{}, error) {
	// Convert interface{} to map[string]interface{}
	data, ok := jsData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for student info")
	}

	// Transform the data as needed
	transformed := make(map[string]interface{})
	
	// Copy basic fields
	for k, v := range data {
		transformed[k] = v
	}

	// Add any specific transformations here if needed
	// For example, converting date strings to time.Time objects
	
	return transformed, nil
}

// TransformDiary transforms diary data from JS client format to Go struct
func (t *NSTransformer) TransformDiary(jsData interface{}) (map[string]interface{}, error) {
	data, ok := jsData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for diary")
	}

	transformed := make(map[string]interface{})
	
	// Copy basic fields
	for k, v := range data {
		transformed[k] = v
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
			transformed["weekDays"] = transformedDays
		}
	}

	return transformed, nil
}

// transformDay transforms a single day from JS client format
func (t *NSTransformer) transformDay(dayData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})
	
	// Copy basic fields
	for k, v := range dayData {
		transformed[k] = v
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

// transformLesson transforms a single lesson from JS client format
func (t *NSTransformer) transformLesson(lessonData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})
	
	// Copy basic fields
	for k, v := range lessonData {
		transformed[k] = v
	}

	// Transform assignments if present
	if assignments, exists := lessonData["assignments"]; exists {
		if assignmentList, ok := assignments.([]interface{}); ok {
			transformedAssignments := make([]interface{}, len(assignmentList))
			for i, assignment := range assignmentList {
				if assignmentMap, ok := assignment.(map[string]interface{}); ok {
					transformedAssignments[i] = t.transformAssignment(assignmentMap)
				} else {
					transformedAssignments[i] = assignment
				}
			}
			transformed["assignments"] = transformedAssignments
		}
	}

	return transformed
}

// transformAssignment transforms a single assignment from JS client format
func (t *NSTransformer) transformAssignment(assignmentData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})
	
	// Copy basic fields
	for k, v := range assignmentData {
		transformed[k] = v
	}

	return transformed
}

// TransformGrades transforms grades data from JS client format to Go struct
func (t *NSTransformer) TransformGrades(jsData interface{}) (map[string]interface{}, error) {
	// For now, just return the data as-is since it's HTML content
	// In a real implementation, this would parse the HTML and extract structured data
	data, ok := jsData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for grades")
	}

	return data, nil
}

// TransformSchedule transforms schedule data from JS client format to Go struct
func (t *NSTransformer) TransformSchedule(jsData interface{}) (map[string]interface{}, error) {
	data, ok := jsData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for schedule")
	}

	transformed := make(map[string]interface{})
	
	// Copy basic fields
	for k, v := range data {
		transformed[k] = v
	}

	return transformed, nil
}

// TransformAssignmentTypes transforms assignment types from JS client format to Go struct
func (t *NSTransformer) TransformAssignmentTypes(jsData interface{}) ([]map[string]interface{}, error) {
	types, ok := jsData.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for assignment types")
	}

	transformed := make([]map[string]interface{}, len(types))
	for i, item := range types {
		if itemMap, ok := item.(map[string]interface{}); ok {
			transformed[i] = make(map[string]interface{})
			
			// Copy fields
			for k, v := range itemMap {
				transformed[i][k] = v
			}
		} else {
			return nil, fmt.Errorf("invalid item format in assignment types")
		}
	}

	return transformed, nil
}

// TransformSchoolInfo transforms school information from JS client format to Go struct
func (t *NSTransformer) TransformSchoolInfo(jsData interface{}) (map[string]interface{}, error) {
	data, ok := jsData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for school info")
	}

	transformed := make(map[string]interface{})
	
	// Copy basic fields
	for k, v := range data {
		transformed[k] = v
	}

	return transformed, nil
}

// TransformContext transforms context information from JS client format to Go struct
func (t *NSTransformer) TransformContext(jsData interface{}) (map[string]interface{}, error) {
	data, ok := jsData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for context")
	}

	transformed := make(map[string]interface{})
	
	// Copy basic fields
	for k, v := range data {
		transformed[k] = v
	}

	return transformed, nil
}

// TransformJournal transforms journal data from JS client format to Go struct
func (t *NSTransformer) TransformJournal(jsData interface{}) (map[string]interface{}, error) {
	// For now, just return the data as-is since it's HTML content
	// In a real implementation, this would parse the HTML and extract structured data
	data, ok := jsData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for journal")
	}

	return data, nil
}

// TransformInfo transforms user info from JS client format to Go struct
func (t *NSTransformer) TransformInfo(jsData interface{}) (map[string]interface{}, error) {
	data, ok := jsData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for info")
	}

	transformed := make(map[string]interface{})
	
	// Copy basic fields
	for k, v := range data {
		transformed[k] = v
	}

	return transformed, nil
}

// TransformAssignmentInfo transforms assignment info from JS client format to Go struct
func (t *NSTransformer) TransformAssignmentInfo(jsData interface{}) (map[string]interface{}, error) {
	data, ok := jsData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for assignment info")
	}

	transformed := make(map[string]interface{})
	
	// Copy basic fields
	for k, v := range data {
		transformed[k] = v
	}

	return transformed, nil
}

// ParseDate parses date string from JS client format to time.Time
func (t *NSTransformer) ParseDate(dateStr string) (time.Time, error) {
	// Try different date formats commonly used by the JS client
	formats := []string{
		"2006-01-02T15:04:05.000Z", // ISO 8601
		"2006-01-02",               // YYYY-MM-DD
		"02.01.06",                 // DD.MM.YY
		"02.01.2006",               // DD.MM.YYYY
	}

	for _, format := range formats {
		if parsed, err := time.Parse(format, dateStr); err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

// FormatDate formats time.Time to string in JS client compatible format
func (t *NSTransformer) FormatDate(date time.Time) string {
	return date.Format("2006-01-02")
}