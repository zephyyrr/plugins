package plugins

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
)

//Loads the plugins at the given path
//Returns a valid plugin or a non-nil error
func Load(path string) (Plugin, error) {
	cmd := exec.Command(path)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err == nil {
		go io.Copy(os.Stderr, stderr)
	}

	go cmd.Start()

	return NewRemotePlugin(stdout, stdin)

}

//
func LoadAll(path string) ([]Plugin, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	dirinfo, err := dir.Stat()
	if err != nil {
		return nil, err
	}

	if !dirinfo.IsDir() {
		return nil, NotDirectory
	}

	files, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}
	plugins := make([]Plugin, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Not a directory. A file. Attempt execution
		pl, pl_err := Load(filepath.Join(path, file.Name()))
		if err == nil {
			plugins = append(plugins, pl)
		} else {
			err = multierror.Append(err, pl_err)
		}
	}

	return plugins, err
}
