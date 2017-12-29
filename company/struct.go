package company

type Company struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

func (company Company) GetAvailableWorkers() {
	//TODO get available workers probably by implementing a graph database solution
}
