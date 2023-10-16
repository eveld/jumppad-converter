package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/jumppad-labs/hclconfig/convert"
	"github.com/jumppad-labs/instruqt-to-jumppad/model"
)

// blueprint
func generateBlueprint(track *model.Track) *hclwrite.Block {
	blueprintBlock := hclwrite.NewBlock("resource", []string{"blueprint", "blueprint"})
	blueprintBody := blueprintBlock.Body()

	// title
	title, err := convert.GoToCtyValue(track.Title)
	if err != nil {
		log.Fatal(err)
	}

	blueprintBody.SetAttributeValue("title", title)

	// slug
	slug, err := convert.GoToCtyValue(track.Slug)
	if err != nil {
		log.Fatal(err)
	}

	blueprintBody.SetAttributeValue("slug", slug)

	// icon
	if track.Icon != "" {
		icon, err := convert.GoToCtyValue(strings.Trim(track.Icon, "\n\r"))
		if err != nil {
			log.Fatal(err)
		}

		blueprintBody.SetAttributeValue("icon", icon)
	}

	// organization
	organization, err := convert.GoToCtyValue(track.Owner)
	if err != nil {
		log.Fatal(err)
	}

	blueprintBody.SetAttributeValue("organization", organization)

	// authors
	if len(track.Developers) > 0 {
		developers := []string{}
		for _, email := range track.Developers {
			parts := strings.Split(email, "@")
			developers = append(developers, fmt.Sprintf("%s <%s>", parts[0], email))
		}

		authors, err := convert.GoToCtyValue(developers)
		if err != nil {
			log.Fatal(err)
		}

		blueprintBody.SetAttributeValue("authors", authors)
	}

	// tags
	tags, err := convert.GoToCtyValue(track.Tags)
	if err != nil {
		log.Fatal(err)
	}

	blueprintBody.SetAttributeValue("tags", tags)

	// summary
	if track.Teaser != "" {
		summary, err := convert.GoToCtyValue(track.Teaser)
		if err != nil {
			log.Fatal(err)
		}

		blueprintBody.SetAttributeValue("summary", summary)
	}

	if strings.Contains(track.Description, "\n") {
		// <<EOF description
		eof := hclwrite.Tokens{}

		lines := []string{"<<-EOF"}
		lines = append(lines, strings.Split(track.Description, "\n")...)
		lines = append(lines, "EOF")

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

		blueprintBody.SetAttributeRaw("description", eof)
	} else {
		// single line
		description, err := convert.GoToCtyValue(track.Description)
		if err != nil {
			log.Fatal(err)
		}

		blueprintBody.SetAttributeValue("description", description)
	}

	return blueprintBlock
}
