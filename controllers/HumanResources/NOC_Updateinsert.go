// Package controllershumanresources provides HTTP handlers for
// NOC-related operations in the Human Resources module.
// --- Creator's Info ---
// Creator: Sridharan

// Created On: 15-07-2025

// Last Modified By: Sridharan

// Last Modified Date: 15-07-2025
package controllershumanresources

import (
	"Hrmodule/auth"
	credentials "Hrmodule/dbconfig"
	"Hrmodule/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb"
)

// NOCRequest represents the structure of the input JSON body
// for the NOC_Updateinsert stored procedure call.
type NOCRequest struct {
	TaskId                 *string `json:"TaskId,omitempty"` // Not used anymore, taken from query param
	EmployeeId             string  `json:"EmployeeId"`
	NOCFor                 string  `json:"NOCFor"`
	NOCPassportFor         string  `json:"NOCPassportFor"`
	PassportNumber         string  `json:"PassportNumber"`
	PassportDateofIssue    string  `json:"PassportDateofIssue"`
	PassportValidTill      string  `json:"PassportValidTill"`
	TentativeTravelFrom    string  `json:"TentativeTravelFrom"`
	TentativeTravelTo      string  `json:"TentativeTravelTo"`
	CountryOfferingVisa    string  `json:"CountryOfferingVisa"`
	PlaceOfVisit           string  `json:"PlaceOfVisit"`
	ReasonforVisiting      string  `json:"ReasonforVisiting"`
	NOCReuiredForDependant string  `json:"NOCReuiredForDependant"`
	Documentdetails        int     `json:"Documentdetails"`
	Remarks                string  `json:"Remarks"`
	ProcessId              int     `json:"ProcessId"`
	SequenceId             int     `json:"SequenceId"`
	IsTaskApproved         int     `json:"IsTaskApproved"`
	IsTaskReturn           int     `json:"IsTaskReturn"`
	Role                   string  `json:"Role"`
	UpdatedBy              string  `json:"UpdatedBy"`
	RequestedUser          string  `json:"RequestedUser"`
	ReturnToRole           string  `json:"ReturnToRole"`
	ReturnToUser           string  `json:"ReturnToUser"`
	SendBackToMe           string  `json:"SendBackToMe"`
	SendBackToUser         string  `json:"SendBackToUser"`
}

// NOCHandler handles HTTP requests for inserting or updating
// NOC records by calling the NOC_Updateinsert stored procedure.
//
// It performs token/IP validation, parses JSON body,
// executes the procedure, and sends back an encrypted response.
func NOCHandler(w http.ResponseWriter, r *http.Request) {
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
		var taskIdParam interface{}
		if taskIdStr == "" {
			taskIdParam = nil // Will pass NULL to SQL Server
		} else {
			taskIdParam = taskIdStr
		}

		// Step 5: Decode JSON body
		var req NOCRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid JSON body: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Convert optional fields to nil if empty or zero
		var requestedUserParam interface{}
		if req.RequestedUser == "" {
			requestedUserParam = nil
		} else {
			requestedUserParam = req.RequestedUser
		}

		var returnToUserParam interface{}
		if req.ReturnToUser == "" {
			returnToUserParam = nil
		} else {
			returnToUserParam = req.ReturnToUser
		}

		var sendBackToMeParam interface{}
		if req.SendBackToMe == "" {
			sendBackToMeParam = nil
		} else {
			sendBackToMeParam = req.SendBackToMe
		}

		var sendBackToUserParam interface{}
		if req.SendBackToUser == "" {
			sendBackToUserParam = nil
		} else {
			sendBackToUserParam = req.SendBackToUser
		}

		var returnToRoleParam interface{}
		if req.ReturnToRole == "" {
			returnToRoleParam = nil
		} else {
			returnToRoleParam = req.ReturnToRole
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

		_, err = db.Exec(`
			EXEC HumanResources..NOC_Updateinsert_New
				@TaskId = ?,
				@EmployeeId = ?,
				@NOCFor = ?,
				@NOCPassportFor = ?,
				@PassportNumber = ?,
				@PassportDateofIssue = ?,
				@PassportValidTill = ?,
				@TentativeTravelFrom = ?,
				@TentativeTravelTo = ?,
				@CountryOfferingVisa = ?,
				@PlaceOfVisit = ?,
				@ReasonforVisiting = ?,
				@NOCReuiredForDependant = ?,
				@Documentdetails = ?,
				@Remarks =?,
				@ProcessId=?,
                @SequenceId=?,
				@IsTaskApproved=?,
				@IsTaskReturn=?,
				@Role=?,
				@UpdatedBy =?,
		        @RequestedUser = ?,
		        @ReturnToRole =?,
		        @ReturnToUser =?,
		        @SendBackToMe = ?,
	         	@SendBackToUser = ?`,
			taskIdParam,
			req.EmployeeId,
			req.NOCFor,
			req.NOCPassportFor,
			req.PassportNumber,
			req.PassportDateofIssue,
			req.PassportValidTill,
			req.TentativeTravelFrom,
			req.TentativeTravelTo,
			req.CountryOfferingVisa,
			req.PlaceOfVisit,
			req.ReasonforVisiting,
			req.NOCReuiredForDependant,
			req.Documentdetails,
			req.Remarks,
			req.ProcessId,
			req.SequenceId,
			req.IsTaskApproved,
			req.IsTaskReturn,
			req.Role,
			req.UpdatedBy,
			requestedUserParam,
			returnToRoleParam,
			returnToUserParam,
			sendBackToMeParam,
			sendBackToUserParam,
		)

		if err != nil {
			http.Error(w, "Failed to execute stored procedure: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Step 7: Send encrypted response
		response := map[string]interface{}{
			"Status":  200,
			"Message": "Stored procedure executed successfully",
			"TaskId":  taskIdParam,
		}

		plainJSON, _ := json.Marshal(response)
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
