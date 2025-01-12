package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
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
	ModifiedTime int64
}

type FileHistory struct {
	fileHistory map[string]HistoryElement
	config      Config
}

func NewFileHistory(config Config) FileHistory {
	obj := new(FileHistory)
	obj.config = config

	if err := obj.loadHistoryFile(); err != nil {
		panic(err)
	}

	return *obj
}

func (fh *FileHistory) loadHistoryFile() error {
	data, err := os.ReadFile(fh.config.HistoryStorePath)
	if err != nil {
		if os.IsNotExist(err) {
			fh.fileHistory = make(map[string]HistoryElement)
			return nil
		} else {
			return fmt.Errorf("failed to create history map, history file may be corrupted: %w", err)
		}
	}

	var content []HistoryElement
	if err := json.Unmarshal(data, &content); err != nil {
		return err
	}

	fh.fileHistory = make(map[string]HistoryElement)
	for _, elm := range content {
		fh.fileHistory[elm.FilePath] = elm
	}

	return nil
}

func (fh *FileHistory) saveHistoryFile() error {
	var content []HistoryElement
	for _, value := range fh.fileHistory {
		content = append(content, value)
	}

	jsonData, err := json.Marshal(content)
	if err != nil {
		return err
	}

	return os.WriteFile(fh.config.HistoryStorePath, jsonData, 0644)
}

func (fh *FileHistory) getElement(filePath string) HistoryElement {
	if val, ok := fh.fileHistory[filePath]; ok {
		return val
	} else {
		return HistoryElement{filePath, "-", 0}
	}
}

func (fh *FileHistory) setElement(historyElement HistoryElement) error {
	fh.fileHistory[historyElement.FilePath] = historyElement
	return fh.saveHistoryFile()
}

func readConfig() Config {
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

func generateHash(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", md5.Sum(data)), nil
}

func copyFile(source string, dest string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return err
	}

	return nil
}

func processFile(filePath string, fileInfo fs.FileInfo, config Config, fileHistory FileHistory) (bool, error) {
	modifiedTime := fileInfo.ModTime().UnixNano()
	fileHash := "NOT CALCULATED"
	if config.CalculateMD5Hash {
		var err error
		fileHash, err = generateHash(filePath)
		if err != nil {
			return false, err
		}
	}

	historicRecord := fileHistory.getElement(filePath)

	fileChanged := true
	fmt.Println("File Info: ", filePath)
	if config.CalculateMD5Hash {
		fmt.Println("\tOld MD5 Hash: ", historicRecord.MD5Hash)
		fmt.Println("\tNew MD5 Hash: ", fileHash)
		fileChanged = historicRecord.MD5Hash != fileHash
	} else {
		fmt.Println("\tOld Timestamp: ", historicRecord.ModifiedTime)
		fmt.Println("\tNew Timestamp: ", modifiedTime)
		fileChanged = historicRecord.ModifiedTime != modifiedTime
	}

	if !fileChanged {
		fmt.Println("File has not changed since it was last copied")
		return false, nil
	}

	outputFileName := fileInfo.Name()
	outputPath := path.Join(config.OutputDir, outputFileName)

	for counter := 0; ; counter++ {
		if _, err := os.Stat(outputPath); err == nil {
			filename := fmt.Sprint("(Copy ", counter, ")", outputFileName)
			outputPath = path.Join(config.OutputDir, filename)
		} else if os.IsNotExist(err) {
			break
		} else {
			return false, fmt.Errorf("error while checking if file exists in output directory: %w", err)
		}
	}

	fmt.Println("Copying file: ", filePath, " ==> ", outputPath)

	if err := copyFile(filePath, outputPath); err != nil {
		return false, fmt.Errorf("failed to copy file to new target dir: %w", err)
	}

	newHistoricRecord := HistoryElement{filePath, fileHash, modifiedTime}
	if err := fileHistory.setElement(newHistoricRecord); err != nil {
		return false, fmt.Errorf("file copied successfuly, but history file could not be updated: %w", err)
	}

	return true, nil
}

func main() {
	config := readConfig()

	filesCopied := 0
	filesUnchanged := 0
	filesInError := 0

	for i := 0; i < len(config.FileExtensions); i++ {
		config.FileExtensions[i] = strings.TrimLeft(strings.ToLower(config.FileExtensions[i]), ".")
	}

	fmt.Println(config.FileExtensions)

	fileHistory := NewFileHistory(config)

	for _, scanPath := range config.ScanPaths {
		filepath.Walk(scanPath, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			extension := strings.ToLower(strings.TrimLeft(filepath.Ext(path), "."))
			if !slices.Contains(config.FileExtensions, extension) {
				fmt.Println("Extension ", extension, " not in list ", config.FileExtensions)
				return nil
			}

			fmt.Println(path)

			result, err := processFile(path, info, config, fileHistory)
			if err != nil {
				filesInError++
			} else if result {
				filesCopied++
			} else {
				filesUnchanged++
			}
			return nil
		})
	}

	fmt.Println("\nCOMPLETE")
	fmt.Println("\tFiles Copied: ", filesCopied)
	fmt.Println("\tFiles Unchanged: ", filesUnchanged)
	fmt.Println("\tFiles With Errors: ", filesInError)
}
