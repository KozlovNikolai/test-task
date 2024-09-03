package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config is ...
type Config struct {
	StorageWR         string        `yaml:"storage_wr" env-required:"true"`
	StorageRO         string        `yaml:"storage_ro" env-required:"false"`
	RepoType          string        `yaml:"repo_type" env-required:"true"`
	TokenTimeDuration time.Duration `yaml:"token_time_duration" env-required:"true"`
	HTTPServer        `yaml:"http_server"`
}

// HTTPServer is ...
type HTTPServer struct {
	Address    string        `yaml:"address" env-default:"localhost:8080"`
	Timeout    time.Duration `yaml:"timeout" env-default:"14s"`
	IdleTimout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

var (
	// Cfg is ...
	Cfg = &Config{}
	// JwtKey is ...
	JwtKey = []byte(os.Getenv("JWT_KEY"))
	// CertFile is ...
	CertFile = os.Getenv("TLS_CERT")
	// KeyFile is ...
	KeyFile   = os.Getenv("TLS_KEY")
	AdminRole = "admin"
)

// MustLoad is ...
func MustLoad() {
	configPath := "./configs/config.yaml"

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}
	//var cfg Config
	if err := cleanenv.ReadConfig(configPath, Cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	//return cfg

}
