package domain

import "time"

type Product struct {
	ID          *string    `db:"id" json:"id,omitempty"`
	DateTime    *time.Time `db:"date_time" json:"dateTime,omitempty"`
	Type        *string    `db:"type" json:"type"`
	ReceptionID *string    `db:"reception_id" json:"receptionId"`
}
