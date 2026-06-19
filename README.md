# chill-crate-api

REST API for managing file storage organized into groups and buckets. Files are stored in S3-compatible object storage (MinIO or Ceph) with metadata tracked in PostgreSQL.

## Stack

- **Go** 1.26 / **Gin** — HTTP server
- **GORM** + **PostgreSQL 16** — metadata storage
- **AWS SDK v2** / **MinIO** — object storage
- **Keycloak** — authentication (in progress)
- **Viper** — configuration

## Concepts

| Entity | Description |
|--------|-------------|
| **Group** | Top-level namespace (team/project). Has members. |
| **Bucket** | Container for files within a group. Names are unique per group. |
| **Object** | A file upload. Supports versioning — uploading the same filename creates a new version; downloads serve the latest. |

## API

```
GET  /ping

POST /groups                              create group
GET  /groups                              list user's groups
GET  /groups/:groupId/buckets             list buckets in group

POST /buckets                             create bucket
GET  /buckets/:bucketId/objects           list objects in bucket
POST /buckets/:bucketId/objects           upload file (multipart/form-data)
GET  /buckets/:bucketId/objects/:filename download file (latest version)

GET  /objects/:id                         get object metadata
```

## Getting Started

**Prerequisites:** Go 1.26+, Docker, Docker Compose

```bash
# 1. Copy and configure environment
cp .env.example .env

# 2. Start dependencies (PostgreSQL, MinIO, Keycloak)
docker-compose up -d

# 3. Run the server
go run cmd/server/main.go
# Listening on http://localhost:8081

# Health check
curl http://localhost:8081/ping
```

## Configuration

Copy `.env.example` to `.env` and set:

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_HOST` | `localhost` | Bind address |
| `SERVER_PORT` | `8081` | Bind port |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | — | PostgreSQL user |
| `DB_PASSWORD` | — | PostgreSQL password |
| `DB_NAME` | — | Database name |
| `DB_SSLMODE` | `disable` | PostgreSQL SSL mode |
| `STORAGE_ENDPOINT` | — | S3/MinIO endpoint URL |
| `STORAGE_REGION` | — | S3 region |
| `STORAGE_ACCESS_KEY` | — | S3 access key |
| `STORAGE_SECRET_KEY` | — | S3 secret key |

Docker Compose services are available at:
- PostgreSQL: `localhost:5432`
- MinIO API: `localhost:9000`, console: `localhost:9001`
- Keycloak: `localhost:8080`

## License

Apache 2.0 — see [LICENSE](LICENSE).
