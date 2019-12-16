package config

import (
	"github.com/urfave/cli"
	"golang.org/x/xerrors"

	"github.com/aquasecurity/trivy/pkg/utils"
)

type Config struct {
	context *cli.Context

	Quiet          bool
	Debug          bool
	CacheDir       string
	Reset          bool
	DownloadDBOnly bool
	SkipUpdate     bool

	Listen      string
	Token       string
	TokenHeader string

	// these variables are generated by Init()
	AppVersion string
}

func New(c *cli.Context) Config {
	debug := c.Bool("debug")
	quiet := c.Bool("quiet")
	return Config{
		context: c,

		Quiet:          quiet,
		Debug:          debug,
		CacheDir:       c.String("cache-dir"),
		Reset:          c.Bool("reset"),
		DownloadDBOnly: c.Bool("download-db-only"),
		SkipUpdate:     c.Bool("skip-update"),
		Listen:         c.String("listen"),
		Token:          c.String("token"),
		TokenHeader:    c.String("token-header"),
	}
}

func (c *Config) Init() (err error) {
	if c.SkipUpdate && c.DownloadDBOnly {
		return xerrors.New("The --skip-update and --download-db-only option can not be specified both")
	}

	c.AppVersion = c.context.App.Version

	// A server always suppresses a progress bar
	utils.Quiet = true

	return nil
}
