package service

import "github.com/dmalyar/pimpmyvocab/domain"

// Vocab provides use cases for vocabs.
type Vocab interface {
	CreateVocab(userID int) (*domain.Vocab, error)
	ClearUserVocab(userID int) error

	AddEntryToUserVocab(entryID, userID int) error
	CheckEntryInUserVocab(entryID, userID int) (bool, error)
	GetRandomEntryFromUserVocab(userID, previousEntryID int) (*domain.VocabEntry, error)
	GetEntriesFromUserVocab(userID int) ([]*domain.VocabEntry, error)
	RemoveEntryFromUserVocab(entryID, userID int) error

	VocabEntry
}

type VocabEntry interface {
	GetVocabEntryByText(text string) (*domain.VocabEntry, error)
	GetVocabEntryByID(id int) (*domain.VocabEntry, error)
}
