package domain

import (
	"time"

	"github.com/google/uuid"
)

// Product represents a cannabis product in a human's personal collection
// Each product is a unique combination of strain + source + chemistry
type Product struct {
	ID           uuid.UUID          `json:"id"`
	HumanID      uuid.UUID          `json:"human_id"`
	StrainName   string             `json:"strain_name"`
	StrainType   StrainType         `json:"strain_type,omitempty"`
	Producer     string             `json:"producer,omitempty"`
	Processor    string             `json:"processor,omitempty"`
	Distributor  string             `json:"distributor,omitempty"`
	HarvestDate  *time.Time         `json:"harvest_date,omitempty"`
	Terpenes     []TerpeneProfile   `json:"terpenes,omitempty"`
	Cannabinoids []CannabinoidProfile `json:"cannabinoids,omitempty"`
	TimesUsed    int                `json:"times_used"`
	LastUsed     *time.Time         `json:"last_used,omitempty"`
	CreatedAt    time.Time          `json:"created_at"`
}

// StrainType represents the cannabis strain classification
type StrainType string

const (
	StrainTypeIndica  StrainType = "indica"
	StrainTypeSativa  StrainType = "sativa"
	StrainTypeHybrid  StrainType = "hybrid"
)

// TerpeneProfile represents a terpene with its concentration
type TerpeneProfile struct {
	Name    string  `json:"name"`
	Percent float64 `json:"percent,omitempty"`
}

// CannabinoidProfile represents a cannabinoid with its concentration
type CannabinoidProfile struct {
	Name    string  `json:"name"`
	Percent float64 `json:"percent,omitempty"`
}

// Common terpenes found in cannabis
var CommonTerpenes = []string{
	"Myrcene",
	"Limonene",
	"Caryophyllene",
	"Pinene",
	"Linalool",
	"Humulene",
	"Terpinolene",
	"Ocimene",
	"Bisabolol",
	"Eucalyptol",
}

// Common cannabinoids
var CommonCannabinoids = []string{
	"THC",
	"CBD",
	"CBN",
	"CBG",
	"THCA",
	"CBDA",
	"THCV",
	"CBC",
	"Delta-8",
}

// Validate performs domain validation on the product
func (p *Product) Validate() error {
	if p.StrainName == "" {
		return ErrInvalidProduct
	}
	return nil
}

// RecordUsage updates usage statistics when product is used
func (p *Product) RecordUsage() {
	p.TimesUsed++
	now := time.Now()
	p.LastUsed = &now
}

// HasTerpeneData returns true if terpene profile is available
func (p *Product) HasTerpeneData() bool {
	return len(p.Terpenes) > 0
}

// HasCannabinoidData returns true if cannabinoid profile is available
func (p *Product) HasCannabinoidData() bool {
	return len(p.Cannabinoids) > 0
}

// GetTotalTHC returns the total THC percentage if available
func (p *Product) GetTotalTHC() float64 {
	for _, c := range p.Cannabinoids {
		if c.Name == "THC" {
			return c.Percent
		}
	}
	return 0
}

// GetTotalCBD returns the total CBD percentage if available
func (p *Product) GetTotalCBD() float64 {
	for _, c := range p.Cannabinoids {
		if c.Name == "CBD" {
			return c.Percent
		}
	}
	return 0
}

// DominantTerpene returns the terpene with highest concentration
func (p *Product) DominantTerpene() string {
	if len(p.Terpenes) == 0 {
		return ""
	}

	dominant := p.Terpenes[0]
	for _, t := range p.Terpenes[1:] {
		if t.Percent > dominant.Percent {
			dominant = t
		}
	}
	return dominant.Name
}

// IsFresh returns true if product was harvested within 6 months
func (p *Product) IsFresh() bool {
	if p.HarvestDate == nil {
		return false
	}
	sixMonthsAgo := time.Now().AddDate(0, -6, 0)
	return p.HarvestDate.After(sixMonthsAgo)
}

// DisplayName returns a human-readable name for the product
func (p *Product) DisplayName() string {
	if p.Producer != "" {
		return p.StrainName + " (" + p.Producer + ")"
	}
	return p.StrainName
}
