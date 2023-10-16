package main

import (
	"log"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/jumppad-labs/hclconfig/convert"
	"github.com/jumppad-labs/instruqt-to-jumppad/model"
)

/*
	resource "book" "terraform_foundations" {
	  title = "Terraform Foundations"

	  chapters = [
	    resource.chapter.terraform_cli
	  ]
	}
*/
func generateBook(track *model.Track) *hclwrite.Block {
	bookBlock := hclwrite.NewBlock("resource", []string{"book", "book"})
	bookBody := bookBlock.Body()

	// title
	title, err := convert.GoToCtyValue(track.Title)
	if err != nil {
		log.Fatal(err)
	}

	bookBody.SetAttributeValue("title", title)
	bookBody.AppendNewline()

	// chapters
	chapters := hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte("[\n"),
		},
		{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte("   resource.chapter.chapter\n"),
		},
		{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte(" ]"),
		},
	}

	bookBody.SetAttributeRaw("chapters", chapters)

	return bookBlock
}
