package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
)

// --- Define constants and variables---//
const filename = "topWines2020.csv"

// var wineType string
var tpl *template.Template

// Create a struct for storing CSV lines and annotate it with JSON struct field tags
type Wine struct {
	Rank        int    `json:"rank"`
	Name        string `json:"name"`
	Country     string `json:"country"`
	Region      string `json:"region"`
	Style       string `json:"style"`
	Description string `json:"description"`
	Grapes      string `json:"grapes"`
	PairWith    string `json:"pairWith"`
	AvgPrice    int    `json:"avgPriceUSD"`
	Vintage     int    `json:"vintage"`
	Score       int    `json:"score"`
}

type Profile struct {
	Header   string
	Grapes   string
	Name     string
	Country  string
	Region   string
	Style    string
	PairWith string
	AvgPrice string
	Vintage  string
	Score    string
	Desc     string
}

func main() {
	mux := http.NewServeMux()
	// --- Serving static css ---//
	mux.HandleFunc("/", index)
	mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServer((http.Dir("./served")))))
	mux.HandleFunc("/search", searchProcess)
	// mux.HandleFunc("/results", resultProfile)
	http.ListenAndServe(":8080", mux)
}

// Set up Endpoints
// func HandleRequests() {
// 	http.HandleFunc("/", index)
// 	// --- Serving static css ---//
// 	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer((http.Dir("./served")))))
// 	http.HandleFunc("/search", searchResult)
// 	http.ListenAndServe(":8080", nil)
// }

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}
func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func searchProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	records, err := OpenFile(filename)
	if err != nil {
		log.Fatal((err))
	}
	wines := WineList(records)
	//---Invoke ParseForm before reading form values ---//
	r.ParseForm()
	// wineType := r.FormValue("winename")
	wineType := strings.Title(r.FormValue("winename"))
	results := SearchFor(wineType, wines)
	WineInfo(results)

	tpl.ExecuteTemplate(w, "results.gohtml", nil)
}

// func resultProfile(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintln(w, "Successful")
// 	tpl.ExecuteTemplate(w, "results.gohtml", nil)
// }

// --- Loading the Wine List ---//
func OpenFile(filename string) ([][]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed open file: %w ", err)
	}
	//Close the file at the end of the program
	defer f.Close()

	// Read CSV file using csv.Reader
	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed read file: %s\n", err)
	}
	return records, nil
}

func WineList(records [][]string) []Wine {
	var wines []Wine
	for _, rec := range records {
		rank, _ := strconv.Atoi(rec[0])
		price, _ := strconv.Atoi(rec[8])
		vintage, _ := strconv.Atoi(rec[9])
		score, _ := strconv.Atoi(rec[10])

		w := Wine{
			Rank:        rank,
			Name:        rec[1],
			Country:     rec[2],
			Region:      rec[3],
			Style:       rec[4],
			Description: rec[5],
			Grapes:      rec[6],
			PairWith:    rec[7],
			AvgPrice:    price,
			Vintage:     vintage,
			Score:       score,
		}
		wines = append(wines, w)
		// fmt.Println(wines)
	}
	return wines
}

// Function performs Search by the Wine type in the Top 100 Wines list
func SearchFor(wineType string, wines []Wine) []Wine {
	match := []Wine{}
	for _, w := range wines {
		if w.Grapes == wineType {
			match = append(match, w)
		}
	}
	if len(match) == 1 {
		fmt.Println("\n Success,", len(match), "item is found")
	} else if len(match) > 1 {
		fmt.Println("\n Success,", len(match), "items are found")
	} else {
		fmt.Printf("Sorry, no %s wine is found \n", wineType)
	}
	return match
}

func WineInfo(results []Wine) []Profile {
	var profile []Profile
	for i := range results {
		p := Profile{
			Grapes:   results[i].Grapes,
			Name:     results[i].Name,
			Country:  results[i].Country,
			Region:   results[i].Region,
			Style:    results[i].Style,
			PairWith: results[i].PairWith,
			// AvgPrice: results[i].AvgPrice,
			// Vintage:  results[i].Vintage,
			// Score:    results[i].Score,
			Desc: results[i].Description,
		}
		profile = append(profile, p)

		fmt.Println("\n Wine:", p.Grapes,
			"\n Name:", p.Name,
			"\n Region:", p.Region,
			"\n Country:", p.Country,
			"\n Style:", p.Style,
			"\n Pairs with:", p.PairWith,
			"\n Review:", p.Desc,
		)
	}
	return profile
}
