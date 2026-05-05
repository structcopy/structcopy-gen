package config

import (
	"bytes"
	"runtime/debug"
	"strings"

	"github.com/spf13/viper"
)

var Version string = "dev"

var CommitHash string = ""

var BuildTime string = ""

var defaultConfig = []byte(`
app: "structcopy-gen"
log_enabled: false
log_level: "info"
log_format: "json"
flag:
  version: false
  standalone: false
  log_enabled: false
  debug_enabled: false
  dry_run: false
`)

type (
	AppConfig struct {
		App        string  `mapstructure:"app"`
		LogEnabled bool    `mapstructure:"log_enabled"`
		LogLevel   string  `mapstructure:"log_level"`
		LogFormat  string  `mapstructure:"log_format"`
		CliFlags   CliFlag `mapstructure:"flag"`
	}

	CliFlag struct {
		Version      bool   `mapstructure:"version"`
		Standalone   bool   `mapstructure:"standalone"`
		OutputPath   string `mapstructure:"output_path"`
		InputPath    string `mapstructure:"input_path"`
		LogEnabled   bool   `mapstructure:"log_enabled"`
		DebugEnabled bool   `mapstructure:"debug_enabled"`
		DryRun       bool   `mapstructure:"dry_run"`
	}
)

func LoadAppConfig(configFilePath, prefix string) (*AppConfig, error) {
	cfg := &AppConfig{}

	// load default config
	v := viper.New()
	v.SetConfigType("yaml")
	err := v.ReadConfig(bytes.NewBuffer(defaultConfig))
	if err != nil {
		return nil, err
	}

	// load config from file
	if configFilePath != "" {
		// viper.SetConfigType("yaml") // same, no need
		v.SetConfigName("config")
		v.AddConfigPath(configFilePath)
		err = v.MergeInConfig()
		if err != nil {
			return nil, err
		}
	}

	if prefix != "" {
		v.SetEnvPrefix(prefix)
	}
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	// load config from env
	v.AutomaticEnv()

	// store to output pointer
	err = v.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	syncEnvKeys(v)

	return cfg, nil
}

func syncEnvKeys(v *viper.Viper) {
	// _ = os.Setenv(envs.EnvKey_EncryptKey, v.GetString(envs.EnvKey_EncryptKey))
	// _ = os.Setenv(envs.EnvKey_DecryptKey, v.GetString(envs.EnvKey_DecryptKey))
}

func GetBuildInfoVersion() string {
	currentVersion := Version
	info, ok := debug.ReadBuildInfo()
	if ok && info.Main.Version != "(devel)" {
		currentVersion = info.Main.Version
	}

	return currentVersion
}

func GetBuildInfoRevision() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			return setting.Value
		}
	}
	return "no-vcs-data"
}
