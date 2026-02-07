# Phase 3: Validation & Cleanup

## Objective

Validate the diagram with flowlint, rework if needed, render PNG, and clean up.

---

## Prerequisites

Build flowlint (one-time):
```bash
cd flow/tools/flowlint
go build -o flowlint .
```

---

## Step 3.1: Run flowlint

**YOU MUST EXECUTE THIS COMMAND VIA BASH:**

```bash
./flow/tools/flowlint/flowlint refine {output}.md .flow-deps.yaml
```

Or if flowlint is in PATH:
```bash
flowlint refine {output}.md .flow-deps.yaml
```

**DO NOT SKIP THIS STEP. Actually run the command and check the output.**

flowlint will:
1. Validate Mermaid syntax with mermaid-cli
2. Check style guide compliance and auto-fix
3. Verify all dependencies are represented
4. Report PASS or FAIL

---

## Step 3.2: Handle Results

**Read the flowlint output carefully.**

### If PASS (exit code 0):

```
✓ Refinement complete - diagram is ready
```

Cleanup:
```bash
rm .flow-deps.yaml
```

Done. Proceed to Step 3.3.

### If FAIL (exit code 1):

**Read the error messages.** Then attempt **one rework cycle**:

| Failure Type | AI Action |
|--------------|-----------|
| Syntax error | Edit {output}.md to fix the Mermaid syntax |
| Missing dependencies | Edit {output}.md to add missing nodes/edges |
| Style issues | Edit {output}.md to fix arrows/shapes/colors |

After editing, re-run flowlint:
```bash
./flow/tools/flowlint/flowlint refine {output}.md .flow-deps.yaml
```

### If still FAIL after rework:

```
⚠️ Validation incomplete - diagram generated with issues.

Output: {output}.md (diagram created, may have issues)
Kept: .flow-deps.yaml (for debugging)

Issues:
- {list of remaining issues}

The diagram is usable but may need manual fixes.
Run: flowlint lint {output}.md --fix
```

**The diagram is ALWAYS generated.** Validation failure = warnings, not blockers.

---

## Step 3.3: Generate PNG (if mermaid-cli available)

After validation, generate a high-resolution PNG:

```bash
npx -p @mermaid-js/mermaid-cli mmdc -i {output}.md -o {output}.png -b white -w 3840 -s 2
```

Options:
- `-w 3840` — 4K width for high resolution
- `-s 2` — 2x scale factor for crisp text
- `-b white` — white background

If mermaid-cli is not installed, this will auto-install via npx. Skip if npm/npx is not available.

---

## Finding flowlint

flowlint is located at:
```
flow/tools/flowlint/flowlint
```

If you get "command not found", use the full path relative to the skill directory.

---

## Completion

### Success

```
Phase 3 Complete: Validation

Status: PASS
Coverage: {found}/{total} (100%)

Cleaned up: .flow-deps.yaml

Output:
- {output}.md (diagram source)
- {output}.png (rendered image, if mermaid-cli available)

If PNG not generated, render manually:
npx -p @mermaid-js/mermaid-cli mmdc -i {output}.md -o {output}.png -b white -w 3840 -s 2
```

### Partial Success

```
Phase 3 Complete: Validation

Status: PARTIAL
Coverage: {found}/{total}

Output:
- {output}.md (diagram source, may have issues)
- .flow-deps.yaml (kept for debugging)

Run manually: flowlint lint {output}.md --fix
```
