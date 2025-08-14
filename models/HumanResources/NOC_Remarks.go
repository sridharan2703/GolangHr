// Package models contains data structure for NOC_Remarks combined with comments table logic for this page.
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On:21-07-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 21-07-2025
//
// Path: Task inbox Page
package modelshumanresources

import (
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb" // MSSQL driver for database/sql
)

const MyQueryNOCremarks = `SELECT 
    UD.DefaultRoleName,
    UD.DisplayName,
    NC.Remarks,
    NC.UpdatedOn,
    NC.UpdatedBy
FROM 
    standard..UserDetails UD
JOIN 
    HumanResources..NOC_Comments NC
    ON UD.EmployeeID = NC.UpdatedBy
WHERE  
    NC.TaskID = ?
ORDER BY 
    NC.UpdatedOn DESC;

 `

// NOCmasterremarksstructure defines the structure of each task in the inbox
type NOCremarksstructure struct {
	DefaultRoleName *string `json:"DefaultRoleName"`
	DisplayName     *string `json:"DisplayName"`
	Remarks         *string `json:"Remarks"`
	UpdatedOn       *string `json:"UpdatedOn"`
	UpdatedBy       *string `json:"UpdatedBy"`
}

// RetrieveNOCmasterremarks retrieves tasks from the inbox based on task ID
func RetrieveNOCremarks(rows *sql.Rows, Taskid string) ([]NOCremarksstructure, error) {
	var NOCremarksapi []NOCremarksstructure

	for rows.Next() {
		var NOC_re NOCremarksstructure
		err := rows.Scan(
			&NOC_re.DefaultRoleName,
			&NOC_re.DisplayName,
			&NOC_re.Remarks,
			&NOC_re.UpdatedOn,
			&NOC_re.UpdatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		NOCremarksapi = append(NOCremarksapi, NOC_re)
	}

	return NOCremarksapi, nil
}
