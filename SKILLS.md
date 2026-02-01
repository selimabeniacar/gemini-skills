# /flow

Generate flow diagrams for debugging, incident investigation, and onboarding.

## Description

Creates Mermaid diagrams to visualize code execution and service communication. Supports three modes:

- **Call Graph**: Trace function calls within a service to find where issues originate
- **Control Flow**: Zoom into a specific function to understand branching logic
- **Service Flow**: Map communication between services showing both **sync (gRPC/HTTP)** and **async (Kafka/queues)** paths

## Critical: Sync vs Async Distinction

For incident investigation and onboarding, it is **essential** to clearly distinguish:

| Type | Examples | Arrow Style | Behavior | Debug Approach |
|------|----------|-------------|----------|----------------|
| **Sync** | gRPC, HTTP, REST | `==>` thick solid | Caller blocks, waits for response | Check latency, timeouts, error codes |
| **Async** | Kafka, RabbitMQ, SQS | `-.->` dotted | Fire and forget, eventual | Check consumer lag, DLQ, offsets |

**This distinction changes everything about how you debug an issue.**

## Arguments

- `<path>` - File or directory to analyze (required)
- `--type=call|control|service` - Type of diagram (default: `call`)
- `--function=<name>` - Specific function to analyze (required for control flow)
- `--depth=<n>` - How deep to trace calls (default: 5, only for call graphs)
- `--topic=<name>` - Filter by specific Kafka topic (optional, for service flows)
- `--service=<name>` - Filter by service name (optional, for service flows)
- `--output=<path>` - Output file path (default: `/tmp/flow-<timestamp>.md`)
- `--verbose` - Include educational explanations for onboarding (default: off)

## Supported Languages

- Go
- Python

## Usage

```bash
# Generate call graph for a directory (outputs to /tmp/flow-<timestamp>.md)
/flow src/payments/ --type=call

# Generate call graph for a specific file
/flow src/payments/processor.go --type=call

# Generate control flow for a specific function
/flow src/payments/processor.go --type=control --function=ProcessPayment

# Short form (defaults to call graph)
/flow src/payments/

# Specify custom output path
/flow src/payments/ --output=./my-diagram.md

# Verbose mode for onboarding (includes educational explanations)
/flow src/payments/ --verbose

# Generate service flow showing sync and async communication
/flow services/ --type=service

# Service flow filtered by topic
/flow services/ --type=service --topic=order.created

# Service flow for a specific service
/flow services/payment/ --type=service --service=payment-service
```

## Visual Convention Reference (MUST Follow)

### Arrow Styles - CRITICAL

```
A ==>|gRPC| B        Sync: caller waits, immediate response or timeout
A -.->|Kafka| B      Async: fire and forget, processed later
A -->|internal| B    Internal: function call within same service
```

### Node Shapes

```
[Service/Handler]    Rectangle: services, handlers, processors
[(topic-name)]       Cylinder: Kafka topics, message queues
([consumer-group])   Stadium: consumer groups
{decision}           Diamond: routing decisions, conditions
[[external]]         Double rectangle: external systems (Stripe, DB)
```

### Colors

| Color | Hex | Meaning |
|-------|-----|---------|
| Green | `#37b24d` | Entry points, success paths |
| Blue | `#4dabf7` | Kafka topics, async infrastructure |
| Red | `#ff6b6b` | Error paths, DLQ, failures |
| Yellow | `#ffd43b` | Warnings, highlights, main flow |
| Purple | `#be4bdb` | gRPC/sync calls between services |
| Orange | `#fd7e14` | Retry paths, circuit breakers |
| Gray | `#868e96` | External systems |

**Default node background**: Use `#ffd43b` (yellow) for main service nodes to make them stand out.

### Service Boundaries (MUST Quote Titles)

Always use subgraphs with **quoted titles** to avoid syntax errors:
```mermaid
subgraph ledger ["Ledger Service"]
    L1[Handler]
end

subgraph payment ["Payment Service"]
    P1[Processor]
end
```

**WRONG** (causes syntax errors):
```
subgraph Ledger Service    ❌ Missing quotes
subgraph ledger [Ledger Service]  ❌ Missing quotes around display name
```

**CORRECT**:
```
subgraph ledger ["Ledger Service"]  ✓
subgraph payment-svc ["Payment Service"]  ✓
```

## ⚠️ CRITICAL: Do NOT Do Broad Text Searches

**NEVER search for keywords across the entire codebase.**

