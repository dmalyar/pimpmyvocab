package service

import (
	"github.com/dmalyar/pimpmyvocab/domain"
	"github.com/dmalyar/pimpmyvocab/mock"
	"sync"
	"testing"
)

func TestVocabConcurrent_CreateVocab(t *testing.T) {
	mockedService := mock.NewVocabServiceConcurrency(func(userID int) (*domain.Vocab, error) {
		return &domain.Vocab{UserID: userID}, nil
	})
	testService := NewConcurrentVocab(mockedService)
	var wg sync.WaitGroup
	for i := 0; i < 300; i++ {
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

func createVocab(wg *sync.WaitGroup, s *ConcurrentVocab, userID int) {
	wg.Add(1)
	defer wg.Done()
	_, _ = s.CreateVocab(userID)
}
