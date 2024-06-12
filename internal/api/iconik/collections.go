package iconik

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/avast/retry-go"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/assets/collections"
	"github.com/rs/zerolog"
	"net/http"
)

// GetCollectionContents makes a request to the GET iconik collection contents endpoint.
func (a *API) GetCollectionContents(ctx context.Context, path, collectionID string) (collections.ContentsDTO, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, a.cfg.OperationTimeout)
	defer cancel()
	zerolog.Ctx(ctxTimeout).Info().Msg("getting collection " + collectionID + " contents from iconik")

	body, statusCode, err := a.req.Do(
		ctxTimeout,
		http.MethodGet,
		a.url+path+"/"+collectionID+"/contents/",
		a.headers,
		nil,
		nil,
	)

	opDelay := a.cfg.OperationRetryDelay

	switch {
	case errors.Is(err, domain.ErrTransformingHeaderValue) || errors.Is(err, domain.ErrTransformingHeaderKey):
		return collections.ContentsDTO{}, err
	case statusCode == nil:
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Msg("status code is nil")
		return collections.ContentsDTO{}, err
	case *statusCode == http.StatusTooManyRequests,
		*statusCode == http.StatusInternalServerError,
		*statusCode == http.StatusServiceUnavailable,
		*statusCode == http.StatusGatewayTimeout:
		f := func() error {
			body, statusCode, err = a.req.Do(
				ctxTimeout,
				http.MethodGet,
				a.url+path+"/"+collectionID+"/contents/",
				a.headers,
				nil,
				nil,
			)
			return err
		}
		onRetry := func(n uint, err error) {
			zerolog.Ctx(ctxTimeout).
				Debug().
				Err(err).
				Uint("attempt", n+1).
				Msg("retrying to get collection contents from iconik")
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
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Int("status code", *statusCode).
			Msg("status code unexpected")
		return collections.ContentsDTO{}, err
	}

	if err != nil {
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Msg("error getting collection contents")
		return collections.ContentsDTO{}, err
	}

	var res collections.Contents
	if err = json.Unmarshal(body, &res); err != nil {
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Msg("error unmarshalling body")
		return collections.ContentsDTO{}, err
	}

	return res.ToContentsDTO(), nil
}
