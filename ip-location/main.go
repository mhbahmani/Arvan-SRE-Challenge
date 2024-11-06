package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type IPRequest struct {
	IP string `json:"ip"`
}

type IPResponse struct {
	Country            string `json:"country"`
	RequestedIP        string `json:"requested_ip"`
	FetchedFromDatabase bool   `json:"fetched_from_database"`
}

type IPInfo struct {
	Country string `json:"country"`
}

func connectDB() (*sql.DB, error) {
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	return sql.Open("postgres", dsn)
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

	country, fromDB, err := getCountry(ipRequest.IP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := IPResponse{
		Country:            country,
		RequestedIP:        ipRequest.IP,
		FetchedFromDatabase: fromDB,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getCountry(ip string) (string, bool, error) {
	db, err := connectDB()
	if err != nil {
		return "", false, err
	}
	defer db.Close()

	var country string
	err = db.QueryRow("SELECT country FROM ip_locations WHERE ip = $1", ip).Scan(&country)
	if err == nil {
		return country, true, nil
	}
	if err != sql.ErrNoRows {
		return "", false, err
	}

	country, err = fetchCountryFromAPI(ip)
	if err != nil {
		return "", false, err
	}

	go func() {
		_, err := db.Exec("INSERT INTO ip_locations (ip, country) VALUES ($1, $2)", ip, country)
		if err != nil {
			log.Printf("failed to insert into db: %v", err)
		}
	}()

	return country, false, nil
}

func fetchCountryFromAPI(ip string) (string, error) {
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	// TODO: Add alternative
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

func main() {
	http.HandleFunc("/get-country", countryHandler)

	port := "8080"
	fmt.Printf("Server is running on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}