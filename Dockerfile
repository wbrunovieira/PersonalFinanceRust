FROM rust:1.72 as builder

WORKDIR /app

# Cache dependencies
COPY Cargo.toml .
RUN cargo fetch

# Copy source files and build the Rust binary
COPY . .
RUN rm -f Cargo.lock
RUN cargo build --release

# Use bookworm-slim (which has glibc 2.36) as the final runtime
FROM debian:bookworm-slim

# Copy the binary from the builder image
COPY --from=builder /app/target/release/financas-rust /usr/local/bin/financas-rust

# Run the Rust binary
CMD ["financas-rust"]
