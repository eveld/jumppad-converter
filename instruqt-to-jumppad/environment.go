package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/instruqt/converter/model"
	"github.com/jumppad-labs/hclconfig/convert"
)

func GenerateNetwork() *hclwrite.Block {
	// Network
	subnet, err := convert.GoToCtyValue("10.0.5.0/16")
	if err != nil {
		log.Fatal(err)
	}

	networkBlock := hclwrite.NewBlock("resource", []string{"network", "main"})
	networkBody := networkBlock.Body()
	networkBody.SetAttributeValue("subnet", subnet)

	return networkBlock
}

func GenerateContainer(container *model.Container) *hclwrite.Block {
	block := hclwrite.NewBlock("resource", []string{"container", container.Name})
	body := block.Body()

	image, err := convert.GoToCtyValue(container.Image)
	if err != nil {
		log.Fatal(err)
	}

	imageBlock := body.AppendNewBlock("image", nil)
	imageBody := imageBlock.Body()
	imageBody.SetAttributeValue("name", image)

	body.AppendNewline()
	networkBlock := body.AppendNewBlock("network", nil)
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

		body.AppendNewline()
		body.SetAttributeValue("entrypoint", entrypoint)
	}

	for _, port := range container.Ports {
		body.AppendNewline()
		portBlock := body.AppendNewBlock("port", nil)
		portBody := portBlock.Body()

		local, err := convert.GoToCtyValue(port)
		if err != nil {
			log.Fatal(err)
		}
		portBody.SetAttributeValue("local", local)

		remote, err := convert.GoToCtyValue(port)
		if err != nil {
			log.Fatal(err)
		}
		portBody.SetAttributeValue("remote", remote)

		host, err := convert.GoToCtyValue(port)
		if err != nil {
			log.Fatal(err)
		}
		portBody.SetAttributeValue("host", host)
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

	body.AppendNewline()
	body.SetAttributeValue("environment", environment)

	if container.Memory != 0 {
		body.AppendNewline()
		resourcesBlock := body.AppendNewBlock("resources", nil)
		resourcesBody := resourcesBlock.Body()

		memory, err := convert.GoToCtyValue(container.Memory)
		if err != nil {
			log.Fatal(err)
		}
		resourcesBody.SetAttributeValue("memory", memory)
	}

	return block
}

func GenerateVM(vm *model.VM) *hclwrite.Block {
	block := hclwrite.NewBlock("resource", []string{"vm", vm.Name})
	body := block.Body()

	configBlock := body.AppendNewBlock("config", nil)
	configBody := configBlock.Body()

	arch, err := convert.GoToCtyValue("x86_64")
	if err != nil {
		log.Fatal(err)
	}
	configBody.SetAttributeValue("arch", arch)

	body.AppendNewline()
	image, err := convert.GoToCtyValue(vm.Image)
	if err != nil {
		log.Fatal(err)
	}
	body.SetAttributeValue("image", image)

	body.AppendNewline()
	resourcesBlock := body.AppendNewBlock("resources", nil)
	resourcesBody := resourcesBlock.Body()

	cpu, err := convert.GoToCtyValue(vm.CPU)
	if err != nil {
		log.Fatal(err)
	}
	resourcesBody.SetAttributeValue("cpu", cpu)

	memory, err := convert.GoToCtyValue(vm.Memory)
	if err != nil {
		log.Fatal(err)
	}
	resourcesBody.SetAttributeValue("memory", memory)

	body.AppendNewline()
	networkBlock := body.AppendNewBlock("network", nil)
	networkBody := networkBlock.Body()

	id := hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte(fmt.Sprintf("resource.network.%s.id", "main")),
		},
	}

	networkBody.SetAttributeRaw("id", id)

	// environment
	// append shell variable
	if vm.Environment == nil {
		vm.Environment = map[string]string{}
	}
	vm.Environment["SHELL"] = vm.Shell

	environment, err := convert.GoToCtyValue(vm.Environment)
	if err != nil {
		log.Fatal(err)
	}

	body.AppendNewline()
	body.SetAttributeValue("environment", environment)

	return block
}
