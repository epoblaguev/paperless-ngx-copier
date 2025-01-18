use std::{arch::aarch64::float32x2_t, error::Error, fmt::Display, fs::File, io::BufReader, path::Path};
use serde_json::{self, Map};
use serde::{Serialize, Deserialize};
use clap::{builder::OsStr, Parser};
use walkdir::WalkDir;

/// Simple program to greet a person
#[derive(Parser, Debug)]
#[command(version, about, long_about = None)]
struct Args {
    /// Path to the config file
    #[arg(index=1)]
    config_path: String,
}

struct HistoryElement {
    file_path: String,
    md5_hash: String,
    modified_time: i64
}

struct FileHistory {
    file_history: Map<String, HistoryElement>,
    config: Config,
}

impl FileHistory {
    fn new(config: Config) {
        
    }
}

#[derive(Serialize, Deserialize, Debug)]
pub struct Config {
    file_extensions: Vec<String>,
    scan_paths: Vec<String>,
    output_dir: String,
    history_store_path: String,
    calculate_md5_hash: bool,
}

fn read_config(config_path: String) -> Result<Config, Box<dyn Error>> {
    let path = Path::new(&config_path);
    let file = File::open(path)?;
    let reader = BufReader::new(file);
    let config = serde_json::from_reader(reader)?;

    Ok(config)
}

fn copy_files(config: &Config) -> Result<(), Box<dyn Error>> {
    let mut files_copied = 0;
    let mut files_unchanged = 0;
    let mut files_in_error = 0;

    let file_extensions: Vec<String> = config.file_extensions.iter().map(|e| e.to_lowercase()).collect();

    println!("{}", file_extensions.join(", "));

    for scan_path in &config.scan_paths {
        for entry in WalkDir::new(scan_path) {
            let entry = entry?;
            if entry.file_type().is_file() {
                println!("{}", entry.path().display());
                println!("{}", entry.path().extension());
            }
        }
    }

    Ok(())
}

fn main() {
    let args: Args = Args::parse();

    println!("Config Path: {}", args.config_path);

    let config = read_config(args.config_path).expect("Error");
    
    dbg!(&config);

    println!("{}", &config.file_extensions.join(","));
    
    copy_files(&config);
    // dbg!(config);
}