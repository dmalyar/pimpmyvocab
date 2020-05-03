package dictionary

import (
	"encoding/json"
	"fmt"
	"github.com/dmalyar/pimpmyvocab/domain"
	"github.com/dmalyar/pimpmyvocab/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const URL = "https://dictionary.yandex.net/api/v1/dicservice.json/lookup?key=%s&lang=en-ru&text="

type Yandex struct {
	logger log.Logger
	client *http.Client
	url    string
}

type Response struct {
	Def []Definition
}

type Definition struct {
	Text string
	Pos  string
	Ts   string
	Tr   []Translation
}

type Translation struct {
	Text string
}

func NewYandexDict(logger log.Logger, client *http.Client, url string) *Yandex {
	return &Yandex{
		logger: logger,
		client: client,
		url:    url,
	}
}

// GetVocabEntryByText returns an entry found in the Yandex.Dictionary service.
// Returns nil if entry was not found.
func (y *Yandex) GetVocabEntryByText(text string) (*domain.VocabEntry, error) {
	logger := y.logger.WithField("text", text)
	logger.Debug("Getting vocab entry from yandex dictionary")
	req, err := http.NewRequest(http.MethodGet, y.url+url.PathEscape(text), nil)
	if err != nil {
		return nil, fmt.Errorf("creating http request: %s", err)
	}
	res, err := y.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling yandex.dictionary: %s", err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading http response body: %s", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("yandex dictionary respond with code %v and body %s", res.StatusCode, string(body))
	}
	logger.Debugf("Yandex dictionary respond with body %v", string(body))
	parsedRes, err := parseResponse(body)
	if err != nil {
		return nil, fmt.Errorf("parsing http response: %s", err)
	}
	return convertToVocabEntry(text, parsedRes), nil
}

func convertToVocabEntry(text string, res *Response) *domain.VocabEntry {
	if len(res.Def) == 0 {
		return nil
	}
	entry := new(domain.VocabEntry)
	entry.Text = text
	position := 0
	for _, d := range res.Def {
		if strings.ToLower(d.Text) != strings.ToLower(text) {
			continue
		}
		if entry.Transcription == "" {
			entry.Transcription = d.Ts
		}
		class := d.Pos
		if class == "" {
			continue
		}
		for _, t := range d.Tr {
			translation := new(domain.Translation)
			entry.Translations = append(entry.Translations, translation)
			translation.Text = t.Text
			translation.Class = class
			translation.Position = position
			position++
		}
	}
	if len(entry.Translations) == 0 {
		return nil
	}
	entry.MainTranslation = entry.Translations[0].Text
	return entry
}

func parseResponse(res []byte) (*Response, error) {
	parsedRes := new(Response)
	err := json.Unmarshal(res, parsedRes)
	if err != nil {
		return nil, err
	}
	return parsedRes, nil
}

func (y *Yandex) GetVocabEntryByID(_ int) (*domain.VocabEntry, error) {
	return nil, fmt.Errorf("GetVocabEntryByID is not supported")
}
