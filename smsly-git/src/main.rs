use actix_web::{get, post, web, App, HttpServer, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use dotenv::dotenv;
use std::env;
use std::path::PathBuf;
use log::{info, error};
use git2::{Repository, ObjectType};

#[derive(Serialize)]
struct Status {
    status: String,
    version: String,
}

#[get("/health")]
async fn health_check() -> impl Responder {
    HttpResponse::Ok().json(Status {
        status: "ok".to_string(),
        version: env!("CARGO_PKG_VERSION").to_string(),
    })
}

#[derive(Deserialize)]
struct CreateRepoRequest {
    name: String,
}

#[post("/repo")]
async fn create_repo(req: web::Json<CreateRepoRequest>) -> impl Responder {
    info!("Creating repo: {}", req.name);

    let base_path = env::var("GIT_DATA_PATH").unwrap_or_else(|_| "./git-data".to_string());
    let mut repo_path = PathBuf::from(base_path);

    if req.name.contains("..") || req.name.contains('/') || req.name.contains('\\') {
         return HttpResponse::BadRequest().json(serde_json::json!({
            "error": "Invalid repository name"
        }));
    }

    repo_path.push(&req.name);

    match Repository::init_bare(&repo_path) {
        Ok(_) => {
            info!("Successfully created bare repository at {:?}", repo_path);
            HttpResponse::Ok().json(serde_json::json!({
                "message": format!("Repository {} created", req.name),
                "path": repo_path.to_string_lossy()
            }))
        },
        Err(e) => {
            error!("Failed to create repository: {}", e);
            HttpResponse::InternalServerError().json(serde_json::json!({
                "error": format!("Failed to create repository: {}", e)
            }))
        }
    }
}

#[derive(Serialize)]
struct Commit {
    id: String,
    message: String,
    author: String,
    date: String,
    mip_verified: bool,
}

#[get("/repo/{name}/commits")]
async fn get_commits(name: web::Path<String>) -> impl Responder {
    let base_path = env::var("GIT_DATA_PATH").unwrap_or_else(|_| "./git-data".to_string());
    let repo_path = PathBuf::from(base_path).join(&*name);

    let repo = match Repository::open(&repo_path) {
        Ok(r) => r,
        Err(_) => return HttpResponse::NotFound().json(serde_json::json!({"error": "Repository not found"})),
    };

    let mut revwalk = match repo.revwalk() {
        Ok(r) => r,
        Err(_) => return HttpResponse::InternalServerError().json(serde_json::json!({"error": "Failed to init revwalk"})),
    };

    if revwalk.push_head().is_err() {
        return HttpResponse::Ok().json(Vec::<Commit>::new());
    }

    let commits: Vec<Commit> = revwalk
        .filter_map(|id| id.ok())
        .filter_map(|id| repo.find_commit(id).ok())
        .take(20)
        .map(|commit| {
            // Check for MIP signature
            // Real implementation would verify cryptographic signature
            // For now, we mock it: if commit message contains "[MIP]", it's verified.
            let message = commit.message().unwrap_or("").to_string();
            let mip_verified = message.contains("[MIP]");

            Commit {
                id: commit.id().to_string(),
                message,
                author: commit.author().name().unwrap_or("").to_string(),
                date: commit.time().seconds().to_string(),
                mip_verified,
            }
        })
        .collect();

    HttpResponse::Ok().json(commits)
}

#[derive(Serialize)]
struct TreeEntry {
    name: String,
    kind: String,
    id: String,
}

#[get("/repo/{name}/tree/{ref_name}/{path:.*}")]
async fn get_tree(path: web::Path<(String, String, String)>) -> impl Responder {
    let (name, ref_name, subpath) = path.into_inner();
    let base_path = env::var("GIT_DATA_PATH").unwrap_or_else(|_| "./git-data".to_string());
    let repo_path = PathBuf::from(base_path).join(&name);

    let repo = match Repository::open(&repo_path) {
        Ok(r) => r,
        Err(_) => return HttpResponse::NotFound().json(serde_json::json!({"error": "Repository not found"})),
    };

    let obj = match repo.revparse_single(&ref_name) {
        Ok(o) => o,
        Err(_) => return HttpResponse::NotFound().json(serde_json::json!({"error": "Ref not found"})),
    };

    let commit = match obj.peel_to_commit() {
        Ok(c) => c,
        Err(_) => return HttpResponse::InternalServerError().json(serde_json::json!({"error": "Not a commit"})),
    };

    let tree = match commit.tree() {
        Ok(t) => t,
        Err(_) => return HttpResponse::InternalServerError().json(serde_json::json!({"error": "Failed to get tree"})),
    };

    let target_tree = if !subpath.is_empty() && subpath != "/" {
         match tree.get_path(std::path::Path::new(&subpath)) {
            Ok(entry) => match entry.to_object(&repo) {
                 Ok(obj) => {
                     match obj.into_tree() {
                        Ok(t) => t,
                        Err(_) => return HttpResponse::BadRequest().json(serde_json::json!({"error": "Path is not a directory"})),
                     }
                 },
                 Err(_) => return HttpResponse::InternalServerError().json(serde_json::json!({"error": "Failed to get object"})),
            },
            Err(_) => return HttpResponse::NotFound().json(serde_json::json!({"error": "Path not found"})),
         }
    } else {
        tree
    };

    let mut entries: Vec<TreeEntry> = Vec::new();
    for entry in target_tree.iter() {
        let kind = match entry.kind() {
            Some(ObjectType::Blob) => "blob",
            Some(ObjectType::Tree) => "tree",
            _ => "unknown",
        };
        entries.push(TreeEntry {
            name: entry.name().unwrap_or("").to_string(),
            kind: kind.to_string(),
            id: entry.id().to_string(),
        });
    }

    HttpResponse::Ok().json(entries)
}


#[actix_web::main]
async fn main() -> std::io::Result<()> {
    dotenv().ok();
    env_logger::init_from_env(env_logger::Env::new().default_filter_or("info"));

    let host = env::var("GIT_HOST").unwrap_or_else(|_| "0.0.0.0".to_string());
    let port = env::var("GIT_PORT").unwrap_or_else(|_| "8081".to_string());
    let addr = format!("{}:{}", host, port);

    let base_path = env::var("GIT_DATA_PATH").unwrap_or_else(|_| "./git-data".to_string());
    std::fs::create_dir_all(&base_path)?;

    info!("Starting SMSLY Git Engine on {}", addr);
    info!("Storing repositories in {}", base_path);

    HttpServer::new(|| {
        App::new()
            .service(health_check)
            .service(create_repo)
            .service(get_commits)
            .service(get_tree)
    })
    .bind(addr)?
    .run()
    .await
}
