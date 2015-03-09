package plugins

import (
	"os"
	"os/exec"
	"path/filepath"
)

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

	go cmd.Start()

	return NewRemotePlugin(stdout, stdin)

}

func LoadAll(path string) ([]Plugin, error) {
	dir, err := os.Open(path)
	dirinfo, err := dir.Stat()
	if !dirinfo.IsDir() {
		return nil, NotDirectory
	}

	files, err := dir.Readdir(0)
	plugins := make([]Plugin, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Not a directory. A file. Attempt execution
		Load(filepath.Join(path, file.Name()))
	}

	return plugins, err
}
