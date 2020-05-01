package repo

import "github.com/dmalyar/pimpmyvocab/domain"

// Vocab provides methods for interacting with vocabs on repository level.
type Vocab interface {
	AddVocab(vocab *domain.Vocab) (*domain.Vocab, error)
	GetVocabByUserID(userID int) (*domain.Vocab, error)

	AddVocabEntry(entry *domain.VocabEntry) (*domain.VocabEntry, error)
	GetVocabEntryByText(text string) (*domain.VocabEntry, error)
	GetVocabEntryByID(id int) (*domain.VocabEntry, error)
	CheckEntryInUserVocab(entryID, userID int) (bool, error)
	AddEntryToUserVocab(entryID, userID int) error
	RemoveEntryFromUserVocab(entryID, userID int) error
	GetVocabEntriesByUserID(userID int) ([]*domain.VocabEntry, error)
}
