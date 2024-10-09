package iconik

import (
	"context"
	"encoding/json"
	"github.com/avast/retry-go"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/search"
	"github.com/rs/zerolog"
	"net/http"
	"strconv"
)

// Search makes a request to the POST iconik search endpoint.
func (a *API) Search(ctx context.Context, path string, payload []byte) (search.ResultsDTO, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, a.cfg.OperationTimeout)
	defer cancel()

	queryParams := make(map[string]string)
	queryParams[PerPage] = strconv.Itoa(a.cfg.PerPage)

	body, statusCode, err := a.req.Do(
		ctxTimeout,
		http.MethodPost,
		a.url+path,
		a.headers,
		queryParams,
		payload,
	)

	opDelay := a.cfg.OperationRetryDelay

	switch {
	case statusCode == nil:
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Msg("status code is nil")
		return search.ResultsDTO{}, err
	case *statusCode == http.StatusTooManyRequests,
		*statusCode == http.StatusInternalServerError,
		*statusCode == http.StatusServiceUnavailable,
		*statusCode == http.StatusGatewayTimeout:
		f := func() error {
			body, statusCode, err = a.req.Do(
				ctxTimeout,
				http.MethodPost,
				a.url+path,
				a.headers,
				queryParams,
				payload,
			)
			return err
		}
		onRetry := func(n uint, err error) {
			zerolog.Ctx(ctxTimeout).
				Debug().
				Err(err).
				Uint("attempt", n+1).
				Msg("retrying to search in iconik")
		}
		if *statusCode != http.StatusTooManyRequests {
			opDelay = 0
		}
		_ = retry.Do(
			f,
			retry.Attempts(a.cfg.OperationRetryAttempts),
			retry.Delay(opDelay),
			retry.OnRetry(onRetry),
		)
	case *statusCode == http.StatusForbidden:
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Int("status code", *statusCode).
			RawJSON("response", body).
			Msg("forbidden when searching assets")
		return search.ResultsDTO{}, domain.ErrForbidden
	case *statusCode == http.StatusUnauthorized:
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Int("status code", *statusCode).
			RawJSON("response", body).
			Msg("unauthorized when searching assets")
		return search.ResultsDTO{}, domain.Err401Search
	case *statusCode != http.StatusOK:
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			RawJSON("response", body).
			Int("status code", *statusCode).
			Msg("status code unexpected")
		return search.ResultsDTO{}, domain.ErrInternalError
	}

	if err != nil {
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Msg("error when using search")
		return search.ResultsDTO{}, err
	}

	var res search.Results
	if err = json.Unmarshal(body, &res); err != nil {
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Msg("error unmarshalling body")
		return search.ResultsDTO{}, err
	}

	return res.ToResultsDTO(), nil
}
