package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/instruqt/converter/model"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Track       model.Track
	Challenges  []model.Challenge
	Environment model.Environment
}

func main() {
	args := os.Args[1:]

	dir, err := filepath.Abs(args[0])
	if err != nil {
		log.Fatal(err)
	}

	track, err := readTrack(dir)
	if err != nil {
		log.Fatal(err)
	}

	challenges, err := readChallenges(dir)
	if err != nil {
		log.Fatal(err)
	}

	env, err := readConfig(dir)
	if err != nil {
		log.Fatal(err)
	}

	config := &Config{
		Track:       *track,
		Challenges:  challenges,
		Environment: *env,
	}

	err = os.RemoveAll("out")
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll("out", 0755)
	if err != nil {
		log.Fatal(err)
	}

	f := hclwrite.NewEmptyFile()
	root := f.Body()

	trackBlock := config.GenerateTrack(track, challenges)
	root.AppendBlock(trackBlock)
	root.AppendNewline()

	allTabs := map[string]model.Tab{}
	for _, challenge := range challenges {
		challengeTabs := map[string]model.Tab{}
		for _, tab := range challenge.Tabs {
			slug := config.GenerateTabSlug(tab)
			challengeTabs[slug] = tab
			allTabs[slug] = tab
		}

		if challenge.Type == "quiz" {
			quizBlock := config.GenerateQuiz(&challenge, challengeTabs)
			root.AppendBlock(quizBlock)
			root.AppendNewline()
		}

		if challenge.Type == "challenge" {
			challengeBlock := config.GenerateChallenge(&challenge, challengeTabs)
			root.AppendBlock(challengeBlock)
			root.AppendNewline()
		}
	}

	for slug, tab := range allTabs {
		tabBlock := config.GenerateTab(&tab, slug)
		root.AppendBlock(tabBlock)
		root.AppendNewline()
	}

	err = os.WriteFile("out/track.hcl", f.Bytes(), 0755)
	if err != nil {
		log.Fatal(err)
	}

	e := hclwrite.NewEmptyFile()
	environment := e.Body()

	networkBlock := GenerateNetwork()
	environment.AppendBlock(networkBlock)
	environment.AppendNewline()

	for _, container := range env.Containers {
		containerBlock := GenerateContainer(&container)
		environment.AppendBlock(containerBlock)
		environment.AppendNewline()
	}

	for _, vm := range env.VirtualMachines {
		vmBlock := GenerateVM(&vm)
		environment.AppendBlock(vmBlock)
		environment.AppendNewline()
	}

	err = os.WriteFile("out/config.hcl", e.Bytes(), 0755)
	if err != nil {
		log.Fatal(err)
	}
}

func readTrack(path string) (*model.Track, error) {
	trackFile := filepath.Join(path, "track.yml")
	trackData, err := os.ReadFile(trackFile)
	if err != nil {
		return nil, err
	}

	track := &model.Track{}
	err = yaml.Unmarshal(trackData, track)
	if err != nil {
		return nil, err
	}

	track.Slug = strings.ReplaceAll(track.Slug, "-", "_")

	return track, nil
}

func readChallenges(path string) ([]model.Challenge, error) {
	challenges := []model.Challenge{}

	entries, err := os.ReadDir(path)
	if err != nil {
		return challenges, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			r := regexp.MustCompile(`^\d.*`)
			if r.MatchString(entry.Name()) {
				challengePath := filepath.Join(path, entry.Name())
				files, err := os.ReadDir(challengePath)
				if err != nil {
					return challenges, err
				}

				parts := strings.SplitN(entry.Name(), "-", 2)

				assignmentPath := filepath.Join(challengePath, "assignment.md")
				data, err := os.ReadFile(assignmentPath)
				if err != nil {
					return challenges, err
				}

				challenge := model.Challenge{
					ID:       parts[0],
					Scripts:  []model.Script{},
					Setups:   map[string]string{},
					Cleanups: map[string]string{},
					Checks:   map[string]string{},
					Solves:   map[string]string{},
				}

				rest, err := frontmatter.Parse(strings.NewReader(string(data)), &challenge)
				if err != nil {
					return challenges, err
				}

				challenge.Slug = strings.ReplaceAll(challenge.Slug, "-", "_")
				challenge.Assignment = string(rest)

				for _, file := range files {
					if !file.IsDir() && file.Name() != "assignment.md" {
						filePath := filepath.Join(challengePath, file.Name())
						data, err := os.ReadFile(filePath)
						if err != nil {
							return challenges, err
						}

						nameParts := strings.Split(file.Name(), "-")

						script := model.Script{
							Type:    nameParts[0],
							Target:  strings.Join(nameParts[1:], "-"),
							Content: string(data),
						}

						challenge.Scripts = append(challenge.Scripts, script)
					}
				}

				challenges = append(challenges, challenge)
			}
		}
	}

	return challenges, nil
}

func readConfig(path string) (*model.Environment, error) {
	envFile := filepath.Join(path, "config.yml")
	envData, err := os.ReadFile(envFile)
	if err != nil {
		return nil, err
	}

	env := &model.Environment{}
	err = yaml.Unmarshal(envData, env)
	if err != nil {
		return nil, err
	}

	return env, nil
}

func (c *Config) GenerateComment(comment string) hclwrite.Tokens {
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

func (c *Config) LookupResource(hostname string) (string, error) {
	for _, vm := range c.Environment.VirtualMachines {
		if vm.Name == hostname {
			return fmt.Sprintf("resource.vm.%s", vm.Name), nil
		}
	}

	for _, container := range c.Environment.Containers {
		if container.Name == hostname {
			return fmt.Sprintf("resource.container.%s", container.Name), nil
		}
	}

	return "", fmt.Errorf("could not find resource for hostname %s", hostname)
}
