package gomake

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Mode uint

const (
	Binary Mode = iota
	SharedLibrary
	StaticLibrary
	ObjectFile
)

type Builder struct {
	files   []string
	mode    Mode
	out_dir string
}

func NewBuilder(mode Mode) Builder {
	return Builder{
		files:   make([]string, 0, 10),
		mode:    mode,
		out_dir: ".",
	}
}

func build(out_dir, file string) {
	cmd := exec.Command("gcc", "-c", file)
	cmd.Dir = out_dir
	fmt.Printf("Compiling %v\n", file)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func (b *Builder) Build() error {
	var out_file string
	for _, file := range b.files {
		out_file = strings.Split(file, ".")[0] + ".o"
		if IsFileUpdated(b.out_dir, file, out_file) {
			build(b.out_dir, file)
		}

	}
	return nil
}

func (b *Builder) Add(file_name string) {
	b.files = append(b.files, file_name)
}
func (b *Builder) CurrentDir(path string) {
	b.out_dir = path
}

func GetTime(name string) (*time.Time, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// reading file state
	file_info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	t := file_info.ModTime()
	return &t, nil

}
func IsFileUpdated(out_dir, file_name, output_name string) bool {
	t_file, err1 := GetTime(out_dir + "/" + file_name)
	t_output, err2 := GetTime(out_dir + "/" + output_name)

	if err1 != nil {
		log.Fatal(err1)
	}
	if err2 != nil {
		build(out_dir, file_name)
		return false
	}

	d := t_file.Sub(*t_output)

	if d.Seconds() > 3 {
		return true
	}
	fmt.Printf("Up to date %v\n", file_name)
	return false
}
