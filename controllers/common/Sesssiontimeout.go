// Package common contains APIs that are commonly used across the application and are grouped together for reusability.
//
// This Api used to mark a user session as inactive upon logout or timeout.
//
// Path:Login Page
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On: 09-07-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 09-07-2025
package commoncontrollers

import (
	"Hrmodule/auth"
	"Hrmodule/dbconfig"
	"Hrmodule/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb"
)

// SessionRequest represents the expected JSON structure for a session timeout request.
// It is used to decode incoming JSON from the /SessionTimeout API endpoint.
type SessionRequest struct {
	SessionID string `json:"session_id"` // SessionID is the identifier of the session to be updated.
}

// APIResponse defines the standard JSON response structure used by API endpoints.
// It contains a status code and a human-readable message.
type APIResponse struct {
	Status  int    `json:"status"`  // HTTP-like status code (e.g., 200 for success, 500 for failure)
	Message string `json:"message"` // Human-readable message describing the result
}

// UpdateSessionLogout updates the `Is_Active` flag to 0 and sets the `Logout_Date`
// to the current timestamp using SQL Server's GETDATE() for the session identified by sessionId.
//
// Parameters:
//   - sessionId: a string representing the session to mark as logged out.
//
// Returns:
//   - error: nil if the update succeeds, or an error describing the failure.
func UpdateSessionLogout(sessionId string) error {

	// Connection string for SQL Server
	connectionString := credentials.GetTestdatabase15()

	// Open connection to SQL Server
	db, err := sql.Open("mssql", connectionString)
	if err != nil {
		return fmt.Errorf("DB open error: %v", err)
	}
	defer db.Close()

	// SQL query to set session inactive and mark logout time
	query := `
		UPDATE Standard..Session_Data 
		SET Is_Active = 0, Logout_Date = GETDATE() 
		WHERE Session_Id = ?`

	// Execute the update query
	_, err = db.Exec(query, sessionId)
	if err != nil {
		return fmt.Errorf("update error: %v", err)
	}

	return nil
}

// SessionTimeoutHandler handles POST requests to the /SessionTimeout endpoint.
// It performs the following:
//   - Authenticates request using IP/token-based validation.
//   - Logs request information.
//   - Parses the incoming JSON body containing a session_id.
//   - Updates the session record in the database by marking it as inactive.
//   - Sends back an encrypted JSON response indicating success or failure.
func SessionTimeoutHandler(w http.ResponseWriter, r *http.Request) {
	// Step 1: Authenticate and validate request origin (IP + token check)
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	// Step 2: Wrap logic with logging
	loggedHandler := auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Step 3: Token validation for user identity
		if err := auth.IsValid_IDFromRequest(r); err != nil {
			http.Error(w, "Invalid Token provided", http.StatusBadRequest)
			return
		}

		// Step 4: Allow only POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Step 5: Decode incoming JSON body
		var req SessionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.SessionID == "" {
			http.Error(w, "Invalid JSON or missing session_id", http.StatusBadRequest)
			return
		}

		// Step 6: Update session logout in database
		err := UpdateSessionLogout(req.SessionID)

		// Build standard API response
		var response APIResponse
		if err != nil {
			response = APIResponse{
				Status:  500,
				Message: "Failed to update session: " + err.Error(),
			}
		} else {
			response = APIResponse{
				Status:  200,
				Message: "Session updated successfully",
			}
		}

		// Step 7: Marshal response to JSON
		responseBytes, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Failed to serialize JSON", http.StatusInternalServerError)
			return
		}

		// Step 8: Encrypt the JSON response using AES-GCM
		encrypted, err := utils.Encrypt(responseBytes)
		if err != nil {
			http.Error(w, "Encryption failed", http.StatusInternalServerError)
			return
		}

		// Step 9: Send encrypted JSON response
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"Data": encrypted,
		})
	}))

	// Step 10: Execute the logged handler
	loggedHandler.ServeHTTP(w, r)
}
