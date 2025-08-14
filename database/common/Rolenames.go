// Package database handles database connections and queries related to DefaultRoleName data.
//
// --- Creator's Info ---
//
// Creator: Sivabala
//
// Creator: Sivabala
//
// Created On:30-07-2025
//
// Last Modified By: Sivabala
//
// Last Modified Date: 30-07-2025
package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelscommon "Hrmodule/models/common"
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb" // SQL Server driver
)

func Rolesdatabase(w http.ResponseWriter, r *http.Request) ([]modelscommon.Rolesstructure, int, error) {

	// Connection string for MYSQL Server
	connectionString := credentials.GetMySQLDatabase17_HR()

	// Open connection to the database
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, 0, fmt.Errorf("error opening database connection: %v", err)
	}
	defer db.Close()
	// Extract the 'RoleName' parameter from the URL query
	RoleName := r.URL.Query().Get("RoleName")
	if RoleName == "" {
		return nil, 0, fmt.Errorf("missing 'RoleName' parameter")
	}

	// Execute the query using the student's RoleName
	rows, err := db.Query(modelscommon.MyQueryRoles, RoleName)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	// Map the query result to this struct
	Rolesapi, err := modelscommon.RetrieveRoles(rows, RoleName)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving data: %v", err)
	}

	// Count total records returned
	totalCount := len(Rolesapi)
	return Rolesapi, totalCount, nil
}
