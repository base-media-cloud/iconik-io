package api

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/rs/zerolog"
	"golang.org/x/sync/syncmap"
)

// Request is a struct that wraps the http Client.
type Request struct {
	client *http.Client
}

// New is a function that returns a new instance of the Request struct.
func New(client *http.Client) *Request {
	return &Request{
		client: client,
	}
}

// Do is a helper function that makes a http request.
func (r *Request) Do(
	ctx context.Context,
	method,
	url string,
	headers *syncmap.Map,
	queryParams map[string]string,
	payload []byte,
) ([]byte, *int, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(payload))
	if err != nil {
		zerolog.Ctx(ctx).Error().
			Err(err).
			Msg("error creating http request")
		return nil, nil, errors.New("error creating http request")
	}

	headers.Range(func(key, value interface{}) bool {
		val, ok := value.(string)
		if !ok {
			err = errors.New("cannot transform header value to string")
			return false
		}
		k, ok := key.(string)
		if !ok {
			err = errors.New("cannot transform header key to string")
			return false
		}
		req.Header.Set(k, val)
		return true
	})
	if err != nil {
		zerolog.Ctx(ctx).Error().
			Err(err).
			Msg("error setting request headers")
		return nil, nil, err
	}

	q := req.URL.Query()
	for key, value := range queryParams {
		q.Set(key, value)
	}
	req.URL.RawQuery = q.Encode()

	res, err := r.client.Do(req)
	if err != nil {
		zerolog.Ctx(ctx).Error().
			Err(err).
			Msg("error sending http request")
		return nil, nil, errors.New("error sending http request")
	}
	defer func() {
		err = errors.Join(err, res.Body.Close())
	}()

	if method == http.MethodDelete && res.StatusCode == http.StatusNoContent {
		return nil, &res.StatusCode, nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		zerolog.Ctx(ctx).Error().
			Err(err).
			Msg("error reading response body")
		return nil, nil, errors.New("error reading response body")
	}

	return body, &res.StatusCode, nil
}
