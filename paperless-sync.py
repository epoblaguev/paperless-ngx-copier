from dataclasses import dataclass, asdict
import json
import sys
from os import path, walk
import hashlib
import shutil

@dataclass
class HistoryElement:
    file_path: str
    md5_hash: str
    modified_time: float

@dataclass
class Config:
    file_extensions: list[str]
    scan_paths: list[str]
    output_dir: str
    history_store_path: str
    calculate_md5_hash: bool

class FileHistory:
    __file_history: dict[str, HistoryElement]
    __config: Config

    def __init__(self, config: Config) -> None:
        self.__config = config
        self.__file_history = self.__load_history_file()

    def getElement(self, file_path):
        return self.__file_history.get(file_path, HistoryElement(file_path, '-', '-'))

    def setElement(self, element: HistoryElement):
        self.__file_history[element.file_path] = element
        self.__save_history_file()

    def __save_history_file(self):
        json_content = [asdict(elm) for elm in self.__file_history.values()]
        with open(self.__config.history_store_path, 'w') as file:
            json.dump(json_content, file)
            
    def __load_history_file(self) -> dict[str, HistoryElement]:
        file_path = self.__config.history_store_path
        if not path.exists(file_path):
            print(f'History file does not exist at "{file_path}"\n ... new file will be created.')
            return {}
        
        with open(file_path, 'r') as file:
            content = json.load(file)

        try:
            history_elements = (HistoryElement(**elm) for elm in content)
            return {elm.file_path: elm for elm in history_elements}
        except:
            raise Exception('Failed to create history dict, history file may be corrupted')


def generate_hash(file_path: str) -> str:
    with open(file_path, "rb") as f:
        file_hash = hashlib.md5()
        while chunk := f.read(8192):
            file_hash.update(chunk)
    return file_hash.hexdigest()


def process_file(file_path: str, config: Config, file_history: FileHistory) -> bool:
    file_hash = generate_hash(file_path) if config.calculate_md5_hash else 'NOT CALCULATED'
    modified_time = path.getmtime(file_path)

    historic_record = file_history.getElement(file_path)

    file_changed = True
    print(f'\nFile Info: {file_path}')
    if config.calculate_md5_hash:
        print(f'\tOld MD5 Hash: {historic_record.md5_hash}')
        print(f'\tNew MD5 Hash: {file_hash}')
        file_changed = historic_record.md5_hash != file_hash
    else:
        print(f'\tOld Timestamp: {historic_record.modified_time}')
        print(f'\tNew Timestamp: {modified_time}')
        file_changed = historic_record.modified_time != modified_time

    if not file_changed:
        print('File has not changed since it was last copied')
        return False

    output_filename = path.basename(file_path)
    output_path = path.join(config.output_dir, output_filename)
    counter = 1
    while path.exists(output_path):
        filename = f'(Copy {counter}) {output_filename}'
        output_path = path.join(config.output_dir, filename)
        counter += 1

    print(f'Copying file: {file_path} ==> {output_path}')
    shutil.copy2(file_path, output_path)

    new_historic_record = HistoryElement(file_path, file_hash, modified_time)
    file_history.setElement(new_historic_record)

    return True


def main(config_path: str):
    files_copied = 0
    files_unchanged = 0
    files_in_error = 0

    # Read Config
    with open(config_path) as config_file:
        config = Config(**json.load(config_file))
    
    file_extensions = tuple(f'.{e.lstrip(".")}'.lower() for e in config.file_extensions)

    file_history = FileHistory(config)

    for scan_path in config.scan_paths:
        for root, dirs, files in walk(scan_path):
            for file in files:
                if not file.lower().endswith(file_extensions):
                    continue
                file_path = path.join(root, file)
                try:
                    result = process_file(file_path, config, file_history)
                    if result:
                        files_copied += 1
                    else:
                        files_unchanged += 1
                except Exception as ex:
                    print(ex)
                    files_in_error += 1

    print(f"""
    \nCOMPLETE:
    \tFiles Copied: {files_copied}
    \tFiles Unchanged: {files_unchanged}
    \tFiles With Errors: {files_in_error}
    """)

if __name__ == '__main__':
    try:
        config_path = sys.argv[1]
    except:
        raise Exception('Please provide valid path to config file')

    main(config_path)
   
    