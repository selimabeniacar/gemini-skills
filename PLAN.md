# /flow Skill Redesign - Comprehensive Plan

## Problem Statement

Current approach fails because:
1. **Single-step generation** - AI tries to understand AND draw simultaneously, leading to inconsistency
2. **No validation** - Diagrams may not even compile
3. **No structure** - Free-form exploration leads to missed dependencies
4. **Inconsistent styling** - Colors, shapes, arrows vary between runs
5. **Poor layouts** - Entangled arrows, weird gaps, unreadable text
6. **Context overflow** - Broad searches kill the context window

## Solution: Multi-Phase Pipeline

Separate responsibilities into distinct, verifiable phases with intermediate files.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           /flow Skill Pipeline                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐    ┌───────────┐ │
│  │   PHASE 1    │    │   PHASE 2    │    │   PHASE 3    │    │  PHASE 4  │ │
│  │  Discovery   │───▶│  Generation  │───▶│  Refinement  │───▶│  Verify   │ │
│  │              │    │              │    │              │    │           │ │
│  │  AI explores │    │  AI generates│    │  Scripts fix │    │ AI checks │ │
│  │  codebase    │    │  from schema │    │  and validate│    │ complete  │ │
│  └──────┬───────┘    └──────┬───────┘    └──────┬───────┘    └─────┬─────┘ │
│         │                   │                   │                  │       │
│         ▼                   ▼                   ▼                  ▼       │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐    ┌───────────┐ │
│  │ dependencies │    │  diagram.md  │    │ diagram.md   │    │  DONE or  │ │
│  │    .yaml     │    │   (draft)    │    │  (refined)   │    │  iterate  │ │
│  └──────────────┘    └──────────────┘    └──────────────┘    └───────────┘ │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Phase 1: Discovery

**Goal**: Extract ALL dependencies into a structured, machine-readable format.

**Input**: Service directory path or topic/journey name
**Output**: `dependencies.yaml` - strict schema, no ambiguity

### Rules for Discovery

1. **Check runbooks/docs FIRST** - they list dependencies explicitly
2. **Read only target directory** - never grep entire codebase
3. **Be exhaustive** - every gRPC client, every Kafka topic, every DB connection
4. **No inference** - only include what is explicitly found in code/docs
5. **Include file references** - every dependency must cite source file:line

### Dependencies Schema (YAML)

```yaml
# dependencies.yaml - strict schema
version: "1.0"
generated: "2024-01-15T14:30:22Z"
target_service: "ledger-service"
target_path: "services/ledger/"

# Metadata from runbooks/docs
documentation:
  runbook: "docs/RUNBOOK.md"
  architecture: "docs/architecture.md"
  notes: "Any important context from docs"

# The service being analyzed
service:
  name: "Ledger Service"
  description: "Manages financial transactions and balances"
  entrypoints:
    - type: "grpc"
      name: "LedgerService"
      proto: "proto/ledger.proto"
      methods:
        - "CreateTransaction"
        - "GetBalance"
        - "ListTransactions"
    - type: "http"
      path: "/api/v1/ledger"
      methods:
        - "GET /health"
        - "POST /webhook"

# All external services this service CALLS (outbound)
dependencies:
  sync:  # gRPC, HTTP - caller blocks
    - name: "Payment Service"
      type: "grpc"
      proto: "PaymentService"
      methods_called:
        - "ProcessPayment"
        - "RefundPayment"
      source_file: "internal/client/payment_client.go"
      source_line: 45
      timeout: "5s"
      retry: true
      circuit_breaker: true

    - name: "PostgreSQL"
      type: "database"
      connection: "ledger_db"
      operations:
        - "read"
        - "write"
      source_file: "internal/repository/ledger_repo.go"
      source_line: 12

  async:  # Kafka, queues - fire and forget
    - name: "order.completed"
      type: "kafka"
      direction: "consume"
      consumer_group: "ledger-consumer"
      source_file: "internal/consumer/order_consumer.go"
      source_line: 27
      dlq: true

    - name: "ledger.transaction.created"
      type: "kafka"
      direction: "produce"
      source_file: "internal/producer/transaction_producer.go"
      source_line: 34

# Services that CALL this service (inbound) - if known
callers:
  - name: "Order Service"
    type: "grpc"
    methods_called:
      - "CreateTransaction"
    source: "from runbook"  # or "from proto import"

# External third-party systems
external:
  - name: "Stripe API"
    type: "https"
    purpose: "Payment processing"
    source_file: "internal/stripe/client.go"
    source_line: 23

# Caches
caches:
  - name: "Redis"
    purpose: "Balance cache"
    source_file: "internal/cache/balance_cache.go"
    source_line: 15
```

### Discovery Checklist

