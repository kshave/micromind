package services

import (
	"testing"
)

func TestZenService(t *testing.T) {
	t.Run("Quote returns a random Quote", func(t *testing.T) {
		spyRepo := &SpyRepo{}
		zs := NewInstanceOfZenService(spyRepo)

		wantQ := "I have a bad feeling about this"
		wantA := "Obi-Wan"

		gotQ, gotA, error := zs.Quote()

		if error != nil {
			t.Errorf("error %s", error)
		}

		if (gotQ != wantQ) || (gotA != wantA) {
			t.Errorf("gotQuote: %s wantQuote: %s gotAuthor: %s wantAuthor: %s", gotQ, wantQ, gotA, wantA)
		}

		if spyRepo.Calls > 1 {
			t.Errorf("Too many calls to repo, want 1 but got %d", spyRepo.Calls)
		}
	})

	t.Run("Question returns a random Question", func(t *testing.T) {
		spyRepo := &SpyRepo{}
		zs := NewInstanceOfZenService(spyRepo)

		wantQ := "Is this a random question?"

		gotQ, error := zs.Question()

		if error != nil {
			t.Errorf("error %s", error)
		}

		if gotQ != wantQ {
			t.Errorf("gotQuestion: %s wantQuestion: %s", gotQ, wantQ)
		}

		if spyRepo.Calls > 1 {
			t.Errorf("Too many calls to repo, want 1 but got %d", spyRepo.Calls)
		}
	})
}

type SpyRepo struct {
	Calls int
}

func (s *SpyRepo) GetRandomQuestion() (string, error) {
	s.Calls++
	return "Is this a random question?", nil
}

func (s *SpyRepo) GetRandomQuote() (string, string, error) {
	s.Calls++
	return "I have a bad feeling about this", "Obi-Wan", nil
}
