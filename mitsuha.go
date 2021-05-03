package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"github.com/slack-go/slack"
)


type WeatherResult struct {
	Name     string `json:"city_name"`
	Data 	[]*Data `json:"data"`
}

type Data struct {
	Datetime string  `json:"datetime"`
    Wind_spd float64  `json:"wind_spd"`
    Wind_cdir string `json:"wind_cdir"`
	Temp float64 	  `json:"temp"`
	Max_temp float64  `json:"max_temp"`
	Min_temp float64  `json:"min_temp"`
	Pop int			 `json:"pop"`//降水確率
	Weather struct {
		Description string `json:"description"`
	} `json:"weather"`
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

func WhetherForecast() (forecast_today string ,forecast_tomorrow string) {
    url := "https://weatherbit-v1-mashape.p.rapidapi.com/forecast/daily?lat=35.026244&lon=135.780822&lang=ja"

    response, err := http.NewRequest("GET", url, nil)

	response.Header.Add("x-rapidapi-key", "a9959c2822mshb7fcf665a4b5bddp16bc17jsn2392fa173805")
	response.Header.Add("x-rapidapi-host", "weatherbit-v1-mashape.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(response)

    if err != nil {   // エラーハンドリング
        log.Fatalf("Connection Error: %v", err)
    }

    // 遅延
    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        log.Fatalf("Connection Error: %v", err)
    }

    jsonBytes := ([]byte)(body)
    data := new(WeatherResult)
    if err := json.Unmarshal(jsonBytes, data); err != nil {
        log.Fatalf("Connection Error: %v", err)
    }

    datetime_today := data.Data[0].Datetime
    weather_today := data.Data[0].Weather.Description

	Temp_today := data.Data[0].Temp
    temp_today := strconv.FormatFloat(Temp_today, 'f', 1, 64)//本日の平均気温
	Max_temp_today := data.Data[0].Max_temp
    max_temp_today := strconv.FormatFloat(Max_temp_today, 'f', 1, 64)//本日の最高気温
	Min_temp_today := data.Data[0].Min_temp
    min_temp_today := strconv.FormatFloat(Min_temp_today, 'f', 1, 64)//本日の最低気温

	Pop_today := data.Data[0].Pop
	pop_today := strconv.Itoa(Pop_today)//本日の降水確率

    Wind_spd_today := data.Data[0].Wind_spd
    wind_spd_today := strconv.FormatFloat(Wind_spd_today, 'f', 1, 64)//本日の平均風速
    wind_cdir_today := data.Data[0].Wind_cdir

    datetime_tomorrow := data.Data[1].Datetime
    weathet_tomorrow := data.Data[1].Weather.Description

	Temt_tomorrow := data.Data[1].Temp
    temt_tomorrow := strconv.FormatFloat(Temt_tomorrow, 'f', 1, 64)//本日の平均気温
	Max_temt_tomorrow := data.Data[1].Max_temp
    max_temt_tomorrow := strconv.FormatFloat(Max_temt_tomorrow, 'f', 1, 64)//本日の最高気温
	Min_temt_tomorrow := data.Data[1].Min_temp
    min_temt_tomorrow := strconv.FormatFloat(Min_temt_tomorrow, 'f', 1, 64)//本日の最低気温

	Pot_tomorrow := data.Data[1].Pop
	pot_tomorrow := strconv.Itoa(Pot_tomorrow)//本日の降水確率

    Wind_spd_tomorrow := data.Data[0].Wind_spd
    wind_spd_tomorrow := strconv.FormatFloat(Wind_spd_tomorrow, 'f', 1, 64)//本日の平均風速
    wind_cdir_tomorrow := data.Data[0].Wind_cdir

    forecast_today = "日付: "+datetime_today+"\n天気: "+weather_today+"\n平均気温: "+temp_today+"°C\n最高気温: "+max_temp_today+"°C\n最低気温: "+min_temp_today+"°C\n降水確率: "+pop_today+"%\n平均風速: "+wind_spd_today+"m/s\n風向き:  "+wind_cdir_today
    forecast_tomorrow = "日付: "+datetime_tomorrow+"\n天気: "+weathet_tomorrow+"\n平均気温: "+temt_tomorrow+"°C\n最高気温: "+max_temt_tomorrow+"°C\n最低気温: "+min_temt_tomorrow+"°C\n降水確率: "+pot_tomorrow+"%\n平均風速: "+wind_spd_tomorrow+"m/s\n風向き:  "+wind_cdir_tomorrow

    return forecast_today,forecast_tomorrow
}

func Reglarsend(){
    morning := time.Date(2021, 5, 1, 7, 00, 0, 0, time.Local)//朝の投稿時間設定
    night := time.Date(2021, 5, 1, 20, 00, 0, 0, time.Local)//の投稿時間設定

    var now = time.Now()
    if (now.Hour() == morning.Hour() &&
    now.Minute() == morning.Minute()){
        forecast_today,_ := WhetherForecast()
        api.PostMessage(
            ChannelID,
            slack.MsgOptionText("朝の天気予報の時間やよ！\n今日の京都市の天気はこんな感じ！\n"+forecast_today, false),
        )
        time.Sleep(time.Minute)//一分間停止
    }else if (now.Hour() == night.Hour() &&
    now.Minute() == night.Minute()){
        _,forecast_tomorrow:= WhetherForecast()
        api.PostMessage(
            ChannelID,
            slack.MsgOptionText("夜の天気予報の時間やよ！\n明日の京都市の天気はこんな感じ！\n"+forecast_tomorrow, false),
        )
        time.Sleep(time.Minute)//一分間停止
    }
}

//ListenTo excute functions under suitable conditions
func ListenTo() {
	switch {
        case strings.Contains(EV.Text,"君の名は。"):
		    RTM.SendMessage(RTM.NewOutgoingMessage("三葉！名前は三葉！", EV.Channel))
		    return
	}
}
//scp -r slackbots/mitsuha.go replica@kmc.gr.jp:~/Program/slackbots/