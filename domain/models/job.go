package models

import "time"

type Job struct{
	Title string
	Description string
	PostedOn time.Time
	Total_applicant uint
	Company_name string
	PostBy User		// if admin
}