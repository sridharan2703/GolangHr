// Package database handles database connections and queries related to inbox data.
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Creator: Sridharan
//
// Created On:15-07-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 15-07-2025
package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	modelscommon "Hrmodule/models/common"
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb" // SQL Server driver
)

func Taskinboxdatabase(w http.ResponseWriter, r *http.Request) ([]modelscommon.Taskinboxstructure, int, error) {
	// Connection string for MSSQL Server
	connectionString := credentials.GetTestdatabasetwo15()

	// Get query parameters
	empIdParam := r.URL.Query().Get("EmployeeID")
	assignedRoleParam := r.URL.Query().Get("AssignedRole")

	// Validate: both can't be empty
	if empIdParam == "" && assignedRoleParam == "" {
		http.Error(w, "EmployeeID and AssignedRole cannot both be empty", http.StatusBadRequest)
		return nil, 0, fmt.Errorf("missing parameters: EmployeeID and AssignedRole are required")
	}

	// Convert to SQL-compatible values (nil if empty)
	var empId interface{}
	if empIdParam == "" {
		empId = nil
	} else {
		empId = empIdParam
	}

	var assignedRole interface{}
	if assignedRoleParam == "" {
		assignedRole = nil
	} else {
		assignedRole = assignedRoleParam
	}

	// Open connection to the database
	db, err := sql.Open("mssql", connectionString)
	if err != nil {
		return nil, 0, fmt.Errorf("error opening database connection: %v", err)
	}
	defer db.Close()

	// Execute the query with parameters
	rows, err := db.Query(modelscommon.MyQuerytaskinbox, empId, assignedRole)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	// Map the query result
	taskinboxapi, err := modelscommon.Retrievetaskinbox(rows, empId, assignedRole)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving data: %v", err)
	}

	// Return result with count
	totalCount := len(taskinboxapi)
	return taskinboxapi, totalCount, nil
}
