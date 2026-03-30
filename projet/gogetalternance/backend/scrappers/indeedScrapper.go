package scrappers

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
)

const IndeedBaseURL = "https://fr.indeed.com"

func RunIndeedScrapper(browser *rod.Browser, keywordsToSearch []string, contractTypes []string) []JobOffer {
	fmt.Println("Démarrage du scraping Indeed...")
	var allIndeedResults []JobOffer

	for _, kw := range keywordsToSearch {
		fmt.Printf("\n--- Recherche Indeed: %s ---\n", kw)

		for page := 1; page <= 1; page++ {
			offres := ScrapeIndeedListingPage(browser, page, kw, contractTypes)
			if len(offres) == 0 {
				break
			}
			allIndeedResults = append(allIndeedResults, offres...)
			randomSleep(3, 6)
		}
	}

	fmt.Printf("\nPhase 1 Indeed terminée: %d offres trouvées.\n", len(allIndeedResults))
	return allIndeedResults
}

func ScrapeIndeedListingPage(browser *rod.Browser, pageNum int, keyword string, contractTypes []string) []JobOffer {
	start := (pageNum - 1) * 10
	kwURL := strings.ReplaceAll(keyword, " ", "+")

	scParam := "0kf%3Aattr%28CPAHG%7CQADT5%7CVDTG7%252COR%29%3B"
	url := fmt.Sprintf("%s/jobs?q=%s&sc=%s&start=%d", IndeedBaseURL, kwURL, scParam, start)

	page := stealth.MustPage(browser)
	defer page.MustClose()

	page.MustNavigate(url)

	err := page.Timeout(20*time.Second).WaitElementsMoreThan(`.job_seen_beacon`, 0)
	if err != nil {
		fmt.Printf("    Page %d: timeout (Probablement un Captcha Cloudflare !)\n", pageNum)
		return nil
	}

	time.Sleep(2 * time.Second)

	elements := page.MustElements(`.job_seen_beacon`)
	fmt.Printf("    Page %d: %d offres trouvées\n", pageNum, len(elements))

	var capitalizedContracts []string
	for _, c := range contractTypes {
		capitalizedContracts = append(capitalizedContracts, strings.Title(c))
	}
	contratFormatte := strings.Join(capitalizedContracts, "/")

	var offres []JobOffer
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

		offres = append(offres, JobOffer{
			Titre:           strings.TrimSpace(titre),
			Entreprise:      strings.TrimSpace(entreprise),
			Contrat:         contratFormatte, // Utilisation de la variable dynamique ici
			Localisation:    strings.TrimSpace(localisation),
			Lien:            lienFinal,
			DateScraping:    time.Now().Format("2006-01-02"),
			Source:          "Indeed",
			MotCleRecherche: keyword,
		})
	}

	return offres
}
