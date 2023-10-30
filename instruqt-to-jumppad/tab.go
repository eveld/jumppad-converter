package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/instruqt/converter/model"
	"github.com/jumppad-labs/hclconfig/convert"
)

func GenerateTabSlug(tab model.Tab) string {
	var slug string
	switch tab.Type {
	case "service":
		slug = fmt.Sprintf("%s_%s_%d", tab.Type, tab.Hostname, tab.Port)
	case "terminal":
		slug = fmt.Sprintf("%s_%s", tab.Type, tab.Hostname)
	case "code":
		slug = fmt.Sprintf("%s_%s", tab.Type, tab.Hostname)
	}
	return slug
}

func GenerateTab(tab *model.Tab, slug string) *hclwrite.Block {
	/*
		Tab naming:
		service: resource.tab.type_hostname_port
		terminal: resource.tab.type_hostname
		code: resource.tab.type_hostname
	*/

	block := hclwrite.NewBlock("resource", []string{"tab", slug})
	body := block.Body()

	tabType, err := convert.GoToCtyValue(tab.Type)
	if err != nil {
		log.Fatal(err)
	}
	body.SetAttributeValue("type", tabType)

	title, err := convert.GoToCtyValue(tab.Title)
	if err != nil {
		log.Fatal(err)
	}
	body.SetAttributeValue("title", title)

	if tab.Hostname != "" {
		hostname, err := convert.GoToCtyValue(tab.Hostname)
		if err != nil {
			log.Fatal(err)
		}
		body.SetAttributeValue("hostname", hostname)
	}

	if tab.Path != "" {
		path, err := convert.GoToCtyValue(tab.Path)
		if err != nil {
			log.Fatal(err)
		}
		body.SetAttributeValue("path", path)
	}

	if tab.Port != 0 {
		port, err := convert.GoToCtyValue(tab.Port)
		if err != nil {
			log.Fatal(err)
		}
		body.SetAttributeValue("port", port)
	}

	if tab.CustomRequestHeaders != nil {
		customRequestHeaders, err := convert.GoToCtyValue(tab.CustomRequestHeaders)
		if err != nil {
			log.Fatal(err)
		}
		body.SetAttributeValue("custom_request_headers", customRequestHeaders)
	}

	if tab.CustomResponseHeaders != nil {
		customResponseHeaders, err := convert.GoToCtyValue(tab.CustomResponseHeaders)
		if err != nil {
			log.Fatal(err)
		}
		body.SetAttributeValue("custom_response_headers", customResponseHeaders)
	}

	if tab.Workdir != "" {
		workdir, err := convert.GoToCtyValue(tab.Workdir)
		if err != nil {
			log.Fatal(err)
		}
		body.SetAttributeValue("workdir", workdir)
	}

	if tab.Command != "" {
		command, err := convert.GoToCtyValue(tab.Command)
		if err != nil {
			log.Fatal(err)
		}
		body.SetAttributeValue("command", command)
	}

	if tab.NewWindow {
		newWindow, err := convert.GoToCtyValue(tab.NewWindow)
		if err != nil {
			log.Fatal(err)
		}
		body.SetAttributeValue("new_window", newWindow)
	}

	return block
}
