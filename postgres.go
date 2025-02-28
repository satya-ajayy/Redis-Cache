package main

import (
	// Go Internal Packages
	"context"
	"log"

	// External Packages
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresSQL struct {
	DBPool *pgxpool.Pool
}

func NewPostgresSQL(URI string) (*PostgresSQL, error) {
	DBPool, err := pgxpool.Connect(context.Background(), URI)
	if err != nil {
		log.Fatalf("Error creating connection DBPool due to %v", err)
	}

	err = DBPool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Could not connect to Postgres due to %v", err)
	}
	return &PostgresSQL{DBPool: DBPool}, nil
}

func (p *PostgresSQL) Close() {
	p.DBPool.Close()
}

func (p *PostgresSQL) GetEmployeeSalary(empID string) (*EmpSalary, error) {
	query := `
		SELECT e.employee_id, e.first_name || ' ' || e.last_name AS name, e.gender, s.job_title, s.salary
		FROM employees e
		JOIN salary s ON e.employee_id = s.employee_id
		WHERE e.employee_id = $1;
	`

	var empSalary EmpSalary
	err := p.DBPool.QueryRow(context.Background(), query, empID).
		Scan(&empSalary.EmployeeID, &empSalary.Name, &empSalary.Gender, &empSalary.JobTitle, &empSalary.Salary)
	if err != nil {
		log.Printf("Error executing query due to %v", err)
		return nil, err
	}
	return &empSalary, nil
}
