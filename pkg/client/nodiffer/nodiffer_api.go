package nodiffer

import (
	"bytes"
	"context"
	"differ-template-engine/pkg/config"
	"encoding/json"
	"errors"
	"github.com/hashicorp/go-retryablehttp"
	"io"
	"net"
	"net/http"
	"time"
)

type API interface {
	HasDiff(ctx context.Context, request HasDiffRequest) (*HasDiffResponse, error)
}

type nodifferApi struct {
	cfg    config.APIConfig
	client http.Client
}

func NewNodifferAPI(cfg config.APIConfig) API {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 60 * time.Second,
		}).DialContext,
		MaxIdleConns:          10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	api := &nodifferApi{
		cfg: cfg,
	}

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = cfg.Retry
	retryClient.Logger = nil
	api.client.Transport = transport
	api.client.Timeout = cfg.Timeout
	api.client = *retryClient.StandardClient()

	return api
}

func (n *nodifferApi) HasDiff(ctx context.Context, request HasDiffRequest) (*HasDiffResponse, error) {
	url := n.cfg.Host + "/differ/execute"

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(body)

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, reader)
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("nodiffer request is not ok")
	}

	if bodyBytes, err := io.ReadAll(resp.Body); err == nil {
		var resp *HasDiffResponse
		if err := json.Unmarshal(bodyBytes, &resp); err == nil {
			return resp, nil
		} else {
			return nil, err
		}
	}

	return nil, err
}
