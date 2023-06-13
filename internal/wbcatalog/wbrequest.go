package wbcatalog

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
)

func RequestCatalogAds(keyword string) (catalogAdsResponse *CatalogAdsResponse, err error) {


	keywordUrlEncoded := url.QueryEscape(keyword)
	u := fmt.Sprintf("https://catalog-ads.wildberries.ru/api/v5/search?keyword=%s", keywordUrlEncoded)
	log.Debugf("keyword '%s' - Request: %s", keyword, u)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, fmt.Errorf("can't create request: %v", err)
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Origin", "https://www.wildberries.ru/")
	req.Header.Set("Host", "catalog-ads.wildberries.ru")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.0 Safari/605.1.15")
	req.Header.Set("Accept-Language", "en-us")
	req.Header.Set("Referer", fmt.Sprintf("https://www.wildberries.ru/catalog/0/search.aspx?search=%s", keywordUrlEncoded))
	req.Header.Set("Connection", "keep-alive")

	//TODO: сделать создание клиента один раз на все вызовы
	client := createHTTPClient()

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't send request: %v", err)
	}
	defer resp.Body.Close()

	log.Debugf("keyword: '%s' - Status code: %d", keyword, resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read response body: %v", err)
	}
	log.Tracef("keyword: '%s' - Response: %s", keyword, string(body))

	err = json.Unmarshal(body, &catalogAdsResponse)
	if err != nil {
		return nil, fmt.Errorf("can't parse json: %v", err)
	}
	
	return catalogAdsResponse, nil
}

func createHTTPClient() *http.Client {
	proxyParam := viper.GetString("Proxy")
	if proxyParam != "" {
		log.Debugf("Using proxy: %s", proxyParam)
		proxyUrl, err := url.Parse(proxyParam)
		if err != nil {
			log.Fatalf("can't parse proxy url: %v, error: %v", proxyParam, err)
		}

		client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
		return client
	}
	return &http.Client{}
}