package company

type Location interface{}

type Company struct {
	Name     string   `json:"name"`
	Location Location `json:"location"`
}
