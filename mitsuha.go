package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

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

//#replica-memoのID
const ChannelID string ="C01SGM52Y6Q"

var api *slack.Client = slack.New(BotToken)

func main() {

	RTM = api.NewRTM()

	go RTM.ManageConnection()

    for{
        for msg := range RTM.IncomingEvents {
            switch ev := msg.Data.(type) {
                case *slack.ConnectedEvent:
                    fmt.Printf("Start connection with Slack\n")
                case *slack.MessageEvent:
                    EV = ev
                    ListenTo()
                default:
                    Reglarsend()
            }
        }
    }
}

func WhetherForecast() (forecast_now string ,forecast_today string) {
    values := url.Values{}
    baseUrl := "http://api.openweathermap.org/data/2.5/weather?"

    // query
    values.Add("appid", "3e2c765871b9be3be5864647fb2d23da")    // OpenWeatherのAPIKey
    values.Add("id","1857910") //京都市
    values.Add("lang","ja") //日本語

    response, err := http.Get(baseUrl + values.Encode())

    description := ""//天気

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
        weather_comment = "雷が鳴っているから気をつけて！"
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
    Temp := data.Main.Temp
    Temp -= 273.15
    temp := strconv.FormatFloat(Temp, 'f', 1, 64)//現在の気温(摂氏)

    Temp_max := data.Main.TempMax
    Temp_min := data.Main.TempMin
    Temp_max -= 273.15
    Temp_min -= 273.15
    temp_max := strconv.FormatFloat(Temp_max, 'f', 1, 64)//最高気温(摂氏)
    temp_min := strconv.FormatFloat(Temp_min, 'f', 1, 64)//最低気温(摂氏)

    Wind := data.Wind.Speed
    wind := strconv.FormatFloat(Wind, 'f', 1, 64)//現在の風速

    forecast_now = "天候: "+description+"\n気温: "+temp+"°C\n"+"風速: "+wind+"m/s\n"+weather_comment
    forecast_today = "天候: "+description+"\n最高気温: "+temp_max+"°C\n最低気温: "+temp_min+"°C\n風速: "+wind+"m/s\n"+weather_comment

    return forecast_now,forecast_today
}

func Reglarsend(){
    t := time.Date(2021, 5, 1, 7, 00, 0, 0, time.Local)//定時投稿の投稿時間設定
    var now = time.Now()
    if now.Hour() != t.Hour() ||
        now.Minute() != t.Minute() {
    }else{
        _,forecast_today := WhetherForecast()
        api.PostMessage(
            ChannelID,
            slack.MsgOptionText("朝の天気予報の時間やよ！\n今日の京都市の天気はこんな感じ！"+forecast_today, false),
        )

        time.Sleep(time.Minute)
    }
}

//ListenTo excute functions under suitable conditions
func ListenTo() {
	switch {
	    case strings.Contains(EV.Text,"天気"):
            forecast_now,_ := WhetherForecast()
		    RTM.SendMessage(RTM.NewOutgoingMessage("こんにちはこんにちは、<@"+EV.User+">さん！\n現在の京都市の天気はこんな感じやよ！\n"+forecast_now, EV.Channel))
		    return
        case strings.Contains(EV.Text,"君の名は。"):
		    RTM.SendMessage(RTM.NewOutgoingMessage("三葉！名前は三葉！", EV.Channel))
		    return
	}
}