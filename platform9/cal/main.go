package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CalculationRequest struct {
	Operand1 float64 `json:"operand1"`
	Operand2 float64 `json:"operand2"`
	Operator string  `json:"operator"`
}

type CalculationResponse struct {
	Result float64 `json:"result"`
	Error  string  `json:"error,omitempty"`
}

func calculate(w http.ResponseWriter, r *http.Request) {
	var req CalculationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var res CalculationResponse
	switch req.Operator {
	case "+":
		res.Result = req.Operand1 + req.Operand2
	case "-":
		res.Result = req.Operand1 - req.Operand2
	case "*":
		res.Result = req.Operand1 * req.Operand2
	case "/":
		if req.Operand2 == 0 {
			res.Error = "division by zero"
		} else {
			res.Result = req.Operand1 / req.Operand2
		}
	default:
		res.Error = "invalid operator"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func main() {
	http.HandleFunc("/calculate", calculate)
	http.Handle("/", http.FileServer(http.Dir(".")))
	fmt.Println("Server is running on http://localhost:8081")
	http.ListenAndServe(":8081", nil)
}

