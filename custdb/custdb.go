package custdb

type CustomerDB struct {
	// ...
}

type Customer struct {
	PhoneNumber string
	Name        string
}

func NewCustomerDB() *CustomerDB {
	return &CustomerDB{}
}
