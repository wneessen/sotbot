package handler

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vascocosta/owm"
)

func GetCurrentWeather(o *owm.Client, loc []string) (string, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.GetCurrentWeather",
	})

	// Bot ist located in Cologne, Germany... so that's a default
	weatherLoc := "Köln"
	if len(loc) > 0 {
		weatherLoc = ""
		for i, loctext := range loc {
			if i > 0 {
				weatherLoc = fmt.Sprintf("%v %v", weatherLoc, loctext)
			} else {
				weatherLoc = loctext
			}
		}
	}

	curWeather, err := o.WeatherByName(weatherLoc, "metric")
	if err != nil {
		l.Errorf("Failed to look up weather data: %v", err)
		return "", err
	}

	responseMsg := fmt.Sprintf("Current weather in %v. Temperatur: %.2f°C (Min: %.1f°C, Max: %.1f°C), "+
		"Humidity: %d%%, Pressure: %.1fhPa.", curWeather.Name, curWeather.Main.Temp, curWeather.Main.TempMin,
		curWeather.Main.TempMax, curWeather.Main.Humidity, curWeather.Main.Pressure)
	return responseMsg, nil
}