| WRONG | RIGHT |
|-------|-------|
| Searching "EOD" across all files → 10k results → context overflow | Find service directory first: `services/eod/` |
| Grep for "payment" everywhere | Read `services/payment/` directory only |

**If the user asks about "EOD service" or similar:**
1. Ask for the path: "Which directory contains the EOD service?"
2. Or locate it: `ls services/` to find `eod/`, `end-of-day/`, etc.
3. STOP and ask if unclear - don't search broadly

## 📖 PRIORITY: Check Runbooks FIRST

**Before reading any code, find and read runbook/docs files:**

1. Look for these files in the service directory:
   - `RUNBOOK.md`, `runbook.md`
   - `README.md`
   - `docs/architecture.md`, `docs/README.md`
   - `ARCHITECTURE.md`, `DESIGN.md`

2. **Runbooks contain the answers** - they typically list:
   - Service dependencies (gRPC, HTTP)
   - Kafka topics (producer/consumer)
   - Database connections
   - Related services
   - Architecture diagrams

3. **Use runbook info to guide code reading** - only read code files to verify or fill gaps

This saves context and gives accurate dependency information.

## Instructions

When this skill is invoked:

1. **Parse the arguments** to determine:
   - The target path (file or directory)
   - The diagram type (`call`, `control`, or `service`)
   - The function name (if control flow)
   - The depth (if call graph)
   - The topic filter (if service flow)
   - The service filter (if service flow)
   - The output path (default: `/tmp/flow-<timestamp>.md`)
   - Whether verbose mode is enabled

2. **MANDATORY: Explore the directory structure FIRST** (before reading any code):
   - List ALL subdirectories and files in the target path
   - Identify the directory structure pattern (monorepo, service-per-folder, etc.)
   - **Report what you found** to the user before proceeding:
     ```
     Found directory structure:
     - services/ledger/api/          (main API handlers)
     - services/ledger/consumer/     (Kafka consumers)
     - services/ledger/client/       (gRPC client)
     - services/ledger/jobs/         (background jobs - will deprioritize)
     - services/ledger/proto/        (gRPC definitions)
     ```
   - **Prioritize main service code** over:
     - Jobs, cron tasks, migrations (deprioritize)
     - Test files (skip)
     - Mocks, fixtures (skip)
     - Disaster recovery scripts (deprioritize)
   - **For `--type=service`**: Also look for neighboring services that communicate with this one

3. **Read the target files systematically**:
   - **Read ALL relevant files**, not just one
   - If a directory, read all `.go`, `.py`, and `.proto` files
   - If a file, read just that file
   - Identify the programming language from file extensions
   - Look for proto files to understand gRPC service definitions
   - **Search for communication patterns**:
     - gRPC client imports and calls
     - Kafka producer/consumer setup
     - HTTP client calls
   - **Track which services communicate with each other**

4. **MANDATORY: Check documentation and runbooks**:
   - Look for `README.md`, `RUNBOOK.md`, `docs/` in the service directory
   - Look for `architecture.md`, `design.md`, `ADR/` (architecture decision records)
   - Search for Kafka topic configurations in:
     - `config/`, `configs/`, `configuration/`
     - `*.yaml`, `*.yml`, `*.json` config files
     - Environment variable definitions
   - **Extract from runbooks**:
     - Service dependencies mentioned
     - Kafka topics (producer/consumer)
     - gRPC/HTTP endpoints called
     - Database dependencies
   - **Cross-reference code with documentation** to ensure nothing is missed

5. **BE THOROUGH - Do not skip dependencies**:

   ⚠️ **IMPORTANT: Take your time. Do not rush. Do not skip.**

   Before generating the diagram, you MUST have identified:
   - [ ] ALL gRPC services this service calls (check client code)
   - [ ] ALL gRPC services that call this service (check proto definitions)
   - [ ] ALL Kafka topics this service publishes to
   - [ ] ALL Kafka topics this service consumes from
   - [ ] ALL HTTP/REST endpoints this service calls
   - [ ] ALL databases this service connects to
   - [ ] ALL caches (Redis, Memcached) this service uses
   - [ ] ALL third-party APIs (Stripe, AWS, etc.)

   **If you are unsure about a dependency, include it with a note.**
   **It is better to include too much than to miss something.**

   Report your findings:
   ```
   Dependencies found for Ledger Service:

   SYNC (gRPC/HTTP):
   - → PaymentService.ProcessPayment (grpc) [from client/payment_client.go:45]
   - → InventoryService.Reserve (grpc) [from client/inventory_client.go:23]
   - → PostgreSQL [from repository/ledger_repo.go:12]

   ASYNC (Kafka):
   - PRODUCES: ledger.transaction.created [from producer/transaction_producer.go:34]
   - PRODUCES: ledger.balance.updated [from producer/balance_producer.go:18]
   - CONSUMES: order.completed [from consumer/order_consumer.go:27]
   - CONSUMES: payment.refunded [from consumer/refund_consumer.go:15]

   FROM RUNBOOK (docs/RUNBOOK.md):
   - Also depends on: AuditService for compliance logging
   - Also consumes: compliance.check.required (not found in code - verify)
   ```

