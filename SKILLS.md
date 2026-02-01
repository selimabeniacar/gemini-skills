# /flow

Generate flow diagrams for debugging, incident investigation, and onboarding.

## Description

Creates Mermaid diagrams to visualize code execution and service communication. Supports three modes:

- **Call Graph**: Trace function calls within a service to find where issues originate
- **Control Flow**: Zoom into a specific function to understand branching logic
- **Service Flow**: Map communication between services showing both **sync (gRPC/HTTP)** and **async (Kafka/queues)** paths

## Critical: Sync vs Async Distinction

For incident investigation and onboarding, it is **essential** to clearly distinguish:

| Type | Examples | Arrow Style | Behavior | Debug Approach |
|------|----------|-------------|----------|----------------|
| **Sync** | gRPC, HTTP, REST | `==>` thick solid | Caller blocks, waits for response | Check latency, timeouts, error codes |
| **Async** | Kafka, RabbitMQ, SQS | `-.->` dotted | Fire and forget, eventual | Check consumer lag, DLQ, offsets |

**This distinction changes everything about how you debug an issue.**

## Arguments

- `<path>` - File or directory to analyze (required)
- `--type=call|control|service` - Type of diagram (default: `call`)
- `--function=<name>` - Specific function to analyze (required for control flow)
- `--depth=<n>` - How deep to trace calls (default: 5, only for call graphs)
- `--topic=<name>` - Filter by specific Kafka topic (optional, for service flows)
- `--service=<name>` - Filter by service name (optional, for service flows)
- `--output=<path>` - Output file path (default: `/tmp/flow-<timestamp>.md`)
- `--verbose` - Include educational explanations for onboarding (default: off)

## Supported Languages

- Go
- Python

## Usage

```bash
# Generate call graph for a directory (outputs to /tmp/flow-<timestamp>.md)
/flow src/payments/ --type=call

# Generate call graph for a specific file
/flow src/payments/processor.go --type=call

# Generate control flow for a specific function
/flow src/payments/processor.go --type=control --function=ProcessPayment

# Short form (defaults to call graph)
/flow src/payments/

# Specify custom output path
/flow src/payments/ --output=./my-diagram.md

# Verbose mode for onboarding (includes educational explanations)
/flow src/payments/ --verbose

# Generate service flow showing sync and async communication
/flow services/ --type=service

# Service flow filtered by topic
/flow services/ --type=service --topic=order.created

# Service flow for a specific service
/flow services/payment/ --type=service --service=payment-service
```

## Visual Convention Reference (MUST Follow)

### Arrow Styles - CRITICAL

```
A ==>|gRPC| B        Sync: caller waits, immediate response or timeout
A -.->|Kafka| B      Async: fire and forget, processed later
A -->|internal| B    Internal: function call within same service
```

### Node Shapes

```
[Service/Handler]    Rectangle: services, handlers, processors
[(topic-name)]       Cylinder: Kafka topics, message queues
([consumer-group])   Stadium: consumer groups
{decision}           Diamond: routing decisions, conditions
[[external]]         Double rectangle: external systems (Stripe, DB)
```

### Colors

| Color | Hex | Meaning |
|-------|-----|---------|
| Green | `#69db7c` | Entry points, success paths |
| Blue | `#4dabf7` | Kafka topics, async infrastructure |
| Red | `#ff6b6b` | Error paths, DLQ, failures |
| Yellow | `#ffd43b` | Warnings: missing consumer, no error handling |
| Purple | `#be4bdb` | gRPC/sync calls between services |
| Orange | `#fd7e14` | Retry paths, circuit breakers |
| Gray | `#868e96` | External systems |

### Service Boundaries

Always use subgraphs to show service boundaries:
```
subgraph ServiceName [Service Name]
    components...
end
```

## Instructions

When this skill is invoked:

1. **Parse the arguments** to determine:
   - The target path (file or directory)
   - The diagram type (`call`, `control`, or `service`)
   - The function name (if control flow)
   - The depth (if call graph)
   - The topic filter (if service flow)
   - The service filter (if service flow)
   - The output path (default: `/tmp/flow-<timestamp>.md`)
   - Whether verbose mode is enabled

