use std::env;
use std::path::PathBuf;

#[derive(Debug, Clone)]
pub struct Config {
    pub host: String,
    pub port: u16,
    pub data_path: PathBuf,
}

impl Config {
    pub fn from_env() -> Self {
        let host = env::var("GIT_HOST").unwrap_or_else(|_| "0.0.0.0".to_string());
        let port = env::var("GIT_PORT")
            .unwrap_or_else(|_| "8081".to_string())
            .parse()
            .unwrap_or(8081);
        let data_path = PathBuf::from(env::var("GIT_DATA_PATH").unwrap_or_else(|_| "./git-data".to_string()));

        Config {
            host,
            port,
            data_path,
        }
    }
}
