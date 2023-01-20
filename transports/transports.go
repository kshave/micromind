package transports

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/micromind/services"

	"github.com/go-kit/kit/endpoint"
)

type QuoteRequest struct {
}

type QuoteResponse struct {
	Q string `json:"quote,omitempty"`
	A string `json:"author,omitempty"`
	E string `json:"error,omitempty"`
}

type QuestionRequest struct {
}

type QuestionResponse struct {
	Q string `json:"question,omitempty"`
	E string `json:"error,omitempty"`
}

type RootResponse struct {
	S string   `json:"service"`
	B string   `json:"bio"`
	E []string `json:"endpoints"`
	V string   `json:"version"`
}

func MakeRootEndpoint() endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		return RootResponse{
			"micromind",
			"a glimpse into mindfullness from the comfort of your command line!",
			[]string{"/quote", "/question"},
			"V0.1",
		}, nil

	}
}

func MakeQuoteEndpoint(svc services.ZenService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		q, a, err := svc.Quote()
		if err != nil {
			return QuoteResponse{q, a, err.Error()}, nil
		}
		return QuoteResponse{q, a, ""}, nil
	}
}

func MakeQuestionEndpoint(svc services.ZenService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		q, err := svc.Question()
		if err != nil {
			return QuestionResponse{q, err.Error()}, nil
		}
		return QuestionResponse{q, ""}, nil
	}
}

func DecodeQuoteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request QuoteRequest
	return request, nil
}

func DecodeQuestionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request QuestionRequest
	return request, nil
}

func DecodeQuestionResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var response QuestionResponse
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	return encoder.Encode(response)
}

func EndcodeRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}
