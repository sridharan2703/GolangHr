// --- Creator's Info ---
//
// Creator: Elakiya
//
// Created On: 04-08-2025
//
// Last Modified By: Elakiya
//
// Last Modified Date: 04-08-2025

package commoncontrollers

import (
	"Hrmodule/auth"
	database "Hrmodule/database/common"
	"Hrmodule/utils"
	"encoding/json"
	"net/http"
)

// APIResponseNOCComment defines the structure of the API response returned
// by the GetReturntouser handler.
type APIResponseNOCComment struct {
	Status  int         `json:"Status"`  // HTTP status code
	Message string      `json:"message"` // Message indicating success or failure
	Data    interface{} `json:"Data"`    // Actual payload, includes record count and list
}

func GetReturntouser(w http.ResponseWriter, r *http.Request) {
	// Step 1: Validate API name, token, and IP address
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	// Step 2: Wrap the logic with logging and validation
	loggedHandler := auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Step 3: Validate ID from token or request
		if err := auth.IsValid_IDFromRequest(r); err != nil {
			http.Error(w, "Invalid Token provided", http.StatusBadRequest)
			return
		}

		// Step 4: Retrieve NOC comments from the database
		data, totalCount, err := database.GetNOCCommentData(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Step 5: Build the success response
		response := APIResponseNOCComment{
			Status:  200,
			Message: "Success",
			Data: map[string]interface{}{
				"No Of Records": totalCount,
				"Records":       data,
			},
		}

		// Step 6: Marshal the response to JSON
		jsonResp, err := json.MarshalIndent(response, "", "    ")
		if err != nil {
			http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
			return
		}

		// Step 7: Encrypt the JSON response
		encrypted, err := utils.Encrypt(jsonResp)
		if err != nil {
			http.Error(w, "Encryption failed", http.StatusInternalServerError)
			return
		}

		// Step 8: Send the encrypted response
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"Data": encrypted,
		})
	}))

	// Step 9: Serve the HTTP request using the wrapped handler
	loggedHandler.ServeHTTP(w, r)
}
