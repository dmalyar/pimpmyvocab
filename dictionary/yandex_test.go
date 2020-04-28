package dictionary

import (
	"github.com/dmalyar/pimpmyvocab/domain"
	"github.com/dmalyar/pimpmyvocab/mock"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

const (
	positiveJson = `{
  "def": [
    {
      "text": "Positive",
      "pos": "adjective",
      "ts": "ˈpɒzɪtɪv",
      "tr": [
        {
          "text": "положительный"
        },
        {
          "text": "уверенный"
        }
      ]
    },
    {
      "text": "Positive",
      "pos": "noun",
      "ts": "ˈpɒzɪtɪv",
      "tr": [
        {
          "text": "позитив"
        }
      ]
    },
    {
      "text": "Positive"
    },
    {
      "text": "not positive",
      "pos": "noun"
    }
  ]
}`
	diffTextJson = `{
  "def": [
    {
      "text": "something",
      "pos": "noun",
      "tr": [
        {
          "text": "что-то"
        }
      ]
    }
  ]
}`
	emptyJson  = `{}`
	brokenJson = `{{}`
)

func TestYandex_GetVocabEntryByText(t *testing.T) {
	testCases := []struct {
		text          string
		expectedEntry *domain.VocabEntry
		expectErr     bool
	}{
		{
			text: "Positive",
			expectedEntry: &domain.VocabEntry{
				ID:            0,
				Text:          "Positive",
				Transcription: "ˈpɒzɪtɪv",
				Translations: []*domain.Translation{
					{
						ID:       0,
						Text:     "положительный",
						Class:    "adjective",
						Position: 0,
					},
					{
						ID:       0,
						Text:     "уверенный",
						Class:    "adjective",
						Position: 1,
					},
					{
						ID:       0,
						Text:     "позитив",
						Class:    "noun",
						Position: 2,
					},
				},
			},
		},
		{
			text: "Different text",
		},
		{
			text: "Empty json",
		},
		{
			text:      "Broken json",
			expectErr: true,
		},
		{
			text:      "Code not 200",
			expectErr: true,
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/Positive":
			rw.Write([]byte(positiveJson))
		case "/Different text":
			rw.Write([]byte(diffTextJson))
		case "/Empty json":
			rw.Write([]byte(emptyJson))
		case "/Broken json":
			rw.Write([]byte(brokenJson))
		case "/Code not 200":
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer mockServer.Close()

	ya := NewYandexDict(
		&mock.Logger{},
		http.DefaultClient,
		mockServer.URL+"/",
	)
	for _, c := range testCases {
		t.Run(c.text, func(t *testing.T) {
			entry, err := ya.GetVocabEntryByText(c.text)
			if c.expectErr == false && err != nil {
				t.Errorf("Expected no error, but got %s", err)
			}
			if c.expectErr && err == nil {
				t.Errorf("Expected error, but got nothing")
			}
			if c.expectedEntry == nil && entry != nil {
				t.Errorf("Nil entry expected")
			}
			if c.expectedEntry != nil && entry == nil {
				t.Errorf("Not nil entry expected")
			}
			if c.expectedEntry != nil && entry != nil && !reflect.DeepEqual(c.expectedEntry, entry) {
				t.Errorf("Expected vocab:%+v;Actual:%+v", c.expectedEntry, entry)
			}
		})
	}
}
