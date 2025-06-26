package main

import (
	"html"
    "html/template"
    "fmt"
    "net/http"
    "os"
    "encoding/json"
    "strings"
    // "time"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

const maxChars = 215

type Lection struct {
    Title     string   `json:"title"`
    Reference string   `json:"reference,omitempty"`
    Chunks    []string `json:"chunks"`
}

type Lectionary struct {
    DayTitle      string    `json:"day_title"`
    FirstReading  Lection 	`json:"first_reading"`
    Psalm         Lection 	`json:"psalm"`
    SecondReading *Lection 	`json:"second_reading,omitempty"`
    Acclamation   Lection 	`json:"acclamation"`
    Gospel        Lection 	`json:"gospel"`
}

func fetchLectionary(url string) (Lectionary, error) {
    res, err := http.Get(url)
    if err != nil {
        return Lectionary{}, err
    }
    defer res.Body.Close()

    if res.StatusCode != 200 {
        return Lectionary{}, fmt.Errorf("status code error: %d", res.StatusCode)
    }

    doc, err := goquery.NewDocumentFromReader(res.Body)
    if err != nil {
        return Lectionary{}, err
    }

    dayTitle := strings.TrimSpace(doc.Find(".b-lectionary h2").First().Text())

    firstReading, err := parseLection(doc, "Primera Lectura")
    if err != nil {
        return Lectionary{}, err
    }

    psalm, err := parseLection(doc, "Salmo Responsorial")
    if err != nil {
        return Lectionary{}, err
    }

    var secondReading *Lection
    if doc.Find("h3.name:contains('Segunda Lectura')").Length() > 0 {
        second, err := parseLection(doc, "Segunda Lectura")
        if err != nil {
            return Lectionary{}, err
        }
        secondReading = &second
    }

    acclamation, err := parseLection(doc, "Aclamación antes del Evangelio")
    if err != nil {
        return Lectionary{}, err
    }

    gospel, err := parseLection(doc, "Evangelio")
    if err != nil {
        return Lectionary{}, err
    }

    return Lectionary{
        DayTitle:      	dayTitle,
        FirstReading:  	firstReading,
        Psalm:         	psalm,
        SecondReading: 	secondReading,
		Acclamation: 	acclamation,
        Gospel:        	gospel,
    }, nil
}

func parseLection(doc *goquery.Document, targetHeading string) (Lection, error) {
	var lection Lection
	found := false

	doc.Find(".b-verse .innerblock").EachWithBreak(func(i int, s *goquery.Selection) bool {
		header := strings.TrimSpace(strings.ToLower(s.Find(".content-header h3.name").Text()))
		target := strings.TrimSpace(strings.ToLower(targetHeading))

		if header == target {
			found = true

			lection.Title = strings.TrimSpace(s.Find(".content-header h3.name").Text())
			lection.Reference = strings.TrimSpace(s.Find(".content-header .address").Text())

			// Get raw HTML inside <p> tag
			html, err := s.Find(".content-body").Html()
			if err != nil {
				return true // Skip this block if error
			}

			lection.Chunks = lectionChunkify(lection.Title, html)

			return false
		}
		return true
	})

	if !found {
		return lection, fmt.Errorf("lection with heading '%s' not found", targetHeading)
	}

	return lection, nil
}

func lectionChunkify(title string, rawHtml string) []string {
	title = strings.ToLower(title)

	if strings.Contains(title, "primera lectura") || strings.Contains(title, "segunda lectura") {
		return splitByPunctuationAndLength(rawHtml)
	} else if strings.Contains(title, "salmo") || strings.Contains(title, "aclamación") {
		return splitByPsalmResponse(rawHtml)
	} else if title == "evangelio" {
		return splitByPunctuationAndLength(rawHtml)
	}

	// Fallback for unknown sections
	return []string{stripHTML(rawHtml)}
}

func splitByPunctuationAndLength(rawHtml string) []string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader("<div>" + rawHtml + "</div>"))
	if err != nil {
		return []string{stripHTML(rawHtml)} // Fallback
	}

	var text string
	doc.Find("p").Each(func(i int, p *goquery.Selection) {
		paragraph := strings.TrimSpace(p.Text())
		if paragraph != "" {
			text += paragraph + " "
		}
	})

	// Split into chunks based on punctuation + maxChars
	var chunks []string
	sentences := strings.FieldsFunc(text, func(r rune) bool {
		return r == '.' || r == ';' || r == '?' || r == '!'
	})

	var currentChunk strings.Builder
	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if sentence == "" {
			continue
		}
		sentence += ". "

		if currentChunk.Len()+len(sentence) > maxChars && currentChunk.Len() > 0 {
			chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
			currentChunk.Reset()
		}

		currentChunk.WriteString(sentence)
	}

	if currentChunk.Len() > 0 {
		chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
	}

	return chunks
}

func splitByPsalmResponse(rawHtml string) []string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader("<div>" + rawHtml + "</div>"))
	if err != nil {
		return []string{stripHTML(rawHtml)} // Fallback
	}

	var chunks []string
	var currentVerse strings.Builder
	inResponse := false

	doc.Find("p").Contents().Each(func(i int, node *goquery.Selection) {
		nodeName := goquery.NodeName(node)

		if nodeName == "#text" {
			text := strings.TrimSpace(node.Text())
			if text == "" {
				return
			}

			if strings.HasPrefix(text, "R.") {
				// ✅ If there was a verse being built, finalize it first
				if currentVerse.Len() > 0 {
					chunks = append(chunks, strings.TrimSpace(currentVerse.String()))
					currentVerse.Reset()
				}

				// ✅ Start new response chunk with the "R." text
				currentVerse.WriteString(text)
				inResponse = true
			} else {
				// ✅ Normal verse text
				currentVerse.WriteString(text + " ")
			}

		} else if nodeName == "strong" && inResponse {
			// ✅ Capture strong text after R.
			responseText := strings.TrimSpace(node.Text())
			if responseText != "" {
				currentVerse.WriteString(" " + responseText)
			}

		} else if nodeName == "br" {
			currentVerse.WriteString("\n")

			// ✅ After a response chunk finishes at line break, save it
			if inResponse && currentVerse.Len() > 0 {
				chunks = append(chunks, strings.TrimSpace(currentVerse.String()))
				currentVerse.Reset()
				inResponse = false
			}
		}
	})

	// ✅ Save final verse at end of loop
	if currentVerse.Len() > 0 {
		chunks = append(chunks, strings.TrimSpace(currentVerse.String()))
	}

	return chunks
}

func stripHTML(input string) string {
    // Quick and dirty HTML tag stripper (you can use regex or a proper HTML parser for more complex cases)
    re := regexp.MustCompile(`<.*?>`)
    return html.UnescapeString(re.ReplaceAllString(input, ""))
}

func main() {
	url := "https://bible.usccb.org/es/bible/lecturas/062625.cfm"

	// ✅ Fetch Lectionary from USCCB
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
		"templates/base.html",
		"templates/blank_slide.html",
		"templates/title.html",
		"templates/first_reading.html",
		"templates/psalm.html",
		"templates/second_reading.html",
		"templates/acclamation.html",
		"templates/gospel_reading.html",
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

