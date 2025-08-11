// Package models contains data structure for TaskSummary_count and database TaskSummary_count logic for the NOCmaster page.
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
//
// Path: Task Summary Page
package modelshumanresources

import (
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb" // MSSQL driver for database/sql
)

const MyQueryTaskSummary_count = `TaskSummary_count  @EmpId=?`

// TaskSummary_countstructure definNOC_M the structure of each task in the inbox
type TaskSummary_countstructure struct {
	InitiatedOngoing      *int `json:"InitiatedOngoing"`
	InitiatedCompleted    *int `json:"InitiatedCompleted"`
	InitiatedDeleted      *int `json:"InitiatedDeleted"`
	ParticipatedOngoing   *int `json:"ParticipatedOngoing"`
	ParticipatedCompleted *int `json:"ParticipatedCompleted"`
	ParticipatedDeleted   *int `json:"ParticipatedDeleted"`
}

// TaskSummary_count retriev tasksummary count  from the inbox based on employee ID
func Retrievetasksummarycount(rows *sql.Rows, EmpId string) ([]TaskSummary_countstructure, error) {
	var TaskSummary_countapi []TaskSummary_countstructure

	for rows.Next() {
		var tsc TaskSummary_countstructure
		err := rows.Scan(
			&tsc.InitiatedOngoing,
			&tsc.InitiatedCompleted,
			&tsc.InitiatedDeleted,
			&tsc.ParticipatedOngoing,
			&tsc.ParticipatedCompleted,
			&tsc.ParticipatedDeleted,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		TaskSummary_countapi = append(TaskSummary_countapi, tsc)
	}

	return TaskSummary_countapi, nil
}
