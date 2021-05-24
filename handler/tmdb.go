package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/ryanbradynd05/go-tmdb"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/sotbot/random"
	"strconv"
)

func TMDbRandMovie(t *tmdb.TMDb) (*discordgo.MessageEmbed, error) {
	l := log.WithFields(log.Fields{
		"action": "handler.TMDbRandMovie",
	})

	randPage, err := random.Number(500)
	if err != nil {
		l.Errorf("Random number generation failed: %v", err)
		return &discordgo.MessageEmbed{}, err
	}

	searchOpts := make(map[string]string, 0)
	searchOpts["region"] = "DE"
	searchOpts["page"] = strconv.FormatInt(int64(randPage), 10)
	movieResult, err := t.DiscoverMovie(searchOpts)
	if err != nil {
		l.Errorf("Failed to look up TMDB: %v", err)
		return &discordgo.MessageEmbed{}, err
	}

	numResults := len(movieResult.Results)
	randResult, err := random.Number(numResults)
	if err != nil {
		l.Errorf("Failed to generate random number from number of results: %v", err)
		return &discordgo.MessageEmbed{}, err
	}

	randMovie := movieResult.Results[randResult]
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
