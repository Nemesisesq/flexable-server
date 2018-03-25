package company

type Company struct {
	Name string `json:"name" bson:"name"`
	UUID string `json:"uuid" bson:"uuid"`
}

func (company Company) GetAvailableWorkers() {
	//TODO get available workers probably by implementing a graph database solution
}
