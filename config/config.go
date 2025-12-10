package config

import (
	"bytes"
	"strings"

	"github.com/spf13/viper"
)

var Version string = "dev"

var CommitHash string = ""

var BuildTime string = ""

var defaultConfig = []byte(`
app: "structcopy-gen"
log_level: "info"
log_format: "json"
flags:
  version: false
  standalone: false
`)

type (
	AppConfig struct {
		App       string   `mapstructure:"app"`
		LogLevel  string   `mapstructure:"log_level"`
		LogFormat string   `mapstructure:"log_format"`
		CliFlags  CliFlags `mapstructure:"flags"`
	}

	CliFlags struct {
		Version    bool `mapstructure:"version"`
		Standalone bool `mapstructure:"standalone"`
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
