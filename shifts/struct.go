package shifts

type Shift struct {
	ID           string `json:"_id"`
	Name         string `json:"name"`
	AbsentWorker struct {
		Name string `json:"name"`
		ID   string `json:"_id"`
	} `json:"absentWorker"`
	Job struct {
		Title        string  `json:"title"`
		Compensation float64 `json:"compensation"`
		ID           string  `json:"_id"`
	} `json:"job"`
	Location   string        `json:"location"`
	Date       string        `json:"date"`
	StartTime  string        `json:"startTime"`
	EndTime    string        `json:"endTime"`
	Volunteers []interface{} `json:"volunteers"`
	V          int           `json:"__v"`
}