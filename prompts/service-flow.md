# Service Flow Generation Prompt

## Overview

Generate a service flow diagram showing **both sync (gRPC/HTTP) and async (Kafka/queues)** communication between services. The distinction between sync and async is **critical** for incident investigation and onboarding.

## ⚠️ CRITICAL REQUIREMENTS

### 🛑 NO BROAD SEARCHES

**NEVER grep or search the entire codebase for keywords.**

If the user asks for a diagram of "EOD service" or similar:
1. Ask for the directory path, OR
2. Check runbooks (`runbooks/`, `docs/`) for the path, OR
3. Use `ls services/` to find the directory name

**Only read files within the target service directory.** Broad searches cause context overflow.

### No Single-Service Diagrams

**NEVER generate a diagram with only one service.**

Service flow mode MUST show:
- The target service AND all services it communicates with
- Minimum 2+ services (target + dependencies)
- Kafka topics as intermediaries between services
- External systems (databases, third-party APIs)

If you only found one service, you have NOT finished investigating. Keep looking for:
- gRPC/HTTP clients (outbound calls)
- Kafka producers/consumers (async communication)
- Database connections
- Config files mentioning other services

### Thoroughness Requirement

**Take your time. Be methodical. Do not skip dependencies.**

This diagram will be used for:
- **Incident investigation** - missing a dependency could mean missing the root cause
- **Onboarding** - new engineers need the complete picture

**You MUST check ALL sources:**
1. Source code (*.go, *.py, *.proto)
2. Configuration files (*.yaml, *.yml, *.json)
3. Documentation (README.md, RUNBOOK.md, docs/)
4. Environment variables and secrets references

**If unsure, include it with a note rather than skip it.**

## CRITICAL: Arrow Style Rules

**You MUST use the correct arrow style. This is non-negotiable.**

| Communication Type | Arrow Syntax | When to Use |
|-------------------|--------------|-------------|
| **Sync (gRPC/HTTP)** | `==>` | gRPC calls, HTTP requests, REST APIs, direct DB queries |
| **Async (Kafka/Queue)** | `-.->` | Kafka, RabbitMQ, SQS, Redis pub/sub, any message queue |
| **Internal (same service)** | `-->` | Function calls within the same service boundary |

**Why this matters:**
- Sync: Caller is BLOCKED waiting for response. Timeout = immediate failure.
- Async: Caller continues immediately. Failure is delayed, shows up as lag or DLQ.

## CRITICAL: Read ALL Files First

**DO NOT cherry-pick one file and generate a diagram. You MUST:**

1. **List the entire directory structure** first
2. **Identify ALL services/components** in the directory
3. **Report what you found** before generating anything:
   ```
   Directory scan results:
   - ledger/api/          → Main API handlers (INCLUDE)
   - ledger/consumer/     → Kafka consumers (INCLUDE)
   - ledger/producer/     → Kafka producers (INCLUDE)
   - ledger/grpc/         → gRPC server/client (INCLUDE)
   - ledger/jobs/         → Background jobs (DEPRIORITIZE - not main flow)
   - ledger/migrations/   → DB migrations (SKIP)
   - ledger/tests/        → Tests (SKIP)
   ```
4. **Read the main service files**, not jobs/utilities/tests
5. **Within the target directory**, look for communication patterns:
   - gRPC client calls → which services does this call?
   - Kafka producers → which topics does this publish to?
   - Kafka consumers → which topics does this consume from?
   - HTTP clients → which services does this call?

   ⚠️ **DO NOT search/grep across the entire codebase. Only read files in the target directory.**

**Files to PRIORITIZE:**
- `main.go`, `app.py` - entry points
- `*_handler.go`, `*_server.go` - API handlers
- `*_client.go` - gRPC/HTTP clients (shows outbound calls)
- `*_consumer.go`, `*_producer.go` - Kafka integration
- `*.proto` - gRPC service definitions

