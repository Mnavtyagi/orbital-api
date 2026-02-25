package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	http.HandleFunc("/crop", cropHandler)
	http.HandleFunc("/crop-json", cropJSONHandler)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func cropHandler(w http.ResponseWriter, r *http.Request) {
	lat := r.URL.Query().Get("lat")
	lng := r.URL.Query().Get("lng")

	// Check missing params
	if lat == "" || lng == "" {
		http.Error(w, "Missing lat or lng", http.StatusBadRequest)
		return
	}

	// Validate numeric input
	_, err1 := strconv.ParseFloat(lat, 64)
	_, err2 := strconv.ParseFloat(lng, 64)

	if err1 != nil || err2 != nil {
		http.Error(w, "Invalid coordinates", http.StatusBadRequest)
		return
	}

	// Call Python script
	cmd := exec.Command("python", "../../processor/process.py", lat, lng)
	output, err := cmd.CombinedOutput()

	if err != nil {
		http.Error(w, string(output), http.StatusInternalServerError)
		return
	}

	filename := strings.TrimSpace(string(output))
	imagePath := "../../data/" + filename

	w.Header().Set("Content-Type", "image/png")
	http.ServeFile(w, r, imagePath)
}

func cropJSONHandler(w http.ResponseWriter, r *http.Request) {
	lat := r.URL.Query().Get("lat")
	lng := r.URL.Query().Get("lng")

	if lat == "" || lng == "" {
		http.Error(w, "Missing lat or lng", http.StatusBadRequest)
		return
	}

	_, err1 := strconv.ParseFloat(lat, 64)
	_, err2 := strconv.ParseFloat(lng, 64)

	if err1 != nil || err2 != nil {
		http.Error(w, "Invalid coordinates", http.StatusBadRequest)
		return
	}

	cmd := exec.Command("python", "../../processor/process.py", lat, lng)
	output, err := cmd.CombinedOutput()

	if err != nil {
		http.Error(w, string(output), http.StatusInternalServerError)
		return
	}

	filename := strings.TrimSpace(string(output))

	response := map[string]interface{}{
		"lat":           lat,
		"lng":           lng,
		"image_file":    filename,
		"buffer_meters": 500,
		"source":        "Sample GeoTIFF",
		"status":        "success",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}