2. **Read the target files**:
   - If a directory, read all `.go`, `.py`, and `.proto` files
   - If a file, read just that file
   - Identify the programming language from file extensions
   - Look for proto files to understand gRPC service definitions

3. **Generate the appropriate diagram** based on type:

### For Call Graphs (`--type=call`)

Follow the instructions in `prompts/call-graph.md`:
- Identify all function definitions
- Trace function calls within each function up to the specified depth
- Build a dependency graph
- Output as Mermaid flowchart
- Highlight suspicious paths (missing error handling, potential issues)

### For Control Flow Graphs (`--type=control`)

Follow the instructions in `prompts/control-flow.md`:
- Locate the specified function
- Parse control structures (if/else, switch, loops, try/catch)
- Identify all branches and exit points
- Output as Mermaid flowchart with decision nodes
- Annotate with line numbers

### For Service Flows (`--type=service`)

Follow the instructions in `prompts/service-flow.md`:
- **Identify sync communication**: gRPC clients, HTTP clients, REST calls
- **Identify async communication**: Kafka producers/consumers, message queues
- **Map service boundaries**: Which code belongs to which service
- **Use correct arrow styles**: `==>` for sync, `-.->` for async (THIS IS CRITICAL)
- **Show failure modes**: Timeouts, DLQs, retries, circuit breakers
- **Highlight issues**: Missing error handling, single points of failure

4. **Write the output file**:
   - Generate a filename using the format `/tmp/flow-<timestamp>.md` (e.g., `/tmp/flow-20240115-143022.md`)
   - Or use the path specified by `--output`
   - **Always include the legend** showing sync vs async distinction
   - Write the Mermaid diagram in a fenced code block
   - Include analysis notes after the diagram
   - Tell the user the output file path so they can open it

5. **Add analysis notes** in the output file explaining:
   - Key observations about the flow
   - Potential issues or areas of concern
   - **For incidents**: What to check based on sync vs async paths
   - Suggestions for debugging if applicable

6. **If `--verbose` is enabled**, add educational content for onboarding:

### Verbose Mode Content

Include these additional sections in the output:

#### How to Read This Diagram
- **Sync vs Async**: Explain thick arrows (==> gRPC) block the caller, dotted arrows (-.-> Kafka) don't
- Explain node shapes and what they represent
- Explain color coding and what each color means
- Describe how to trace a request through the system

#### Why Sync vs Async Matters
- Sync failures: Immediate, caller gets error, check timeouts and retries
- Async failures: Delayed, check consumer lag, DLQ, message schemas
- How this affects incident debugging

#### Glossary
- Define key terms specific to the codebase
- Explain domain-specific concepts visible in the diagram
- Define infrastructure terms (consumer group, partition, DLQ, circuit breaker)

#### Architecture Context
- Explain where this code fits in the overall system
- Describe which services own which data
- Note any important patterns (saga, outbox, CQRS)

#### Incident Investigation Guide
- For sync path issues: What to check (latency, error rates, circuit breaker state)
- For async path issues: What to check (consumer lag, DLQ depth, partition assignment)
- Common failure modes in this architecture

#### Next Steps for Learning
- Suggest related files or functions to explore
- Recommend what to trace next to understand the full flow
- Point to relevant documentation if known

## Examples

See the examples directory for detailed output examples:
- `examples/call-graph.md` - Call graph examples
- `examples/control-flow.md` - Control flow examples
- `examples/service-flow.md` - Service flow examples with sync/async distinction

## Output Reliability Checklist

Before finalizing output, verify:
- [ ] All sync calls use `==>` arrows
- [ ] All async calls use `-.->` arrows
- [ ] Service boundaries are shown with subgraphs
- [ ] Legend is included showing arrow meanings
- [ ] Colors are applied consistently
- [ ] Issues are highlighted with appropriate colors
- [ ] Line numbers or file references are included where helpful
