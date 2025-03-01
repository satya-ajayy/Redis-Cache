package main

import "time"

type Employee struct {
	EmployeeID  int       `json:"employee_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Gender      string    `json:"gender"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Salary struct {
	EmployeeID int       `json:"employee_id"`
	JobTitle   string    `json:"job_title"`
	Salary     float64   `json:"salary"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type EmpSalary struct {
	EmployeeID int     `json:"employee_id"`
	Name       string  `json:"name"`
	Gender     string  `json:"gender"`
	JobTitle   string  `json:"job_title"`
	Salary     float64 `json:"salary"`
}

type Error struct {
	Message string `json:"error"`
}
