package models

type Customer struct {
	ID int
	*User
	Capital int
}

type CustomerAndAllLoaders struct {
	Customer *Customer
	Loaders  []Loader
}
