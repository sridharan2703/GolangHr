// Package credentials provides functions to generate and verify
// database connection strings using environment variables.
// It supports both MSSQL and MySQL.
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Creator: Sridharan
//
// Created On:30-07-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 30-07-2025
package credentials

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// init loads environment variables from a .env file.
// It panics if the file cannot be loaded.
func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Failed to load .env: " + err.Error())
	}
}

// getDBConnectionString constructs and verifies a database connection string
// for the given driver (e.g., "mssql" or "mysql"). It opens and pings the
// database to ensure the connection is valid. It returns the connection string
// or panics on failure.
func getDBConnectionString(driver, server, user, password, database, port string) string {
	var connStr string

	switch driver {
	case "mssql":
		connStr = fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;port=%s", server, user, password, database, port)
	case "mysql":
		connStr = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, server, port, database)
	default:
		panic("Unsupported DB driver: " + driver)
	}

	db, err := sql.Open(driver, connStr)
	if err != nil {
		panic("Failed to open DB connection: " + err.Error())
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic("Database connection failed: " + err.Error())
	}

	return connStr
}

// GetTestdatabase15 returns a verified MSSQL connection string for
// the test credentials with IITMAcademics database from environment variables
func GetTestdatabase15() string {
	return getDBConnectionString(
		"mssql",
		os.Getenv("server1"),
		os.Getenv("userId1"),
		os.Getenv("password1"),
		os.Getenv("database1"),
		os.Getenv("port1"),
	)
}

// GetTestdatabasetwo15 returns a verified MSSQL connection string for
// the test credentials with Human Resources database from environment variables
func GetTestdatabasetwo15() string {
	return getDBConnectionString(
		"mssql",
		os.Getenv("server1"),
		os.Getenv("userId1"),
		os.Getenv("password1"),
		os.Getenv("database2"),
		os.Getenv("port1"),
	)
}

// GetMySQLDatabase17 returns a verified Mysql connection string for
// the test credentials with api_new database from environment variables
func GetMySQLDatabase17() string {
	return getDBConnectionString(
		"mysql",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
}

// GetMySQLDatabase17 returns a verified Mysql connection string for
// the test credentials with api_new database from environment variables
func GetMySQLDatabase17_HR() string {
	return getDBConnectionString(
		"mysql",
		os.Getenv("DB_HOST_HR"),
		os.Getenv("DB_USER_HR"),
		os.Getenv("DB_PASSWORD_HR"),
		os.Getenv("DB_NAME_HR"),
		os.Getenv("DB_PORT_HR"),
	)
}
