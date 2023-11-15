package models

type Loader struct {
	ID int `db:"id"`
	*User
	MaxPortableWeight int `db:"max_portable_weight"`
	Fatigue           int `db:"fatigue"`
	Salary            int `db:"salary"`
	Drunkenness       int `db:"drunkenness"`
}
