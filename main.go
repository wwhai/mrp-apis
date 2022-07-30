package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var httpC http.Client = http.Client{}

func HttpGet(client http.Client, url string) (string, error) {
	var err error
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

type WeatherInfo struct {
	Data struct {
		Forecast []struct {
			Date      string `json:"date"`
			High      string `json:"high"`
			Fengli    string `json:"fengli"`
			Low       string `json:"low"`
			Fengxiang string `json:"fengxiang"`
			Type      string `json:"type"`
		} `json:"forecast"`
	} `json:"data"`
	Status int `json:"status"`
}

type mobileInfo struct {
	Ret    string   `json:"ret"`
	Mobile string   `json:"mobile"`
	Data   []string `json:"data"`
}

//
func BuildWeatherInfoString(w WeatherInfo) string {
	s := ""
	for _, v := range w.Data.Forecast {
		s += v.Date + ", "
		s += v.Type + ", "
		s += v.Low[strings.Index(v.Low, " "):] + " ~ "
		s += v.High[strings.Index(v.High, " "):] + ", "
		s += v.Fengxiang + "\n"
	}
	return s
}

//{"ret":"ok","mobile":"13209760000","data":["青海","海西","联通","0977",""]}
func BuildMobileInfoString(m mobileInfo) string {
	s := "省  份: " + m.Data[0] + "\n"
	s += "市  区: " + m.Data[1] + "\n"
	s += "运营商: " + m.Data[2] + "\n"
	s += "区  号: " + m.Data[3]
	return s
}

func getWeather(city string) (weatherInfo WeatherInfo, err error) {
	// http://wthrcdn.etouch.cn/weather_mini?city=广州
	result, err1 := HttpGet(httpC, "http://wthrcdn.etouch.cn/weather_mini?city="+city)
	if err1 != nil {
		err = err1
		return
	}
	if err = json.Unmarshal([]byte(result), &weatherInfo); err != nil {
		return
	}
	return
}

func getMobile(mobile string) (r mobileInfo, err error) {
	// find({"ret":"ok","mobile":"13209760000","data":["青海","海西","联通","0977",""]})
	// http://wthrcdn.etouch.cn/weather_mini?city=广州
	result, err1 := HttpGet(httpC,
		"https://api.ip138.com/mobile/?datatype=jsonp&callback=find&token=126aa3a2be7b606dd727fd38d7d71f07&mobile="+mobile)
	if err1 != nil {
		err = err1
		return
	}
	if err = json.Unmarshal([]byte(result[5:len(result)-1]), &r); err != nil {
		return
	}
	return
}
func Weather(w http.ResponseWriter, r *http.Request) {
	log.Println("Request From", r.Host)
	weatherInfo, err := getWeather(r.URL.Query().Get("city"))
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	if weatherInfo.Status == 1000 {
		w.Write([]byte(BuildWeatherInfoString(weatherInfo)))
		return
	}
	w.Write([]byte("GET ERROR"))
}

// https://api.ip138.com/mobile/?datatype=jsonp&callback=find&token=126aa3a2be7b606dd727fd38d7d71f07&mobile=13209760000
func Phone(w http.ResponseWriter, r *http.Request) {
	log.Println("Request From", r.Host)
	mobile, err := getMobile(r.URL.Query().Get("mobile"))
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	if mobile.Ret == "ok" {
		w.Write([]byte(BuildMobileInfoString(mobile)))
		return
	}
	w.Write([]byte("GET ERROR"))
}

// CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build main.go
func main() {
	log.Println("服务器启动成功:9901")
	http.HandleFunc("/", Weather)
	http.HandleFunc("/p", Phone)
	http.ListenAndServe(":9901", nil)
}
