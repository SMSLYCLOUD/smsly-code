use actix_web::{get, post, web, App, HttpServer, HttpResponse, Responder, middleware};
use serde::{Deserialize, Serialize};
use dotenv::dotenv;
use log::{info, error};

use smsly_git::config::Config;
use smsly_git::repo;
use smsly_git::transport::http;
use smsly_git::error::GitError;

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
    owner: String,
    name: String,
}

#[post("/repo")]
async fn create_repo(req: web::Json<CreateRepoRequest>, config: web::Data<Config>) -> impl Responder {
    info!("Creating repo: {}/{}", req.owner, req.name);

    match repo::init_bare(&config.data_path, &req.owner, &req.name) {
        Ok(handle) => {
            info!("Successfully created bare repository at {:?}", handle.path);
            HttpResponse::Ok().json(serde_json::json!({
                "message": format!("Repository {}/{} created", req.owner, req.name),
                "path": handle.path.to_string_lossy()
            }))
        },
        Err(e) => {
            error!("Failed to create repository: {}", e);
            let status = match e {
                GitError::AlreadyExists(_) => 409,
                GitError::InvalidName(_) => 400,
                GitError::NotFound(_) => 404,
                GitError::PermissionDenied(_) => 403,
                _ => 500,
            };
            HttpResponse::build(actix_web::http::StatusCode::from_u16(status).unwrap())
                .json(serde_json::json!({
                    "error": e.to_string()
                }))
        }
    }
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    dotenv().ok();
    env_logger::init_from_env(env_logger::Env::new().default_filter_or("info"));

    let config = Config::from_env();
    let addr = format!("{}:{}", config.host, config.port);

    std::fs::create_dir_all(&config.data_path)?;

    info!("Starting SMSLY Git Engine on {}", addr);
    info!("Storing repositories in {:?}", config.data_path);

    let config_data = web::Data::new(config.clone());

    HttpServer::new(move || {
        App::new()
            .app_data(config_data.clone())
            .wrap(middleware::Logger::default())
            .service(health_check)
            .service(create_repo)
            // Register Git Smart HTTP routes
            .service(http::info_refs)
            .service(http::service_rpc)
    })
    .bind(addr)?
    .run()
    .await
}
