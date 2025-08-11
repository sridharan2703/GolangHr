// Package models contains data structure for NOC_Master and database NOC_Master logic for the NOCmaster page.
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
package modelshumanresources

import (
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb" // MSSQL driver for database/sql
)

const MyQueryNOCmaster = `SELECT 
    n.*,
   -- u.EmployeeId,
    u.FirstName,
    u.Department,
    LEFT(u.Department, CHARINDEX(',', u.Department + ',') - 1) AS DepartmentCode,
    d.DepartmentName,
    a.DOJ,
    a.EmployeeDesignationId,
    desig.EmployeeDesignationName
FROM HumanResources..noc_master n
LEFT JOIN HumanResources..UserDetails u
    ON n.EmployeeId = u.EmployeeId
LEFT JOIN HumanResources..DepartmentMaster d 
    ON LEFT(u.Department, CHARINDEX(',', u.Department + ',') - 1) = d.DepartmentCode
LEFT JOIN HumanResources..AppointmentDetails a 
    ON u.EmployeeId = a.EmployeeId
LEFT JOIN HumanResources..DesignationMaster desig 
    ON a.EmployeeDesignationId = desig.EmployeeDesignationId
WHERE n.TaskId = ?;
`

// NOCmasterstructure definNOC_M the structure of each task in the inbox
type NOCmasterstructure struct {
	Id                      *int    `json:"Id"`
	Taskid                  *string `json:"Taskid"`
	EmployeeId              *string `json:"EmployeeId"`
	NOCFor                  *string `json:"NOCFor"`
	NOCPassportFor          *string `json:"NOCPassportFor"`
	PassportNumber          *string `json:"PassportNumber"`
	PassportDateofIssue     *string `json:"PassportDateofIssue"`
	PassportValidTill       *string `json:"PassportValidTill"`
	TentativeTravelFrom     *string `json:"TentativeTravelFrom"`
	TentativeTravelTo       *string `json:"TentativeTravelTo"`
	CountryOfferingVisa     *string `json:"CountryOfferingVisa"`
	PlaceOfVisit            *string `json:"PlaceOfVisit"`
	ReasonforVisiting       *string `json:"ReasonforVisiting"`
	NOCReuiredForDependant  *string `json:"NOCReuiredForDependant"`
	UpdatedBy               *string `json:"UpdatedBy"`
	UpdatedOn               *string `json:"UpdatedOn"` // optionally use *time.Time if parsed
	AssignTo                *string `json:"AssignTo"`
	Documentdetails         *string `json:"Documentdetails"`
	TaskStatusId            *int    `json:"TaskStatusId"`
	ActivitySeqNo           *int    `json:"ActivitySeqNo"`
	IsTaskReturn            *string `json:"IsTaskReturn"`
	Istaskapproved          *string `json:"Istaskapproved"`
	CoverPageNo             *string `json:"CoverPageNo"`
	ProcessId               *int    `json:"ProcessId"`
	AssignedRole            *string `json:"AssignedRole"`
	InitiatedBy             *string `json:"InitiatedBy"`
	InitiatedOn             *string `json:"InitiatedOn"`
	FirstName               *string `json:"FirstName"`
	Department              *string `json:"Department"`
	DepartmentCode          *string `json:"DepartmentCode"`
	DepartmentName          *string `json:"DepartmentName"`
	DOJ                     *string `json:"DOJ"`
	EmployeeDesignationId   *int    `json:"EmployeeDesignationId"`
	EmployeeDesignationName *string `json:"EmployeeDesignationName"`
}

// Retrievetaskinbox retrievNOC_M tasks from the inbox based on employee ID
func RetrieveNOCmaster(rows *sql.Rows, Taskid string) ([]NOCmasterstructure, error) {
	var NOCmasterapi []NOCmasterstructure

	for rows.Next() {
		var NOC_M NOCmasterstructure
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
			&NOC_M.Istaskapproved,
			&NOC_M.CoverPageNo,
			&NOC_M.ProcessId,
			&NOC_M.AssignedRole,
			&NOC_M.InitiatedBy,
			&NOC_M.InitiatedOn,
			&NOC_M.FirstName,
			&NOC_M.Department,
			&NOC_M.DepartmentCode,
			&NOC_M.DepartmentName,
			&NOC_M.DOJ,
			&NOC_M.EmployeeDesignationId,
			&NOC_M.EmployeeDesignationName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		NOCmasterapi = append(NOCmasterapi, NOC_M)
	}

	return NOCmasterapi, nil
}
