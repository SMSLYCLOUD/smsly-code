use thiserror::Error;

#[derive(Error, Debug)]
pub enum GitError {
    #[error("Repository not found: {0}")]
    NotFound(String),

    #[error("Repository already exists: {0}")]
    AlreadyExists(String),

    #[error("Invalid repository name: {0}")]
    InvalidName(String),

    #[error("Git error: {0}")]
    GitError(#[from] git2::Error),

    #[error("IO error: {0}")]
    IoError(#[from] std::io::Error),

    #[error("Permission denied: {0}")]
    PermissionDenied(String),
}
