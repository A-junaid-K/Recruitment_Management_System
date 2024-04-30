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
