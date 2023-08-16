package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/samber/lo"
)

const (
	AppConfigTOML = "app_config.toml"
)

var (
	app_config      *AppConfig
	app_root        = filepath.Dir(executable())
	app_config_path = filepath.Join(app_root, AppConfigTOML)
)

func init() {
	readAppConfig()
	fmt.Println(app_config)
}

// executable returns lo.Must1(os.Executable())
func executable() string {
	return lo.Must1(os.Executable())
}

func readAppConfig() {
	data := lo.Must(os.ReadFile(app_config_path))
	lo.Must0(toml.Unmarshal(data, &app_config))
}
