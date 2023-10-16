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

func generateTask(challenge *model.Challenge, previous *model.Challenge) *hclwrite.Block {
	/*
			prerequisites = []

		  config {
		    user = "root"
		    target = variable.terraform_target
		  }

		  condition "binary_exists" {
		    description = "Terraform installed on path"

		    check {
		      script = file("checks/installation/manual_installation/binary_exists")
		      failure_message = "terraform binary not found on the PATH"
		    }

		    solve {
		      script = file("checks/installation/manual_installation/solve")
		      timeout = 60
		    }
		  }
	*/
	taskBlock := hclwrite.NewBlock("resource", []string{"task", challenge.Slug})
	taskBody := taskBlock.Body()

	if previous != nil {
		prerequisites := hclwrite.Tokens{
			{
				Type:  hclsyntax.TokenIdent,
				Bytes: []byte("[\n"),
			},
			{
				Type:  hclsyntax.TokenIdent,
				Bytes: []byte("   resource.task." + previous.Slug + ".id\n"),
			},
			{
				Type:  hclsyntax.TokenIdent,
				Bytes: []byte(" ]"),
			},
		}

		taskBody.SetAttributeRaw("prerequisites", prerequisites)
	} else {
		prerequisites := hclwrite.Tokens{
			{
				Type:  hclsyntax.TokenIdent,
				Bytes: []byte("[]"),
			},
		}
		taskBody.SetAttributeRaw("prerequisites", prerequisites)
	}

	// prerequisites

	taskBody.AppendNewline()

	// condition
	conditionBlock := hclwrite.NewBlock("condition", []string{challenge.Slug})
	conditionBody := conditionBlock.Body()

	// description
	description, err := convert.GoToCtyValue(challenge.Teaser)
	if err != nil {
		log.Fatal(err)
	}

	conditionBody.SetAttributeValue("description", description)
	conditionBody.AppendNewline()

	// scripts
	err = os.MkdirAll(fmt.Sprintf("out/scripts/%s", challenge.Slug), 0755)
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range challenge.Scripts {
		block := conditionBody.AppendNewBlock(s.Type, nil)
		body := block.Body()

		setupfile := fmt.Sprintf("out/scripts/%s/%s_%s", challenge.Slug, s.Type, s.Target)
		err := os.WriteFile(setupfile, []byte(s.Content), 0755)
		if err != nil {
			log.Fatal(err)
		}

		// script
		script, err := convert.GoToCtyValue(setupfile)
		if err != nil {
			log.Fatal(err)
		}

		body.SetAttributeValue("script", script)

		// target
		target := hclwrite.Tokens{
			{
				Type:  hclsyntax.TokenIdent,
				Bytes: []byte("resource.container." + s.Target + ".id"),
			},
		}

		body.SetAttributeRaw("target", target)

		conditionBody.AppendNewline()
	}

	taskBody.AppendBlock(conditionBlock)

	return taskBlock
}
