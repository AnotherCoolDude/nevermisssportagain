package main

import (
	"encoding/json"
	"fmt"

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
	mJSON, err := json.Marshal(requestData)

	if err != nil {
		fmt.Printf("could not marshal struct: %s", err)
	} else {
		fmt.Println(string(mJSON))
	}
	// url = https://anmeldung.hochschulsport-koeln.de/inc/methods.php
	_, body, error := request.Post("127.0.0.1:8080/test").Send(string(mJSON)).End()

	if error != nil {
		fmt.Println(request.Data)
		for e := range error {
			fmt.Printf("Error from request: %v\n", e)
		}
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

func newRequest(rD *registerData) RequestData {
	return RequestData{
		State:             "studentAnmelden",
		TypeStudent:       "student",
		OfferCourseID:     30,
		Vorname:           rD.vorname,
		Nachname:          rD.nachname,
		Matrikel:          rD.matrikel,
		Email:             rD.email,
		Hochschulen:       rD.hochschule,
		Hochschulenextern: "null",
		Office:            "null",
	}
}
