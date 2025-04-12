package domain

type ReceptionWithProducts struct {
	Reception Reception `json:"reception"`
	Products  []Product `json:"products,omitempty"`
}
