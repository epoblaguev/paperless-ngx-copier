# paperless-ngx-copier
A script for copying new files to the paperless-ngx consume directory

## Quick Start ##
1. Download proper binary from the [releases](https://github.com/epoblaguev/paperless-ngx-copier/releases).
2. Create a JSON configuration file based on the [example config](config-example.json) in this repo.
3. Run the binary with the path to the config file as the only parameter. For example
      ```bash
      ./paperless-ngx-copier_linux_x86 /path/to/config.json
      ```
      
## Config File ##
```json
{
    "file_extensions": ["pdf", "png", "jpg"],  // List of file extensions to be copied
    "scan_paths": [  // List of paths to scan for documents to copy
        "/path/to/documents/",
        "/another/path/to/documents/"
    ],
    "output_dir": "/path/to/paperless-ngx/consume/folder",  // Path to the Paperless-ngx 'consume' folder
    "history_store_path": "/path/to/history.json",  // Path to the JSON file that contains a history of copied files. This file stores a list of all copied files and their md5 hash / timestamp
    "calculate_md5_hash": true  // If true - use MD5 hash to check if a file has changed since it was last copied. If false - use modified date
}
```
