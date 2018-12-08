package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	unis = `
University								   : Index

Universität zu Köln 				   	   : 1
Technische Hochschule Köln 				   : 2
Cologne Business School					   : 3
Deutsche Sporthochschule Köln 			   : 4
Europäische Fachhochschule 				   : 5
Hochschule Fresenius Köln 				   : 6
Hochschule für Musik und Tanz Köln 		   : 7
Hochschule Macromedia 					   : 8
Internationale Filmschule Köln 			   : 9
Katholische Hochschule Nordrhein-Westfalen : 10
Kunsthochschule für Medien Köln 		   : 11
Rheinische Fachhochschule Köln 			   : 12
Fachhochschule der Wirtschaft 			   : 13
`
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
	app            = kingpin.New("sport", "Register for Volleyball made easy")
	register       = app.Command("register", "register player for today's course")
	registerPlayer = register.Arg("player", "player to register").Required().Strings()

	list    = app.Command("list", "list all available player")
	listUni = list.Flag("list universities", "list all available universities").Short('u').Bool()

	new         = app.Command("new", "add a new player")
	newfName    = new.Flag("firstname", "first name of new player").Required().Short('f').String()
	newlName    = new.Flag("lastname", "last name of new player").Required().Short('l').String()
	newMatrikel = new.Flag("matrikel", "matrikel number of new player").Required().Short('m').String()
	newEmail    = new.Flag("email", "email of new player").Required().Short('e').String()
	newUni      = new.Flag("university", "index of university of new player \n see list -u for details").Required().Short('u').Int()
)

func main() {
	player := loadPlayer()

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case register.FullCommand():
		fmt.Print(registerPlayer)
	case list.FullCommand():
		if *listUni {
			fmt.Print(unis)
		} else {
			for _, user := range *player.Data {
				fmt.Println(user.Vorname)
			}
		}
	case new.FullCommand():
		rData := RegisterData{
			Vorname:    *newfName,
			Nachname:   *newlName,
			Matrikel:   *newMatrikel,
			Email:      *newEmail,
			Hochschule: *newUni,
		}
		*player.Data = append(*player.Data, rData)
		writePlayer(&player)
		fmt.Printf("new Player added: \n %+v", rData)
	}

	/*
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
	*/
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

//RegisterData holds the necessary player specific data needed to register
type RegisterData struct {
	Vorname    string
	Nachname   string
	Matrikel   string
	Email      string
	Hochschule int
}

// Player holds the registerData of all created player
type Player struct {
	Data *[]RegisterData
}

func newRequest(rD *RegisterData) RequestData {
	return RequestData{
		State:             "studentAnmelden",
		TypeStudent:       "student",
		OfferCourseID:     1733,
		Vorname:           rD.Vorname,
		Nachname:          rD.Nachname,
		Matrikel:          rD.Matrikel,
		Email:             rD.Email,
		Hochschulen:       rD.Hochschule,
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

func loadPlayer() Player {
	p := Player{Data: &[]RegisterData{}}

	bPlayer, err := ioutil.ReadFile("player.json")
	if err != nil {
		handleError(err)
		return p
	}
	handleError(err)
	err = json.Unmarshal(bPlayer, &p)
	handleError(err)
	return p
}

func writePlayer(p *Player) {
	pJSON, err := json.Marshal(*p)
	handleError(err)
	err = ioutil.WriteFile("player.json", pJSON, 0600)
	handleError(err)
}

func handleError(e error) {
	if e != nil {
		fmt.Println(e.Error())
	}
}
