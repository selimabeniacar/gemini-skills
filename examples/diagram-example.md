# Service Flow: Ledger Service

> Generated: 2024-01-15T14:35:00Z
> Source: services/ledger/

---

## Diagram

```mermaid
flowchart TD
    %% ========================================
    %% Style Definitions (muted, professional colors)
    %% ========================================
    classDef service fill:#a5d8ff,stroke:#339af0,color:#1864ab
    classDef entry fill:#b2f2bb,stroke:#51cf66,color:#2b8a3e
    classDef kafka fill:#96f2d7,stroke:#38d9a9,color:#087f5b
    classDef database fill:#ffec99,stroke:#fcc419,color:#e67700
    classDef cache fill:#d0bfff,stroke:#9775fa,color:#6741d9
    classDef external fill:#dee2e6,stroke:#adb5bd,color:#495057

    %% ========================================
    %% Entry Points (Callers)
    %% ========================================
    subgraph entry ["Entry Points"]
        direction LR
        E1[Order Service]
        E2[Reporting Service]
        E3[Admin Dashboard]
    end

    %% ========================================
    %% Ledger Service
    %% ========================================
    subgraph target ["Ledger Service"]
        T1[gRPC Server]
        T2[HTTP Server]
        T3[Transaction Handler]
        T4[Balance Handler]
    end

    %% ========================================
    %% Dependent Services
    %% ========================================
    subgraph deps ["Dependencies"]
        direction LR
        D1[Payment Service]
        D2[Account Service]
        D3[Audit Service]
    end

    %% ========================================
    %% Message Bus
    %% ========================================
    subgraph kafka ["Message Bus"]
        direction LR
        K1[(order.completed)]
        K2[(payment.refunded)]
        K3[(ledger.transaction.created)]
        K4[(ledger.balance.updated)]
    end

    %% ========================================
    %% Data Stores
    %% ========================================
    subgraph data ["Data Stores"]
        direction LR
        DB1[(PostgreSQL)]
        C1(Redis Cache)
    end

    %% ========================================
    %% External Systems
    %% ========================================
    subgraph ext ["External Systems"]
        X1[[Stripe API]]
    end

    %% ========================================
    %% Sync Connections (gRPC, HTTP, SQL)
    %% ========================================
    E1 ==>|gRPC: CreateTransaction| T1
    E2 ==>|gRPC: ListTransactions| T1
    E3 ==>|HTTP| T2
    T1 --> T3
    T1 --> T4
    T2 --> T3
    T3 ==>|gRPC: ProcessPayment| D1
    T3 ==>|gRPC: GetAccount| D2
    T3 ==>|gRPC: LogAudit| D3
    T3 ==>|SQL| DB1
    T4 ==>|SQL| DB1
    T4 ==>|cache| C1
    T2 ==>|HTTPS| X1

    %% ========================================
    %% Async Connections (Kafka)
    %% ========================================
    K1 -.->|consume| T3
    K2 -.->|consume| T3
    T3 -.->|publish| K3
    T4 -.->|publish| K4

    %% ========================================
    %% Apply Styles
    %% ========================================
    class E1,E2,E3 entry
    class T1,T2,T3,T4 service
    class D1,D2,D3 service
    class K1,K2,K3,K4 kafka
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
| Ledger Service | Payment Service | gRPC | ProcessPayment, RefundPayment | 5s | Yes | internal/client/payment_client.go:45 |
| Ledger Service | Account Service | gRPC | GetAccount, ValidateAccount | 3s | Yes | internal/client/account_client.go:32 |
| Ledger Service | Audit Service | gRPC | LogAudit | - | - | internal/audit/client.go:18 |
| Ledger Service | PostgreSQL | SQL | read/write | - | - | internal/repository/ledger_repo.go:12 |
| Ledger Service | Redis | cache | read/write | - | - | internal/cache/balance_cache.go:15 |
| Ledger Service | Stripe API | HTTPS | webhook verify | - | - | internal/webhook/stripe_handler.go:56 |

---

## Async Dependencies

| Topic | Direction | Consumer Group | DLQ | Source |
|-------|-----------|----------------|-----|--------|
| order.completed | consume | ledger-order-consumer | Yes | internal/consumer/order_consumer.go:27 |
| payment.refunded | consume | ledger-refund-consumer | Yes | internal/consumer/refund_consumer.go:15 |
| ledger.transaction.created | produce | - | - | internal/producer/transaction_producer.go:34 |
| ledger.balance.updated | produce | - | - | internal/producer/balance_producer.go:22 |

---

## External Systems

| System | Type | Purpose | Source |
|--------|------|---------|--------|
| Stripe API | HTTPS | Payment webhook verification | internal/webhook/stripe_handler.go:56 |
| Audit Service | gRPC | Compliance logging | internal/audit/client.go:18 |

---

## Data Stores

| Store | Type | Operations | Source |
|-------|------|------------|--------|
| PostgreSQL | Database | read/write | internal/repository/ledger_repo.go:12 |
| Redis | Cache | Balance cache | internal/cache/balance_cache.go:15 |

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
- `services/ledger/docs/RUNBOOK.md` - Service runbook

---

## Render Commands

Generate PNG - high resolution (recommended for documentation):
```bash
npx -p @mermaid-js/mermaid-cli mmdc -i diagram.md -o diagram.png -b white -w 3840 -s 2
```

Generate SVG (for web, scalable):
```bash
npx -p @mermaid-js/mermaid-cli mmdc -i diagram.md -o diagram.svg -b white
```

Generate PDF:
```bash
npx -p @mermaid-js/mermaid-cli mmdc -i diagram.md -o diagram.pdf -b white
```
