package appyaml

import (
	"os"
	"path/filepath"

	yaml "github.com/goccy/go-yaml"
)

type (
	AppYAML struct {
		Runtime      string
		EnvVariables map[string]string `yaml:"env_variables"`
		Includes     []Include
		Handlers     []Handler
		children     map[string]*AppYAML
		intent       intent
	}
	Handler struct {
		URL         string
		StaticDir   string
		StaticFiles string
		Script      string
	}
	intent struct {
		Entrypoint string
		FullPath   string
	}
)

type Include string

func Load(fpath string) (*AppYAML, error) {
	app := &AppYAML{
		intent:   intent{Entrypoint: fpath},
		children: map[string]*AppYAML{},
	}

	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}

	app.intent.FullPath, err = filepath.Abs(f.Name())
	if err != nil {
		return nil, err
	}

	if err := yaml.NewDecoder(f).Decode(app); err != nil {
		return nil, err
	}

	if err := app.Resolve(); err != nil {
		return nil, err
	}

	if _, err := app.LoadEnv(); err != nil {
		return nil, err
	}

	return app, err
}

func (app *AppYAML) Resolve() error {
	dir := filepath.Dir(app.intent.FullPath)
	for _, include := range app.Includes {
		fpath := filepath.Join(dir, string(include))
		child, err := Load(fpath)
		if err != nil {
			return err
		}
		if err := child.Resolve(); err != nil {
			return err
		}
		app.children[fpath] = child
	}
	return nil
}

func (app *AppYAML) LoadEnv() (map[string]string, error) {
	loaded := map[string]string{}
	for key, value := range app.EnvVariables {
		if err := os.Setenv(key, value); err != nil {
			return loaded, err
		}
		loaded[key] = value
	}
	for _, child := range app.children {
		if loadedByChild, err := child.LoadEnv(); err != nil {
			return loaded, err
		} else {
			for key, value := range loadedByChild {
				loaded[key] = value
			}
		}
	}
	return loaded, nil
}
