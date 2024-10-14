package httpClient

import (
	"bytes"
	"io"
	"net/url"
	"strings"
	"time"

	http "github.com/bogdanfinn/fhttp"
	"github.com/bogdanfinn/fhttp/cookiejar"
	tlsclient "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

type HttpClient struct {
	cookieJar http.CookieJar
	Client    tlsclient.HttpClient
	Proxy     string
}

type Response struct {
	Response *http.Response

	Body       string
	StatusCode int

	Cookies []*http.Cookie
	Headers map[string]string

	Content []byte
}

func NewHttpClient(proxy string, timeout time.Duration, followRedirects bool) (*HttpClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	options := []tlsclient.HttpClientOption{
		tlsclient.WithTimeoutSeconds(int(timeout.Seconds())),
		tlsclient.WithClientProfile(profiles.DefaultClientProfile),
		tlsclient.WithCookieJar(jar),
	}

	if !followRedirects {
		options = append(options, tlsclient.WithNotFollowRedirects())
	}

	if proxy != "" {
		options = append(options, tlsclient.WithProxyUrl(proxy))
	}

	client, err := tlsclient.NewHttpClient(tlsclient.NewNoopLogger(), options...)
	if err != nil {
		return nil, err
	}

	return &HttpClient{Client: client, cookieJar: jar, Proxy: proxy}, nil
}

func (h *HttpClient) Get(url string, headers map[string]string) (*Response, error) {
	return h.Request("GET", url, nil, "", headers)
}

func (h *HttpClient) Post(url string, body []byte, contentType string, headers map[string]string) (*Response, error) {
	return h.Request("POST", url, body, contentType, headers)
}

func (h *HttpClient) Cookies(url *url.URL) []*http.Cookie {
	return h.cookieJar.Cookies(url)
}

func (h *HttpClient) Request(method, url string, body []byte, contentType string, headers map[string]string) (*Response, error) {
	request, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		request.Header.Add(key, value)
	}

	if strings.TrimSpace(contentType) != "" {
		request.Header.Set("Content-Type", contentType)
	}

	response, err := h.Client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	responseHeaders := make(map[string]string)
	for key, value := range response.Header {
		if len(value) == 0 {
			continue
		}

		responseHeaders[key] = value[0]
	}

	rBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		Response: response,

		Body:       string(rBody),
		StatusCode: response.StatusCode,

		Cookies: response.Cookies(),
		Headers: responseHeaders,

		Content: rBody,
	}, nil
}
