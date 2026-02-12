use actix_web::{get, post, web, App, HttpServer, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use dotenv::dotenv;
use std::env;
use log::info;

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
    // TODO: Implement libgit2 logic using git2 crate
    // let repo = git2::Repository::init(&req.name);

    HttpResponse::Ok().json(serde_json::json!({
        "message": format!("Repository {} created (mock)", req.name)
    }))
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    dotenv().ok();
    env_logger::init_from_env(env_logger::Env::new().default_filter_or("info"));

    let host = env::var("GIT_HOST").unwrap_or_else(|_| "0.0.0.0".to_string());
    let port = env::var("GIT_PORT").unwrap_or_else(|_| "8081".to_string());
    let addr = format!("{}:{}", host, port);

    info!("Starting SMSLY Git Engine on {}", addr);

    HttpServer::new(|| {
        App::new()
            .service(health_check)
            .service(create_repo)
    })
    .bind(addr)?
    .run()
    .await
}
