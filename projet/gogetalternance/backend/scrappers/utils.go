package scrappers

import (
	"encoding/csv"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func createCSVWithHeaders(path string) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatal("Impossible de créer le CSV:", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"titre", "entreprise", "contrat", "localisation", "lien", "date_scraping", "source", "mot_cle_recherche", "description"})
}

func appendSingleToCSV(path string, offre JobOffer) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Erreur ouverture CSV append:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	desc := strings.ReplaceAll(offre.Description, "\n", " ")

	writer.Write([]string{
		offre.Titre, offre.Entreprise, offre.Contrat, offre.Localisation,
		offre.Lien, offre.DateScraping, offre.Source, offre.MotCleRecherche, desc,
	})
}

func saveToCSV(path string, offres []JobOffer, includeDesc bool) {
	createCSVWithHeaders(path)
	for _, o := range offres {
		appendSingleToCSV(path, o)
	}
}

// --- Autres Utilitaires ---

// C'est aussi le bon endroit pour mettre ton randomSleep !
func randomSleep(min, max float64) {
	sleepTime := min + rand.Float64()*(max-min)
	time.Sleep(time.Duration(sleepTime * float64(time.Second)))
}
