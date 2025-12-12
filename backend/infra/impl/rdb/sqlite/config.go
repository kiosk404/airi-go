package sqlite

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/kiosk404/airi-go/backend/pkg/lang/conv"
)

type Config struct {
	Timeout       time.Duration `yaml:"timeout"`
	Loc           string        `yaml:"loc"`
	DBName        string        `yaml:"db_name"`
	WithReturning bool          `yaml:"with_returning"`
}

func (cfg *Config) buildDSN() string {
	dsn := cfg.DBName

	args := []string{
		"cache=shared",
		"_pragma=foreign_keys(1)",
		"_busy_timeout=" + conv.Int64ToStr(cfg.Timeout.Milliseconds()),
	}

	if cfg.Loc != "" {
		args = append(args, "_loc="+url.QueryEscape(cfg.Loc))
	}

	return fmt.Sprintf("file:%s?%s", dsn, strings.Join(args, "&"))
}
