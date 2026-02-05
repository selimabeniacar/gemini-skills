# Phase 2: Diagram Generation

## Objective

Generate a professional Mermaid diagram from `dependencies.yaml`.

**CRITICAL: In this phase, you read ONLY the dependencies.yaml file. Do NOT read the codebase.**

## Input

- `.flow-deps.yaml` from Phase 1 (in working directory)

## Output

- `{output}.md` - Mermaid diagram with metadata (default: `flow-diagram.md`)

## Process

### Step 1: Read Dependencies File

Read `{output_dir}/dependencies.yaml` and parse:
- Service metadata
- All sync dependencies
- All async dependencies
- External systems
- Databases and caches

### Step 2: Plan the Diagram Layout

Based on dependencies, plan node placement:

```
LEFT → → → → → → → → → → → → → → → → → → → → RIGHT

[Entry Points] → [Target Service] → [Dependencies] → [Message Bus] → [Data/External]
```

Group nodes:
1. **Entry Points** - How traffic enters (API gateway, load balancer, direct gRPC)
2. **Target Service** - The service being documented (handlers, processors)
3. **Dependent Services** - Services this one calls synchronously
4. **Message Bus** - Kafka topics
5. **Data Stores** - Databases, caches
6. **External Systems** - Third-party APIs

### Step 3: Generate Mermaid Code

Follow the style guide in `styles/diagram-styles.yaml` EXACTLY.

#### Mermaid Structure

```mermaid
flowchart LR
    %% ========================================
    %% Style Definitions (REQUIRED)
    %% ========================================
    classDef service fill:#228be6,stroke:#1971c2,color:#fff
    classDef entry fill:#40c057,stroke:#2f9e44,color:#fff
    classDef kafka fill:#12b886,stroke:#099268,color:#fff
    classDef database fill:#fab005,stroke:#f59f00,color:#000
    classDef cache fill:#be4bdb,stroke:#9c36b5,color:#fff
    classDef external fill:#868e96,stroke:#495057,color:#fff

    %% ========================================
    %% Entry Points (leftmost)
    %% ========================================
    subgraph entry ["Entry Points"]
        E1[gRPC Server]
        E2[HTTP Server]
    end

    %% ========================================
    %% Target Service
    %% ========================================
    subgraph target ["Service Name"]
        T1[Handler]
        T2[Processor]
    end

    %% ========================================
    %% Dependent Services
    %% ========================================
    subgraph deps ["Dependencies"]
        D1[Payment Service]
        D2[Inventory Service]
    end

    %% ========================================
    %% Message Bus
    %% ========================================
    subgraph kafka ["Message Bus"]
        K1[(order.created)]
        K2[(payment.completed)]
    end

    %% ========================================
    %% Data Stores
    %% ========================================
    subgraph data ["Data Stores"]
        DB1[(PostgreSQL)]
        C1(Redis Cache)
    end

    %% ========================================
    %% External Systems (rightmost)
    %% ========================================
    subgraph ext ["External Systems"]
        X1[[Stripe API]]
    end

    %% ========================================
    %% Connections - Sync (thick arrows)
    %% ========================================
    E1 ==>|gRPC| T1
    T1 --> T2
    T2 ==>|gRPC: ProcessPayment| D1
    T2 ==>|SQL| DB1
    T2 ==>|cache| C1
    D1 ==>|HTTPS| X1

    %% ========================================
    %% Connections - Async (dotted arrows)
    %% ========================================
    T2 -.->|publish| K1
    K2 -.->|consume| T1

    %% ========================================
    %% Apply Styles
    %% ========================================
    class E1,E2 entry
    class T1,T2 service
    class D1,D2 service
    class K1,K2 kafka
    class DB1 database
    class C1 cache
    class X1 external
```

### Step 4: Node Naming Rules

| Type | Format | Example |
|------|--------|---------|
| Service | Title Case | `Payment Service` |
| Handler | Title Case | `Order Handler` |
| Kafka Topic | lowercase.dot | `order.created` |
| Database | Title Case | `PostgreSQL` |
| Cache | Title Case | `Redis Cache` |
| External | Title Case | `Stripe API` |
| Method | PascalCase | `ProcessPayment` |

**NO ABBREVIATIONS.** Use full names always.

### Step 5: Arrow Rules

| Connection Type | Arrow | Label Format |
|-----------------|-------|--------------|
| gRPC call | `==>` | `\|gRPC: MethodName\|` |
| HTTP call | `==>` | `\|HTTP\|` or `\|REST\|` |
| SQL query | `==>` | `\|SQL\|` |
| Cache read/write | `==>` | `\|cache\|` |
| Kafka produce | `-.->` | `\|publish\|` |
| Kafka consume | `-.->` | `\|consume\|` |
| Internal call | `-->` | (no label or `\|internal\|`) |

### Step 6: Subgraph Rules

1. **Always quote titles**: `subgraph id ["Display Name"]`
2. **Use lowercase IDs**: `subgraph kafka ["Message Bus"]`
3. **Group by type**: All Kafka topics together, all databases together
4. **Single service per subgraph** for clarity

### Step 7: Write Complete Markdown File

Use the template from `templates/diagram-template.md`:

```markdown
# Service Flow: {service_name}

Generated: {timestamp}
Source: {target_path}
Dependencies: {output_dir}/dependencies.yaml

## Diagram

\`\`\`mermaid
{mermaid_code}
\`\`\`

## Legend

| Symbol | Meaning |
|--------|---------|
| `==>` | Synchronous (gRPC/HTTP) - caller blocks |
| `-.->` | Asynchronous (Kafka) - fire and forget |
| `-->` | Internal call |
| Blue rectangle | Service |
| Green cylinder | Kafka topic |
| Yellow cylinder | Database |
| Gray double rectangle | External system |
| Purple rounded | Cache |

## Sync Dependencies

| From | To | Type | Method/Endpoint | Timeout | Retry |
|------|-----|------|-----------------|---------|-------|
{sync_table}

## Async Dependencies

| Topic | Direction | Consumer Group | DLQ |
|-------|-----------|----------------|-----|
{async_table}

## External Systems

| System | Type | Purpose |
|--------|------|---------|
{external_table}

## Source References

{source_references}
```

## Quality Checklist

Before completing Phase 2, verify:

- [ ] All services from dependencies.yaml are in the diagram
- [ ] All Kafka topics are in the diagram
- [ ] All databases are in the diagram
- [ ] All external systems are in the diagram
- [ ] Sync calls use `==>` arrows
- [ ] Async calls use `-.->` arrows
- [ ] All subgraph titles are quoted
- [ ] classDef styles are applied to all nodes
- [ ] No abbreviations in node names
- [ ] Legend is included
- [ ] Summary tables are included

## Output

After completing Phase 2, report:

```
Phase 2 Complete: Diagram Generation

Output: {output_dir}/diagram.md

Diagram contains:
- Services: {count}
- Kafka topics: {count}
- Databases: {count}
- External systems: {count}
- Sync connections: {count}
- Async connections: {count}

Proceeding to Phase 3: Refinement
```
