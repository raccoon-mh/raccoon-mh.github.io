package common

import (
	"api/src/handler/models"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func CommonCaller(callMethod string, targetUrl string, endPoint string, commonRequest *models.CommonRequest, auth string) (*models.CommonResponse, error) {
	pathParamsUrl := mappingUrlPathParams(endPoint, commonRequest)
	queryParamsUrl := mappingQueryParams(pathParamsUrl, commonRequest)
	requestUrl := targetUrl + queryParamsUrl
	commonResponse, err := CommonHttpToCommonResponse(requestUrl, commonRequest.Request, callMethod, auth)
	return commonResponse, err
}

func CommonCallerWithoutToken(callMethod string, targetUrl string, endPoint string, commonRequest *models.CommonRequest) (*models.CommonResponse, error) {
	pathParamsUrl := mappingUrlPathParams(endPoint, commonRequest)
	queryParamsUrl := mappingQueryParams(pathParamsUrl, commonRequest)
	requestUrl := targetUrl + queryParamsUrl
	commonResponse, err := CommonHttpToCommonResponse(requestUrl, commonRequest.Request, callMethod, "")
	return commonResponse, err
}

func mappingUrlPathParams(endPoint string, commonRequest *models.CommonRequest) string {
	u := endPoint
	for k, r := range commonRequest.PathParams {
		u = strings.Replace(u, "{"+k+"}", r, -1)
	}
	return u
}

func mappingQueryParams(targeturl string, commonRequest *models.CommonRequest) string {
	u, err := url.Parse(targeturl)
	if err != nil {
		return ""
	}
	q := u.Query()
	for k, v := range commonRequest.QueryParams {
		q.Set(string(k), v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func CommonHttpToCommonResponse(url string, s interface{}, httpMethod string, auth string) (*models.CommonResponse, error) {
	log.Println("CommonHttp - METHOD:" + httpMethod + " => url:" + url)
	log.Println("isauth:", auth)

	jsonData, err := json.Marshal(s)
	if err != nil {
		log.Println("commonPostERR : json.Marshal : ", err.Error())
		return nil, err
	}

	req, err := http.NewRequest(httpMethod, url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error CommonHttp creating request:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if auth != "" {
		req.Header.Add("Authorization", auth)
	}

	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Println("Error CommonHttp creating httputil.DumpRequest:", err)
	}
	log.Println("\n", string(requestDump))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error CommonHttp request:", err)
		return models.CommonResponseStatusInternalServerError(err), err
	}
	defer resp.Body.Close()

	respBody, ioerr := io.ReadAll(resp.Body)
	if ioerr != nil {
		log.Println("Error CommonHttp reading response:", ioerr)
	}

	commonResponse := &models.CommonResponse{}
	commonResponse.Status.Message = resp.Status
	commonResponse.Status.StatusCode = resp.StatusCode
	if len(respBody) > 0 {
		jsonerr := json.Unmarshal(respBody, &commonResponse.ResponseData)
		if jsonerr != nil {
			log.Println("Error CommonHttp Unmarshal response:", jsonerr.Error())
			return commonResponse, jsonerr
		}
	}
	return commonResponse, nil
}
