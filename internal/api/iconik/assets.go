package iconik

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/assets/assets"
	"net/http"

	"github.com/avast/retry-go"
	"github.com/rs/zerolog"
)

// GetAsset makes a request to the GET iconik asset endpoint.
func (a *API) GetAsset(ctx context.Context, path, assetID string) (assets.DTO, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, a.cfg.OperationTimeout)
	defer cancel()
	zerolog.Ctx(ctxTimeout).Info().Msg("getting asset " + assetID + " from iconik")

	body, statusCode, err := a.req.Do(
		ctxTimeout,
		http.MethodGet,
		a.url+path+"/"+assetID+"/",
		a.headers,
		nil,
		nil,
	)

	opDelay := a.cfg.OperationRetryDelay

	switch {
	case errors.Is(err, domain.ErrTransformingHeaderValue) || errors.Is(err, domain.ErrTransformingHeaderKey):
		return assets.DTO{}, err
	case statusCode == nil:
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Msg("status code is nil")
		return assets.DTO{}, err
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
				Msg("retrying to get asset from iconik")
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
		return assets.DTO{}, err
	}

	if err != nil {
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Msg("error getting asset")
		return assets.DTO{}, err
	}

	var res assets.Asset
	if err = json.Unmarshal(body, &res); err != nil {
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Msg("error unmarshalling body")
		return assets.DTO{}, err
	}

	return res.ToDTO(), nil
}

// PatchAsset makes a request to the PATCH iconik asset endpoint.
func (a *API) PatchAsset(ctx context.Context, path, assetID string, payload []byte) (assets.DTO, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, a.cfg.OperationTimeout)
	defer cancel()
	zerolog.Ctx(ctxTimeout).Info().Msg("updating asset " + assetID + " in iconik")

	body, statusCode, err := a.req.Do(
		ctxTimeout,
		http.MethodPatch,
		a.url+path+"/"+assetID+"/",
		a.headers,
		nil,
		payload,
	)

	opDelay := a.cfg.OperationRetryDelay

	switch {
	case errors.Is(err, domain.ErrTransformingHeaderValue) || errors.Is(err, domain.ErrTransformingHeaderKey):
		return assets.DTO{}, err
	case statusCode == nil:
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Msg("status code is nil")
		return assets.DTO{}, err
	case *statusCode == http.StatusTooManyRequests,
		*statusCode == http.StatusInternalServerError,
		*statusCode == http.StatusServiceUnavailable,
		*statusCode == http.StatusGatewayTimeout:
		f := func() error {
			body, statusCode, err = a.req.Do(
				ctxTimeout,
				http.MethodPatch,
				a.url+path+"/"+assetID+"/",
				a.headers,
				nil,
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
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Int("status code", *statusCode).
			Msg("status code unexpected")
		return assets.DTO{}, err
	}

	if err != nil {
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Msg("error updating asset")
		return assets.DTO{}, err
	}

	var res assets.Asset
	if err = json.Unmarshal(body, &res); err != nil {
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Msg("error unmarshalling body")
		return assets.DTO{}, err
	}

	return res.ToDTO(), nil
}
