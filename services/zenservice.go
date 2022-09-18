package services

import (
	"errors"
	"log"

	db "github.com/micromind/repositories"
)

// ZenService
type ZenService interface {
	Quote() (string, string, error)
	Question() (string, error)
}

type zenService struct {
	zr db.ZenRepository
}

func NewInstanceOfZenService(zr db.ZenRepository) zenService {
	return zenService{zr}
}

func (svc zenService) Quote() (string, string, error) {
	quote, author, err := svc.zr.GetRandomQuote()
	if err != nil {
		log.Printf("Failed to get random quote: %v", err)
		return "", "", errBad
	}
	return quote, author, nil
}

func (svc zenService) Question() (string, error) {
	question, err := svc.zr.GetRandomQuestion()
	if err != nil {
		log.Printf("Failed to get a random question: %v", err)
		return "", errBad
	}
	return question, nil
}

// ServiceMiddleware is a chainable behavior modifier for ZenService
type ServiceMiddleware func(ZenService) ZenService

// errBad is returned when something has gone wrong!
var errBad = errors.New("i have a bad feeling about this")
