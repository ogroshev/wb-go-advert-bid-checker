package wbcatalog

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

func CollectAdvertCompaniesInfo(keywords []string) (keywordsAdvert []KeywordCollectedData, err error) {
	start := time.Now()
	totalTimeout := time.After(kFindMaxBetTimeoutMSec * time.Millisecond)
	ch := make(chan KeywordCollectedData)
	wg := sync.WaitGroup{}
	wg.Add(len(keywords))
	for _, keyword := range keywords {
		go func(keyword string) {
			defer wg.Done()
			for i := 0; i < kRetriesCount; i++ {
				resp, err := RequestCatalogAds(keyword)
				if err != nil {
					log.Warnf("keyword: %s, error getting first place bet: %v. Retry after %d ms", keyword, err, kIntervalMsec)
					time.Sleep(kIntervalMsec * time.Millisecond)
					continue
				}
				ch <- KeywordCollectedData{
					Keyword: keyword,
					CatalogAds: *resp,
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

loop:
	for {
		select {
		case res, ok := <-ch:
			if !ok {
				log.Debugf("all keywords handled")
				break loop
			}
			keywordsAdvert = append(keywordsAdvert, res)
		case <-totalTimeout:
			log.Warnf("timeout %d ms was reached. %d of %d keywords handled", kFindMaxBetTimeoutMSec, len(keywordsAdvert), len(keywords))
			break loop
		}
	}
	log.Infof("%d keywords handled in %s\n", len(keywordsAdvert), time.Since(start))

	return
}

