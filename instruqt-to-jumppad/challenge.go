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

func GenerateChallenge(challenge *model.Challenge) *hclwrite.Block {
	challengeBlock := hclwrite.NewBlock("resource", []string{"challenge", challenge.Slug})
	challengeBody := challengeBlock.Body()

	// title
	title, err := convert.GoToCtyValue(challenge.Title)
	if err != nil {
		log.Fatal(err)
	}
	challengeBody.SetAttributeValue("title", title)

	// teaser
	teaser, err := convert.GoToCtyValue(strings.TrimRight(challenge.Teaser, "\n"))
	if err != nil {
		log.Fatal(err)
	}
	challengeBody.SetAttributeValue("teaser", teaser)

	// assignment
	err = os.WriteFile(fmt.Sprintf("out/assignments/%s.mdx", challenge.Slug), []byte(challenge.Assignment), 0755)
	if err != nil {
		log.Fatal(err)
	}

	assignment := hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte(fmt.Sprintf("file(\"assignments/%s.mdx\")", challenge.Slug)),
		},
	}
	challengeBody.SetAttributeRaw("assignment", assignment)
	challengeBody.AppendNewline()

	// tabs
	tabs := hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte("{\n"),
		},
	}

	for index, _ := range challenge.Tabs {
		separator := ","
		if index == len(challenge.Tabs)-1 {
			separator = ""
		}

		tabs = append(tabs, &hclwrite.Token{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte(fmt.Sprintf("   resource.tab.%s_%d%s\n", challenge.Slug, index, separator)),
		})
	}

	tabs = append(tabs, &hclwrite.Token{
		Type:  hclsyntax.TokenIdent,
		Bytes: []byte(" }"),
	})

	challengeBody.SetAttributeRaw("tabs", tabs)

	// note
	for _, note := range challenge.Notes {
		challengeBody.AppendNewline()
		block := hclwrite.NewBlock("note", []string{})
		body := block.Body()

		noteType, err := convert.GoToCtyValue(note.Type)
		if err != nil {
			log.Fatal(err)
		}
		body.SetAttributeValue("type", noteType)

		if note.Contents != "" {
			if strings.Contains(note.Contents, "\n") {
				// <<EOF contents
				eof := hclwrite.Tokens{}

				lines := []string{"<<-EOF"}
				lines = append(lines, strings.Split(note.Contents, "\n")...)
				lines = append(lines, "   EOF")

				for index, line := range lines {
					space := "     "
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

				body.SetAttributeRaw("contents", eof)
			} else {
				// single line
				description, err := convert.GoToCtyValue(note.Contents)
				if err != nil {
					log.Fatal(err)
				}

				body.SetAttributeValue("contents", description)
			}
		}

		if note.URL != "" {
			url, err := convert.GoToCtyValue(note.URL)
			if err != nil {
				log.Fatal(err)
			}
			body.SetAttributeValue("url", url)
		}

		challengeBody.AppendBlock(block)
	}

	for _, s := range challenge.Scripts {
		challengeBody.AppendNewline()
		block := challengeBody.AppendNewBlock(s.Type, nil)
		body := block.Body()

		target, err := convert.GoToCtyValue(s.Target)
		if err != nil {
			log.Fatal(err)
		}
		body.SetAttributeValue("target", target)

		setupfile := fmt.Sprintf("out/scripts/%s/%s_%s.sh", challenge.Slug, s.Type, s.Target)
		err = os.WriteFile(setupfile, []byte(s.Content), 0755)
		if err != nil {
			log.Fatal(err)
		}

		script, err := convert.GoToCtyValue(setupfile)
		if err != nil {
			log.Fatal(err)
		}
		body.SetAttributeValue("script", script)
	}

	return challengeBlock
}