6. **Generate the appropriate diagram** based on type:

### For Call Graphs (`--type=call`)

Follow the instructions in `prompts/call-graph.md`:
- Identify all function definitions
- Trace function calls within each function up to the specified depth
- Build a dependency graph
- Output as Mermaid flowchart
- Highlight suspicious paths (missing error handling, potential issues)

### For Control Flow Graphs (`--type=control`)

Follow the instructions in `prompts/control-flow.md`:
- Locate the specified function
- Parse control structures (if/else, switch, loops, try/catch)
- Identify all branches and exit points
- Output as Mermaid flowchart with decision nodes
- Annotate with line numbers

### For Service Flows (`--type=service`)

⚠️ **NEVER generate a single-service diagram. Service flow MUST show:**
- The target service AND all services it communicates with
- Minimum 2 services in the diagram (target + at least one dependency)
- If no dependencies found, something is wrong - investigate further

Follow the instructions in `prompts/service-flow.md`:
- **Include ALL services found in step 2**, not just one
- **Include neighboring services**: Any service this one calls or is called by
- **Identify sync communication**: gRPC clients, HTTP clients, REST calls
- **Identify async communication**: Kafka producers/consumers, message queues
- **Map service boundaries**: Which code belongs to which service
- **Use correct arrow styles**: `==>` for sync, `-.->` for async (THIS IS CRITICAL)
- **Show failure modes**: Timeouts, DLQs, retries, circuit breakers
- **Highlight issues**: Missing error handling, single points of failure
- **Show neighboring services**: Any service that this service calls or is called by

7. **Write the output file**:
   - Generate a filename using the format `/tmp/flow-<timestamp>.md` (e.g., `/tmp/flow-20240115-143022.md`)
   - Or use the path specified by `--output`
   - **Always include the legend** showing sync vs async distinction
   - Write the Mermaid diagram in a fenced code block
   - Include analysis notes after the diagram
   - Tell the user the output file path so they can open it

8. **Add analysis notes** in the output file explaining:
   - Key observations about the flow
   - Potential issues or areas of concern
   - **For incidents**: What to check based on sync vs async paths
   - Suggestions for debugging if applicable

9. **If `--verbose` is enabled**, add educational content for onboarding:

### Verbose Mode Content

Include these additional sections in the output:

#### How to Read This Diagram
- **Sync vs Async**: Explain thick arrows (==> gRPC) block the caller, dotted arrows (-.-> Kafka) don't
- Explain node shapes and what they represent
- Explain color coding and what each color means
- Describe how to trace a request through the system

#### Why Sync vs Async Matters
- Sync failures: Immediate, caller gets error, check timeouts and retries
- Async failures: Delayed, check consumer lag, DLQ, message schemas
- How this affects incident debugging

#### Glossary
- Define key terms specific to the codebase
- Explain domain-specific concepts visible in the diagram
- Define infrastructure terms (consumer group, partition, DLQ, circuit breaker)

#### Architecture Context
- Explain where this code fits in the overall system
- Describe which services own which data
- Note any important patterns (saga, outbox, CQRS)

#### Incident Investigation Guide
- For sync path issues: What to check (latency, error rates, circuit breaker state)
- For async path issues: What to check (consumer lag, DLQ depth, partition assignment)
- Common failure modes in this architecture

#### Next Steps for Learning
- Suggest related files or functions to explore
- Recommend what to trace next to understand the full flow
- Point to relevant documentation if known

## Examples

See the examples directory for detailed output examples:
- `examples/call-graph.md` - Call graph examples
- `examples/control-flow.md` - Control flow examples
- `examples/service-flow.md` - Service flow examples with sync/async distinction

## Output Reliability Checklist

Before finalizing output, verify:
- [ ] All sync calls use `==>` arrows
- [ ] All async calls use `-.->` arrows
- [ ] Service boundaries are shown with subgraphs
- [ ] Legend is included showing arrow meanings
- [ ] Colors are applied consistently
- [ ] Issues are highlighted with appropriate colors
- [ ] Line numbers or file references are included where helpful
