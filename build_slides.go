package main

import (
    "html/template"
    "fmt"
    "net/http"
    "os"
    "encoding/json"
    "strings"
    // "time"

	"github.com/PuerkitoBio/goquery"
)

type Reading struct {
	Title		string		`json:"title"`
	Reference	string		`json:"reference"`
	Paragraphs	[]string	`json:"paragraphs"`
}

func fetchFirstReading(url string) (Reading, error) {
    res, err := http.Get(url)
    if err != nil {
        return Reading{}, err
    }
    defer res.Body.Close()

    if res.StatusCode != 200 {
        return Reading{}, fmt.Errorf("status code error: %d", res.StatusCode)
    }

    doc, err := goquery.NewDocumentFromReader(res.Body)
    if err != nil {
        return Reading{}, err
    }

    var reading Reading

    doc.Find(".innerblock").Each(func(i int, s *goquery.Selection) {
        header := s.Find(".content-header h3.name").Text()
		if strings.Contains(header, "Primera Lectura") {
			addressText := strings.TrimSpace(s.Find(".address").Text())
			addressFields := strings.Fields(addressText)
			if len(addressFields) > 0 {
				reading.Title = addressFields[0]
			}

            reading.Reference = addressText

            s.Find(".content-body p").Each(func(j int, p *goquery.Selection) {
                paragraph := strings.TrimSpace(p.Text())
                if paragraph != "" {
                    reading.Paragraphs = append(reading.Paragraphs, paragraph)
                }
            })
        }
    })

    return reading, nil
}

func main() {
	url := "https://bible.usccb.org/es/bible/lecturas/062525.cfm"
    reading, err := fetchFirstReading(url)
    if err != nil {
        panic(err)
    }

    jsonData, err := json.MarshalIndent(reading, "", "  ")
    if err != nil {
        panic(err)
    }

    os.WriteFile("reading_first.json", jsonData, 0644)
    fmt.Println("First reading saved to reading_first.json")

    tmpl, err := template.ParseFiles(
        "templates/base.html",
		"templates/title.html",
        "templates/first_reading.html",
        "templates/second_reading.html",
        "templates/gospel_reading.html",
    )
    if err != nil {
        panic(err)
    }

	data := struct {
		FirstReadingTitle     	string
		FirstReadingReference 	string
		FirstReadingParagraphs 	[]string
		SecondReadingTitle    	string
		SecondReadingReference 	string
		SecondReadingParagraphs []string
		GospelTitle           	string
		GospelReference         string
		GospelParagraphs      	[]string
	}{
		FirstReadingTitle:     		"Lectura del libro del Génesis",
		FirstReadingReference: 		"Génesis 15, 1-12. 17-18",
		FirstReadingParagraphs: 	[]string{"Primera parte del texto...", "Segunda parte del texto..."},
		SecondReadingTitle:    		"Lectura de la carta del apóstol san Pablo a los Romanos",
		SecondReadingReference: 	"Romans 12:16",
		SecondReadingParagraphs:	[]string{"Contenido de la segunda lectura..."},
		GospelTitle:           		"Lectura del santo Evangelio según san Mateo",
		GospelReference: 			"Mateo 10:2",
		GospelParagraphs:      		[]string{"Parte 1 del Evangelio...", "Parte 2 del Evangelio..."},
	}

    // Generate output file
    outputFile, err := os.Create("slides.html")
    if err != nil {
        panic(err)
    }
    defer outputFile.Close()

    tmpl.ExecuteTemplate(outputFile, "base.html", data)
}

