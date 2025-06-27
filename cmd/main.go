package main

import (
    "html/template"
    "fmt"
    "os"
	"time"
	"flag"
    "encoding/json"
)

func main() {
	// ✅ Command line flag
	dateFlag := flag.String("date", "", "Optional: specify date as MMDDYY (e.g., 062625 for June 26, 2025)")
	flag.Parse()

	var url string
	if *dateFlag != "" {
		// ✅ If user provides --date
		url = fmt.Sprintf("https://bible.usccb.org/es/bible/lecturas/%s.cfm", *dateFlag)
	} else {
		// ✅ Default: use today's date
		today := time.Now()
		url = fmt.Sprintf("https://bible.usccb.org/es/bible/lecturas/%02d%02d%02d.cfm", today.Month(), today.Day(), today.Year()%100)
	}

	fmt.Println("Fetching URL:", url)

	// ✅ Proceed as normal
	lectionary, err := fetchLectionary(url)
	if err != nil {
		panic(err)
	}

	// ✅ Save JSON for debugging
	jsonData, err := json.MarshalIndent(lectionary, "", "  ")
	if err != nil {
		panic(err)
	}
	os.WriteFile("lectionary.json", jsonData, 0644)
	fmt.Println("Lectionary saved to lectionary.json")

	// ✅ Parse templates
	tmpl, err := template.ParseFiles(
		"ui/templates/base.html",
		"ui/templates/blank_slide.html",
		"ui/templates/title.html",
		"ui/templates/first_reading.html",
		"ui/templates/psalm.html",
		"ui/templates/second_reading.html",
		"ui/templates/acclamation.html",
		"ui/templates/gospel_reading.html",
	)
	if err != nil {
		panic(err)
	}

	// ✅ Create output HTML file
	outputFile, err := os.Create("slides.html")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	// ✅ Execute main template, passing in the full Lectionary struct
	err = tmpl.ExecuteTemplate(outputFile, "base.html", lectionary)
	if err != nil {
		panic(err)
	}

	fmt.Println("Slides generated as slides.html")
}

