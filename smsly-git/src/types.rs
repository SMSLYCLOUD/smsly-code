// ============= THESE ARE CANONICAL — DO NOT CHANGE =============

use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use std::path::PathBuf;

/// Handle to an open Git repository
pub struct RepoHandle {
    pub path: PathBuf,
    pub owner: String,
    pub name: String,
    pub(crate) inner: git2::Repository,
}

/// Git commit information
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CommitInfo {
    pub sha: String,
    pub short_sha: String,
    pub message: String,
    pub body: Option<String>,
    pub author: Signature,
    pub committer: Signature,
    pub parents: Vec<String>,
    pub tree_sha: String,
    pub timestamp: DateTime<Utc>,
    pub is_merge: bool,
}

/// Git author/committer signature
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Signature {
    pub name: String,
    pub email: String,
    pub timestamp: DateTime<Utc>,
}

/// Branch information
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BranchInfo {
    pub name: String,
    pub sha: String,
    pub is_default: bool,
    pub is_protected: bool,
}

/// Tag information
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TagInfo {
    pub name: String,
    pub sha: String,
    pub is_annotated: bool,
    pub message: Option<String>,
    pub tagger: Option<Signature>,
    pub target_sha: String,
}

/// File tree entry
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TreeEntry {
    pub name: String,
    pub path: String,
    pub sha: String,
    pub entry_type: EntryType,
    pub size: Option<u64>,
    pub mode: u32,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum EntryType {
    Blob,
    Tree,
    Submodule,
    Symlink,
}

/// Diff result
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DiffResult {
    pub files: Vec<DiffFile>,
    pub stats: DiffStats,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DiffStats {
    pub additions: usize,
    pub deletions: usize,
    pub files_changed: usize,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DiffFile {
    pub old_path: Option<String>,
    pub new_path: Option<String>,
    pub status: DiffStatus,
    pub hunks: Vec<DiffHunk>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum DiffStatus {
    Added,
    Modified,
    Deleted,
    Renamed,
    Copied,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DiffHunk {
    pub old_start: u32,
    pub old_lines: u32,
    pub new_start: u32,
    pub new_lines: u32,
    pub header: String,
    pub lines: Vec<DiffLine>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DiffLine {
    pub line_type: LineType,
    pub content: String,
    pub old_lineno: Option<u32>,
    pub new_lineno: Option<u32>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum LineType {
    Add,
    Delete,
    Context,
}

/// Repository information
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RepoInfo {
    pub owner: String,
    pub name: String,
    pub default_branch: String,
    pub size_bytes: u64,
    pub branch_count: usize,
    pub tag_count: usize,
    pub last_commit_at: Option<DateTime<Utc>>,
    pub is_empty: bool,
}

/// File content
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct FileContent {
    pub content: Vec<u8>,
    pub size: u64,
    pub encoding: String,
    pub is_binary: bool,
    pub mime_type: String,
    pub sha: String,
}

/// Blame line
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BlameLine {
    pub line_number: usize,
    pub content: String,
    pub commit_sha: String,
    pub author: Signature,
    pub date: DateTime<Utc>,
}

/// Hook types
pub enum HookType {
    PreReceive,
    Update,
    PostReceive,
}

/// Ref update (used in hooks)
#[derive(Debug, Clone)]
pub struct RefUpdate {
    pub ref_name: String,
    pub old_sha: String,
    pub new_sha: String,
}
