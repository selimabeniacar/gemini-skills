# Phase 1: Discovery

## Objective

Extract ALL dependencies from target service directories into a structured `.flow-deps.yaml` file.

After extracting dependencies, you MUST self-review by listing every file in the target path and accounting for each one.

---

## ⛔ CRITICAL: NO FULL REPO SEARCHES

**STOP. READ THIS BEFORE DOING ANYTHING.**

```
╔═══════════════════════════════════════════════════════════════════════╗
║  🚫 NEVER grep/search/find across the entire repository              ║
║  🚫 NEVER use glob patterns like **/*.go or **/*.py at repo root     ║
║  🚫 NEVER read files outside the target directory                     ║
║                                                                        ║
║  ✅ ONLY read files within: {target_path}/                            ║
║  ✅ ONLY search within: {target_path}/                                ║
╚═══════════════════════════════════════════════════════════════════════╝
```

If you catch yourself about to search the whole repo, **STOP IMMEDIATELY**.

---

## Critical Rules

### ❌ DO NOT

- **DO NOT** run `grep -r` or `rg` at repo root
- **DO NOT** run `find . -name "*.go"` at repo root
- **DO NOT** search for patterns like `**/*.proto` without a target path
- **DO NOT** infer or guess dependencies
- **DO NOT** abbreviate names
- **DO NOT** skip source file references
- **DO NOT** read files outside the target directory

### ✅ DO

- **DO** read runbooks/docs FIRST
- **DO** read ONLY files in `{target_path}/`
- **DO** search ONLY within `{target_path}/`
- **DO** cite source file and line for EVERY dependency
- **DO** use full names (Title Case for services, lowercase for topics)
- **DO** ask for the directory path if not provided
- **DO** follow the mandatory file checklist below

---

### 📋 MANDATORY FILE CHECKLIST

You MUST attempt to read ALL of these file patterns within `{target_path}/`. Do not skip any category.

| Category | File Patterns | What to Extract |
|----------|---------------|-----------------|
| **1. Docs** | `README.md`, `RUNBOOK.md`, `docs/*.md` | Service description, listed dependencies |
| **2. Config** | `config/*.yaml`, `*.yaml`, `.env.example` | Topic names, DB strings, service URLs |
| **3. gRPC Clients** | `*client*.go`, `*client*.py`, `client/*.go` | Outbound gRPC calls |
| **4. HTTP Clients** | `*client*.go`, `*http*.go`, `api/*.go` | Outbound HTTP calls |
| **5. Consumers** | `*consumer*.go`, `*consumer*.py`, `consumer/*.go` | Kafka topics consumed |
| **6. Producers** | `*producer*.go`, `*producer*.py`, `producer/*.go` | Kafka topics produced |
| **7. Proto** | `*.proto`, `proto/*.proto` | gRPC service definitions |
| **8. Repository** | `*repo*.go`, `*repository*.go`, `dal/*.go`, `store/*.go` | Database connections |
| **9. Cache** | `*cache*.go`, `*redis*.go` | Cache connections |
| **10. Handlers** | `*handler*.go`, `*server*.go`, `api/*.go` | Entry points |

**After reading files, report which categories had files and which were empty.**

### Multi-Service Discovery

If analyzing multiple services, repeat the discovery process for each `target_path`. Each service becomes a separate entry in the `services[]` array.

### Step-by-Step Process

#### Step 1: Confirm Target Directory

If the user provides a service name without a path:
1. Ask: "What is the path to the {service} directory?"
2. Or check for runbooks: `ls runbooks/` or `ls docs/`
3. Or list services: `ls services/` or `ls cmd/`

**You MUST have a concrete directory path before proceeding.**

#### Step 2: Explore Directory Structure

```bash
ls -la {target_path}
```

Identify and report:
```
Directory structure for {service}:
- {path}/internal/      → Internal packages (INCLUDE)
- {path}/consumer/      → Kafka consumers (INCLUDE)
- {path}/producer/      → Kafka producers (INCLUDE)
- {path}/client/        → gRPC/HTTP clients (INCLUDE)
- {path}/proto/         → Proto definitions (INCLUDE)
- {path}/config/        → Configuration (INCLUDE)
- {path}/docs/          → Documentation (READ FIRST)
- {path}/migrations/    → DB migrations (SKIP)
- {path}/tests/         → Tests (SKIP)
```

#### Step 3: Read Documentation FIRST

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

**⚠️ CODE IS THE SOURCE OF TRUTH, NOT DOCUMENTATION.**

Runbooks and READMEs are a starting point — they give you a map. But docs go stale. After reading code (Steps 4–8), you MUST cross-check:

1. **Dependencies in docs but NOT in code** — the dependency may have been removed. Do not include it.
2. **Dependencies in code but NOT in docs** — this is the common case. Docs are incomplete. Include it, trace it to the exact source file and line.

If there is a conflict between docs and code, **code wins**. Every dependency in the final YAML must be backed by a `source_file` and `source_line` in actual code.

#### Step 4: Read Configuration Files

Read config files to find:
- Kafka topic names, database connection strings, external service URLs

Files: `{target_path}/config/*.yaml`, `{target_path}/*.yaml`, `{target_path}/.env.example`

#### Step 5: Read Client Files (Outbound Calls)

**gRPC Clients** — look for: `grpc.Dial(`, `pb.New*Client(`, `client.*Method*(ctx,`

