package api

import (
	"encoding/json"
	"github.com/NathanPr03/price-control/pkg/db"
	"net/http"
	"strconv"
	"time"
)

type FinanceResponse struct {
	Eligible bool `json:"eligible"`
}

func ApproveFinance(w http.ResponseWriter, r *http.Request) {
	customerIDStr := r.URL.Query().Get("customer_id")
	if customerIDStr == "" {
		http.Error(w, "customer_id query parameter is required", http.StatusBadRequest)
		return
	}

	customerID, err := strconv.Atoi(customerIDStr)
	if err != nil {
		http.Error(w, "customer_id must be an integer", http.StatusBadRequest)
		return
	}

	doesCustExist, err := doesCustomerExist(customerID)
	if !doesCustExist {
		http.Error(w, "customer_id does not exist", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(FinanceResponse{Eligible: mockFinanceApprovalSystem(customerID)})
}

func doesCustomerExist(customerID int) (bool, error) {
	dbConnection, err := db.ConnectToDb()
	if err != nil {
		println("Error connecting to database: " + err.Error())
		return false, err
	}
	query := "SELECT COUNT(id) FROM customer WHERE id = $1"
	var count int
	err = dbConnection.QueryRow(query, customerID).Scan(&count)
	if err != nil {
		return false, err
	}
	return true, nil
}

func mockFinanceApprovalSystem(customerID int) bool {
	time.Sleep(10 * time.Millisecond)
	return customerID%2 == 0
}

func init() {
	http.HandleFunc("/finance", ApproveFinance)
}
