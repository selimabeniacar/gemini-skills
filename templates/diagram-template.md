# Service Flow: {SERVICE_NAME}

> Generated: {TIMESTAMP}
> Source: {TARGET_PATH}

---

## Diagram

```mermaid
flowchart LR
    %% ========================================
    %% Style Definitions
    %% ========================================
    classDef service fill:#228be6,stroke:#1971c2,color:#fff
    classDef entry fill:#40c057,stroke:#2f9e44,color:#fff
    classDef kafka fill:#12b886,stroke:#099268,color:#fff
    classDef database fill:#fab005,stroke:#f59f00,color:#000
    classDef cache fill:#be4bdb,stroke:#9c36b5,color:#fff
    classDef external fill:#868e96,stroke:#495057,color:#fff

    %% ========================================
    %% Entry Points
    %% ========================================
    subgraph entry ["Entry Points"]
        {ENTRY_NODES}
    end

    %% ========================================
    %% {SERVICE_NAME}
    %% ========================================
    subgraph target ["{SERVICE_NAME}"]
        {TARGET_NODES}
    end

    %% ========================================
    %% Dependent Services
    %% ========================================
    subgraph deps ["Dependencies"]
        {DEPENDENCY_NODES}
    end

    %% ========================================
    %% Message Bus
    %% ========================================
    subgraph kafka ["Message Bus"]
        {KAFKA_NODES}
    end

    %% ========================================
    %% Data Stores
    %% ========================================
    subgraph data ["Data Stores"]
        {DATA_NODES}
    end

    %% ========================================
    %% External Systems
    %% ========================================
    subgraph ext ["External Systems"]
        {EXTERNAL_NODES}
    end

    %% ========================================
    %% Sync Connections (gRPC, HTTP, SQL)
    %% ========================================
    {SYNC_CONNECTIONS}

    %% ========================================
    %% Async Connections (Kafka)
    %% ========================================
    {ASYNC_CONNECTIONS}

    %% ========================================
    %% Apply Styles
    %% ========================================
    {CLASS_APPLICATIONS}
```

---

## Legend

| Symbol | Meaning | Debug Approach |
|--------|---------|----------------|
| `==>` | **Synchronous** (gRPC/HTTP) - caller blocks | Check latency, timeouts, error codes |
| `-.->` | **Asynchronous** (Kafka) - fire and forget | Check consumer lag, DLQ, offsets |
| `-->` | Internal call | Check logs, traces |

### Node Shapes

| Shape | Meaning |
|-------|---------|
| `[Rectangle]` | Service, Handler |
| `[(Cylinder)]` | Database, Kafka Topic |
| `([Stadium])` | Consumer Group |
| `[[Double Rect]]` | External System |
| `(Rounded)` | Cache |

### Colors

| Color | Meaning |
|-------|---------|
| Blue | Services |
| Green | Entry Points |
| Teal | Kafka Topics |
| Yellow | Databases |
| Purple | Caches |
| Gray | External Systems |

---

## Sync Dependencies

| From | To | Type | Method/Endpoint | Timeout | Retry | Source |
|------|-----|------|-----------------|---------|-------|--------|
{SYNC_TABLE}

---

## Async Dependencies

| Topic | Direction | Consumer Group | DLQ | Source |
|-------|-----------|----------------|-----|--------|
{ASYNC_TABLE}

---

## External Systems

| System | Type | Purpose | Source |
|--------|------|---------|--------|
{EXTERNAL_TABLE}

---

## Data Stores

| Store | Type | Operations | Source |
|-------|------|------------|--------|
{DATA_TABLE}

---

## Source References

All dependencies traced from:

{SOURCE_REFERENCES}

---

## Render Commands

Generate PNG (recommended for documentation):
```bash
mmdc -i diagram.md -o diagram.png -b white -w 1920 -H 1080
```

Generate SVG (for web):
```bash
mmdc -i diagram.md -o diagram.svg -b white
```

Generate PDF:
```bash
mmdc -i diagram.md -o diagram.pdf -b white
```
