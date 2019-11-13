package file

import (
	"io/ioutil"
	"log"
	"path"
	"os"
)

type DirectoryCleaner struct {
	Directories []string
}

func NewDirectoryCleaner(directories []string) DirectoryCleaner {
	return DirectoryCleaner{
		Directories: directories,
	}
}

func (f DirectoryCleaner) RemoveFiles() {
	for _, directory := range f.Directories {
		files, err := ioutil.ReadDir(directory)
		if err != nil {
			log.Printf("failed to read directory: %s\n", err.Error())
			continue
		}
		log.Printf("removing all files from directory: %s\n", directory)

		for _, file := range files {
			err := os.Remove(path.Join(directory,file.Name()))
			if err != nil {
				log.Printf("failed to delete file: %s %s\n", file.Name(), err.Error())
				continue
			}
		}
	}
}
