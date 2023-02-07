package model

type HassDevice struct {
	Identifier          string `json:"identifiers"`
	Name                string `json:"name"`
	Model               string `json:"model"`
	Manufacturer        string `json:"manufacturer"`
	Namespace           string `json:"-"`
	EntityPrefix        string `json:"-"`
	AvailabilityTopic   string `json:"-"`
	PayloadAvailable    string `json:"-"`
	PayloadNotAvailable string `json:"-"`
}

func NewHassDevice() *HassDevice {
	return &HassDevice{
		Identifier:          "vegetable-soaker",
		Name:                "Vegetable Soaker",
		Model:               "VegSoak 3000",
		Manufacturer:        "Josh Bonfield",
		Namespace:           "vegetable-soaker",
		EntityPrefix:        "vegetable-soaker-",
		AvailabilityTopic:   "availability",
		PayloadAvailable:    "online",
		PayloadNotAvailable: "offline",
	}
}

func (d HassDevice) GetFqAvailabilityTopic() string {
	return d.Namespace + "/availability"
}
