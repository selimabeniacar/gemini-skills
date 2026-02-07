# /flow

Generate professional architectural diagrams through a multi-phase pipeline.

---

## Usage

```
/flow <path> [path2] [options]
```

### Options

| Option | Description | Default |
|--------|-------------|---------|
| `--output=<name>` | Output filename (without extension) | `flow-diagram` |
| `--keep-deps` | Keep .flow-deps.yaml after completion | false |

### Output Files

| File | Phase | Description |
|------|-------|-------------|
| `{output}.md` | Phase 2 | **Always created** - the diagram source |
| `{output}.png` | Phase 3 | Rendered image (if mermaid-cli installed) |
| `.flow-deps.yaml` | Phase 1 | Intermediate file, deleted on success |

### Supported Languages

- Go
- Python

---

## Pipeline Overview

```
Phase 1              Phase 2              Phase 3
Discovery    →       Generation   →       Validation + Render
     ↓                    ↓                     ↓
Read code            Read YAML only         Run flowlint
Coverage review      Generate Mermaid        ↓
     ↓               Write {output}.md      PASS → PNG → cleanup
.flow-deps.yaml                             FAIL → rework → retry
                                            Still FAIL → warn
```

---

## Phase 1: Discovery

**Goal**: Extract ALL dependencies into `.flow-deps.yaml`

**Instructions**: Follow `prompts/phase1-discovery.md`

Includes mandatory coverage self-review: every non-test file in `{target_path}` must be accounted for before completing this phase.

---

## Phase 2: Generation

**Goal**: Generate Mermaid diagram from `.flow-deps.yaml`

**Instructions**: Follow `prompts/phase2-generation.md`

**DO NOT read the codebase** — only read `.flow-deps.yaml`.

---

## Phase 3: Validation & Cleanup

**Goal**: Validate diagram, rework if needed, render PNG, cleanup

**Instructions**: Follow `prompts/phase3-validation.md`

---

## Examples

### Single Service
```bash
/flow services/ledger/
```

### Multi-Service Pipeline
```bash
/flow services/order/ services/payment/ --output=order-payment-flow
```

### Custom Output
```bash
/flow services/ledger/ --output=ledger-arch --keep-deps
```

---

## File Structure

```
flow/
├── SKILLS.md                        # This file (pipeline router)
├── schemas/
│   └── dependencies.schema.yaml     # Schema definition
├── styles/
│   └── diagram-styles.yaml          # Color/shape/arrow reference
├── prompts/
│   ├── phase1-discovery.md          # Phase 1: Discovery + Coverage review
│   ├── phase2-generation.md         # Phase 2: Diagram generation
│   └── phase3-validation.md         # Phase 3: Validation + Render
├── templates/
│   └── diagram-template.md          # Markdown output template
├── tools/
│   └── flowlint/                    # Go CLI tool
└── examples/
    ├── dependencies-example.yaml    # Example YAML
    └── diagram-example.md           # Example diagram output
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

### "No services found"
- Ensure `.flow-deps.yaml` uses `services[]` array format
- See `schemas/dependencies.schema.yaml` for the expected structure
