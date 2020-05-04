package repo

import "github.com/dmalyar/pimpmyvocab/domain"

// Vocab provides methods for interacting with vocabs on repository level.
type Vocab interface {
	AddVocab(vocab *domain.Vocab) (*domain.Vocab, error)
	GetVocabByUserID(userID int) (*domain.Vocab, error)
	ClearVocabByUserID(userID int) error

	AddVocabEntry(entry *domain.VocabEntry) (*domain.VocabEntry, error)
	GetVocabEntryByText(text string) (*domain.VocabEntry, error)
	GetVocabEntryByID(id int) (*domain.VocabEntry, error)

	AddEntryToUserVocab(entryID, userID int) error
	CheckEntryInUserVocab(entryID, userID int) (bool, error)
	GetEntryIDsByUserID(userID int) ([]int, error)
	GetEntriesByUserID(userID int) ([]*domain.VocabEntry, error)
	RemoveEntryFromUserVocab(entryID, userID int) error
}
