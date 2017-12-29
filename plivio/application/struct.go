package application

type Application struct {
	Status string `json:"message" bson:"status"`
	AppID  string `json:"app_id" bson:"app_id"`
	APIID  string `json:"api_id" bson:"api_id"`
}
