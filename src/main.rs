use std::env;
use tokio_postgres::{NoTls};
use actix_web::{get, App, HttpServer, Responder};
use chrono::NaiveDateTime; 

#[get("/")]
async fn index() -> impl Responder {
    format!("Hello, Your Rust app is live! 🦀")
}

#[get("/db")]
async fn test_db() -> impl Responder {
    let database_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set");
    let (client, connection) = tokio_postgres::connect(&database_url, NoTls).await.unwrap();

 
    tokio::spawn(async move {
        if let Err(e) = connection.await {
            eprintln!("connection error: {}", e);
        }
    });

   
    let rows = client.query("SELECT NOW()", &[]).await.unwrap();
    let current_time: NaiveDateTime = rows[0].get(0);
    format!("Database is working! Current time: {:?}", current_time)
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    let port = env::var("PORT").unwrap_or_else(|_| "8080".to_string());
    println!("Starting server at http://localhost:{}", port);
    HttpServer::new(|| App::new().service(index).service(test_db))
        .bind(("0.0.0.0", port.parse::<u16>().unwrap()))?
        .run()
        .await
}