Before completing Phase 1, AI MUST confirm:
- [ ] Read all runbook/doc files in service directory
- [ ] Read all `*_client.go` / `*_client.py` files (outbound calls)
- [ ] Read all `*_consumer.go` / `*_consumer.py` files (Kafka consumers)
- [ ] Read all `*_producer.go` / `*_producer.py` files (Kafka producers)
- [ ] Read all `*.proto` files (gRPC definitions)
- [ ] Read all config files (`*.yaml`, `*.json`)
- [ ] Listed ALL sync dependencies with source references
- [ ] Listed ALL async dependencies with source references
- [ ] Listed ALL external systems
- [ ] Listed ALL databases and caches

---

## Phase 2: Diagram Generation

**Goal**: Generate a consistent, professional Mermaid diagram from the dependency file.

**Input**: `dependencies.yaml` from Phase 1
**Output**: `diagram.md` - Mermaid diagram following strict style guide

### Key Principle

**The AI does NOT read the codebase in this phase.**
It ONLY reads `dependencies.yaml` and applies the style guide.

This ensures consistency because:
- Same input format → same output style
- No "creative interpretation"
- Deterministic transformation

### Style Guide (Non-Negotiable)

#### Layout
- **Direction**: `flowchart TD` (left-to-right) for service flows
- **Grouping**: Related services in subgraphs
- **Ordering**: Entry points on left, external systems on right
- **Spacing**: Use invisible nodes if needed for alignment

#### Node Shapes (STRICT)

| Type | Shape | Syntax | Example |
|------|-------|--------|---------|
| Service/Handler | Rectangle | `[Name]` | `[Order Service]` |
| Kafka Topic | Cylinder | `[(topic-name)]` | `[(order.created)]` |
| Consumer Group | Stadium | `([group-name])` | `([ledger-consumer])` |
| Database | Cylinder | `[(Database)]` | `[(PostgreSQL)]` |
| External System | Double Rectangle | `[[Name]]` | `[[Stripe API]]` |
| Decision/Router | Diamond | `{condition}` | `{retry?}` |
| Cache | Rectangle with rounded | `(Cache)` | `(Redis Cache)` |

#### Colors (STRICT - No Deviation)

```
classDef service fill:#228be6,stroke:#1971c2,color:#fff
classDef kafka fill:#12b886,stroke:#099268,color:#fff
classDef database fill:#fab005,stroke:#f59f00,color:#000
classDef external fill:#868e96,stroke:#495057,color:#fff
classDef cache fill:#be4bdb,stroke:#9c36b5,color:#fff
classDef error fill:#fa5252,stroke:#e03131,color:#fff
classDef entry fill:#40c057,stroke:#2f9e44,color:#fff
```

#### Arrow Styles (STRICT)

| Type | Syntax | Label Format |
|------|--------|--------------|
| Sync (gRPC) | `==>` | `\|gRPC: MethodName\|` |
| Sync (HTTP) | `==>` | `\|HTTP: endpoint\|` |
| Sync (SQL) | `==>` | `\|SQL\|` |
| Async (Kafka produce) | `-.->` | `\|publish\|` |
| Async (Kafka consume) | `-.->` | `\|consume\|` |
| Internal | `-->` | `\|internal\|` |

#### Subgraph Rules

1. **Always quote titles**: `subgraph id ["Display Name"]`
2. **One subgraph per service boundary**
3. **Kafka topics in their own "Message Bus" subgraph**
4. **External systems in "External" subgraph**
5. **Databases in "Data Stores" subgraph**

#### Naming Rules

1. **NO abbreviations** - "Order Service" not "OrdSvc"
2. **NO inference** - only use names from dependencies.yaml
3. **Full method names** - "ProcessPayment" not "Process"
4. **Consistent casing** - Title Case for services, lowercase for topics

### Diagram Template

