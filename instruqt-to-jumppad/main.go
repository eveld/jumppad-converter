package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"

	"github.com/jumppad-labs/hclconfig/convert"
	"github.com/jumppad-labs/instruqt-to-jumppad/model"
)

func main() {
	args := os.Args[1:]

	generateEnvironment := true

	dir, err := filepath.Abs(args[0])
	if err != nil {
		log.Fatal(err)
	}

	// track
	track, err := readTrack(dir)
	if err != nil {
		log.Fatal(err)
	}

	// challenges
	challenges, err := readChallenges(dir)
	if err != nil {
		log.Fatal(err)
	}

	// config
	env, err := readConfig(dir)
	if err != nil {
		log.Fatal(err)
	}

	err = os.RemoveAll("out")
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll("out/assets", 0755)
	if err != nil {
		log.Fatal(err)
	}

	f := hclwrite.NewEmptyFile()
	root := f.Body()

	blueprint := generateBlueprint(track)
	root.AppendBlock(blueprint)
	root.AppendNewline()

	docsComment := generateComment("Docs")
	root.AppendUnstructuredTokens(docsComment)

	docs := generateDocs(track)
	root.AppendBlock(docs)
	root.AppendNewline()

	bookComment := generateComment("Book")
	root.AppendUnstructuredTokens(bookComment)

	book := generateBook(track)
	root.AppendBlock(book)
	root.AppendNewline()

	chapterComment := generateComment("Chapter")
	root.AppendUnstructuredTokens(chapterComment)

	chapter := generateChapter(track, challenges)
	root.AppendBlock(chapter)
	root.AppendNewline()

	tasksComment := generateComment("Tasks")
	root.AppendUnstructuredTokens(tasksComment)

	for index, challenge := range challenges {
		var previous *model.Challenge
		if index == 0 {
			previous = nil
		} else {
			previous = &challenges[index-1]
		}

		task := generateTask(&challenge, previous)
		root.AppendBlock(task)

		if index != len(challenges)-1 {
			root.AppendNewline()
		}
	}

	if generateEnvironment {
		root.AppendNewline()

		envComment := generateComment("Environment")
		root.AppendUnstructuredTokens(envComment)

		// Network
		subnet, err := convert.GoToCtyValue("10.0.5.0/16")
		if err != nil {
			log.Fatal(err)
		}

		networkBlock := root.AppendNewBlock("resource", []string{"network", "main"})
		networkBody := networkBlock.Body()
		networkBody.SetAttributeValue("subnet", subnet)

		// Containers
		for _, container := range env.Containers {
			root.AppendNewline()

			containerBlock := root.AppendNewBlock("resource", []string{"container", container.Name})
			containerBody := containerBlock.Body()

			image, err := convert.GoToCtyValue(container.Image)
			if err != nil {
				log.Fatal(err)
			}

			imageBlock := containerBody.AppendNewBlock("image", nil)
			imageBody := imageBlock.Body()
			imageBody.SetAttributeValue("name", image)

			containerBody.AppendNewline()
			networkBlock := containerBody.AppendNewBlock("network", nil)
			networkBody := networkBlock.Body()

			id := hclwrite.Tokens{
				{
					Type:  hclsyntax.TokenIdent,
					Bytes: []byte(fmt.Sprintf("resource.network.%s.id", "main")),
				},
			}

			networkBody.SetAttributeRaw("id", id)

			// entrypoint
			if container.Entrypoint != "" {
				entrypoint, err := convert.GoToCtyValue([]string{container.Entrypoint})
				if err != nil {
					log.Fatal(err)
				}

				containerBody.AppendNewline()
				containerBody.SetAttributeValue("entrypoint", entrypoint)
			}

			// environment
			// append shell variable
			if container.Environment == nil {
				container.Environment = map[string]string{}
			}
			container.Environment["SHELL"] = container.Shell

			environment, err := convert.GoToCtyValue(container.Environment)
			if err != nil {
				log.Fatal(err)
			}

			containerBody.AppendNewline()
			containerBody.SetAttributeValue("environment", environment)
		}
	}

	// Output file
	err = os.WriteFile("out/main.hcl", f.Bytes(), 0755)
	if err != nil {
		log.Fatal(err)
	}
}

func generateComment(comment string) hclwrite.Tokens {
	return hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenComment,
			Bytes: []byte("//\n"),
		},
		{
			Type:  hclsyntax.TokenComment,
			Bytes: []byte(fmt.Sprintf("// %s\n", comment)),
		},
		{
			Type:  hclsyntax.TokenComment,
			Bytes: []byte("//\n"),
		},
	}
}
