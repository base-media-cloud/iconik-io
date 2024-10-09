package iconik

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/assets/collections"
	"github.com/rs/zerolog"
	"net/http"
)

// GetCollectionContents makes a request to the GET iconik collection contents endpoint.
func (a *API) GetCollectionContents(ctx context.Context, path, collectionID string, queryParams map[string]string) (collections.ContentsDTO, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, a.cfg.OperationTimeout)
	defer cancel()

	body, statusCode, err := a.req.Do(
		ctxTimeout,
		http.MethodGet,
		fmt.Sprintf("%v%v%v/contents/", a.url, path, collectionID),
		a.headers,
		queryParams,
		nil,
	)

	opDelay := a.cfg.OperationRetryDelay

	switch {
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
				fmt.Sprintf("%v%v%v/contents/", a.url, path, collectionID),
				a.headers,
				queryParams,
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
	case *statusCode == http.StatusForbidden:
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Int("status code", *statusCode).
			RawJSON("response", body).
			Msg("forbidden when getting collection contents")
		return collections.ContentsDTO{}, domain.ErrForbidden
	case *statusCode == http.StatusUnauthorized:
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Int("status code", *statusCode).
			RawJSON("response", body).
			Msg("unauthorized when getting collection contents")
		return collections.ContentsDTO{}, domain.Err401CollectionContents
	case *statusCode != http.StatusOK:
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			RawJSON("response", body).
			Int("status code", *statusCode).
			Msg("status code unexpected")
		return collections.ContentsDTO{}, domain.ErrInternalError
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

// GetCollection makes a request to the GET iconik collection endpoint.
func (a *API) GetCollection(ctx context.Context, path, collectionID string) (collections.CollectionDTO, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, a.cfg.OperationTimeout)
	defer cancel()

	body, statusCode, err := a.req.Do(
		ctxTimeout,
		http.MethodGet,
		fmt.Sprintf("%v%v%v/", a.url, path, collectionID),
		a.headers,
		nil,
		nil,
	)

	opDelay := a.cfg.OperationRetryDelay

	switch {
	case statusCode == nil:
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Msg("status code is nil")
		return collections.CollectionDTO{}, err
	case *statusCode == http.StatusTooManyRequests,
		*statusCode == http.StatusInternalServerError,
		*statusCode == http.StatusServiceUnavailable,
		*statusCode == http.StatusGatewayTimeout:
		f := func() error {
			body, statusCode, err = a.req.Do(
				ctxTimeout,
				http.MethodGet,
				fmt.Sprintf("%v%v%v/", a.url, path, collectionID),
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
				Msg("retrying to get collection from iconik")
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
			Msg("forbidden when getting collection")
		return collections.CollectionDTO{}, domain.ErrForbidden
	case *statusCode == http.StatusUnauthorized:
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Int("status code", *statusCode).
			RawJSON("response", body).
			Msg("unauthorized when getting collection")
		return collections.CollectionDTO{}, domain.Err401Collection
	case *statusCode != http.StatusOK:
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			RawJSON("response", body).
			Int("status code", *statusCode).
			Msg("status code unexpected")
		return collections.CollectionDTO{}, domain.ErrInternalError
	}

	if err != nil {
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Msg("error getting collection")
		return collections.CollectionDTO{}, err
	}

	var res collections.Collection
	if err = json.Unmarshal(body, &res); err != nil {
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Msg("error unmarshalling body")
		return collections.CollectionDTO{}, err
	}

	return res.ToCollectionDTO(), nil
}
