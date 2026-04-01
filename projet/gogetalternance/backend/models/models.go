package models

type SearchRequest struct {
	Keywords      []string `json:"keywords"`
	ContractTypes []string `json:"contractTypes"`
	Platforms     []string `json:"platforms"`
}

type JobOffer struct {
	Titre           string `json:"title"`
	Entreprise      string `json:"company"`
	Contrat         string `json:"contract"`
	Localisation    string `json:"location"`
	Lien            string `json:"url"`
	DateScraping    string `json:"scrapeDate"`
	Source          string `json:"source"`
	MotCleRecherche string `json:"keyword"`
	Description     string `json:"description"`
}
