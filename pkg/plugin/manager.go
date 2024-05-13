package plugin

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/samber/lo"
	"golang.org/x/xerrors"
	"gopkg.in/yaml.v3"

	"github.com/aquasecurity/trivy/pkg/downloader"
	ftypes "github.com/aquasecurity/trivy/pkg/fanal/types"
	"github.com/aquasecurity/trivy/pkg/log"
	"github.com/aquasecurity/trivy/pkg/utils/fsutils"
)

const configFile = "plugin.yaml"

var (
	pluginsRelativeDir = filepath.Join(".trivy", "plugins")

	_defaultManager *Manager
)

type ManagerOption func(indexer *Manager)

func WithWriter(w io.Writer) ManagerOption {
	return func(indexer *Manager) {
		indexer.w = w
	}
}

func WithIndexURL(indexURL string) ManagerOption {
	return func(indexer *Manager) {
		indexer.indexURL = indexURL
	}
}

// Manager manages the plugins
type Manager struct {
	w          io.Writer
	indexURL   string
	logger     *log.Logger
	pluginRoot string
	indexPath  string
}

func NewManager(opts ...ManagerOption) *Manager {
	m := &Manager{
		w:          os.Stdout,
		indexURL:   indexURL,
		logger:     log.WithPrefix("plugin"),
		pluginRoot: filepath.Join(fsutils.HomeDir(), pluginsRelativeDir),
		indexPath:  filepath.Join(fsutils.CacheDir(), "plugin", "index.yaml"),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func defaultManager() *Manager {
	if _defaultManager == nil {
		_defaultManager = NewManager()
	}
	return _defaultManager
}

func Install(ctx context.Context, name string, opts Options) (Plugin, error) {
	return defaultManager().Install(ctx, name, opts)
}
func Start(ctx context.Context, name string, opts Options) (Wait, error) {
	return defaultManager().Start(ctx, name, opts)
}
func RunWithURL(ctx context.Context, name string, opts Options) error {
	return defaultManager().RunWithURL(ctx, name, opts)
}
func Upgrade(ctx context.Context, names []string) error { return defaultManager().Upgrade(ctx, names) }
func Uninstall(ctx context.Context, name string) error  { return defaultManager().Uninstall(ctx, name) }
func Information(name string) error                     { return defaultManager().Information(name) }
func List(ctx context.Context) error                    { return defaultManager().List(ctx) }
func Update(ctx context.Context) error                  { return defaultManager().Update(ctx) }
func Search(ctx context.Context, args []string) error   { return defaultManager().Search(ctx, args) }

// Install installs a plugin
func (m *Manager) Install(ctx context.Context, name string, opts Options) (Plugin, error) {
	src := m.tryIndex(ctx, name)

	// If the plugin is already installed, it skips installing the plugin.
	if p, installed := m.isInstalled(ctx, src); installed {
		m.logger.InfoContext(ctx, "The plugin is already installed", log.String("name", p.Name))
		return p, nil
	}

	m.logger.InfoContext(ctx, "Installing the plugin...", log.String("src", src))
	return m.install(ctx, src, opts)
}

func (m *Manager) install(ctx context.Context, src string, opts Options) (Plugin, error) {
	tempDir, err := downloader.DownloadToTempDir(ctx, src)
	if err != nil {
		return Plugin{}, xerrors.Errorf("download failed: %w", err)
	}
	defer os.RemoveAll(tempDir)

	m.logger.DebugContext(ctx, "Loading the plugin metadata...")
	plugin, err := m.loadMetadata(tempDir)
	if err != nil {
		return Plugin{}, xerrors.Errorf("failed to load the plugin metadata: %w", err)
	}

	if err = plugin.install(ctx, plugin.Dir(), tempDir, opts); err != nil {
		return Plugin{}, xerrors.Errorf("failed to install the plugin: %w", err)
	}

	// Copy plugin.yaml into the plugin dir
	f, err := os.Create(filepath.Join(plugin.Dir(), configFile))
	if err != nil {
		return Plugin{}, xerrors.Errorf("failed to create plugin.yaml: %w", err)
	}
	defer f.Close()

	if err = yaml.NewEncoder(f).Encode(plugin); err != nil {
		return Plugin{}, xerrors.Errorf("yaml encode error: %w", err)
	}

	return plugin, nil
}

// Uninstall installs the plugin
func (m *Manager) Uninstall(ctx context.Context, name string) error {
	pluginDir := filepath.Join(m.pluginRoot, name)
	if !fsutils.DirExists(pluginDir) {
		m.logger.ErrorContext(ctx, "No such plugin")
		return nil
	}
	return os.RemoveAll(pluginDir)
}

// Information gets the information about an installed plugin
func (m *Manager) Information(name string) error {
	plugin, err := m.load(name)
	if err != nil {
		return xerrors.Errorf("plugin load error: %w", err)
	}

	_, err = fmt.Fprintf(m.w, `
Plugin: %s
  Description: %s
  Version:     %s
  Usage:       %s
`, plugin.Name, plugin.Description, plugin.Version, plugin.Usage)

	return err
}

// List gets a list of all installed plugins
func (m *Manager) List(ctx context.Context) error {
	s, err := m.list(ctx)
	if err != nil {
		return xerrors.Errorf("unable to list plugins: %w", err)
	}
	_, err = fmt.Fprintf(m.w, "%s\n", s)
	return err
}

func (m *Manager) list(ctx context.Context) (string, error) {
	if _, err := os.Stat(m.pluginRoot); err != nil {
		if os.IsNotExist(err) {
			return "No Installed Plugins", nil
		}
		return "", xerrors.Errorf("stat error: %w", err)
	}
	plugins, err := m.LoadAll(ctx)
	if err != nil {
		return "", xerrors.Errorf("unable to load plugins: %w", err)
	} else if len(plugins) == 0 {
		return "No Installed Plugins", nil
	}
	pluginList := []string{"Installed Plugins:"}
	for _, plugin := range plugins {
		pluginList = append(pluginList, fmt.Sprintf("  Name:    %s\n  Version: %s\n", plugin.Name, plugin.Version))
	}

	return strings.Join(pluginList, "\n"), nil
}

// Upgrade upgrades an existing plugins
func (m *Manager) Upgrade(ctx context.Context, names []string) error {
	if len(names) == 0 {
		plugins, err := m.LoadAll(ctx)
		if err != nil {
			return xerrors.Errorf("unable to load plugins: %w", err)
		} else if len(plugins) == 0 {
			m.logger.InfoContext(ctx, "No installed plugins")
			return nil
		}
		names = lo.Map(plugins, func(p Plugin, _ int) string { return p.Name })
	}
	for _, name := range names {
		if err := m.upgrade(ctx, name); err != nil {
			return xerrors.Errorf("unable to upgrade '%s' plugin: %w", name, err)
		}
	}
	return nil
}

func (m *Manager) upgrade(ctx context.Context, name string) error {
	plugin, err := m.load(name)
	if err != nil {
		return xerrors.Errorf("plugin load error: %w", err)
	}

	logger := log.With("name", name)
	logger.InfoContext(ctx, "Upgrading plugin...")
	updated, err := m.install(ctx, plugin.Repository, Options{
		// Use the current installed platform
		Platform: ftypes.Platform{
			Platform: &v1.Platform{
				OS:           plugin.Installed.Platform.OS,
				Architecture: plugin.Installed.Platform.Arch,
			},
		},
	})
	if err != nil {
		return xerrors.Errorf("unable to perform an upgrade installation: %w", err)
	}

	if plugin.Version == updated.Version {
		logger.InfoContext(ctx, "The plugin is up-to-date", log.String("version", plugin.Version))
	} else {
		logger.InfoContext(ctx, "Plugin upgraded",
			log.String("from", plugin.Version), log.String("to", updated.Version))
	}
	return nil
}

// LoadAll loads all plugins
func (m *Manager) LoadAll(ctx context.Context) ([]Plugin, error) {
	dirs, err := os.ReadDir(m.pluginRoot)
	if err != nil {
		return nil, xerrors.Errorf("failed to read %s: %w", m.pluginRoot, err)
	}

	var plugins []Plugin
	for _, d := range dirs {
		if !d.IsDir() {
			continue
		}
		plugin, err := m.loadMetadata(filepath.Join(m.pluginRoot, d.Name()))
		if err != nil {
			m.logger.WarnContext(ctx, "Plugin load error", log.Err(err))
			continue
		}
		plugins = append(plugins, plugin)
	}
	return plugins, nil
}

// Start starts the plugin
func (m *Manager) Start(ctx context.Context, name string, opts Options) (Wait, error) {
	plugin, err := m.load(name)
	if err != nil {
		return nil, xerrors.Errorf("plugin load error: %w", err)
	}

	wait, err := plugin.Start(ctx, opts)
	if err != nil {
		return nil, xerrors.Errorf("unable to run %s plugin: %w", plugin.Name, err)
	}
	return wait, nil
}

// RunWithURL runs the plugin
func (m *Manager) RunWithURL(ctx context.Context, name string, opts Options) error {
	plugin, err := m.Install(ctx, name, opts)
	if err != nil {
		return xerrors.Errorf("plugin install error: %w", err)
	}

	if err = plugin.Run(ctx, opts); err != nil {
		return xerrors.Errorf("unable to run %s plugin: %w", plugin.Name, err)
	}
	return nil
}

func (m *Manager) load(name string) (Plugin, error) {
	pluginDir := filepath.Join(m.pluginRoot, name)
	if _, err := os.Stat(pluginDir); err != nil {
		if os.IsNotExist(err) {
			return Plugin{}, xerrors.Errorf("could not find a plugin called '%s', did you install it?", name)
		}
		return Plugin{}, xerrors.Errorf("plugin stat error: %w", err)
	}

	plugin, err := m.loadMetadata(pluginDir)
	if err != nil {
		return Plugin{}, xerrors.Errorf("unable to load plugin metadata: %w", err)
	}

	return plugin, nil
}

func (m *Manager) loadMetadata(dir string) (Plugin, error) {
	filePath := filepath.Join(dir, configFile)
	f, err := os.Open(filePath)
	if err != nil {
		return Plugin{}, xerrors.Errorf("file open error: %w", err)
	}
	defer f.Close()

	var plugin Plugin
	if err = yaml.NewDecoder(f).Decode(&plugin); err != nil {
		return Plugin{}, xerrors.Errorf("yaml decode error: %w", err)
	}

	if plugin.Name == "" {
		return Plugin{}, xerrors.Errorf("'name' is empty")
	}

	// e.g. ~/.trivy/plugins/kubectl
	plugin.dir = filepath.Join(m.pluginRoot, plugin.Name)

	return plugin, nil
}

func (m *Manager) isInstalled(ctx context.Context, url string) (Plugin, bool) {
	installedPlugins, err := m.LoadAll(ctx)
	if err != nil {
		return Plugin{}, false
	}

	for _, plugin := range installedPlugins {
		if plugin.Repository == url {
			return plugin, true
		}
	}
	return Plugin{}, false
}
