package pexels

import (
	"api/src/handler/common"
	"api/src/handler/models"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
)

var (
	PEXELS_API_TOKEN string
	PEXELS_HOST      string
)

func init() {
	PEXELS_API_TOKEN = os.Getenv("PEXELS_API_TOKEN")
	PEXELS_HOST = "https://api.pexels.com"
}

func generateRandomIntInRange(min, max int) int {
	randomNum := rand.Intn(max - min + 1)
	return min + randomNum
}

func PexelsRandomRaccon(c echo.Context) *models.CommonResponse {
	searchEndPoint := "/v1/search"
	reqFirst := &models.CommonRequest{
		QueryParams: map[string]string{
			"query":    "raccoon",
			"page":     "1",
			"per_page": "1",
		},
	}
	commonResponseFirst, _ := common.CommonCaller(http.MethodGet, PEXELS_HOST, searchEndPoint, reqFirst, PEXELS_API_TOKEN)

	totalResults := commonResponseFirst.ResponseData.(map[string]interface{})["total_results"].(float64)
	randomNum := rand.Intn(int(totalResults) + 1)

	reqSecond := &models.CommonRequest{
		QueryParams: map[string]string{
			"query":    "raccoon",
			"page":     strconv.Itoa(int(randomNum)),
			"per_page": "1",
		},
	}
	commonResponseSecond, _ := common.CommonCaller(http.MethodGet, PEXELS_HOST, searchEndPoint, reqSecond, PEXELS_API_TOKEN)

	return commonResponseSecond
}
