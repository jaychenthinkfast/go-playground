package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"go-playground/pkg/runner"
	"go-playground/pkg/sandbox"
)

type RunRequest struct {
	Code    string `json:"code"`
	Version string `json:"version"`
}

type RunResponse struct {
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

type FormatRequest struct {
	Code string `json:"code"`
}

type FormatResponse struct {
	FormattedCode string `json:"formattedCode"`
	Error         string `json:"error,omitempty"`
}

func main() {
	r := mux.NewRouter()

	// API routes
	r.HandleFunc("/api/run", handleRun).Methods("POST")
	r.HandleFunc("/api/format", handleFormat).Methods("POST")
	r.HandleFunc("/api/health", handleHealth).Methods("GET")

	// Create a CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3002", "http://localhost:3003"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Wrap the router with the CORS middleware
	handler := c.Handler(r)

	// Start the server
	port := 3001
	fmt.Printf("Starting server on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}

func handleRun(w http.ResponseWriter, r *http.Request) {
	var req RunRequest
	var resp RunResponse

	// Parse the JSON request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate the request
	if req.Code == "" {
		http.Error(w, "Code cannot be empty", http.StatusBadRequest)
		return
	}

	// Validate the version
	if !runner.IsValidVersion(req.Version) {
		desc := fmt.Sprintf("Unsupported Go version. Available versions: %s, %s, %s, %s",
			runner.Go125, runner.Go124, runner.Go123, runner.Go122)
		resp.Error = desc
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Add version information to output
	versionDesc := runner.Versions[req.Version]

	// Create a sandbox to run the code
	s := sandbox.NewSandbox()

	// Set a timeout for the execution
	ctx, cancel := sandbox.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Run the code in the sandbox
	output, err := runner.Run(ctx, s, req.Code, req.Version)
	if err != nil {
		resp.Error = err.Error()
		// Even if there's an error, include any output that was produced
		resp.Output = output
	} else {
		// Prepend version description to output
		if output != "" && versionDesc != "" {
			resp.Output = fmt.Sprintf("%s\n%s", versionDesc, output)
		} else {
			resp.Output = output
		}
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleFormat(w http.ResponseWriter, r *http.Request) {
	var req FormatRequest
	var resp FormatResponse

	// Parse the JSON request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate the request
	if req.Code == "" {
		http.Error(w, "Code cannot be empty", http.StatusBadRequest)
		return
	}

	// Format the code
	formattedCode, err := runner.Format(req.Code)
	if err != nil {
		resp.Error = err.Error()
	} else {
		resp.FormattedCode = formattedCode
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
