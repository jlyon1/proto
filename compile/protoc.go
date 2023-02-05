package compile

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/jlyon1/proto/repo"
)

var COMMAND = "protoc"

type ProtocBuilder struct {
	repo repo.Repo
}

func NewBuilder(r repo.Repo) ProtocBuilder {
	return ProtocBuilder{
		repo: r,
	}
}

func (p *ProtocBuilder) GetCommand() ([]string, error) {
	cmd := COMMAND
	args := []string{cmd}

	srcPath, err := p.repo.GetAbsolutePath()
	if err != nil {
		return nil, err
	}

	args = append(args, fmt.Sprintf("-I%s", srcPath))

	for _, dep := range p.repo.Build.Deps {
		depPath, err := dep.GetDependencyPath()
		if err != nil {
			return nil, err
		}
		absDep, err := filepath.Abs(path.Join(srcPath, depPath))

		args = append(args, fmt.Sprintf("-I%s", absDep))
	}

	for _, plugin := range p.repo.Build.Plugins {
		absOut, err := filepath.Abs(path.Join(srcPath, plugin.Out))
		if err != nil {
			return nil, err
		}
		base := fmt.Sprintf("--%s_out=%s", plugin.Name, absOut)
		args = append(args, base)
		optBase := fmt.Sprintf("--%s_opt=%s", plugin.Name, plugin.Opt)
		args = append(args, optBase)
	}

	files, err := p.repo.GetAllLocalProtoFiles()
	if err != nil {
		return nil, err
	}
	args = append(args, files...)

	return args, nil

}
