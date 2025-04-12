package domain

type User struct {
	ID       *string `db:"id" json:"id,omitempty"`
	Email    *string `db:"email" json:"email"`
	Password *string `db:"password" json:"-"`
	Role     *string `db:"role" json:"role"`
}
