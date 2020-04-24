package service

import (
	"fmt"
	"github.com/dmalyar/pimpmyvocab/domain"
	"github.com/dmalyar/pimpmyvocab/log"
	"github.com/dmalyar/pimpmyvocab/repo"
)

// VocabWithLocalRepo implements service.Vocab interface for working with local repository.
type VocabWithLocalRepo struct {
	logger    log.Logger
	localRepo repo.Vocab
}

func NewVocabWithLocalRepo(logger log.Logger, localRepo repo.Vocab) *VocabWithLocalRepo {
	return &VocabWithLocalRepo{
		logger:    logger,
		localRepo: localRepo,
	}
}

// CreateVocab creates a vocab in local localRepo for user.
// Returns the created vocab entity.
// Returns nil and skips vocab creation if user already has a vocab.
func (s *VocabWithLocalRepo) CreateVocab(userID int) (*domain.Vocab, error) {
	contextLog := s.logger.WithFields(map[string]interface{}{
		"userID": userID,
	})
	contextLog.Info("Creating vocab")
	vocab, err := s.localRepo.GetVocabByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("error getting vocab: %s", err)
	}
	if vocab != nil {
		contextLog.Info("User already has vocab, no need to create a new one")
		return nil, nil
	}
	vocab, err = s.localRepo.AddVocab(&domain.Vocab{UserID: userID})
	if err != nil {
		return nil, fmt.Errorf("error creating vocab: %s", err)
	}
	contextLog.Info("Vocab created")
	return vocab, nil
}
