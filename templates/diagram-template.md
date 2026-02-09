# Service Flow: __SERVICE_NAME__

> Generated: __TIMESTAMP__
> Source: __TARGET_PATH__

---

## Diagram

```mermaid
flowchart TD
    %% ========================================
    %% Style Definitions
    %% ========================================
    classDef service fill:#a5d8ff,stroke:#339af0,color:#1864ab
    classDef entry fill:#b2f2bb,stroke:#51cf66,color:#2b8a3e
    classDef kafka fill:#96f2d7,stroke:#38d9a9,color:#087f5b
    classDef database fill:#ffec99,stroke:#fcc419,color:#e67700
    classDef cache fill:#d0bfff,stroke:#9775fa,color:#6741d9
    classDef external fill:#dee2e6,stroke:#adb5bd,color:#495057
    classDef step fill:#e8f4f8,stroke:#4a9ebb,color:#2c5f7c

    %% ========================================
    %% Entry Points — the service's own endpoints
    %% ========================================
    subgraph entry ["Entry Points"]
        E1([Endpoint 1])
        E2([Endpoint 2])
    end

    %% ========================================
    %% Target Service — with internal steps if present
    %% ========================================
    subgraph target ["__SERVICE_NAME__"]
        direction TB
        S1_step1[Step 1] --> S1_step2[Step 2]
        S1_step2 --> S1_step3[Step 3]
    end

    %% ========================================
    %% Dependencies — gRPC/HTTP services this one calls
    %% Group by service context, max 4 per subgroup
    %% ========================================
    subgraph deps ["Dependencies"]
        D1[Dependency Service 1]
        D2[Dependency Service 2]
    end

    %% ========================================
    %% Data Stores — this service's database + cache
    %% Use logical names (e.g., "Order DB"), NOT technology names
    %% ========================================
    subgraph data ["Data Stores"]
        DB1[(__SERVICE_SHORT__ DB)]
        C1(__SERVICE_SHORT__ Cache)
    end

    %% ========================================
    %% Consumed Topics — Kafka topics this service consumes
    %% ========================================
    subgraph kafka-in ["Consumed Topics"]
        KI1[(topic.consumed.1)]
        KI2[(topic.consumed.2)]
    end

    %% ========================================
    %% Produced Topics — Kafka topics this service produces
    %% ========================================
    subgraph kafka-out ["Produced Topics"]
        KO1[(topic.produced.1)]
        KO2[(topic.produced.2)]
    end

    %% ========================================
    %% External Systems — third-party APIs
    %% ========================================
    subgraph ext ["External"]
        X1[External API]
    end

    %% ========================================
    %% Arrows — to SUBGROUPS, not individual nodes
    %% One arrow per subgroup, never per node inside
    %% ========================================
    entry ==> target
    kafka-in -.-> target
    S1_step2 ==> data
    S1_step2 ==> deps
    S1_step3 -.-> kafka-out
    S1_step3 ==> ext

    %% ========================================
    %% Apply Styles
    %% ========================================
    class E1,E2 entry
    class S1_step1,S1_step2,S1_step3 step
    class D1,D2 service
    class DB1 database
    class C1 cache
    class KI1,KI2,KO1,KO2 kafka
    class X1 external
```

---

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

---

## Sync Dependencies

| From | To | Type | Source |
|------|-----|------|--------|
| _fill from .flow-deps.yaml_ | | | |

---

## Async Dependencies

| Topic | Direction | Source |
|-------|-----------|--------|
| _fill from .flow-deps.yaml_ | | |

---

## Source References

All dependencies traced from:

- _list source files from .flow-deps.yaml_

---

## Render Commands

```bash
# PNG - high resolution (for documentation)
npx -p @mermaid-js/mermaid-cli mmdc -i flow-diagram.md -o flow-diagram.png -b white -w 3840 -s 2

# SVG (for web, scalable)
npx -p @mermaid-js/mermaid-cli mmdc -i flow-diagram.md -o flow-diagram.svg -b white
```
