# Service Flow: __SERVICE_NAME__

> Generated: __TIMESTAMP__
> Source: __TARGET_PATH__

---

## Diagram

```mermaid
flowchart TD
    %% ========================================
    %% Style Definitions (copy exactly)
    %% ========================================
    classDef service fill:#228be6,stroke:#1971c2,color:#fff
    classDef entry fill:#40c057,stroke:#2f9e44,color:#fff
    classDef kafka fill:#12b886,stroke:#099268,color:#fff
    classDef database fill:#fab005,stroke:#f59f00,color:#000
    classDef cache fill:#be4bdb,stroke:#9c36b5,color:#fff
    classDef external fill:#868e96,stroke:#495057,color:#fff

    %% ========================================
    %% Entry Points (TOP) - services that call this one
    %% Replace with actual caller services
    %% ========================================
    subgraph entry ["Entry Points"]
        E1[Caller Service 1]
        E2[Caller Service 2]
    end

    %% ========================================
    %% Target Service - the service being documented
    %% Replace with actual handlers/components
    %% ========================================
    subgraph target ["__SERVICE_NAME__"]
        T1[Handler 1]
        T2[Handler 2]
    end

    %% ========================================
    %% Dependencies - services this one calls
    %% Replace with actual dependencies
    %% ========================================
    subgraph deps ["Dependencies"]
        D1[Dependency Service 1]
        D2[Dependency Service 2]
    end

    %% ========================================
    %% Message Bus - Kafka topics
    %% Use cylinder shape: [(topic.name)]
    %% ========================================
    subgraph kafka ["Message Bus"]
        K1[(topic.consumed)]
        K2[(topic.produced)]
    end

    %% ========================================
    %% Data Stores - databases and caches
    %% Database: [(Name)], Cache: (Name)
    %% ========================================
    subgraph data ["Data Stores"]
        DB1[(PostgreSQL)]
        C1(Redis Cache)
    end

    %% ========================================
    %% External Systems (BOTTOM)
    %% Use double brackets: [[Name]]
    %% ========================================
    subgraph ext ["External Systems"]
        X1[[External API]]
    end

    %% ========================================
    %% Sync Connections - use ==> for gRPC/HTTP/SQL
    %% ========================================
    E1 ==>|gRPC| T1
    E2 ==>|HTTP| T1
    T1 --> T2
    T2 ==>|gRPC: Method| D1
    T2 ==>|SQL| DB1
    T2 ==>|cache| C1
    D1 ==>|HTTPS| X1

    %% ========================================
    %% Async Connections - use -.-> for Kafka
    %% ========================================
    K1 -.->|consume| T1
    T2 -.->|publish| K2

    %% ========================================
    %% Apply Styles (list all node IDs)
    %% ========================================
    class E1,E2 entry
    class T1,T2 service
    class D1,D2 service
    class K1,K2 kafka
    class DB1 database
    class C1 cache
    class X1 external
```

---

## Legend

| Symbol | Meaning | Debug Approach |
|--------|---------|----------------|
| `==>` | **Synchronous** (gRPC/HTTP) - caller blocks | Check latency, timeouts, error codes |
| `-.->` | **Asynchronous** (Kafka) - fire and forget | Check consumer lag, DLQ, offsets |
| `-->` | Internal call | Check logs, traces |

### Node Shapes

| Shape | Meaning | Syntax |
|-------|---------|--------|
| Rectangle | Service, Handler | `[Name]` |
| Cylinder | Database, Kafka Topic | `[(Name)]` |
| Stadium | Consumer Group | `([Name])` |
| Double Rectangle | External System | `[[Name]]` |
| Rounded | Cache | `(Name)` |

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
| _fill from .flow-deps.yaml_ | | | | | | |

---

## Async Dependencies

| Topic | Direction | Consumer Group | DLQ | Source |
|-------|-----------|----------------|-----|--------|
| _fill from .flow-deps.yaml_ | | | | |

---

## Source References

All dependencies traced from:

- _list source files from .flow-deps.yaml_

---

## Render Commands

```bash
# PNG (for documentation)
mmdc -i flow-diagram.md -o flow-diagram.png -b white -w 1920 -H 1080

# SVG (for web)
mmdc -i flow-diagram.md -o flow-diagram.svg -b white
```
