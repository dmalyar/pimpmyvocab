package service

import "github.com/dmalyar/pimpmyvocab/domain"

// Vocab provides use cases for vocabs.
type Vocab interface {
	CreateVocab(userID int) (*domain.Vocab, error)
}
