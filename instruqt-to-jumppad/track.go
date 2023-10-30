package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/instruqt/converter/model"
	"github.com/jumppad-labs/hclconfig/convert"
)

func (c *Config) GenerateTrack(track *model.Track, challenges []model.Challenge) *hclwrite.Block {
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

	// description
	err = os.WriteFile("out/description.mdx", []byte(track.Description), 0755)
	if err != nil {
		log.Fatal(err)
	}

	description := hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte("file(\"description.mdx\")"),
		},
	}
	trackBody.SetAttributeRaw("description", description)
	trackBody.AppendNewline()

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
