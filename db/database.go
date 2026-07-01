package db

import (
	"context"
	"log"
	"os" // Interacts with target machine operational environment variables

	"github.com/jackc/pgx/v5/pgxpool" // Imports high-performance PostgreSQL driver framework pool manager
)

var Pool *pgxpool.Pool

func CreateTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS department (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL UNIQUE
	);

	CREATE TABLE IF NOT EXISTS employees (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) UNIQUE,
		department_id INT,
		CONSTRAINT fk_department
			FOREIGN KEY (department_id)
			REFERENCES departments(id)

	);

	CREATE TABLE IF NOT EXISTS employee_details (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100),
		email VARCHAR(100),
		department_name VARCHAR(100)
	);
	`

	_, err := Pool.Exec(context.Background(), query)
	return err
}

func InitDB() {
	var err error // instantiates tracker checking functional operation error occurrences
	connstr := os.Getenv("DATABASE_URL")
	if connstr == "" {
		connstr = "postgres://postgres:Admin%40123@localhost:5432/postgres?sslmode=disable"
	}
	Pool, err = pgxpool.New(context.Background(), connstr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	if err := CreateTables(); err != nil {
		log.Fatalf("CreateTables failed: %v", err)
	}
}