**HTTP Clients** — look for: `http.Get(`, `http.Post(`, `client.Do(req)`

#### Step 6: Read Consumer/Producer Files

**Consumers** — look for: `sarama.NewConsumerGroup(`, `consumer.Subscribe(`

**Producers** — look for: `producer.SendMessage(`, topic name strings

#### Step 7: Read Proto Files

Extract service names, RPC methods, imports.

#### Step 8: Read Repository/Cache Files

**Repository** — look for: `sql.Open(`, `db.Query(`, `gorm.Open(`

**Cache** — look for: `redis.NewClient(`, `cache.Get(`, `cache.Set(`

**⚠️ NAMING: Use logical names, NOT technology names.**

Each service has its own database and cache — even if they use the same technology. Two services using PostgreSQL are two separate databases.

```
❌ WRONG:  name: "PostgreSQL"     — ambiguous, merges into one node
✅ RIGHT:  name: "Ledger DB"      — unique per service
           type: "postgresql"

❌ WRONG:  name: "Redis"          — ambiguous
✅ RIGHT:  name: "Ledger Cache"   — unique per service
           type: "redis"
```

#### Step 9: Discover Internal Steps (Optional)

Read handler/entrypoint code to identify major processing stages. These are **optional** — only include when there are clear sequential stages worth showing.

Look for patterns like:
- Validation logic (auth checks, schema validation)
- Data fetching (DB reads, cache lookups)
- Business logic (processing, calculations)
- Persistence (DB writes, cache updates)
- Event publishing (Kafka produces)

Include `internal_steps` when:
- Single-service diagram where more detail is valuable
- User explicitly asks for a "detailed" diagram
- There are 2+ clear sequential stages in the handler

#### Step 10: Compile and Write .flow-deps.yaml

Write the YAML file following the schema in `schemas/dependencies.schema.yaml`.

Output location: `.flow-deps.yaml` (in working directory)

### Output Format

```yaml
generated: "2024-01-15T14:30:22Z"
services:
  - name: "Ledger Service"
    target_path: "services/ledger/"
    description: "Manages financial transactions and balances"
    documentation:
      runbook: "services/ledger/docs/RUNBOOK.md"
    entrypoints:
      - type: "grpc"
        name: "LedgerService"
        methods:
          - "CreateTransaction"
          - "GetBalance"
    dependencies:
      sync:
        - name: "Payment Service"
          type: "grpc"
          source_file: "internal/client/payment_client.go"
          source_line: 45
      async:
        - name: "order.completed"
          direction: "consume"
          source_file: "internal/consumer/order_consumer.go"
          source_line: 27
    databases:
      - name: "Ledger DB"
        type: "postgresql"
        source_file: "internal/repository/ledger_repo.go"
        source_line: 12
    caches:
      - name: "Ledger Cache"
        type: "redis"
        purpose: "Balance cache"
        source_file: "internal/cache/balance_cache.go"
        source_line: 15
    external:
      - name: "Stripe API"
        type: "https"
        source_file: "internal/webhook/stripe_handler.go"
        source_line: 56
    internal_steps:
      - name: "Validate Request"
        description: "Auth and schema validation"
      - name: "Load Account"
        description: "Fetch account from DB + cache"
      - name: "Process Transaction"
        description: "Execute business logic, call payment service"
      - name: "Commit"
        description: "Write to DB, publish event"
```

---

## Coverage Self-Review

**MANDATORY: Do this AFTER writing `.flow-deps.yaml` and BEFORE completing Phase 1.**

### Step 1: List all non-test files

```bash
find {target_path} -type f \( -name "*.go" -o -name "*.py" -o -name "*.proto" \) \
  -not -path "*/vendor/*" \
  -not -name "*_test.go" \
  -not -path "*/mock*" \
  -not -path "*/testdata/*" \
  -not -path "*/migrations/*"
```

### Step 2: Account for every file

For each file in the list, report one of:
- **Referenced** — appears as a `source_file` in `.flow-deps.yaml`
- **Read, no dependency** — you read it and it contains no external calls (explain briefly why)
- **Skipped** — you did not read this file

### Step 3: Read any skipped files

If any files were skipped, read them now. If they contain dependencies, update `.flow-deps.yaml`.

**You must have zero skipped files before completing Phase 1.**

### Coverage Report Format

```
Coverage Report:
- Total files: N
- Referenced in YAML: N
- Read, no dependency: N
- Newly added after review: N

Files read with no dependency:
- internal/util/helpers.go — utility functions, no external calls
- internal/middleware/auth.go — middleware, no outbound dependencies
```

---

## Checklist Before Completing Phase 1

Confirm you have:
- [ ] Read documentation files first
- [ ] Read all client files (outbound gRPC/HTTP)
- [ ] Read all consumer/producer files (Kafka)
- [ ] Read all proto files
- [ ] Read all repository/database/cache files
- [ ] Every dependency has source_file and source_line
- [ ] No abbreviations in names
- [ ] Output uses `services[]` array format
- [ ] Completed coverage self-review (zero skipped files)
- [ ] Written .flow-deps.yaml to working directory

## Output

```
Phase 1 Complete: Discovery

Output: .flow-deps.yaml

Summary:
- Services: {count}
- Total sync dependencies: {count}
- Total async dependencies: {count}
- External systems: {count}
- Databases: {count}
- Caches: {count}

Coverage: {referenced + no-dep} / {total files} (100%)

Proceeding to Phase 2: Diagram Generation
```
