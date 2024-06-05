package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type Customer struct {
	ID             int     `json:"id"`
	Email          string  `json:"email"`
	HasLoyaltyCard bool    `json:"has_loyalty_card"`
	TotalPurchases float64 `json:"total_purchases"`
}

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

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(FinanceResponse{Eligible: mockFinanceApprovalSystem(customerID)})
}

func mockFinanceApprovalSystem(customerID int) bool {
	time.Sleep(10 * time.Millisecond)
	return customerID%2 == 0
}

func init() {
	http.HandleFunc("/finance", ApproveFinance)
}
