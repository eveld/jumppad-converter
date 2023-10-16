package main

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/jumppad-labs/instruqt-to-jumppad/model"
)

/*
	resource "docs" "docs" {
	  network {
	    id = resource.network.main.id
	  }

	  image {
	    name = "ghcr.io/jumppad-labs/docs:v0.4.0"
	  }

	  content = [
	    resource.book.terraform_foundations
	  ]

	 assets = "assets"
	}
*/
func generateDocs(track *model.Track) *hclwrite.Block {
	docsBlock := hclwrite.NewBlock("resource", []string{"docs", "docs"})
	docsBody := docsBlock.Body()

	// network
	networkBlock := docsBody.AppendNewBlock("network", nil)
	networkBody := networkBlock.Body()
	docsBody.AppendNewline()

	id := hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte(fmt.Sprintf("resource.network.%s.id", "main")),
		},
	}

	networkBody.SetAttributeRaw("id", id)

	// content
	content := hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte("[\n"),
		},
		{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte("   resource.book.book\n"),
		},
		{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte(" ]"),
		},
	}

	docsBody.SetAttributeRaw("content", content)
	docsBody.AppendNewline()

	// assets
	assets := hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte("\"${dir()}/assets\""),
		},
	}

	docsBody.SetAttributeRaw("assets", assets)

	return docsBlock
}
