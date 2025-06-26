# 📖 Lectionary Slide Generator

This Go project automatically generates Reveal.js-compatible HTML slides for the Catholic Mass readings, psalms, acclamations, and Gospel for any given day.

---

## ✨ Features

- ✅ Fetches **Spanish daily Mass readings** directly from [usccb.org](https://bible.usccb.org)
- ✅ Parses:
  - First Reading
  - Psalm
  - Second Reading (if present)
  - Acclamation before the Gospel
  - Gospel Reading
- ✅ Splits long readings into multiple slide-sized chunks
- ✅ Separates Psalms and Acclamations cleanly by **Response (R.) markers**
- ✅ Outputs clean Reveal.js-compatible HTML slide decks (`slides.html`)
- ✅ Exports a JSON file (`lectionary.json`) for debugging or other uses

---

## 🏗️ Project Structure

```
.
├── build_slides.go          # Main program (entry point)
├── templates/               # HTML templates for each slide section
│   ├── base.html
│   ├── title.html
│   ├── first_reading.html
│   ├── psalm.html
│   ├── second_reading.html
│   ├── acclamation.html
│   └── gospel_reading.html
├── slides.html              # ✅ Generated slide deck
├── lectionary.json          # ✅ JSON debug export
└── README.md
```

---

## How It Works

1. Fetches the Mass reading page from **usccb.org** for the given date.
2. Parses each Lectionary section (Readings, Psalm, Acclamation, Gospel).
3. Splits content into multiple slide-friendly chunks.
4. Renders all sections into a single Reveal.js slide deck.

---

## Running the Program

```bash
go run build_slides.go
```

This will generate:

- `slides.html` → Your finished slide deck
- `lectionary.json` → Debug/export of parsed lectionary data

---

## Customizing the Date

Right now, the target date URL is hardcoded inside `main()`:

```go
url := "https://bible.usccb.org/es/bible/lecturas/062525.cfm"
```

You can change the URL to the desired day from the USCCB site.

---

## Upcoming Improvements

- [ ] **English language support**
- [ ] Dynamic date generation (get today’s date automatically)
- [ ] CLI flags for choosing output file and date
- [ ] Support for other liturgical seasons or special Mass types
- [ ] Improved styling and layout customization for slides

---

## Requirements

- Go 1.18+
- Internet connection (to fetch live readings from USCCB)

---

## License

MIT License – For non-commercial, educational, and liturgical use.
