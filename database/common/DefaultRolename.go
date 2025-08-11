// Package database handles database connections and queries related to DefaultRoleName data.
//
// --- Creator's Info ---
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
	"Hrmodule/models/common"
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb" // SQL Server driver
)


func DefaultRoleNamedatabase(w http.ResponseWriter, r *http.Request) ([]modelscommon.DefaultRoleNamestructure, int, error) {
	// Connection string for MSSQL Server
	connectionString := credentials.GetTestdatabasetwo15()

	// Open connection to the database
	db, err := sql.Open("mssql", connectionString)
	if err != nil {
		return nil, 0, fmt.Errorf("error opening database connection: %v", err)
	}
	defer db.Close()

	// Extract the 'LoginName' parameter from the URL query
	LoginName := r.URL.Query().Get("LoginName")
	if LoginName == "" {
		return nil, 0, fmt.Errorf("missing 'LoginName' parameter")
	}

	// Execute the query using the student's LoginName
	rows, err := db.Query(modelscommon.MyQueryDefaultRoleName, LoginName)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	// Map the query result to this struct
	DefaultRoleNameapi, err := modelscommon.RetrieveDefaultRoleName(rows, LoginName)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving data: %v", err)
	}

	// Count total records returned
	totalCount := len(DefaultRoleNameapi)
	return DefaultRoleNameapi, totalCount, nil
}
