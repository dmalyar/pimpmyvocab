package mock

import (
	"github.com/dmalyar/pimpmyvocab/domain"
	"github.com/dmalyar/pimpmyvocab/log"
	"sync"
)

// VocabRepo is a mock struct implementing repo.Vocab interface.
type VocabRepo struct {
	AddVocabFn      func(vocab *domain.Vocab) (*domain.Vocab, error)
	AddVocabInvoked bool

	GetVocabByUserIDFn      func(userID int) (*domain.Vocab, error)
	GetVocabByUserIDInvoked bool
}

// AddVocab registers invocation of AddVocab func and calls it.
func (r *VocabRepo) AddVocab(vocab *domain.Vocab) (*domain.Vocab, error) {
	r.AddVocabInvoked = true
	return r.AddVocabFn(vocab)
}

// GetVocabByUserID registers invocation of GetVocabByUserID func and calls it.
func (r *VocabRepo) GetVocabByUserID(userID int) (*domain.Vocab, error) {
	r.GetVocabByUserIDInvoked = true
	return r.GetVocabByUserIDFn(userID)
}

// Reset resets functions invocation.
func (r *VocabRepo) Reset() {
	r.AddVocabInvoked = false
	r.GetVocabByUserIDInvoked = false
}

// vocabService is a mock struct implementing service.Vocab interface.
type VocabService struct {
	CreateVocabFn      func(userID int) (*domain.Vocab, error)
	CreateVocabInvoked bool
}

// CreateVocab registers invocation of CreateVocab func and calls it.
func (s *VocabService) CreateVocab(userID int) (*domain.Vocab, error) {
	s.CreateVocabInvoked = true
	return s.CreateVocabFn(userID)
}

// Reset resets functions invocation.
func (s *VocabService) Reset() {
	s.CreateVocabInvoked = false
}

// VocabServiceConcurrency is a mock struct implementing service.Vocab interface.
// You can test concurrent execution of methods with this struct.
type VocabServiceConcurrency struct {
	CreateVocabFn func(userID int) (*domain.Vocab, error)

	mu                             sync.Mutex
	CreateVocabInvoked             bool
	CreateVocabConcurrentlyInvoked bool
	createVocabConcurrencyCheck    map[int]struct{}
}

// NewVocabServiceConcurrency returns ready to use VocabServiceConcurrency.
func NewVocabServiceConcurrency(fn func(userID int) (*domain.Vocab, error)) *VocabServiceConcurrency {
	return &VocabServiceConcurrency{
		CreateVocabFn:               fn,
		createVocabConcurrencyCheck: make(map[int]struct{}),
	}
}

// CreateVocab registers invocation of CreateVocab func and calls it.
// Also registers if it was called concurrently by the same user.
func (s *VocabServiceConcurrency) CreateVocab(userID int) (*domain.Vocab, error) {
	s.mu.Lock()
	s.CreateVocabInvoked = true
	_, ok := s.createVocabConcurrencyCheck[userID]
	if ok {
		s.CreateVocabConcurrentlyInvoked = true
	} else {
		s.createVocabConcurrencyCheck[userID] = struct{}{}
	}
	s.mu.Unlock()
	defer func() {
		s.mu.Lock()
		delete(s.createVocabConcurrencyCheck, userID)
		s.mu.Unlock()
	}()
	return s.CreateVocabFn(userID)
}

// Reset resets functions invocation and concurrent functions invocation.
func (s *VocabServiceConcurrency) Reset() {
	s.CreateVocabInvoked = false
	s.CreateVocabConcurrentlyInvoked = false
}

type Logger struct {
}

func (_ Logger) Debugf(_ string, _ ...interface{}) {
}

func (_ Logger) Infof(_ string, _ ...interface{}) {
}

func (_ Logger) Warnf(_ string, _ ...interface{}) {
}

func (_ Logger) Errorf(_ string, _ ...interface{}) {
}

func (_ Logger) Panicf(_ string, _ ...interface{}) {
}

func (_ Logger) Debug(_ ...interface{}) {
}

func (_ Logger) Info(_ ...interface{}) {
}

func (_ Logger) Warn(_ ...interface{}) {
}

func (_ Logger) Error(_ ...interface{}) {
}

func (_ Logger) Panic(_ ...interface{}) {
}

func (_ Logger) Println(_ ...interface{}) {
}

func (_ Logger) Printf(_ string, _ ...interface{}) {
}

func (_ Logger) WithFields(_ map[string]interface{}) log.Logger {
	return Logger{}
}
