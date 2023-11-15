package models

type Job struct {
	ID          int    `db:"id"`
	CustomerID  int    `db:"customer_id"`
	Name        string `db:"name"`
	CargoWeight int    `db:"cargo_weight"`
}

type JobAndLoaders struct {
	JobId   int   `json:"job_id"`
	Loaders []int `json:"loaders"`
}
