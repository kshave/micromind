// Based on tests from: https://github.com/go-kit/kit/blob/master/endpoint/endpoint_example_test.go
package transports

import (
	"context"
	"fmt"
	"testing"
)

var (
	ctx = context.Background()
	req = struct{}{}
)

func TestEndpoints(t *testing.T) {
	t.Run("MakeQuoteEndpoint returns a valid endpoint", func(t *testing.T) {
		// Arrange
		spyService := &SpySerivce{}

		// Act
		e := MakeQuoteEndpoint(spyService)

		// Assert
		if resp, err := e(ctx, req); err != nil {
			t.FailNow()
		} else {
			fmt.Println(resp)
		}

	})

	t.Run("MakeQuestionEndpoint returns a valid endpoint", func(t *testing.T) {
		// Arrange
		spyService := &SpySerivce{}

		// Act
		e := MakeQuestionEndpoint(spyService)

		// Assert
		if resp, err := e(ctx, req); err != nil {
			t.FailNow()
		} else {
			fmt.Println(resp)
		}

	})
}

type SpySerivce struct {
	Calls int
}

func (s *SpySerivce) Question() (string, error) {
	s.Calls++
	return "Is this a random question?", nil
}

func (s *SpySerivce) Quote() (string, string, error) {
	s.Calls++
	return "I have a bad feeling about this", "Obi-Wan", nil
}
