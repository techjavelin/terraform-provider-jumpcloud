package apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type (
	Client struct {
		ApiKey          string
		httpClient      *http.Client
		ProviderVersion string
		Context         context.Context
	}

	reusableReader struct {
		io.Reader
		readBuf *bytes.Buffer
		backBuf *bytes.Buffer
	}
)

const (
	JUMPCLOUD_API_BASE_URL = "https://console.jumpcloud.com/api"
	CONTENT_TYPE           = "application/json"
	SUBSYSTEM_NAME         = "apiclient.Client"
)

func (c *Client) ReusableReader(r io.Reader) io.Reader {
	readBuf := bytes.Buffer{}
	readBuf.ReadFrom(r)
	backBuf := bytes.Buffer{}

	return reusableReader{
		io.TeeReader(&readBuf, &backBuf),
		&readBuf,
		&backBuf,
	}
}

func (r reusableReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	if err == io.EOF {
		r.Reset()
	}
	return n, err
}

func (r reusableReader) Reset() {
	io.Copy(r.readBuf, r.backBuf)
}

func (c *Client) ReadBody(r io.Reader) string {
	b, _ := io.ReadAll(r)
	return string(b)
}

func New(ctx context.Context, apikey string, providerVersion string) Client {
	client := &http.Client{}

	tflog.Info(ctx, fmt.Sprintf("Initializing %s Logging Subsystem", SUBSYSTEM_NAME))

	return Client{
		ApiKey:          apikey,
		httpClient:      client,
		ProviderVersion: providerVersion,
		Context:         tflog.NewSubsystem(ctx, SUBSYSTEM_NAME, tflog.WithLevelFromEnv("TF_LOG_PROVIDER_JUMPCLOUD_CLIENT")),
	}
}

func (c *Client) getPayloadAsString(payload interface{}) string {
	var buf *bytes.Buffer = &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(payload)
	if err != nil {
		tflog.SubsystemError(c.Context, SUBSYSTEM_NAME, spew.Sprintf("Error while encoding payload as string: %s", err), map[string]interface{}{
			"method": "getPayloadAsString",
		})
		return "<nil>"
	}

	return buf.String()
}

func (c *Client) prepareRequest(
	method string,
	apiVersion string, endpoint string,
	postBody interface{},
	headerParams map[string]string,
	queryParams url.Values) (request *http.Request, err error) {
	request_url := fmt.Sprintf("%s/%s/%s", JUMPCLOUD_API_BASE_URL, apiVersion, endpoint)

	var body *bytes.Buffer

	if postBody != nil {
		body = &bytes.Buffer{}
		err = json.NewEncoder(body).Encode(postBody)
	}

	if err != nil {
		return nil, err
	}

	url, err := url.Parse(request_url)
	if err != nil {
		return nil, err
	}

	query := url.Query()
	for k, v := range queryParams {
		for _, iv := range v {
			query.Add(k, iv)
		}
	}

	url.RawQuery = query.Encode()

	if body != nil {
		request, err = http.NewRequest(method, url.String(), body)
	} else {
		request, err = http.NewRequest(method, url.String(), nil)
	}

	if err != nil {
		return nil, err
	}

	if len(headerParams) > 0 {
		headers := http.Header{}
		for h, v := range headerParams {
			headers.Set(h, v)
		}
		request.Header = headers
	}

	request.Header.Add("x-api-key", c.ApiKey)
	request.Header.Add("User-Agent", "registry.terraform.io/techjavelin/jumpcloud@V"+c.ProviderVersion)
	request.Header.Add("Accept", CONTENT_TYPE)
	request.Header.Add("Content-Type", CONTENT_TYPE)

	return request, nil
}
