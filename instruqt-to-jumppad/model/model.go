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
	Type       string `yaml:"type"`
	Title      string `yaml:"title"`
	Teaser     string `yaml:"teaser"`
	Notes      []Note `yaml:"notes"`
	Tabs       []Tab  `yaml:"tabs"`
	Assignment string
	Scripts    []Script
	Setups     map[string]string
	Checks     map[string]string
	Solves     map[string]string
	Cleanups   map[string]string
	Answers    []string `yaml:"answers"`
	Solutions  []int    `yaml:"solution"`
}

type Note struct {
	Type     string `yaml:"type"`
	Contents string `yaml:"contents"`
	URL      string `yaml:"url"`
}

type Tab struct {
	Slug                  string
	Type                  string            `yaml:"type"`
	Title                 string            `yaml:"title"`
	URL                   string            `yaml:"url"`
	Hostname              string            `yaml:"hostname"`
	Path                  string            `yaml:"path"`
	Port                  int               `yaml:"port"`
	CustomRequestHeaders  map[string]string `yaml:"custom_request_headers"`
	CustomResponseHeaders map[string]string `yaml:"custom_response_headers"`
	Workdir               string            `yaml:"workdir"`
	Command               string            `yaml:"command"`
	NewWindow             bool              `yaml:"new_window"`
}

type Script struct {
	Type    string
	Target  string
	Content string
}

type Environment struct {
	Version         string      `yaml:"version"`
	Containers      []Container `yaml:"containers"`
	VirtualMachines []VM        `yaml:"virtualmachines"`
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
	Memory               int               `yaml:"memory"`
	CPU                  int               `yaml:"cpus"`
}
