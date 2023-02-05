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

type Repo struct {
	Remote  string `yaml:"remote,omitempty"`
	Version string
	Commit  string `yaml:"commit,omitempty"`
	Deps    []Repo `yaml:"deps,omitempty"`
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
		[]Repo{},
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

func (r *Repo) FetchAndCacheDeps() error {
	for _, dep := range r.Deps {
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

	url := fmt.Sprintf("https://%s", r.Remote)
	if strings.Contains(r.Remote, "https://") {
		return nil, errors.New("invalid remote, should be of form (dns.name)/org/repo")
	}

	localPath = path.Join(localPath, r.Remote)

	v.CreateAtRev(localPath, url, r.Commit)
	// // TODO: Auth
	// repo, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
	// 	URL: r.Remote,
	// })

	// fmt.Println("Cloned", repo)
	// if err != nil {
	// 	return nil, err
	// }

	// w, err := repo.Worktree()
	// if err != nil {
	// 	return nil, err
	// }

	// err = w.Checkout(&git.CheckoutOptions{
	// 	Hash: plumbing.NewHash(r.Commit),
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// files, _ := w.Filesystem.ReadDir(".")

	// for _, f := range files {

	// }

	// obj, err := repo.CommitObject(plumbing.NewHash(r.Commit))

	// fmt.Println(obj.Message)
	return nil, err
}

func (r *Repo) Validate() error {
	path := r.local
	os.Chdir(path)

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
		r.Deps = existingRepo.Deps
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