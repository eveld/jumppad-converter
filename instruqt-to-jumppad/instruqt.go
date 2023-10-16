package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/jumppad-labs/instruqt-to-jumppad/model"
	"gopkg.in/yaml.v3"
)

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
