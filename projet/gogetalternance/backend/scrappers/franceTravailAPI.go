package scrappers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	FTTokenURL  = "https://entreprise.francetravail.fr/connexion/oauth2/access_token?realm=/partenaire"
	FTAPIURL    = "https://api.francetravail.io/partenaire/offresdemploi/v2/offres/search"
	FTScope     = "api_offresdemploiv2 o2dsoffre"
	FTPageSize  = 100
	FTMaxOffres = 300
)

type ftAuthResponse struct {
	AccessToken string `json:"access_token"`
}

type ftOffreResponse struct {
	Resultats []ftOffre `json:"resultats"`
}

type ftOffre struct {
	Id           string `json:"id"`
	Intitule     string `json:"intitule"`
	Description  string `json:"description"`
	TypeContrat  string `json:"typeContrat"`
	DateCreation string `json:"dateCreation"`
	Entreprise   struct {
		Nom string `json:"nom"`
	} `json:"entreprise"`
	LieuTravail struct {
		Libelle string `json:"libelle"`
	} `json:"lieuTravail"`
	OrigineOffre struct {
		UrlOrigine string `json:"urlOrigine"`
	} `json:"origineOffre"`
}

func RunFranceTravailScrapper(keywordsToSearch []string) {
	fmt.Println("\n--- Démarrage du scraping France Travail ---")

	clientID := os.Getenv("FT_CLIENT_ID")
	clientSecret := os.Getenv("FT_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		fmt.Println("ERREUR : FT_CLIENT_ID ou FT_CLIENT_SECRET manquants.")
		fmt.Println("Crée un compte sur francetravail.io pour obtenir tes clés.")
		return
	}

	token, err := getFTToken(clientID, clientSecret)
	if err != nil {
		fmt.Printf("Erreur d'authentification France Travail : %v\n", err)
		return
	}
	fmt.Println("Authentification réussie !")

	var allFTOffers []JobOffer
	seenIDs := make(map[string]bool)

	for _, kw := range keywordsToSearch {
		fmt.Printf("\nRecherche France Travail: %s\n", kw)

		for start := 0; start < FTMaxOffres; start += FTPageSize {
			end := start + FTPageSize - 1
			fmt.Printf("  Requête page %d-%d...\n", start, end)

			offresBrutes := fetchFTOffers(token, kw, start, end)
			if len(offresBrutes) == 0 {
				fmt.Println("  Fin des résultats pour ce mot-clé.")
				break
			}

			newInPage := 0
			for _, o := range offresBrutes {
				if !seenIDs[o.Id] {
					seenIDs[o.Id] = true

					allFTOffers = append(allFTOffers, JobOffer{
						Titre:           o.Intitule,
						Entreprise:      o.Entreprise.Nom,
						Contrat:         o.TypeContrat,
						Localisation:    o.LieuTravail.Libelle,
						Lien:            o.OrigineOffre.UrlOrigine,
						Description:     o.Description,
						DateScraping:    time.Now().Format("2006-01-02"),
						Source:          "France Travail",
						MotCleRecherche: kw,
					})
					newInPage++
				}
			}

			fmt.Printf("  + %d nouvelles offres.\n", newInPage)

			if len(offresBrutes) < FTPageSize {
				break
			}

			time.Sleep(1 * time.Second)
		}
	}

	fmt.Printf("\nPhase France Travail terminée: %d offres trouvées.\n", len(allFTOffers))

	if len(allFTOffers) > 0 {
		fmt.Println("Création du fichier CSV France Travail...")
		os.MkdirAll("data", os.ModePerm)
		saveToCSV("data/offres_francetravail.csv", allFTOffers, true)
		fmt.Println("Données France Travail sauvegardées avec succès !")
	} else {
		fmt.Println("Aucune offre France Travail à sauvegarder.")
	}
}

func getFTToken(clientID, clientSecret string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("scope", FTScope)

	req, err := http.NewRequest("POST", FTTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("statut %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var authResp ftAuthResponse
	json.NewDecoder(resp.Body).Decode(&authResp)
	return authResp.AccessToken, nil
}

func fetchFTOffers(token, keyword string, start, end int) []ftOffre {
	reqURL, _ := url.Parse(FTAPIURL)
	q := reqURL.Query()
	q.Add("motsCles", keyword)
	q.Add("range", fmt.Sprintf("%d-%d", start, end))
	reqURL.RawQuery = q.Encode()

	req, _ := http.NewRequest("GET", reqURL.String(), nil)
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("  Erreur de requête:", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return nil // 204 No Content (Aucun résultat trouvé)
	}

	if resp.StatusCode != 200 && resp.StatusCode != 206 {
		fmt.Printf("  Erreur API: Statut %d\n", resp.StatusCode)
		return nil
	}

	var offreResp ftOffreResponse
	json.NewDecoder(resp.Body).Decode(&offreResp)
	return offreResp.Resultats
}
