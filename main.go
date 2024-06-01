package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// PredictionRequest defines the structure for incoming prediction requests
type PredictionRequest struct {
	ModelName string            `json:"model_name"`
	InputData map[string]interface{} `json:"input_data"`
}

// PredictionResponse defines the structure for outgoing prediction responses
type PredictionResponse struct {
	Model     string      `json:"model"`
	Input     map[string]interface{} `json:"input"`
	Prediction interface{} `json:"prediction"`
	Confidence float64     `json:"confidence,omitempty"`
	Error     string      `json:"error,omitempty"`
}

// HealthCheckHandler provides a simple health check endpoint
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Go ML Inference Server is up and running!")
}

// PredictHandler handles incoming prediction requests
func PredictHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req PredictionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Simulate model inference
	// In a real application, this would involve loading and running an actual ML model
	log.Printf("Received prediction request for model: %s with data: %v", req.ModelName, req.InputData)

	resp := PredictionResponse{
		Model: req.ModelName,
		Input: req.InputData,
	}

	switch req.ModelName {
	case "sentiment_model":
		text, ok := req.InputData["text"].(string)
		if !ok {
			resp.Error = "Missing or invalid \"text\" in input data"
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}
		if len(text) > 0 && (text == "bad" || text == "negative") {
			resp.Prediction = "negative"
			resp.Confidence = 0.95
		} else if len(text) > 0 && (text == "good" || text == "positive") {
			resp.Prediction = "positive"
			resp.Confidence = 0.90
		} else {
			resp.Prediction = "neutral"
			resp.Confidence = 0.75
		}
	case "fraud_model":
		amount, ok := req.InputData["amount"].(float64)
		if !ok {
			resp.Error = "Missing or invalid \"amount\" in input data"
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}
		if amount > 1000.0 && random.Float64() > 0.7 { // Simulate some fraud logic
			resp.Prediction = "fraudulent"
			resp.Confidence = 0.88
		} else {
			resp.Prediction = "legitimate"
			resp.Confidence = 0.99
		}
	default:
		resp.Prediction = "unknown"
		resp.Error = fmt.Sprintf("Model ", req.ModelName, " not found or not supported")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/health", HealthCheckHandler).Methods("GET")
	router.HandleFunc("/predict", PredictHandler).Methods("POST")

	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
