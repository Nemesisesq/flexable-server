package auth0

import (
	"net/http"
)

var baseUri string = "flexable.auth0.com"

func Login(username string, password string) interface{} {
	client := http.Client{}

	params := map[string]interface{}{
		"response_type": "code|token",
		"client_id":     "h0z5X7yRIhMgSCwh4jyw1b6DXYb4zNP5",
		"connection":    "CONNECTION",
		"redirect_uri":  "https://auth.expo.io/@nemesisesq/expo-flexable",
		"state":         "STATE",
		//ADDITIONAL_PARAMETERS

	}
}
