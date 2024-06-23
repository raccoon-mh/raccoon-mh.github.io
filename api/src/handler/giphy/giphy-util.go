package giphy

import (
	"api/src/handler/common"
	"api/src/handler/models"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

var (
	GIPHY_API_TOKEN string
	GIPHY_HOST      string
)

func init() {
	GIPHY_API_TOKEN = os.Getenv("GIPHY_API_TOKEN")
	GIPHY_HOST = "http://api.giphy.com"
}

func GiphyRandomRaccon(c echo.Context) *models.CommonResponse {
	endPoint := "/v1/gifs/random"

	req := &models.CommonRequest{
		QueryParams: map[string]string{
			"api_key":   GIPHY_API_TOKEN,
			"tag":       "raccoon",
			"rating":    "g",
			"random_id": "060860e33c864f5791b1714fac3ae028",
		},
	}

	commonResponse, _ := common.CommonCallerWithoutToken(http.MethodGet, GIPHY_HOST, endPoint, req)
	return commonResponse
}

func GiphsearchRaccon(c echo.Context) *models.CommonResponse {
	endPoint := "/v1/gifs/random"

	req := &models.CommonRequest{
		QueryParams: map[string]string{
			"api_key":   GIPHY_API_TOKEN,
			"tag":       "raccoon animals cute",
			"rating":    "g",
			"random_id": "060860e33c864f5791b1714fac3ae028",
		},
	}

	commonResponse, _ := common.CommonCallerWithoutToken(http.MethodGet, GIPHY_HOST, endPoint, req)
	return commonResponse
}
