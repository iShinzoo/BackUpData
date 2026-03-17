# BackUpData — Distributed Database Backup System
### Go + gRPC + Worker Pool + Streaming Backups + Prometheus + Docker

![Go](https://img.shields.io/badge/Go-1.22-blue)
![gRPC](https://img.shields.io/badge/gRPC-RPC-green)
![Docker](https://img.shields.io/badge/Docker-enabled-blue)
![Prometheus](https://img.shields.io/badge/Prometheus-metrics-orange)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-database-blue)

A production-style backend system built in Go that demonstrates how a **distributed backup platform** can be architected using modern backend patterns.

This project implements:

- Clean Architecture
- gRPC client/server communication
- Worker pool concurrency
- Streaming database backups
- Pluggable storage adapters (local / S3 ready)
- Slack notifications
- Scheduler (cron-based)
- Dockerized development environment
- Prometheus metrics
- Graceful shutdown
- Context propagation
- Retry + backoff strategy

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
├── Scheduler (cron jobs)
├── Worker Pool
├── Retry + Backoff
├── Metrics
│
▼
Backup Pipeline
(pg_dump → gzip → storage)
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
+---------------+                +----------------+
|  PostgreSQL   |                |   MinIO/S3     |
|     :5432     |                |  Object Store  |
+---------------+                +----------------+
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
| Storage | Local / S3 Adapter |
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
│   │   ├── s3/
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
├── backups/                    # Generated backups
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── go.sum
```

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

Large databases are backed up using **streaming I/O**.

```
PostgreSQL
│
▼
pg_dump
│
▼
gzip compression
│
▼
storage adapter
```

This prevents loading the entire backup into memory.

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

Backups can be automated using cron expressions.

Example:
```
0 2 * * *
```

Runs backups **every day at 2 AM**. Currently it runs every 10 seconds.

---

# Notifications

Slack notifications are sent after backup completion.

Example message:
```
Backup completed
Database: postgres-db
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

Start databases:
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

# Verify Backup File

```bash
ls backups
```

Example:
```
postgres-db.sql.gz
```

Inspect backup:
```bash
gunzip -c backups/postgres-db.sql.gz | head
```

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
+---------------+
| Storage Layer |
+---------------+
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
```

## Concurrency Model

```
          Job Queue
              │
   ┌──────────┼──────────┐
   ▼          ▼          ▼
 Worker1    Worker2    Worker3
```


# Author

**Krishna Thakur**

Backend engineering learning project focused on building production-style infrastructure systems in Go.

## Contact

- X (Twitter): [https://x.com/i_krsna4](https://x.com/i_krsna4)
- LinkedIn: [https://www.linkedin.com/in/krishnathakur1/](https://www.linkedin.com/in/krishnathakur1/)

--- 

⭐ Star this repository if you find it helpful!
