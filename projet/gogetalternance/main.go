package main

import (
	"fmt"
	"gogetalternance/backend/scrappers"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func main() {
	fmt.Println("=== Initialisation de l'application ===")

	// 1. Initialisation GLOBALE du navigateur (Visible pour passer Cloudflare)
	navigateurLocal := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	u := launcher.New().
		Bin(navigateurLocal).
		Leakless(false).
		Headless(false). // Garde false pour l'instant pour voir ce qu'il se passe
		NoSandbox(true).
		Set("disable-dev-shm-usage").
		MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose() // On fermera le navigateur à la toute fin

	// 2. Lancement des scrapers en leur "prêtant" le navigateur
	scrappers.RunWTTJScrapper(browser)

	fmt.Println("---------------------------------------------------")

	scrappers.RunIndeedScrapper(browser)

	fmt.Println("=== Fin de l'application ===")
}
