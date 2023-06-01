package model

type Employee struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Address        string `json:"address"`
	EmployeeNumber string `json:"employeeNumber"`
}

type Address struct {
	Address string `json:"address"`
}
