package command

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

var codeToWeather = map[string]string{
	"T1H": "기온(℃)",
	"RN1": "1시간 강수량(mm)",
	"SKY": "하늘상태",
	"UUU": "동서바람(m/s)",
	"VVV": "남북바람(m/s)",
	"REH": "습도(%)",
	"PTY": "강수형태",
	"LGT": "낙뢰(kA)",
	"VEC": "풍향(deg)",
	"WSD": "풍속(m/s)",
}

var ptyMap = map[string]string{
	"0": "☀️ 맑음",
	"1": "🌧️ 비",
	"2": "🌧️ 🌨️ 비/눈",
	"3": "🌨️ 눈",
	"5": "🌧️ 빗방울",
	"6": "🌧️ 🌨️ 빗방울눈날림",
	"7": "🌨️ 눈날림",
}

func weatherCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	margs := make([]interface{}, 0, len(options))
	msgformat := ""

	if option, ok := optionMap["지역"]; ok {
		locationData, err := getLocationInfo(option.StringValue())
		if err != nil || len(locationData.Documents) == 0 {
			log.Println(err)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "지역 정보를 가져오는 동안 문제가 생겼구리.. 진짜 너구리를 찾아구리!",
				},
			})
			return
		}
		margs = append(margs, option.StringValue())
		margs = append(margs, locationData.Documents[0].AddressName)

		lon, _ := strconv.ParseFloat(locationData.Documents[0].X, 64)
		lat, _ := strconv.ParseFloat(locationData.Documents[0].Y, 64)
		X, Y := convertLonLatToGrid(lon, lat)
		// margs = append(margs, fmt.Sprintf("%d %d", X, Y))

		weatherData, err := getWeather(X, Y)
		if err != nil {
			log.Println(err)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "날씨 정보를 가져오는 동안 문제가 생겼구리.. 진짜 너구리를 찾아구리!",
				},
			})
			return
		}

		// margs = append(margs, fmt.Sprintf("%d", len(weatherData.Body.Items.Item)))
		var ptyTemp string
		weatherDataParsed := ""
		for _, item := range weatherData.Body.Items.Item {
			if item.Category == "PTY" {
				ptyTemp = strconv.FormatFloat(item.ObsrValue, 'f', -1, 64)
				weatherDataParsed += fmt.Sprintf("### %s\n", ptyMap[ptyTemp])
			} else if item.Category == "RN1" {
				if ptyTemp != "0" {
					weatherDataParsed += fmt.Sprintf("> %s : %g\n", codeToWeather[item.Category], item.ObsrValue)
				}
			} else if item.Category == "WSD" || item.Category == "REH" || item.Category == "T1H" {
				weatherDataParsed += fmt.Sprintf("> %s : %g\n", codeToWeather[item.Category], item.ObsrValue)
			} else {
				continue
			}
		}

		margs = append(margs, weatherDataParsed)

		msgformat += "### (%s)초단기예보-%s\n%s"
		msgformat += "\n`날씨구리 사용법 /날씨 <지역>`"
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// Ignore type for now, they will be discussed in "responses"
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				msgformat,
				margs...,
			),
		},
	})
}

// LocationQueryDataResponsse START
type Document struct {
	AddressName       string `json:"address_name"`
	CategoryGroupCode string `json:"category_group_code"`
	CategoryGroupName string `json:"category_group_name"`
	CategoryName      string `json:"category_name"`
	Distance          string `json:"distance"`
	ID                string `json:"id"`
	Phone             string `json:"phone"`
	PlaceName         string `json:"place_name"`
	PlaceURL          string `json:"place_url"`
	RoadAddressName   string `json:"road_address_name"`
	X                 string `json:"x"`
	Y                 string `json:"y"`
}

type Meta struct {
	IsEnd         bool `json:"is_end"`
	PageableCount int  `json:"pageable_count"`
	SameName      struct {
		Keyword        string   `json:"keyword"`
		Region         []string `json:"region"`
		SelectedRegion string   `json:"selected_region"`
	} `json:"same_name"`
	TotalCount int `json:"total_count"`
}

type LocationQueryDataResponsse struct {
	Documents []Document `json:"documents"`
	Meta      Meta       `json:"meta"`
}

// LocationQueryDataResponsse END

