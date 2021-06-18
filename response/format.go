package response

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"sort"
)

func FormatEmFields(vo map[string]int64) []*discordgo.MessageEmbedField {
	// Prepare the output
	p := message.NewPrinter(language.German)
	var emFields []*discordgo.MessageEmbedField
	var keyNames []string
	for k := range vo {
		keyNames = append(keyNames, k)
	}
	sort.Strings(keyNames)
	for _, k := range keyNames {
		v := vo[k]
		if v != 0 {
			emFields = append(emFields, &discordgo.MessageEmbedField{
				Name:   fmt.Sprintf("%v %v", Icon(k), IconKey(k)),
				Value:  fmt.Sprintf("%v**%v** %v", BalanceIcon(k, v), p.Sprintf("%d", v), IconValue(k)),
				Inline: true,
			})
		}
	}
	for len(emFields)%3 != 0 {
		emFields = append(emFields, &discordgo.MessageEmbedField{
			Value:  "\U0000FEFF",
			Name:   "\U0000FEFF",
			Inline: true,
		})
	}

	return emFields
}
