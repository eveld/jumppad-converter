package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/compose-spec/compose-go/cli"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/jumppad-labs/hclconfig/convert"
)

func main() {
	args := os.Args[1:]

	file, err := filepath.Abs(args[0])
	if err != nil {
		log.Fatal(err)
	}

	dirname := filepath.Base(filepath.Dir(file))

	projectOptions, err := cli.NewProjectOptions([]string{file}, cli.WithOsEnv, cli.WithWorkingDirectory(filepath.Dir(file)), cli.WithInterpolation(true))
	if err != nil {
		log.Fatal(err)
	}

	project, err := cli.ProjectFromOptions(projectOptions)
	if err != nil {
		log.Fatal(err)
	}

	f := hclwrite.NewEmptyFile()
	root := f.Body()

	// networks
	for _, network := range project.Networks {
		name := strings.TrimPrefix(network.Name, fmt.Sprintf("%s_", dirname))

		if name == "default" {
			continue
		}

		cidr := "10.0.5.0/16"
		if len(network.Ipam.Config) != 0 {
			cidr = network.Ipam.Config[0].Subnet

		}

		subnet, err := convert.GoToCtyValue(cidr)
		if err != nil {
			log.Fatal(err)
		}

		block := root.AppendNewBlock("resource", []string{"network", name})
		body := block.Body()
		body.SetAttributeValue("subnet", subnet)

		root.AppendNewline()
	}

	// services
	for _, service := range project.Services {
		block := root.AppendNewBlock("resource", []string{"container", service.Name})
		body := block.Body()

		if len(service.DependsOn) > 0 {
			names := []string{}
			for _, name := range service.GetDependencies() {
				names = append(names, fmt.Sprintf("resource.container.%s", name))
			}

			dependencies, err := convert.GoToCtyValue(names)
			if err != nil {
				log.Fatal(err)
			}

			body.SetAttributeValue("depends_on", dependencies)
			body.AppendNewline()
		}

		if service.Image != "" {
			/*
				image {
					name     = "consul:1.6.1"
					username = "repo_username"
					password = "repo_password"
				}
			*/
			image, err := convert.GoToCtyValue(service.Image)
			if err != nil {
				log.Fatal(err)
			}

			imageBlock := body.AppendNewBlock("image", nil)
			imageBody := imageBlock.Body()
			imageBody.SetAttributeValue("name", image)
		}

		if len(service.Networks) > 0 {
			/*
				network {
					id         = resource.network.cloud.id
					ip_address = "10.16.0.200"
					aliases    = ["my_unique_name_ip_address"]
				}
			*/
			for name := range service.Networks {
				if name == "default" {
					continue
				}

				body.AppendNewline()
				networkBlock := body.AppendNewBlock("network", nil)
				networkBody := networkBlock.Body()

				id := hclwrite.Tokens{
					{
						Type:  hclsyntax.TokenIdent,
						Bytes: []byte(fmt.Sprintf("resource.network.%s.id", name)),
					},
				}

				networkBody.SetAttributeRaw("id", id)

				aliases, err := convert.GoToCtyValue([]string{service.ContainerName})
				if err != nil {
					log.Fatal(err)
				}

				networkBody.SetAttributeValue("aliases", aliases)
			}
		}

		if len(service.Command) > 0 {
			/*
				command = [
					"consul",
					"agent"
				]
			*/
			command, err := convert.GoToCtyValue(service.Command)
			if err != nil {
				log.Fatal(err)
			}

			body.AppendNewline()
			body.SetAttributeValue("command", command)
		}

		if len(service.Ports) > 0 {
			/*
				port {
					local  = 8500
					remote = 8500
					host   = 18500
				}
			*/

			for _, port := range service.Ports {
				body.AppendNewline()
				portBlock := body.AppendNewBlock("port", nil)
				portBody := portBlock.Body()

				if port.Protocol != "tcp" {
					protocol, err := convert.GoToCtyValue(port.Protocol)
					if err != nil {
						log.Fatal(err)
					}

					portBody.SetAttributeValue("protocol", protocol)
					portBody.AppendNewline()
				}

				local, err := convert.GoToCtyValue(int(port.Target))
				if err != nil {
					log.Fatal(err)
				}

				portBody.SetAttributeValue("local", local)

				published, _ := strconv.Atoi(port.Published)
				remote, err := convert.GoToCtyValue(published)
				if err != nil {
					log.Fatal(err)
				}

				portBody.SetAttributeValue("remote", remote)
				portBody.SetAttributeValue("host", remote)
			}
		}

		if len(service.Expose) > 0 {
			for _, expose := range service.Expose {
				number, _ := strconv.Atoi(expose)
				port, err := convert.GoToCtyValue(number)
				if err != nil {
					log.Fatal(err)
				}

				body.AppendNewline()
				portBlock := body.AppendNewBlock("port", nil)
				portBody := portBlock.Body()
				portBody.SetAttributeValue("local", port)
				portBody.SetAttributeValue("remote", port)
			}
		}

		if len(service.Environment) > 0 {
			/*
				environment = {
					CONSUL_HTTP_ADDR = "http://localhost:8500"
				}
			*/

			//
			// Can we safely just convert map[string]*string?
			// It seems to work...
			//
			environment, err := convert.GoToCtyValue(service.Environment.RemoveEmpty())
			if err != nil {
				log.Fatal(err)
			}

			body.AppendNewline()
			body.SetAttributeValue("environment", environment)
		}

		if len(service.Volumes) > 0 {
			/*
			 	volume {
			    source      = data("config")
			    destination = "/config"
			  }
			*/
			for _, volume := range service.Volumes {
				body.AppendNewline()
				volumeBlock := body.AppendNewBlock("volume", nil)
				volumeBody := volumeBlock.Body()

				if strings.Index(volume.Source, "/") != 0 {
					source := hclwrite.Tokens{
						{
							Type:  hclsyntax.TokenIdent,
							Bytes: []byte(fmt.Sprintf("data(\"%s\")", volume.Source)),
						},
					}

					volumeBody.SetAttributeRaw("source", source)
				} else {
					source, err := convert.GoToCtyValue(volume.Source)
					if err != nil {
						log.Fatal(err)
					}

					volumeBody.SetAttributeValue("source", source)
				}

				destination, err := convert.GoToCtyValue(volume.Target)
				if err != nil {
					log.Fatal(err)
				}
				volumeBody.SetAttributeValue("destination", destination)
			}
		}

		if service.HealthCheck != nil {
			/*
				health_check {
					timeout = "30s"

					exec {
						script = <<-EOF
							#!/bin/bash

							curl "http://localhost:9090"
						EOF
					}
				}
			*/

			// Apparently the timeout is not in seconds but "10s"
			duration, err := time.ParseDuration(service.HealthCheck.Timeout.String())
			if err != nil {
				log.Fatal(err)
			}

			interval, err := time.ParseDuration(service.HealthCheck.Interval.String())
			if err != nil {
				log.Fatal(err)
			}

			seconds := int((duration.Seconds() + interval.Seconds()) * float64(*service.HealthCheck.Retries))

			timeout, err := convert.GoToCtyValue(fmt.Sprintf("%ds", seconds))
			if err != nil {
				log.Fatal(err)
			}

			// Cut off CMD / CMD-SHELL / NONE
			args := service.HealthCheck.Test[1:]
			script, err := convert.GoToCtyValue(strings.Join(args, " "))
			if err != nil {
				log.Fatal(err)
			}

			body.AppendNewline()
			healthcheckBlock := body.AppendNewBlock("health_check", nil)
			healthcheckBody := healthcheckBlock.Body()
			healthcheckBody.SetAttributeValue("timeout", timeout)
			healthcheckBody.AppendNewline()
			execBlock := healthcheckBody.AppendNewBlock("exec", nil)
			execBody := execBlock.Body()
			execBody.SetAttributeValue("script", script)
		}

		if len(service.CapAdd) > 0 || service.Privileged {
			privileged, err := convert.GoToCtyValue(true)
			if err != nil {
				log.Fatal(err)
			}

			body.AppendNewline()
			body.SetAttributeValue("privileged", privileged)
		}

		root.AppendNewline()
	}

	err = os.WriteFile("out/main.hcl", f.Bytes(), 0755)
	if err != nil {
		log.Fatal(err)
	}
}
