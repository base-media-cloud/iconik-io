package iconik

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/avast/retry-go"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/search"
	"github.com/rs/zerolog"
	"net/http"
	"strconv"
)

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
	case errors.Is(err, domain.ErrTransformingHeaderValue) || errors.Is(err, domain.ErrTransformingHeaderKey):
		return search.ResultsDTO{}, err
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
				Msg("retrying to update asset in iconik")
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
	case *statusCode != http.StatusOK:
		return search.ResultsDTO{}, err
	}

	if err != nil {
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Msg("error updating asset")
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
