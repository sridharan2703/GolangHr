// Package models contains data structures and database access logic for the task inbox page.
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Creator: Sridharan
//
// Created On:15-07-2025
//
// Last Modified By: Elakiya
//
// Last Modified Date: 04-08-2025
//
// Path: Task inbox Page
package modelscommon

import (
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb" // MSSQL driver for database/sql
)

const MyQuerytaskinbox = `GETInboxTasks_role  @EmpId=?, @AssignedRole=?`

// Taskinboxstructure defines the structure of each task in the inbox
type Taskinboxstructure struct {
	TaskId         *string `json:"TaskId"`
	EmployeeId     *string `json:"EmployeeId"`
	UpdatedOn      *string `json:"UpdatedOn"`
	UpdatedBy      *string `json:"UpdatedBy"`
	ActivitySeqNo  *int    `json:"ActivitySeqNo"`
	Remarks        *string `json:"Remarks"`
	ProcessName    *string `json:"ProcessName"`
	ProcessKeyword *string `json:"ProcessKeyword"`
	Path           *string `json:"Path"`
	Component      *string `json:"Component"`
	CoverPageNo    *string `json:"CoverPageNo"`
	ProcessId      *string `json:"ProcessId"`
	//SendBackToMe   string  `json:"SendBackToMe"`
	//IsActive       string  `json:"IsActive"`
}

// Retrievetaskinbox retrieves tasks from the inbox based on employee ID
func Retrievetaskinbox(rows *sql.Rows, EmpId interface{}, AssignedRole interface{}) ([]Taskinboxstructure, error) {
	var Taskinboxapi []Taskinboxstructure

	for rows.Next() {
		var TI Taskinboxstructure
		err := rows.Scan(
			&TI.TaskId,
			&TI.EmployeeId,
			&TI.UpdatedOn,
			&TI.UpdatedBy,
			&TI.ActivitySeqNo,
			&TI.Remarks,
			&TI.ProcessName,
			&TI.ProcessKeyword,
			&TI.Path,
			&TI.Component,
			&TI.CoverPageNo,
			&TI.ProcessId,
			//&TI.SendBackToMe,
			//&TI.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		Taskinboxapi = append(Taskinboxapi, TI)
	}

	return Taskinboxapi, nil
}
