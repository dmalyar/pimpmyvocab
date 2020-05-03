package service

import (
	"github.com/dmalyar/pimpmyvocab/domain"
	"github.com/dmalyar/pimpmyvocab/mock"
	"sync"
	"testing"
)

func TestVocabConcurrent_CreateVocab(t *testing.T) {
	mockedService := mock.NewVocabServiceConcurrencyCheck()
	mockedService.CreateVocabFn = func(userID int) (*domain.Vocab, error) {
		return &domain.Vocab{UserID: userID}, nil
	}
	testService := NewConcurrentVocab(mockedService)
	var wg sync.WaitGroup
	for i := 0; i < 300; i++ {
		wg.Add(3)
		go createVocab(&wg, testService, 1)
		go createVocab(&wg, testService, 2)
		go createVocab(&wg, testService, 3)
	}
	wg.Wait()
	if !mockedService.CreateVocabInvoked {
		t.Error("CreateVocab wasn't invoked")
	}
	if mockedService.UserIDConcurrentlyInvoked {
		t.Error("Underlying service was invoked concurrently by the same user")
	}
}

func TestVocabConcurrent_ClearUserVocab(t *testing.T) {
	mockedService := mock.NewVocabServiceConcurrencyCheck()
	mockedService.ClearUserVocabFn = func(userID int) error {
		return nil
	}
	testService := NewConcurrentVocab(mockedService)
	var wg sync.WaitGroup
	for i := 0; i < 300; i++ {
		wg.Add(3)
		go clearVocab(&wg, testService, 1)
		go clearVocab(&wg, testService, 2)
		go clearVocab(&wg, testService, 3)
	}
	wg.Wait()
	if !mockedService.ClearUserVocabInvoked {
		t.Error("ClearUserVocab wasn't invoked")
	}
	if mockedService.UserIDConcurrentlyInvoked {
		t.Error("Underlying service was invoked concurrently by the same user")
	}
}

func TestConcurrentVocab_AddEntryToUserVocab(t *testing.T) {
	mockedService := mock.NewVocabServiceConcurrencyCheck()
	mockedService.AddEntryToUserVocabFn = func(entryID, userID int) error {
		return nil
	}
	testService := NewConcurrentVocab(mockedService)
	var wg sync.WaitGroup
	for i := 0; i < 300; i++ {
		wg.Add(3)
		go addEntryToUserVocab(&wg, testService, 1, 1)
		go addEntryToUserVocab(&wg, testService, 2, 2)
		go addEntryToUserVocab(&wg, testService, 3, 1)
	}
	wg.Wait()
	if !mockedService.AddEntryToUserVocabInvoked {
		t.Error("AddEntryToUserVocab wasn't invoked")
	}
	if mockedService.UserIDConcurrentlyInvoked {
		t.Error("Underlying service was invoked concurrently by the same user")
	}
}

func TestConcurrentVocab_CheckEntryInUserVocab(t *testing.T) {
	mockedService := mock.NewVocabServiceConcurrencyCheck()
	mockedService.CheckEntryInUserVocabFn = func(entryID, userID int) (bool, error) {
		return true, nil
	}
	testService := NewConcurrentVocab(mockedService)
	var wg sync.WaitGroup
	for i := 0; i < 300; i++ {
		wg.Add(3)
		go checkEntryInUserVocab(&wg, testService, 1, 1)
		go checkEntryInUserVocab(&wg, testService, 2, 2)
		go checkEntryInUserVocab(&wg, testService, 3, 1)
	}
	wg.Wait()
	if !mockedService.CheckEntryInUserVocabInvoked {
		t.Error("CheckEntryInUserVocab wasn't invoked")
	}
	if mockedService.UserIDConcurrentlyInvoked {
		t.Error("Underlying service was invoked concurrently by the same user")
	}
}

func TestConcurrentVocab_RemoveEntryFromUserVocab(t *testing.T) {
	mockedService := mock.NewVocabServiceConcurrencyCheck()
	mockedService.RemoveEntryFromUserVocabFn = func(entryID, userID int) error {
		return nil
	}
	testService := NewConcurrentVocab(mockedService)
	var wg sync.WaitGroup
	for i := 0; i < 300; i++ {
		wg.Add(3)
		go removeEntryFromUserVocab(&wg, testService, 1, 1)
		go removeEntryFromUserVocab(&wg, testService, 2, 2)
		go removeEntryFromUserVocab(&wg, testService, 3, 1)
	}
	wg.Wait()
	if !mockedService.RemoveEntryFromUserVocabInvoked {
		t.Error("RemoveEntryFromUserVocab wasn't invoked")
	}
	if mockedService.UserIDConcurrentlyInvoked {
		t.Error("Underlying service was invoked concurrently by the same user")
	}
}

