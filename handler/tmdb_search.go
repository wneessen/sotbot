package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/ryanbradynd05/go-tmdb"
	log "github.com/sirupsen/logrus"
)

func TMDbSearchMovie(t *tmdb.TMDb, q []string) (*discordgo.MessageEmbed, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.TMDbSearchMovie",
	})

	var searchString string
	for _, keyWord := range q {
		searchString = fmt.Sprintf("%v %v", searchString, keyWord)
	}
	movieResult, err := t.SearchMovie(searchString, nil)
	if err != nil {
		l.Errorf("Failed to look up TMDB: %v", err)
		return &discordgo.MessageEmbed{}, err
	}

	if len(movieResult.Results) == 0 {
		return &discordgo.MessageEmbed{}, fmt.Errorf("No matching movie found")
	}

	randMovie := movieResult.Results[0]
	responseEmbed := &discordgo.MessageEmbed{
		Title: randMovie.Title,
		Description: fmt.Sprintf("%v\n\n**Release date:** %v\n**Score:** %.0f%%", randMovie.Overview,
			randMovie.ReleaseDate, (randMovie.VoteAverage * 10)),
		Type: discordgo.EmbedTypeImage,
		Image: &discordgo.MessageEmbedImage{
			URL:   fmt.Sprintf("https://image.tmdb.org/t/p/w300%v", randMovie.PosterPath),
			Width: 300,
		},
		URL: fmt.Sprintf("https://www.themoviedb.org/movie/%v", randMovie.ID),
	}

	return responseEmbed, nil
}
