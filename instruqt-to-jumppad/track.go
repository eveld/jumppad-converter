package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/instruqt/converter/model"
	"github.com/jumppad-labs/hclconfig/convert"
)

func GenerateTrack(track *model.Track, challenges []model.Challenge) *hclwrite.Block {
	trackBlock := hclwrite.NewBlock("resource", []string{"track", track.Slug})
	trackBody := trackBlock.Body()

	title, err := convert.GoToCtyValue(track.Title)
	if err != nil {
		log.Fatal(err)
	}
	trackBody.SetAttributeValue("title", title)

	owner, err := convert.GoToCtyValue(track.Owner)
	if err != nil {
		log.Fatal(err)
	}
	trackBody.SetAttributeValue("owner", owner)

	teaser, err := convert.GoToCtyValue(strings.TrimRight(track.Teaser, "\n"))
	if err != nil {
		log.Fatal(err)
	}
	trackBody.SetAttributeValue("teaser", teaser)

	if strings.Contains(track.Description, "\n") {
		// <<EOF description
		eof := hclwrite.Tokens{}

		lines := []string{"<<-EOF"}
		lines = append(lines, strings.Split(track.Description, "\n")...)
		lines = append(lines, " EOF")

		for index, line := range lines {
			space := "   "
			newline := "\n"

			if index == 0 {
				space = ""
			}

			if index == len(lines)-2 && line == "" {
				continue
			}

			if index == len(lines)-1 {
				space = ""
				newline = ""
			}

			eof = append(eof, &hclwrite.Token{
				Type:  hclsyntax.TokenIdent,
				Bytes: []byte(space + line + newline),
			})
		}

		trackBody.SetAttributeRaw("description", eof)
	} else {
		// single line
		description, err := convert.GoToCtyValue(track.Description)
		if err != nil {
			log.Fatal(err)
		}

		trackBody.SetAttributeValue("description", description)
	}

	tags, err := convert.GoToCtyValue(track.Tags)
	if err != nil {
		log.Fatal(err)
	}
	trackBody.SetAttributeValue("tags", tags)

	developers, err := convert.GoToCtyValue(track.Developers)
	if err != nil {
		log.Fatal(err)
	}
	trackBody.SetAttributeValue("developers", developers)

	icon, err := convert.GoToCtyValue(track.Icon)
	if err != nil {
		log.Fatal(err)
	}
	trackBody.SetAttributeValue("icon", icon)

	trackBody.AppendNewline()
	challengeList := hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte("{\n"),
		},
	}

	for index, challenge := range challenges {
		separator := ","
		if index == len(challenges)-1 {
			separator = ""
		}

		challengeList = append(challengeList, &hclwrite.Token{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte(fmt.Sprintf("   resource.challenge.%s%s\n", challenge.Slug, separator)),
		})
	}

	challengeList = append(challengeList, &hclwrite.Token{
		Type:  hclsyntax.TokenIdent,
		Bytes: []byte(" }"),
	})

	trackBody.SetAttributeRaw("challenges", challengeList)

	return trackBlock
}
