# Phase 1: Discovery

## Objective

Extract ALL dependencies from the target service directory into a structured `dependencies.yaml` file.

## Critical Rules

### DO NOT

- **DO NOT** search/grep the entire codebase
- **DO NOT** infer or guess dependencies
- **DO NOT** abbreviate names
- **DO NOT** skip source file references
- **DO NOT** read files outside the target directory (unless following an import)

### DO

- **DO** read runbooks/docs FIRST
- **DO** read ONLY files in the target directory
- **DO** cite source file and line for EVERY dependency
- **DO** use full names (Title Case for services, lowercase for topics)
- **DO** ask for the directory path if not provided

## Step-by-Step Process

### Step 1: Confirm Target Directory

If the user provides a service name without a path:
1. Ask: "What is the path to the {service} directory?"
2. Or check for runbooks: `ls runbooks/` or `ls docs/`
3. Or list services: `ls services/` or `ls cmd/`

**You MUST have a concrete directory path before proceeding.**

### Step 2: Explore Directory Structure

List the target directory contents:

```bash
ls -la {target_path}
```

Identify and report:
```
Directory structure for {service}:
- {path}/api/           → API handlers (INCLUDE)
- {path}/internal/      → Internal packages (INCLUDE)
- {path}/consumer/      → Kafka consumers (INCLUDE)
- {path}/producer/      → Kafka producers (INCLUDE)
- {path}/client/        → gRPC/HTTP clients (INCLUDE)
- {path}/proto/         → Proto definitions (INCLUDE)
- {path}/config/        → Configuration (INCLUDE)
- {path}/docs/          → Documentation (READ FIRST)
- {path}/jobs/          → Background jobs (DEPRIORITIZE)
- {path}/migrations/    → DB migrations (SKIP)
- {path}/tests/         → Tests (SKIP)
- {path}/*_test.go      → Test files (SKIP)
```

### Step 3: Read Documentation FIRST

Before reading any code, read these files if they exist:

1. `{target_path}/README.md`
2. `{target_path}/RUNBOOK.md`
3. `{target_path}/docs/architecture.md`
4. `{target_path}/docs/README.md`
5. `runbooks/{service}.md` (if exists at repo root)

**Extract from documentation:**
- Service description
- Listed dependencies
- Kafka topics (producer/consumer)
- Database connections
- External APIs

Report what you found:
```
From documentation:
- Description: {description}
- Dependencies mentioned: {list}
- Kafka topics: {list}
- Databases: {list}
```

### Step 4: Read Configuration Files

Read config files to find:
- Kafka topic names
- Database connection strings
- External service URLs
- Timeout/retry settings

Files to check:
- `{target_path}/config/*.yaml`
- `{target_path}/config/*.json`
- `{target_path}/*.yaml`
- `{target_path}/.env.example`

### Step 5: Read Client Files (Outbound Calls)

These files show what this service CALLS:

**gRPC Clients:**
- `*_client.go`, `*_client.py`
- `client/*.go`, `client/*.py`

Look for:
```go
// Go gRPC
grpc.Dial(
grpc.NewClient(
pb.New{Service}Client(
client.{Method}(ctx, request)
```

```python
# Python gRPC
grpc.insecure_channel(
{Service}Stub(channel)
stub.{Method}(request)
```

**HTTP Clients:**
```go
// Go HTTP
http.Get(
http.Post(
client.Do(req)
```

```python
# Python HTTP
requests.get(
requests.post(
httpx.get(
```

### Step 6: Read Consumer Files (Inbound Async)

These files show what Kafka topics this service CONSUMES:

- `*_consumer.go`, `*_consumer.py`
- `consumer/*.go`, `consumer/*.py`

Look for:
```go
// Go Kafka (Sarama)
sarama.NewConsumerGroup(
consumer.Subscribe(
for msg := range claim.Messages()
```

```python
# Python Kafka
KafkaConsumer(
@app.agent(topic)  # Faust
for msg in consumer:
```

Record:
- Topic name
- Consumer group name
- Whether DLQ exists

### Step 7: Read Producer Files (Outbound Async)

These files show what Kafka topics this service PRODUCES:

- `*_producer.go`, `*_producer.py`
- `producer/*.go`, `producer/*.py`

Look for:
```go
// Go Kafka
producer.SendMessage(&sarama.ProducerMessage{Topic: "topic.name"})
producer.Produce(&kafka.Message{TopicPartition: ...})
```

```python
# Python Kafka
producer.send('topic.name', value=data)
```

### Step 8: Read Proto Files (gRPC Definitions)

Proto files define the gRPC interface:

- `{target_path}/proto/*.proto`
- `{target_path}/*.proto`

Extract:
- Service name
- RPC methods
- Imports (dependencies on other protos)

```protobuf
service LedgerService {
    rpc CreateTransaction(CreateRequest) returns (CreateResponse);
    rpc GetBalance(BalanceRequest) returns (BalanceResponse);
}
```

### Step 9: Read Repository Files (Database)

These files show database connections:

- `*_repo.go`, `*_repository.go`
- `repository/*.go`
- `dal/*.go`, `store/*.go`

Look for:
```go
sql.Open(
db.Query(
db.Exec(
gorm.Open(
```

### Step 10: Compile Dependencies

After reading all files, compile the complete list:

```
Dependencies found for {service}:

SYNC (gRPC/HTTP):
- → {ServiceName}.{Method} (grpc) [from {file}:{line}]
- → {Database} [from {file}:{line}]
- → {ExternalAPI} (https) [from {file}:{line}]

ASYNC (Kafka):
- PRODUCES: {topic.name} [from {file}:{line}]
- CONSUMES: {topic.name} [from {file}:{line}]

FROM DOCUMENTATION:
- {any additional dependencies from docs}
```

### Step 11: Write .flow-deps.yaml

Write the YAML file following the schema in `schemas/dependencies.schema.yaml`.

Output location: `.flow-deps.yaml` (in working directory)

**Note**: This file will be automatically cleaned up after Phase 4 completes successfully.

## Checklist Before Completing Phase 1

Confirm you have:
- [ ] Read documentation files first
- [ ] Read all client files (outbound gRPC/HTTP)
- [ ] Read all consumer files (Kafka consumers)
- [ ] Read all producer files (Kafka producers)
- [ ] Read all proto files
- [ ] Read all repository/database files
- [ ] Read configuration files
- [ ] Every dependency has source_file and source_line
- [ ] No abbreviations in names
- [ ] Written dependencies.yaml to output directory

## Output

After completing Phase 1, report:

```
Phase 1 Complete: Discovery

Output: {output_dir}/dependencies.yaml

Summary:
- Service: {service_name}
- Sync dependencies: {count}
- Async dependencies: {count}
- External systems: {count}
- Databases: {count}
- Caches: {count}

Proceeding to Phase 2: Diagram Generation
```
