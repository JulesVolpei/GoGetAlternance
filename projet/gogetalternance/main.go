package main

import (
	"fmt"
	"net/http"
	"sync"

	"gogetalternance/backend/models"
	"gogetalternance/backend/scrappers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Fichier .env non trouvé")
	}

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowMethods = []string{"POST", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type"}
	r.Use(cors.New(config))

	r.POST("/api/cherche", func(c *gin.Context) {
		var req models.SearchRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON invalide"})
			return
		}

		fmt.Println("-------------------------------------------------")
		fmt.Printf("RECHERCHE : %v | CONTRATS : %v | PLATEFORMES : %v\n", req.Keywords, req.ContractTypes, req.Platforms)
		fmt.Println("-------------------------------------------------")

		var wg sync.WaitGroup
		var mu sync.Mutex
		var allOffers []models.JobOffer

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

		if contains(req.Platforms, "FranceTravail") {
			wg.Add(1)
			go func() {
				defer wg.Done()
				offers := scrappers.RunFranceTravailScrapper(req.Keywords, req.ContractTypes)
				mu.Lock()
				allOffers = append(allOffers, offers...)
				mu.Unlock()
			}()
		}

		if contains(req.Platforms, "Welcome to the Jungle") {
			wg.Add(1)
			go func() {
				defer wg.Done()
				offers := scrappers.RunWTTJScrapper(browser, req.Keywords, req.ContractTypes)
				mu.Lock()
				allOffers = append(allOffers, offers...)
				mu.Unlock()
			}()
		}

		if contains(req.Platforms, "HelloWork") {
			wg.Add(1)
			go func() {
				defer wg.Done()

				offers := scrappers.RunHelloWorkScrapper(browser, req.Keywords, req.ContractTypes)

				mu.Lock()
				allOffers = append(allOffers, offers...)
				mu.Unlock()
			}()
		}

		if contains(req.Platforms, "Indeed") {
			wg.Add(1)
			go func() {
				defer wg.Done()

				offers := scrappers.RunIndeedScrapper(browser, req.Keywords, req.ContractTypes)

				mu.Lock()
				allOffers = append(allOffers, offers...)
				mu.Unlock()
			}()
		}

		wg.Wait()

		fmt.Printf("RECHERCHE TERMINÉE : %d offres trouvées au total.\n", len(allOffers))

		c.JSON(http.StatusOK, allOffers)
	})

	fmt.Println("Serveur Go lancé sur http://localhost:8080")
	r.Run(":8080")
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