**Files to SKIP:**
- `*_test.go`, `test_*.py` - tests
- `*_mock.go`, `mock_*.py` - mocks
- `migrations/` - database migrations
- `fixtures/` - test data

## Steps

### 1. Identify Service Boundaries

Scan for:
- Directory structure: `services/payment/`, `cmd/order-service/`
- Go modules: `module github.com/company/payment-service`
- Python packages: separate `requirements.txt` or `pyproject.toml`
- Kubernetes/Docker: `deployment.yaml`, `Dockerfile`
- Proto files: `service PaymentService { ... }`

### 2. Identify Sync Communication (gRPC/HTTP)

#### Go gRPC Patterns
```go
// Client creation
conn, err := grpc.Dial(address, opts...)
client := pb.NewPaymentServiceClient(conn)

// Calling another service (THIS IS SYNC - use ==>)
response, err := client.ProcessPayment(ctx, request)
```

#### Go HTTP Patterns
```go
// HTTP client calls (THIS IS SYNC - use ==>)
resp, err := http.Get("http://other-service/api/...")
resp, err := client.Do(req)
```

#### Python gRPC Patterns
```python
# Client creation
channel = grpc.insecure_channel('payment-service:50051')
stub = payment_pb2_grpc.PaymentServiceStub(channel)

# Calling another service (THIS IS SYNC - use ==>)
response = stub.ProcessPayment(request)
```

#### Python HTTP Patterns
```python
# HTTP calls (THIS IS SYNC - use ==>)
response = requests.get("http://other-service/api/...")
response = await client.get("http://other-service/api/...")
```

### 3. Identify Async Communication (Kafka/Queues)

#### Go Kafka Patterns
```go
// Sarama producer (THIS IS ASYNC - use -.->)
producer.SendMessage(&sarama.ProducerMessage{Topic: "orders"})

// Confluent producer (THIS IS ASYNC - use -.->)
p.Produce(&kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &topic}})

// Consumer (receives async messages)
for msg := range claim.Messages() { ... }
```

#### Python Kafka Patterns
```python
# kafka-python (THIS IS ASYNC - use -.->)
producer.send('orders', value=data)

# Consumer
for msg in consumer:
    process(msg)

# Faust (THIS IS ASYNC - use -.->)
@app.agent(orders_topic)
async def process(orders):
    async for order in orders:
        ...
```

#### Other Async Patterns
```go
// RabbitMQ (ASYNC - use -.->)
ch.Publish(exchange, routingKey, mandatory, immediate, msg)

// Redis pub/sub (ASYNC - use -.->)
client.Publish(ctx, "channel", message)
```

### 4. Identify External Systems

Mark external dependencies distinctly:
- Databases: PostgreSQL, MySQL, MongoDB
- Caches: Redis, Memcached
- Third-party APIs: Stripe, Twilio, AWS services

### 5. Identify Neighboring Services

**This is critical for architectural diagrams.** Look for:

#### From gRPC clients:
```go
// This tells you the service calls "payment-service"
client := pb.NewPaymentServiceClient(conn)
```

#### From Kafka topics:
```go
// Publishing to "order.created" - who consumes this?
producer.SendMessage(&sarama.ProducerMessage{Topic: "order.created"})

// Consuming from "payment.completed" - who produces this?
consumer.Subscribe("payment.completed")
```

#### From HTTP clients:
```go
// This tells you the service calls "inventory-service"
resp, err := http.Get("http://inventory-service/api/stock")
```

#### From proto imports:
```protobuf
import "payment/payment.proto";  // Depends on payment service
```

**Build a communication map:**
```
Ledger Service:
  CALLS (sync):
    - PaymentService.ProcessPayment (gRPC)
    - InventoryService.CheckStock (HTTP)
  PUBLISHES (async):
    - ledger.transaction.created
    - ledger.balance.updated
  CONSUMES (async):
    - order.completed
    - payment.refunded
```

