package scrappers

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
)

const HelloWorkBaseURL = "https://www.hellowork.com"
const HelloWorkPagesPerKeyword = 3

func getHelloWorkContractParams(contractTypes []string) string {
	var params []string
	for _, c := range contractTypes {
		cLower := strings.ToLower(c)
		switch cLower {
		case "alternance":
			params = append(params, "c=Alternance")
		case "stage":
			params = append(params, "c=Stage")
		case "cdi":
			params = append(params, "c=CDI")
		case "cdd":
			params = append(params, "c=CDD")
		case "freelance":
			params = append(params, "c=Freelance")
		}
	}

	if len(params) > 0 {
		return "&" + strings.Join(params, "&")
	}
	return ""
}

func RunHelloWorkScrapper(browser *rod.Browser, keywordsToSearch []string, contractTypes []string) []JobOffer {
	var allResults []JobOffer
	seenLinks := make(map[string]bool)

	for _, kw := range keywordsToSearch {
		for pageNum := 1; pageNum <= HelloWorkPagesPerKeyword; pageNum++ {
			offres := scrapeHelloWorkPage(browser, pageNum, kw, contractTypes)

			if len(offres) == 0 {
				break
			}

			for _, off := range offres {
				if !seenLinks[off.Lien] {
					seenLinks[off.Lien] = true
					allResults = append(allResults, off)
				}
			}

			if len(offres) < 15 {
				break
			}

			randomSleep(2, 5)
		}
	}

	fmt.Printf("[HelloWork] %d offres trouvées\n", len(allResults))

	return allResults
}

func scrapeHelloWorkPage(browser *rod.Browser, pageNum int, keyword string, contractTypes []string) []JobOffer {
	kwURL := strings.ReplaceAll(keyword, " ", "+")
	contractQuery := getHelloWorkContractParams(contractTypes)

	url := fmt.Sprintf("%s/fr-fr/emploi/recherche.html?k=%s%s&st=relevance", HelloWorkBaseURL, kwURL, contractQuery)

	if pageNum > 1 {
		url += fmt.Sprintf("&p=%d", pageNum)
	}

	page := stealth.MustPage(browser)
	defer page.MustClose()

	page.MustNavigate(url)

	err := page.Timeout(15*time.Second).WaitElementsMoreThan(`ul[aria-label="liste des offres"] > li`, 0)
	if err != nil {
		return nil
	}

	time.Sleep(2 * time.Second)

	elements := page.MustElements(`ul[aria-label="liste des offres"] > li`)

	var offres []JobOffer
	for _, el := range elements {
		linkEl, err := el.Element(`a[data-cy="offerTitle"]`)
		if err != nil {
			continue
		}

		hrefPtr, _ := linkEl.Attribute("href")
		lienFinal := ""
		if hrefPtr != nil {
			lienFinal = *hrefPtr
			if strings.HasPrefix(lienFinal, "/") {
				lienFinal = HelloWorkBaseURL + lienFinal
			}
		}

		titre := "Titre inconnu"
		entreprise := "Non précisée"

		pElements, err := linkEl.Elements(`h3 p`)
		if err == nil && len(pElements) >= 2 {
			titre, _ = pElements[0].Text()
			entreprise, _ = pElements[1].Text()
		} else if err == nil && len(pElements) == 1 {
			titre, _ = pElements[0].Text()
		}

		localisation := "France"
		if locEl, err := el.Element(`[data-cy="localisationCard"]`); err == nil {
			localisation, _ = locEl.Text()
		}

		contrat := "Non précisé"
		if contractEl, err := el.Element(`[data-cy="contractCard"]`); err == nil {
			contrat, _ = contractEl.Text()
		}

		offres = append(offres, JobOffer{
			Titre:           strings.TrimSpace(titre),
			Entreprise:      strings.TrimSpace(entreprise),
			Contrat:         strings.TrimSpace(contrat),
			Localisation:    strings.TrimSpace(localisation),
			Lien:            lienFinal,
			DateScraping:    time.Now().Format("2006-01-02"),
			Source:          "HelloWork",
			MotCleRecherche: keyword,
		})
	}

	return offres
}
