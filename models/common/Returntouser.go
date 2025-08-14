// Package models contains data structures and database access logic for the Return to user page.
//
// --- Creator's Info ---
//
// Creator: Elakiya
//
// Creator: Elakiya
//
// Created On:04-08-2025
//
// Last Modified By: Elakiya
//
// Last Modified Date: 04-08-2025
//
// Path: Return to user
package modelscommon

import (
	"database/sql"
	"fmt"
)

// Query for NOC comments with latest distinct Assign_to_role using CTE and ROW_NUMBER
const NOCCommentQuery = `
WITH RankedComments AS (
    SELECT 
        B.Assign_to_role,
        A.UpdatedBy,
        C.FirstName,
        A.UpdatedOn,
        ROW_NUMBER() OVER (PARTITION BY B.Assign_to_role ORDER BY A.UpdatedOn DESC) AS rn
    FROM HumanResources..NOC_comments A
    JOIN Standard..ActivityMaster_new B 
        ON A.ActivitySeqNo = B.Activityid
    JOIN HumanResources..NOC_Master D 
        ON A.TaskId = D.TaskId
    JOIN HumanResources..UserDetails C 
        ON A.UpdatedBy = C.EmployeeId
    WHERE 
        A.TaskId = ? 
        AND D.ProcessId = ? 
        AND A.ActivitySeqNo < ?
)
SELECT 
    Assign_to_role,
    UpdatedBy,
    FirstName,
    UpdatedOn
FROM RankedComments
WHERE rn = 1
ORDER BY UpdatedOn ASC
`

// NOCComment defines the structure of each task row
type NOCComment struct {
	AssignToRole *string `json:"AssignToRole"`
	UpdatedBy    *string `json:"UpdatedBy"`
	FirstName    *string `json:"FirstName"`
	UpdatedOn    *string `json:"UpdatedOn"`
}

// ScanNOCComments reads rows returned by the query into a slice of NOCComment
func ScanNOCComments(rows *sql.Rows) ([]NOCComment, error) {
	var results []NOCComment

	for rows.Next() {
		var item NOCComment
		if err := rows.Scan(&item.AssignToRole, &item.UpdatedBy, &item.FirstName, &item.UpdatedOn); err != nil {
			return nil, fmt.Errorf("error scanning NOCComment: %v", err)
		}
		results = append(results, item)
	}

	return results, nil
}
