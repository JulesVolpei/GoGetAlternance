package scrappers

import (
	"fmt"
	"os"
	_ "path/filepath"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

const BaseURL = "https://www.welcometothejungle.com"

var Keywords = []string{
	"data",
}

var urlContractMap = map[string]string{
	"alternance": "apprenticeship",
	"stage":      "internship",
	"cdi":        "full_time",
	"cdd":        "temporary",
	"freelance":  "freelance",
}

const PagesPerKeyword = 5

var Contrats = []string{"alternance", "stage"}

type JobOffer struct {
	Titre           string
	Entreprise      string
	Contrat         string
	Localisation    string
	Lien            string
	DateScraping    string
	Source          string
	MotCleRecherche string
	Description     string
}

func RunWTTJScrapper(browser *rod.Browser) {
	fmt.Println("Démarrage du scraping WTTJ...")

	var allResults []JobOffer
	seenLinks := make(map[string]bool)

	for _, kw := range Keywords {
		fmt.Printf("\n--- Recherche: %s ---\n", kw)
		for page := 1; page <= PagesPerKeyword; page++ {
			offres := scrapeListingPage(browser, page, kw)
			if len(offres) == 0 {
				break
			}

			for _, off := range offres {
				if !seenLinks[off.Lien] {
					seenLinks[off.Lien] = true
					allResults = append(allResults, off)
				}
			}
			randomSleep(2, 4)
		}
	}

	fmt.Printf("\nPhase 1 terminée: %d offres uniques trouvées.\n", len(allResults))

	// Sauvegarde intermédiaire
	os.MkdirAll("data", os.ModePerm)
	saveToCSV("data/offres_wttj_listings_only.csv", allResults, false)
	fmt.Println("Sauvegarde intermédiaire effectuée.")
}

func scrapeListingPage(browser *rod.Browser, pageNum int, keyword string) []JobOffer {
	kwURL := strings.ReplaceAll(keyword, " ", "+")

	contractQuery := ""
	for i, c := range Contrats {
		if val, ok := urlContractMap[strings.ToLower(c)]; ok {
			contractQuery += fmt.Sprintf("&refinementList%%5Bcontract_type%%5D%%5B%d%%5D=%s", i, val)
		}
	}

	url := fmt.Sprintf("%s/fr/jobs?refinementList%%5Boffices.country_code%%5D%%5B0%%5D=FR%s&query=%s&page=%d",
		BaseURL, contractQuery, kwURL, pageNum)

	page := browser.MustPage(url)
	defer page.MustClose()

	err := page.Timeout(10*time.Second).WaitElementsMoreThan(`li[data-testid="search-results-list-item-wrapper"]`, 0)
	if err != nil {
		fmt.Printf("    Page %d: timeout ou vide\n", pageNum)
		return nil
	}

	time.Sleep(2 * time.Second)

	elements := page.MustElements(`li[data-testid="search-results-list-item-wrapper"]`)
	fmt.Printf("    Page %d: %d offres trouvées\n", pageNum, len(elements))

	var offres []JobOffer
	for _, el := range elements {
		lienEl, err := el.Element(`a[aria-label]`)
		if err != nil {
			continue
		}

		ariaLabel, _ := lienEl.Attribute("aria-label")
		href, _ := lienEl.Attribute("href")
		carteText, _ := el.Text()

		titre, contratAria := parseAriaLabel(*ariaLabel)

		contratFinal := detectContrat(titre, contratAria)
		if contratFinal == "Non précisé" {
			continue // On passe directement à l'offre suivante
		}

		lienFinal := *href
		if strings.HasPrefix(lienFinal, "/") {
			lienFinal = BaseURL + lienFinal
		}

		offres = append(offres, JobOffer{
			Titre:           titre,
			Entreprise:      parseEntreprise(*href),
			Contrat:         contratFinal, // On utilise la variable vérifiée
			Localisation:    parseLocalisationDynamique(carteText),
			Lien:            lienFinal,
			DateScraping:    time.Now().Format("2006-01-02"),
			Source:          "WTTJ",
			MotCleRecherche: keyword,
		})
	}
	return offres
}

func scrapeJobDetails(browser *rod.Browser, url string) string {
	page := browser.MustPage(url)
	defer page.MustClose()

	_, err := page.Timeout(8 * time.Second).Element("main")
	if err != nil {
		return ""
	}
	time.Sleep(1 * time.Second)

	mainContent, err := page.Element("main")
	if err != nil {
		return ""
	}

	text, _ := mainContent.Text()
	return text
}

func parseAriaLabel(ariaLabel string) (string, string) {
	texte := strings.ReplaceAll(ariaLabel, "Consultez l'offre ", "")
	parts := strings.SplitN(texte, " | ", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}
	return strings.TrimSpace(texte), ""
}

func detectContrat(titre, contratAria string) string {
	if contratAria != "" {
		return contratAria
	}
	titreLower := strings.ToLower(titre)
	for _, c := range Contrats {
		if strings.Contains(titreLower, c) {
			return strings.Title(c)
		}
	}
	return "Non précisé"
}

func parseEntreprise(href string) string {
	parties := strings.Split(href, "/")
	for i, part := range parties {
		if part == "companies" && i+1 < len(parties) {
			nom := strings.ReplaceAll(parties[i+1], "-", " ")
			return strings.Title(nom)
		}
	}
	return "N/A"
}

func parseLocalisationDynamique(texteCarte string) string {
	lines := strings.Split(texteCarte, "\n")
	var cleanLines []string

	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if trimmed != "" {
			cleanLines = append(cleanLines, trimmed)
		}
	}

	contratsExacts := []string{"stage", "alternance", "cdi", "cdd", "freelance", "apprentissage"}

	for i, line := range cleanLines {
		lineLower := strings.ToLower(line)

		isContrat := false
		for _, c := range contratsExacts {
			if lineLower == c {
				isContrat = true
				break
			}
		}

		if isContrat && i+1 < len(cleanLines) {
			ville := cleanLines[i+1]

			villeLower := strings.ToLower(ville)
			if strings.Contains(villeLower, "télétravail") || strings.Contains(villeLower, "remote") || strings.Contains(villeLower, "jours") {
				return "Non précisée"
			}

			return strings.Title(ville)
		}
	}

	return "France"
}

