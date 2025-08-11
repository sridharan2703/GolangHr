// Package controllers for the Activityseq handles HTTP request routing, authentication,
// and API response formatting for the application.
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On:15-07-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 15-07-2025
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

// Activityseq struct for request body parameters
type Activityseq struct {
	IsTaskApproved int    `json:"isTaskApproved"`
	IsTaskReturn   int    `json:"isTaskReturn"`
	Remarks        string `json:"remarks"`
}

// Activitysequpdate handles calling the stored procedure UpdateNextActivitySeqNo
func Activitysequpdate(w http.ResponseWriter, r *http.Request) {
	// Step 1: Authenticate and validate token/IP
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	// Step 2: Log request and proceed
	loggedHandler := auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Step 3: Validate token format
		if err := auth.IsValid_IDFromRequest(r); err != nil {
			http.Error(w, "Invalid Token provided", http.StatusBadRequest)
			return
		}

		// Step 4: Extract TaskId from query string
		taskIdStr := r.URL.Query().Get("TaskId")
		if taskIdStr == "" {
			http.Error(w, "TaskId is required in query parameter", http.StatusBadRequest)
			return
		}

		// Step 5: Decode JSON body
		var req Activityseq
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid JSON body: "+err.Error(), http.StatusBadRequest)
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

		// Step 6: Call stored procedure using global DB
		_, err = db.Exec(`
			EXEC HumanResources..UpdateNextActivitySeqNo 
			@TaskId = ?, 
			@IsTaskApproved = ?, 
			@IsTaskReturn = ?, 
			@Remarks = ?`,
			taskIdStr,
			req.IsTaskApproved,
			req.IsTaskReturn,
			req.Remarks,
		)

		if err != nil {
			log.Printf("Error executing stored procedure: %v", err)
			http.Error(w, "Failed to execute stored procedure: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Step 7: Send encrypted response
		response := map[string]interface{}{
			"Status":  200,
			"Message": "Stored procedure executed successfully",
			"TaskId":  taskIdStr,
			"Data":    map[string]interface{}{},
		}

		plainJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "JSON marshalling failed", http.StatusInternalServerError)
			return
		}

		encrypted, err := utils.Encrypt(plainJSON)
		if err != nil {
			http.Error(w, "Encryption failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"Data": encrypted,
		})
	}))

	loggedHandler.ServeHTTP(w, r)
}
