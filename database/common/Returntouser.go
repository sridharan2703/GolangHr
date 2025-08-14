// Package databasecommon handles database connections and queries
// related to NOC comment data in the Human Resources module.
//
// --- Creator's Info ---
//
// Creator: Elakiya
//
// Created On: 04-08-2025
//
// Last Modified By: Elakiya
//
// Last Modified Date: 04-08-2025

package databasecommon

import (
	credentials "Hrmodule/dbconfig"
	models "Hrmodule/models/common"
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb" // SQL Server driver
)

func GetNOCCommentData(w http.ResponseWriter, r *http.Request) ([]models.NOCComment, int, error) {
	// Establish database connection
	connStr := credentials.GetTestdatabasetwo15()
	db, err := sql.Open("mssql", connStr)
	if err != nil {
		return nil, 0, fmt.Errorf("DB connection error: %v", err)
	}
	defer db.Close()

	// Extract query parameters (case-insensitive fallback handling)
	taskId := r.URL.Query().Get("TaskId")
	if taskId == "" {
		taskId = r.URL.Query().Get("taskId")
	}

	processId := r.URL.Query().Get("ProcessId")
	if processId == "" {
		processId = r.URL.Query().Get("processId")
	}

	activitySeqNo := r.URL.Query().Get("ActivitySeqNo")
	if activitySeqNo == "" {
		activitySeqNo = r.URL.Query().Get("activitySeqNo")
	}

	// Validate required parameters
	if taskId == "" || processId == "" || activitySeqNo == "" {
		return nil, 0, fmt.Errorf("Missing required parameters")
	}

	// Run the NOC comment query
	rows, err := db.Query(models.NOCCommentQuery, taskId, processId, activitySeqNo)
	if err != nil {
		return nil, 0, fmt.Errorf("Query execution error: %v", err)
	}
	defer rows.Close()

	// Scan the query results into model
	data, err := models.ScanNOCComments(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("Data scan error: %v", err)
	}

	return data, len(data), nil
}
