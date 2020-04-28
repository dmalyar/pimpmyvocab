package service

import "github.com/dmalyar/pimpmyvocab/domain"

// Vocab provides use cases for vocabs.
type Vocab interface {
	VocabEntry
	CreateVocab(userID int) (*domain.Vocab, error)
}

type VocabEntry interface {
	GetVocabEntryByText(text string) (*domain.VocabEntry, error)
	GetVocabEntryByID(id int) (*domain.VocabEntry, error)
}
