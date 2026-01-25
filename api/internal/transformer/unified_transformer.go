package transformer

import (
	"fmt"
	"time"
)


type UnifiedTransformer struct{}


func NewUnifiedTransformer() *UnifiedTransformer {
	return &UnifiedTransformer{}
}


func (t *UnifiedTransformer) TransformStudentInfo(apiData interface{}) (map[string]interface{}, error) {
	
	data, ok := apiData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for student info")
	}

	
	transformed := make(map[string]interface{})
	
	
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

	
	transformed["created_at"] = time.Now().UTC()

	return transformed, nil
}


func (t *UnifiedTransformer) TransformDiary(apiData interface{}) (map[string]interface{}, error) {
	data, ok := apiData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for diary")
	}

	transformed := make(map[string]interface{})
	
	
	basicFields := []string{"termName", "term_name", "className", "class_name", "weekStart", "week_start", "weekEnd", "week_end"}
	for _, field := range basicFields {
		if value, exists := data[field]; exists {
			standardField := standardizeFieldName(field)
			transformed[standardField] = value
		}
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
			transformed["week_days"] = transformedDays
		}
	} else if weekDays, exists := data["week_days"]; exists {
		
		transformed["week_days"] = weekDays
	}

	transformed["created_at"] = time.Now().UTC()
	return transformed, nil
}


func (t *UnifiedTransformer) transformDay(dayData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})
	
	
	if dateVal, exists := dayData["date"]; exists {
		transformed["date"] = dateVal
	} else if dateVal, exists := dayData["dt"]; exists {
		transformed["date"] = dateVal
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


func (t *UnifiedTransformer) transformLesson(lessonData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})
	
	
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


func (t *UnifiedTransformer) transformAssignment(assignmentData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})
	
	
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


func (t *UnifiedTransformer) TransformGrades(apiData interface{}) (map[string]interface{}, error) {
	
	
	data, ok := apiData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for grades")
	}

	transformed := make(map[string]interface{})
	
	
	for k, v := range data {
		transformed[standardizeFieldName(k)] = v
	}

	transformed["created_at"] = time.Now().UTC()
	return transformed, nil
}


func (t *UnifiedTransformer) TransformSchedule(apiData interface{}) (map[string]interface{}, error) {
	data, ok := apiData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for schedule")
	}

	transformed := make(map[string]interface{})
	
	
	for k, v := range data {
		transformed[standardizeFieldName(k)] = v
	}

	
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


func (t *UnifiedTransformer) TransformAssignmentTypes(apiData interface{}) ([]map[string]interface{}, error) {
	types, ok := apiData.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for assignment types")
	}

	transformed := make([]map[string]interface{}, len(types))
	for i, item := range types {
		if itemMap, ok := item.(map[string]interface{}); ok {
			transformed[i] = make(map[string]interface{})
			
			
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


func (t *UnifiedTransformer) TransformSchoolInfo(apiData interface{}) (map[string]interface{}, error) {
	data, ok := apiData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for school info")
	}

	transformed := make(map[string]interface{})
	
	
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


func (t *UnifiedTransformer) TransformContext(apiData interface{}) (map[string]interface{}, error) {
	data, ok := apiData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for context")
	}

	transformed := make(map[string]interface{})
	
	
	for k, v := range data {
		transformed[standardizeFieldName(k)] = v
	}

	transformed["created_at"] = time.Now().UTC()
	return transformed, nil
}


func (t *UnifiedTransformer) TransformJournal(apiData interface{}) (map[string]interface{}, error) {
	
	
	data, ok := apiData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for journal")
	}

	transformed := make(map[string]interface{})
	
	
	for k, v := range data {
		transformed[standardizeFieldName(k)] = v
	}

	transformed["created_at"] = time.Now().UTC()
	return transformed, nil
}


func (t *UnifiedTransformer) TransformInfo(apiData interface{}) (map[string]interface{}, error) {
	data, ok := apiData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for info")
	}

	transformed := make(map[string]interface{})
	
	
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


func (t *UnifiedTransformer) TransformAssignmentInfo(apiData interface{}) (map[string]interface{}, error) {
	data, ok := apiData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format for assignment info")
	}

	transformed := make(map[string]interface{})
	
	
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


func (t *UnifiedTransformer) ParseDate(dateStr string) (time.Time, error) {
	
	formats := []string{
		"2006-01-02T15:04:05.000Z", 
		"2006-01-02T15:04:05Z",     
		"2006-01-02T15:04:05",      
		"2006-01-02",               
		"02.01.06",                 
		"02.01.2006",               
		"02/01/2006",               
		"01/02/2006",               
		"Jan 2, 2006",              
		"January 2, 2006",          
		"2006-01-02 15:04:05",      
		"02.01.2006 15:04:05",     
	}

	for _, format := range formats {
		if parsed, err := time.Parse(format, dateStr); err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}


func (t *UnifiedTransformer) FormatDate(date time.Time) string {
	return date.Format("2006-01-02T15:04:05Z")
}


func standardizeFieldName(fieldName string) string {
	
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