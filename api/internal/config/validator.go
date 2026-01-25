package config

import (
	"errors"
	"fmt"
	"strings"
)


func (cfg *Config) Validate() error {
	var errs []string

	
	if cfg.Server.Port == "" {
		errs = append(errs, "server.port is required")
	}

	
	if cfg.Database.URL == "" {
		if cfg.Database.Host == "" {
			errs = append(errs, "database.host is required")
		}
		if cfg.Database.Name == "" {
			errs = append(errs, "database.name is required")
		}
		if cfg.Database.User == "" {
			errs = append(errs, "database.user is required")
		}
		if cfg.Database.Password == "" {
			errs = append(errs, "database.password is required")
		}
	}

	

	
	if cfg.JWT.Secret == "" || cfg.JWT.Secret == "default_secret_key_change_me" {
		errs = append(errs, "jwt.secret is required and should not be default value")
	}

	if len(errs) > 0 {
		return errors.New("validation errors: " + strings.Join(errs, ", "))
	}

	return nil
}


func (cfg *Config) GetDatabaseURL() string {
	if cfg.Database.URL != "" {
		return cfg.Database.URL
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)
}