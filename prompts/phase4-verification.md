# Phase 4: Verification

## Objective

Verify the final diagram is complete, accurate, and ready for documentation.

## Input

- `{output}.md` (refined from Phase 3, default: `flow-diagram.md`)
- `.flow-deps.yaml` (from Phase 1)

## Output

- `{output}-report.md` (default: `flow-diagram-report.md`)
- Final status: PASS or FAIL
- Cleanup: Delete `.flow-deps.yaml` on success (unless `--keep-deps`)

## Process

### Step 1: Read Both Files

1. Read `{output_dir}/dependencies.yaml` - source of truth
2. Read `{output_dir}/diagram.md` - diagram to verify

### Step 2: Extract Entities from Dependencies

Parse dependencies.yaml and list:

**Services:**
- Target service name
- All sync dependency names
- All caller names

**Kafka Topics:**
- All produced topics
- All consumed topics

**Data Stores:**
- All databases
- All caches

**External Systems:**
- All external API names

### Step 3: Extract Entities from Diagram

Parse the Mermaid diagram and extract:

**Nodes:**
- All node IDs and labels
- Node types (from shapes)

**Edges:**
- All connections (source → target)
- Edge types (sync `==>`, async `-.->`)

**Subgraphs:**
- All subgraph names
- Contents of each subgraph

### Step 4: Compare Coverage

For each entity in dependencies.yaml, verify it appears in diagram:

```
Coverage Check:

Services:
✓ Ledger Service - found as node T1
✓ Payment Service - found as node D1
✗ Inventory Service - MISSING

Kafka Topics:
✓ order.created - found as node K1
✓ ledger.transaction.created - found as node K2
✗ payment.refunded - MISSING

Databases:
✓ PostgreSQL - found as node DB1

External:
✓ Stripe API - found as node X1
```

### Step 5: Verify Arrow Styles

Check that connection types match:

```
Arrow Style Check:

✓ T1 ==> D1 (gRPC call uses sync arrow)
✓ T2 -.-> K1 (Kafka produce uses async arrow)
✗ T2 --> DB1 (Database should use ==> not -->)
```

### Step 6: Check Style Compliance

Verify against `styles/diagram-styles.yaml`:

```
Style Compliance:

classDef declarations:
✓ service class defined
✓ kafka class defined
✓ database class defined
✓ external class defined

Class applications:
✓ All service nodes have 'service' class
✗ K2 missing 'kafka' class

Subgraph titles:
✓ All titles are quoted

Node naming:
✓ No abbreviations found
✗ "PaySvc" should be "Payment Service"
```

### Step 7: Check Layout Quality

Visual quality checks:

```
Layout Quality:

✓ Entry points on left
✓ External systems on right
✓ Logical grouping in subgraphs
✗ Too many nodes in "deps" subgraph (consider splitting)
✗ Orphan node found: X2 has no connections
```

### Step 8: Generate Verification Report

Write `{output_dir}/verification-report.md`:

```markdown
# Verification Report

Generated: {timestamp}
Service: {service_name}

## Status: PASS / FAIL

## Coverage Summary

| Category | Expected | Found | Missing |
|----------|----------|-------|---------|
| Services | 5 | 5 | 0 |
| Kafka Topics | 3 | 2 | 1 |
| Databases | 1 | 1 | 0 |
| External Systems | 2 | 2 | 0 |
| **Total** | **11** | **10** | **1** |

## Missing Items

| Item | Type | Source Reference |
|------|------|------------------|
| payment.refunded | Kafka Topic | consumer/refund_consumer.go:15 |

## Arrow Style Issues

| Connection | Current | Expected | Fix |
|------------|---------|----------|-----|
| T2 --> DB1 | internal | sync (==>) | Change --> to ==> |

## Style Issues

| Issue | Location | Recommendation |
|-------|----------|----------------|
| Missing class | Node K2 | Add `class K2 kafka` |
| Abbreviation | Node label "PaySvc" | Change to "Payment Service" |

## Layout Issues

| Issue | Recommendation |
|-------|----------------|
| Orphan node X2 | Add connection or remove |
| Large subgraph "deps" | Split into multiple subgraphs |

## Recommendations

1. Add missing Kafka topic `payment.refunded`
2. Fix arrow style for database connection
3. Apply kafka class to node K2
4. Rename "PaySvc" to "Payment Service"

## Final Checklist

- [x] All services represented
- [ ] All Kafka topics represented (1 missing)
- [x] All databases represented
- [x] All external systems represented
- [ ] All sync calls use ==> arrows (1 incorrect)
- [x] All async calls use -.-> arrows
- [x] All subgraph titles quoted
- [ ] All nodes have correct class (1 missing)
- [ ] No abbreviations (1 found)
- [x] No orphan nodes
- [x] Legend included

## Render Command

To generate PNG:
\`\`\`bash
mmdc -i diagram.md -o diagram.png -b white -w 1920
\`\`\`

To generate SVG:
\`\`\`bash
mmdc -i diagram.md -o diagram.svg -b white
\`\`\`
```

## Decision Logic

### PASS Criteria (ALL must be true)

1. 100% coverage - all dependencies in diagram
2. No arrow style errors
3. No missing class applications
4. No abbreviations
5. No orphan nodes
6. Legend present

### FAIL Criteria (ANY triggers fail)

1. Missing service, topic, database, or external system
2. Wrong arrow style (sync shown as async or vice versa)
3. Unquoted subgraph titles
4. Critical style violations

## Cleanup

**On PASS**: Delete the intermediate `.flow-deps.yaml` file:
```bash
rm .flow-deps.yaml
```

**Skip cleanup if**: `--keep-deps` was specified or status is FAIL.

## Output

After completing Phase 4, report:

```
Phase 4 Complete: Verification

Status: PASS / FAIL

Coverage: {found}/{expected} ({percentage}%)

Issues found: {count}
- Critical: {count}
- Warnings: {count}

Report: {output}-report.md
Diagram: {output}.md

{if PASS}
Cleaned up: .flow-deps.yaml
Diagram is ready for documentation.
{endif}

{if FAIL}
Kept: .flow-deps.yaml (for debugging)
Please review the verification report and fix the issues.
Re-run Phase 2 or manually edit diagram.md.
{endif}
```
