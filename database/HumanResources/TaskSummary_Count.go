// Package database handles database connections and queries related to inbox data.
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Creator: Sridharan
//
// Created On:22-07-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 22-07-2025
package databasehumanresources

import (
	"Hrmodule/dbconfig"
	modelshumanresources "Hrmodule/models/HumanResources"
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb" // SQL Server driver
)

func Tasksummarydatabase(w http.ResponseWriter, r *http.Request) ([]modelshumanresources.TaskSummary_countstructure, int, error) {

	// Connection string for MSSQL Server
	connectionString := credentials.GetTestdatabasetwo15()

	// Open connection to the database
	db, err := sql.Open("mssql", connectionString)
	if err != nil {
		return nil, 0, fmt.Errorf("error opening database connection: %v", err)
	}
	defer db.Close()

	// Extract the 'EmpId' parameter from the URL query
	EmpId := r.URL.Query().Get("EmpId")
	if EmpId == "" {
		return nil, 0, fmt.Errorf("missing 'EmpId' parameter")
	}

	// Execute the query using the Employee's EmpId
	rows, err := db.Query(modelshumanresources.MyQueryTaskSummary_count, EmpId)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	// Map the query result to this struct
	tasksummarycountapi, err := modelshumanresources.Retrievetasksummarycount(rows, EmpId)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving data: %v", err)
	}

	// Count total records returned
	totalCount := len(tasksummarycountapi)
	return tasksummarycountapi, totalCount, nil
}
