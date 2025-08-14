// Package models contains data structures and database access logic for the DefaultRoleName page.
//
// --- Creator's Info ---
//
// Creator: Sivabala
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

const MyQueryRoles = `SELECT * FROM Roles WHERE isactive = '1' and role_names =?;`

// DefaultRoleNamestructure defines the structure of DefaultRoleName
type Rolesstructure struct {
	Id                      *int    `json:"Id"`
	Role_Names              *string `json:"Role_Names"`
	Rights                  *string `json:"Rights"`
	Menu                    *string `json:"Menu"`
	Module                  *string `json:"Module"`
	IsEnabled               *string `json:"IsEnabled"`
	Path                    *string `json:"Path"`
	Element                 *string `json:"Element"`
	Icons                   *string `json:"Icons"`
	IsActive                *string `json:"IsActive"`
	Created_By              *string `json:"Created_By"`
	Created_On              *string `json:"Created_On"`
	Last_UpdatedBy          *string `json:"Last_UpdatedBy"`
	Last_UpdatedOn          *string `json:"Last_UpdatedOn"`
}

// RetrieveDefaultRoleName retrieves tasks from the DefaultRoleName based on RoleName
func RetrieveRoles(rows *sql.Rows, RoleName string) ([]Rolesstructure, error) {
	var Rolesapi []Rolesstructure

	for rows.Next() {
		var Role Rolesstructure
		err := rows.Scan(
			&Role.Id,
			&Role.Role_Names,
			&Role.Rights,
			&Role.Menu,
			&Role.Module,
			&Role.IsEnabled,
			&Role.Path,
			&Role.Element,
			&Role.Icons,
			&Role.IsActive,
			&Role.Created_By,
			&Role.Created_On,
			&Role.Last_UpdatedBy,
			&Role.Last_UpdatedOn,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		Rolesapi = append(Rolesapi, Role)
	}

	return Rolesapi, nil
}
