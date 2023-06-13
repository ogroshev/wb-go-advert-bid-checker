package wbcatalog


type CatalogAdsResponse struct {
	Pages []struct {
		Positions []int `json:"positions"`
		Page int `json:"page"`
		Count int `json:"count"`
	} `json:"pages"`
	PrioritySubjects []int `json:"prioritySubjects"`
	Adverts []struct {
		Code string `json:"code"`
		AdvertId uint64 `json:"advertId"`
		Id uint64 `json:"id"`
		Cpm uint64 `json:"cpm"`
		Subject uint64 `json:"subject"`
	}
	MinCpm uint64 `json:"minCpm"`
}

type KeywordCollectedData struct {
	Keyword string
	CatalogAds CatalogAdsResponse
}

