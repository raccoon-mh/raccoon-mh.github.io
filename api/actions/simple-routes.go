package actions

import (
	"log"
	"strings"

	"github.com/gobuffalo/buffalo"
)

func PostRouteController(c buffalo.Context) error {
	log.Println("#### Post RouteController ")
	commonRequest := &CommonRequest{}
	c.Bind(commonRequest)
	targetController := strings.ToLower(c.Param("targetController"))
	log.Printf("== targetController\t:[ %s ]\n", targetController)
	log.Printf("== commonRequest\t:\n%+v\n\n", commonRequest)

	commonResponse := &CommonResponse{}
	switch targetController {
	default:
		commonResponse = CommonResponseStatusNotFound("NO MATCH targetController")
		return c.Render(commonResponse.Status.StatusCode, r.JSON(commonResponse))
	}
}

// Get으로 전송되는 data 처리를 위하여
func GetRouteController(c buffalo.Context) error {
	log.Println("#### Get RouteController ")
	targetController := strings.ToLower(c.Param("targetController"))
	log.Printf("== targetController\t:[ %s ]\n", targetController)

	commonResponse := &CommonResponse{}
	switch targetController {
	default:
		commonResponse = CommonResponseStatusNotFound("NO MATCH targetController")
		return c.Render(commonResponse.Status.StatusCode, r.JSON(commonResponse))
	}
}
