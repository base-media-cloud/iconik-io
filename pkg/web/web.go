package web

import (
	"io"
	"net/http"

	"github.com/base-media-cloud/pd-iconik-io-rd/app/services/config"
	"go.uber.org/zap"
)

func GetResponseBody(method, uri string, params io.Reader, cfg *config.Conf, log *zap.SugaredLogger) (*http.Response, []byte, error) {
	log.Infow(uri)

	client := &http.Client{}
	req, err := http.NewRequest(method, uri, params)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Add("App-ID", cfg.AppID)
	req.Header.Add("Auth-Token", cfg.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	return res, resBody, nil
}
