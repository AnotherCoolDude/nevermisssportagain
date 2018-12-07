package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/parnurzeal/gorequest"
)

/*
01 Universität zu Köln 						  : 1
02 Technische Hochschule Köln 				  : 2
06 Cologne Business School					  : 3
07 Deutsche Sporthochschule Köln 			  : 4
09 Europäische Fachhochschule 				  : 5
08 Hochschule Fresenius Köln 				  : 6
04 Hochschule für Musik und Tanz Köln 		  : 7
10 Hochschule Macromedia 					  : 8
13 Internationale Filmschule Köln 			  : 9
05 Katholische Hochschule Nordrhein-Westfalen : 10
12 Kunsthochschule für Medien Köln 			  : 11
03 Rheinische Fachhochschule Köln 			  : 12
11 Fachhochschule der Wirtschaft 			  : 13
*/

var (
	app          = kingpin.New("sport", "Register for Volleyball made easy")
	register     = app.Command("register", "register player for today's course")
	registerUser = register.Arg("user", "user to register").Required().Strings()
)

func main() {

	rd := registerData{
		vorname:    "Christian",
		nachname:   "Hovenbitzer",
		matrikel:   "271172024",
		email:      "christian.hovenbitzer@gmx.net",
		hochschule: 12,
	}

	requestData := newRequest(&rd)

	request := gorequest.New()
	request.SetDebug(true)

	// url = https://anmeldung.hochschulsport-koeln.de/inc/methods.php
	_, body, errors := request.Post("127.0.0.1:8080/test").Send(requestData.jsonString()).End()

	if errors != nil {
		fmt.Println(request.Data)
		fmt.Println(errors)
	} else {
		fmt.Println(body)
	}
}

// RequestData wraps the necessary data into a struct
type RequestData struct {
	State             string `json:"state"`
	TypeStudent       string `json:"type"`
	OfferCourseID     int    `json:"offerCourseID"`
	Vorname           string `json:"vorname"`
	Nachname          string `json:"nachname"`
	Matrikel          string `json:"matrikel"`
	Email             string `json:"email"`
	Hochschulen       int    `json:"hochschulen"`
	Hochschulenextern string `json:"hochschulenextern"`
	Office            string `json:"office"`
}

type registerData struct {
	vorname    string
	nachname   string
	matrikel   string
	email      string
	hochschule int
}

type player struct {
	data []registerData
}

func newRequest(rD *registerData) RequestData {
	return RequestData{
		State:             "studentAnmelden",
		TypeStudent:       "student",
		OfferCourseID:     1733,
		Vorname:           rD.vorname,
		Nachname:          rD.nachname,
		Matrikel:          rD.matrikel,
		Email:             rD.email,
		Hochschulen:       rD.hochschule,
		Hochschulenextern: "null",
		Office:            "null",
	}
}

func (requestData RequestData) jsonString() string {
	mJSON, err := json.Marshal(requestData)
	if err != nil {
		fmt.Printf("could not marshal struct: %s", err)
	}
	return string(mJSON)
}

func loadPlayer() player {
	p := player{}
	if _, err := os.Stat("player.json"); os.IsNotExist(err) {
		writePlayer(&p)
		return p
	}
	bPlayer, err := ioutil.ReadFile("player.json")
	handleError(err)
	err = json.Unmarshal(bPlayer, p)
	handleError(err)
	return p
}

func writePlayer(p *player) {
	pJSON, err := json.Marshal(p)
	handleError(err)
	err = ioutil.WriteFile("player.json", pJSON, 0644)
	handleError(err)
}

func handleError(e error) {
	if e != nil {
		fmt.Println(e.Error())
	}
}
