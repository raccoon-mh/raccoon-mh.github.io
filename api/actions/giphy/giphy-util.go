package giphy

import (
	"api/util"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
)

var (
	GIPHY_API_TOKEN string
	GIPHY_HOST      string
)

func init() {
	GIPHY_API_TOKEN = envy.Get("GIPHY_API_TOKEN", "")
	GIPHY_HOST = "http://api.giphy.com"
}

func GiphyRandomRaccon(c buffalo.Context) *util.CommonResponse {
	endPoint := "/v1/gifs/random"

	req := &util.CommonRequest{
		QueryParams: map[string]string{
			"api_key":   GIPHY_API_TOKEN,
			"tag":       "raccoon",
			"rating":    "g",
			"random_id": "060860e33c864f5791b1714fac3ae028",
		},
	}

	commonResponse, _ := util.CommonCallerWithoutToken(http.MethodGet, GIPHY_HOST, endPoint, req)
	return commonResponse
}

func GiphyearchRaccon(c buffalo.Context) *util.CommonResponse {
	endPoint := "/v1/gifs/random"

	req := &util.CommonRequest{
		QueryParams: map[string]string{
			"api_key":   GIPHY_API_TOKEN,
			"tag":       "raccoon animals cute",
			"rating":    "g",
			"random_id": "060860e33c864f5791b1714fac3ae028",
		},
	}

	commonResponse, _ := util.CommonCallerWithoutToken(http.MethodGet, GIPHY_HOST, endPoint, req)
	return commonResponse
}
