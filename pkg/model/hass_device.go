package model

type HassDevice struct {
	Identifier   string `json:"identifiers"`
	Name         string `json:"name"`
	Model        string `json:"model"`
	Manufacturer string `json:"manufacturer"`
	Namespace    string
}

func NewHassDevice() *HassDevice {
	return &HassDevice{
		Identifier:   "vegetable-soaker",
		Name:         "Vegatable Soaker",
		Model:        "vegetable-soaker",
		Manufacturer: "Josh Bonfield",
		Namespace:    "vegetable-soaker",
	}
}
