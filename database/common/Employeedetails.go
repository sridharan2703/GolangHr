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
		"Hrmodule/dbconfig"
	 "Hrmodule/models/common"
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb" // SQL Server driver

)

func Employeedetailsdatabase(w http.ResponseWriter, r *http.Request) ([]modelscommon.Employeedetailsstructure, int, error) {


	// Connection string for MSSQL Server
	connectionString := credentials.GetTestdatabase15()


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
	rows, err := db.Query(modelscommon.MyQueryEmployeedetails, LoginName)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	// Map the query result to this struct
	employeedetailsapi, err := modelscommon.Retrieveemployeedetails(rows, LoginName)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving data: %v", err)
	}

	// Count total records returned
	totalCount := len(employeedetailsapi)
	return employeedetailsapi, totalCount, nil
}
