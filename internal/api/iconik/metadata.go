package iconik

import (
	"context"
	"encoding/json"
	"github.com/avast/retry-go"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/metadata/views"
	"github.com/rs/zerolog"
	"net/http"
)

func (a *API) GetMetadataViews(ctx context.Context, path, viewID string) ([]views.ViewFieldDTO, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, a.cfg.OperationTimeout)
	defer cancel()
	zerolog.Ctx(ctx).Info().Msg("getting metadata views from iconik")

	body, statusCode, err := a.req.Do(
		ctxTimeout,
		http.MethodGet,
		a.url+path,
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
				Msg("retrying to get metadata views from iconik")
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
			Msg("error getting metadata views")
		return nil, domain.ErrInternalError
	}

	var res views.Views
	if err = json.Unmarshal(body, &res); err != nil {
		zerolog.Ctx(ctxTimeout).Error().
			Err(err).
			Msg("error unmarshalling body")
		return nil, domain.ErrInternalError
	}

	dtos := make([]views.ViewFieldDTO, len(res.ViewFields))
	for i, o := range res.ViewFields {
		dtos[i] = o.ToDTO()
	}

	zerolog.Ctx(ctxTimeout).Info().
		Int("metadata views", len(dtos)).
		Msg("successfully got metadata views from iconik api")

	return dtos, nil
}
