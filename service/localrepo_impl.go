package service

import (
	"fmt"
	"github.com/dmalyar/pimpmyvocab/domain"
	"github.com/dmalyar/pimpmyvocab/log"
	"github.com/dmalyar/pimpmyvocab/repo"
)

// VocabWithLocalRepo implements service.Vocab interface for working with local repository.
type VocabWithLocalRepo struct {
	logger       log.Logger
	localRepo    repo.Vocab
	entryService VocabEntry
}

func NewVocabWithLocalRepo(logger log.Logger, localRepo repo.Vocab, entryService VocabEntry) *VocabWithLocalRepo {
	return &VocabWithLocalRepo{
		logger:       logger,
		localRepo:    localRepo,
		entryService: entryService,
	}
}

// CreateVocab creates a vocab in local localRepo for user.
// Returns the created vocab entity.
// Returns nil and skips vocab creation if user already has a vocab.
// Return error if it occurs.
func (v *VocabWithLocalRepo) CreateVocab(userID int) (*domain.Vocab, error) {
	contextLog := v.logger.WithField("userID", userID)
	contextLog.Debug("Checking if user already has a vocab")
	vocab, err := v.localRepo.GetVocabByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("getting vocab by user ID: %s", err)
	}
	if vocab != nil {
		contextLog.Info("User already has vocab, no need to create a new one")
		return nil, nil
	}
	contextLog.Debug("Creating vocab")
	vocab, err = v.localRepo.AddVocab(&domain.Vocab{UserID: userID})
	if err != nil {
		return nil, fmt.Errorf("creating vocab: %s", err)
	}
	contextLog.Infof("Vocab created: %s", vocab)
	return vocab, nil
}

// GetVocabEntryByText looks for vocab entry in the local repo by the given text.
// If it's found then returns it. If not then calls entry service method. If entry is found there then adds it
// to the local repo.
// If it's not found there then returns nil.
// Returns error if it occurs.
func (v *VocabWithLocalRepo) GetVocabEntryByText(text string) (*domain.VocabEntry, error) {
	contextLog := v.logger.WithField("text", text)
	contextLog.Info("Getting vocab entry")
	entry, err := v.localRepo.GetVocabEntryByText(text)
	if err != nil {
		return nil, fmt.Errorf("getting vocab entry by text in the local repo: %s", err)
	}
	if entry != nil {
		contextLog.Info("Found vocab entry in the local repo: %s", entry)
		return entry, nil
	}
	contextLog.Info("Vocab entry not found in the local repo")
	entry, err = v.entryService.GetVocabEntryByText(text)
	if err != nil {
		return nil, fmt.Errorf("getting vocab entry from the vocab entry service: %s", err)
	}
	if entry == nil {
		contextLog.Info("Vocab entry not found in the vocab entry service")
		return nil, nil
	}
	contextLog.Info("Vocab entry found in the vocab entry service")
	entry, err = v.localRepo.AddVocabEntry(entry)
	if err != nil {
		return nil, fmt.Errorf("adding vocab entry to the local repo: %s", err)
	}
	contextLog.Info("Vocab entry added to the local repo: %s", entry)
	return entry, nil
}

// GetVocabEntryByID returns vocab entry found in the local repo by ID.
// Returns nil if entry is not found.
// Returns error if it occurs.
func (v *VocabWithLocalRepo) GetVocabEntryByID(id int) (*domain.VocabEntry, error) {
	contextLog := v.logger.WithField("ID", id)
	contextLog.Info("Getting vocab entry")
	ve, err := v.localRepo.GetVocabEntryByID(id)
	if err != nil {
		return nil, fmt.Errorf("error getting vocab entry: %s", err)
	}
	return ve, nil
}
