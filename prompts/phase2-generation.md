# Phase 2: Diagram Generation

## Objective

Generate a professional Mermaid diagram from `.flow-deps.yaml`.

**CRITICAL: In this phase, you read ONLY the .flow-deps.yaml file. Do NOT read the codebase.**

## Input

- `.flow-deps.yaml` from Phase 1 (in working directory, `services[]` format)

## Output

- `{output}.md` - Mermaid diagram with metadata (default: `flow-diagram.md`)

---

## Process

### Step 1: Read Dependencies File

Read `.flow-deps.yaml` and parse:
- All services in the `services[]` array
- Each service's sync/async dependencies, databases, caches, external systems
- Internal steps (if present)

### Step 2: Plan the Diagram Layout

**Flow direction: Always `flowchart TD` (top-down).**

For **single-service** diagrams:
```
                    TOP
                     ↓
            [Entry Points]
                     ↓
             [Target Service]
                     ↓
    [Sync Dependencies]  [Data Stores]
                     ↓
    [Consumed Topics]    [Produced Topics]
                     ↓
            [External Systems]
                     ↓
                   BOTTOM
```

All groups are scoped to THIS service. See Step 6.5 for grouping rules.

For **multi-service** diagrams:
- Each service becomes its own subgraph
- Use node ID prefixes (`S1_`, `S2_`, `S3_`) to avoid collisions
- Inter-service edges: if Service A depends on "Payment Service" and there's a Service B named "Payment Service", draw an edge between them
- Only share a subgraph for infrastructure genuinely used by multiple services (same topic, same DB)
- Maximum 3 services per diagram — suggest splitting if more

For **pipeline/journey** diagrams:
- Start at TOP, end at BOTTOM
- Add explicit markers: `START([Start: Request Received])` and `END([End: Response Sent])`
- Use dedicated classDef: `classDef startEnd fill:#d4edda,stroke:#28a745,color:#155724,stroke-width:2px`
- Verify visual reading order matches logical flow: top → bottom
- Never place entry point below its downstream dependencies

### Step 3: Generate Mermaid Code

Follow the style guide in `styles/diagram-styles.yaml` EXACTLY.

#### Style Definitions (REQUIRED)

```
classDef service fill:#a5d8ff,stroke:#339af0,color:#1864ab
classDef entry fill:#b2f2bb,stroke:#51cf66,color:#2b8a3e
classDef kafka fill:#96f2d7,stroke:#38d9a9,color:#087f5b
classDef database fill:#ffec99,stroke:#fcc419,color:#e67700
classDef cache fill:#d0bfff,stroke:#9775fa,color:#6741d9
classDef external fill:#dee2e6,stroke:#adb5bd,color:#495057
classDef step fill:#e8f4f8,stroke:#4a9ebb,color:#2c5f7c
classDef startEnd fill:#d4edda,stroke:#28a745,color:#155724,stroke-width:2px
```

### Step 4: Node Naming Rules

| Type | Format | Example |
|------|--------|---------|
| Service | Title Case | `Payment Service` |
| Kafka Topic | lowercase.dot | `order.created` |
| Database | Title Case | `Ledger DB` |
| Cache | Title Case | `Ledger Cache` |
| External | Title Case | `Stripe API` |
| Internal Step | Title Case | `Validate Request` |

**CRITICAL LABEL RULES:**
- **NO NEWLINES** — all labels must be single-line
- **NO ABBREVIATIONS** — use full names always

```
❌ WRONG:
    A[Commit Stage
    (Write)]

✅ CORRECT:
    A[Commit Stage Write]
```

### Step 5: Arrow Rules

| Connection Type | Arrow | Label |
|-----------------|-------|-------|
| gRPC call | `==>` | `\|MethodName\|` (only if useful) |
| HTTP call | `==>` | (no label) |
| Database | `==>` | (no label — the node shape says it all) |
| Cache | `==>` | (no label) |
| Kafka produce | `-.->` | (no label — dotted arrow + cylinder = obvious) |
| Kafka consume | `-.->` | (no label) |
| Internal call | `-->` | (no label) |
| Internal step chain | `-->` | (no label) |

**Labels are for information the diagram doesn't already show.** The arrow style (`==>` vs `-.->`) shows sync vs async. The node shape (cylinder, rounded) shows the target type. Don't repeat what's already visible.

Only add a label when it provides specific context:
- `\|ProcessPayment\|` — tells you which RPC method (useful)
- `\|SQL\|` on an arrow to a database cylinder — redundant, skip it

### Step 6: Subgraph Rules

