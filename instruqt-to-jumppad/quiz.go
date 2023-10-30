package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/instruqt/converter/model"
	"github.com/jumppad-labs/hclconfig/convert"
)

/*
	resource "quiz" "name" {
		title = "title"
		teaser = "teaser"

		option {
			answer = "answer"
			solution = true
		}

		option {
			answer = "answer"
			solution = false
		}

		option {
			answer = "answer"
		}
	}
*/
func (c *Config) GenerateQuiz(quiz *model.Challenge, tabs map[string]model.Tab) *hclwrite.Block {
	// create a directory for the quiz
	quizPath := fmt.Sprintf("out/%s", quiz.Slug)
	err := os.MkdirAll(quizPath, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// create a directory for the quiz notes
	err = os.MkdirAll(filepath.Join(quizPath, "notes"), 0755)
	if err != nil {
		log.Fatal(err)
	}

	quizBlock := hclwrite.NewBlock("resource", []string{"quiz", quiz.Slug})
	quizBody := quizBlock.Body()

	// title
	title, err := convert.GoToCtyValue(quiz.Title)
	if err != nil {
		log.Fatal(err)
	}

	quizBody.SetAttributeValue("title", title)

	// teaser
	teaser, err := convert.GoToCtyValue(quiz.Teaser)
	if err != nil {
		log.Fatal(err)
	}

	quizBody.SetAttributeValue("teaser", teaser)

	// assignment
	err = os.WriteFile(fmt.Sprintf("out/%s/assignment.mdx", quiz.Slug), []byte(quiz.Assignment), 0755)
	if err != nil {
		log.Fatal(err)
	}

	assignment := hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte(fmt.Sprintf("file(\"%s/assignment.mdx\")", quiz.Slug)),
		},
	}
	quizBody.SetAttributeRaw("assignment", assignment)

	// note
	for index, note := range quiz.Notes {
		quizBody.AppendNewline()
		block := hclwrite.NewBlock("note", []string{})
		body := block.Body()

		noteType, err := convert.GoToCtyValue(note.Type)
		if err != nil {
			log.Fatal(err)
		}
		body.SetAttributeValue("type", noteType)

		if note.Contents != "" {
			noteFile := fmt.Sprintf("%s/notes/note_%d.mdx", quiz.Slug, index)
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

		quizBody.AppendBlock(block)
	}

	// tabs
	quizBody.AppendNewline()
	quizTabs := hclwrite.Tokens{
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

		quizTabs = append(quizTabs, &hclwrite.Token{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte(fmt.Sprintf("   resource.tab.%s%s\n", slug, separator)),
		})

		index++
	}

	quizTabs = append(quizTabs, &hclwrite.Token{
		Type:  hclsyntax.TokenIdent,
		Bytes: []byte(" }"),
	})

	quizBody.SetAttributeRaw("tabs", quizTabs)

	// options
	for index, answer := range quiz.Answers {
		quizBody.AppendNewline()
		answerBlock := hclwrite.NewBlock("answer", nil)
		answerBody := answerBlock.Body()

		value, err := convert.GoToCtyValue(answer)
		if err != nil {
			log.Fatal(err)
		}

		answerBody.SetAttributeValue("value", value)

		for _, solution := range quiz.Solutions {
			if index == solution {
				solutionValue, err := convert.GoToCtyValue(true)
				if err != nil {
					log.Fatal(err)
				}

				answerBody.SetAttributeValue("solution", solutionValue)
				break
			}
		}

		quizBody.AppendBlock(answerBlock)
	}

	return quizBlock
}
