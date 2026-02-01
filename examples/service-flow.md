# Service Flow Examples

## Example 1: E-Commerce Order Flow (Mixed Sync/Async)

**Command**: `/flow services/ --type=service`

This example shows a typical order processing flow with both gRPC (sync) and Kafka (async) communication.

```mermaid
flowchart LR
    subgraph Client [Client]
        C1[Mobile App]
    end

    subgraph APIGateway [API Gateway]
        GW[Gateway]
    end

    subgraph OrderService [Order Service]
        O1[OrderHandler]
        O2[OrderProcessor]
        O3[KafkaPublisher]
    end

    subgraph PaymentService [Payment Service]
        P1[gRPC Server]
        P2[PaymentProcessor]
    end

    subgraph InventoryService [Inventory Service]
        I1([inventory-consumer])
        I2[StockManager]
    end

    subgraph NotificationService [Notification Service]
        N1([notification-consumer])
        N2[EmailSender]
    end

    subgraph Kafka [Kafka]
        T1[(order.created)]
        T2[(order.completed)]
        DLQ[(order.created.dlq)]
    end

    subgraph External [External Systems]
        DB[[PostgreSQL]]
        STRIPE[[Stripe API]]
        SES[[AWS SES]]
    end

    %% SYNC paths (gRPC/HTTP) - thick arrows ==>
    C1 ==>|HTTPS| GW
    GW ==>|HTTP| O1
    O2 ==>|gRPC| P1
    P2 ==>|HTTPS| STRIPE
    O2 ==>|SQL| DB
    P2 ==>|SQL| DB
    I2 ==>|SQL| DB
    N2 ==>|HTTPS| SES

    %% ASYNC paths (Kafka) - dotted arrows -.->
    O3 -.->|publish| T1
    T1 -.->|consume| I1
    T1 -.->|consume| N1
    I1 -.->|on failure| DLQ
    I2 -.->|publish| T2

    %% Internal calls - thin arrows -->
    O1 --> O2
    O2 --> O3
    P1 --> P2
    I1 --> I2
    N1 --> N2

    %% Styling
    style GW fill:#69db7c
    style T1 fill:#4dabf7
    style T2 fill:#4dabf7
    style DLQ fill:#ff6b6b
    style DB fill:#868e96
    style STRIPE fill:#868e96
    style SES fill:#868e96
```

**Legend:**

| Arrow | Meaning | Debug Approach |
|-------|---------|----------------|
| `==>` | **Sync** (gRPC/HTTP) - caller blocks | Check latency, timeouts, error codes |
| `-.->` | **Async** (Kafka) - fire & forget | Check consumer lag, DLQ, offsets |
| `-->` | Internal call | Check logs, traces |

**Communication Summary:**

### Sync Calls (gRPC/HTTP)
| From | To | Protocol | Timeout | Notes |
|------|-----|----------|---------|-------|
| Client | Gateway | HTTPS | 30s | Client-facing |
| Gateway | OrderService | HTTP | 10s | Internal |
| OrderService | PaymentService | gRPC | 5s | Critical path |
| PaymentService | Stripe | HTTPS | 30s | External API |
| All Services | PostgreSQL | SQL | 5s | Database |

### Async Messages (Kafka)
| Topic | Producers | Consumers | DLQ |
|-------|-----------|-----------|-----|
| `order.created` | OrderService | InventoryService, NotificationService | Yes |
| `order.completed` | InventoryService | - | No |

