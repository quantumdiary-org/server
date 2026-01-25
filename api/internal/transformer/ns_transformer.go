package transformer

import (
	"fmt"
	"time"
)


type NSTransformer struct{}


func NewNSTransformer() *NSTransformer {
	return &NSTransformer{}
}


func (t *NSTransformer) TransformStudentInfo(jsData interface{}) (map[string]interface{}, error) {
	
	data, ok := jsData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for student info")
	}

	
	transformed := make(map[string]interface{})
	
	
	for k, v := range data {
		transformed[k] = v
	}

	
	
	
	return transformed, nil
}


func (t *NSTransformer) TransformDiary(jsData interface{}) (map[string]interface{}, error) {
	data, ok := jsData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for diary")
	}

	transformed := make(map[string]interface{})
	
	
	for k, v := range data {
		transformed[k] = v
	}

	
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


func (t *NSTransformer) transformDay(dayData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})
	
	
	for k, v := range dayData {
		transformed[k] = v
	}

	
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


func (t *NSTransformer) transformLesson(lessonData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})
	
	
	for k, v := range lessonData {
		transformed[k] = v
	}

	
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


func (t *NSTransformer) transformAssignment(assignmentData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})
	
	
	for k, v := range assignmentData {
		transformed[k] = v
	}

	return transformed
}


func (t *NSTransformer) TransformGrades(jsData interface{}) (map[string]interface{}, error) {
	
	
	data, ok := jsData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for grades")
	}

	return data, nil
}


func (t *NSTransformer) TransformSchedule(jsData interface{}) (map[string]interface{}, error) {
	data, ok := jsData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for schedule")
	}

	transformed := make(map[string]interface{})
	
	
	for k, v := range data {
		transformed[k] = v
	}

	return transformed, nil
}


func (t *NSTransformer) TransformAssignmentTypes(jsData interface{}) ([]map[string]interface{}, error) {
	types, ok := jsData.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for assignment types")
	}

	transformed := make([]map[string]interface{}, len(types))
	for i, item := range types {
		if itemMap, ok := item.(map[string]interface{}); ok {
			transformed[i] = make(map[string]interface{})
			
			
			for k, v := range itemMap {
				transformed[i][k] = v
			}
		} else {
			return nil, fmt.Errorf("invalid item format in assignment types")
		}
	}

	return transformed, nil
}


func (t *NSTransformer) TransformSchoolInfo(jsData interface{}) (map[string]interface{}, error) {
	data, ok := jsData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for school info")
	}

	transformed := make(map[string]interface{})
	
	
	for k, v := range data {
		transformed[k] = v
	}

	return transformed, nil
}


func (t *NSTransformer) TransformContext(jsData interface{}) (map[string]interface{}, error) {
	data, ok := jsData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for context")
	}

	transformed := make(map[string]interface{})
	
	
	for k, v := range data {
		transformed[k] = v
	}

	return transformed, nil
}


func (t *NSTransformer) TransformJournal(jsData interface{}) (map[string]interface{}, error) {
	
	
	data, ok := jsData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for journal")
	}

	return data, nil
}


func (t *NSTransformer) TransformInfo(jsData interface{}) (map[string]interface{}, error) {
	data, ok := jsData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for info")
	}

	transformed := make(map[string]interface{})
	
	
	for k, v := range data {
		transformed[k] = v
	}

	return transformed, nil
}


func (t *NSTransformer) TransformAssignmentInfo(jsData interface{}) (map[string]interface{}, error) {
	data, ok := jsData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for assignment info")
	}

	transformed := make(map[string]interface{})
	
	
	for k, v := range data {
		transformed[k] = v
	}

	return transformed, nil
}


func (t *NSTransformer) ParseDate(dateStr string) (time.Time, error) {
	
	formats := []string{
		"2006-01-02T15:04:05.000Z", 
		"2006-01-02",               
		"02.01.06",                 
		"02.01.2006",               
	}

	for _, format := range formats {
		if parsed, err := time.Parse(format, dateStr); err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}


func (t *NSTransformer) FormatDate(date time.Time) string {
	return date.Format("2006-01-02")
}