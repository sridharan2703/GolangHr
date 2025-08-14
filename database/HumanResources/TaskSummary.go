// Package database handles database connections and queries related to task summary data.
//
// --- Creator's Info ---
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
	credentials "Hrmodule/dbconfig"
	modelshumanresources "Hrmodule/models/HumanResources"
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb" // SQL Server driver
)

// RetrieveTaskSummaryFromDB retrieves task summary data from the database
func RetrieveTaskSummaryFromDB(w http.ResponseWriter, r *http.Request) ([]modelshumanresources.TaskSummarystructure, int, error) {

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
		return nil, 0, fmt.Errorf("missing required 'EmpId' parameter")
	}

	// Extract the 'Type' parameter from the URL query
	Type := r.URL.Query().Get("Type")
	if Type == "" {
		return nil, 0, fmt.Errorf("missing required 'Type' parameter")
	}

	rows, err := db.Query(modelshumanresources.MyQueryTaskSummary, EmpId, Type)
	if err != nil {
		return nil, 0, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	// Map the query result to this struct
	tasksummaryapi, err := modelshumanresources.Retrievetasksummary(rows, EmpId, Type)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving task summary data: %v", err)
	}

	// Count total records returned
	totalCount := len(tasksummaryapi)
	return tasksummaryapi, totalCount, nil
}
