package factory

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sync"
)

// Config - входные данные для фабрики
type Config struct {
	Sources []struct {
		URL string `json:"url"`
		Tag string `json:"tag"`
	} `json:"sources"`
}

// GitClient - пример зависимости
type GitClient interface {
	Clone(url string) error
}

// Реализация GitClient
type gitClientImpl struct {
	cfg *Config
}

func (c *gitClientImpl) Clone(url string) error {
	// Используем данные из конфига
	for _, source := range c.cfg.Sources {
		if source.URL == url {
			// Логика клонирования
			return nil
		}
	}
	return errors.New("source not found")
}

// Singleton Factory
var (
	instance *Factory
	once     sync.Once
)

type Factory struct {
	cfg *Config
}

func ConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "tempro", "config.json"), nil
}

func loadConfig() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	fmt.Printf("cfg path is %s \n", path)

	// Создаем директорию если нужно
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	// Чтение или создание конфига
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		cfg := defaultConfig()
		if err := saveConfig(cfg); err != nil {
			return nil, err
		}
		return cfg, nil
	}
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func saveConfig(cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	configPath, err := ConfigPath()
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func defaultConfig() *Config {
	return &Config{
		Sources: []struct {
			URL string `json:"url"`
			Tag string `json:"tag"`
		}{
			{
				URL: "https://github.com/End1essRage/tempro-templates",
				Tag: "",
			},
		},
	}
}

// Init инициализирует фабрику (должна вызываться один раз при старте приложения)
func Init() error {
	var initErr error
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	once.Do(func() {
		if cfg == nil {
			initErr = errors.New("config cannot be nil")
			return
		}
		instance = &Factory{cfg: cfg}
	})
	return initErr
}

// GetGitClient создает GitClient (глобальный доступ)
func GetGitClient() (GitClient, error) {
	if instance == nil {
		return nil, errors.New("factory not initialized")
	}

	// Валидируем URL из конфига
	for _, source := range instance.cfg.Sources {
		if _, err := url.Parse(source.URL); err != nil {
			return nil, err
		}
	}

	return &gitClientImpl{cfg: instance.cfg}, nil
}