Include ALL these services in the diagram, even if they're not in the current directory.

### 6. Check Runbooks and Documentation

**MANDATORY: Do not skip this step.**

Look for and read these files:
- `README.md` - service overview, dependencies
- `RUNBOOK.md` - operational info, dependencies, topics
- `docs/*.md` - architecture docs
- `docs/architecture.md` - system design
- `config/*.yaml` - Kafka topics, service URLs
- `.env.example` - environment variables showing dependencies

**Extract from documentation:**
```
From RUNBOOK.md:
- "This service publishes to ledger.audit.events"
- "Depends on PaymentService for refund processing"
- "Consumes from compliance.check.required topic"

From config/kafka.yaml:
- topics:
    - ledger.transaction.created
    - ledger.balance.updated
```

**Cross-reference code with docs** - if docs mention a dependency not in code:
- Include it in the diagram with a note: `[from runbook - verify]`

### 7. Generate Mermaid Diagram

**CRITICAL: Always quote subgraph display names to avoid syntax errors.**

```mermaid
flowchart LR
    subgraph api-gateway ["API Gateway"]
        GW[Gateway Handler]
    end

    subgraph order-service ["Order Service"]
        O1[OrderHandler]
        O2[OrderProcessor]
        O3[KafkaPublisher]
    end

    subgraph payment-service ["Payment Service"]
        P1[gRPC Server]
        P2[PaymentProcessor]
        P3([kafka-consumer])
    end

    subgraph kafka ["Kafka Topics"]
        T1[(order.created)]
        T2[(payment.completed)]
    end

    subgraph external ["External Systems"]
        DB1[[PostgreSQL]]
        STRIPE[[Stripe API]]
    end

    %% Sync calls (gRPC/HTTP) - thick arrows
    GW ==>|HTTP| O1
    O2 ==>|gRPC| P1
    P2 ==>|HTTPS| STRIPE
    O2 ==>|SQL| DB1
    P2 ==>|SQL| DB1

    %% Async calls (Kafka) - dotted arrows
    O3 -.->|publish| T1
    T1 -.->|consume| P3
    P3 --> P2
    P2 -.->|publish| T2

    %% Internal calls - thin arrows
    O1 --> O2
    O2 --> O3
    P1 --> P2

    %% Styling
    style GW fill:#2f9e44
    style T1 fill:#1c7ed6
    style T2 fill:#1c7ed6
    style DB1 fill:#495057
    style STRIPE fill:#495057
```

### 6. Visual Convention Checklist

**Arrow Styles (MUST be correct):**
- [ ] `==>` for ALL gRPC calls
- [ ] `==>` for ALL HTTP/REST calls
- [ ] `==>` for ALL database queries
- [ ] `-.->` for ALL Kafka publishes
- [ ] `-.->` for ALL message queue operations
- [ ] `-->` for internal function calls only

**Node Shapes:**
- [ ] `[name]` rectangles for services/handlers
- [ ] `[(name)]` cylinders for Kafka topics
- [ ] `([name])` stadiums for consumer groups
- [ ] `[[name]]` double rectangles for external systems

**Colors:**
- [ ] Green `#2f9e44` for entry points
- [ ] Blue `#1c7ed6` for Kafka topics
- [ ] Purple `#9c36b5` for gRPC endpoints (optional)
- [ ] Gray `#495057` for external systems
- [ ] Red `#e03131` for DLQ/error paths
- [ ] Yellow `#f08c00` for warnings

### 7. Always Include Legend

Every diagram MUST include this legend:

```markdown
**Legend:**
| Arrow | Meaning | Debug Approach |
|-------|---------|----------------|
| `==>` | **Sync** (gRPC/HTTP) - caller blocks | Check latency, timeouts, error codes |
| `-.->` | **Async** (Kafka) - fire & forget | Check consumer lag, DLQ, offsets |
| `-->` | Internal call | Check logs, traces |
```

