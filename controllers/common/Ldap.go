// This package handles LDAP-based login authentication, secure credential decryption,
// session initialization and logging, and communication with MSSQL for session recording.
// Responses are optionally encrypted using AES-GCM encryption from the utils package.
//
// It integrates with:
//   - auth: to handle request authentication, logging, and IP/token validation
//   - utils: for encrypting responses and cryptographic utilities
//   - MSSQL and LDAP for backend data sources
//
// Core features:
//   - Decrypts and validates user credentials (ONLY ACCEPTS ENCRYPTED INPUT)
//   - Performs multi-domain LDAP lookup (staff, faculty, project)
//   - Inserts session data into the MSSQL `Session_Data` table
//   - Encrypts API responses for enhanced security
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
	credentials "Hrmodule/dbconfig"
	"Hrmodule/utils"
	"crypto/aes"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	ldap "github.com/go-ldap/ldap/v3"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Valid      bool   `json:"valid"`
	UserId     string `json:"userId,omitempty"`
	Username   string `json:"username,omitempty"`
	Role       string `json:"role,omitempty"`
	EmployeeId string `json:"EmployeeId"`
}

type AuthResponsefalse struct {
	Valid    bool   `json:"valid"`
	Username string `json:"username,omitempty"`
	Error    string `json:"error,omitempty"`
}

// Helper function to check if string is hex-encoded
func isHexString(s string) bool {
	if len(s)%2 != 0 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}

// decryptData decrypts hex-encoded encrypted data using AES
func decryptData(encryptedData, key string) (string, error) {
	keyBytes := []byte(key)
	encryptedBytes, err := hex.DecodeString(encryptedData)
	if err != nil {
		return "", fmt.Errorf("invalid hex encoding: %v", err)
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}

	if len(encryptedBytes)%aes.BlockSize != 0 {
		return "", fmt.Errorf("encrypted data is not a multiple of the block size")
	}

	decrypted := make([]byte, len(encryptedBytes))
	for i := 0; i < len(encryptedBytes); i += aes.BlockSize {
		block.Decrypt(decrypted[i:i+aes.BlockSize], encryptedBytes[i:i+aes.BlockSize])
	}

	// Remove padding
	decrypted = PKCS5Unpad(decrypted)

	return string(decrypted), nil
}

// decryptDataStrict only accepts encrypted (hex-encoded) data
func decryptDataStrict(data, key string) (string, error) {
	// Check if data is hex-encoded (encrypted)
	if !isHexString(data) {
		return "", fmt.Errorf("invalid input: data must be encrypted (hex-encoded)")
	}

	// Only decrypt if it's valid hex
	return decryptData(data, key)
}

// validateEncryptedCredentials validates that both username and password are encrypted
func validateEncryptedCredentials(username, password string) (bool, string) {
	if username == "" || password == "" {
		return false, "Missing username or password"
	}

	if !isHexString(username) {
		return false, "Invalid username format - must be encrypted (hex-encoded)"
	}

	if !isHexString(password) {
		return false, "Invalid password format - must be encrypted (hex-encoded)"
	}

	return true, ""
}

// PKCS5Unpad removes padding from decrypted data
func PKCS5Unpad(data []byte) []byte {
	pad := int(data[len(data)-1])
	return data[:len(data)-pad]
}

