package service

import (
	"github.com/dmalyar/pimpmyvocab/domain"
	"github.com/dmalyar/pimpmyvocab/mock"
	"sync"
	"testing"
)

func TestVocabConcurrent_CreateVocab(t *testing.T) {
	mockedService := mockWrappedService()
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
	if mockedService.CreateVocabConcurrentlyInvoked {
		t.Error("Underlying service was invoked concurrently by the same user")
	}
}

func TestConcurrentVocab_CheckEntryInUserVocab(t *testing.T) {
	mockedService := mockWrappedService()
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
	if mockedService.CheckEntryInUserVocabConcurrentlyInvoked {
		t.Error("Underlying service was invoked concurrently by the same user")
	}
}

func TestConcurrentVocab_AddEntryToUserVocab(t *testing.T) {
	mockedService := mockWrappedService()
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
	if mockedService.AddEntryToUserVocabConcurrentlyInvoked {
		t.Error("Underlying service was invoked concurrently by the same user")
	}
}

func TestConcurrentVocab_RemoveEntryFromUserVocab(t *testing.T) {
	mockedService := mockWrappedService()
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
	if mockedService.RemoveEntryFromUserVocabConcurrentlyInvoked {
		t.Error("Underlying service was invoked concurrently by the same user")
	}
}

func TestConcurrentVocab_GetVocabEntryByText(t *testing.T) {
	mockedService := mockWrappedService()
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
	if mockedService.GetVocabEntryByTextConcurrentlyInvoked {
		t.Error("Underlying service was invoked concurrently for the same text")
	}
}

func TestConcurrentVocab_GetVocabEntryByID(t *testing.T) {
	mockedService := mockWrappedService()
	testService := NewConcurrentVocab(mockedService)
	_, err := testService.GetVocabEntryByID(1)
	if err != nil {
		t.Errorf("Error calling method: %s", err)
	}
	if !mockedService.GetVocabEntryByIDInvoked {
		t.Error("GetVocabEntryByID wasn't invoked")
	}
}

func mockWrappedService() *mock.VocabServiceConcurrency {
	return mock.NewVocabServiceConcurrency(
		func(userID int) (*domain.Vocab, error) {
			return &domain.Vocab{UserID: userID}, nil
		},
		func(text string) (*domain.VocabEntry, error) {
			return &domain.VocabEntry{Text: text}, nil
		},
		func(ID int) (*domain.VocabEntry, error) {
			return &domain.VocabEntry{ID: ID}, nil
		},
		func(entryID, userID int) (bool, error) {
			return true, nil
		},
		func(entryID, userID int) error {
			return nil
		},
		func(entryID, userID int) error {
			return nil
		},
	)
}

func createVocab(wg *sync.WaitGroup, s *ConcurrentVocab, userID int) {
	_, _ = s.CreateVocab(userID)
	wg.Done()
}

func checkEntryInUserVocab(wg *sync.WaitGroup, s *ConcurrentVocab, entryID int, userID int) {
	_, _ = s.CheckEntryInUserVocab(entryID, userID)
	wg.Done()
}

func addEntryToUserVocab(wg *sync.WaitGroup, s *ConcurrentVocab, entryID int, userID int) {
	_ = s.AddEntryToUserVocab(entryID, userID)
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
