// Package models contains data structure for TaskSummary and database TaskSummary logic for the Task Summary page.
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
//
// Path: Task Summary Page
package modelshumanresources

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/denisenkom/go-mssqldb" // MSSQL driver for database/sql
)

const MyQueryTaskSummary = `Task_Summary @EmpId=?, @Type=?`

// TaskSummarystructure defines the structure of each task in the summary
type TaskSummarystructure struct {
	TaskId            *string    `json:"TaskId"`
	CoverPageNo       *string    `json:"CoverPageNo"`
	Type              *string    `json:"Type"`
	TaskStatusId      *int       `json:"TaskStatusId"`
	EmployeeId        *string    `json:"EmployeeId"`
	UpdatedOn         *time.Time `json:"UpdatedOn"`
	DateDiff          *int       `json:"DateDiff"`
	AssignTo          *string    `json:"AssignTo"`
	Inbox             *string    `json:"Inbox"`
	ProcessName       *string    `json:"ProcessName"`
	CurrentActivity   *string    `json:"CurrentActivity"`
	ProcessKeyword    *string    `json:"ProcessKeyword"`
	StatusDescription *string    `json:"StatusDescription"`
	LoginName         *string    `json:"LoginName"`
	ProcessId         *int       `json:"ProcessId"`
	ActivitySeqNo     *int       `json:"ActivitySeqNo"`
}

// Retrievetasksummary retrieves task summary from the database based on employee ID and type
func Retrievetasksummary(rows *sql.Rows, EmpId string, Type string) ([]TaskSummarystructure, error) {
	var TaskSummaryapi []TaskSummarystructure

	rowCount := 0
	for rows.Next() {
		var ts TaskSummarystructure
		err := rows.Scan(
			&ts.TaskId,
			&ts.CoverPageNo,
			&ts.Type,
			&ts.TaskStatusId,
			&ts.EmployeeId,
			&ts.UpdatedOn,
			&ts.DateDiff,
			&ts.AssignTo,
			&ts.Inbox,
			&ts.ProcessName,
			&ts.CurrentActivity,
			&ts.ProcessKeyword,
			&ts.StatusDescription,
			&ts.LoginName,
			&ts.ProcessId,
			&ts.ActivitySeqNo,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row %d: %v", rowCount+1, err)
		}
		TaskSummaryapi = append(TaskSummaryapi, ts)
		rowCount++
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	log.Printf("Successfully scanned %d rows for EmpId: %s, Type: %s", rowCount, EmpId, Type)
	return TaskSummaryapi, nil
}