1. **Always quote titles**: `subgraph id ["Display Name"]`
2. **Use lowercase IDs**: `subgraph kafka-in ["Consumed Topics"]`
3. **Max 4 nodes per subgroup** — if more, split into multiple subgroups stacked vertically
4. **Arrows go to subgroups, NOT individual nodes** — one arrow to a subgroup covers all nodes inside

### Step 6.5: Grouping Rules — Group by Service Context, NOT by Type

**⛔ DO NOT group all topics together, all databases together, etc.**

Group by what's relevant to the same service and the same direction:

| Group | Contains | Example Title |
|-------|----------|---------------|
| Service's consumed topics | Kafka topics this service consumes | `"Consumed Topics"` |
| Service's produced topics | Kafka topics this service produces | `"Produced Topics"` |
| Service's sync dependencies | gRPC/HTTP services it calls | `"Dependencies"` |
| Service's data stores | Its database + cache | `"Data Stores"` |
| Service's external systems | Third-party APIs it calls | `"External"` |

```
❌ WRONG — grouping all topics from all services:
    subgraph topics ["All Kafka Topics"]
        K1[(order.created)]    %% produced by Order Service
        K2[(payment.done)]     %% consumed by Order Service
        K3[(ledger.updated)]   %% produced by Ledger Service
    end

✅ RIGHT — grouping by service + direction:
    subgraph order-out ["Order Service: Produced"]
        K1[(order.created)]
    end
    subgraph order-in ["Order Service: Consumed"]
        K2[(payment.done)]
    end
```

### Step 6.6: Subgroup Sizing and Stacking

**Max 4 nodes per subgroup.** If a service has 7 sync dependencies, split into two subgroups stacked vertically:

```mermaid
flowchart TD
    subgraph target ["Order Service"]
        T1[Handler]
    end

    subgraph deps1 ["Dependencies"]
        D1[Payment Service]
        D2[Inventory Service]
        D3[Auth Service]
    end

    subgraph deps2 ["Dependencies"]
        D4[Notification Service]
        D5[Audit Service]
    end

    subgraph data ["Data Stores"]
        DB1[(Order DB)]
        C1(Order Cache)
    end

    T1 ==> deps1
    T1 ==> deps2
    T1 ==> data
```

**Key rules:**
- Subgroups stack vertically (default `TD` flow)
- **One arrow per subgroup** — NOT one arrow per node inside
- Nodes inside a subgroup that receives an arrow are NOT orphans

---

## Internal Steps Rendering

When a service has `internal_steps`, show them as connected nodes inside its subgraph.

### When to Include Internal Steps

- **Include** when: single-service diagram, OR user asks for "detailed" diagram
- **Skip** for multi-service diagrams unless explicitly requested — keeps things clean

### How to Render

Steps are chained with `-->` arrows in order inside the service subgraph. **All subgroup rules STILL apply** — dependencies go into subgroups, and arrows from steps go to subgroups.

```
╔══════════════════════════════════════════════════════════════════╗
║  ⛔ DO NOT draw arrows from internal steps to individual nodes  ║
║  ✅ Draw arrows from the relevant step to the SUBGROUP          ║
╚══════════════════════════════════════════════════════════════════╝
```

```mermaid
flowchart TD
    subgraph service ["Order Service"]
        direction TB
        S1_step1[Validate Request] --> S1_step2[Fetch Data]
        S1_step2 --> S1_step3[Process Payment]
        S1_step3 --> S1_step4[Commit]
    end

    subgraph deps ["Dependencies"]
        D1[Payment Service]
        D2[Inventory Service]
    end

    subgraph data ["Data Stores"]
        DB1[(Order DB)]
        C1(Order Cache)
    end

    subgraph kafka-out ["Produced Topics"]
        KO1[(order.created)]
    end

    S1_step2 ==> data
    S1_step3 ==> deps
    S1_step4 -.-> kafka-out
```

```
❌ WRONG — arrows to individual nodes (creates arrow explosion):
    S1_step2 ==> DB1[(Order DB)]
    S1_step2 ==> C1(Order Cache)
    S1_step3 ==> D1[Payment Service]
    S1_step3 ==> D2[Inventory Service]
    S1_step4 -.-> KO1[(order.created)]
    — 5 arrows for 5 nodes

✅ RIGHT — arrows to subgroups (clean, minimal):
    S1_step2 ==> data
    S1_step3 ==> deps
    S1_step4 -.-> kafka-out
    — 3 arrows covering all nodes
```

