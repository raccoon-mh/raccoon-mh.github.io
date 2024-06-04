package pexels

import (
	"api/util"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
)

var (
	PEXELS_API_TOKEN string
	PEXELS_HOST      string
)

func init() {
	PEXELS_API_TOKEN = envy.Get("PEXELS_API_TOKEN", "")
	PEXELS_HOST = "https://api.pexels.com"
}

func generateRandomIntInRange(min, max int) int {
	randomNum := rand.Intn(max - min + 1)
	return min + randomNum
}

func PexelsRandomRaccon(c buffalo.Context) *util.CommonResponse {
	searchEndPoint := "/v1/search"
	reqFirst := &util.CommonRequest{
		QueryParams: map[string]string{
			"query":    "raccoon",
			"page":     "1",
			"per_page": "1",
		},
	}
	commonResponseFirst, _ := util.CommonCaller(http.MethodGet, PEXELS_HOST, searchEndPoint, reqFirst, PEXELS_API_TOKEN)

	totalResults := commonResponseFirst.ResponseData.(map[string]interface{})["total_results"].(float64)
	randomNum := rand.Intn(int(totalResults) + 1)

	reqSecond := &util.CommonRequest{
		QueryParams: map[string]string{
			"query":    "raccoon",
			"page":     strconv.Itoa(int(randomNum)),
			"per_page": "1",
		},
	}
	commonResponseSecond, _ := util.CommonCaller(http.MethodGet, PEXELS_HOST, searchEndPoint, reqSecond, PEXELS_API_TOKEN)

	return commonResponseSecond
}
