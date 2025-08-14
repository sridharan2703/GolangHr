// Package models contains data structures and database access logic for the DefaultRoleName page.
//
// --- Creator's Info ---
//
// Creator: Sivabala
//
// Created On:30-07-2025
//
// Last Modified By: Sivabala
//
// Last Modified Date: 30-07-2025
//
// Path: Task inbox Page
package modelscommon

import (
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb" // MSSQL driver for database/sql
)

const MyQueryDefaultRoleName = `
SELECT RM.RoleName as DefaultRoleName,ADS.Department,ADS.EmployeeId FROM        
standard..UserMaster UM WITH(NOLOCK) INNER JOIN        
standard..UserDetails ADS WITH(NOLOCK) ON ADS.LoginName=UM.UserName INNER JOIN        
standard..OrgUnitUserMapping OUM WITH(NOLOCK) ON       
OUM.UserId=UM.UserId INNER JOIN        
standard..OrgUnitRoleMapping ORM WITH(NOLOCK) ON               
ORM.RoleMapId=OUM.RoleMapId INNER JOIN      
HumanResources..RoleMaster RM WITH(NOLOCK) ON RM.RoleId=ORM.RoleId WHERE        
OUM.IsActive=1 and ADS.LoginName=?`

// DefaultRoleNamestructure defines the structure of DefaultRoleName
type DefaultRoleNamestructure struct {
	DefaultRoleName *string `json:"DefaultRoleName"`
	Department      *string `json:"Department"`
	EmployeeId      *string `json:"EmployeeId"`
}

// RetrieveDefaultRoleName retrieves tasks from the DefaultRoleName based on LoginName
func RetrieveDefaultRoleName(rows *sql.Rows, LoginName string) ([]DefaultRoleNamestructure, error) {
	var DefaultRoleNameapi []DefaultRoleNamestructure

	for rows.Next() {
		var DRN DefaultRoleNamestructure
		err := rows.Scan(
			&DRN.DefaultRoleName,
			&DRN.Department,
			&DRN.EmployeeId,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		DefaultRoleNameapi = append(DefaultRoleNameapi, DRN)
	}

	return DefaultRoleNameapi, nil
}
