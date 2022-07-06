package plugin

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/xerrors"

	"github.com/aquasecurity/trivy/pkg/log"
	"github.com/aquasecurity/trivy/pkg/plugin"
)

// Install installs a plugin
func Install(cmd *cobra.Command, args []string) error {
	url := args[0]
	if _, err := plugin.Install(cmd.Context(), url, true); err != nil {
		return xerrors.Errorf("plugin install error: %w", err)
	}

	return nil
}

// Uninstall uninstalls the plugin
func Uninstall(_ *cobra.Command, args []string) error {
	if err := plugin.Uninstall(args[0]); err != nil {
		return xerrors.Errorf("plugin uninstall error: %w", err)
	}
	return nil
}

// Information displays information about the plugin
func Information(_ *cobra.Command, args []string) error {
	info, err := plugin.Information(args[0])
	if err != nil {
		return xerrors.Errorf("plugin information display error: %w", err)
	}

	if _, err = fmt.Fprintf(os.Stdout, info); err != nil {
		return xerrors.Errorf("print error: %w", err)
	}

	return nil
}

// List displays a list of all of installed plugins
func List(_ *cobra.Command, _ []string) error {
	info, err := plugin.List()
	if err != nil {
		return xerrors.Errorf("plugin list display error: %w", err)
	}

	if _, err = fmt.Fprintf(os.Stdout, info); err != nil {
		return xerrors.Errorf("print error: %w", err)
	}

	return nil
}

// Update updates an existing plugin
func Update(_ *cobra.Command, args []string) error {
	if err := plugin.Update(args[0]); err != nil {
		return xerrors.Errorf("plugin update error: %w", err)
	}
	return nil
}

// Run runs the plugin
func Run(cmd *cobra.Command, args []string) error {
	url := args[0]
	pluginArgs := args[1:]
	return RunWithArgs(cmd.Context(), url, pluginArgs)
}

// RunWithArgs runs the plugin with arguments
func RunWithArgs(ctx context.Context, url string, args []string) error {
	pl, err := plugin.Install(ctx, url, false)
	if err != nil {
		return xerrors.Errorf("plugin install error: %w", err)
	}

	if err = pl.Run(ctx, args); err != nil {
		return xerrors.Errorf("unable to run %s plugin: %w", pl.Name, err)
	}
	return nil
}

// LoadCommands loads plugins as subcommands
func LoadCommands() []*cobra.Command {
	var commands []*cobra.Command
	plugins, err := plugin.LoadAll()
	if err != nil {
		log.Logger.Debugf("no plugins were loaded")
		return nil
	}
	for _, p := range plugins {
		p := p
		cmd := &cobra.Command{
			Use:   fmt.Sprintf("%s [flags]", p.Name),
			Short: p.Usage,
			RunE: func(cmd *cobra.Command, args []string) error {
				if err = p.Run(cmd.Context(), args); err != nil {
					return xerrors.Errorf("plugin error: %w", err)
				}
				return nil
			},
		}
		commands = append(commands, cmd)
	}
	return commands
}
