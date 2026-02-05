# /flow

Generate professional architectural diagrams through a multi-phase pipeline.

---

## Usage

```
/flow <path> [options]
```

## Options

| Option | Description | Default |
|--------|-------------|---------|
| `--output=<name>` | Output filename (without extension) | `flow-diagram` |
| `--phase=<n>` | Run only up to phase N (1-4) | 4 (all phases) |
| `--no-tools` | Skip Phase 3 (if flowlint not installed) | false |
| `--keep-deps` | Keep dependencies.yaml after completion | false |

## Output Files

All files are written to the **current working directory**:

| File | Phase | Cleaned Up? |
|------|-------|-------------|
| `.flow-deps.yaml` | Phase 1 | Yes (after Phase 4) |
| `{output}.md` | Phase 2 | No (final diagram) |
| `{output}-report.md` | Phase 4 | No (verification report) |

The intermediate `.flow-deps.yaml` file is automatically deleted after successful completion unless `--keep-deps` is specified.

## Supported Languages

- Go
- Python

---

## Pipeline Overview

```
Phase 1          Phase 2          Phase 3          Phase 4
Discovery   →    Generation  →    Refinement  →    Verification
    ↓                ↓                ↓                ↓
dependencies    diagram.md      diagram.md       PASS/FAIL
   .yaml         (draft)        (refined)
```

**Key Principle**: Each phase has ONE job. No phase does another's work.

---

## Phase 1: Discovery

**Goal**: Extract ALL dependencies into `dependencies.yaml`

**Instructions**: Follow `prompts/phase1-discovery.md`

### Critical Rules

1. **NEVER grep/search the entire codebase**
2. **Read runbooks/docs FIRST** - they list dependencies explicitly
3. **Only read files in the target directory**
4. **Every dependency MUST have source_file and source_line**
5. **No abbreviations** - use full names

### Process

1. Confirm target directory path (ask if not provided)
2. List directory structure
3. Read documentation files (README, RUNBOOK, docs/)
4. Read configuration files (*.yaml, *.json)
5. Read client files (*_client.go) - outbound gRPC/HTTP
6. Read consumer files (*_consumer.go) - inbound Kafka
7. Read producer files (*_producer.go) - outbound Kafka
8. Read proto files (*.proto) - gRPC definitions
9. Read repository files (*_repo.go) - database
10. Write `.flow-deps.yaml` to working directory

### Output

```yaml
version: "1.0"
generated: "2024-01-15T14:30:22Z"
target_service: "ledger-service"
target_path: "services/ledger/"

service:
  name: "Ledger Service"
  description: "..."

dependencies:
  sync:
    - name: "Payment Service"
      type: "grpc"
      source_file: "internal/client/payment_client.go"
      source_line: 45
  async:
    - name: "order.completed"
      type: "kafka"
      direction: "consume"
      source_file: "internal/consumer/order_consumer.go"
      source_line: 27

external:
  - name: "Stripe API"
    source_file: "internal/stripe/client.go"
    source_line: 23

caches:
  - name: "Redis"
    source_file: "internal/cache/balance_cache.go"
    source_line: 15
```

### Completion

Report:
```
Phase 1 Complete: Discovery

Output: .flow-deps.yaml

Summary:
- Service: {name}
- Sync dependencies: {count}
- Async dependencies: {count}
- External systems: {count}
- Caches: {count}

Proceeding to Phase 2...
```

---

## Phase 2: Generation

**Goal**: Generate Mermaid diagram from `.flow-deps.yaml`

**Instructions**: Follow `prompts/phase2-generation.md`

### Critical Rules

1. **DO NOT read the codebase** - only read .flow-deps.yaml
2. **Follow the style guide exactly** - see `styles/diagram-styles.yaml`
3. **No abbreviations** - use full names from dependencies.yaml
4. **Correct arrows** - `==>` sync, `-.->` async
5. **Quote subgraph titles** - `subgraph id ["Title"]`

### Style Reference

**Colors** (non-negotiable):
```
classDef service fill:#228be6,stroke:#1971c2,color:#fff
classDef entry fill:#40c057,stroke:#2f9e44,color:#fff
classDef kafka fill:#12b886,stroke:#099268,color:#fff
classDef database fill:#fab005,stroke:#f59f00,color:#000
classDef cache fill:#be4bdb,stroke:#9c36b5,color:#fff
classDef external fill:#868e96,stroke:#495057,color:#fff
```

**Shapes**:
| Type | Shape | Syntax |
|------|-------|--------|
| Service | Rectangle | `[Service Name]` |
| Kafka Topic | Cylinder | `[(topic.name)]` |
| Database | Cylinder | `[(PostgreSQL)]` |
| External | Double Rectangle | `[[Stripe API]]` |
| Cache | Rounded | `(Redis Cache)` |

