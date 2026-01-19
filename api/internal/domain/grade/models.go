package grade

// Grade represents a student's grade
type Grade struct {
	ID          string `json:"id"`
	StudentID   string `json:"student_id"`
	SubjectID   string `json:"subject_id"`
	Value       string `json:"value"`
	Date        string `json:"date"`
	Description string `json:"description"`
	TeacherID   string `json:"teacher_id"`
	Weight      int    `json:"weight"`
}

// GradeService defines the interface for grade-related operations
type GradeService interface {
	GetGradesForStudent(studentID string) ([]*Grade, error)
	GetGradesForSubject(studentID, subjectID string) ([]*Grade, error)
	AddGrade(grade *Grade) error
	UpdateGrade(gradeID string, grade *Grade) error
	DeleteGrade(gradeID string) error
}