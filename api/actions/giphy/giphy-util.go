package giphy

import (
	"api/util"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
)

func RandomRaccon(c buffalo.Context) *util.CommonResponse {
	giphyHost := "http://api.giphy.com"
	endPoint := "/v1/gifs/random"

	req := &util.CommonRequest{
		QueryParams: map[string]string{
			"api_key": envy.Get("GIPHY_API_TOKEN", ""),
			"tag":     "raccoon",
			"rating":  "g",
		},
	}

	commonResponse, _ := util.CommonCallerWithoutToken(http.MethodGet, giphyHost, endPoint, req)
	return commonResponse
}
