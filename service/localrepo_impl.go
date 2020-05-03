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
func (v *VocabWithLocalRepo) CreateVocab(userID int) (*domain.Vocab, error) {
	logger := v.logger.WithField("userID", userID)
	logger.Debug("Checking if user already has a vocab")
	vocab, err := v.localRepo.GetVocabByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("getting vocab by user ID: %s", err)
	}
	if vocab != nil {
		logger.Info("User already has vocab, no need to create a new one")
		return nil, nil
	}
	logger.Debug("Creating vocab")
	vocab, err = v.localRepo.AddVocab(&domain.Vocab{UserID: userID})
	if err != nil {
		return nil, fmt.Errorf("creating vocab: %s", err)
	}
	logger.Infof("Vocab created: %s", vocab)
	return vocab, nil
}

// ClearUserVocab clears the user's vocab by removing all entries from it.
func (v *VocabWithLocalRepo) ClearUserVocab(userID int) error {
	logger := v.logger.WithField("userID", userID)
	logger.Debugf("Clearing the user's vocab")
	err := v.localRepo.ClearVocabByUserID(userID)
	if err != nil {
		return fmt.Errorf("removing all entries from the user's vocab: %s", err)
	}
	return nil
}

// AddEntryToUserVocab adds the vocab entry to the user's vocab.
// If user already has this entry added to the vocab then do nothing.
func (v *VocabWithLocalRepo) AddEntryToUserVocab(entryID, userID int) error {
	logger := v.logger.WithFields(map[string]interface{}{
		"entryID": entryID,
		"userID":  userID,
	})
	inVocab, err := v.CheckEntryInUserVocab(entryID, userID)
	if err != nil {
		return err
	}
	if inVocab {
		return nil
	}
	logger.Debug("Adding the entry to the user's vocab")
	err = v.localRepo.AddEntryToUserVocab(entryID, userID)
	if err != nil {
		return fmt.Errorf("adding entry to user's vocab: %s", err)
	}
	logger.Info("Vocab entry added to the user's vocab")
	return nil
}

// CheckEntryInUserVocab checks if user already has the entry added to the vocab.
func (v *VocabWithLocalRepo) CheckEntryInUserVocab(entryID, userID int) (bool, error) {
	logger := v.logger.WithFields(map[string]interface{}{
		"entryID": entryID,
		"userID":  userID,
	})
	logger.Debug("Checking if the entry is added to the user's vocab")
	inVocab, err := v.localRepo.CheckEntryInUserVocab(entryID, userID)
	if err != nil {
		return false, fmt.Errorf("checking if entry is added to vocab: %s", err)
	}
	if inVocab {
		logger.Info("Vocab entry is in the user's vocab")
	} else {
		logger.Info("Vocab entry is not in the user's vocab")
	}
	return inVocab, nil
}

// RemoveEntryFromUserVocab removes the vocab entry from the user's vocab.
// If the entry is not in the user's vocab then do nothing.
func (v *VocabWithLocalRepo) RemoveEntryFromUserVocab(entryID, userID int) error {
	logger := v.logger.WithFields(map[string]interface{}{
		"entryID": entryID,
		"userID":  userID,
	})
	inVocab, err := v.CheckEntryInUserVocab(entryID, userID)
	if err != nil {
		return err
	}
	if !inVocab {
		return nil
	}
	logger.Debug("Removing the entry from the user's vocab")
	err = v.localRepo.RemoveEntryFromUserVocab(entryID, userID)
	if err != nil {
		return fmt.Errorf("removing entry from user's vocab: %s", err)
	}
	logger.Info("Vocab entry removed from the user's vocab")
	return nil
}

// GetEntriesByUserID returns all entries linked to the user's vocab.
func (v *VocabWithLocalRepo) GetEntriesFromUserVocab(userID int) ([]*domain.VocabEntry, error) {
	logger := v.logger.WithField("userID", userID)
	logger.Debugf("Getting vocab entries")
	entries, err := v.localRepo.GetEntriesByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("getting vocab entries by user ID: %s", err)
	}
	logger.Infof("Found %v entry(-ies)", len(entries))
	return entries, nil
}

// GetVocabEntryByText looks for vocab entry in the local repo by the given text.
// If it's found then returns it. If not then calls entry service method. If entry is found there then adds it
// to the local repo.
// If it's not found there then returns nil.
func (v *VocabWithLocalRepo) GetVocabEntryByText(text string) (*domain.VocabEntry, error) {
	logger := v.logger.WithField("text", text)
	logger.Info("Getting vocab entry")
	entry, err := v.localRepo.GetVocabEntryByText(text)
	if err != nil {
		return nil, fmt.Errorf("getting vocab entry by text in the local repo: %s", err)
	}
	if entry != nil {
		logger.Info("Found vocab entry in the local repo: %s", entry)
		return entry, nil
	}
	logger.Info("Vocab entry not found in the local repo")
	entry, err = v.entryService.GetVocabEntryByText(text)
	if err != nil {
		return nil, fmt.Errorf("getting vocab entry from the vocab entry service: %s", err)
	}
	if entry == nil {
		logger.Info("Vocab entry not found in the vocab entry service")
		return nil, nil
	}
	logger.Info("Vocab entry found in the vocab entry service")
	entry, err = v.localRepo.AddVocabEntry(entry)
	if err != nil {
		return nil, fmt.Errorf("adding vocab entry to the local repo: %s", err)
	}
	logger.Info("Vocab entry added to the local repo: %s", entry)
	return entry, nil
}

// GetVocabEntryByID returns vocab entry found in the local repo by ID.
// Returns nil if entry is not found.
func (v *VocabWithLocalRepo) GetVocabEntryByID(id int) (*domain.VocabEntry, error) {
	logger := v.logger.WithField("id", id)
	logger.Info("Getting vocab entry")
	ve, err := v.localRepo.GetVocabEntryByID(id)
	if err != nil {
		return nil, fmt.Errorf("error getting vocab entry: %s", err)
	}
	return ve, nil
}
