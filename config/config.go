package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App    `yaml:"server-cmd"`
		HTTP   `yaml:"http"`
		Log    `yaml:"logger"`
		PG     `yaml:"postgres"`
		JWT    `yaml:"jwt"`
		WeChat `yaml:"we-chat"`
		File   `yaml:"file"`
		//RMQ  `yaml:"rabbitmq"`
	}

	// App -.
	App struct {
		Name    string ` yaml:"name"    env:"APP_NAME"`
		Version string ` yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string ` yaml:"port" env:"HTTP_PORT"`
	}

	// Log -.
	Log struct {
		Level string ` yaml:"log_level"   env:"LOG_LEVEL"`
	}

	// PG -.
	PG struct {
		PoolMax int    `yaml:"pool_max" env:"PG_POOL_MAX" `
		URL     string `               env:"PG_URL"`
	}

	// JWT -.
	JWT struct {
		Secret string `yaml:"secret" env:"JWT_SECRET"`
		Expire int    `yaml:"expire" env:"JWT_EXPIRE"`
	}

	// WeChat -.
	WeChat struct {
		AppID     string `yaml:"app_id"     env:"WECHAT_APP_ID"`
		AppSecret string `yaml:"app_secret" env:"WECHAT_APP_SECRET"`
		Token     string `yaml:"token"      env:"WECHAT_TOKEN"`
		Encoding  string `yaml:"encoding"   env:"WECHAT_ENCODING"`
	}

	// File -.
	File struct {
		RootDir string `yaml:"root_dir" env:"FILE_ROOT_DIR"`
	}
)

// NewConfig returns server-cmd config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	//data, e := ioutil.ReadFile("./config/config.yml")
	//if e != nil {
	//	return nil, fmt.Errorf("config error: %w", e)
	//}
	//fmt.Println(string(data))

	//err := cleanenv.ReadConfig("./config/config.yml", cfg)
	//if err != nil {
	//	return nil, fmt.Errorf("config --- error: %w", err)
	//}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	if cfg.PG.PoolMax <= 0 {
		cfg.PG.PoolMax = 10
	}

	return cfg, nil
}
