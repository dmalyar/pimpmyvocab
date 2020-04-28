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

	AddVocabEntryFn      func(vocab *domain.VocabEntry) (*domain.VocabEntry, error)
	AddVocabEntryInvoked bool

	GetVocabEntryByTextFn      func(text string) (*domain.VocabEntry, error)
	GetVocabEntryByTextInvoked bool

	GetVocabEntryByIDFn      func(ID int) (*domain.VocabEntry, error)
	GetVocabEntryByIDInvoked bool
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

// AddVocabEntry registers invocation of AddVocabEntry func and calls it.
func (r *VocabRepo) AddVocabEntry(entry *domain.VocabEntry) (*domain.VocabEntry, error) {
	r.AddVocabEntryInvoked = true
	return r.AddVocabEntryFn(entry)
}

// GetVocabEntryByText registers invocation of GetVocabEntryByText func and calls it.
func (r *VocabRepo) GetVocabEntryByText(text string) (*domain.VocabEntry, error) {
	r.GetVocabEntryByTextInvoked = true
	return r.GetVocabEntryByTextFn(text)
}

// GetVocabEntryByID registers invocation of GetVocabEntryByID func and calls it.
func (r *VocabRepo) GetVocabEntryByID(ID int) (*domain.VocabEntry, error) {
	r.GetVocabEntryByIDInvoked = true
	return r.GetVocabEntryByIDFn(ID)
}

// Reset resets functions invocation.
func (r *VocabRepo) Reset() {
	r.AddVocabInvoked = false
	r.GetVocabByUserIDInvoked = false
	r.AddVocabEntryInvoked = false
	r.GetVocabEntryByTextInvoked = false
	r.GetVocabEntryByIDInvoked = false
}

// vocabService is a mock struct implementing service.Vocab interface.
type VocabService struct {
	CreateVocabFn      func(userID int) (*domain.Vocab, error)
	CreateVocabInvoked bool

	VocabEntryService
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

type VocabEntryService struct {
	GetVocabEntryByTextFn      func(text string) (*domain.VocabEntry, error)
	GetVocabEntryByTextInvoked bool

	GetVocabEntryByIDFn      func(id int) (*domain.VocabEntry, error)
	GetVocabEntryByIDInvoked bool
}

func (v *VocabEntryService) GetVocabEntryByText(text string) (*domain.VocabEntry, error) {
	v.GetVocabEntryByTextInvoked = true
	return v.GetVocabEntryByTextFn(text)
}

func (v *VocabEntryService) GetVocabEntryByID(id int) (*domain.VocabEntry, error) {
	v.GetVocabEntryByIDInvoked = true
	return v.GetVocabEntryByIDFn(id)
}

func (v *VocabEntryService) Reset() {
	v.GetVocabEntryByTextInvoked = false
	v.GetVocabEntryByIDInvoked = false
}

// VocabServiceConcurrency is a mock struct implementing service.Vocab interface.
// You can test concurrent execution of methods with this struct.
type VocabServiceConcurrency struct {
	CreateVocabFn          func(userID int) (*domain.Vocab, error)
	GetVocabEntryForWordFn func(word string) (*domain.VocabEntry, error)
	GetVocabEntryByIDFn    func(ID int) (*domain.VocabEntry, error)

	createVocabMu                  sync.Mutex
	CreateVocabInvoked             bool
	CreateVocabConcurrentlyInvoked bool
	createVocabConcurrencyCheck    map[int]struct{}

	getVocabEntryForWordMu                 sync.Mutex
	GetVocabEntryByTextInvoked             bool
	GetVocabEntryByTextConcurrentlyInvoked bool
	getVocabEntryForWordConcurrencyCheck   map[string]struct{}

	GetVocabEntryByIDInvoked bool
}

// NewVocabServiceConcurrency returns ready to use VocabServiceConcurrency.
func NewVocabServiceConcurrency(createVocabFn func(userID int) (*domain.Vocab, error),
	getVocabEntryForWordFn func(word string) (*domain.VocabEntry, error),
	getVocabEntryByID func(ID int) (*domain.VocabEntry, error)) *VocabServiceConcurrency {
	return &VocabServiceConcurrency{
		CreateVocabFn:                        createVocabFn,
		createVocabConcurrencyCheck:          make(map[int]struct{}),
		GetVocabEntryForWordFn:               getVocabEntryForWordFn,
		getVocabEntryForWordConcurrencyCheck: make(map[string]struct{}),
		GetVocabEntryByIDFn:                  getVocabEntryByID,
	}
}

// CreateVocab registers invocation of CreateVocab func and calls it.
// Also registers if it was called concurrently by the same user.
func (s *VocabServiceConcurrency) CreateVocab(userID int) (*domain.Vocab, error) {
	s.createVocabMu.Lock()
	s.CreateVocabInvoked = true
	_, ok := s.createVocabConcurrencyCheck[userID]
	if ok {
		s.CreateVocabConcurrentlyInvoked = true
	} else {
		s.createVocabConcurrencyCheck[userID] = struct{}{}
	}
	s.createVocabMu.Unlock()
	defer func() {
		s.createVocabMu.Lock()
		delete(s.createVocabConcurrencyCheck, userID)
		s.createVocabMu.Unlock()
	}()
	return s.CreateVocabFn(userID)
}

// GetVocabEntryByText registers invocation of GetVocabEntryByText func and calls it.
// Also registers if it was called concurrently for the same word.
func (s *VocabServiceConcurrency) GetVocabEntryByText(word string) (*domain.VocabEntry, error) {
	s.getVocabEntryForWordMu.Lock()
	s.GetVocabEntryByTextInvoked = true
	_, ok := s.getVocabEntryForWordConcurrencyCheck[word]
	if ok {
		s.GetVocabEntryByTextConcurrentlyInvoked = true
	} else {
		s.getVocabEntryForWordConcurrencyCheck[word] = struct{}{}
	}
	s.getVocabEntryForWordMu.Unlock()
	defer func() {
		s.getVocabEntryForWordMu.Lock()
		delete(s.getVocabEntryForWordConcurrencyCheck, word)
		s.getVocabEntryForWordMu.Unlock()
	}()
	return s.GetVocabEntryForWordFn(word)
}

// GetVocabEntryByID registers invocation of GetVocabEntryByID func and calls it.
func (s *VocabServiceConcurrency) GetVocabEntryByID(ID int) (*domain.VocabEntry, error) {
	s.GetVocabEntryByIDInvoked = true
	return s.GetVocabEntryByIDFn(ID)
}

// Reset resets functions invocation and concurrent functions invocation.
func (s *VocabServiceConcurrency) Reset() {
	s.CreateVocabInvoked = false
	s.CreateVocabConcurrentlyInvoked = false
	s.createVocabConcurrencyCheck = make(map[int]struct{})
	s.GetVocabEntryByTextInvoked = false
	s.GetVocabEntryByTextConcurrentlyInvoked = false
	s.getVocabEntryForWordConcurrencyCheck = make(map[string]struct{})
	s.GetVocabEntryByIDInvoked = false
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

func (_ Logger) WithField(_ string, _ interface{}) log.Logger {
	return Logger{}
}

func (_ Logger) WithFields(_ map[string]interface{}) log.Logger {
	return Logger{}
}