### 8. Identify and Highlight Issues

#### Sync Path Issues
- **No timeout configured**: gRPC/HTTP call without deadline
- **No retry logic**: Single attempt that can fail
- **No circuit breaker**: Can cascade failures
- **Blocking in async context**: Sync call inside Kafka consumer

#### Async Path Issues
- **No DLQ**: Failed messages have nowhere to go
- **No consumer**: Topic has no subscribers
- **Single consumer**: No redundancy for critical topics
- **No idempotency**: Consumer can't handle duplicates

#### Highlight with colors:
```mermaid
%% Warning: No timeout on gRPC call
style GRPC_CALL fill:#f08c00

%% Error: No DLQ configured
style CONSUMER fill:#e03131

%% Warning: Single consumer for critical topic
style SINGLE_CONSUMER fill:#f08c00
```

### 9. Output Structure

```markdown
# Service Flow: [Area/Feature Name]

## Diagram

\`\`\`mermaid
flowchart LR
    ...
\`\`\`

## Legend

| Arrow | Meaning | Debug Approach |
|-------|---------|----------------|
| `==>` | **Sync** (gRPC/HTTP) | Check latency, timeouts, error codes |
| `-.->` | **Async** (Kafka) | Check consumer lag, DLQ, offsets |
| `-->` | Internal | Check logs, traces |

## Communication Summary

### Sync Calls (gRPC/HTTP)
| From | To | Method | Timeout | Retry |
|------|-----|--------|---------|-------|
| OrderService | PaymentService | gRPC ProcessPayment | 5s | 3x |
| PaymentService | Stripe | HTTPS | 30s | 3x exp backoff |

### Async Messages (Kafka)
| Topic | Producers | Consumers | DLQ |
|-------|-----------|-----------|-----|
| order.created | OrderService | PaymentService, InventoryService | Yes |
| payment.completed | PaymentService | NotificationService | Yes |

## Issues Found

### Critical
- ❌ [Description of critical issue]

### Warnings
- ⚠️ [Description of warning]

## Incident Investigation Guide

### If sync path is failing:
1. Check error rates on [service] dashboard
2. Check latency percentiles
3. Verify circuit breaker state
4. Check downstream service health

### If async path is delayed:
1. Check consumer lag for [topic]
2. Check DLQ depth
3. Verify consumer group membership
4. Check for poison messages
```

### 10. Language-Specific Detection Patterns

#### Go - gRPC Client Detection
```go
// Look for these patterns - ALL ARE SYNC
grpc.Dial(
grpc.NewClient(
pb.New*Client(
client.*( // where client is a gRPC stub
```

#### Go - Kafka Detection
```go
// Look for these patterns - ALL ARE ASYNC
sarama.NewSyncProducer
sarama.NewAsyncProducer
kafka.NewProducer
producer.SendMessage
producer.Produce
sarama.NewConsumerGroup
```

#### Python - gRPC Detection
```python
# Look for these patterns - ALL ARE SYNC
grpc.insecure_channel
grpc.secure_channel
stub = *Stub(channel)
stub.*( # method calls on stub
```

#### Python - Kafka Detection
```python
# Look for these patterns - ALL ARE ASYNC
KafkaProducer
KafkaConsumer
producer.send
@app.agent # Faust
```

### 11. Proto File Analysis

When `.proto` files are present:
1. Extract service definitions to understand gRPC interface
2. Map RPC methods to actual calls in code
3. Include service methods in the diagram

```protobuf
service PaymentService {
    rpc ProcessPayment(PaymentRequest) returns (PaymentResponse);
    rpc RefundPayment(RefundRequest) returns (RefundResponse);
}
```

Maps to:
```mermaid
subgraph payment-service ["Payment Service - gRPC"]
    P1[ProcessPayment]
    P2[RefundPayment]
end
```
