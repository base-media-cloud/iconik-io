package iconik

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/assets/collections"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/assets/collections/contents"
	"github.com/rs/zerolog"
	"net/http"
	"strconv"
)

func (a *API) GetCollection(ctx context.Context, path, collectionID string) (collections.DTO, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, a.cfg.OperationTimeout)
	defer cancel()
	zerolog.Ctx(ctx).Info().Msg("getting collection from iconik")

	body, statusCode, err := a.req.Do(
		ctxTimeout,
		http.MethodGet,
		a.url+path+collectionID,
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
		return collections.DTO{}, domain.ErrInternalError

	case *statusCode == http.StatusTooManyRequests,
		*statusCode == http.StatusInternalServerError,
		*statusCode == http.StatusServiceUnavailable,
		*statusCode == http.StatusGatewayTimeout:
		f := func() error {
			body, statusCode, err = a.req.Do(
				ctxTimeout,
				http.MethodGet,
				a.url+path,
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

	case *statusCode != http.StatusOK:
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Int("status code", *statusCode).
			Msg("status code unexpected")
		return collections.DTO{}, domain.ErrInternalError
	}

	if err != nil {
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Msg("error getting collection")
		return collections.DTO{}, domain.ErrInternalError
	}

	var res collections.Collection
	if err = json.Unmarshal(body, &res); err != nil {
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Msg("error unmarshalling body")
		return collections.DTO{}, domain.ErrInternalError
	}

	zerolog.Ctx(ctxTimeout).Info().
		Msg("successfully got collection contents from iconik api")

	return res.ToDTO(), nil
}

func (a *API) GetCollContents(ctx context.Context, path, collectionID string, pageNo int) ([]contents.ObjectDTO, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, a.cfg.OperationTimeout)
	defer cancel()
	zerolog.Ctx(ctx).Info().Msg("getting collection contents from iconik")

	a.queryParams["page"] = strconv.Itoa(pageNo)

	endpoint := a.url + path + collectionID + "/contents/"

	body, statusCode, err := a.req.Do(
		ctxTimeout,
		http.MethodGet,
		endpoint,
		a.headers,
		a.queryParams,
		nil,
	)

	opDelay := a.cfg.OperationRetryDelay

	switch {
	case statusCode == nil:
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Msg("status code is nil")
		return nil, domain.ErrInternalError

	case *statusCode == http.StatusTooManyRequests,
		*statusCode == http.StatusInternalServerError,
		*statusCode == http.StatusServiceUnavailable,
		*statusCode == http.StatusGatewayTimeout:
		f := func() error {
			body, statusCode, err = a.req.Do(
				ctxTimeout,
				http.MethodGet,
				a.url+path,
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
		return nil, domain.ErrInternalError
	}

	if err != nil {
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Msg("error getting collection contents")
		return nil, domain.ErrInternalError
	}

	var res contents.Contents
	if err = json.Unmarshal(body, &res); err != nil {
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Msg("error unmarshalling body")
		return nil, domain.ErrInternalError
	}

	if res.Errors != nil {
		fmt.Println(res.Errors, endpoint, collectionID)
		// return NewWrappedErrs(res.Errors)
	}

	dtos := make([]contents.ObjectDTO, len(res.Objects))
	for i, o := range res.Objects {
		dtos[i] = o.ToDTO()
	}

	zerolog.Ctx(ctxTimeout).Info().
		Int("collection contents", len(dtos)).
		Msg("successfully got collection contents from iconik api")

	return dtos, nil
}
