package student


type Student struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	BirthDate string `json:"birth_date"`
	Class     string `json:"class"`
	SchoolID  int    `json:"school_id"`
}


type StudentService interface {
	GetStudentInfo(studentID string) (*Student, error)
	GetStudentsByClass(classID string) ([]*Student, error)
	UpdateStudentProfile(studentID string, profile *Student) error
}