```markdown
# Service Flow: {service_name}

Generated: {timestamp}
Source: {target_path}

## Diagram

\`\`\`mermaid
flowchart TD
    %% Style definitions
    classDef service fill:#228be6,stroke:#1971c2,color:#fff
    classDef kafka fill:#12b886,stroke:#099268,color:#fff
    classDef database fill:#fab005,stroke:#f59f00,color:#000
    classDef external fill:#868e96,stroke:#495057,color:#fff
    classDef cache fill:#be4bdb,stroke:#9c36b5,color:#fff
    classDef entry fill:#40c057,stroke:#2f9e44,color:#fff

    %% Entry points (leftmost)
    subgraph entry ["Entry Points"]
        ...
    end

    %% Target service
    subgraph target ["Target Service Name"]
        ...
    end

    %% Dependent services
    subgraph deps ["Dependent Services"]
        ...
    end

    %% Message bus
    subgraph kafka ["Message Bus"]
        ...
    end

    %% Data stores
    subgraph data ["Data Stores"]
        ...
    end

    %% External systems (rightmost)
    subgraph ext ["External Systems"]
        ...
    end

    %% Sync connections (thick arrows)
    ...

    %% Async connections (dotted arrows)
    ...

    %% Apply styles
    class ... service
    class ... kafka
    class ... database
    class ... external
\`\`\`

## Legend

| Symbol | Meaning |
|--------|---------|
| `==>` | Synchronous call (gRPC/HTTP) - caller blocks |
| `-.->` | Asynchronous message (Kafka) - fire and forget |
| `-->` | Internal call |
| Blue rectangle | Service |
| Green cylinder | Kafka topic |
| Yellow cylinder | Database |
| Gray double rectangle | External system |
| Purple rounded | Cache |

## Dependencies Summary

### Sync Dependencies
| From | To | Type | Method |
|------|-----|------|--------|
...

### Async Dependencies
| Topic | Direction | Consumer Group |
|-------|-----------|----------------|
...

## Source References

All dependencies traced from:
- {list of source files with line numbers}
```

---

## Phase 3: Refinement (Scripts)

**Goal**: Validate and automatically fix common diagram issues.

**Input**: `diagram.md` from Phase 2
**Output**: `diagram.md` (refined) + validation report

### Go CLI Tool: `flowlint`

A single Go binary with subcommands for all refinement operations.

### Installation

```bash
cd flow/tools/flowlint
go build -o flowlint .
# Move to PATH or use directly
```

### Commands

```bash
# Validate Mermaid syntax (shells out to mmdc)
flowlint validate diagram.md

# Lint against style guide, auto-fix violations
flowlint lint diagram.md --fix --output diagram-fixed.md

# Check completeness against dependencies
flowlint check diagram.md dependencies.yaml

# Run full refinement pipeline
flowlint refine diagram.md dependencies.yaml --output diagram-final.md
```

### Subcommand Details

#### `flowlint validate`

- Extracts Mermaid code block from markdown
- Runs `mmdc` (Mermaid CLI) to compile
- Parses error output, returns structured report
- Exit code 0 = valid, 1 = invalid

#### `flowlint lint`

Checks:
- All nodes have correct shape for their type
- All arrows use correct style (sync `==>` vs async `-.->`)
- All subgraph titles are quoted
- Colors applied via classDef
- No abbreviations in node labels (warns on short names)
- No orphan nodes (everything connected)
- No duplicate node IDs

Auto-fixes (with `--fix`):
- Quote unquoted subgraph titles
- Add missing classDef declarations
- Normalize arrow styles

#### `flowlint check`

- Parses `dependencies.yaml`
- Parses diagram nodes and edges
- Reports missing services, topics, connections
- Exit code 0 = complete, 1 = missing items

#### `flowlint refine`

Orchestrates full pipeline:
1. Validate syntax
2. Lint and fix
3. Check completeness
4. Output final diagram

### Project Structure

```
flow/tools/flowlint/
├── go.mod
├── go.sum
├── main.go                 # CLI entrypoint with cobra
├── cmd/
│   ├── validate.go         # validate subcommand
│   ├── lint.go             # lint subcommand
│   ├── check.go            # check subcommand
│   └── refine.go           # refine subcommand
├── internal/
│   ├── parser/
│   │   ├── markdown.go     # Extract mermaid from markdown
│   │   ├── mermaid.go      # Parse mermaid syntax
│   │   └── yaml.go         # Parse dependencies.yaml
│   ├── linter/
│   │   ├── rules.go        # Linting rules
│   │   └── fixer.go        # Auto-fix logic
│   └── styles/
│       └── styles.go       # Color/shape definitions
└── testdata/
    ├── valid.md
    ├── invalid.md
    └── dependencies.yaml
```

### Dependencies

```go
// go.mod
module github.com/user/flowlint

go 1.21

require (
    github.com/spf13/cobra v1.8.0      // CLI framework
    gopkg.in/yaml.v3 v3.0.1            // YAML parsing
)
```

### External Requirement

`mmdc` (Mermaid CLI) must be installed for the `validate` command:
```bash
npm install -g @mermaid-js/mermaid-cli
# or use npx
```

---

## Phase 4: Verification

**Goal**: Final human-in-the-loop verification with AI assistance.

**Input**: Refined `diagram.md` + `dependencies.yaml`
**Output**: Approval or iteration request

### Verification Steps

1. **Render the diagram** - generate PNG/SVG for visual inspection
2. **Cross-reference** - AI compares diagram to dependencies.yaml
3. **Report discrepancies** - list anything missing or incorrect
4. **Suggest improvements** - note any unclear areas

