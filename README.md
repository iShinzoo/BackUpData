# BackUpData - Distributed Database Backup System
### Go + gRPC + Worker Pool + Streaming Backups + Prometheus + Docker

![Go](https://img.shields.io/badge/Go-1.22-blue)
![gRPC](https://img.shields.io/badge/gRPC-RPC-green)
![Docker](https://img.shields.io/badge/Docker-enabled-blue)
![Prometheus](https://img.shields.io/badge/Prometheus-metrics-orange)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-database-blue)
![MinIO](https://img.shields.io/badge/MinIO-S3%20Storage-red)

A production-style backend system built in Go that demonstrates how a **distributed backup platform** can be architected using modern backend patterns.

This project implements:

- Clean Architecture
- gRPC client/server communication
- Worker pool concurrency
- Streaming database backups → MinIO (S3-compatible object storage)
- Pluggable storage adapters (local / S3 / GCS)
- Slack notifications
- Scheduler (cron-based, triggers real gRPC backup jobs)
- Dockerized development environment
- Prometheus metrics
- Graceful shutdown
- Context propagation
- Retry + backoff strategy
- Timestamp-based backup naming

The goal of this project is to simulate how **real infrastructure backup tools** are designed and implemented.

---

# Architecture Overview
```
backup-cli
│
▼
gRPC Client
│
▼
backup-daemon (gRPC Server)
│
├── Scheduler (cron jobs → triggers real gRPC RunBackup calls)
├── Worker Pool
├── Retry + Backoff
├── Metrics
│
▼
Backup Pipeline
(pg_dump → gzip → S3/MinIO)
│
▼
Notifications (Slack)
```

---

# Containerized Architecture
```
           +--------------------+
           |     backup-cli     |
           |   gRPC Client      |
           +----------+---------+
                      |
                      v
           +--------------------+
           |    backup-daemon   |
           |   gRPC Server      |
           |      :50051        |
           | Metrics :9090      |
           +----------+---------+
                      |
    +-----------------+------------------+
    |                                    |
    v                                    v
+---------------+                +--------------------+
|  PostgreSQL   |                |   MinIO            |
|     :5432     |                | S3-Compatible Store|
+---------------+                |  API  :9000        |
                                 |  Console :9001     |
                                 +--------------------+
```

All services run locally through **Docker Compose** to simulate a production-like environment.

---

# Tech Stack

| Layer | Technology |
|------|-------------|
| Language | Go |
| RPC Transport | gRPC |
| Concurrency | Goroutines + Channels |
| Scheduler | robfig/cron |
| Compression | gzip |
| Storage | MinIO (S3-compatible) / Local / GCS Adapter |
| Notifications | Slack Webhook |
| Metrics | Prometheus |
| Logging | Zap |
| Containerization | Docker + Docker Compose |

---

# Project Structure
```
BackUpData/
│
├── cmd/
│   ├── backup-cli/           # CLI client
│   │   ├── cmd/
│   │   │   ├── backup.go
│   │   │   ├── root.go
│   │   │   └── schedule.go
│   │   └── main.go
│   │
│   └── backup-daemon/         # gRPC backup server
│       └── main.go
│
├── internal/
│   ├── compression/
│   │   └── gzip.go
│   │
│   ├── core/
│   │   ├── worker/
│   │   ├── backup_service.go
│   │   ├── full_backup.go
│   │   ├── interfaces.go
│   │   ├── job.go
│   │   ├── job_handler.go
│   │   └── strategy.go
│   │
│   ├── db/
│   │   └── postgres/
│   │       ├── backup.go
│   │       ├── executor.go
│   │       └── postgres.go
│   │
│   ├── storage/
│   │   ├── local/
│   │   ├── s3/           # Active: used by MinIO adapter
│   │   └── gcs/
│   │
│   ├── notification/
│   │   └── slack/
│   │       └── slack.go
│   │
│   ├── scheduler/
│   │   └── scheduler.go
│   │
│   └── metrics/
│       └── metrics.go
│
├── pkg/
│   ├── config/
│   └── logger/
│
├── proto/
│   ├── backup.proto
│   ├── backup.pb.go
│   └── backup_grpc.pb.go
│
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── go.sum
```

> **Note:** The `backups/` local directory is no longer the primary destination. All backups are streamed directly to **MinIO** using the S3 adapter. Local storage remains available as a fallback adapter.

---

# Core System Concepts

## Clean Architecture

Business logic is isolated from external systems.
```
CLI / gRPC Layer
│
▼
Application Layer
│
▼
Core Business Logic
│
▼
Adapters (DB / Storage / Notifications)
```

This makes the system:
- extensible
- testable
- loosely coupled

---

## Worker Pool Concurrency

Multiple backups can run concurrently.
```
          +-------------+
          | Job Channel |
          +-------------+
           │    │    │
           ▼    ▼    ▼
           W1   W2   W3
```

Each worker processes backup jobs asynchronously.

---

## Streaming Backup Pipeline

Large databases are backed up using **streaming I/O** — no intermediate files, no full in-memory loads.
```
PostgreSQL
│
▼
pg_dump (streaming stdout)
│
▼
gzip compression (streaming)
│
▼
MinIO / S3 storage adapter (multipart upload)
```

Backup objects are named using a **timestamp-based convention** to avoid collisions and enable chronological sorting:
```
postgres-db_2024-01-15T02-00-00Z.sql.gz
```

---

# gRPC Communication

The CLI communicates with the daemon using **gRPC**.
```
backup-cli
│
▼
RunBackup RPC
│
▼
backup-daemon
```

The daemon executes the backup pipeline and returns a response.

---

# Scheduler

Backups are automated using cron expressions. The scheduler **directly invokes the gRPC `RunBackup` RPC** internally — it does not use shell commands or indirect triggers.

Example cron expression:
```
0 2 * * *
```

Runs a real backup job **every day at 2 AM** by calling `RunBackup` through the gRPC layer. During development, the interval is set to **every 10 seconds** for rapid iteration.

---

# MinIO Setup

MinIO runs as a containerized S3-compatible object store. It is started automatically via Docker Compose alongside PostgreSQL.

## Access the MinIO Console
```
http://localhost:9001
```

Default credentials (development only):
```
Username: minio
Password: minio123
```

## Create a Backup Bucket

Using the MinIO console:

1. Open `http://localhost:9001` in your browser
2. Navigate to **Buckets → Create Bucket**
3. Name the bucket: `backups`
4. Click **Create Bucket**

Using the MinIO CLI (`mc`):
```bash
# Configure the local MinIO alias
mc alias set local http://localhost:9000 minioadmin minioadmin

# Create the bucket
mc mb local/backups

# Verify
mc ls local/
```

## S3 Adapter Configuration

The daemon uses the following environment variables to connect to MinIO:
```env
S3_ENDPOINT=http://localhost:9000
S3_BUCKET=backup-bucket
S3_ACCESS_KEY=minio123
S3_SECRET_KEY=minio
S3_USE_PATH_STYLE=true
```

> `S3_USE_PATH_STYLE=true` is required for MinIO — it uses path-style URLs instead of virtual-hosted style.

---

# Notifications

Slack notifications are sent after backup completion.

Example message:
```
Backup completed
Database: postgres-db
File: postgres-db_2024-01-15T02-00-00Z.sql.gz
Storage: MinIO (s3://backups/)
Duration: 3s
Size: 12000 bytes
```

---

# Observability

## Prometheus Metrics

The daemon exposes metrics at:
```
http://localhost:9090/metrics
```

Available metrics:
```
backup_success_total
backup_failure_total
backup_duration_seconds
```

These metrics allow integration with **Grafana dashboards and alerts**.

---

# Running the System

The system can be run in **two ways**:
1. Local development
2. Docker deployment

---

# Local Development Setup

## Requirements

- Go 1.22+
- Docker
- Docker Compose
- protoc
- protoc-gen-go
- protoc-gen-go-grpc

---

# Start Infrastructure

Start all services (PostgreSQL + MinIO):
```bash
docker compose up -d
```

Verify containers:
```bash
docker ps
```

Expected services:
```
backup-postgres
backup-mysql
backup-minio
```

---

# Insert Sample Data

Connect to the database container:
```bash
docker exec -it backup-postgres psql -U backup -d testdb
```

Create sample table:
```sql
CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    name TEXT
);
```

Insert sample data:
```sql
INSERT INTO users(name) VALUES ('alice'), ('bob');
```

Verify:
```sql
SELECT * FROM users;
```

---

# Create MinIO Bucket

Before running the daemon, ensure the `backups` bucket exists in MinIO:
```bash
mc alias set local http://localhost:9000 minioadmin minioadmin
mc mb local/backups
```

Or use the console at `http://localhost:9001`.

---

# Start Backup Daemon
```bash
go run ./cmd/backup-daemon
```

Daemon runs on:
- gRPC: `:50051`
- Metrics: `:9090`

---

# Trigger Backup

In another terminal:
```bash
go run ./cmd/backup-cli backup
```

Expected output:
```
Daemon Response: backup completed
```

---

# Verify Backup in MinIO

List backup objects in the MinIO bucket:
```bash
mc ls local/backups
```

Example output:
```
[2024-01-15 02:00:03]   12345  postgres-db_2024-01-15T02-00-00Z.sql.gz
```

Inspect backup content:
```bash
mc cat local/backups/postgres-db_2024-01-15T02-00-00Z.sql.gz | gunzip | head
```

Or via the MinIO console at `http://localhost:9001` → **Object Browser → backups**.

---

# Verify Metrics

Open browser:
```
http://localhost:9090/metrics
```

Example output:
```
backup_success_total 1
backup_failure_total 0
```

---

# Docker Deployment

Run the full system with:
```bash
docker compose up --build
```

Stop the system:
```bash
docker compose down
```

---

# UML Diagrams

## Component Diagram
```
+----------------+
|   backup-cli   |
+----------------+
        |
        v
+----------------+
| backup-daemon  |
|   gRPC Server  |
+----------------+
     |       |
     v       v
+---------+  +--------------+
| Worker  |  | PostgreSQL   |
+---------+  +--------------+
     |
     v
+------------------------+
| Storage Layer          |
| (S3 Adapter → MinIO)   |
+------------------------+
     |
     v
+--------------+
| Notifications|
+--------------+
```

## Class Diagram
```
+---------------------+
| BackupService       |
|---------------------|
| RunBackup()         |
| StreamProgress()    |
+---------------------+
          |
          v
+----------------------+
| BackupExecutor       |
|----------------------|
| Run()                |
+----------------------+
          |
          v
+----------------------+
| PostgresExecutor     |
+----------------------+
          |
          v
+----------------------+
| S3StorageAdapter     |
| (MinIO-compatible)   |
+----------------------+
```

## Concurrency Model
```
          Job Queue
              │
   ┌──────────┼──────────┐
   ▼          ▼          ▼
 Worker1    Worker2    Worker3
   │          │          │
   └──────────┴──────────┘
              │
     Streaming → MinIO
```

---

# Backup Naming Convention

All backup objects follow a **timestamp-based naming scheme**:
```
{database-name}_{ISO8601-timestamp}.sql.gz
```

Example:
```
postgres-db_2024-01-15T02-00-00Z.sql.gz
```

This ensures:
- No filename collisions across scheduled runs
- Chronological sorting in object storage
- Easy retention policy management by date prefix

---

# Author

**Krishna Thakur**

Backend engineering learning project focused on building production-style infrastructure systems in Go.

## Contact

- X (Twitter): [https://x.com/i_krsna4](https://x.com/i_krsna4)
- LinkedIn: [https://www.linkedin.com/in/krishnathakur1/](https://www.linkedin.com/in/krishnathakur1/)

---

⭐ Star this repository if you find it helpful!
