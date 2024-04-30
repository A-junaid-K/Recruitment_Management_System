package models

type Job struct {
	JobID          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Title          string `json:"title,omitempty"`
	Description    string `json:"description,omitempty"`
	PostedOn       string `json:"posted_on,omitempty"`
	TotalApplicant uint   `json:"total_applicant,omitempty"`
	CompanyName    string `json:"company_name,omitempty"`
	PostBy         string `json:"post_by,omitempty"`
}

type Apply struct {
	ApplyID      uint   `json:"job_id,omitempty"`
	ApplicantID  uint   `json:"applicant_id,omitempty"`
	Name         string `json:"name,omitempty"`
	Email        string `json:"email,omitempty"`
	Phone        string `json:"phone,omitempty"`
	Resume       string `json:"resume"`
	Skills       string `json:"skills,omitempty"`
	Education    string `json:"education,omitempty"`
	Experience   string
	CurrentCTC   uint `json:"current_ctc"`
	ExpectedCTC  uint `json:"expected_ctc"`
	NoticePeriod string `json:"notice_period,omitempty"`
}
