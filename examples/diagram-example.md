# Service Flow: Ledger Service

> Generated: 2024-01-15T14:35:00Z
> Source: services/ledger/

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
        E1([LedgerService])
        E2([/api/v1/ledger])
    end

    %% ========================================
    %% Ledger Service — with internal processing steps
    %% ========================================
    subgraph target ["Ledger Service"]
        direction TB
        S1_step1[Validate Request] --> S1_step2[Load Account]
        S1_step2 --> S1_step3[Process Transaction]
        S1_step3 --> S1_step4[Commit]
    end

    %% ========================================
    %% Dependencies — gRPC services Ledger calls
    %% ========================================
    subgraph deps ["Dependencies"]
        D1[Payment Service]
        D2[Account Service]
    end

    %% ========================================
    %% Data Stores — Ledger's own database + cache
    %% Logical names, NOT technology names
    %% ========================================
    subgraph data ["Data Stores"]
        DB1[(Ledger DB)]
        C1(Ledger Cache)
    end

    %% ========================================
    %% Consumed Topics — Kafka topics Ledger consumes
    %% ========================================
    subgraph kafka-in ["Consumed Topics"]
        KI1[(order.completed)]
        KI2[(payment.refunded)]
    end

    %% ========================================
    %% Produced Topics — Kafka topics Ledger produces
    %% ========================================
    subgraph kafka-out ["Produced Topics"]
        KO1[(ledger.transaction.created)]
        KO2[(ledger.balance.updated)]
    end

    %% ========================================
    %% External Systems
    %% ========================================
    subgraph ext ["External"]
        EX1[Stripe API]
        EX2[Audit Service]
    end

    %% ========================================
    %% Arrows — to SUBGROUPS, not individual nodes
    %% 6 arrows total for 15 nodes
    %% ========================================
    entry ==> target
    kafka-in -.-> target
    S1_step2 ==> data
    S1_step3 ==> deps
    S1_step4 -.-> kafka-out
    S1_step4 ==> ext

    %% ========================================
    %% Apply Styles
    %% ========================================
    class E1,E2 entry
    class S1_step1,S1_step2,S1_step3,S1_step4 step
    class D1,D2 service
    class DB1 database
    class C1 cache
    class KI1,KI2,KO1,KO2 kafka
    class EX1,EX2 external
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
| Ledger Service | Payment Service | gRPC | internal/client/payment_client.go:45 |
| Ledger Service | Account Service | gRPC | internal/client/account_client.go:32 |

---

## Async Dependencies

| Topic | Direction | Source |
|-------|-----------|--------|
| order.completed | consume | internal/consumer/order_consumer.go:27 |
| payment.refunded | consume | internal/consumer/refund_consumer.go:15 |
| ledger.transaction.created | produce | internal/producer/transaction_producer.go:34 |
| ledger.balance.updated | produce | internal/producer/balance_producer.go:22 |

---

## Data Stores

| Store | Type | Source |
|-------|------|--------|
| Ledger DB | postgresql | internal/repository/ledger_repo.go:12 |
| Ledger Cache | redis | internal/cache/balance_cache.go:15 |

---

## External Systems

| System | Type | Source |
|--------|------|--------|
| Stripe API | https | internal/webhook/stripe_handler.go:56 |
| Audit Service | grpc | internal/audit/client.go:18 |

---

## Source References

All dependencies traced from:

- `internal/client/payment_client.go:45` - Payment Service gRPC client
- `internal/client/account_client.go:32` - Account Service gRPC client
- `internal/consumer/order_consumer.go:27` - Order completed consumer
- `internal/consumer/refund_consumer.go:15` - Payment refunded consumer
- `internal/producer/transaction_producer.go:34` - Transaction created producer
- `internal/producer/balance_producer.go:22` - Balance updated producer
- `internal/repository/ledger_repo.go:12` - Database repository
- `internal/cache/balance_cache.go:15` - Redis cache client
- `internal/webhook/stripe_handler.go:56` - Stripe webhook handler
- `internal/audit/client.go:18` - Audit service client

---

## Render Commands

```bash
# PNG - high resolution (for documentation)
npx -p @mermaid-js/mermaid-cli mmdc -i flow-diagram.md -o flow-diagram.png -b white -w 3840 -s 2

# SVG (for web, scalable)
npx -p @mermaid-js/mermaid-cli mmdc -i flow-diagram.md -o flow-diagram.svg -b white
```
