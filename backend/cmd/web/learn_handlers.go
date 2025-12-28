package web

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

// CannabinoidData represents the structure from our Supabase API
type CannabinoidData struct {
	Name                string                 `json:"name"`
	FullName           string                 `json:"full_name"`
	Description        string                 `json:"description"`
	Psychoactive       bool                   `json:"psychoactive"`
	ReportedExperiences map[string]interface{} `json:"reported_experiences"`
	CompoundNotes      map[string]interface{} `json:"compound_notes"`
}

// CannabinoidsHandler serves the educational cannabinoids page
func CannabinoidsHandler(c *gin.Context) {
	// Get configuration from environment variables
	supabaseURL := os.Getenv("SUPABASE_URL")
	apiKey := os.Getenv("SUPABASE_ANON_KEY")
	
	if supabaseURL == "" || apiKey == "" {
		log.Printf("Error: Missing Supabase configuration")
		c.String(http.StatusInternalServerError, "Configuration error")
		return
	}
	
	// Create HTTP client and request
	client := &http.Client{}
	fullURL := supabaseURL + "/rest/v1/cannabinoids?select=name,full_name,description,psychoactive,reported_experiences,compound_notes&order=name"
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		c.String(http.StatusInternalServerError, "Error fetching data")
		return
	}

	// Add headers
	req.Header.Add("apikey", apiKey)
	req.Header.Add("Authorization", "Bearer "+apiKey)

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		c.String(http.StatusInternalServerError, "Error fetching data")
		return
	}
	defer resp.Body.Close()

	// Parse response
	var cannabinoids []CannabinoidData
	if err := json.NewDecoder(resp.Body).Decode(&cannabinoids); err != nil {
		log.Printf("Error decoding response: %v", err)
		c.String(http.StatusInternalServerError, "Error parsing data")
		return
	}

	// Convert to map format for template
	var cannabinoidMaps []map[string]interface{}
	for _, cannabinoid := range cannabinoids {
		cannabinoidMap := map[string]interface{}{
			"name":                 cannabinoid.Name,
			"full_name":           cannabinoid.FullName,
			"description":         cannabinoid.Description,
			"psychoactive":        cannabinoid.Psychoactive,
			"reported_experiences": cannabinoid.ReportedExperiences,
			"compound_notes":      cannabinoid.CompoundNotes,
		}
		cannabinoidMaps = append(cannabinoidMaps, cannabinoidMap)
	}

	// Render template
	component := CannabinoidsPage(cannabinoidMaps)
	c.Header("Content-Type", "text/html")
	component.Render(c.Request.Context(), c.Writer)
}