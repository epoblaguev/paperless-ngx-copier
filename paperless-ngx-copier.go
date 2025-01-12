package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	FileExtensions   []string `json:"file_extensions"`
	ScanPaths        []string `json:"scan_paths"`
	OutputDir        string   `json:"output_dir"`
	HistoryStorePath string   `json:"history_store_path"`
	CalculateMD5Hash bool     `json:"calculate_md5_hash"`
}

type HistoryElement struct {
	FilePath     string
	MD5Hash      string
	ModifiedTime float32
}

type FileHistory struct {
	file_history map[string]HistoryElement
	config       Config
}

func NewFileHistory(file_history map[string]HistoryElement, config Config) FileHistory {
	obj := new(FileHistory)
	obj.config = config

	return *obj

}

func read_config() Config {
	args := os.Args[1:]
	fmt.Println(args)

	if len(args) < 1 {
		panic("Please provide valid path to config file")
	}

	data, err := os.ReadFile(args[0])
	if err != nil {
		panic(err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		panic(err)
	}

	return config
}

func main() {
	config := read_config()

	// files_copied = 0
	// files_unchanged = 0
	// files_in_error = 0

	// for _, ext := range config.FileExtensions {
	// 	file_extensions = append(file_extensions, strings.ToLower(strings.TrimLeft(ext, ".")))
	// }

	for i := 0; i < len(config.FileExtensions); i++ {
		config.FileExtensions[i] = strings.TrimLeft(strings.ToLower(config.FileExtensions[i]), ".")
	}

	fmt.Println(config.FileExtensions)

	for i := 0; i < len(config.FileExtensions); i++ {

	}
}
