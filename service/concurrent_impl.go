package service

import (
	"github.com/dmalyar/pimpmyvocab/domain"
	"sync"
)

// ConcurrentVocab wraps another implementation of service.Vocab and adds logic for correct concurrency work.
// Implements service.Vocab itself.
type ConcurrentVocab struct {
	wrappedService Vocab
	vocabSync      *userIDSync
	vocabEntrySync *textSync
}

// Returns ready to use ConcurrentVocab.
func NewConcurrentVocab(wrappedService Vocab) *ConcurrentVocab {
	return &ConcurrentVocab{
		wrappedService: wrappedService,
		vocabSync: &userIDSync{
			mu:     sync.Mutex{},
			inWork: make(map[int]chan struct{}),
		},
		vocabEntrySync: &textSync{
			mu:     sync.Mutex{},
			inWork: make(map[string]chan struct{}),
		},
	}
}

// CreateVocab calls CreateVocab of wrapped vocabService with concurrent safe logic.
// Makes one call of the wrapped method at a time per user.
func (v *ConcurrentVocab) CreateVocab(userID int) (*domain.Vocab, error) {
	v.vocabSync.startWork(userID)
	defer v.vocabSync.endWork(userID)
	return v.wrappedService.CreateVocab(userID)
}

func (v *ConcurrentVocab) GetVocabEntryByText(text string) (*domain.VocabEntry, error) {
	v.vocabEntrySync.startWork(text)
	defer v.vocabEntrySync.endWork(text)
	return v.wrappedService.GetVocabEntryByText(text)
}

func (v *ConcurrentVocab) GetVocabEntryByID(id int) (*domain.VocabEntry, error) {
	return v.wrappedService.GetVocabEntryByID(id)
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

type textSync struct {
	mu     sync.Mutex
	inWork map[string]chan struct{}
}

func (s *textSync) startWork(text string) {
	s.mu.Lock()
	ch, ok := s.inWork[text]
	if !ok {
		s.inWork[text] = make(chan struct{})
		s.mu.Unlock()
		return
	}
	s.mu.Unlock()
	<-ch
	s.startWork(text)
}

func (s *textSync) endWork(text string) {
	s.mu.Lock()
	ch := s.inWork[text]
	delete(s.inWork, text)
	close(ch)
	s.mu.Unlock()
}
