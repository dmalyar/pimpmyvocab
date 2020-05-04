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

	ClearVocabByUserIDFn      func(userID int) error
	ClearVocabByUserIDInvoked bool

	AddVocabEntryFn      func(vocab *domain.VocabEntry) (*domain.VocabEntry, error)
	AddVocabEntryInvoked bool

	GetVocabEntryByTextFn      func(text string) (*domain.VocabEntry, error)
	GetVocabEntryByTextInvoked bool

	GetVocabEntryByIDFn      func(id int) (*domain.VocabEntry, error)
	GetVocabEntryByIDInvoked bool

	AddEntryToUserVocabFn      func(entryID, userID int) error
	AddEntryToUserVocabInvoked bool

	CheckEntryInUserVocabFn      func(entryID, userID int) (bool, error)
	CheckEntryInUserVocabInvoked bool

	GetEntryIDsByUserIDFn      func(userID int) ([]int, error)
	GetEntryIDsByUserIDInvoked bool

	GetEntriesByUserIDFn      func(userID int) ([]*domain.VocabEntry, error)
	GetEntriesByUserIDInvoked bool

	RemoveEntryFromUserVocabFn      func(entryID, userID int) error
	RemoveEntryFromUserVocabInvoked bool
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

// ClearVocabByUserID registers invocation of ClearVocabByUserID func and calls it.
func (r *VocabRepo) ClearVocabByUserID(userID int) error {
	r.ClearVocabByUserIDInvoked = true
	return r.ClearVocabByUserIDFn(userID)
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
func (r *VocabRepo) GetVocabEntryByID(id int) (*domain.VocabEntry, error) {
	r.GetVocabEntryByIDInvoked = true
	return r.GetVocabEntryByIDFn(id)
}

// AddEntryToUserVocab registers invocation of AddEntryToUserVocab func and calls it.
func (r *VocabRepo) AddEntryToUserVocab(entryID, userID int) error {
	r.AddEntryToUserVocabInvoked = true
	return r.AddEntryToUserVocabFn(entryID, userID)
}

// CheckEntryInUserVocab registers invocation of CheckEntryInUserVocab func and calls it.
func (r *VocabRepo) CheckEntryInUserVocab(entryID, userID int) (bool, error) {
	r.CheckEntryInUserVocabInvoked = true
	return r.CheckEntryInUserVocabFn(entryID, userID)
}

// GetEntryIDsByUserID registers invocation of GetEntryIDsByUserID func and calls it.
func (r *VocabRepo) GetEntryIDsByUserID(userID int) ([]int, error) {
	r.GetEntryIDsByUserIDInvoked = true
	return r.GetEntryIDsByUserIDFn(userID)
}

// GetEntriesByUserID registers invocation of GetEntriesByUserID func and calls it.
func (r *VocabRepo) GetEntriesByUserID(userID int) ([]*domain.VocabEntry, error) {
	r.GetEntriesByUserIDInvoked = true
	return r.GetEntriesByUserIDFn(userID)
}

// RemoveEntryFromUserVocab registers invocation of RemoveEntryFromUserVocab func and calls it.
func (r *VocabRepo) RemoveEntryFromUserVocab(entryID, userID int) error {
	r.RemoveEntryFromUserVocabInvoked = true
	return r.RemoveEntryFromUserVocabFn(entryID, userID)
}

// Reset resets functions invocation.
func (r *VocabRepo) Reset() {
	r.AddVocabInvoked = false
	r.GetVocabByUserIDInvoked = false
	r.ClearVocabByUserIDInvoked = false
	r.AddVocabEntryInvoked = false
	r.GetVocabEntryByTextInvoked = false
	r.GetVocabEntryByIDInvoked = false
	r.AddEntryToUserVocabInvoked = false
	r.CheckEntryInUserVocabInvoked = false
	r.GetEntriesByUserIDInvoked = false
	r.RemoveEntryFromUserVocabInvoked = false
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

// VocabServiceConcurrencyCheck is a mock struct implementing service.Vocab interface.
// You can test concurrent execution of methods with this struct.
type VocabServiceConcurrencyCheck struct {
	CreateVocabFn      func(userID int) (*domain.Vocab, error)
	CreateVocabInvoked bool

	ClearUserVocabFn      func(userID int) error
	ClearUserVocabInvoked bool

	AddEntryToUserVocabFn      func(entryID, userID int) error
	AddEntryToUserVocabInvoked bool

	CheckEntryInUserVocabFn      func(entryID, userID int) (bool, error)
	CheckEntryInUserVocabInvoked bool

	GetRandomEntryFromUserVocabFn      func(userId, previousEntryID int) (*domain.VocabEntry, error)
	GetRandomEntryFromUserVocabInvoked bool

	GetEntriesFromUserVocabFn      func(userID int) ([]*domain.VocabEntry, error)
	GetEntriesFromUserVocabInvoked bool

	RemoveEntryFromUserVocabFn      func(entryID, userID int) error
	RemoveEntryFromUserVocabInvoked bool

	GetVocabEntryByTextFn      func(text string) (*domain.VocabEntry, error)
	GetVocabEntryByTextInvoked bool

	GetVocabEntryByIDFn      func(ID int) (*domain.VocabEntry, error)
	GetVocabEntryByIDInvoked bool

	textConcurrencyCheckMu  sync.Mutex
	textConcurrencyCheck    map[string]struct{}
	TextConcurrentlyInvoked bool

	userIDConcurrencyCheckMu  sync.Mutex
	userIDConcurrencyCheck    map[int]struct{}
	UserIDConcurrentlyInvoked bool
}

// NewVocabServiceConcurrencyCheck returns ready to use VocabServiceConcurrencyCheck.
func NewVocabServiceConcurrencyCheck() *VocabServiceConcurrencyCheck {
	return &VocabServiceConcurrencyCheck{
		textConcurrencyCheck:   make(map[string]struct{}),
		userIDConcurrencyCheck: make(map[int]struct{}),
	}
}

// CreateVocab registers invocation of CreateVocab func and calls it.
// Also registers if it was called concurrently by the same user.
func (s *VocabServiceConcurrencyCheck) CreateVocab(userID int) (*domain.Vocab, error) {
	s.startWorkSyncedByUserID(userID, &s.CreateVocabInvoked)
	defer s.endWorkSyncedByUserID(userID)
	return s.CreateVocabFn(userID)
}

// ClearUserVocab registers invocation of ClearUserVocab func and calls it.
// Also registers if it was called concurrently by the same user.
func (s *VocabServiceConcurrencyCheck) ClearUserVocab(userID int) error {
	s.startWorkSyncedByUserID(userID, &s.ClearUserVocabInvoked)
	defer s.endWorkSyncedByUserID(userID)
	return s.ClearUserVocabFn(userID)
}

// AddEntryToUserVocab registers invocation of AddEntryToUserVocab func and calls it.
// Also registers if it was called concurrently by the same user.
func (s *VocabServiceConcurrencyCheck) AddEntryToUserVocab(entryID, userID int) error {
	s.startWorkSyncedByUserID(userID, &s.AddEntryToUserVocabInvoked)
	defer s.endWorkSyncedByUserID(userID)
	return s.AddEntryToUserVocabFn(entryID, userID)
}

// CheckEntryInUserVocab registers invocation of CheckEntryInUserVocab func and calls it.
// Also registers if it was called concurrently by the same user.
func (s *VocabServiceConcurrencyCheck) CheckEntryInUserVocab(entryID, userID int) (bool, error) {
	s.startWorkSyncedByUserID(userID, &s.CheckEntryInUserVocabInvoked)
	defer s.endWorkSyncedByUserID(userID)
	return s.CheckEntryInUserVocabFn(entryID, userID)
}

// GetRandomEntryFromUserVocab registers invocation of GetRandomEntryFromUserVocab func and calls it.
// Also registers if it was called concurrently by the same user.
func (s *VocabServiceConcurrencyCheck) GetRandomEntryFromUserVocab(userID, previousEntryID int) (*domain.VocabEntry, error) {
	s.startWorkSyncedByUserID(userID, &s.GetRandomEntryFromUserVocabInvoked)
	defer s.endWorkSyncedByUserID(userID)
	return s.GetRandomEntryFromUserVocabFn(userID, previousEntryID)
}

// GetEntriesByUserID registers invocation of GetEntriesByUserID func and calls it.
// Also registers if it was called concurrently by the same user.
func (s *VocabServiceConcurrencyCheck) GetEntriesFromUserVocab(userID int) ([]*domain.VocabEntry, error) {
	s.startWorkSyncedByUserID(userID, &s.GetEntriesFromUserVocabInvoked)
	defer s.endWorkSyncedByUserID(userID)
	return s.GetEntriesFromUserVocabFn(userID)
}

// RemoveEntryFromUserVocab registers invocation of RemoveEntryFromUserVocab func and calls it.
// Also registers if it was called concurrently by the same user.
func (s *VocabServiceConcurrencyCheck) RemoveEntryFromUserVocab(entryID, userID int) error {
	s.startWorkSyncedByUserID(userID, &s.RemoveEntryFromUserVocabInvoked)
	defer s.endWorkSyncedByUserID(userID)
	return s.RemoveEntryFromUserVocabFn(entryID, userID)
}

// GetVocabEntryByText registers invocation of GetVocabEntryByText func and calls it.
// Also registers if it was called concurrently for the same text.
func (s *VocabServiceConcurrencyCheck) GetVocabEntryByText(text string) (*domain.VocabEntry, error) {
	s.startWorkSyncedByText(text, &s.GetVocabEntryByTextInvoked)
	defer s.endWorkSyncedByText(text)
	return s.GetVocabEntryByTextFn(text)
}

// GetVocabEntryByID registers invocation of GetVocabEntryByID func and calls it.
func (s *VocabServiceConcurrencyCheck) GetVocabEntryByID(ID int) (*domain.VocabEntry, error) {
	s.GetVocabEntryByIDInvoked = true
	return s.GetVocabEntryByIDFn(ID)
}

func (s *VocabServiceConcurrencyCheck) startWorkSyncedByUserID(userID int, invocation *bool) {
	s.userIDConcurrencyCheckMu.Lock()
	*invocation = true
	_, ok := s.userIDConcurrencyCheck[userID]
	if ok {
		s.UserIDConcurrentlyInvoked = true
	} else {
		s.userIDConcurrencyCheck[userID] = struct{}{}
	}
	s.userIDConcurrencyCheckMu.Unlock()
}

func (s *VocabServiceConcurrencyCheck) endWorkSyncedByUserID(userID int) {
	s.userIDConcurrencyCheckMu.Lock()
	delete(s.userIDConcurrencyCheck, userID)
	s.userIDConcurrencyCheckMu.Unlock()
}

func (s *VocabServiceConcurrencyCheck) startWorkSyncedByText(text string, invocation *bool) {
	s.textConcurrencyCheckMu.Lock()
	*invocation = true
	_, ok := s.textConcurrencyCheck[text]
	if ok {
		s.UserIDConcurrentlyInvoked = true
	} else {
		s.textConcurrencyCheck[text] = struct{}{}
	}
	s.textConcurrencyCheckMu.Unlock()
}

func (s *VocabServiceConcurrencyCheck) endWorkSyncedByText(text string) {
	s.textConcurrencyCheckMu.Lock()
	delete(s.textConcurrencyCheck, text)
	s.textConcurrencyCheckMu.Unlock()
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
