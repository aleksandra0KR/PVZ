package domain

import "time"

type PVZ struct {
	ID               *string    `db:"id" json:"id,omitempty"`
	RegistrationDate *time.Time `db:"registration_date" json:"registrationDate,omitempty"`
	City             *string    `db:"city" json:"city"`
}
