# Go Audit Partition Purger

This project is an example of a Terabyte-scale Audit Log management system. It focuses on Data Lifecycle Management using PostgreSQL Declarative Partitioning alongside the `pg_partman` extension. This prevents database bloat and severe performance degradation over time. The codebase follows a Clean Architecture approach and exposes a REST API built with the Gin Framework, fully documented via Swagger.

---

## How to Run and Test the API

### 1. Spin Up the Database

I use a custom PostgreSQL 16 image that has the `pg_partman` extension pre-installed. Run the following command via Docker to start the infrastructure:

```bash
docker-compose up -d --build
```

### 1.5. Prepare the Environment Variables

This system relies on environment variables for local database connections. Copy the example file and rename it to `.env`:

```bash
cp .env.example .env
```

### 2. Run Database Migrations

Load the migration script into the database. This creates the parent tables and configures the automated partition rules:

```bash
docker exec -i audit_db psql -U root -d audit_logs_db < ./migrations/000001_init_audit_log.up.sql
```

### 3. Start the Go Application (REST API Server)

Once the database is fully initialized, install the Go dependencies and start the web server:

```bash
go mod tidy
go run main.go
```

**Result**: You should see server logs indicating that the application has successfully started at `http://localhost:8080`.

### 4. Test via Swagger UI

Open your browser and navigate to:
**[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**

From the API documentation page, you can easily test the endpoint:

1. Click on the `POST /api/v1/audit-logs` endpoint.
2. Click `Try it out`.
3. Modify the JSON payload to anything you want. For example:
   ```json
   {
     "action": "UPLOAD_FILE",
     "details": "{\"file_size\": \"5MB\"}",
     "user_id": "USR-999"
   }
   ```
4. Click `Execute` and verify that the response returns a 201 Created status.

_(Note: If you accidentally send a "details" field that isn't a valid JSON structure, the API will reject it with a 400 Bad Request. I intentionally placed validation layers across the app to prevent runtime panics at the database level.)_

### 5. Inspect the Database Directly

**Method 1: Open instantly with TablePlus (Mac only)**
If you are using macOS and have the TablePlus app installed, you can simply run this command in your terminal. It will open TablePlus and connect to the database immediately.

```bash
open "postgresql://root:password@localhost:5432/audit_logs_db"
```

**Method 2: Use other Database Clients (DBeaver, DataGrip, pgAdmin)**
Enter the strict connection details below into your preferred database client:

- **Host**: `127.0.0.1` or `localhost`
- **Port**: `5432`
- **User**: `root`
- **Password**: `password`
- **Database**: `audit_logs_db`

**Method 3: Connect via Terminal Shell**
For those who prefer the standard command line:

```bash
docker exec -it audit_db psql -U root -d audit_logs_db
```

Once inside the shell, type `\dt` to list all the auto-generated child partitions, or type `SELECT * FROM audit_logs;` to see the data you just sent via Swagger.

---

## The Theory: Why DROP PARTITION is Better Than DELETE

When a table holds massive amounts of data, running a traditional command like `DELETE FROM audit_logs WHERE ...` causes severe negative impacts:

1. **Database Bloat**: PostgreSQL does not immediately return deleted file space to the operating system. Instead, it creates dead tuples and forces you to rely efficiently on VACUUM processes, which consumes enormous CPU resources.
2. **Locking and Performance Drops**: Querying a single massive table degrades index performance and slows down the entire system.

**The Solution: Partitioning + Drop**
Instead of storing everything in one giant bucket, I split the data sequentially into monthly buckets. When it's time to delete old data according to the retention policy, I simply decouple the bucket and throw it away entirely.

```sql
-- For example, to drop the January 2026 partition:
ALTER TABLE audit_logs DETACH PARTITION audit_logs_p20260101;
DROP TABLE audit_logs_p20260101;
```

This method reclaims terabytes of physical storage instantaneously. In this project, the `pg_partman` extension silently handles the creation and destruction of these data buckets in the background so you never have to do it manually.

---

## Lessons Learned from AI Code Review (CodeRabbit)

I intentionally planted a few "security flaws" in the Go codebase to test the accuracy and usefulness of my AI Code Reviewer. Here is exactly what the AI caught and taught me during the initial Pull Request:

1. **SQL Injection (CWE-89)**: It found a critical vulnerability where raw input was being injected directly into an `ALTER TABLE %s` formatting query. The AI instructed me how to correctly wrap PostgreSQL identifiers with double quotes and fix the DDL syntax order.
2. **Missing Input Validation**: It noticed I was trusting raw strings for formatted dates. It recommended robust Regex validation at the usecase layer to prevent users from executing a purge command on malformed partition names.
3. **Data Integrity Mismatch (jsonb)**: This review prevented severe runtime panics. My database schema utilized native `jsonb`, but the Go struct used a raw string format. If non-JSON text successfully hit the database driver, the app would crash. I solved this by adding `json.Valid()` checks right at the Gin API boundary.

---

## The Next Step: Partitioning vs. Sharding

While this project is a showcasing PostgreSQL Declarative Partitioning, it is important to understand how it differs from Sharding and how they can be combined for massive scale.

- **Partitioning (What I did)**: Data is divided into smaller chunks (like monthly tables) but remains on a **single physical server**. This solves the issue of bloated indexes and slow deletions perfectly.
- **Sharding**: Data is distributed across **multiple physical servers** (e.g., using extensions like Citus). This solves the issue of a single machine not having enough CPU, RAM, or Disk Space to handle extreme global traffic.

### Application-Level vs. Database-Level Sharding

If I were to implement Sharding, I could do it at two levels:

1. **Application-Level**: I would maintain multiple database connections in my Go code and use `if/else` logic to route data based on criteria (like a Country ID). This makes the application logic highly complex.
2. **Database-Level (Citus)**: I would install an extension like Citus on the PostgreSQL cluster. My Go application would remain exactly as it is—connecting to only one Coordinator Node. Under the hood, the Coordinator magically analyzes the "Sharding Key" and routes the payload over the network to the correct Worker Node.

### The End-Game Architecture: Sharding + Partitioning

Can partitioning and sharding be used together? Absolutely. It is the ultimate architecture for enterprise-level time-series data (often called Distributed Time-Series).

In a distributed environment, my Go application sends an insert request to the Coordinator Node. The Coordinator checks the Sharding Key (e.g., Country) and throws the data across the network to the correct physical Server Worker. Once the data arrives at that worker, its local PostgreSQL engine checks the Partition Key (e.g., Date) and instantly routes it into a highly-optimized, monthly partition bucket. This provides infinite computing scale alongside lightning-fast disk I/O and instant data purging capacities.
