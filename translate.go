package gtfs

import "io"

// A Translation is a single translation between an original string and another
// language.
//
// For example, the following Translation:
//
// 	Translation{
// 		ID:          "station-001",
// 		Language:    "en",
// 		Translation: "City Center",
// 	}
//
// would translate "station-001" to "City Center" in English.
//
// Fields correspond directly to columns in translations.txt.
type Translation struct {
	ID          string
	Language    string
	Translation string
}

// Translate translates a source string into the specified language.
//
// If no translation is available, the original source string is returned.
func (g *GTFS) Translate(sourceStr, lang string) string {
	trs, ok := g.translationsByID[sourceStr]
	if !ok {
		return sourceStr
	}

	tr, ok := trs[lang]
	if !ok {
		return sourceStr
	}

	return tr.Translation
}

var translationFields = map[string]bool{
	"trans_id":    true,
	"lang":        true,
	"translation": true,
}

func (g *GTFS) processTranslations(r io.Reader) error {
	res, err := readCSVWithHeadings(r, translationFields, g.strictMode)
	if err != nil {
		return err
	}

	g.translationsByID = map[string]map[string]*Translation{}

	for _, row := range res {
		t := &Translation{
			ID:          row["trans_id"],
			Language:    row["lang"],
			Translation: row["translation"],
		}

		g.Translations = append(g.Translations, t)

		_, ok := g.translationsByID[t.ID]
		if !ok {
			g.translationsByID[t.ID] = map[string]*Translation{}
		}

		g.translationsByID[t.ID][t.Language] = t
	}

	return nil
}
