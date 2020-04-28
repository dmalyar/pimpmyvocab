package service

import (
	"errors"
	"fmt"
	"github.com/dmalyar/pimpmyvocab/domain"
	"github.com/dmalyar/pimpmyvocab/mock"
	"reflect"
	"testing"
)

func TestVocabWithLocalRepo_CreateVocabTable(t *testing.T) {
	testCases := []struct {
		name                      string
		userID                    int
		expectedVocab             *domain.Vocab
		expectGetVocabByUserIDInv bool
		expectAddVocabInv         bool
		expectErr                 bool
	}{
		{
			name:   "Positive",
			userID: 1,
			expectedVocab: &domain.Vocab{
				ID:     1,
				UserID: 1,
			},
			expectGetVocabByUserIDInv: true,
			expectAddVocabInv:         true,
		},
		{
			name:                      "GetVocabByUserID returns error",
			userID:                    2,
			expectErr:                 true,
			expectGetVocabByUserIDInv: true,
		},
		{
			name:                      "GetVocabByUserID returns vocab",
			userID:                    3,
			expectGetVocabByUserIDInv: true,
		},
		{
			name:                      "AddVocab returns error",
			userID:                    4,
			expectErr:                 true,
			expectGetVocabByUserIDInv: true,
			expectAddVocabInv:         true,
		},
	}

	mockedRepo := &mock.VocabRepo{
		GetVocabByUserIDFn: func(userID int) (*domain.Vocab, error) {
			switch userID {
			case 1:
				return nil, nil
			case 2:
				return nil, errors.New("err")
			case 3:
				return &domain.Vocab{}, nil
			case 4:
				return nil, nil
			}
			return nil, nil
		},
		AddVocabFn: func(vocab *domain.Vocab) (*domain.Vocab, error) {
			switch vocab.UserID {
			case 1:
				return &domain.Vocab{
					ID:     1,
					UserID: 1,
				}, nil
			case 4:
				return nil, errors.New("err")
			}
			return nil, nil
		},
	}

	vocabService := NewVocabWithLocalRepo(mock.Logger{}, mockedRepo, &mock.VocabEntryService{})
	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			vocab, err := vocabService.CreateVocab(c.userID)
			if c.expectErr == false && err != nil {
				t.Errorf("Expected no error, but got %s", err)
			}
			if c.expectErr && err == nil {
				t.Errorf("Expected error, but got nothing")
			}
			if c.expectAddVocabInv != mockedRepo.AddVocabInvoked {
				t.Errorf("Actual invocation of AddVocab(%v) doesn't match expectations", mockedRepo.AddVocabInvoked)
			}
			if c.expectGetVocabByUserIDInv != mockedRepo.GetVocabByUserIDInvoked {
				t.Errorf("Actual invocation of GetVocabByUserID(%v) doesn't match expectations", mockedRepo.GetVocabByUserIDInvoked)
			}
			if c.expectedVocab == nil && vocab != nil {
				t.Errorf("Nil vocab expected")
			}
			if c.expectedVocab != nil && vocab == nil {
				t.Errorf("Not nil vocab expected")
			}
			if c.expectedVocab != nil && vocab != nil && !reflect.DeepEqual(c.expectedVocab, vocab) {
				t.Errorf("Expected vocab:%+v;Actual:%+v", c.expectedVocab, vocab)
			}
			mockedRepo.Reset()
		})
	}

}

func TestVocabWithLocalRepo_GetVocabEntryByText(t *testing.T) {
	testCases := []struct {
		text                      string
		expectedEntry             *domain.VocabEntry
		expectLocalGetByTextInv   bool
		expectServiceGetByTextInv bool
		expectLocalAddInv         bool
		expectErr                 bool
	}{
		{
			text: "Positive: found in local repo",
			expectedEntry: &domain.VocabEntry{
				Text: "Positive: found in local repo",
			},
			expectLocalGetByTextInv: true,
		},
		{
			text:                      "Positive: not found in the entry service",
			expectLocalGetByTextInv:   true,
			expectServiceGetByTextInv: true,
		},
		{
			text: "Positive: found in the entry service",
			expectedEntry: &domain.VocabEntry{
				Text: "Positive: found in the entry service",
			},
			expectLocalGetByTextInv:   true,
			expectServiceGetByTextInv: true,
			expectLocalAddInv:         true,
		},
		{
			text:                    "Local GetVocabEntryByText returns error",
			expectLocalGetByTextInv: true,
			expectErr:               true,
		},
		{
			text:                      "Entry service GetVocabEntryByText returns error",
			expectLocalGetByTextInv:   true,
			expectServiceGetByTextInv: true,
			expectErr:                 true,
		},
		{
			text:                      "Local AddVocabEntry returns error",
			expectLocalGetByTextInv:   true,
			expectServiceGetByTextInv: true,
			expectLocalAddInv:         true,
			expectErr:                 true,
		},
	}

	mockedRepo := &mock.VocabRepo{
		GetVocabEntryByTextFn: func(text string) (*domain.VocabEntry, error) {
			switch text {
			case "Positive: found in local repo":
				return &domain.VocabEntry{Text: text}, nil
			case "Local GetVocabEntryByText returns error":
				return nil, fmt.Errorf("error")
			default:
				return nil, nil
			}
		},
		AddVocabEntryFn: func(entry *domain.VocabEntry) (*domain.VocabEntry, error) {
			switch entry.Text {
			case "Positive: found in the entry service":
				return &domain.VocabEntry{Text: entry.Text}, nil
			case "Local AddVocabEntry returns error":
				return nil, fmt.Errorf("error")
			default:
				return nil, nil
			}
		},
	}
	mockedEntryService := &mock.VocabEntryService{
		GetVocabEntryByTextFn: func(text string) (*domain.VocabEntry, error) {
			switch text {
			case "Positive: found in the entry service", "Local AddVocabEntry returns error":
				return &domain.VocabEntry{Text: text}, nil
			case "Entry service GetVocabEntryByText returns error":
				return nil, fmt.Errorf("error")
			default:
				return nil, nil
			}
		},
	}

	vocabService := NewVocabWithLocalRepo(mock.Logger{}, mockedRepo, mockedEntryService)
	for _, c := range testCases {
		t.Run(c.text, func(t *testing.T) {
			entry, err := vocabService.GetVocabEntryByText(c.text)
			if c.expectErr == false && err != nil {
				t.Errorf("Expected no error, but got %s", err)
			}
			if c.expectErr && err == nil {
				t.Errorf("Expected error, but got nothing")
			}
			if c.expectLocalGetByTextInv != mockedRepo.GetVocabEntryByTextInvoked {
				t.Errorf(
					"Actual invocation of local repo GetVocabEntryByText(%v) doesn't match expectations",
					mockedRepo.GetVocabEntryByTextInvoked,
				)
			}
			if c.expectServiceGetByTextInv != mockedEntryService.GetVocabEntryByTextInvoked {
				t.Errorf(
					"Actual invocation of entry service GetVocabEntryByText(%v) doesn't match expectations",
					mockedEntryService.GetVocabEntryByTextInvoked,
				)
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
			mockedRepo.Reset()
			mockedEntryService.Reset()
		})
	}
}
