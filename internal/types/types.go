package types

import "gopkg.in/yaml.v2"

type (
	// Cluster instance
	Cluster struct {
		Cluster map[string]string `yaml:"cluster"`
		Name    string            `yaml:"name"`
	}

	// User instance
	User struct {
		Name string `yaml:"name"`
	}

	// Context instance
	Context struct {
		Context map[string]string `yaml:"context"`
		Name    string            `yaml:"name"`
	}

	// KubeConfig instance
	KubeConfig struct {
		CurrentContext string    `yaml:"current-context"`
		Clusters       []Cluster `yaml:"clusters"`
		Users          []User    `yaml:"users"`
		Contexts       []Context `yaml:"contexts"`
	}
)

// Parse data into KubeConfig instance.
func (k *KubeConfig) Parse(data []byte) error {
	err := yaml.Unmarshal(data, k)
	if err != nil {
		return err
	}
	return nil
}
