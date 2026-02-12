use actix_web::{web, HttpRequest, HttpResponse, Responder};
use futures::StreamExt;
use tokio::process::Command;
use std::process::Stdio;
use std::env;
use std::fs;
use std::path::PathBuf;
use log::{error, warn};
use tokio::io::AsyncWriteExt;
use jsonwebtoken::{decode, DecodingKey, Validation, Algorithm};
use serde::{Deserialize, Serialize};
use base64::{Engine as _, engine::general_purpose};

#[derive(Debug, Serialize, Deserialize)]
struct Claims {
    // We make fields optional or check if they exist to be robust
    user_id: Option<u64>,
    exp: Option<usize>,
}

fn check_auth(req: &HttpRequest) -> bool {
    let auth_header = match req.headers().get("Authorization") {
        Some(h) => h,
        None => return false,
    };

    let auth_str = match auth_header.to_str() {
        Ok(s) => s,
        Err(_) => return false,
    };

    if !auth_str.starts_with("Basic ") {
        return false;
    }

    let token = auth_str.trim_start_matches("Basic ");
    let decoded = match general_purpose::STANDARD.decode(token) {
        Ok(d) => d,
        Err(_) => return false,
    };

    let creds = match String::from_utf8(decoded) {
        Ok(s) => s,
        Err(_) => return false,
    };

    let parts: Vec<&str> = creds.splitn(2, ':').collect();
    if parts.len() != 2 {
        return false;
    }

    let password = parts[1];

    if password.is_empty() {
        return false;
    }

    let secret = env::var("JWT_SECRET").unwrap_or_else(|_| "secret".to_string());

    // Validate JWT
    let mut validation = Validation::new(Algorithm::HS256);
    validation.validate_exp = false; // We can enable this if we want strict checking

    // For MVP, if we can decode it with the secret, it's valid.
    match decode::<Claims>(password, &DecodingKey::from_secret(secret.as_bytes()), &validation) {
        Ok(_) => true,
        Err(e) => {
            warn!("JWT validation failed: {}", e);
            false
        }
    }
}

pub async fn get_info_refs(req: HttpRequest, path: web::Path<String>) -> impl Responder {
    let name = path.into_inner();

    if !check_auth(&req) {
        return HttpResponse::Unauthorized()
            .append_header(("WWW-Authenticate", "Basic realm=\"SMSLY Code\""))
            .body("Unauthorized");
    }

    let query_str = req.query_string();
    let service = if query_str.contains("service=git-upload-pack") {
        "git-upload-pack"
    } else if query_str.contains("service=git-receive-pack") {
        "git-receive-pack"
    } else {
        return HttpResponse::BadRequest().body("Invalid service or missing service parameter");
    };

    let base_path = env::var("GIT_DATA_PATH").unwrap_or_else(|_| "./git-data".to_string());
    let repo_path = PathBuf::from(base_path).join(&name);

    if !repo_path.exists() {
        return HttpResponse::NotFound().body("Repository not found");
    }

    let repo_path = fs::canonicalize(&repo_path).unwrap_or(repo_path);

    // Run git <service> --stateless-rpc --advertise-refs .
    // The service name in command usually doesn't have "git-" prefix if we use "upload-pack" or "receive-pack" subcommand.
    let subcommand = service.trim_start_matches("git-");

    let output = Command::new("git")
        .arg(subcommand)
        .arg("--stateless-rpc")
        .arg("--advertise-refs")
        .arg(&repo_path)
        .output()
        .await;

    match output {
        Ok(output) => {
            let mut body = Vec::new();
            // Packet line for service
            let service_line = format!("# service={}\n", service);
            let len = service_line.len() + 4;
            let header = format!("{:04x}{}", len, service_line);
            body.extend_from_slice(header.as_bytes());
            body.extend_from_slice(b"0000"); // Flush packet
            body.extend_from_slice(&output.stdout);

            HttpResponse::Ok()
                .content_type(format!("application/x-{}-advertisement", service))
                .append_header(("Cache-Control", "no-cache"))
                .body(body)
        },
        Err(e) => {
            error!("Failed to execute git: {}", e);
            HttpResponse::InternalServerError().body("Internal Server Error")
        }
    }
}

pub async fn git_service(req: HttpRequest, path: web::Path<String>, mut payload: web::Payload) -> impl Responder {
    let name = path.into_inner();
    let path_str = req.uri().path();

    let service_name = if path_str.ends_with("/git-upload-pack") {
        "git-upload-pack"
    } else if path_str.ends_with("/git-receive-pack") {
        "git-receive-pack"
    } else {
        return HttpResponse::BadRequest().body("Invalid service endpoint");
    };

    if !check_auth(&req) {
        return HttpResponse::Unauthorized()
            .append_header(("WWW-Authenticate", "Basic realm=\"SMSLY Code\""))
            .body("Unauthorized");
    }

    let base_path = env::var("GIT_DATA_PATH").unwrap_or_else(|_| "./git-data".to_string());
    let repo_path = PathBuf::from(base_path).join(&name);

    if !repo_path.exists() {
        return HttpResponse::NotFound().body("Repository not found");
    }

    let repo_path = fs::canonicalize(&repo_path).unwrap_or(repo_path);

    let subcommand = service_name.trim_start_matches("git-");

    let mut child = match Command::new("git")
        .current_dir(&repo_path)
        .arg("-c")
        .arg("core.hooksPath=hooks")
        .arg(subcommand)
        .arg("--stateless-rpc")
        .arg(".")
        .stdin(Stdio::piped())
        .stdout(Stdio::piped())
        .spawn() {
            Ok(c) => c,
            Err(e) => {
                error!("Failed to spawn git: {}", e);
                return HttpResponse::InternalServerError().body("Internal Server Error");
            }
        };

    let mut stdin = child.stdin.take().expect("Failed to open stdin");

    // Write payload to stdin
    while let Some(chunk) = payload.next().await {
        match chunk {
            Ok(bytes) => {
                if let Err(e) = stdin.write_all(&bytes).await {
                    error!("Failed to write to git stdin: {}", e);
                    // We can't return error easily here because we are in the middle of processing
                    // But we can break and let git fail
                    break;
                }
            }
            Err(e) => {
                error!("Payload error: {}", e);
                break;
            }
        }
    }
    // stdin is dropped here, closing it
    drop(stdin);

    let output = match child.wait_with_output().await {
        Ok(o) => o,
        Err(e) => {
             error!("Failed to wait on git: {}", e);
             return HttpResponse::InternalServerError().body("Internal Server Error");
        }
    };

    HttpResponse::Ok()
        .content_type(format!("application/x-{}-result", service_name))
        .append_header(("Cache-Control", "no-cache"))
        .body(output.stdout)
}
