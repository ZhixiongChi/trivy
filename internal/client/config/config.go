package config

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/genuinetools/reg/registry"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"golang.org/x/xerrors"

	dbTypes "github.com/aquasecurity/trivy-db/pkg/types"
	"github.com/aquasecurity/trivy/pkg/log"
	"github.com/aquasecurity/trivy/pkg/utils"
)

type Config struct {
	context *cli.Context
	logger  *zap.SugaredLogger

	Quiet bool
	Debug bool

	CacheDir   string
	ClearCache bool

	Input    string
	output   string
	Format   string
	Template string

	Timeout       time.Duration
	vulnType      string
	severities    string
	IgnoreFile    string
	IgnoreUnfixed bool
	ExitCode      int

	RemoteAddr    string
	token         string
	tokenHeader   string
	customHeaders []string
	CustomHeaders http.Header

	// these variables are generated by Init()
	ImageName  string
	VulnType   []string
	Output     *os.File
	Severities []dbTypes.Severity
	AppVersion string
}

func New(c *cli.Context) (Config, error) {
	debug := c.Bool("debug")
	quiet := c.Bool("quiet")
	logger, err := log.NewLogger(debug, quiet)
	if err != nil {
		return Config{}, xerrors.New("failed to create a logger")
	}
	return Config{
		context: c,
		logger:  logger,

		Quiet: quiet,
		Debug: debug,

		CacheDir:   c.String("cache-dir"),
		ClearCache: c.Bool("clear-cache"),

		Input:    c.String("input"),
		output:   c.String("output"),
		Format:   c.String("format"),
		Template: c.String("template"),

		Timeout:       c.Duration("timeout"),
		vulnType:      c.String("vuln-type"),
		severities:    c.String("severity"),
		IgnoreFile:    c.String("ignorefile"),
		IgnoreUnfixed: c.Bool("ignore-unfixed"),
		ExitCode:      c.Int("exit-code"),

		RemoteAddr:    c.String("remote"),
		token:         c.String("token"),
		tokenHeader:   c.String("token-header"),
		customHeaders: c.StringSlice("custom-headers"),
	}, nil
}

func (c *Config) Init() (err error) {
	c.Severities = c.splitSeverity(c.severities)
	c.VulnType = strings.Split(c.vulnType, ",")
	c.AppVersion = c.context.App.Version
	c.CustomHeaders = splitCustomHeaders(c.customHeaders)

	// add token to custom headers
	if c.token != "" {
		c.CustomHeaders.Set(c.tokenHeader, c.token)
	}

	if c.Quiet {
		utils.Quiet = true
	}

	// --clear-cache doesn't conduct the scan
	if c.ClearCache {
		return nil
	}

	args := c.context.Args()
	if c.Input == "" && len(args) == 0 {
		c.logger.Error(`trivy requires at least 1 argument or --input option`)
		cli.ShowAppHelp(c.context)
		return xerrors.New("arguments error")
	}

	c.Output = os.Stdout
	if c.output != "" {
		if c.Output, err = os.Create(c.output); err != nil {
			return xerrors.Errorf("failed to create an output file: %w", err)
		}
	}

	if c.Input == "" {
		c.ImageName = args[0]
	}

	// Check whether 'latest' tag is used
	if c.ImageName != "" {
		image, err := registry.ParseImage(c.ImageName)
		if err != nil {
			return xerrors.Errorf("invalid image: %w", err)
		}
		if image.Tag == "latest" {
			c.logger.Warn("You should avoid using the :latest tag as it is cached. You need to specify '--clear-cache' option when :latest image is changed")
		}
	}

	return nil
}

func (c *Config) splitSeverity(severity string) []dbTypes.Severity {
	c.logger.Debugf("Severities: %s", severity)
	var severities []dbTypes.Severity
	for _, s := range strings.Split(severity, ",") {
		severity, err := dbTypes.NewSeverity(s)
		if err != nil {
			c.logger.Warnf("unknown severity option: %s", err)
		}
		severities = append(severities, severity)
	}
	return severities
}

func splitCustomHeaders(headers []string) http.Header {
	result := make(http.Header)
	for _, header := range headers {
		// e.g. x-api-token:XXX
		s := strings.SplitN(header, ":", 2)
		if len(s) != 2 {
			continue
		}
		result.Set(s[0], s[1])
	}
	return result
}
