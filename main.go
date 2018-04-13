package main

import (
	"encoding/json"
	"fmt"
	flag "github.com/ogier/pflag"
	"io/ioutil"
	"os"
)

type Page struct {
	Submitted string `json:"Submitted"`
	LastName  string `json:"Last-Name"`
	Email     string `json:"from_email"`
	IpAddr    string `json:"ip_address"`
}

func (p Page) toString() string {
	return toJson(p)
}

func toJson(p interface{}) string {
	bytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return string(bytes)
}

func getApps() []Page {
	raw, err := ioutil.ReadFile(infile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c []Page
	json.Unmarshal(raw, &c)
	return c
}

var (
	infile string
)

func main() {
	flag.Parse()

	if flag.NFlag() == 0 {
		printUsage()
	}

	fmt.Printf("Loading: %s\n", infile)

	apps := getApps()
	for _, p := range apps {
		fmt.Println(p.toString())
	}

}

func init() {
	flag.StringVarP(&infile, "infile", "i", "", "Input JSON file")
}

func printUsage() {
	fmt.Printf("Usage: %s [options]\n", os.Args[0])
	fmt.Println("Options:")
	flag.PrintDefaults()
	os.Exit(1)
}
