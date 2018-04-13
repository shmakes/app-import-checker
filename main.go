package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	flag "github.com/ogier/pflag"
	"io/ioutil"
	"net/http"
	"os"
    "strings"
)

type AppData struct {
	Submitted string `json:"Submitted"`
	LastName  string `json:"Last-Name"`
	Email     string `json:"from_email"`
	IpAddress string `json:"ip_address"`
}

type Row struct {
	Id    string   `json:"id"`
	Key   []string `json:"key"`
	Value []string `json:"value"`
}

type ReviewData struct {
	TotalRows int   `json:"total_rows"`
	Rows      []Row `json:"rows"`
}

func (appData AppData) toString() string {
	return toJson(appData)
}

func (reviewData Row) toString() string {
	return toJson(reviewData)
}

func toJson(d interface{}) string {
	bytes, err := json.Marshal(d)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return string(bytes)
}

func getApps() []AppData {
	raw, err := ioutil.ReadFile(infile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c []AppData
	json.Unmarshal(raw, &c)
	return c
}

func getReviews() ReviewData {
	rs, err := http.Get(dburl)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer rs.Body.Close()

	bodyBytes, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var r ReviewData
	json.Unmarshal(bodyBytes, &r)
	return r
}

var (
	infile string
	dburl  string
)

func main() {
	flag.Parse()

	if flag.NFlag() == 0 {
		printUsage()
	}

	var m map[string]int
	m = make(map[string]int)
	var revKeyBuf bytes.Buffer

	reviews := getReviews()
	for _, r := range reviews.Rows {
		revDate := r.Key[0][:10]
		revTime := r.Key[0][11:15]
		revName := strings.ToLower(strings.TrimSpace(r.Key[1]))
		revEmail := strings.ToLower(strings.TrimSpace(r.Value[0]))
	    revIpAddr := strings.ToLower(strings.TrimSpace(r.Value[1]))
		fmt.Println(revDate, revTime, revName, revEmail, revIpAddr)
		revKeyBuf.WriteString(revDate)
		//revKeyBuf.WriteString(revTime)
		//revKeyBuf.WriteString(revName)
		//revKeyBuf.WriteString(revEmail)
		//revKeyBuf.WriteString(revIpAddr)
		m[revKeyBuf.String()] = 1
		//fmt.Println(revKeyBuf.String())
		revKeyBuf.Reset()
	}

	fmt.Println()
	fmt.Println("======================================")
	fmt.Println()

	var appKeyBuf bytes.Buffer
	cnt := 0
	apps := getApps()
	for _, a := range apps {
		appDate := a.Submitted[:10]
		appTime := a.Submitted[11:15]
        appName := strings.ToLower(strings.TrimSpace(a.LastName))
		appEmail := strings.ToLower(strings.TrimSpace(a.Email))
		appIpAddr := strings.ToLower(strings.TrimSpace(a.IpAddress))
		fmt.Println(appDate, appTime, appName, appEmail, appIpAddr)
		appKeyBuf.WriteString(appDate)
		//appKeyBuf.WriteString(appTime)
		//appKeyBuf.WriteString(appName)
		//appKeyBuf.WriteString(appEmail)
		//appKeyBuf.WriteString(appIpAddr)
		_, ok := m[appKeyBuf.String()]
		if !ok {
			//fmt.Println(a.toString())
			cnt++
		}
		appKeyBuf.Reset()
	}

	fmt.Println()
	fmt.Println("======================================")
	fmt.Println()

	fmt.Println(fmt.Sprintf("There were %v reviews.", len(reviews.Rows)))
	fmt.Println(fmt.Sprintf("There were %v applications.", len(apps)))
	fmt.Println(fmt.Sprintf("There were %v missing reviews.", cnt))

}

func init() {
	flag.StringVarP(&infile, "infile", "i", "", "Input JSON file")
	flag.StringVarP(&dburl, "dburl", "d", "", "Review Database URL")
}

func printUsage() {
	fmt.Printf("Usage: %s [options]\n", os.Args[0])
	fmt.Println("Options:")
	flag.PrintDefaults()
	os.Exit(1)
}