**Analysis:**
- Order creation is sync through payment (user waits for payment confirmation)
- Inventory and notification are async (user doesn't wait)
- DLQ configured for inventory failures
- `order.completed` has no consumers - verify if this is intentional

---

## Example 2: Incident Investigation View

**Command**: `/flow services/ --type=service`

**Scenario**: Orders are being placed but customers aren't receiving confirmation emails.

```mermaid
flowchart LR
    subgraph OrderService [Order Service]
        O1[OrderHandler]
        O2[Publisher]
    end

    subgraph Kafka [Kafka]
        T1[(order.created<br/>lag: 50,000)]
        DLQ[(order.dlq<br/>depth: 12,847)]
    end

    subgraph NotificationService [Notification Service]
        N1([email-consumer<br/>1 instance])
        N2[EmailSender]
    end

    subgraph External [External]
        SES[[AWS SES<br/>rate limited]]
    end

    O1 --> O2
    O2 -.->|publish| T1
    T1 -.->|consume| N1
    N1 --> N2
    N2 ==>|HTTPS| SES
    N1 -.->|failures| DLQ

    style T1 fill:#ff6b6b
    style DLQ fill:#ff6b6b
    style N1 fill:#ffd43b
    style SES fill:#ff6b6b
```

**Legend:**

| Arrow | Meaning | Debug Approach |
|-------|---------|----------------|
| `==>` | **Sync** (gRPC/HTTP) - caller blocks | Check latency, timeouts, error codes |
| `-.->` | **Async** (Kafka) - fire & forget | Check consumer lag, DLQ, offsets |
| `-->` | Internal call | Check logs, traces |

**Issues Identified:**

### Critical ❌
1. **Consumer lag: 50,000 messages** on `order.created`
   - Messages are backing up faster than they're processed
   - Action: Scale consumers or investigate processing bottleneck

2. **DLQ depth: 12,847 messages**
   - Many messages are failing processing
   - Action: Investigate DLQ messages for common error pattern

3. **AWS SES rate limited**
   - External dependency is throttling requests
   - Action: Implement backoff, request limit increase

### Warnings ⚠️
1. **Single consumer instance** for critical topic
   - No redundancy, single point of failure
   - Action: Scale to at least 3 instances

**Incident Investigation Steps:**

1. ✅ Async path issue confirmed (not sync)
2. Check consumer lag: `kafka-consumer-groups --describe --group email-consumer`
3. Check DLQ: Sample messages for error patterns
4. Check SES: CloudWatch metrics for throttling
5. Scale consumers: Increase replicas from 1 to 3
6. Implement exponential backoff for SES calls

---

## Example 3: gRPC-Heavy Microservices

**Command**: `/flow services/ --type=service`

This example shows a system where most communication is sync gRPC.

```mermaid
flowchart TD
    subgraph Gateway [API Gateway]
        GW[GraphQL Gateway]
    end

    subgraph UserService [User Service - gRPC]
        U1[GetUser]
        U2[UpdateUser]
    end

    subgraph ProductService [Product Service - gRPC]
        PR1[GetProduct]
        PR2[ListProducts]
    end

    subgraph CartService [Cart Service - gRPC]
        CA1[GetCart]
        CA2[AddToCart]
    end

    subgraph PricingService [Pricing Service - gRPC]
        PI1[CalculatePrice]
    end

    subgraph InventoryService [Inventory Service - gRPC]
        IN1[CheckStock]
        IN2[ReserveStock]
    end

    subgraph Kafka [Async Events]
        T1[(cart.updated)]
        T2[(inventory.low)]
    end

    subgraph Analytics [Analytics Service]
        A1([analytics-consumer])
    end

    %% ALL gRPC calls are SYNC - use ==>
    GW ==>|gRPC| U1
    GW ==>|gRPC| PR1
    GW ==>|gRPC| CA1
    CA2 ==>|gRPC| PR1
    CA2 ==>|gRPC| PI1
    CA2 ==>|gRPC| IN1
    PI1 ==>|gRPC| PR1

    %% Async events - use -.->
    CA2 -.->|publish| T1
    IN2 -.->|publish| T2
    T1 -.->|consume| A1
    T2 -.->|consume| A1

    %% Styling
    style GW fill:#69db7c
    style T1 fill:#4dabf7
    style T2 fill:#4dabf7
    linkStyle 0,1,2,3,4,5,6 stroke:#be4bdb,stroke-width:2px
```

**Legend:**

| Arrow | Meaning | Debug Approach |
|-------|---------|----------------|
| `==>` | **Sync** (gRPC) - caller blocks | Check latency, timeouts, error codes, circuit breaker |
| `-.->` | **Async** (Kafka) - fire & forget | Check consumer lag, DLQ |

**Sync Call Chain Analysis:**

```
Gateway -> CartService -> ProductService -> (response)
                       -> PricingService -> ProductService -> (response)
                       -> InventoryService -> (response)
```

**Critical Path**: AddToCart requires 4 sync calls to complete
- Total latency = sum of all service latencies
- Single failure = entire request fails

**Recommendations:**
1. Add circuit breakers between services
2. Consider caching ProductService responses
3. Add timeout budgets (if Gateway timeout is 5s, downstream must be <5s total)

---

## Example 4: Saga Pattern with Compensation

**Command**: `/flow services/checkout/ --type=service`

```mermaid
flowchart LR
    subgraph CheckoutService [Checkout Orchestrator]
        C1[SagaCoordinator]
    end

    subgraph PaymentService [Payment Service]
        P1[ChargeCard]
        P2[RefundCard]
    end

    subgraph InventoryService [Inventory Service]
        I1[ReserveStock]
        I2[ReleaseStock]
    end

    subgraph ShippingService [Shipping Service]
        S1[CreateShipment]
        S2[CancelShipment]
    end

    subgraph Kafka [Saga Events]
        T1[(saga.payment.request)]
        T2[(saga.payment.completed)]
        T3[(saga.payment.failed)]
        T4[(saga.inventory.request)]
        T5[(saga.inventory.completed)]
        T6[(saga.inventory.failed)]
        T7[(saga.shipping.request)]
        T8[(saga.compensate)]
    end

    %% Saga orchestration is ASYNC
    C1 -.->|1. request| T1
    T1 -.-> P1
    P1 -.->|success| T2
    P1 -.->|failure| T3
    T2 -.-> C1

    C1 -.->|2. request| T4
    T4 -.-> I1
    I1 -.->|success| T5
    I1 -.->|failure| T6
    T5 -.-> C1

    C1 -.->|3. request| T7
    T7 -.-> S1

    %% Compensation flow
    T3 -.-> C1
    T6 -.-> C1
    C1 -.->|compensate| T8
    T8 -.-> P2
    T8 -.-> I2
    T8 -.-> S2

    %% Styling
    style C1 fill:#69db7c
    style T1 fill:#4dabf7
    style T2 fill:#69db7c
    style T3 fill:#ff6b6b
    style T4 fill:#4dabf7
    style T5 fill:#69db7c
    style T6 fill:#ff6b6b
    style T8 fill:#fd7e14
```

**Legend:**

| Arrow | Meaning | Debug Approach |
|-------|---------|----------------|
| `-.->` | **Async** (Kafka) - all saga steps are async | Check consumer lag, message ordering |
| Orange | Compensation/rollback path | Check if compensations completed |

**Saga Flow:**
1. All steps are **async** (Kafka-based saga)
2. Each step publishes success or failure
3. Coordinator decides next step or compensation
4. Compensation reverses completed steps on failure

**Debugging Saga Issues:**
1. Check saga state in database
2. Verify all step events were published
3. Check for stuck sagas (no progress)
4. Verify compensation messages were processed

---

## Example 5: Verbose Mode (Onboarding)

**Command**: `/flow services/ --type=service --verbose`

```mermaid
flowchart LR
    subgraph OrderService [Order Service]
        O1[HTTP Handler]
        O2[Business Logic]
        O3[Kafka Publisher]
    end

    subgraph PaymentService [Payment Service]
        P1[gRPC Server]
        P2[Payment Logic]
    end

    subgraph Kafka [Message Bus]
        T1[(order.events)]
    end

    subgraph External [External]
        DB[[Database]]
    end

    %% Sync
    O1 ==>|gRPC| P1
    O2 ==>|SQL| DB
    P2 ==>|SQL| DB

    %% Async
    O3 -.->|Kafka| T1

    %% Internal
    O1 --> O2 --> O3
    P1 --> P2

    style O1 fill:#69db7c
    style T1 fill:#4dabf7
    style DB fill:#868e96
```

---

## How to Read This Diagram

### Understanding Arrow Types (CRITICAL)

The most important thing to understand is the **difference between arrow types**:

#### Thick Solid Arrows (`==>`) = SYNC
```
ServiceA ==>|gRPC| ServiceB
```
- **What happens**: ServiceA calls ServiceB and **WAITS** for a response
- **If it fails**: ServiceA gets an error immediately
- **Timeout**: If ServiceB is slow, ServiceA might timeout
- **Example**: `response, err := paymentClient.Charge(ctx, request)`

#### Dotted Arrows (`-.->`) = ASYNC
```
ServiceA -.->|Kafka| Topic
```
- **What happens**: ServiceA publishes a message and **CONTINUES immediately**
- **If it fails**: ServiceA doesn't know (unless publish itself fails)
- **Processing**: Happens later, by some consumer
- **Example**: `producer.Send("orders", orderEvent)`

#### Thin Solid Arrows (`-->`) = INTERNAL
```
Handler --> Processor
```
- **What happens**: Function call within the same service
- **Same process**: No network involved

### Why This Matters for Debugging

**Scenario**: User complains order confirmation email never arrived

**If the path is SYNC (`==>`)**:
- Error happened immediately
- Check error logs at time of request
- Look for timeout or error response
- User probably saw an error

**If the path is ASYNC (`-.->`)**:
- Error happened later (or message is stuck)
- Check Kafka consumer lag
- Check Dead Letter Queue
- User saw "success" but processing failed later

### Node Shapes

| Shape | Syntax | Meaning |
|-------|--------|---------|
| Rectangle | `[Name]` | Service, handler, or component |
| Cylinder | `[(name)]` | Kafka topic or message queue |
| Stadium | `([name])` | Consumer group |
| Double Rectangle | `[[Name]]` | External system (DB, 3rd party API) |

### Colors

| Color | Meaning |
|-------|---------|
| 🟢 Green | Entry point - where requests come in |
| 🔵 Blue | Kafka topic - async message storage |
| 🔴 Red | Error path - DLQ, failures |
| 🟡 Yellow | Warning - potential issue |
| 🟣 Purple | gRPC calls (optional highlight) |
| ⚫ Gray | External systems |
| 🟠 Orange | Retry/compensation paths |

---

## Glossary

### Sync Communication
| Term | Definition |
|------|------------|
| **gRPC** | High-performance RPC framework using HTTP/2 and Protocol Buffers |
| **Timeout** | Maximum time to wait for sync response before failing |
| **Circuit Breaker** | Pattern to stop calling failing services temporarily |
| **Deadline** | gRPC concept - absolute time by which call must complete |

### Async Communication
| Term | Definition |
|------|------------|
| **Topic** | Named channel in Kafka where messages are stored |
| **Partition** | Subdivision of topic for parallelism |
| **Consumer Group** | Set of consumers sharing work on a topic |
| **Offset** | Position in topic - where consumer has read up to |
| **Lag** | How far behind a consumer is (unprocessed messages) |
| **DLQ** | Dead Letter Queue - where failed messages go |

### General
| Term | Definition |
|------|------------|
| **Saga** | Pattern for distributed transactions using async events |
| **Idempotency** | Processing same message twice has same effect as once |
| **At-least-once** | Messages delivered at least once (may duplicate) |
| **Exactly-once** | Messages delivered exactly once (hard to achieve) |

---

## Incident Investigation Guide

### Sync Path Issues (gRPC/HTTP)

**Symptoms**: Immediate errors, timeouts, high latency

**What to check**:
1. **Error rates**: Dashboard showing 5xx or gRPC errors
2. **Latency**: P99 latency spikes
3. **Circuit breaker**: Is it open? Check circuit breaker dashboard
4. **Downstream health**: Is the called service healthy?
5. **Logs**: Error messages at time of failure

**Common causes**:
- Downstream service down
- Network issues
- Resource exhaustion (CPU, memory, connections)
- Timeout too aggressive

### Async Path Issues (Kafka)

**Symptoms**: Delayed processing, missing data, growing lag

**What to check**:
1. **Consumer lag**: `kafka-consumer-groups --describe --group <group>`
2. **DLQ depth**: How many messages in dead letter queue?
3. **Consumer health**: Are consumer pods running?
4. **Partition assignment**: Are all partitions assigned?
5. **Message inspection**: Sample DLQ messages for patterns

**Common causes**:
- Consumer crashed or stuck
- Processing too slow (scale consumers)
- Poison message blocking partition
- Schema mismatch
- External dependency in consumer failing

---

## Architecture Context

This diagram shows a **microservices architecture** with:

### Communication Patterns
- **Sync (gRPC)**: Used for operations where caller needs immediate response
- **Async (Kafka)**: Used for decoupling, event-driven flows, eventual consistency

### When to Use Sync
- User is waiting for response
- Need immediate confirmation
- Transaction must be atomic

### When to Use Async
- User doesn't need to wait
- Decoupling between services
- High throughput requirements
- Event sourcing / audit trails

---

## Next Steps for Learning

1. **Trace a request**: Follow a single order from API to completion
2. **Find gRPC definitions**: Look in `proto/` directory for `.proto` files
3. **Find Kafka topics**: Search for `topic` in config files
4. **Check consumer groups**: `kafka-consumer-groups --list`
5. **Review error handling**: Search for timeout and retry configuration
6. **Understand circuit breakers**: Look for hystrix, resilience4j, or similar

**Key files to explore**:
- `proto/*.proto` - gRPC service definitions
- `config/kafka.yaml` - Kafka topic configuration
- `internal/client/*.go` - gRPC client wrappers
- `internal/consumer/*.go` - Kafka consumers
