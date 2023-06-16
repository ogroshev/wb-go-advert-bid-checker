package handler

import (
	"encoding/json"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"

	"gitlab.com/wb-dynamics/wb-go-advert-bid-checker/internal/wbcatalog"
)

type AdvertBetInfoRequest struct {
	AdvertID                  uint64   `json:"advert_id"`
	SubjectID                 uint64   `json:"subject_id"`
	TargetPlace               uint32   `json:"target_place"`
	ValidateSubjectIDPriority bool     `json:"validate_subject_id_priority"`
	Keywords                  []string `json:"keywords"`
}

type AdvertWarning struct {
	Keyword           string `json:"keyword"`
	PrioritySubjectID uint64 `json:"priority_subject_id"`
}

type BetInfo struct {
	Keyword   string `json:"keyword"`
	SubjectID uint64 `json:"subject_id"`
	AdvertID  uint64 `json:"advert_id"`
	Place     uint32 `json:"place"`
	Bet       uint64 `json:"bet"`
}

type Bets struct {
	TargetPlace BetInfo `json:"target_place"`
	NextPlace   BetInfo `json:"next_place"`
	MyPlace     uint32  `json:"my_place"`
}

type AdvertBetInfoResponse struct {
	Bets     Bets            `json:"bets"`
	Warnings []AdvertWarning `json:"warnings"`
}

func HandleAdvertBetInfo(w http.ResponseWriter, r *http.Request) {
	log.Infof("request received")
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Could not read request body: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Debugf("request body: %s", string(b))
	req := AdvertBetInfoRequest{}
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

	var respBody AdvertBetInfoResponse
	if req.ValidateSubjectIDPriority {
		respBody.Warnings = findWarnings(keywordsData, req.SubjectID)
	}
	respBody.Bets.TargetPlace = findBetForPlace(keywordsData, req.SubjectID, req.TargetPlace)
	respBody.Bets.NextPlace = findBetForPlace(keywordsData, req.SubjectID, req.TargetPlace+1)
	respBody.Bets.MyPlace = findMaxPlaceByAdvertID(keywordsData, uint32(req.AdvertID))

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

func findWarnings(keywordsData []wbcatalog.KeywordCollectedData, subjectID uint64) []AdvertWarning {
	var warnings []AdvertWarning
	for _, kwData := range keywordsData {
		if uint64(kwData.CatalogAds.PrioritySubjects[0]) != subjectID {
			warnings = append(warnings, AdvertWarning{
				Keyword:           kwData.Keyword,
				PrioritySubjectID: uint64(kwData.CatalogAds.PrioritySubjects[0]),
			})
		}
	}
	return warnings
}

func findBetForPlace(keywordsData []wbcatalog.KeywordCollectedData, subjectID uint64,
	targetPlace uint32) (maxBet BetInfo) {
	for _, kwData := range keywordsData {
		var elemIdxWithSubjectID uint32 = 0
		for idx, advertInfo := range kwData.CatalogAds.Adverts {
			// в массиве Adverts порядок определяет место. Берем targetPlace - 1 элемент, который совпадает с нашим subjectId
			if advertInfo.Subject == subjectID {
				if targetPlace-1 == elemIdxWithSubjectID && advertInfo.Cpm > maxBet.Bet {
					maxBet = BetInfo{
						Keyword:   kwData.Keyword,
						SubjectID: subjectID,
						AdvertID:  advertInfo.AdvertId,
						Place:     uint32(idx + 1),
						Bet:       advertInfo.Cpm,
					}
					break
				}
				elemIdxWithSubjectID++

			}
		}
	}
	return maxBet
}

func findMaxPlaceByAdvertID(keywordsData []wbcatalog.KeywordCollectedData, advertID uint32) (maxPlace uint32) {
	maxPlace = 0
	for _, kwData := range keywordsData {
		for idx, advertInfo := range kwData.CatalogAds.Adverts {
			if uint32(advertInfo.AdvertId) == advertID {
				if uint32(idx + 1)  > maxPlace {
					maxPlace = uint32(idx + 1)
				}
			}
		}
	}
	return maxPlace
}
