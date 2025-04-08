package domain

import "time"

type Reception struct {
	ID       string    `db:"id"`
	DateTime time.Time `db:"date_time"`
	PVZID    string    `db:"pvz_id"`
	Status   string    `db:"status"`
}
