package domain

import "time"

type PVZ struct {
	ID               *string    `db:"id" json:"id"`
	RegistrationDate *time.Time `db:"registration_date" json:"registrationDate"`
	City             *string    `db:"city" json:"city"`
}
