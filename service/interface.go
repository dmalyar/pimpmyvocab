package service

import "github.com/dmalyar/pimpmyvocab/domain"

// Vocab provides use cases for vocabs.
type Vocab interface {
	CreateVocab(userID int) (*domain.Vocab, error)
	ClearUserVocab(userID int) error

	AddEntryToUserVocab(entryID, userID int) error
	CheckEntryInUserVocab(entryID, userID int) (bool, error)
	RemoveEntryFromUserVocab(entryID, userID int) error

	GetEntriesFromUserVocab(userID int) ([]*domain.VocabEntry, error)

	VocabEntry
}

type VocabEntry interface {
	GetVocabEntryByText(text string) (*domain.VocabEntry, error)
	GetVocabEntryByID(id int) (*domain.VocabEntry, error)
}
