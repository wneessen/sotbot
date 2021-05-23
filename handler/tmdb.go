package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/api"
	"github.com/wneessen/sotbot/random"
	"net/http"
)

func TmdbRandMovie(h *http.Client, a string) (*discordgo.MessageEmbed, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.MovieDb",
	})

	randPage, err := random.Number(500)
	if err != nil {
		l.Errorf("Random number generation failed: %v", err)
		return &discordgo.MessageEmbed{}, err
	}
	u := fmt.Sprintf("/3/discover/movie?page=%d&region=DE&api_key=%v", randPage, a)
	movieObj, err := api.GetTMDbMovie(h, u)
	if err != nil {
		l.Errorf("Could not fetch random movie: %v", err)
		return &discordgo.MessageEmbed{}, err
	}

	responseEmbed := &discordgo.MessageEmbed{
		Title: movieObj.Title,
		Description: fmt.Sprintf("%v\n\nRelease date: %v\nScore: %.0f%%", movieObj.Overview,
			movieObj.ReleaseDate, (movieObj.AvgVote * 10)),
		Type: discordgo.EmbedTypeImage,
		Image: &discordgo.MessageEmbedImage{
			URL:   fmt.Sprintf("https://image.tmdb.org/t/p/w300%v", movieObj.PosterPath),
			Width: 300,
		},
		URL: fmt.Sprintf("https://www.themoviedb.org/movie/%v", movieObj.Id),
	}

	return responseEmbed, nil
}
