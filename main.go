package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/apoorvam/goterminal"
	tm "github.com/buger/goterm"
	"github.com/parnurzeal/gorequest"
)

var (
	app               = kingpin.New("sport", "Register for Volleyball made easy")
	cmdRegister       = app.Command("register", "register player for today's course")
	argRegisterPlayer = cmdRegister.Arg("player", "player to register").Required().Strings()

	list    = app.Command("list", "list all available player")
	listUni = list.Flag("list universities", "list all available universities").Short('u').Bool()

	new         = app.Command("new", "add a new player")
	newfName    = new.Flag("firstname", "first name of new player").Required().Short('f').String()
	newlName    = new.Flag("lastname", "last name of new player").Required().Short('l').String()
	newMatrikel = new.Flag("matrikel", "matrikel number of new player").Required().Short('m').String()
	newEmail    = new.Flag("email", "email of new player").Required().Short('e').String()
	newUni      = new.Flag("university", "index of university of new player \n see list -u for details").Required().Short('u').Int()

	uniMap = map[string]int{
		"Universität zu Köln":                        1,
		"Technische Hochschule Köln":                 2,
		"Cologne Business School":                    3,
		"Deutsche Sporthochschule Köln":              4,
		"Europäische Fachhochschule":                 5,
		"Hochschule Fresenius Köln":                  6,
		"Hochschule für Musik und Tanz Köln":         7,
		"Hochschule Macromedia":                      8,
		"Internationale Filmschule Köln":             9,
		"Katholische Hochschule Nordrhein-Westfalen": 10,
		"Kunsthochschule für Medien Köln":            11,
		"Rheinische Fachhochschule Köln":             12,
		"Fachhochschule der Wirtschaft":              13,
	}
)

// RequestData wraps the necessary data into a struct
type RequestData struct {
	State             string `json:"state"`
	TypeStudent       string `json:"type"`
	OfferCourseID     string `json:"offerCourseID"`
	Vorname           string `json:"vorname"`
	Nachname          string `json:"nachname"`
	Matrikel          string `json:"matrikel"`
	Email             string `json:"email"`
	Hochschulen       int    `json:"hochschulen"`
	Hochschulenextern string `json:"hochschulenextern"`
	Office            string `json:"office"`
}

//RegisterData holds the necessary, player specific data needed to register
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

func main() {
	player := loadPlayer()

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case cmdRegister.FullCommand():
		defer fmt.Println("leaving main..")
		ready := make(chan bool)

		scheduleRegistration(ready)
		<-ready
		player.register(*argRegisterPlayer)
	case list.FullCommand():
		if *listUni {
			printUniversities()
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
		fmt.Printf("new Player added: \n%s", rData.string())
	}
}

func newRequest(rD *RegisterData) RequestData {
	return RequestData{
		State:             "studentAnmelden",
		TypeStudent:       "student",
		OfferCourseID:     "+17+",
		Vorname:           rD.Vorname,
		Nachname:          rD.Nachname,
		Matrikel:          rD.Matrikel,
		Email:             rD.Email,
		Hochschulen:       rD.Hochschule,
		Hochschulenextern: "null",
		Office:            "null",
	}
}

func (rd RequestData) jsonString() string {
	mJSON, err := json.Marshal(rd)
	if err != nil {
		fmt.Printf("could not marshal struct: %s", err)
	}
	return string(mJSON)
}

func (rd *RequestData) formEncoded() url.Values {
	form := url.Values{
		"state":             {rd.State},
		"type":              {rd.TypeStudent},
		"offerCourseID":     {rd.OfferCourseID},
		"vorname":           {rd.Vorname},
		"nachname":          {rd.Nachname},
		"matrikel":          {rd.Matrikel},
		"email":             {rd.Email},
		"hochschulen":       {string(rd.Hochschulen)},
		"hochschulenextern": {rd.Hochschulenextern},
		"office":            {rd.Office},
	}
	fmt.Print(form)
	return form
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

func printUniversities() {
	t := tm.NewTable(0, 10, 5, ' ', 0)
	fmt.Fprintf(t, "University\tIndex\n")
	fmt.Fprint(t, "\n")
	for uni, index := range uniMap {
		fmt.Fprintf(t, "%s\t%d\n", uni, index)
	}
	fmt.Print(t.String())
}

func (rd *RegisterData) string() string {
	uni, _ := mapkey(uniMap, rd.Hochschule)
	t := tm.NewTable(0, 10, 5, ' ', 0)
	fmt.Fprintf(t, "Vorname\t%s\n", rd.Vorname)
	fmt.Fprintf(t, "Nachname\t%s\n", rd.Nachname)
	fmt.Fprintf(t, "Matrikel\t%s\n", rd.Matrikel)
	fmt.Fprintf(t, "Email\t%s\n", rd.Email)
	fmt.Fprintf(t, "University\t%s\n", uni)
	return t.String()
}

func mapkey(m map[string]int, value int) (key string, ok bool) {
	for k, v := range m {
		if v == value {
			key = k
			ok = true
			return
		}
	}
	return
}

func (p *Player) register(names []string) {
	fmt.Printf("names: %v\n", names)
	pToRegister := filter(*p.Data, func(i int, rd RegisterData) bool {
		return rd.Vorname == names[i]
	})
	fmt.Printf("pToRegister: %v\n", pToRegister)
	for _, p := range pToRegister {
		fmt.Printf("registering %s\n", p.Vorname)
		requestData := newRequest(&p)

		request := gorequest.New()
		request.SetDebug(true)
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		request.FormData = requestData.formEncoded()
		// url = https://anmeldung.hochschulsport-koeln.de/inc/methods.php
		_, body, errors := request.Post("/t/kntu3-1544613998/post").End()

		if errors != nil {
			fmt.Println(request.Data)
			fmt.Println(errors)
		} else {
			fmt.Println(body)
		}
	}

}

func filter(vs []RegisterData, f func(int, RegisterData) bool) []RegisterData {
	vsf := make([]RegisterData, 0)
	fmt.Printf("vs: %v\n", vs)
	for i, v := range vs {

		if f(i, v) {
			fmt.Printf("index: %d, filtered Name: %v\n", i, v)
			vsf = append(vsf, v)
		}
		fmt.Printf("vsf: %v\n", vsf)
	}

	return vsf
}

func calculateRegisterStart() time.Time {
	tuesday := time.Tuesday
	now := time.Now()
	wDay := now.Weekday()
	diff := 0
	if wDay < tuesday {
		diff = int(tuesday) - int(wDay)
	}
	if wDay > tuesday {
		diff = 7 - (int(wDay) - int(tuesday))
	}
	return time.Date(now.Year(), now.Month(), now.Day(), 15, 30, 10, 0, now.Location()).AddDate(0, 0, diff)
}

func scheduleRegistration(ready chan bool) {
	writer := goterminal.New(os.Stdout)
	ticker := time.NewTicker(time.Second)
	done := make(chan bool)
	test := time.Now().Add(3 * time.Second)
	defer writer.Reset()

	for {
		select {
		case <-done:
			go func() {
				fmt.Fprintln(writer, "registering now...")
				writer.Print()
				ready <- true
			}()
			return
		case t := <-ticker.C:
			go func() {
				countdown := test.Sub(t) //calculateRegisterStart().Sub(t)
				writer.Clear()
				fmt.Fprint(writer, fmtDuration(countdown))
				writer.Print()
				if countdown.Seconds() <= 0 {
					ticker.Stop()
					done <- true
				}
			}()
		}
	}

}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d until registration\n", h, m, s)
}
