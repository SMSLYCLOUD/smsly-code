use std::path::{Path, PathBuf};
use std::fs;
use crate::types::{RepoHandle, RepoInfo};
use crate::error::GitError;
use git2::{Repository, BranchType};
use chrono::{TimeZone, Utc};

fn get_repo_path(base_path: &Path, owner: &str, name: &str) -> PathBuf {
    base_path.join(owner).join(format!("{}.git", name))
}

fn validate_name(name: &str) -> Result<(), GitError> {
    if name.contains("..") || name.contains('/') || name.contains('\\') || name.chars().any(|c| c.is_control()) {
        return Err(GitError::InvalidName(name.to_string()));
    }
    Ok(())
}

pub fn init_bare(base_path: &Path, owner: &str, name: &str) -> Result<RepoHandle, GitError> {
    validate_name(name)?;
    let path = get_repo_path(base_path, owner, name);

    if path.exists() {
        return Err(GitError::AlreadyExists(name.to_string()));
    }

    if let Some(parent) = path.parent() {
        fs::create_dir_all(parent)?;
    }

    let mut init_opts = git2::RepositoryInitOptions::new();
    init_opts.bare(true);
    init_opts.initial_head("main");
    let repo = Repository::init_opts(&path, &init_opts)?;

    Ok(RepoHandle {
        path,
        owner: owner.to_string(),
        name: name.to_string(),
        inner: repo,
    })
}

pub fn open(base_path: &Path, owner: &str, name: &str) -> Result<RepoHandle, GitError> {
    validate_name(name)?;
    let path = get_repo_path(base_path, owner, name);

    if !path.exists() {
        return Err(GitError::NotFound(name.to_string()));
    }

    let repo = Repository::open(&path)?;

    Ok(RepoHandle {
        path,
        owner: owner.to_string(),
        name: name.to_string(),
        inner: repo,
    })
}

pub fn delete(base_path: &Path, owner: &str, name: &str) -> Result<(), GitError> {
    validate_name(name)?;
    let path = get_repo_path(base_path, owner, name);

    if !path.exists() {
        return Err(GitError::NotFound(name.to_string()));
    }

    fs::remove_dir_all(&path)?;
    Ok(())
}

pub fn exists(base_path: &Path, owner: &str, name: &str) -> bool {
    if validate_name(name).is_err() {
        return false;
    }
    let path = get_repo_path(base_path, owner, name);
    path.exists()
}

pub fn get_info(handle: &RepoHandle) -> Result<RepoInfo, GitError> {
    let repo = &handle.inner;

    let is_empty = repo.is_empty()?;
    let mut default_branch = "HEAD".to_string();
    let mut last_commit_at = None;

    if !is_empty {
        if let Ok(head) = repo.head() {
            if let Some(name) = head.name() {
                default_branch = name.replace("refs/heads/", "");
            }
            if let Ok(commit) = head.peel_to_commit() {
                let time = commit.time();
                last_commit_at = Some(Utc.timestamp_opt(time.seconds(), 0).unwrap());
            }
        }
    }

    let branch_count = repo.branches(Some(BranchType::Local))?.count();
    let tag_count = repo.tag_names(None)?.len();

    // Calculate size (simplified: just sum file sizes in .git dir)
    let size_bytes = calculate_dir_size(&handle.path).unwrap_or(0);

    Ok(RepoInfo {
        owner: handle.owner.clone(),
        name: handle.name.clone(),
        default_branch,
        size_bytes,
        branch_count,
        tag_count,
        last_commit_at,
        is_empty,
    })
}

fn calculate_dir_size(path: &Path) -> std::io::Result<u64> {
    let mut total_size = 0;
    if path.is_dir() {
        for entry in fs::read_dir(path)? {
            let entry = entry?;
            let metadata = entry.metadata()?;
            if metadata.is_dir() {
                total_size += calculate_dir_size(&entry.path())?;
            } else {
                total_size += metadata.len();
            }
        }
    }
    Ok(total_size)
}

pub fn set_description(handle: &RepoHandle, desc: &str) -> Result<(), GitError> {
    let desc_path = handle.path.join("description");
    fs::write(desc_path, desc)?;
    Ok(())
}

pub fn set_default_branch(handle: &RepoHandle, branch: &str) -> Result<(), GitError> {
    // branch name should be full ref or just name? Usually just name like "main"
    // But set_head expects refs/heads/main
    let ref_name = if branch.starts_with("refs/") {
        branch.to_string()
    } else {
        format!("refs/heads/{}", branch)
    };

    handle.inner.set_head(&ref_name)?;
    Ok(())
}

pub fn fork(base_path: &Path, source: &RepoHandle, new_owner: &str, new_name: &str) -> Result<RepoHandle, GitError> {
    validate_name(new_name)?;
    let new_path = get_repo_path(base_path, new_owner, new_name);

    if new_path.exists() {
        return Err(GitError::AlreadyExists(new_name.to_string()));
    }

    if let Some(parent) = new_path.parent() {
        fs::create_dir_all(parent)?;
    }

    // Clone bare
    let _ = git2::build::RepoBuilder::new()
        .bare(true)
        .clone(source.path.to_str().unwrap(), &new_path)?;

    Ok(RepoHandle {
        path: new_path.clone(),
        owner: new_owner.to_string(),
        name: new_name.to_string(),
        inner: Repository::open(&new_path)?,
    })
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_validate_name() {
        assert!(validate_name("valid-name").is_ok());
        assert!(validate_name("invalid/name").is_err());
        assert!(validate_name("..").is_err());
        assert!(validate_name("with\\slash").is_err());
    }
}
