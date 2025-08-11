// Package controllers for the NOCTaskreturn handles HTTP request routing, authentication,
// and API response formatting for the application.
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
package controllershumanresources

import (
	"Hrmodule/auth"
	"Hrmodule/dbconfig"
	"Hrmodule/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb"
)

// Handler that calls NOC_TaskReturnUpdate with TaskId
func NOCTaskreturn(w http.ResponseWriter, r *http.Request) {
	// Step 1: Authenticate request (IP, Token etc.)
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	// Step 2: Wrap in logger and token validator
	loggedHandler := auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Step 3: Validate token format
		if err := auth.IsValid_IDFromRequest(r); err != nil {
			http.Error(w, "Invalid token provided", http.StatusBadRequest)
			return
		}

		// Step 4: Extract TaskId from query parameters
		taskId := r.URL.Query().Get("TaskId")
		if taskId == "" {
			http.Error(w, "Missing TaskId parameter", http.StatusBadRequest)
			return
		}
		// Step 6: Database connection and operation
		// Connection string for SQL Server
		connectionString := credentials.GetTestdatabasetwo15()

		// Open connection to SQL Server
		db, err := sql.Open("mssql", connectionString)
		if err != nil {
			log.Printf("DB open error: %v", err)
			http.Error(w, fmt.Sprintf("DB open error: %v", err), http.StatusInternalServerError)
			return
		}
		defer db.Close()
		// Step 5: Call stored procedure
		_, err = db.Exec(`EXEC NOC_TaskReturnUpdate @TaskId = ?`, taskId)
		if err != nil {
			http.Error(w, "Failed to execute stored procedure: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Step 6: Prepare encrypted response
		response := map[string]interface{}{
			"Status":  200,
			"Message": "NOC_TaskReturnUpdate executed successfully",
			"TaskId":  taskId,
		}

		plainJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}

		encrypted, err := utils.Encrypt(plainJSON)
		if err != nil {
			http.Error(w, "Encryption failed", http.StatusInternalServerError)
			return
		}

		// Step 7: Send encrypted response
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"Data": encrypted,
		})
	}))

	loggedHandler.ServeHTTP(w, r)
}
