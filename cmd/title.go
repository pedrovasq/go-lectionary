package main

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var titleCaser = cases.Title(language.Spanish)

var bookTypes = map[string]string{
	// Old Testament Books
	"génesis":        "book",
	"éxodo":          "book",
	"levítico":       "book",
	"números":        "book",
	"deuteronomio":   "book",
	"isaías":         "prophet",
	"jeremías":       "prophet",
	"ezequiel":       "prophet",
	"daniel":         "prophet",
	"salmo":          "psalm",

	// New Testament Letters
	"romanos":        "letter",
	"corintios":      "letter",
	"gálatas":        "letter",
	"efesios":        "letter",
	"filipenses":     "letter",
	"colosenses":     "letter",
	"tesalonicenses": "letter",
	"timoteo":        "letter",
	"tito":           "letter",
	"hebreos":        "letter",
	"santiago":       "letter",
	"pedro":          "letter",
	"juan":           "letter",
	"judas":          "letter",

	// Gospels
	"mateo":          "gospel",
	"marcos":         "gospel",
	"lucas":          "gospel",
	"juan evangelio": "gospel",  // Helps avoid conflict with John's letters

	// Other Books
	"hechos":         "acts",
	"apocalipsis":    "apocalypse",
}

var typeFormats = map[string]string{
	"book":       "Lectura del libro del %s",
	"prophet":    "Lectura del libro del profeta %s",
	"letter":     "Lectura de la carta del apóstol san Pablo a los %s",
	"gospel":     "Lectura del santo Evangelio según san %s",
	"psalm":      "Salmo Responsorial",
	"acts":       "Lectura del libro de los Hechos de los Apóstoles",
	"apocalypse": "Lectura del libro del Apocalipsis del apóstol san Juan",
}

func formatTitle(address string) string {
	addressLower := strings.ToLower(address)

	for book, bookType := range bookTypes {
		if strings.Contains(addressLower, book) {
			format := typeFormats[bookType]
			bookTitle := titleCaser.String(book) 
			return fmt.Sprintf(format, bookTitle)
		}
	}

	return "Lectura"
}


