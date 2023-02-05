package model

type HassDevice struct {
	Identifier        string `json:"identifiers"`
	Name              string `json:"name"`
	Model             string `json:"model"`
	Manufacturer      string `json:"manufacturer"`
	Namespace         string
	EntityPrefix      string
	AvailabilityTopic string
}

func NewHassDevice() *HassDevice {
	return &HassDevice{
		Identifier:        "vegetable-soaker",
		Name:              "Vegatable Soaker",
		Model:             "VegSoak 3000",
		Manufacturer:      "Josh Bonfield",
		Namespace:         "vegetable-soaker",
		EntityPrefix:      "vegetable-soaker-",
		AvailabilityTopic: "availability",
	}
}
