package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	App        AppConfig        `yaml:"app" env-prefix:"APP_"`
	Server     ServerConfig     `yaml:"server" env-prefix:"SERVER_"`
	Database   DatabaseConfig   `yaml:"database" env-prefix:"DB_"`
	Cache      CacheConfig      `yaml:"cache" env-prefix:"CACHE_"`
	NetSchool  NetSchoolConfig  `yaml:"netschool" env-prefix:"NETSCHOOL_"`
	JWT        JWTConfig        `yaml:"jwt" env-prefix:"JWT_"`
	Logging    LoggingConfig    `yaml:"logging" env-prefix:"LOGGING_"`
}

type AppConfig struct {
	Name    string `yaml:"name" env:"NAME" env-default:"netschool-proxy/api"`
	Version string `yaml:"version" env:"VERSION" env-default:"1.0.0"`
	Debug   bool   `yaml:"debug" env:"DEBUG" env-default:"false"`
}

type ServerConfig struct {
	Port         string        `yaml:"port" env:"PORT" env-default:"8080"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env:"READ_TIMEOUT" env-default:"10s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env:"WRITE_TIMEOUT" env-default:"10s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env:"IDLE_TIMEOUT" env-default:"30s"`
	CORS         CORSConfig    `yaml:"cors" env-prefix:"CORS_"`
}

type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins" env:"ALLOWED_ORIGINS" env-default:"*"`
	AllowedMethods []string `yaml:"allowed_methods" env:"ALLOWED_METHODS" env-default:"GET,POST,PUT,DELETE,OPTIONS"`
	AllowedHeaders []string `yaml:"allowed_headers" env:"ALLOWED_HEADERS" env-default:"Origin,Content-Type,Accept,Authorization"`
}

type DatabaseConfig struct {
	Type     string `yaml:"type" env:"TYPE" env-default:"postgres"` // "postgres", "mariadb", "mysql", "sqlite"
	Host     string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port     int    `yaml:"port" env:"PORT" env-default:"5432"`
	Name     string `yaml:"name" env:"NAME" env-default:"netschool_proxy"`
	User     string `yaml:"user" env:"USER" env-default:"proxy_user"`
	Password string `yaml:"password" env:"PASSWORD" env-default:"dev_password"`
	SSLMode  string `yaml:"sslmode" env:"SSLMODE" env-default:"disable"`
	URL      string `yaml:"url" env:"URL"` // Alternative: full connection string
	SQLitePath string `yaml:"sqlite_path" env:"SQLITE_PATH" env-default:"./db.sqlite"`
}

type CacheConfig struct {
	Type       string `yaml:"type" env:"TYPE" env-default:"memory"` // "redis" or "memory"
	RedisAddr  string `yaml:"redis_addr" env:"REDIS_ADDR" env-default:"localhost:6379"`
	TTL        int    `yaml:"ttl" env:"TTL" env-default:"300"` // seconds
	MemorySize int    `yaml:"memory_size" env:"MEMORY_SIZE" env-default:"1000"`
}

type NetSchoolConfig struct {
	Mode      string        `yaml:"mode" env:"MODE" env-default:"ns-webapi"` // ns-webapi, ns-mobileapi, dev-mockapi
	Timeout   time.Duration `yaml:"timeout" env:"TIMEOUT" env-default:"30s"`
	RetryMax  int           `yaml:"retry_max" env:"RETRY_MAX" env-default:"3"`
	RetryWait time.Duration `yaml:"retry_wait" env:"RETRY_WAIT" env-default:"1s"`
}

type JWTConfig struct {
	Secret    string        `yaml:"secret" env:"SECRET" env-default:"default_secret_key_change_me"`
	ExpiresIn time.Duration `yaml:"expires_in" env:"EXPIRES_IN" env-default:"24h"`
}

type LoggingConfig struct {
	Level string `yaml:"level" env:"LEVEL" env-default:"info"`
	File  string `yaml:"file" env:"FILE"`
}

// LoadConfig loads configuration from YAML file and environment variables
func LoadConfig(configPath string) (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// If .env file doesn't exist, continue without it
		// This allows the application to work in production environments
	}

	var cfg Config

	if configPath != "" {
		// Load from YAML file
		err := cleanenv.ReadConfig(configPath, &cfg)
		if err != nil {
			return nil, err
		}
	} else {
		// Load only from environment variables
		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			return nil, err
		}
	}

	// Override with environment variables if they exist
	if port := os.Getenv("PORT"); port != "" {
		cfg.Server.Port = port
	}

	return &cfg, nil
}