package service

import (
	"errors"
	"github.com/dmalyar/pimpmyvocab/domain"
	"github.com/dmalyar/pimpmyvocab/mock"
	"testing"
)

type createVocabCase struct {
	userID                    int
	expectedVocab             *domain.Vocab
	expectErr                 bool
	expectGetVocabByUserIDInv bool
	expectAddVocabInv         bool
}

func TestVocabDB_CreateVocab(t *testing.T) {
	cases := []createVocabCase{
		{ // All good
			userID: 1,
			expectedVocab: &domain.Vocab{
				ID:     1,
				UserID: 1,
			},
			expectGetVocabByUserIDInv: true,
			expectAddVocabInv:         true,
		},
		{ // GetVocabByUserID returns error
			userID:                    2,
			expectErr:                 true,
			expectGetVocabByUserIDInv: true,
		},
		{ // GetVocabByUserID returns vocab
			userID:                    3,
			expectGetVocabByUserIDInv: true,
		},
		{ // AddVocab returns error
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
	vocabService := NewVocabWithLocalRepo(mock.Logger{}, mockedRepo)
	for _, c := range cases {
		vocab, err := vocabService.CreateVocab(c.userID)
		if c.expectErr == false && err != nil {
			t.Errorf("Expected no error, but got %s\n", err)
		}
		if c.expectErr && err == nil {
			t.Errorf("Expected error, but got nothing\n")
		}
		if c.expectAddVocabInv != mockedRepo.AddVocabInvoked {
			t.Errorf("Actual invocation of AddVocab(%v) doesn't match expectations\n", mockedRepo.AddVocabInvoked)
		}
		if c.expectGetVocabByUserIDInv != mockedRepo.GetVocabByUserIDInvoked {
			t.Errorf("Actual invocation of GetVocabByUserID(%v) doesn't match expectations\n", mockedRepo.GetVocabByUserIDInvoked)
		}
		if c.expectedVocab == nil && vocab != nil {
			t.Errorf("Nil vocab expected\n")
		}
		if c.expectedVocab != nil && vocab == nil {
			t.Errorf("Not nil vocab expected\n")
		}
		if c.expectedVocab != nil && vocab != nil && *c.expectedVocab != *vocab {
			t.Errorf("Expected vocab:%+v\nActual:%+v", c.expectedVocab, vocab)
		}
		mockedRepo.Reset()
	}

}
