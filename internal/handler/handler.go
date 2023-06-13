package handler

import (
	"encoding/json"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"

	"gitlab.com/wb-dynamics/wb-go-advert-bid-checker/internal/wbcatalog"
)

type FindMaxBetRequest struct {
	SubjectID uint64   `json:"subject_id"`
	Keywords  []string `json:"keywords"`
}

type AdvertWarning struct {
	Keyword           string `json:"keyword"`
	PrioritySubjectID uint64 `json:"priority_subject_id"`
}

type MaxBet struct {
	Keyword   string `json:"keyword"`
	SubjectID uint64 `json:"subject_id"`
	AdvertID  uint64 `json:"advert_id"`
	Place     uint32 `json:"place"`
	Bet       uint64 `json:"bet"`
}
type FindMaxBetResponse struct {
	SubjectMaxBet MaxBet          `json:"subject_max_bet"`
	Warnings      []AdvertWarning `json:"warnings"`
}

func HandleFindMaxBet(w http.ResponseWriter, r *http.Request) {
	log.Infof("request received")
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Could not read request body: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Debugf("request body: %s", string(b))
	req := FindMaxBetRequest{}
	err = json.Unmarshal(b, &req)
	if err != nil {
		log.Errorf("Could not parse request json: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	keywordsData, err := wbcatalog.CollectAdvertCompaniesInfo(req.Keywords)
	if err != nil {
		log.Errorf("Could not collect advert companies info: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	respBody := findMaxBet(keywordsData, req.SubjectID)

	respByte, err := json.Marshal(respBody)
	if err != nil {
		log.Errorf("Could not marshal response: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respByte)
	log.Debugf("send response: %s", string(respByte))
}

func findMaxBet(keywordsData []wbcatalog.KeywordCollectedData, subjectID uint64) (resp FindMaxBetResponse) {
	var maxBet MaxBet
	for _, kwData := range keywordsData {
		if subjectID != uint64(kwData.CatalogAds.PrioritySubjects[0]) {
			resp.Warnings = append(resp.Warnings, AdvertWarning{
				Keyword:           kwData.Keyword,
				PrioritySubjectID: uint64(kwData.CatalogAds.PrioritySubjects[0]),
			})
		}
		for idx, advertInfo := range kwData.CatalogAds.Adverts {
			// в массиве Adverts порядок определяет место. Берем первый элемент, который совпадает с нашим subjectId
			if advertInfo.Subject == subjectID {
				if advertInfo.Cpm > maxBet.Bet {
					maxBet = MaxBet{
						Keyword:   kwData.Keyword,
						SubjectID: subjectID,
						AdvertID:  advertInfo.AdvertId,
						Place:     uint32(idx + 1),
						Bet:       advertInfo.Cpm,
					}
				}
				break
			}
		}
	}
	resp.SubjectMaxBet = maxBet
	return resp
}
