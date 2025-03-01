package main

import (
	// Go Internal Packages
	"context"
	"fmt"
	"log"
	"math/rand"

	// External Packages
	"github.com/jackc/pgx/v4"
	"github.com/jaswdr/faker"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "postgres"
	PG_HOST     = "localhost"
	PG_PORT     = "5432"
)

var JobTitles = []struct {
	Title     string
	SalaryMin float64
	SalaryMax float64
}{
	{"Software Engineer", 60000, 150000},
	{"Data Scientist", 70000, 160000},
	{"Product Manager", 80000, 180000},
	{"DevOps Engineer", 65000, 140000},
	{"HR Manager", 50000, 100000},
}

func connectDB() (*pgx.Conn, error) {
	URI := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", DB_USER, DB_PASSWORD, PG_HOST, PG_PORT, DB_NAME)
	conn, err := pgx.Connect(context.Background(), URI)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func CreateTables(conn *pgx.Conn) error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS employees (
			employee_id SERIAL PRIMARY KEY,
			first_name VARCHAR(50),
			last_name VARCHAR(50),
			email VARCHAR(100) UNIQUE,
			phone_number VARCHAR(30),
			gender VARCHAR(10),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS salary (
			id SERIAL PRIMARY KEY,
			employee_id INT REFERENCES employees(employee_id) ON DELETE CASCADE,
			job_title VARCHAR(100),
			salary NUMERIC(10, 2),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, query := range tables {
		_, err := conn.Exec(context.Background(), query)
		if err != nil {
			return err
		}
	}
	fmt.Println("Tables created successfully.")
	return nil
}

func InsertRandomData(conn *pgx.Conn, numRecords int) error {
	fake := faker.New()

	for i := 0; i < numRecords; i++ {
		firstName := fake.Person().FirstName()
		lastName := fake.Person().LastName()
		email := fake.Internet().Email()
		phoneNumber := fake.Phone().Number()
		gender := []string{"Male", "Female"}[rand.Intn(2)]

		var employeeID int
		err := conn.QueryRow(context.Background(),
			`INSERT INTO employees (first_name, last_name, email, phone_number, gender)
			 VALUES ($1, $2, $3, $4, $5) RETURNING employee_id`,
			firstName, lastName, email, phoneNumber, gender).Scan(&employeeID)

		if err != nil {
			log.Println("Error inserting employee due to: ", err)
			continue
		}

		// Assign job and salary
		job := JobTitles[rand.Intn(len(JobTitles))]
		salary := job.SalaryMin + rand.Float64()*(job.SalaryMax-job.SalaryMin)

		_, err = conn.Exec(context.Background(),
			`INSERT INTO salary (employee_id, job_title, salary)
			 VALUES ($1, $2, $3)`,
			employeeID, job.Title, salary)

		if err != nil {
			log.Println("Error inserting salary due to: ", err)
		}
	}
	fmt.Println("Data Populated Successfully")
	return nil
}

func main() {
	conn, err := connectDB()
	if err != nil {
		log.Fatal("Failed to connect to database due to: ", err)
	}
	defer conn.Close(context.Background())

	err = CreateTables(conn)
	if err != nil {
		log.Fatal("Failed to create tables due to: ", err)
	}

	err = InsertRandomData(conn, 1000)
	if err != nil {
		log.Fatal("Failed to insert data due to: ", err)
	}
}
