package service

import (
	"errors"
	"fmt"
	"github.com/dmalyar/pimpmyvocab/domain"
	"github.com/dmalyar/pimpmyvocab/mock"
	"reflect"
	"testing"
)

func TestVocabWithLocalRepo_CreateVocab(t *testing.T) {
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

func TestVocabWithLocalRepo_CheckEntryInUserVocab(t *testing.T) {
	testCases := []struct {
		name            string
		entryID, userID int
		expectedRes     bool
		expectErr       bool
	}{
		{
			name:        "Positive",
			entryID:     1,
			userID:      1,
			expectedRes: true,
		},
		{
			name:      "Repo returns error",
			entryID:   2,
			userID:    2,
			expectErr: true,
		},
	}

	mockedRepo := &mock.VocabRepo{
		CheckEntryInUserVocabFn: func(entryID, userID int) (bool, error) {
			switch {
			case entryID == 1 && userID == 1:
				return true, nil
			case entryID == 2 && userID == 2:
				return false, fmt.Errorf("error")
			default:
				return false, nil
			}
		},
	}

	vocabService := NewVocabWithLocalRepo(mock.Logger{}, mockedRepo, &mock.VocabEntryService{})
	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			res, err := vocabService.CheckEntryInUserVocab(c.entryID, c.userID)
			if c.expectErr == false && err != nil {
				t.Errorf("Expected no error, but got %s", err)
			}
			if c.expectErr && err == nil {
				t.Errorf("Expected error, but got nothing")
			}
			if c.expectedRes != res {
				t.Errorf("Expected res:%+v;Actual:%+v", c.expectedRes, res)
			}
			if !mockedRepo.CheckEntryInUserVocabInvoked {
				t.Errorf("CheckEntryInUserVocab was not invoked")
			}
			mockedRepo.Reset()
		})
	}
}

func TestVocabWithLocalRepo_AddEntryToUserVocab(t *testing.T) {
	testCases := []struct {
		name                string
		entryID, userID     int
		expectErr           bool
		expectCheckEntryInv bool
		expectAddEntryInv   bool
	}{
		{
			name:                "Positive added to vocab",
			entryID:             1,
			userID:              1,
			expectCheckEntryInv: true,
			expectAddEntryInv:   true,
		},
		{
			name:                "Positive was in vocab",
			entryID:             2,
			userID:              2,
			expectCheckEntryInv: true,
		},
		{
			name:                "Check returns error",
			entryID:             3,
			userID:              3,
			expectCheckEntryInv: true,
			expectErr:           true,
		},
		{
			name:                "Add returns error",
			entryID:             4,
			userID:              4,
			expectCheckEntryInv: true,
			expectAddEntryInv:   true,
			expectErr:           true,
		},
	}

	mockedRepo := &mock.VocabRepo{
		CheckEntryInUserVocabFn: func(entryID, userID int) (bool, error) {
			switch {
			case entryID == 2 && userID == 2:
				return true, nil
			case entryID == 3 && userID == 3:
				return false, fmt.Errorf("error")
			default:
				return false, nil
			}
		},
		AddEntryToUserVocabFn: func(entryID, userID int) error {
			switch {
			case entryID == 4 && userID == 4:
				return fmt.Errorf("error")
			default:
				return nil
			}
		},
	}

	vocabService := NewVocabWithLocalRepo(mock.Logger{}, mockedRepo, &mock.VocabEntryService{})
	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			err := vocabService.AddEntryToUserVocab(c.entryID, c.userID)
			if c.expectErr == false && err != nil {
				t.Errorf("Expected no error, but got %s", err)
			}
			if c.expectCheckEntryInv != mockedRepo.CheckEntryInUserVocabInvoked {
				t.Errorf("Actual invocation of CheckEntryInUserVocab(%v) doesn't match expectations", mockedRepo.CheckEntryInUserVocabInvoked)
			}
			if c.expectAddEntryInv != mockedRepo.AddEntryToUserVocabInvoked {
				t.Errorf("Actual invocation of AddEntryToUserVocab(%v) doesn't match expectations", mockedRepo.AddEntryToUserVocabInvoked)
			}
			if c.expectErr && err == nil {
				t.Errorf("Expected error, but got nothing")
			}
			mockedRepo.Reset()
		})
	}
}

