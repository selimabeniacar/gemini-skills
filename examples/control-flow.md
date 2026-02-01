# Control Flow Examples

## Example 1: Go Error Handling Function

**Command**: `/flow payments/processor.go --type=control --function=ChargeCustomer`

**Source Code**:
```go
// Line 25
func ChargeCustomer(customer *Customer, amount float64) error {
    // Line 26
    if amount <= 0 {
        return errors.New("invalid amount")
    }

    // Line 30
    card, err := customer.GetPaymentMethod()
    if err != nil {
        return fmt.Errorf("payment method error: %w", err)
    }

    // Line 35
    if !card.IsValid() {
        return errors.New("invalid card")
    }

    // Line 40
    charge, err := stripe.CreateCharge(card, amount)
    if err != nil {
        return fmt.Errorf("charge failed: %w", err)
    }

    // Line 45
    if err := db.SaveTransaction(charge); err != nil {
        // Note: Charge succeeded but save failed - potential issue!
        return fmt.Errorf("save failed: %w", err)
    }

    // Line 50
    return nil
}
```

**Output**:

```mermaid
flowchart TD
    A[Start: ChargeCustomer<br/>Line 25] --> B{amount <= 0?<br/>Line 26}
    B -->|Yes| C[Return Error:<br/>invalid amount]
    B -->|No| D[GetPaymentMethod<br/>Line 30]
    D --> E{err != nil?}
    E -->|Yes| F[Return Error:<br/>payment method error]
    E -->|No| G{card.IsValid?<br/>Line 35}
    G -->|No| H[Return Error:<br/>invalid card]
    G -->|Yes| I[stripe.CreateCharge<br/>Line 40]
    I --> J{err != nil?}
    J -->|Yes| K[Return Error:<br/>charge failed]
    J -->|No| L[db.SaveTransaction<br/>Line 45]
    L --> M{err != nil?}
    M -->|Yes| N[Return Error:<br/>save failed]
    M -->|No| O[Return nil<br/>Line 50]

    style A fill:#69db7c
    style O fill:#69db7c
    style C fill:#ff6b6b
    style F fill:#ff6b6b
    style H fill:#ff6b6b
    style K fill:#ff6b6b
    style N fill:#ffd43b
```

**Analysis**:
- **Location**: `payments/processor.go:25-50`
- **Cyclomatic Complexity**: 6 (5 decision points + 1)
- **Branches**: 6 possible exit points (5 errors, 1 success)
- **Issue Found**: Line 45 - If `SaveTransaction` fails, the charge has already succeeded with Stripe. This could lead to:
  - Customer charged but no record
  - Potential need for refund logic
  - Consider: wrap in transaction or add compensation logic

---

## Example 2: Python Async Handler

**Command**: `/flow api/handlers.py --type=control --function=handle_upload`

**Source Code**:
```python
# Line 15
async def handle_upload(request):
    # Line 16
    if not request.user.is_authenticated:
        raise AuthenticationError("Login required")

    # Line 20
    file = request.files.get('document')
    if not file:
        raise ValidationError("No file provided")

    # Line 25
    try:
        # Line 26
        validated = await validate_file(file)

        # Line 28
        if validated.size > MAX_SIZE:
            raise ValidationError("File too large")

        # Line 32
        stored = await storage.upload(validated)

        # Line 34
        record = await db.create_document(
            user=request.user,
            path=stored.path
        )

        # Line 39
        return {"id": record.id, "url": stored.url}

    except StorageError as e:
        # Line 42
        logger.error(f"Storage failed: {e}")
        raise ServerError("Upload failed")

    except DatabaseError as e:
        # Line 46
        await storage.delete(stored.path)  # Cleanup
        logger.error(f"DB failed: {e}")
        raise ServerError("Save failed")
```

**Output**:

```mermaid
flowchart TD
    A[Start: handle_upload<br/>Line 15] --> B{user.is_authenticated?<br/>Line 16}
    B -->|No| C[Raise AuthenticationError]
    B -->|Yes| D[Get file from request<br/>Line 20]
    D --> E{file exists?}
    E -->|No| F[Raise ValidationError:<br/>No file provided]
    E -->|Yes| G[try block<br/>Line 25]
    G --> H[validate_file<br/>Line 26]
    H --> I{size > MAX_SIZE?<br/>Line 28}
    I -->|Yes| J[Raise ValidationError:<br/>File too large]
    I -->|No| K[storage.upload<br/>Line 32]
    K --> L[db.create_document<br/>Line 34]
    L --> M[Return success<br/>Line 39]

    G -.->|StorageError| N[Log error<br/>Line 42]
    N --> O[Raise ServerError:<br/>Upload failed]

    G -.->|DatabaseError| P[storage.delete cleanup<br/>Line 46]
    P --> Q[Log error]
    Q --> R[Raise ServerError:<br/>Save failed]

    style A fill:#69db7c
    style M fill:#69db7c
    style C fill:#ff6b6b
    style F fill:#ff6b6b
    style J fill:#ff6b6b
    style O fill:#ff6b6b
    style R fill:#ffd43b
```

**Analysis**:
- **Location**: `api/handlers.py:15-47`
- **Cyclomatic Complexity**: 7
- **Exception Handling**: Good - catches specific exceptions
- **Cleanup Logic**: Present for DatabaseError (deletes uploaded file)
- **Potential Issue**: StorageError handler doesn't need cleanup (upload didn't complete), but verify `validate_file` doesn't create temporary files

---

## Example 3: Go Switch Statement

**Command**: `/flow router/handler.go --type=control --function=HandleMethod`

**Output**:

```mermaid
flowchart TD
    A[Start: HandleMethod<br/>Line 10] --> B{switch method}
    B -->|GET| C[handleGet<br/>Line 13]
    B -->|POST| D[handlePost<br/>Line 15]
    B -->|PUT| E[handlePut<br/>Line 17]
    B -->|DELETE| F{user.IsAdmin?<br/>Line 19}
    B -->|default| G[Return 405<br/>Method Not Allowed]

    F -->|Yes| H[handleDelete<br/>Line 21]
    F -->|No| I[Return 403<br/>Forbidden]

    C --> J[Return response]
    D --> J
    E --> J
    H --> J

    style A fill:#69db7c
    style G fill:#ffd43b
    style I fill:#ff6b6b
```

**Analysis**:
- **Location**: `router/handler.go:10-25`
- **All HTTP methods handled**: GET, POST, PUT, DELETE + default
- **Authorization**: DELETE requires admin rights
- **Note**: Consider if PUT should also require authorization
