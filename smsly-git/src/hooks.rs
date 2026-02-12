use std::io::{self, BufRead};
use git2::Oid;
use std::process::{exit, Command};

pub fn run_pre_receive() {
    let stdin = io::stdin();
    for line in stdin.lock().lines() {
        let line = match line {
            Ok(l) => l,
            Err(_) => break,
        };

        let parts: Vec<&str> = line.split_whitespace().collect();
        if parts.len() < 3 {
            continue;
        }
        let old_oid_str = parts[0];
        let new_oid_str = parts[1];
        let refname = parts[2];

        // Check if branch is protected
        if refname == "refs/heads/main" || refname == "refs/heads/master" {
            // Deletion check
            if new_oid_str.chars().all(|c| c == '0') {
                eprintln!("Error: Deletion of protected branch '{}' is not allowed.", refname);
                exit(1);
            }

            // Creation is allowed (old_oid is zero)
            if old_oid_str.chars().all(|c| c == '0') {
                continue;
            }

            // Parse OIDs to validate format
            if Oid::from_str(old_oid_str).is_err() {
                 eprintln!("Error: Invalid old OID {}", old_oid_str);
                 exit(1);
            }
            if Oid::from_str(new_oid_str).is_err() {
                 eprintln!("Error: Invalid new OID {}", new_oid_str);
                 exit(1);
            }

            // Check fast-forward using git command (respects quarantine)
            let status = Command::new("git")
                .arg("merge-base")
                .arg("--is-ancestor")
                .arg(old_oid_str)
                .arg(new_oid_str)
                .status();

            match status {
                Ok(s) => {
                    if !s.success() {
                        eprintln!("Error: Non-fast-forward update to protected branch '{}' is rejected.", refname);
                        exit(1);
                    }
                },
                Err(e) => {
                    eprintln!("Error checking commit graph: {}", e);
                    exit(1);
                }
            }
        }
    }
}
