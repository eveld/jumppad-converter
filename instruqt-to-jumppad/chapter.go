package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/jumppad-labs/hclconfig/convert"
	"github.com/jumppad-labs/instruqt-to-jumppad/model"
)

func generateChapter(track *model.Track, challenges []model.Challenge) *hclwrite.Block {
	chapterBlock := hclwrite.NewBlock("resource", []string{"chapter", "chapter"})
	chapterBody := chapterBlock.Body()

	// title
	title, err := convert.GoToCtyValue(track.Title)
	if err != nil {
		log.Fatal(err)
	}

	chapterBody.SetAttributeValue("title", title)
	chapterBody.AppendNewline()

	// tasks
	tasks := hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte("{\n"),
		},
	}

	for _, challenge := range challenges {
		tasks = append(tasks, &hclwrite.Token{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte(fmt.Sprintf("  %s = resource.task.%s\n", challenge.Slug, challenge.Slug)),
		})
	}

	tasks = append(tasks, &hclwrite.Token{
		Type:  hclsyntax.TokenIdent,
		Bytes: []byte("}"),
	})

	chapterBody.SetAttributeRaw("tasks", tasks)
	chapterBody.AppendNewline()

	// pages
	err = os.MkdirAll("out/content", 0755)
	if err != nil {
		log.Fatal(err)
	}

	for index, challenge := range challenges {
		pageBlock := hclwrite.NewBlock("page", []string{challenge.Slug})
		pageBody := pageBlock.Body()

		// content
		assignmentfile := fmt.Sprintf("out/content/%s.md", challenge.Slug)
		err := os.WriteFile(assignmentfile, []byte(challenge.Assignment), 0755)
		if err != nil {
			log.Fatal(err)
		}

		content := hclwrite.Tokens{
			{
				Type:  hclsyntax.TokenIdent,
				Bytes: []byte("file(\"content/" + challenge.Slug + ".md\")"),
			},
		}

		pageBody.SetAttributeRaw("content", content)

		chapterBody.AppendBlock(pageBlock)

		if index != len(challenges)-1 {
			chapterBody.AppendNewline()
		}
	}

	return chapterBlock
}
