package iconikio

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func (i *Iconik) getResponseBody(method, uri string, params io.Reader) (*http.Response, []byte, error) {
	log.Println(uri)

	client := &http.Client{}
	req, err := http.NewRequest(method, uri, params)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Add("App-ID", i.IconikClient.Config.AppID)
	req.Header.Add("Auth-Token", i.IconikClient.Config.AuthToken)
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

func removeNullJSON(m map[string]interface{}) map[string]interface{} {
	for k, v := range m {
		if v == nil {
			delete(m, k)
			continue
		}
		switch val := v.(type) {
		case map[string]interface{}:
			m[k] = removeNullJSON(val)
		case []interface{}:
			for i := 0; i < len(val); i++ {
				if _, ok := val[i].(map[string]interface{}); ok {
					val[i] = removeNullJSON(val[i].(map[string]interface{}))
				}
			}
		}
	}
	return m
}

func (i *Iconik) joinURL(endpoint, path string, index int) (*url.URL, error) {

	var paths []string

	path1 := i.IconikClient.Config.APIConfig.Endpoints[endpoint].([]interface{})[index].(map[string]interface{})["path"].([]string)
	path2 := i.IconikClient.Config.APIConfig.Endpoints[endpoint].([]interface{})[index].(map[string]interface{})["path2"].([]string)

	paths = append(paths, path1...)
	paths = append(paths, path)
	paths = append(paths, path2...)

	result, err := url.JoinPath(i.IconikClient.Config.APIConfig.Host, paths...)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(result)
	if err != nil {
		return nil, fmt.Errorf("invalid url")
	}

	u.Scheme = i.IconikClient.Config.APIConfig.Scheme

	return u, nil
}
