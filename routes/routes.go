// Package routes defines CORS settings and registers HTTP routes for the API server.
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On:07-07-2025
//
// Last Modified By: Sridharan
//
// Last Modified Date: 09-07-2025
package routes

import (
	controllershumanresources "Hrmodule/controllers/HumanResources"
	commoncontrollers "Hrmodule/controllers/common"

	"fmt"
	"net/http"

	"github.com/rs/cors"
)

// Registerroutes sets up the HTTPS server with CORS support,
func Registerroutes() {
	// Create a new ServeMux router
	router := http.NewServeMux()

	// Register your API routes
	router.Handle("/HRldap", http.HandlerFunc(commoncontrollers.HandleLDAPAuth))
	router.Handle("/SessionTimeout", http.HandlerFunc(commoncontrollers.SessionTimeoutHandler))
	router.Handle("/Menu", http.HandlerFunc(commoncontrollers.Getmenu))
	router.Handle("/Taskinbox", http.HandlerFunc(commoncontrollers.GetTaskinbox))
	router.Handle("/Employeedetails", http.HandlerFunc(commoncontrollers.GetEmployeedetails))
	router.Handle("/NOCupdateinsert", http.HandlerFunc(controllershumanresources.NOCHandler))
	router.Handle("/NOCmaster", http.HandlerFunc(controllershumanresources.GetNOCmaster))
	router.Handle("/Activitysequpdate", http.HandlerFunc(controllershumanresources.Activitysequpdate))
	router.Handle("/NOCmasterremarks", http.HandlerFunc(controllershumanresources.GetNOCmasterremarks))
	router.Handle("/Tasksummarycount", http.HandlerFunc(controllershumanresources.Gettasksummarycount))
	router.Handle("/Tasksummary", http.HandlerFunc(controllershumanresources.Gettasksummary))
	router.Handle("/NOCTaskreturn", http.HandlerFunc(controllershumanresources.NOCTaskreturn))
	router.Handle("/NOCcommentsremarks", http.HandlerFunc(controllershumanresources.GetNOCremarks))
	router.Handle("/DefaultRoleName", http.HandlerFunc(commoncontrollers.DefaultRoleName))
	router.Handle("/UserRoles", http.HandlerFunc(commoncontrollers.GetRoles))
	router.Handle("/Returntouser", http.HandlerFunc(commoncontrollers.GetReturntouser))
	router.HandleFunc("/DynamicActivitySequence", http.HandlerFunc(commoncontrollers.DynamicActivitySequence))
	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Use specific origin(s) in production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Apply CORS middleware to the router
	handler := c.Handler(router)

	fmt.Println("Server starting on port 5000")

	// TLS certificate and key
	certFile := "certificate.pem"
	keyFile := "key.pem"

	// Start the HTTPS server with CORS-enabled handler
	err := http.ListenAndServeTLS(":5000", certFile, keyFile, handler)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
