package domain

type PVZWithReceptions struct {
	PVZ        PVZ                     `json:"pvz"`
	Receptions []ReceptionWithProducts `json:"receptions,omitempty"`
}
