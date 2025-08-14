// Package common contains APIs that are commonly used across the application and are grouped together for reusability.
//
// This API returns the route path mapped to the authenticated user.
//
// Path:Login Page
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On: 2025-07-04
//
// Last Modified By: Sridharan
//
// Last Modified Date: 2025-07-04
package commoncontrollers

import (
	"encoding/json"
	"net/http"

	"Hrmodule/auth"
	"Hrmodule/utils"
)

// Menu represents a single menu item.
type Menu struct {
	Path string `json:"Path"`
	Icon string `json:"Icon"`
}

// Getmenu is an HTTP handler that returns hardcoded
// menu data as a JSON array, encrypted using AES-GCM.
func Getmenu(w http.ResponseWriter, r *http.Request) {
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

		menus := []Menu{
			{
				Path: `import Dashboard from "../../pages/dashboard/Dashboard"`,
				Icon: "fa-solid fa-list-check",
			},
			{
				Path: `import Profile from "../../pages/profile/Profile"`,
				Icon: "fa-solid fa-circle-check",
			},
		}

		// Step 4: Marshal JSON response
		jsonResponse, err := json.MarshalIndent(menus, "", "  ")
		if err != nil {
			http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
			return
		}

		// Step 5: Encrypt JSON response using AES-GCM
		encrypted, err := utils.Encrypt(jsonResponse)
		if err != nil {
			http.Error(w, "Encryption failed", http.StatusInternalServerError)
			return
		}

		// Step 6: Send encrypted response
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"Data": encrypted,
		})
	}))

	// Step 7: Serve the request
	loggedHandler.ServeHTTP(w, r)
}
