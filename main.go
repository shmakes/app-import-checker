package main

import (
	"encoding/json"
	"fmt"
	flag "github.com/ogier/pflag"
	"io/ioutil"
	"net/http"
	"os"
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

	apps := getApps()
	for _, a := range apps {
		appDate := a.Submitted[:10]
		appTime := a.Submitted[11:16]
		appName := a.LastName
		appEmail := a.Email
		appIpAddr := a.IpAddress
		fmt.Println(appDate, appTime, appName, appEmail, appIpAddr)
	}

    fmt.Println()
	fmt.Println("======================================")
    fmt.Println()

	reviews := getReviews()
	for _, r := range reviews.Rows {
		revDate := r.Key[0][:10]
		revTime := r.Key[0][11:16]
		revName := r.Key[1]
		revEmail := r.Value[0]
		revIpAddr := r.Value[1]
		fmt.Println(revDate, revTime, revName, revEmail, revIpAddr)
	}

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
