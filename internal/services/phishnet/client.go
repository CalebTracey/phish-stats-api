package phishnet

import "net/url"

type ApiClient struct {
	BaseUrl   *url.URL
	ApiKeyUri string
}
