package scrappers

import (
	"fmt"
	"gogetalternance/backend/models"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

const BaseURL = "https://www.welcometothejungle.com"

var urlContractMap = map[string]string{
	"alternance": "apprenticeship",
	"stage":      "internship",
	"cdi":        "full_time",
	"cdd":        "temporary",
	"freelance":  "freelance",
}

const PagesPerKeyword = 5

func RunWTTJScrapper(browser *rod.Browser, keywordsToSearch []string, contractTypes []string) []models.JobOffer {
	var allResults []models.JobOffer
	seenLinks := make(map[string]bool)

	for _, kw := range keywordsToSearch {
		for page := 1; page <= PagesPerKeyword; page++ {
			offres := scrapeListingPage(browser, page, kw, contractTypes)
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

	fmt.Printf("[WTTJ] %d offres trouvées\n", len(allResults))

	return allResults
}

func scrapeListingPage(browser *rod.Browser, pageNum int, keyword string, contractTypes []string) []models.JobOffer {
	kwURL := strings.ReplaceAll(keyword, " ", "+")

	contractQuery := ""
	for i, c := range contractTypes {
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
		return nil
	}

	time.Sleep(2 * time.Second)

	elements := page.MustElements(`li[data-testid="search-results-list-item-wrapper"]`)

	var offres []models.JobOffer
	for _, el := range elements {
		lienEl, err := el.Element(`a[aria-label]`)
		if err != nil {
			continue
		}

		ariaLabel, _ := lienEl.Attribute("aria-label")
		href, _ := lienEl.Attribute("href")
		carteText, _ := el.Text()

		titre, contratAria := parseAriaLabel(*ariaLabel)

		contratFinal := detectContrat(titre, contratAria, contractTypes)
		if contratFinal == "Non précisé" {
			continue
		}

		lienFinal := *href
		if strings.HasPrefix(lienFinal, "/") {
			lienFinal = BaseURL + lienFinal
		}

		offres = append(offres, models.JobOffer{
			Titre:           titre,
			Entreprise:      parseEntreprise(*href),
			Contrat:         contratFinal,
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

func detectContrat(titre, contratAria string, contractTypes []string) string {
	if contratAria != "" {
		return contratAria
	}
	titreLower := strings.ToLower(titre)
	for _, c := range contractTypes {
		if strings.Contains(titreLower, strings.ToLower(c)) {
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
