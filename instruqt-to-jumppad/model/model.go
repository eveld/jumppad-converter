package model

type Track struct {
	Slug        string   `yaml:"slug"`
	ID          string   `yaml:"id"`
	Title       string   `yaml:"title"`
	Teaser      string   `yaml:"teaser"`
	Description string   `yaml:"description"`
	Icon        string   `yaml:"icon"`
	Tags        []string `yaml:"tags"`
	Owner       string   `yaml:"owner"`
	Developers  []string `yaml:"developers"`
	ShowTimer   bool     `yaml:"show_timer"`
	Checksum    string   `yaml:"checksum"`
}

type Challenge struct {
	ID         string
	Slug       string `yaml:"slug"`
	Title      string `yaml:"title"`
	Teaser     string `yaml:"teaser"`
	Assignment string
	Scripts    []Script
	Setups     map[string]string
	Checks     map[string]string
	Solves     map[string]string
	Cleanups   map[string]string
}

type Script struct {
	Type    string
	Target  string
	Content string
}

type Environment struct {
	Version    string      `yaml:"version"`
	Containers []Container `yaml:"containers"`
}

type Container struct {
	Name        string            `yaml:"name"`
	Image       string            `yaml:"image"`
	Entrypoint  string            `yaml:"entrypoint"`
	Command     string            `yaml:"cmd"`
	Shell       string            `yaml:"shell"`
	Ports       []int             `yaml:"ports"`
	Memory      int               `yaml:"memory"`
	Environment map[string]string `yaml:"environment"`
}

type VM struct {
	Name                 string            `yaml:"name"`
	Image                string            `yaml:"image"`
	MachineType          string            `yaml:"machine_type"`
	Environment          map[string]string `yaml:"environment"`
	Shell                string            `yaml:"shell"`
	AllowExternalIngress []string          `yaml:"allow_external_ingress"`
}