// HandleLDAPAuth processes an HTTP request for LDAP authentication.
// It ONLY accepts encrypted credentials, validates them against LDAP servers (staff, faculty, project),
// inserts session data into the database, and returns an encrypted JSON response.
func HandleLDAPAuth(w http.ResponseWriter, r *http.Request) {

	// handles sp validation
	authorized := auth.HandleRequestfor_apiname_ipaddress_token(w, r)
	if !authorized {
		return
	}

	// For getting clientipaddress
	loggedHandler := auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// For id parameter
		err := auth.IsValid_IDFromRequest(r)
		if err != nil {
			http.Error(w, "Invalid Token provided", http.StatusBadRequest)
			return
		}

		username := r.URL.Query().Get("username")
		password := r.URL.Query().Get("password")

		// Validate that credentials are encrypted
		valid, errorMsg := validateEncryptedCredentials(username, password)
		if !valid {
			log.Printf("Validation error: %s", errorMsg)
			resp := AuthResponsefalse{
				Valid: false,
				Error: errorMsg,
			}
			jsonResponse, err := json.Marshal(resp)
			if err != nil {
				http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
				return
			}
			encrypted, err := utils.Encrypt(jsonResponse)
			if err != nil {
				http.Error(w, "Encryption failed", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{
				"Data": encrypted,
			})
			return
		}

		// Print encrypted credentials for debugging
		fmt.Println("Encrypted username:", username)
		fmt.Println("Encrypted password:", password)

		// Decrypt username and password - ONLY accept encrypted data
		key := "7xPz!qL3vNc#eRb9Wm@f2Zh8Kd$gYp1B"

		fmt.Println("Decryption key:", key)

		// Decrypt username using strict decryption
		decodedUsername, err := decryptDataStrict(username, key)
		if err != nil {
			log.Printf("Error decrypting username: %v", err)
			resp := AuthResponsefalse{
				Valid:    false,
				Username: "Invalid",
				Error:    "Username decryption failed",
			}
			jsonResponse, err := json.Marshal(resp)
			if err != nil {
				http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
				return
			}
			encrypted, err := utils.Encrypt(jsonResponse)
			if err != nil {
				http.Error(w, "Encryption failed", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{
				"Data": encrypted,
			})
			return
		}

		fmt.Println("Decrypted Username:", decodedUsername)

		// Decrypt password using strict decryption
		decodedPassword, err := decryptDataStrict(password, key)
		if err != nil {
			log.Printf("Error decrypting password: %v", err)
			resp := AuthResponsefalse{
				Valid:    false,
				Username: "Invalid",
				Error:    "Password decryption failed",
			}
			jsonResponse, err := json.Marshal(resp)
			if err != nil {
				http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
				return
			}
			encrypted, err := utils.Encrypt(jsonResponse)
			if err != nil {
				http.Error(w, "Encryption failed", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{
				"Data": encrypted,
			})
			return
		}
		fmt.Println("Decrypted Password:", decodedPassword)

		// Continue with LDAP authentication using decodedUsername and decodedPassword...
		dn := "cn=academicbind,ou=bind,dc=ldap,dc=iitm,dc=ac,dc=in"
		pass := "1@iIL~0K"
		ldapUserFilter := "(&(objectclass=*)(uid=" + decodedUsername + "))"
		searchBaseStaff := "ou=staff,ou=people,dc=ldap,dc=iitm,dc=ac,dc=in"
		searchBaseFaculty := "ou=faculty,ou=people,dc=ldap,dc=iitm,dc=ac,dc=in"
		searchbase_project := "ou=project,ou=employee,dc=ldap,dc=iitm,dc=ac,dc=in"
		//	searchBaseStudent := "ou=student,dc=ldap,dc=iitm,dc=ac,dc=in"  //comment for later use
		//fmt.Println(decodedUsername)
		ldapURL := "ldap://ldap.iitm.ac.in:389"

		conn, err := ldap.DialURL(ldapURL)
		if err != nil {
			log.Printf("Failed to connect to LDAP server: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		err = conn.Bind(dn, pass)
		if err != nil {
			log.Printf("Server DN Bind Failed: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var ou string
		var responseSent bool

		performSearch := func(searchBase, userType string) {
			req := ldap.NewSearchRequest(
				searchBase,
				ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
				ldapUserFilter,
				nil,
				nil,
			)

			sr, err := conn.Search(req)
			if err != nil {
				log.Printf("Search Failed: %v", err)
				return
			}

			for _, entry := range sr.Entries {
				dn := entry.DN
				bindCredentials := decodedPassword

				if dn != "" && bindCredentials != "" {
					err = conn.Bind(dn, bindCredentials)
					if err != nil {
						log.Printf("%s Bind Failed: %v", userType, err)
					} else {
						log.Printf("%s Bind Successful", userType)
						ou = userType

						if !responseSent {
							responseSent = true

							userId := generateUserId()
							role, empId, err := getDefaultRole(decodedUsername)
							if err != nil {
								log.Printf("Error retrieving default role: %v", err)
								http.Error(w, "Internal Server Error", http.StatusInternalServerError)
								return
							}

							err = insertSessionData(userId, decodedUsername, ou,empId)
							if err != nil {
								log.Printf("Error inserting session data: %v", err)
								http.Error(w, "Internal Server Error", http.StatusInternalServerError)
								return
							}

							resp := AuthResponse{
								Valid:      true,
								UserId:     userId,
								Username:   decodedUsername,
								Role:       role,
								EmployeeId: empId,
							}

							jsonResponse, err := json.Marshal(resp)
							if err != nil {
								http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
								return
							}
							encrypted, err := utils.Encrypt(jsonResponse)
							if err != nil {
								http.Error(w, "Encryption failed", http.StatusInternalServerError)
								return
							}
							w.Header().Set("Content-Type", "application/json")
							_ = json.NewEncoder(w).Encode(map[string]string{
								"Data": encrypted,
							})

						}
					}
				} else {
					log.Println("DN or password is null or undefined")
				}
			}
		}

		// Check staff first
		performSearch(searchBaseStaff, "staff")
		// If not staff, check faculty
		performSearch(searchBaseFaculty, "faculty")
		// If not faculty, check project
		performSearch(searchbase_project, "project")
		// If not faculty, check student
		//		performSearch(searchBaseStudent, "student")  //comment for later use

		if ou == "" {
			log.Println("LDAP Entries Mismatch")
			if !responseSent {
				responseSent = true

				userId := generateUserId()

				// Try to fetch role and employeeId even if bind fails
				role, empId, err := getDefaultRole(decodedUsername)
				if err != nil {
					log.Printf("Error retrieving role for failed LDAP: %v", err)
					role = ""
					empId = ""
				}

				ou := "unauthorized" // Track that this was a failed login attempt
				_ = insertSessionData(userId, decodedUsername, ou,empId)

				resp := AuthResponse{
					Valid:      true,
					UserId:     userId,
					Username:   decodedUsername,
					Role:       role,
					EmployeeId: empId,
				}

				jsonResponse, err := json.Marshal(resp)
				if err != nil {
					http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
					return
				}
				encrypted, err := utils.Encrypt(jsonResponse)
				if err != nil {
					http.Error(w, "Encryption failed", http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(map[string]string{
					"Data": encrypted,
				})

			} else {
				// If LDAP authentication succeeds for any user type,
				// reset the responseSent flag to false
				responseSent = false
			}
		}

	}))
	loggedHandler.ServeHTTP(w, r)
}

// generateUserId creates and returns a new UUID string.
func generateUserId() string {
	return uuid.New().String()
}

// insertSessionData inserts a session record into the MSSQL Session_Data table.
func insertSessionData(userId, username, ou,employeeId string) error {

	// Connection string for MSSQL Server
	connectionString := credentials.GetTestdatabase15()

	// Open MSSQL database connection
	db, err := sql.Open("mssql", connectionString)
	if err != nil {
		return fmt.Errorf("DB open error: %v", err)
	}
	defer db.Close()
	// Insert the datas in db
	query := `INSERT INTO Standard..Session_Data (Session_Id, Logout_Date, Username, Role, Is_Active,Department,User_id,Employee_id,Login_Date) VALUES (?,  ?, ?, ?, ?, ?,?,?,GETDATE())`
	_, err = db.Exec(query, userId, nil, username, ou, "1", "department", "User_id", employeeId)
	if err != nil {
		return err
	}
	return nil
}

// getDefaultRole queries the Axon.ADSUserDetails table to retrieve a default role and employee ID for a given user.
func getDefaultRole(username string) (role string, employeeId string, err error) {

	// Connection string for MSSQL Server
	connectionString := credentials.GetTestdatabase15()

	// Open connection to the database
	db, err := sql.Open("mssql", connectionString)
	if err != nil {
		return "", "", err
	}
	defer db.Close()

	query := `SELECT DefaultRoleName, EmployeeId FROM Axon..ADSUserDetails WHERE LoginName = ?`
	row := db.QueryRow(query, username)

	// Scan both role and employee ID
	err = row.Scan(&role, &employeeId)
	if err != nil {
		return "", "", err
	}

	return role, employeeId, nil
}