//func randomSleep(min, max float64) {
//	sleepTime := min + rand.Float64()*(max-min)
//	time.Sleep(time.Duration(sleepTime * float64(time.Second)))
//}
//
//func createCSVWithHeaders(path string) {
//	file, err := os.Create(path)
//	if err != nil {
//		log.Fatal("Impossible de créer le CSV:", err)
//	}
//	defer file.Close()
//
//	writer := csv.NewWriter(file)
//	defer writer.Flush()
//	writer.Write([]string{"titre", "entreprise", "contrat", "localisation", "lien", "date_scraping", "source", "mot_cle_recherche", "description"})
//}
//
//func appendSingleToCSV(path string, offre JobOffer) {
//	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
//	if err != nil {
//		log.Println("Erreur ouverture CSV append:", err)
//		return
//	}
//	defer file.Close()
//
//	writer := csv.NewWriter(file)
//	defer writer.Flush()
//
//	desc := strings.ReplaceAll(offre.Description, "\n", " ")
//
//	writer.Write([]string{
//		offre.Titre, offre.Entreprise, offre.Contrat, offre.Localisation,
//		offre.Lien, offre.DateScraping, offre.Source, offre.MotCleRecherche, desc,
//	})
//}
//
//func saveToCSV(path string, offres []JobOffer, includeDesc bool) {
//	createCSVWithHeaders(path)
//	for _, o := range offres {
//		appendSingleToCSV(path, o)
//	}
//}
