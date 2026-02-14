use actix_web::{get, post, web, HttpResponse, Responder, HttpRequest};
use crate::repo;
use crate::config::Config;
use std::process::{Command, Stdio};
use std::io::Write;
use crate::error::GitError;

fn check_repo(config: &Config, owner: &str, repo_name: &str) -> Result<std::path::PathBuf, GitError> {
    let name = repo_name.trim_end_matches(".git");
    let handle = repo::open(&config.data_path, owner, name)?;
    Ok(handle.path)
}

#[get("/{owner}/{repo}/info/refs")]
pub async fn info_refs(
    req: HttpRequest,
    path: web::Path<(String, String)>,
    config: web::Data<Config>,
) -> impl Responder {
    let (owner, repo_name) = path.into_inner();
    let query = req.query_string();
    // Simple parsing, assuming service is the first or only param
    let service = if query.contains("service=git-upload-pack") {
        "git-upload-pack"
    } else if query.contains("service=git-receive-pack") {
        "git-receive-pack"
    } else {
        return HttpResponse::Forbidden().body("Only git-upload-pack and git-receive-pack are supported");
    };

    let repo_path = match check_repo(&config, &owner, &repo_name) {
        Ok(p) => p,
        Err(_) => return HttpResponse::NotFound().body("Repository not found"),
    };

    let subcommand = service.trim_start_matches("git-");
    let output = Command::new("git")
        .arg(subcommand)
        .arg("--stateless-rpc")
        .arg("--advertise-refs")
        .arg(&repo_path)
        .output();

    match output {
        Ok(out) => {
            let mut body = Vec::new();
            let packet = format!("# service={}\n", service);
            let len = packet.len() + 4;
            // Write length prefix (hex, 4 chars)
            let _ = write!(body, "{:04x}", len);
            body.write_all(packet.as_bytes()).unwrap();
            body.write_all(b"0000").unwrap(); // Flush packet
            body.write_all(&out.stdout).unwrap();

            HttpResponse::Ok()
                .content_type(format!("application/x-{}-advertisement", service))
                .append_header(("Cache-Control", "no-cache"))
                .body(body)
        },
        Err(e) => {
            HttpResponse::InternalServerError().body(format!("Git error: {}", e))
        }
    }
}

#[post("/{owner}/{repo}/{service}")]
pub async fn service_rpc(
    path: web::Path<(String, String, String)>,
    body: web::Bytes,
    config: web::Data<Config>,
) -> impl Responder {
    let (owner, repo_name, service) = path.into_inner();

    if service != "git-upload-pack" && service != "git-receive-pack" {
        return HttpResponse::Forbidden().body("Invalid service");
    }

    let repo_path = match check_repo(&config, &owner, &repo_name) {
        Ok(p) => p,
        Err(_) => return HttpResponse::NotFound().body("Repository not found"),
    };

    // Note: If body is gzipped, this might fail if git doesn't handle it or we don't decompress.
    // Standard git-http-backend handles it?
    // Usually the web server (Apache/Nginx) handles decompression or passes it through.
    // If git receives compressed data on stdin but expects raw, it breaks.
    // But git-receive-pack expects the pack stream.
    // We will pass body as is.

    let subcommand = service.trim_start_matches("git-");
    let child = Command::new("git")
        .arg(subcommand)
        .arg("--stateless-rpc")
        .arg(&repo_path)
        .stdin(Stdio::piped())
        .stdout(Stdio::piped())
        .spawn();

    match child {
        Ok(mut child) => {
            if let Some(mut stdin) = child.stdin.take() {
                let _ = stdin.write_all(&body);
            }
            let output = child.wait_with_output();
            match output {
                 Ok(out) => {
                     HttpResponse::Ok()
                        .content_type(format!("application/x-{}-result", service))
                        .append_header(("Cache-Control", "no-cache"))
                        .body(out.stdout)
                 },
                 Err(e) => HttpResponse::InternalServerError().body(format!("Process error: {}", e))
            }
        },
        Err(e) => HttpResponse::InternalServerError().body(format!("Spawn error: {}", e))
    }
}
