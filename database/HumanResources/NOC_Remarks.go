// Package database handles database connections and queries related to inbox data.
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Creator: Sridharan
//
// Created On:21-07-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 21-07-2025
package databasehumanresources

import (
	"Hrmodule/dbconfig"
	modelshumanresources "Hrmodule/models/HumanResources"
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb" // SQL Server driver
)

func NOCremarksdatabase(w http.ResponseWriter, r *http.Request) ([]modelshumanresources.NOCremarksstructure, int, error) {

	// Connection string for MSSQL Server
	connectionString := credentials.GetTestdatabase15()

	// Open connection to the database
	db, err := sql.Open("mssql", connectionString)
	if err != nil {
		return nil, 0, fmt.Errorf("error opening database connection: %v", err)
	}
	defer db.Close()
	// Extract the 'Taskid' parameter from the URL query
	Taskid := r.URL.Query().Get("Taskid")
	if Taskid == "" {
		return nil, 0, fmt.Errorf("missing 'Taskid' parameter")
	}

	// Execute the query using the student's Taskid
	rows, err := db.Query(modelshumanresources.MyQueryNOCremarks, Taskid)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	// Map the query result to this struct
	NOCremarksapi, err := modelshumanresources.RetrieveNOCremarks(rows, Taskid)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving data: %v", err)
	}

	// Count total records returned
	totalCount := len(NOCremarksapi)
	return NOCremarksapi, totalCount, nil
}
