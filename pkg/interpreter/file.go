package interpreter

import (
	"io/ioutil"
	"strings"

	"github.com/mat285/interpreter/pkg/parser"
)

func save(file string, ctx parser.Context) error {
	data := []byte(ctx.Source())
	return ioutil.WriteFile(file, data, 0777)
}

func load(file string) ([]string, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(data), "\n"), nil
}
