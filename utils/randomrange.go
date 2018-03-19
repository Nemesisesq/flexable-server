package utils

import (
	"github.com/kr/pretty"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func RandomRange(min, max int) int {
	rand.Seed(time.Now().Unix())
	myrand := random(min, max)
	return myrand
}

func PrintBody(r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err)
	}
	pretty.Println(body)
}