func TestConcurrentVocab_GetEntriesFromUserVocab(t *testing.T) {
	mockedService := mock.NewVocabServiceConcurrencyCheck()
	mockedService.GetEntriesFromUserVocabFn = func(userID int) ([]*domain.VocabEntry, error) {
		return []*domain.VocabEntry{}, nil
	}
	testService := NewConcurrentVocab(mockedService)
	var wg sync.WaitGroup
	for i := 0; i < 300; i++ {
		wg.Add(3)
		go getEntriesFromUserVocab(&wg, testService, 1)
		go getEntriesFromUserVocab(&wg, testService, 2)
		go getEntriesFromUserVocab(&wg, testService, 3)
	}
	wg.Wait()
	if !mockedService.GetEntriesFromUserVocabInvoked {
		t.Error("GetEntriesFromUserVocab wasn't invoked")
	}
	if mockedService.UserIDConcurrentlyInvoked {
		t.Error("Underlying service was invoked concurrently for the same text")
	}
}

func TestConcurrentVocab_GetVocabEntryByText(t *testing.T) {
	mockedService := mock.NewVocabServiceConcurrencyCheck()
	mockedService.GetVocabEntryByTextFn = func(text string) (*domain.VocabEntry, error) {
		return &domain.VocabEntry{Text: text}, nil
	}
	testService := NewConcurrentVocab(mockedService)
	var wg sync.WaitGroup
	for i := 0; i < 300; i++ {
		wg.Add(3)
		go getVocabEntryByText(&wg, testService, "text")
		go getVocabEntryByText(&wg, testService, "another")
		go getVocabEntryByText(&wg, testService, "one")
	}
	wg.Wait()
	if !mockedService.GetVocabEntryByTextInvoked {
		t.Error("GetVocabEntryByText wasn't invoked")
	}
	if mockedService.TextConcurrentlyInvoked {
		t.Error("Underlying service was invoked concurrently for the same text")
	}
}

func TestConcurrentVocab_GetVocabEntryByID(t *testing.T) {
	mockedService := mock.NewVocabServiceConcurrencyCheck()
	mockedService.GetVocabEntryByIDFn = func(id int) (*domain.VocabEntry, error) {
		return &domain.VocabEntry{ID: id}, nil
	}
	testService := NewConcurrentVocab(mockedService)
	_, err := testService.GetVocabEntryByID(1)
	if err != nil {
		t.Errorf("Error calling method: %s", err)
	}
	if !mockedService.GetVocabEntryByIDInvoked {
		t.Error("GetVocabEntryByID wasn't invoked")
	}
}

func TestConcurrentVocab_SyncByUserIDTest(t *testing.T) {
	mockedService := mock.NewVocabServiceConcurrencyCheck()
	mockedService.CreateVocabFn = func(userID int) (*domain.Vocab, error) {
		return &domain.Vocab{UserID: userID}, nil
	}
	mockedService.ClearUserVocabFn = func(userID int) error {
		return nil
	}
	mockedService.AddEntryToUserVocabFn = func(entryID, userID int) error {
		return nil
	}
	mockedService.CheckEntryInUserVocabFn = func(entryID, userID int) (bool, error) {
		return true, nil
	}
	mockedService.RemoveEntryFromUserVocabFn = func(entryID, userID int) error {
		return nil
	}
	mockedService.GetEntriesFromUserVocabFn = func(userID int) ([]*domain.VocabEntry, error) {
		return []*domain.VocabEntry{}, nil
	}
	testService := NewConcurrentVocab(mockedService)
	var wg sync.WaitGroup
	for i := 0; i < 300; i++ {
		wg.Add(6)
		go createVocab(&wg, testService, 1)
		go clearVocab(&wg, testService, 1)
		go addEntryToUserVocab(&wg, testService, 1, 1)
		go checkEntryInUserVocab(&wg, testService, 1, 1)
		go removeEntryFromUserVocab(&wg, testService, 1, 1)
		go getEntriesFromUserVocab(&wg, testService, 1)
	}
	wg.Wait()
	if mockedService.UserIDConcurrentlyInvoked {
		t.Error("Underlying service was invoked concurrently by the same user")
	}
}

func createVocab(wg *sync.WaitGroup, s *ConcurrentVocab, userID int) {
	_, _ = s.CreateVocab(userID)
	wg.Done()
}

func clearVocab(wg *sync.WaitGroup, s *ConcurrentVocab, userID int) {
	_ = s.ClearUserVocab(userID)
	wg.Done()
}

func addEntryToUserVocab(wg *sync.WaitGroup, s *ConcurrentVocab, entryID int, userID int) {
	_ = s.AddEntryToUserVocab(entryID, userID)
	wg.Done()
}

func checkEntryInUserVocab(wg *sync.WaitGroup, s *ConcurrentVocab, entryID int, userID int) {
	_, _ = s.CheckEntryInUserVocab(entryID, userID)
	wg.Done()
}

func getEntriesFromUserVocab(wg *sync.WaitGroup, s *ConcurrentVocab, userID int) {
	_, _ = s.GetEntriesFromUserVocab(userID)
	wg.Done()
}

func removeEntryFromUserVocab(wg *sync.WaitGroup, s *ConcurrentVocab, entryID int, userID int) {
	_ = s.RemoveEntryFromUserVocab(entryID, userID)
	wg.Done()
}

func getVocabEntryByText(wg *sync.WaitGroup, s *ConcurrentVocab, word string) {
	_, _ = s.GetVocabEntryByText(word)
	wg.Done()
}
