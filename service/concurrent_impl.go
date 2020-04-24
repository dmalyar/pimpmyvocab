package service

import (
	"github.com/dmalyar/pimpmyvocab/domain"
	"sync"
)

// ConcurrentVocab wraps another implementation of service.Vocab and adds logic for correct concurrency work.
// Implements service.Vocab itself.
type ConcurrentVocab struct {
	wrappedService  Vocab
	vocabCreateSync *userIDSync
}

// Returns ready to use ConcurrentVocab.
func NewConcurrentVocab(wrappedService Vocab) *ConcurrentVocab {
	return &ConcurrentVocab{
		wrappedService: wrappedService,
		vocabCreateSync: &userIDSync{
			mu:     sync.Mutex{},
			inWork: make(map[int]chan struct{}),
		},
	}
}

// CreateVocab calls CreateVocab of wrapped vocabService with concurrent safe logic.
// Makes one call of the wrapped method at a time per user.
func (v *ConcurrentVocab) CreateVocab(userID int) (*domain.Vocab, error) {
	v.vocabCreateSync.startWork(userID)
	defer v.vocabCreateSync.endWork(userID)
	return v.wrappedService.CreateVocab(userID)
}

type userIDSync struct {
	mu     sync.Mutex
	inWork map[int]chan struct{}
}

func (s *userIDSync) startWork(userID int) {
	s.mu.Lock()
	ch, ok := s.inWork[userID]
	if !ok {
		s.inWork[userID] = make(chan struct{})
		s.mu.Unlock()
		return
	}
	s.mu.Unlock()
	<-ch
	s.startWork(userID)
}

func (s *userIDSync) endWork(userID int) {
	s.mu.Lock()
	ch := s.inWork[userID]
	delete(s.inWork, userID)
	close(ch)
	s.mu.Unlock()
}
