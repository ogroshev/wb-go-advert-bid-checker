package handler

import (
	"encoding/json"
	"net/http"
	"io"

	log "github.com/sirupsen/logrus"

	"gitlab.com/wb-dynamics/wb-go-advert-bid-checker/internal/bet"
)

type FindMaxBetRequest struct {
	Keywords []string `json:"keywords"`
}

func HandleFindMaxBet(w http.ResponseWriter, r *http.Request) {
	log.Debugf("request received")
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
	

	maxBet, advertID, err :=  bet.FindMaxBet(req.Keywords)
	if err != nil {
		log.Errorf("Could not find max bet: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respBody := bet.BetInfo{
		AdvertID: advertID,
		Bet: maxBet,
	}
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