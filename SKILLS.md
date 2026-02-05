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
| `--no-tools` | Skip flowlint validation (if not installed) | false |
| `--keep-deps` | Keep .flow-deps.yaml after completion | false |

## Output Files

All files are written to the **current working directory**:

| File | Phase | Description |
|------|-------|-------------|
| `{output}.md` | Phase 2 | **Always created** - the diagram source |
| `{output}.png` | Phase 3 | Rendered image (if mmdc installed) |
| `.flow-deps.yaml` | Phase 1 | Intermediate file, deleted on success |

- **Diagram is always generated**, even if validation has issues
- PNG is generated if `mmdc` (Mermaid CLI) is installed
- `.flow-deps.yaml` is deleted on success, kept on failure (for debugging)
- Use `--keep-deps` to always keep the intermediate file

## Supported Languages

- Go
- Python

---

## Pipeline Overview

```
Phase 1              Phase 2              Phase 3
Discovery    →       Generation   →       Validation + Render + Cleanup
   ↓                     ↓                     ↓
.flow-deps.yaml    flow-diagram.md      Run flowlint
                   (ALWAYS created)       ↓
                                        PASS → render PNG → cleanup
                                          ↓
                                        FAIL → AI reworks → re-validate
                                          ↓
                                        Still FAIL → warn, keep .flow-deps.yaml
```

**Key Principle**:
- Each phase has ONE job
- **Diagram is ALWAYS generated** (Phase 2)
- PNG rendered if mmdc installed
- Validation issues = warnings, not blockers

---

## Phase 1: Discovery

**Goal**: Extract ALL dependencies into `.flow-deps.yaml`

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
3. **No abbreviations** - use full names from .flow-deps.yaml
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

## Phase 3: Validation & Cleanup

**Goal**: Validate diagram, rework if needed, cleanup

### Prerequisites

Build flowlint (one-time):
```bash
cd flow/tools/flowlint
go build -o flowlint .
```

Install Mermaid CLI (optional, for syntax validation):
```bash
npm install -g @mermaid-js/mermaid-cli
```

### Step 3.1: Run flowlint

**YOU MUST EXECUTE THIS COMMAND VIA BASH:**

```bash
./flow/tools/flowlint/flowlint refine {output}.md .flow-deps.yaml --skip-validate
```

Or if flowlint is in PATH:
```bash
flowlint refine {output}.md .flow-deps.yaml --skip-validate
```

**DO NOT SKIP THIS STEP. Actually run the command and check the output.**

flowlint will:
1. Check style guide compliance and auto-fix
2. Verify all dependencies are represented
3. Report PASS or FAIL

### Step 3.2: Handle Results

**Read the flowlint output carefully.**

#### If flowlint output shows PASS (exit code 0):
```
✓ Refinement complete - diagram is ready
```

Then cleanup:
```bash
rm .flow-deps.yaml
```

Done.

#### If flowlint output shows FAIL (exit code 1):

**Read the error messages.** Then attempt **one rework cycle**:

| Failure Type | AI Action |
|--------------|-----------|
| Syntax error | Edit {output}.md to fix the Mermaid syntax |
| Missing dependencies | Edit {output}.md to add missing nodes |
| Style issues | Edit {output}.md to fix arrows/shapes/colors |

After editing the diagram, re-run flowlint:
```bash
./flow/tools/flowlint/flowlint refine {output}.md .flow-deps.yaml --skip-validate
```

#### If still FAIL after rework:
```
⚠️ Validation incomplete - diagram generated with issues.

Output: {output}.md (diagram created, may have issues)
Kept: .flow-deps.yaml (for debugging)

Issues:
- {list of remaining issues}

The diagram is usable but may need manual fixes.
Run: flowlint lint {output}.md --fix
```

**The diagram is ALWAYS generated.** Validation failure just means it may not be perfect.

### Finding flowlint

flowlint is located at:
```
flow/tools/flowlint/flowlint
```

If you get "command not found", use the full path relative to the skill directory.

### Skip Validation

If flowlint not available or not built:
```
/flow services/ledger/ --no-tools
```

This skips validation entirely - use with caution.

### Step 3.3: Generate PNG (if mmdc available)

After validation passes, generate a PNG:

```bash
mmdc -i {output}.md -o {output}.png -b white -w 1920
```

If mmdc is not installed, skip this step and inform the user how to render manually.

### Completion (Success)

```
Phase 3 Complete: Validation

Status: PASS
Coverage: 10/10 (100%)

Cleaned up: .flow-deps.yaml

Output:
- {output}.md (diagram source)
- {output}.png (rendered image, if mmdc available)

If PNG not generated, render manually:
mmdc -i {output}.md -o {output}.png -b white
```

---

## Example Usage

```bash
# Full pipeline (outputs flow-diagram.md)
/flow services/ledger/

# Custom output name (outputs ledger-arch.md)
/flow services/ledger/ --output=ledger-arch

# Skip validation (if flowlint not installed)
/flow services/ledger/ --no-tools

# Keep intermediate dependencies file
/flow services/ledger/ --keep-deps
```

---

## File Structure

```
flow/
├── SKILLS.md                 # This file
├── schemas/
│   └── dependencies.schema.yaml
├── styles/
│   └── diagram-styles.yaml
├── prompts/
│   ├── phase1-discovery.md
│   └── phase2-generation.md
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
- Run `flowlint validate {output}.md` for details.

### "Missing dependencies"
- Run `flowlint check {output}.md .flow-deps.yaml`
- Add missing nodes to the diagram
- Re-run flowlint

### "Wrong arrow style"
- Sync calls (gRPC, HTTP, SQL) use `==>`
- Async calls (Kafka, queues) use `-.->`
- Run `flowlint lint {output}.md --fix` to auto-fix
