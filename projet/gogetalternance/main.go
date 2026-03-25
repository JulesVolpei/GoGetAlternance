package main

import (
	"fmt"
	"gogetalternance/backend/scrappers"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func main() {
	fmt.Println("=== Initialisation de l'application ===")

	navigateurLocal := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	u := launcher.New().
		Bin(navigateurLocal).
		Leakless(false).
		Headless(false).
		NoSandbox(true).
		Set("disable-dev-shm-usage").
		MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	motsCles := []string{"développeur python"}

	//scrappers.RunWTTJScrapper(browser, motsCles)
	//scrappers.RunIndeedScrapper(browser, motsCles)
	scrappers.RunFranceTravailScrapper(motsCles)

	fmt.Println("=== Fin de l'application ===")
}
