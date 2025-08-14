// Package models contains data structure for NOC_Master combined with masterhistory istask rejected remarks logic for this page.
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

const MyQueryNOCmasterremarks = `select *from humanresources..noc_masterhistory where Taskid=?;
`

// NOCmasterremarksstructure defines the structure of each task in the inbox
type NOCmasterremarksstructure struct {
	Id                     *int    `json:"Id"`
	Taskid                 *string `json:"Taskid"`
	EmployeeId             *string `json:"EmployeeId"`
	NOCFor                 *string `json:"NOCFor"`
	NOCPassportFor         *string `json:"NOCPassportFor"`
	PassportNumber         *string `json:"PassportNumber"`
	PassportDateofIssue    *string `json:"PassportDateofIssue"`
	PassportValidTill      *string `json:"PassportValidTill"`
	TentativeTravelFrom    *string `json:"TentativeTravelFrom"`
	TentativeTravelTo      *string `json:"TentativeTravelTo"`
	CountryOfferingVisa    *string `json:"CountryOfferingVisa"`
	PlaceOfVisit           *string `json:"PlaceOfVisit"`
	ReasonforVisiting      *string `json:"ReasonforVisiting"`
	NOCReuiredForDependant *string `json:"NOCReuiredForDependant"`
	UpdatedBy              *string `json:"UpdatedBy"`
	UpdatedOn              *string `json:"UpdatedOn"` // or use *time.Time
	AssignTo               *string `json:"AssignTo"`
	Documentdetails        *string `json:"Documentdetails"`
	TaskStatusId           *int    `json:"TaskStatusId"`
	ActivitySeqNo          *int    `json:"ActivitySeqNo"`
	IsTaskReturn           *string `json:"IsTaskReturn"`
	Remarks                *string `json:"Remarks"`
	Istaskapproved         *string `json:"IsTaskApproved"`
	CoverPageNo            *string `json:"CoverPageNo"`
	AssignedRole           *string `json:"AssignedRole"`
	ProcessId              *int    `json:"ProcessId"`
	InitiatedBy            *string `json:"InitiatedBy"`
	InitiatedOn            *string `json:"InitiatedOn"`
}

// RetrieveNOCmasterremarks retrieves tasks from the inbox based on task ID
func RetrieveNOCmasterremarks(rows *sql.Rows, Taskid string) ([]NOCmasterremarksstructure, error) {
	var NOCmasterremarksapi []NOCmasterremarksstructure

	for rows.Next() {
		var NOC_M NOCmasterremarksstructure
		err := rows.Scan(
			&NOC_M.Id,
			&NOC_M.Taskid,
			&NOC_M.EmployeeId,
			&NOC_M.NOCFor,
			&NOC_M.NOCPassportFor,
			&NOC_M.PassportNumber,
			&NOC_M.PassportDateofIssue,
			&NOC_M.PassportValidTill,
			&NOC_M.TentativeTravelFrom,
			&NOC_M.TentativeTravelTo,
			&NOC_M.CountryOfferingVisa,
			&NOC_M.PlaceOfVisit,
			&NOC_M.ReasonforVisiting,
			&NOC_M.NOCReuiredForDependant,
			&NOC_M.UpdatedBy,
			&NOC_M.UpdatedOn,
			&NOC_M.AssignTo,
			&NOC_M.Documentdetails,
			&NOC_M.TaskStatusId,
			&NOC_M.ActivitySeqNo,
			&NOC_M.IsTaskReturn,
			&NOC_M.Remarks,
			&NOC_M.Istaskapproved,
			&NOC_M.CoverPageNo,
			&NOC_M.AssignedRole,
			&NOC_M.ProcessId,
			&NOC_M.InitiatedBy,
			&NOC_M.InitiatedOn,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		NOCmasterremarksapi = append(NOCmasterremarksapi, NOC_M)
	}

	return NOCmasterremarksapi, nil
}