**Arrows**:
| Type | Arrow | Use For |
|------|-------|---------|
| Sync | `==>` | gRPC, HTTP, SQL |
| Async | `-.->` | Kafka, queues |
| Internal | `-->` | Function calls |

### Output

Write `{output}.md` (default: `flow-diagram.md`) using template from `templates/diagram-template.md`

### Completion

Report:
```
Phase 2 Complete: Generation

Output: {output}.md

Diagram contains:
- Services: {count}
- Kafka topics: {count}
- Databases: {count}
- External systems: {count}

Proceeding to Phase 3...
```

---

## Phase 3: Refinement

**Goal**: Validate and fix the diagram using `flowlint`

### Prerequisites

Build flowlint (one-time):
```bash
cd flow/tools/flowlint
go build -o flowlint .
```

Install Mermaid CLI (for validation):
```bash
npm install -g @mermaid-js/mermaid-cli
```

### Commands

```bash
# Run full refinement
flowlint refine diagram.md dependencies.yaml --output diagram-final.md

# Or run steps individually:
flowlint validate diagram.md          # Check syntax
flowlint lint diagram.md --fix        # Fix style issues
flowlint check diagram.md deps.yaml   # Check completeness
```

### What flowlint Checks

1. **Syntax** - Mermaid compiles without errors
2. **Style** - Correct arrows, quoted subgraphs, classDefs
3. **Completeness** - All dependencies represented
4. **No orphans** - All nodes connected

### Skip Phase 3

If flowlint not available, use `--no-tools` flag:
```
/flow services/ledger/ --no-tools
```

### Completion

Report:
```
Phase 3 Complete: Refinement

Applied fixes: {count}
Remaining issues: {count}

Proceeding to Phase 4...
```

---

## Phase 4: Verification

**Goal**: Final verification that diagram is complete and correct

**Instructions**: Follow `prompts/phase4-verification.md`

### Checks

1. All dependencies from YAML appear in diagram
2. Arrow styles match communication types
3. Style classes applied correctly
4. No orphan nodes
5. Legend present

### Output

Write `{output}-report.md` (default: `flow-diagram-report.md`):

```markdown
# Verification Report

## Status: PASS / FAIL

## Coverage
| Category | Expected | Found | Missing |
|----------|----------|-------|---------|
| Services | 5 | 5 | 0 |
| Kafka | 3 | 3 | 0 |
| Databases | 1 | 1 | 0 |

## Issues Found
- None

## Render Commands
mmdc -i diagram.md -o diagram.png -b white
```

### Cleanup

After successful verification (Status: PASS), **delete the intermediate file**:

```bash
rm .flow-deps.yaml
```

Skip cleanup if `--keep-deps` was specified.

### Completion

```
Phase 4 Complete: Verification

Status: PASS

All 11 dependencies represented.
No style issues.

Cleaned up: .flow-deps.yaml

Output files:
- {output}.md (diagram)
- {output}-report.md (verification)

Render with:
mmdc -i {output}.md -o diagram.png -b white
```

---

## Example Usage

```bash
# Full pipeline (outputs flow-diagram.md)
/flow services/ledger/

# Custom output name (outputs ledger-arch.md)
/flow services/ledger/ --output=ledger-arch

# Stop after discovery (inspect .flow-deps.yaml)
/flow services/ledger/ --phase=1

# Skip tooling validation
/flow services/ledger/ --no-tools

# Keep intermediate dependencies file
/flow services/ledger/ --keep-deps
```

---

## File Structure

```
flow/
├── SKILLS.md                 # This file
├── PLAN.md                   # Implementation plan
├── schemas/
│   └── dependencies.schema.yaml
├── styles/
│   └── diagram-styles.yaml
├── prompts/
│   ├── phase1-discovery.md
│   ├── phase2-generation.md
│   └── phase4-verification.md
├── templates/
│   └── diagram-template.md
├── tools/
│   └── flowlint/             # Go CLI tool
└── examples/
    ├── dependencies-example.yaml
    └── diagram-example.md
```

---

## Troubleshooting

### "Context window exceeded"
- You searched the entire codebase. **DO NOT DO THIS.**
- Start over and only read files in the target directory.
- Read runbooks first to find dependency information.

### "Diagram has syntax errors"
- Check subgraph titles are quoted: `subgraph id ["Title"]`
- Check arrow syntax: `==>`, `-.->`, `-->`
- Run `flowlint validate diagram.md` for details.

### "Missing dependencies"
- Run `flowlint check diagram.md dependencies.yaml`
- Add missing nodes to the diagram
- Re-run Phase 4 verification

### "Wrong arrow style"
- Sync calls (gRPC, HTTP, SQL) use `==>`
- Async calls (Kafka, queues) use `-.->`
- Run `flowlint lint diagram.md --fix` to auto-fix
