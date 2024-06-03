package actions

import (
	"log"
	"strings"

	"api/actions/giphy"
	"api/util"

	"github.com/gobuffalo/buffalo"
)

func PostRouteController(c buffalo.Context) error {
	log.Println("#### Post RouteController ")
	commonRequest := &util.CommonRequest{}
	c.Bind(commonRequest)
	targetController := strings.ToLower(c.Param("targetController"))
	log.Printf("== targetController\t:[ %s ]\n", targetController)
	log.Printf("== commonRequest\t:\n%+v\n\n", commonRequest)

	commonResponse := &util.CommonResponse{}
	switch targetController {
	default:
		commonResponse = util.CommonResponseStatusNotFound("NO MATCH targetController")
		return c.Render(commonResponse.Status.StatusCode, r.JSON(commonResponse))
	}
}

func GetRouteController(c buffalo.Context) error {
	log.Println("#### Get RouteController ")
	targetController := strings.ToLower(c.Param("targetController"))
	log.Printf("== targetController\t:[ %s ]\n", targetController)

	commonResponse := &util.CommonResponse{}
	switch targetController {
	case "getrandomraccoon":
		commonResponse = giphy.RandomRaccon(c)
	default:
		commonResponse = util.CommonResponseStatusNotFound("NO MATCH targetController")
		return c.Render(commonResponse.Status.StatusCode, r.JSON(commonResponse))
	}

	return c.Render(commonResponse.Status.StatusCode, r.JSON(commonResponse))
}
