package repo

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/vcs"
	"gopkg.in/yaml.v3"
)

type Plugin struct {
	Name string
	Out  string
	Opt  string
}

type Build struct {
	Plugins []Plugin
	Output  string
	Deps    []Repo `yaml:"deps,omitempty"`
}

type Repo struct {
	Remote  string `yaml:"remote,omitempty"`
	Version string `yaml:"version,omitempty"`
	Commit  string `yaml:"commit,omitempty"`
	Build   Build  `yaml:"build,omitempty"`
	local   string
}

func NewRepo(local string, remote string, version string) (*Repo, error) {
	path, err := filepath.Abs(local)
	if err != nil {
		return nil, err
	}
	return &Repo{
		remote,
		version,
		"",
		Build{
			Plugins: []Plugin{
				{
					Name: "go",
					Out:  "gen/",
				},
			},
			Deps: []Repo{},
		},
		path,
	}, nil
}

func FromFile(file string) (*Repo, error) {
	var existingRepo Repo
	existingFile, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(existingFile, &existingRepo)
	existingRepo.local = path.Dir(file)
	if err != nil {
		return nil, err
	} else {
		return &existingRepo, nil
	}
}

// FetchAndCacheDeps removes all locally cached dependencies, then tries to re-download them
// eventually it will cache them and only download if necessary
func (r *Repo) FetchAndCacheDeps() error {
	err := os.RemoveAll(path.Join(r.local, ".proto/"))
	if err != nil {
		return err
	}

	for _, dep := range r.Build.Deps {
		dep.local = r.local
		_, err := dep.FetchAndCache()
		fmt.Println(err)
		if err != nil {
			return err
		}
	}
	return nil
}

// FetchAndCache fetches a repo from its remote and stores it in the .proto/ directory to be used
// for compilation. Must be run from a directory containing a proto.yaml
func (r *Repo) FetchAndCache() (*Repo, error) {
	localPath := path.Join(r.local, ".proto/")
	print(localPath)
	_, err := os.Stat(localPath)
	if err != nil && err == os.ErrNotExist {
		err := os.Mkdir(localPath, 0644)
		if err != nil {
			return nil, err
		}
	}
	v := vcs.ByCmd("git")

	root, err := vcs.RepoRootForImportPath(r.Remote, false)
	if err != nil {
		return nil, err
	}

	localPath = path.Join(localPath, root.Root)

	v.CreateAtRev(localPath, root.Repo, r.Commit)

	return nil, err
}

func (r *Repo) Validate() error {
	path := r.local

	err := os.Chdir(path)
	if err != nil {
		return err
	}

	_, err = os.Stat("proto.yaml")
	if err != nil {
		return errors.New("no proto.yaml detected")
	}

	return nil
}

func (r *Repo) GetAbsolutePath() (string, error) {
	return filepath.Abs(r.local)
}

// If this repo is a dependency get the path to it's root
func (r *Repo) GetDependencyPath() (string, error) {
	root, err := vcs.RepoRootForImportPath(r.Remote, false)
	if err != nil {
		return "", err
	}
	return path.Join(".proto", root.Root), nil
}

func (r *Repo) GetAllLocalProtoFiles() ([]string, error) {
	files := []string{}
	err := filepath.Walk(".",
		func(fullFilename string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// ignore remote deps
			if strings.Contains(path.Dir(fullFilename), ".proto") || info.IsDir() {
				return nil
			}

			if strings.Contains(info.Name(), ".proto") {
				path, err := filepath.Abs(path.Join(path.Dir(fullFilename), info.Name()))
				if err != nil {
					return err
				}
				files = append(files, path)
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (r *Repo) Init() error {
	path, err := filepath.Abs(r.local)
	if err != nil {
		return err
	}

	os.Chdir(path)
	info, err := os.Stat("proto.yaml")

	var existingRepo *Repo
	reinitialize := false

	if err == nil {
		existingRepo, err = FromFile(info.Name())
		if err != nil {
			return err
		}
		reinitialize = true
	}

	if reinitialize {
		r.Remote = existingRepo.Remote
		r.Build = existingRepo.Build
		fmt.Printf("Reinitializing existing repo in %s\n", r.local)
	}
	yamlData, err := yaml.Marshal(r)
	if err != nil {
		return err
	}

	err = os.WriteFile("proto.yaml", yamlData, 0644)
	if err != nil {
		return err
	}

	if err := r.Validate(); err != nil {
		return err
	}

	return nil
}
