package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/vascocosta/owm"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"time"
)

func GetCurrentWeather(o *owm.Client, loc string) (*discordgo.MessageEmbed, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.GetCurrentWeather",
	})

	l.Debugf("Trying to fetch weather conditions for: %v", loc)
	curWeather, err := o.WeatherByName(loc, "metric")
	if err != nil {
		l.Errorf("Failed to look up weather data: %v", err)
		return &discordgo.MessageEmbed{}, err
	}

	var emFields []*discordgo.MessageEmbedField
	p := message.NewPrinter(language.German)
	emFields = append(emFields, &discordgo.MessageEmbedField{
		Name:   "Current conditions",
		Value:  curWeather.Weather[0].Description,
		Inline: false,
	})
	emFields = append(emFields, &discordgo.MessageEmbedField{
		Name:   "Temperature",
		Value:  p.Sprintf("%.1f°C", curWeather.Main.Temp),
		Inline: true,
	})
	emFields = append(emFields, &discordgo.MessageEmbedField{
		Name:   "Min./Max.",
		Value:  p.Sprintf("%.1f°C / %.1f°C", curWeather.Main.TempMin, curWeather.Main.TempMax),
		Inline: true,
	})
	emFields = append(emFields, &discordgo.MessageEmbedField{
		Name:   "Humidity",
		Value:  p.Sprintf("%d%%", curWeather.Main.Humidity),
		Inline: true,
	})
	emFields = append(emFields, &discordgo.MessageEmbedField{
		Name:   "Atmosph. pressure",
		Value:  p.Sprintf("%.0fhPa", curWeather.Main.Pressure),
		Inline: true,
	})
	emFields = append(emFields, &discordgo.MessageEmbedField{
		Name:   "Sunrise",
		Value:  time.Unix(int64(curWeather.Sys.Sunrise), 0).Format("15:04:05 MST"),
		Inline: true,
	})
	emFields = append(emFields, &discordgo.MessageEmbedField{
		Name:   "Sunset",
		Value:  time.Unix(int64(curWeather.Sys.Sunset), 0).Format("15:04:05 MST"),
		Inline: true,
	})

	responseEmbed := &discordgo.MessageEmbed{
		Type:  discordgo.EmbedTypeRich,
		Title: fmt.Sprintf("Weather summary for %s, %s", curWeather.Name, curWeather.Sys.Country),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: fmt.Sprintf("https://openweathermap.org/img/wn/%s@2x.png",
				curWeather.Weather[0].Icon),
		},
		Fields: emFields,
	}
	return responseEmbed, nil
}
