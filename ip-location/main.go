package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type IPRequest struct {
	IP string `json:"ip"`
}

type IPInfo struct {
	Country string `json:"country"`
}

func getCountryByIP(ip string) (string, error) {
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to get IP info: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var info IPInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	return info.Country, nil
}

func countryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var ipRequest IPRequest
	if err := json.NewDecoder(r.Body).Decode(&ipRequest); err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	country, err := getCountryByIP(ipRequest.IP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"country": country}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/get-country", countryHandler)

	port := "8080"
	fmt.Printf("Server is running on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

