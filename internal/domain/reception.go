package domain

import "time"

type Reception struct {
	ID       *string    `db:"id" json:"id,omitempty"`
	DateTime *time.Time `db:"date_time" json:"dateTime"`
	PvzId    *string    `db:"pvz_id" json:"pvzId"`
	Status   *string    `db:"status" json:"status"`
}