### Verification Checklist

AI must confirm:
- [ ] All services from dependencies.yaml appear in diagram
- [ ] All Kafka topics appear with correct direction
- [ ] All databases appear
- [ ] All external systems appear
- [ ] Sync vs async arrows are correct
- [ ] No orphan nodes
- [ ] Diagram renders without errors
- [ ] Layout is readable (no overlapping text)
- [ ] Colors are consistent

### Output

```markdown
## Verification Report

### Status: PASS / FAIL

### Coverage
- Services: 5/5 ✓
- Kafka topics: 3/3 ✓
- Databases: 1/1 ✓
- External systems: 2/2 ✓

### Issues Found
- None / List of issues

### Recommendations
- None / List of suggestions

### Final Diagram Location
`/tmp/flow-{timestamp}/diagram.md`

To render as PNG:
\`\`\`
mmdc -i diagram.md -o diagram.png -b transparent
\`\`\`
```

---

## File Structure

```
flow/
├── SKILLS.md                      # Main skill definition (step-based)
├── schemas/
│   └── dependencies.schema.yaml   # Schema documentation for dependencies.yaml
├── styles/
│   └── diagram-styles.yaml        # Color palette, shapes, conventions
├── tools/
│   └── flowlint/                  # Go CLI tool
│       ├── go.mod
│       ├── main.go
│       ├── cmd/
│       │   ├── validate.go
│       │   ├── lint.go
│       │   ├── check.go
│       │   └── refine.go
│       └── internal/
│           ├── parser/
│           ├── linter/
│           └── styles/
├── prompts/
│   ├── phase1-discovery.md        # Discovery phase instructions
│   ├── phase2-generation.md       # Diagram generation instructions
│   └── phase4-verification.md     # Verification instructions
├── templates/
│   └── diagram-template.md        # Template for diagram output
└── examples/
    ├── dependencies-example.yaml  # Example dependency file
    └── diagram-example.md         # Example final diagram
```

---

## SKILLS.md Structure (High-Level)

```markdown
# /flow

Generate professional architectural diagrams through a multi-phase pipeline.

## Usage

/flow <path> [options]

## Options

- `--type=service|call|control` - Diagram type (default: service)
- `--output=<dir>` - Output directory (default: /tmp/flow-{timestamp}/)
- `--phase=<n>` - Run only up to phase N (for debugging)
- `--no-scripts` - Skip Phase 3 scripts (if not installed)

## Pipeline

This skill executes in 4 phases:

### Phase 1: Discovery
- Explore the target directory
- Extract ALL dependencies
- Write to `dependencies.yaml`

### Phase 2: Generation
- Read `dependencies.yaml`
- Apply strict style guide
- Generate `diagram.md`

### Phase 3: Refinement
- Run validation scripts
- Auto-fix style issues
- Optimize layout

### Phase 4: Verification
- Cross-check diagram vs dependencies
- Report coverage
- Confirm completion

## Output Files

After completion:
- `dependencies.yaml` - Structured dependency data
- `diagram.md` - Final Mermaid diagram
- `verification-report.md` - Coverage report
```

---

## Implementation Order

1. **Create schemas/** - Define the dependency schema
2. **Create styles/** - Define the style guide
3. **Create prompts/** - Write phase-specific instructions
4. **Create templates/** - Diagram template
5. **Create tools/flowlint/** - Go CLI tool
   - Set up Go module
   - Implement `validate` command
   - Implement `lint` command with auto-fix
   - Implement `check` command
   - Implement `refine` command (orchestrator)
6. **Update SKILLS.md** - Orchestrate the pipeline
7. **Create examples/** - Reference examples
8. **Test end-to-end** - Verify the pipeline works

---

## Success Criteria

A successful run produces:
1. A `dependencies.yaml` that is complete and accurate
2. A `diagram.md` that:
   - Compiles without errors
   - Uses consistent colors (always the same palette)
   - Uses correct shapes (service=rectangle, kafka=cylinder, etc.)
   - Uses correct arrows (sync=thick, async=dotted)
   - Has readable text (no abbreviations, no tiny fonts)
   - Has clean layout (no entangled arrows, no weird gaps)
   - Groups related items in subgraphs
   - Can be used directly in documentation
3. A verification report confirming 100% coverage

---

## Notes

- `flowlint` is a single Go binary - build once, use anywhere
- `mmdc` (Mermaid CLI) required only for `validate` command (can skip with `--skip-validate`)
- If `flowlint` unavailable, Phase 3 can be skipped with `--no-tools`
- Each phase can be run independently for debugging
- Intermediate files allow manual inspection and correction
- Go 1.21+ required to build `flowlint`