func TestVocabWithLocalRepo_RemoveEntryFromUserVocab(t *testing.T) {
	testCases := []struct {
		name                 string
		entryID, userID      int
		expectErr            bool
		expectCheckEntryInv  bool
		expectRemoveEntryInv bool
	}{
		{
			name:                 "Positive removed to vocab",
			entryID:              1,
			userID:               1,
			expectCheckEntryInv:  true,
			expectRemoveEntryInv: true,
		},
		{
			name:                "Positive was not in vocab",
			entryID:             2,
			userID:              2,
			expectCheckEntryInv: true,
		},
		{
			name:                "Check returns error",
			entryID:             3,
			userID:              3,
			expectCheckEntryInv: true,
			expectErr:           true,
		},
		{
			name:                 "Remove returns error",
			entryID:              4,
			userID:               4,
			expectCheckEntryInv:  true,
			expectRemoveEntryInv: true,
			expectErr:            true,
		},
	}

	mockedRepo := &mock.VocabRepo{
		CheckEntryInUserVocabFn: func(entryID, userID int) (bool, error) {
			switch {
			case entryID == 2 && userID == 2:
				return false, nil
			case entryID == 3 && userID == 3:
				return false, fmt.Errorf("error")
			default:
				return true, nil
			}
		},
		RemoveEntryFromUserVocabFn: func(entryID, userID int) error {
			switch {
			case entryID == 4 && userID == 4:
				return fmt.Errorf("error")
			default:
				return nil
			}
		},
	}

	vocabService := NewVocabWithLocalRepo(mock.Logger{}, mockedRepo, &mock.VocabEntryService{})
	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			err := vocabService.RemoveEntryFromUserVocab(c.entryID, c.userID)
			if c.expectErr == false && err != nil {
				t.Errorf("Expected no error, but got %s", err)
			}
			if c.expectCheckEntryInv != mockedRepo.CheckEntryInUserVocabInvoked {
				t.Errorf("Actual invocation of CheckEntryInUserVocab(%v) doesn't match expectations", mockedRepo.CheckEntryInUserVocabInvoked)
			}
			if c.expectRemoveEntryInv != mockedRepo.RemoveEntryFromUserVocabInvoked {
				t.Errorf("Actual invocation of RemoveEntryFromUserVocab(%v) doesn't match expectations", mockedRepo.RemoveEntryFromUserVocabInvoked)
			}
			if c.expectErr && err == nil {
				t.Errorf("Expected error, but got nothing")
			}
			mockedRepo.Reset()
		})
	}
}

func TestVocabWithLocalRepo_GetVocabEntriesByUserID(t *testing.T) {
	testCases := []struct {
		name            string
		userID          int
		expectedEntries []*domain.VocabEntry
		expectErr       bool
	}{
		{
			name:   "Positive",
			userID: 1,
			expectedEntries: []*domain.VocabEntry{
				{ID: 1, Text: "One"},
				{ID: 2, Text: "Two"},
			},
		},
		{
			name:   "Positive no entries found",
			userID: 2,
		},
		{
			name:      "Get vocab entries returns err",
			userID:    3,
			expectErr: true,
		},
	}

	mockedRepo := &mock.VocabRepo{
		GetVocabEntriesByUserIDFn: func(userID int) ([]*domain.VocabEntry, error) {
			switch userID {
			case 1:
				return []*domain.VocabEntry{
					{ID: 1, Text: "One"},
					{ID: 2, Text: "Two"},
				}, nil
			case 3:
				return nil, fmt.Errorf("error")
			default:
				return nil, nil
			}
		},
	}

	vocabService := NewVocabWithLocalRepo(mock.Logger{}, mockedRepo, &mock.VocabEntryService{})
	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			entries, err := vocabService.GetVocabEntriesByUserID(c.userID)
			if c.expectErr && err == nil {
				t.Errorf("Expected error, but got nothing")
			}
			if !reflect.DeepEqual(c.expectedEntries, entries) {
				t.Errorf("Expected res:%+v;Actual:%+v", c.expectedEntries, entries)
			}
			if !mockedRepo.GetVocabEntriesByUserIDInvoked {
				t.Errorf("GetVocabEntriesByUserIDInvoked was not invoked")
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
