package main // import "github.com/uphy/java-source-collector"

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/uphy/java-source-collector/enc"
	"github.com/uphy/java-source-collector/java"
)

type Config struct {
	Input  []Input `json:"input"`
	Output Output  `json:"output"`
}

type Input struct {
	Path     string  `json:"path"`
	Encoding *string `json:"encoding"`
}

type Output Input

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "config.yml", "")
	flag.Parse()
	if err := run(configFile); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(configFile string) error {
	config, err := readConfig(configFile)
	if err != nil {
		return err
	}
	defaultEncoding := "utf8"
	for _, input := range config.Input {
		if input.Encoding == nil {
			input.Encoding = &defaultEncoding
		}
		processJavaFiles(input, config.Output)
	}
	return nil
}

func processJavaFiles(input Input, output Output) error {
	return filepath.Walk(input.Path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if ext != ".java" {
			return nil
		}
		// java file
		if err := processJavaFile(path, *input.Encoding, output); err != nil {
			log.Println("failed to process java file:", err, ":", path)
		}
		return nil
	})
}

func processJavaFile(path string, encoding string, output Output) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	r, err := enc.NewReader(f, encoding)
	if err != nil {
		return err
	}
	pkg, err := java.DetectPackage(r)
	if err != nil {
		return err
	}
	_, filename := filepath.Split(path)
	dstDir := filepath.Join(output.Path, strings.Replace(pkg, ".", "/", -1))
	if _, err := os.Stat(dstDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dstDir, 0777); err != nil {
			return err
		}
	}
	dst := filepath.Join(dstDir, filename)
	if err := syncFile(path, dst); err != nil {
		return err
	}
	return nil
}

func syncFile(src string, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	return err
}

func readConfig(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
