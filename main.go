package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
)

type Location struct {
	Name    string `json:"name"`
	Region  string `json:"region"`
	Country string `json:"country"`
}

type Condition struct {
	Text string `json:"text"`
}

type Current struct {
	TempC     float64   `json:"temp_c"`
	Condition Condition `json:"condition"`
}

type Hour []struct {
	TimeEpoch    int64     `json:"time_epoch"`
	TempC        float64   `json:"temp_c"`
	Condition    Condition `json:"condition"`
	ChanceOfRain float64   `json:"chance_of_rain"`
}

type Forecastday []struct {
	Hour Hour `json:"hour"`
}

type Forecast struct {
	Forecastday Forecastday `json:"forecastday"`
}

type Weather struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
	Forecast Forecast `json:"forecast"`
}

func main() {

	loc := "Mersin"

	// first is argument itself
	if len(os.Args) >= 2 {
		loc = os.Args[1]
	}

	res, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=ed43a058b84f45259de131426230907&q=" + loc + "&days=1&aqi=no&alerts=no")

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic(res.Status)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)

	if err != nil {
		panic(err)
	}

	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour

	fmt.Printf(
		"%s, %s, %s, %.0f C°, %s\n\n",
		location.Name,
		location.Region,
		location.Country,
		current.TempC,
		current.Condition.Text,
	)

	fmt.Println("Hour - Temperature  Chance of Rain     Condition")
	for _, hour := range hours {

		date := time.Unix(hour.TimeEpoch, 0)

		if date.Before(time.Now()) {
			continue
		}

		message := fmt.Sprintf(
			"%s - %.0f C° \t %.0f%% \t\t %s \t\n",
			date.Format("15:04:05"),
			hour.TempC,
			hour.ChanceOfRain,
			hour.Condition,
		)

		if hour.ChanceOfRain > 50 {
			color.Blue(message)
		} else {
			color.Red(message)
		}
	}
}
