package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ServerAddress  string        `yaml:"server_address"`
	Username       string        `yaml:"username"`
	PrivateKeyPath string        `yaml:"private_key_path"`
	RemoteRoot     string        `yaml:"remote_root"`
	Games          []*GameConfig `yaml:"games"`

	WatchTargetCount int `yaml:"-"`
}

type GameConfig struct {
	Name         string   `yaml:"name"`
	LocalDir     string   `yaml:"local_dir"`
	FilePatterns []string `yaml:"file_patterns"`
	ProgramName  string   `yaml:"program_name"`

	RemoteRoot string           `yaml:"-"`
	AltName    string           `yaml:"-"`
	FileRegExp []*regexp.Regexp `yaml:"-"`
}

var (
	Value *Config

	ReAltName = regexp.MustCompile(`[^a-zA-Z0-9_]`)
)

func LoadConfig() error {
	bt, err := os.ReadFile("./config.yaml")
	if err != nil {
		if os.IsNotExist(err) {
			log.Println(`config.yaml file not found.`)
			log.Println(`Run "scpsave.exe -c" to create a "config.sample.yaml" file.`)
			log.Println(`After editing it, rename it to "config.yaml" and run the program again.`)
		}
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(bt, &config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}
	if len(config.Games) == 0 {
		return errors.New("empty games config")
	}

	testAltNames := make(map[string]struct{})
	for _, game := range config.Games {
		game.RemoteRoot = config.RemoteRoot
		game.ProgramName = strings.ToLower(game.ProgramName)
		if game.ProgramName != "" {
			config.WatchTargetCount++
		}

		game.AltName = ReAltName.ReplaceAllString(game.Name, "")
		if game.AltName == "" {
			return fmt.Errorf("game '%s' has an invalid name that results in an empty AltName", game.Name)
		}
		if _, exists := testAltNames[game.AltName]; exists {
			return fmt.Errorf("game '%s' has a duplicate AltName '%s'", game.Name, game.AltName)
		}
		testAltNames[game.AltName] = struct{}{}

		for _, pattern := range game.FilePatterns {
			re, err := regexp.Compile(strings.ToLower(pattern))
			if err != nil {
				return fmt.Errorf("failed to compile regex pattern '%s' for game '%s': %w", pattern, game.Name, err)
			}
			game.FileRegExp = append(game.FileRegExp, re)
		}
	}
	Value = &config
	return nil
}

func MakeSampleConfig() error {
	config := Config{
		ServerAddress:  "example.com:22",
		Username:       "user",
		PrivateKeyPath: `C:\Users\user\.ssh\id_rsa`,
		RemoteRoot:     "/remote/path",
		Games: []*GameConfig{
			{
				Name:         "Game1",
				LocalDir:     `C:\Users\user\Games\Game1`,
				FilePatterns: []string{`some.+\\.+\.save`, `.+\.dat`},
				ProgramName:  "game1.exe",
			},
			{
				Name:         "Game2",
				LocalDir:     `C:\Users\user\Games\Game2`,
				FilePatterns: []string{`.+\.sav`, `.+\.dat`},
				ProgramName:  `C:\Game\Folder\game2.exe`,
			},
		},
	}
	bt, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("failed to marshal sample config: %w", err)
	}
	if err := os.WriteFile("./config.sample.yaml", bt, 0644); err != nil {
		return fmt.Errorf("failed to write sample config file: %w", err)
	}
	return nil
}

func (g *GameConfig) BaseMetaFilePath() string {
	return filepath.Join(".", "working", g.AltName, "base.yaml")
}

func (g *GameConfig) RemoteMetaFileLocalPath() string {
	return filepath.Join(".", "working", g.AltName, "remote.yaml")
}

func (g *GameConfig) LocalFilePath(relPath string) string {
	return filepath.Join(g.LocalDir, relPath)
}

func (g *GameConfig) LocalFileDownloadPath(relPath string) string {
	return filepath.Join(".", "working", g.AltName, relPath)
}

func (g *GameConfig) RemoteMetaFileRemotePath() string {
	return path.Join(g.RemoteRoot, "meta", g.AltName, "remote.yaml")
}

func (g *GameConfig) RemoteMetaFileUploadPath() string {
	return path.Join(g.RemoteRoot, "meta", g.AltName, "remote.temp.yaml")
}

func (g *GameConfig) RemoteFileUploadPath(relPath string) string {
	return path.Join(g.RemoteRoot, "upload", g.AltName, relPath)
}

func (g *GameConfig) RemoteFilePath(relPath string) string {
	return path.Join(g.RemoteRoot, "save", g.AltName, relPath)
}
