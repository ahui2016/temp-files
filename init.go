package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/samber/lo"
)

const (
	MB = 1024 * 1024
)

const (
	AppConfigTOML    = "app_config.toml"
	PublicFolderName = "public"
	FilesFolderName  = "files"
)

var (
	Separator       = string(filepath.Separator)
	app_root        = filepath.Dir(executable())
	app_config      *AppConfig
	app_config_path = filepath.Join(app_root, AppConfigTOML)
	public_folder   = filepath.Join(app_root, PublicFolderName)
	files_folder    = filepath.Join(app_root, FilesFolderName)
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