func getLocationInfo(query string) (LocationQueryDataResponsse, error) {
	baseURL := "https://dapi.kakao.com/v2/local/search/keyword.json"
	contentType := "application/json"
	authorization := KAKAO_TOKEN

	params := url.Values{}
	params.Add("query", query)
	params.Add("size", "1")

	finalURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequest("GET", finalURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return LocationQueryDataResponsse{}, err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", authorization)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return LocationQueryDataResponsse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return LocationQueryDataResponsse{}, err
	}

	fmt.Println("kogpt Response Status:", resp.Status)
	if resp.StatusCode != 200 {
		return LocationQueryDataResponsse{}, fmt.Errorf("통신문제발생")
	}

	fmt.Println("kogpt Response Body:", string(body))

	var data LocationQueryDataResponsse

	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		fmt.Println(err)
		return LocationQueryDataResponsse{}, err
	}

	fmt.Println("@@@@@@@ data", data)

	return data, err
}

// WeatherResponse START
type WeatherResponse struct {
	XMLName xml.Name `xml:"response"`
	Header  Header   `xml:"header"`
	Body    Body     `xml:"body"`
}

type Header struct {
	ResultCode string `xml:"resultCode"`
	ResultMsg  string `xml:"resultMsg"`
}

type Body struct {
	DataType   string `xml:"dataType"`
	Items      Items  `xml:"items"`
	NumOfRows  int    `xml:"numOfRows"`
	PageNo     int    `xml:"pageNo"`
	TotalCount int    `xml:"totalCount"`
}

type Items struct {
	Item []Item `xml:"item"`
}

type Item struct {
	BaseDate  string  `xml:"baseDate"`
	BaseTime  string  `xml:"baseTime"`
	Category  string  `xml:"category"`
	Nx        int     `xml:"nx"`
	Ny        int     `xml:"ny"`
	ObsrValue float64 `xml:"obsrValue"`
}

// WeatherResponse END

func getWeather(nx int, ny int) (WeatherResponse, error) {
	baseURL := "http://apis.data.go.kr/1360000/VilageFcstInfoService_2.0/getUltraSrtNcst"
	authorization := KISANG_TOKEN

	currentTime := time.Now()
	base_date := currentTime.Format("20060102")
	base_time := currentTime.Format("1504")

	params := url.Values{}
	params.Add("serviceKey", authorization)
	params.Add("base_date", base_date)
	params.Add("base_time", base_time)
	params.Add("nx", strconv.Itoa(nx))
	params.Add("ny", strconv.Itoa(ny))

	finalURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	log.Println("@@@@@@@@@@@@@@@ finalURL:", finalURL)

	req, err := http.NewRequest("GET", finalURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return WeatherResponse{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return WeatherResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return WeatherResponse{}, err
	}

	fmt.Println("kogpt Response Status:", resp.Status)
	if resp.StatusCode != 200 {
		return WeatherResponse{}, fmt.Errorf("통신문제발생")
	}

	fmt.Println("kogpt Response Body:", string(body))

	var data WeatherResponse
	err = xml.Unmarshal([]byte(body), &data)
	if err != nil {
		fmt.Println(err)
		return WeatherResponse{}, err
	}

	fmt.Println("@@@@@@@ data", data)

	return data, err
}

func convertLonLatToGrid(lon, lat float64) (int, int) {
	const (
		Re    = 6371.00877   // 지도반경
		Grid  = 5.0          // 격자간격 (km)
		Slat1 = 30.0         // 표준위도 1
		Slat2 = 60.0         // 표준위도 2
		Olon  = 126.0        // 기준점 경도
		Olat  = 38.0         // 기준점 위도
		Xo    = 210.0 / Grid // 기준점 X좌표
		Yo    = 675.0 / Grid // 기준점 Y좌표
	)
	const PI = math.Pi
	const DEGRAD = PI / 180.0

	re := Re / Grid
	slat1 := Slat1 * DEGRAD
	slat2 := Slat2 * DEGRAD
	olon := Olon * DEGRAD
	olat := Olat * DEGRAD

	sn := math.Tan(PI*0.25+slat2*0.5) / math.Tan(PI*0.25+slat1*0.5)
	sn = math.Log(math.Cos(slat1)/math.Cos(slat2)) / math.Log(sn)
	sf := math.Tan(PI*0.25 + slat1*0.5)
	sf = math.Pow(sf, sn) * math.Cos(slat1) / sn
	ro := math.Tan(PI*0.25 + olat*0.5)
	ro = re * sf / math.Pow(ro, sn)

	ra := math.Tan(PI*0.25 + (lat * DEGRAD * 0.5))
	ra = re * sf / math.Pow(ra, sn)
	theta := lon*DEGRAD - olon
	if theta > PI {
		theta -= 2.0 * PI
	}
	if theta < -PI {
		theta += 2.0 * PI
	}
	theta *= sn
	x := ra*math.Sin(theta) + Xo
	y := ro - ra*math.Cos(theta) + Yo

	return int(math.Ceil(x)), int(math.Ceil(y))
}
