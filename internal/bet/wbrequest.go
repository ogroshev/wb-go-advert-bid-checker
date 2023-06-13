package bet

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

func GetFirstPlaceBet(keyword string) (bet uint32, advertID uint64, err error) {
	keywordUrlEncoded := url.QueryEscape(keyword)
	u := fmt.Sprintf("https://catalog-ads.wildberries.ru/api/v5/search?keyword=%s", keywordUrlEncoded)
	log.Debugf("Request: %s", u)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("can't create request: %v", err)
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Origin", "https://www.wildberries.ru/")
	req.Header.Set("Host", "catalog-ads.wildberries.ru")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.0 Safari/605.1.15")
	req.Header.Set("Accept-Language", "en-us")
	req.Header.Set("Referer", fmt.Sprintf("https://www.wildberries.ru/catalog/0/search.aspx?search=%s", keywordUrlEncoded))
	req.Header.Set("Connection", "keep-alive")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, fmt.Errorf("can't send request: %v", err)
	}
	defer resp.Body.Close()

	log.Debugf("Status code: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("can't read response body: %v", err)
	}
	log.Tracef("Response: %s", string(body))

	var catalogAdsResponse CatalogAdsResponse
	err = json.Unmarshal(body, &catalogAdsResponse)
	if err != nil {
		return 0, 0, fmt.Errorf("can't parse json: %v", err)
	}
	b, id := findFirstPlaceBet(catalogAdsResponse)
	return b, id, nil
}

func findFirstPlaceBet(catalogAdsResponse CatalogAdsResponse) (bet uint32, advertID uint64) {
	aa := catalogAdsResponse.Adverts
	if len(aa) == 0 {
		return 0, 0
	}
	return uint32(aa[0].Cpm), uint64(aa[0].AdvertId)
}