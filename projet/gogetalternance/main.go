package main

import (
	"encoding/json"
	"fmt"
	"gogetalternance/backend/scrappers"
	"log"
	"os"
	"sync"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("=== Initialisation de l'application ===")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erreur lors du chargement du fichier .env")
	}

	navigateurLocal := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	u := launcher.New().
		Bin(navigateurLocal).
		Leakless(false).
		Headless(true).
		NoSandbox(true).
		Set("disable-dev-shm-usage").
		MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	motsCles := []string{"Data analyst"}
	typesContrats := []string{"alternance", "stage"}

	var wg sync.WaitGroup

	wg.Add(3)

	var wttjOffres []scrappers.JobOffer
	var indeedOffres []scrappers.JobOffer
	var ftOffres []scrappers.JobOffer

	fmt.Println("Scrapping en parallèle ...")

	go func() {
		defer wg.Done()
		wttjOffres = scrappers.RunWTTJScrapper(browser, motsCles, typesContrats)
	}()

	go func() {
		defer wg.Done()
		indeedOffres = scrappers.RunIndeedScrapper(browser, motsCles, typesContrats)
	}()

	go func() {
		defer wg.Done()
		ftOffres = scrappers.RunFranceTravailScrapper(motsCles, typesContrats)
	}()

	wg.Wait()

	//scrappers.RunWTTJScrapper(browser, motsCles)
	//scrappers.RunIndeedScrapper(browser, motsCles)
	//scrappers.RunFranceTravailScrapper(motsCles)

	fmt.Println("=== Fin de l'application ===")

	var toutesLesOffres []scrappers.JobOffer
	toutesLesOffres = append(toutesLesOffres, wttjOffres...)
	toutesLesOffres = append(toutesLesOffres, indeedOffres...)
	toutesLesOffres = append(toutesLesOffres, ftOffres...)

	fmt.Printf("Total global des offres récupérées : %d\n", len(toutesLesOffres))

	sauvegarderEnJSON(toutesLesOffres)
}

func sauvegarderEnJSON(offres []scrappers.JobOffer) {
	if len(offres) == 0 {
		fmt.Println("Aucune offre à sauvegarder.")
		return
	}

	os.MkdirAll("data", os.ModePerm)

	file, err := os.Create("data/toutes_les_offres.json")
	if err != nil {
		log.Fatalf("Erreur lors de la création du fichier : %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(offres)
	if err != nil {
		log.Fatalf("Erreur lors de l'encodage JSON : %v", err)
	}

	fmt.Println("Toutes les données ont été sauvegardées avec succès !")
}
