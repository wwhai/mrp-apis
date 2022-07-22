package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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

func BuildWeatherInfoString(w WeatherInfo) string {
	s := ""
	for _, v := range w.Data.Forecast {
		s += v.Date + ", "
		s += v.High + ", "
		s += v.Low + ", "
		s += v.Fengxiang + "\n"
	}
	return s
}
func getWeather() (weatherInfo WeatherInfo, err error) {
	// http://wthrcdn.etouch.cn/weather_mini?city=广州
	result, err1 := HttpGet(httpC, "http://wthrcdn.etouch.cn/weather_mini?city=%E5%B9%BF%E5%B7%9E")
	if err1 != nil {
		err = err1
		return
	}
	if err = json.Unmarshal([]byte(result), &weatherInfo); err != nil {
		return
	}
	return
}
func index(w http.ResponseWriter, r *http.Request) {
	weatherInfo, err := getWeather()
	if err != nil {
		return
	}
	w.Write([]byte(BuildWeatherInfoString(weatherInfo)))
}
// CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build main.go
func main() {
	log.Println("服务器启动成功:9901")
	http.HandleFunc("/", index)
	http.ListenAndServe(":9901", nil)
}
