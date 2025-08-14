// Package models contains data structures and database access logic for the EMployeedetails page.
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
//
// Path: Task inbox Page
package modelscommon

import (
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb" // MSSQL driver for database/sql
)

const MyQueryEmployeedetails = `SELECT 
    a.EmployeeId,
    a.FirstName,
    a.Department,
    LEFT(a.Department, CHARINDEX(',', a.Department + ',') - 1) AS DepartmentCode,
    d.DepartmentName,
    e.DOJ,
    e.EmployeeDesignationId,
    f.EmployeeDesignationName
FROM axon..ADSUserDetails a
LEFT JOIN iitm..DepartmentMaster d 
    ON LEFT(a.Department, CHARINDEX(',', a.Department + ',') - 1) = d.DepartmentCode
LEFT JOIN iitm..EmployeeAppointmentDetails e 
    ON a.EmployeeId = e.EmployeeId
LEFT JOIN iitm..EmployeeDesignationMaster f 
    ON e.EmployeeDesignationId = f.EmployeeDesignationId
WHERE a.LoginName = ?`

// Employeedetailsstructure defines the structure of each task in the inbox
type Employeedetailsstructure struct {
	EmployeeId              *string `json:"EmployeeId"`
	FirstName               *string `json:"FirstName"`
	Department              *string `json:"Department"`
	DepartmentCode          *string `json:"DepartmentCode"`
	DepartmentName          *string `json:"DepartmentName"`
	DOJ                     *string `json:"DOJ"`
	EmployeeDesignationId   *int    `json:"EmployeeDesignationId"`
	EmployeeDesignationName *string `json:"EmployeeDesignationName"`
}

// Retrievetaskinbox retrieves tasks from the inbox based on employee ID
func Retrieveemployeedetails(rows *sql.Rows, LoginName string) ([]Employeedetailsstructure, error) {
	var Employeedetailsapi []Employeedetailsstructure

	for rows.Next() {
		var ES Employeedetailsstructure
		err := rows.Scan(
			&ES.EmployeeId,
			&ES.FirstName,
			&ES.Department,
			&ES.DepartmentCode,
			&ES.DepartmentName,
			&ES.DOJ,
			&ES.EmployeeDesignationId,
			&ES.EmployeeDesignationName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		Employeedetailsapi = append(Employeedetailsapi, ES)
	}

	return Employeedetailsapi, nil
}
