package scrappers

import (
	"fmt"
	"gogetalternance/backend/models"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
)

const IndeedBaseURL = "https://fr.indeed.com"
const IndeedPagesPerKeyword = 3

func RunIndeedScrapper(browser *rod.Browser, keywordsToSearch []string, contractTypes []string) []models.JobOffer {
	var allIndeedResults []models.JobOffer
	seenLinks := make(map[string]bool)

	for _, kw := range keywordsToSearch {
		for page := 1; page <= IndeedPagesPerKeyword; page++ {
			offres := ScrapeIndeedListingPage(browser, page, kw, contractTypes)

			if len(offres) == 0 {
				break
			}

			for _, off := range offres {
				if !seenLinks[off.Lien] {
					seenLinks[off.Lien] = true
					allIndeedResults = append(allIndeedResults, off)
				}
			}

			if len(offres) < 10 {
				break
			}

			randomSleep(3, 6)
		}
	}

	fmt.Printf("[Indeed] %d offres trouvées\n", len(allIndeedResults))
	return allIndeedResults
}

func ScrapeIndeedListingPage(browser *rod.Browser, pageNum int, keyword string, contractTypes []string) []models.JobOffer {
	start := (pageNum - 1) * 10
	kwURL := strings.ReplaceAll(keyword, " ", "+")

	scParam := "0kf%3Aattr%28CPAHG%7CQADT5%7CVDTG7%252COR%29%3B"
	url := fmt.Sprintf("%s/jobs?q=%s&sc=%s&start=%d", IndeedBaseURL, kwURL, scParam, start)

	page := stealth.MustPage(browser)
	defer page.MustClose()

	page.MustNavigate(url)

	err := page.Timeout(20*time.Second).WaitElementsMoreThan(`.job_seen_beacon`, 0)
	if err != nil {
		return nil
	}

	time.Sleep(2 * time.Second)

	elements := page.MustElements(`.job_seen_beacon`)

	var capitalizedContracts []string
	for _, c := range contractTypes {
		capitalizedContracts = append(capitalizedContracts, strings.Title(c))
	}
	contratFormatte := strings.Join(capitalizedContracts, "/")

	var offres []models.JobOffer
	for _, el := range elements {
		titleEl, err := el.Element(`a.jcs-JobTitle`)
		if err != nil {
			continue
		}

		titre, _ := titleEl.Text()
		hrefPtr, _ := titleEl.Attribute("href")

		lienFinal := ""
		if hrefPtr != nil {
			lienFinal = *hrefPtr
			if strings.HasPrefix(lienFinal, "/") {
				lienFinal = IndeedBaseURL + lienFinal
			}
		}

		entreprise := "Non précisée"
		if compEl, err := el.Element(`[data-testid="company-name"]`); err == nil {
			entreprise, _ = compEl.Text()
		}

		localisation := "France"
		if locEl, err := el.Element(`[data-testid="text-location"]`); err == nil {
			localisation, _ = locEl.Text()
		}

		offres = append(offres, models.JobOffer{
			Titre:           strings.TrimSpace(titre),
			Entreprise:      strings.TrimSpace(entreprise),
			Contrat:         contratFormatte,
			Localisation:    strings.TrimSpace(localisation),
			Lien:            lienFinal,
			DateScraping:    time.Now().Format("2006-01-02"),
			Source:          "Indeed",
			MotCleRecherche: keyword,
		})
	}

	return offres
}
