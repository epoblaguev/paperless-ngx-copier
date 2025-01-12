import { walk } from "jsr:@std/fs/walk";
import { crypto } from "jsr:@std/crypto";
import { encodeHex } from "jsr:@std/encoding/hex";

interface Config {
  file_extensions: string[];
  scan_paths: string[];
  output_dir: string[];
  history_store_path: string;
  calculate_md5_hash: boolean;
}

interface HistoryElement {
  filePath: string;
  md5Hash: string;
  modifiedTime: number;
}

class FileHistory {
  private config: Config;
  private fileHistory: Map<string, HistoryElement>;

  constructor(config: Config) {
    this.config = config;
    this.fileHistory = this.loadHistoryFile();
  }

  getElement(filePath: string) {
    return this.fileHistory.get(filePath) ??
      { filePath, md5Hash: "-", modifiedTime: "-" };
  }

  setElement(element: HistoryElement) {
    this.fileHistory.set(element.filePath, element);
    this.saveHistoryFile();
  }

  private saveHistoryFile() {
    const jsonContent = JSON.stringify(this.fileHistory.values().toArray());
    Deno.writeTextFile(this.config.history_store_path, jsonContent);
  }

  private loadHistoryFile(): Map<string, HistoryElement> {
    const filePath = this.config.history_store_path;

    let fileContent: string;
    try {
      fileContent = Deno.readTextFileSync(filePath);
    } catch (_) {
      console.warn(`History file does not exist at ${filePath}`);
      console.warn("... new file will be created.");
      return new Map();
    }

    try {
      const fileJson: HistoryElement[] = JSON.parse(fileContent);
      return fileJson.reduce((acc: Map<string, HistoryElement>, item) => {
        acc.set(item.filePath, item);
        return acc;
      }, new Map());
    } catch (_) {
      throw "Failed to create history dict, history file may be corrupted";
    }
  }
}

async function generateHash(filePath: string) {
  const file = await Deno.readFile(filePath);
  const fileHashBuffer = await crypto.subtle.digest("MD5", file);
  return encodeHex(fileHashBuffer);
}

async function processFile(filePath: string, config: Config, fileHistory: FileHistory): Promise<boolean> {
  const fileHash$ = config.calculate_md5_hash ? generateHash(filePath) : Promise.resolve('NOT CALCULATED');
  const fileStat$ = Deno.stat(filePath);

  const historicRecord = fileHistory.getElement(filePath)

  console.log(`Before await all ${filePath}`)
  const [fileHash, fileStat] = await Promise.all([fileHash$, fileStat$]);
  console.log(`After await all ${filePath} -> ${fileHash}`)

  let fileChanged: boolean;
  console.log(`File Info: ${filePath}`)
  if(config.calculate_md5_hash) {
    console.log(`\tOld MD5 Hash: ${historicRecord.md5Hash}`)
    console.log(`\tNew MD5 Hash: ${fileHash}`)
    fileChanged = historicRecord.md5Hash != fileHash
  } else {
    console.log(`\tOld Timestamp: ${historicRecord.modifiedTime}`)
    console.log(`\tNew Timestamp: ${fileStat.mtime?.getTime()}`)
    fileChanged = historicRecord.modifiedTime != fileStat.mtime?.getTime()
  }

  if(!fileChanged) {
    console.log('File has not changed since it was last copied')
    return false
  }


}

async function main(configPath: string) {
  let filesCopied = 0;
  let filesUnchanged = 0;
  let filesInError = 0;

  const config: Config = JSON.parse(await Deno.readTextFile(configPath));

  const fileHistory = new FileHistory(config);

  for (const scanPath of config.scan_paths) {
    const results$ = []
    for await (const dirEntry of walk(scanPath, { exts: config.file_extensions })) {
      console.log("Processing: ", dirEntry.path);
      const result$ = processFile(dirEntry.path, config, fileHistory);
      results$.push(result$)
    }

    for (const result of await Promise.all(results$)) {
      
    }
  }
}

if (import.meta.main) {
  const configPath = Deno.args[0];

  if (!configPath) {
    throw "Please provide valid path to config file";
  }

  main(configPath);
}
