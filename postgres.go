package main

import (
	// Go Internal Packages
	"context"
	"fmt"
	"log"

	// External Packages
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "postgres"
	PG_HOST     = "localhost"
	PG_PORT     = "5432"
)

type PostgresSQL struct {
	DBPool *pgxpool.Pool
}

func NewPostgresSQL() (*PostgresSQL, error) {
	URI := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", DB_USER, DB_PASSWORD, PG_HOST, PG_PORT, DB_NAME)
	DBPool, err := pgxpool.Connect(context.Background(), URI)
	if err != nil {
		return nil, err
	}

	err = DBPool.Ping(context.Background())
	if err != nil {
		return nil, err
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
