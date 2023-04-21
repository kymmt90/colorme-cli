package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type UserConfig struct {
	LoginID     string `yaml:"login_id"`
	AccessToken string `yaml:"access_token"`
}

func (c *UserConfig) Save() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("saveUserConfig: %w", err)
	}

	colormeCfgDir := filepath.Join(home, ".config/colorme")
	if err := os.MkdirAll(colormeCfgDir, 0751); err != nil {
		return fmt.Errorf("saveUserConfig: %w", err)
	}

	d, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("saveUserConfig: %w", err)
	}

	f, err := os.Create(filepath.Join(colormeCfgDir, "users.yml"))
	if err != nil {
		return fmt.Errorf("saveUserConfig: %w", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	n, err := w.Write(d)
	if err != nil {
		return fmt.Errorf("saveUserConfig: %w", err)
	}
	if n != len(d) {
		return fmt.Errorf("saveUserConfig: failed to write all bytes")
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("saveUserConfig: %w", err)
	}

	return nil
}
