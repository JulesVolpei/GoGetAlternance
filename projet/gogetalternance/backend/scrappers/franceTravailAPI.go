package scrappers

import (
	"encoding/json"
	"fmt"
	"gogetalternance/backend/models"
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

var contractTypeMap = map[string]string{
	"alternance": "E1",
	"stage":      "E2",
}

func RunFranceTravailScrapper(keywordsToSearch []string, contractTypes []string) []models.JobOffer {
	clientID := os.Getenv("FT_CLIENT_ID")
	clientSecret := os.Getenv("FT_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		fmt.Println("[France Travail] ERREUR : FT_CLIENT_ID ou FT_CLIENT_SECRET manquants.")
		return nil
	}

	token, err := getFTToken(clientID, clientSecret)
	if err != nil {
		fmt.Printf("[France Travail] Erreur d'authentification : %v\n", err)
		return nil
	}

	var allFTOffers []models.JobOffer
	seenIDs := make(map[string]bool)

	for _, kw := range keywordsToSearch {
		for start := 0; start < FTMaxOffres; start += FTPageSize {
			end := start + FTPageSize - 1

			offresBrutes := fetchFTOffers(token, kw, start, end, contractTypes)

			if len(offresBrutes) == 0 {
				break
			}

			for _, o := range offresBrutes {
				if !seenIDs[o.Id] {
					seenIDs[o.Id] = true

					allFTOffers = append(allFTOffers, models.JobOffer{
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
				}
			}

			if len(offresBrutes) < FTPageSize {
				break
			}
			time.Sleep(1 * time.Second)
		}
	}

	// Affichage unique et clair pour le parallélisme
	fmt.Printf("[France Travail] %d offres trouvées\n", len(allFTOffers))

	return allFTOffers
}

func fetchFTOffers(token, keyword string, start, end int, contractTypes []string) []ftOffre {
	reqURL, _ := url.Parse(FTAPIURL)
	q := reqURL.Query()
	q.Add("motsCles", keyword)
	q.Add("range", fmt.Sprintf("%d-%d", start, end))

	for _, ct := range contractTypes {
		if mapped, ok := contractTypeMap[strings.ToLower(ct)]; ok {
			q.Add("natureContrat", mapped)
		}
	}

	reqURL.RawQuery = q.Encode()

	req, _ := http.NewRequest("GET", reqURL.String(), nil)
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return nil
	}

	var offreResp ftOffreResponse
	json.NewDecoder(resp.Body).Decode(&offreResp)
	return offreResp.Resultats
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

	var authResp ftAuthResponse
	json.NewDecoder(resp.Body).Decode(&authResp)
	return authResp.AccessToken, nil
}
