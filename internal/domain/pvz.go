package domain

import "time"

type PVZ struct {
	ID               string    `db:"id"`
	RegistrationDate time.Time `db:"registration_date"`
	City             string    `db:"city"`
}
