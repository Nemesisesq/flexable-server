package employee

type GeoLocation struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type Position struct {
	Title        string  `json:"title"`
	Compensation float32 `json:"compensation"`
	Rate         string  `json:"rate"`
}

type Employee struct {
	Name     string      `json:"name"`
	Number   string      `json:"number"`
	Email    string      `json:"email"`
	Location GeoLocation `json:"location"`
	Position Position
}
