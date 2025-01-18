package main

import (
	"reflect"
	"testing"
)

func TestFileHistory_getElement(t *testing.T) {
	type fields struct {
		fileHistory map[string]HistoryElement
		config      Config
	}
	type args struct {
		filePath string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   HistoryElement
	}{
		{
			name: "Element Exists",
			fields: fields{
				fileHistory: map[string]HistoryElement{
					"/path/to/file": {
						FilePath:     "/path/to/file",
						MD5Hash:      "",
						ModifiedTime: 12345,
					},
				},
				config: Config{},
			},
			args: args{filePath: "/path/to/file"},
			want: HistoryElement{
				FilePath:     "/path/to/file",
				MD5Hash:      "",
				ModifiedTime: 12345,
			},
		},
		{
			name: "Element Doesn't Exist",
			fields: fields{
				fileHistory: map[string]HistoryElement{
					"/path/to/file": {
						FilePath:     "/path/to/file",
						MD5Hash:      "",
						ModifiedTime: 12345,
					},
				},
				config: Config{},
			},
			args: args{
				filePath: "/invalid/path",
			},
			want: HistoryElement{
				FilePath:     "/invalid/path",
				MD5Hash:      "-",
				ModifiedTime: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fh := &FileHistory{
				fileHistory: tt.fields.fileHistory,
				config:      tt.fields.config,
			}
			if got := fh.getElement(tt.args.filePath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FileHistory.getElement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_formatExtensions(t *testing.T) {
	type fields struct {
		FileExtensions   []string
		ScanPaths        []string
		OutputDir        string
		HistoryStorePath string
		CalculateMD5Hash bool
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "lowercase and trim period",
			fields: fields{
				FileExtensions: []string{".txt", ".PDF", ".xYz"},
			},
			want: []string{"txt", "pdf", "xyz"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				FileExtensions:   tt.fields.FileExtensions,
				ScanPaths:        tt.fields.ScanPaths,
				OutputDir:        tt.fields.OutputDir,
				HistoryStorePath: tt.fields.HistoryStorePath,
				CalculateMD5Hash: tt.fields.CalculateMD5Hash,
			}
			config.formatExtensions()
			if !reflect.DeepEqual(config.FileExtensions, tt.want) {
				t.Errorf("Config.formatExtensions() = %v, want %v", config.FileExtensions, tt.want)
			}
		})
	}
}
