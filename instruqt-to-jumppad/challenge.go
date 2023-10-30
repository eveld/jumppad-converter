package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/instruqt/converter/model"
	"github.com/jumppad-labs/hclconfig/convert"
)

func (c *Config) GenerateChallenge(challenge *model.Challenge, tabs map[string]model.Tab) *hclwrite.Block {
	challengeBlock := hclwrite.NewBlock("resource", []string{"challenge", challenge.Slug})
	challengeBody := challengeBlock.Body()

	// create a directory for the challenge
	challengePath := fmt.Sprintf("out/%s", challenge.Slug)
	err := os.MkdirAll(challengePath, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// create a directory for the challenge scripts
	err = os.MkdirAll(filepath.Join(challengePath, "scripts"), 0755)
	if err != nil {
		log.Fatal(err)
	}

	// create a directory for the challenge notes
	err = os.MkdirAll(filepath.Join(challengePath, "notes"), 0755)
	if err != nil {
		log.Fatal(err)
	}

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
	err = os.WriteFile(fmt.Sprintf("out/%s/assignment.mdx", challenge.Slug), []byte(challenge.Assignment), 0755)
	if err != nil {
		log.Fatal(err)
	}

	assignment := hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte(fmt.Sprintf("file(\"%s/assignment.mdx\")", challenge.Slug)),
		},
	}
	challengeBody.SetAttributeRaw("assignment", assignment)

	// note
	for index, note := range challenge.Notes {
		challengeBody.AppendNewline()
		block := hclwrite.NewBlock("note", []string{})
		body := block.Body()

		noteType, err := convert.GoToCtyValue(note.Type)
		if err != nil {
			log.Fatal(err)
		}
		body.SetAttributeValue("type", noteType)

		if note.Contents != "" {
			noteFile := fmt.Sprintf("%s/notes/note_%d.mdx", challenge.Slug, index)
			err = os.WriteFile(filepath.Join("out", noteFile), []byte(note.Contents), 0755)
			if err != nil {
				log.Fatal(err)
			}

			note := hclwrite.Tokens{
				{
					Type:  hclsyntax.TokenIdent,
					Bytes: []byte(fmt.Sprintf("file(\"%s\")", noteFile)),
				},
			}
			body.SetAttributeRaw("note", note)
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

	// tabs
	challengeBody.AppendNewline()
	challengeTabs := hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte("{\n"),
		},
	}

	index := 0
	for slug := range tabs {
		separator := ","
		if index == len(tabs)-1 {
			separator = ""
		}

		challengeTabs = append(challengeTabs, &hclwrite.Token{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte(fmt.Sprintf("   resource.tab.%s%s\n", slug, separator)),
		})

		index++
	}

	challengeTabs = append(challengeTabs, &hclwrite.Token{
		Type:  hclsyntax.TokenIdent,
		Bytes: []byte(" }"),
	})

	challengeBody.SetAttributeRaw("tabs", challengeTabs)

	for _, s := range challenge.Scripts {
		challengeBody.AppendNewline()
		block := challengeBody.AppendNewBlock(s.Type, nil)
		body := block.Body()

		resource, err := c.LookupResource(s.Target)
		if err != nil {
			log.Fatal(err)
		}

		target := hclwrite.Tokens{
			{
				Type:  hclsyntax.TokenIdent,
				Bytes: []byte(resource),
			},
		}
		body.SetAttributeRaw("target", target)

		scriptFile := fmt.Sprintf("%s/scripts/%s_%s.sh", challenge.Slug, s.Type, s.Target)
		err = os.WriteFile(filepath.Join("out", scriptFile), []byte(s.Content), 0755)
		if err != nil {
			log.Fatal(err)
		}

		script := hclwrite.Tokens{
			{
				Type:  hclsyntax.TokenIdent,
				Bytes: []byte(fmt.Sprintf("file(\"%s\")", scriptFile)),
			},
		}
		body.SetAttributeRaw("script", script)
	}

	return challengeBlock
}
