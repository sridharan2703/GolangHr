// Package controllers handles HTTP request routing, authentication,
// and API response formatting for the application.
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On:15-07-2025
//
// Last Modified By: Sridharan11111111111
//
// Last Modified Date: 15-07-202511111111111111
package commoncontrollers

import (
	"Hrmodule/auth"
	database "Hrmodule/database/common"
	"Hrmodule/utils"
	"encoding/json"
	"net/http"
)

// APIResponseforemployeedetails defines the standard structure of the API response.
type APIResponseforemployeedetails struct {
	Status  int         `json:"Status"`
	Message string      `json:"message"`
	Data    interface{} `json:"Data"`
}

// GetEmployeedetails handles the HTTP GET request to fetch Employeedetails data for a Employees.
//
// Flow:
//  1. Authenticates the request using token/IP validation
//  2. Logs the incoming request
//  3. Fetches Employeedetails data from the database based on Employeeid
//  4. Constructs a structured JSON response
//  5. Encrypts the response using AES-GCM
//  6. Sends the encrypted string in the "encrypted" field of the response
//
// Response:
//   - JSON object with one field "encrypted" containing the AES-encrypted payload
func GetEmployeedetails(w http.ResponseWriter, r *http.Request) {
	// Step 1: Authenticate and validate request origin (token/IP)
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	// Step 2: Wrap in logging middleware
	loggedHandler := auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Step 3: Validate token and extract ID from request
		if err := auth.IsValid_IDFromRequest(r); err != nil {
			http.Error(w, "Invalid Token provided", http.StatusBadRequest)
			return
		}

		// Step 4: Retrieve Employeedetails data from DB
		EmployeedetailsData, totalCount, err := database.Employeedetailsdatabase(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Step 5: Structure API response
		var response APIResponseforemployeedetails
		{
			// Case: Successful data retrieval
			response = APIResponseforemployeedetails{
				Status:  200,
				Message: "Success",
				Data: map[string]interface{}{
					"No Of Records": totalCount,
					"Records":       EmployeedetailsData,
				},
			}
		}

		// Step 6: Convert response to JSON
		jsonResponse, err := json.MarshalIndent(response, "", "    ")
		if err != nil {
			http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
			return
		}

		// Step 7: Encrypt JSON response using AES-GCM
		encrypted, err := utils.Encrypt(jsonResponse)
		if err != nil {
			http.Error(w, "Encryption failed", http.StatusInternalServerError)
			return
		}

		// Step 8: Send encrypted response
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"Data": encrypted,
		})
	}))

	// Step 9: Execute the handler with logging
	loggedHandler.ServeHTTP(w, r)
}
