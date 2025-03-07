package iconik

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/metadata"
	"github.com/rs/zerolog"
	"net/http"
)

// GetMetadataView makes a request to the GET iconik metadata view endpoint.
func (a *API) GetMetadataView(ctx context.Context, path, viewID string) (metadata.DTO, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, a.cfg.OperationTimeout)
	defer cancel()

	body, statusCode, err := a.req.Do(
		ctxTimeout,
		http.MethodGet,
		fmt.Sprintf("%v%v%v/", a.url, path, viewID),
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
		return metadata.DTO{}, err
	case *statusCode == http.StatusTooManyRequests,
		*statusCode == http.StatusInternalServerError,
		*statusCode == http.StatusServiceUnavailable,
		*statusCode == http.StatusGatewayTimeout:
		f := func() error {
			body, statusCode, err = a.req.Do(
				ctxTimeout,
				http.MethodGet,
				fmt.Sprintf("%v%v%v/", a.url, path, viewID),
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
				Msg("retrying to get metadata view from iconik")
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
			Msg("forbidden when getting metadata view")
		return metadata.DTO{}, domain.ErrForbidden
	case *statusCode == http.StatusUnauthorized:
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Int("status code", *statusCode).
			RawJSON("response", body).
			Msg("unauthorized when getting metadata view")
		return metadata.DTO{},
			fmt.Errorf("you do not have the correct permissions to get the metadata view %s", viewID)
	case *statusCode != http.StatusOK:
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			RawJSON("response", body).
			Int("status code", *statusCode).
			Msg("status code unexpected")
		return metadata.DTO{}, domain.ErrInternalError
	}

	if err != nil {
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Msg("error getting metadata view")
		return metadata.DTO{}, err
	}

	var res metadata.Metadata
	if err = json.Unmarshal(body, &res); err != nil {
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Msg("error unmarshalling body")
		return metadata.DTO{}, err
	}

	return res.ToDTO(), nil
}

// UpdateMetadataInAsset makes a request to the PUT iconik metadata endpoint.
func (a *API) UpdateMetadataInAsset(ctx context.Context, path, viewID, assetID string, payload []byte) (metadata.DTO, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, a.cfg.OperationTimeout)
	defer cancel()

	body, statusCode, err := a.req.Do(
		ctxTimeout,
		http.MethodPut,
		fmt.Sprintf("%v%v%v/views/%v/", a.url, path, assetID, viewID),
		a.headers,
		nil,
		payload,
	)

	opDelay := a.cfg.OperationRetryDelay

	switch {
	case statusCode == nil:
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Msg("status code is nil")
		return metadata.DTO{}, err
	case *statusCode == http.StatusTooManyRequests,
		*statusCode == http.StatusInternalServerError,
		*statusCode == http.StatusServiceUnavailable,
		*statusCode == http.StatusGatewayTimeout:
		f := func() error {
			body, statusCode, err = a.req.Do(
				ctxTimeout,
				http.MethodPut,
				fmt.Sprintf("%v%v%v/views/%v/", a.url, path, assetID, viewID),
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
				Msg("retrying to update metadata in iconik")
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
			Msg("forbidden when updating metadata")
		return metadata.DTO{}, domain.ErrForbidden
	case *statusCode == http.StatusUnauthorized:
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Int("status code", *statusCode).
			RawJSON("response", body).
			Msg("unauthorized when updating metadata")
		return metadata.DTO{},
			fmt.Errorf("you do not have the correct permissions to update the metadata for asset %s", assetID)
	case *statusCode != http.StatusOK:
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			RawJSON("response", body).
			Int("status code", *statusCode).
			Msg("status code unexpected")
		return metadata.DTO{}, domain.ErrInternalError
	}

	if err != nil {
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Msg("error updating metadata")
		return metadata.DTO{}, err
	}

	var res metadata.Metadata
	if err = json.Unmarshal(body, &res); err != nil {
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Msg("error unmarshalling body")
		return metadata.DTO{}, err
	}

	return res.ToDTO(), nil
}
