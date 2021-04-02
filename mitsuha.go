package main

import (
	"fmt"
    "strconv"
	"strings"
	"encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
	"github.com/slack-go/slack"
)

type WeatherResult struct {
    Coord struct {
        Lon float64 `json:"lon"`
        Lat float64 `json:"lat"`
    } `json:"coord"`
    Weather []struct {
        ID          int    `json:"id"`
        Main        string `json:"main"`
        Description string `json:"description"`
        Icon        string `json:"icon"`
    } `json:"weather"`
    Base string `json:"base"`
    Main struct {
        Temp      float64 `json:"temp"`
        FeelsLike float64 `json:"feels_like"`
        TempMin   float64 `json:"temp_min"`
        TempMax   float64 `json:"temp_max"`
        Pressure  int     `json:"pressure"`
        Humidity  int     `json:"humidity"`
    } `json:"main"`
    Visibility int `json:"visibility"`
    Wind       struct {
        Speed float64 `json:"speed"`
    } `json:"wind"`
    Clouds struct {
        All int `json:"all"`
    } `json:"clouds"`
    Dt  int `json:"dt"`
    Sys struct {
        Type    int    `json:"type"`
        ID      int    `json:"id"`
        Country string `json:"country"`
        Sunrise int    `json:"sunrise"`
        Sunset  int    `json:"sunset"`
    } `json:"sys"`
    Timezone int    `json:"timezone"`
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Cod      json.Number    `json:"cod"`
}

//EV put new slack events
var EV *slack.MessageEvent

//RTM use for sending events to slack
var RTM *slack.RTM

//BotToken Put your slackbot token here
const BotToken string = "xoxb-3069876617-1921979534688-IOyPwRPt882DUliK2d2WXITe"

//DefaultChannel Put your default channel
const DefaultChannel string = "#replica-memo"

func main() {
	var api *slack.Client = slack.New(BotToken)

	RTM = api.NewRTM()

	go RTM.ManageConnection()

	for msg := range RTM.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			fmt.Printf("Start connection with Slack\n")
		case *slack.MessageEvent:
			EV = ev
			ListenTo()
		}
	}
}

//ListenTo excute functions under suitable conditions
func ListenTo() {
    values := url.Values{}
    baseUrl := "http://api.openweathermap.org/data/2.5/weather?"

    // query
    values.Add("appid", "3e2c765871b9be3be5864647fb2d23da")    // OpenWeatherのAPIKey
    values.Add("id","1857910") //京都市
    values.Add("lang","ja") //日本語

    response, err := http.Get(baseUrl + values.Encode())

    description := ""

    if err != nil {   // エラーハンドリング
        log.Fatalf("Connection Error: %v", err)
         description="不明"
    }

   // 遅延
    defer response.Body.Close()

    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatalf("Connection Error: %v", err)
        description="不明"
    }

    jsonBytes := ([]byte)(body)
    data := new(WeatherResult)
    if err := json.Unmarshal(jsonBytes, data); err != nil {
        log.Fatalf("Connection Error: %v", err)
    }

    if data.Weather != nil {
        description= data.Weather[0].Description
    }

    weather_comment := "" //天気についてのコメント
    weather := data.Weather[0].Main

    if weather == "Thunderstorm" {
        weather_comment = "雷怖い！"
    }else if weather == "Drizzle" {
        weather_comment = "周りが見えづらいから気をつけて！"
    }else if weather == "Rain" {
        weather_comment = "外出する時は傘を忘れないでね！"
    }else if weather == "Snow" {
        weather_comment = "雪だるま作ろう！"
    }else if weather == "Clear" {
        weather_comment = "洗濯日和！"
    }else if weather == "Clouds" {
        weather_comment = "雲がいっぱい！"
    }else {
        weather_comment = "外は危険がいっぱい！"
    }

    Temp := data.Main.Temp //現在の気温
    Temp -= 273.15
    temp := strconv.FormatFloat(Temp, 'f', 1, 64)

    temp_comment := ""

    if Temp < 0 {
        temp_comment = "寒すぎ！"
    }else if Temp < 5 {
        temp_comment = "寒い！"
    }else if Temp < 10 {
        temp_comment = "少し寒い！"
    }else if Temp < 15 {
        temp_comment ="肌寒いね！"
    }else if Temp < 20 {
        temp_comment = "涼しい！"
    }else if Temp < 25 {
        temp_comment = "暖かくて快適！"
    }else if Temp < 30 {
        temp_comment = "少し暑〜い！"
    }else if Temp < 35 {
        temp_comment = "暑い〜!"
    }else {
        temp_comment = "暑すぎて死にそう！"
    }

    wind_comment := ""
    wind := data.Wind.Speed

    if wind < 1.5 {
        wind_comment = "風もほとんど吹いてなさそう！"
    }else if wind < 5.0 {
        wind_comment = "そよ風が吹いてるよ！"
    }else if wind < 10 {
        wind_comment = "少し風が強いかも！"
    }else if wind < 15 {
        wind_comment = "強い風が吹いてるから気をつけて！"
    }else if wind < 20 {
        wind_comment = "風が強いから外出は危険！"
    }else if wind < 25 {
        wind_comment = "風で家が壊れそう！"
    }else{
        wind_comment = "強風で大きな被害が出そう！避難して！"
    }

	switch {
	    case strings.Contains("こんにちはこんにちは", EV.Text):
		    RTM.SendMessage(RTM.NewOutgoingMessage("こんにちはこんにちは、<@"+EV.User+">さん！今の京都市の天気は"+description+"やよ！"+weather_comment+"気温は"+temp+"度。"+temp_comment+wind_comment, EV.Channel))
		    return
        case strings.Contains("君の名は。", EV.Text):
		    RTM.SendMessage(RTM.NewOutgoingMessage("三葉！名前は三葉！", EV.Channel))
		    return
	}
}