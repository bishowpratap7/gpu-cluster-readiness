# Architecture

## Pipeline

```
kubeconfig ──► client-go ──► checks (one goroutine per node) ──► results ──► report (JSON + Markdown)
```

Every check implements the same contract:

```go
type Result struct {
    Check    string            // "node-inventory"
    Status   string            // PASS | FAIL | SKIP
    Evidence map[string]string // the numbers behind the verdict
    Detail   string            // one human sentence
}
```

A check that cannot run returns `SKIP` with the reason in `Detail`. It never
guesses, and it never converts missing evidence into a pass.

## The checks

1. **node-inventory** — Lists Nodes via the typed clientset. For each node,
   reads `status.capacity["nvidia.com/gpu"]` and
   `status.allocatable["nvidia.com/gpu"]`. A node whose capacity claims GPUs
   but whose allocatable is 0 is flagged: something below the API (usually
   the device plugin) is broken.
2. **device-plugin** — Verifies the GPU Operator's device-plugin DaemonSet
   has a ready pod on every GPU node.
3. **driver-cuda** — Reads the `nvidia.com/gpu` product label and driver
   version from node labels; flags version skew across the fleet.
4. **vllm-smoke** — Deploys a small vLLM pod with a tiny model, sends a fixed
   prompt, records prompt + generation tokens/sec from the server's own
   counters. Deletes the pod afterwards.
5. **nccl-smoke** — Runs a two-GPU all-reduce smoke job where GPUs exist;
   `SKIP` (cleanly) everywhere else.

## The operator (Weeks 3–4)

A kubebuilder-built controller with one CRD:

```yaml
apiVersion: readiness.bishow.dev/v1alpha1
kind: ClusterReadiness
spec:
  schedule: "0 */6 * * *"
status:
  conditions: [...]   # per-check verdicts, same Result contract
  lastReport: "..."   # link/name of the generated report
```

The operator reuses the CLI's check package unchanged. The CLI proves the
checks; the operator makes them continuous.

## Why Go

The checks are hundreds of concurrent, mostly-idle API calls. Goroutines
(one per node) keep that fan-out cheap and the code sequential-looking.
client-go is the native client; there is no second choice worth defending.