**Rules:**
- Steps use `classDef step` styling
- Arrows from steps go to **subgroups** (same rule as Step 6, rule 4 — no exceptions)
- Each step connects to the subgroup(s) it interacts with
- When no `internal_steps`, service remains a single node (default behavior)

---

## Complete Single-Service Example

This shows the Ledger Service (from the example YAML) with ALL rules applied: subgroups, service-context grouping, internal steps, arrows to subgroups only.

**Input:** Ledger Service with 2 sync deps, 2 consumed topics, 2 produced topics, 1 DB, 1 cache, 2 external, 4 internal steps.

```mermaid
flowchart TD
    subgraph entry ["Entry Points"]
        E1([LedgerService])
        E2([/api/v1/ledger])
    end

    subgraph service ["Ledger Service"]
        direction TB
        S1_step1[Validate Request] --> S1_step2[Load Account]
        S1_step2 --> S1_step3[Process Transaction]
        S1_step3 --> S1_step4[Commit]
    end

    subgraph deps ["Dependencies"]
        D1[Payment Service]
        D2[Account Service]
    end

    subgraph data ["Data Stores"]
        DB1[(Ledger DB)]
        C1(Ledger Cache)
    end

    subgraph kafka-in ["Consumed Topics"]
        KI1[(order.completed)]
        KI2[(payment.refunded)]
    end

    subgraph kafka-out ["Produced Topics"]
        KO1[(ledger.transaction.created)]
        KO2[(ledger.balance.updated)]
    end

    subgraph ext ["External"]
        EX1[Stripe API]
        EX2[Audit Service]
    end

    entry ==> service
    kafka-in -.-> service
    S1_step2 ==> data
    S1_step3 ==> deps
    S1_step4 -.-> kafka-out
    S1_step4 ==> ext

    class E1,E2 entry
    class S1_step1,S1_step2,S1_step3,S1_step4 step
    class D1,D2 service
    class DB1 database
    class C1 cache
    class KI1,KI2,KO1,KO2 kafka
    class EX1,EX2 external
```

**Result: 6 arrows** connect all 15 nodes. Without subgroups this would be 12+ arrows.

---

## Multi-Service Layout

### Node ID Prefixes

Each service gets a prefix to avoid ID collisions:
- Service 1: `S1_handler`, `S1_db`, `S1_step1`
- Service 2: `S2_handler`, `S2_db`, `S2_step1`

### Inter-Service Edges

If Service A has a sync dependency named "Payment Service" and Service B is named "Payment Service", draw:
```
S1_handler ==> S2_handler
```

### Shared Infrastructure

Only group infrastructure as "shared" if the SAME topic/DB is used by multiple services. Otherwise keep it in the service's own subgroup.

---

## Step 7: Write Complete Markdown File

Use the template from `templates/diagram-template.md`:

```markdown
# Service Flow: {service_name}

Generated: {timestamp}
Source: {target_path}

## Diagram

\`\`\`mermaid
{mermaid_code}
\`\`\`

## Legend

| Symbol | Meaning |
|--------|---------|
| `==>` | **Synchronous** (gRPC/HTTP) |
| `-.->` | **Asynchronous** (Kafka) |
| `-->` | Internal call / step chain |

### Colors

| Color | Meaning |
|-------|---------|
| Blue | Services |
| Green | Entry Points |
| Teal | Kafka Topics |
| Yellow | Databases |
| Purple | Caches |
| Gray | External Systems |
| Light Blue | Internal Steps |

## Dependencies

{dependency_tables}

## Source References

{source_references}
```

---

## Quality Checklist

Before completing Phase 2, verify:

- [ ] All services from .flow-deps.yaml are in the diagram
- [ ] All Kafka topics, databases, caches, external systems are in the diagram
- [ ] Sync calls use `==>`, async use `-.->`, internal use `-->`
- [ ] All subgraph titles are quoted
- [ ] Max 4 nodes per subgroup — split and stack vertically if more
- [ ] Arrows go to subgroups, not individual nodes inside them
- [ ] Grouped by service context (consumed topics, produced topics, deps), NOT by type
- [ ] No redundant labels on arrows
- [ ] No abbreviations in node names
- [ ] classDef styles applied to all nodes
- [ ] Internal steps: arrows go to subgroups, NOT individual nodes (if present)
- [ ] No arrow explosion — count total arrows, compare to subgroup count (should be similar)

## Output

```
Phase 2 Complete: Diagram Generation

Output: {output}.md

Diagram contains:
- Services: {count}
- Kafka topics: {count}
- Databases: {count}
- External systems: {count}
- Internal steps: {count}

Proceeding to Phase 3: Validation
```
