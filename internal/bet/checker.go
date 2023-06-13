package bet

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	kRetriesCount          = 3
	kIntervalMsec          = 200
	kFindMaxBetTimeoutMSec = 11000
)

func FindMaxBet(keywords []string) (bet uint32, advertID uint64, err error) {
	start := time.Now()
	totalTimeout := time.After(kFindMaxBetTimeoutMSec * time.Millisecond)
	ch := make(chan BetInfo)
	wg := sync.WaitGroup{}
	wg.Add(len(keywords))
	for _, keyword := range keywords {
		go func(keyword string) {
			defer wg.Done()
			for i := 0; i < kRetriesCount; i++ {
				b, id, err := GetFirstPlaceBet(keyword)
				if err != nil {
					log.Warnf("keyword: %s, error getting first place bet: %v. Retry after %d ms", keyword, err, kIntervalMsec)
					time.Sleep(kIntervalMsec * time.Millisecond)
					continue
				}
				ch <- BetInfo{
					AdvertID: id,
					Bet:      b,
				}
				return
			}
			log.Errorf("could not handle keyword: %s after %d retries", keyword, kRetriesCount)
		}(keyword)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	results := []BetInfo{}
loop:
	for {
		select {
		case res, ok := <-ch:
			if !ok {
				log.Infof("all keywords handled")
				break loop
			}
			results = append(results, res)
		case <-totalTimeout:
			log.Warnf("timeout %d ms was reached. %d of %d keywords handled", kFindMaxBetTimeoutMSec, len(results), len(keywords))
			break loop
		}
	}
	log.Infof("%d keywords handled in %s\n", len(results), time.Since(start))
	maxBet := max(results)
	log.Infof("max bet is: %d", maxBet.Bet)
	return maxBet.Bet, maxBet.AdvertID, nil
}

func max(bets []BetInfo) BetInfo {
	m := BetInfo{}
	for _, b := range bets {
		if b.Bet > m.Bet {
			m = b
		}
	}
	return m
}
