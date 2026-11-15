# gpu-cluster-readiness

**A Go tool that certifies an NVIDIA GPU Kubernetes cluster as production-ready, and produces the report that proves it.**

Companies are buying GPU clusters faster than anyone can verify them. This tool runs the checks a platform engineer runs before sign-off, in minutes, with evidence attached.

## The four answers

- **What problem it solves:** "Installed" is not "ready." A node can report `nvidia.com/gpu: 8` in its capacity and still be unable to run a single workload, because the device plugin is down, the driver is mismatched, or the GPUs throttle under load. Today that gap is checked by hand, cluster by cluster, if it is checked at all.
- **What I am building:** A CLI first, then a Kubernetes operator. It checks node GPU inventory (capacity vs allocatable), GPU Operator and device-plugin health, driver/CUDA sanity, runs a vLLM smoke benchmark (tokens/sec), runs an NCCL smoke test that skips cleanly when no GPUs are present, and emits a certification report in JSON + Markdown.
- **One technical decision, and why:** Go + client-go, one goroutine per node. Readiness checking is hundreds of parallel Kubernetes API calls that mostly wait on the network. Goroutines make that fan-out cheap and keep the code straight-line, and Go is the language the platform teams who run these clusters already write.
- **A result / what I learned:** This repo is the Day 1 scaffold of a 30-day build. The first real check (node GPU inventory) lands in Week 1. Measured results will be written here as they exist, never before.

## The checks

| # | Check | What it proves | Status |
|---|------|---------------|--------|
| 1 | Node GPU inventory | Every node's `nvidia.com/gpu` capacity vs allocatable, per node and totalled | Week 1 |
| 2 | Device plugin + GPU Operator health | The plugin DaemonSet is running everywhere, so claimed GPUs are actually schedulable | Week 1 |
| 3 | Driver / CUDA sanity | Driver version present and consistent across nodes | Week 2 |
| 4 | vLLM smoke benchmark | A real model serves on the cluster; tokens/sec recorded as evidence | Week 2 |
| 5 | NCCL smoke test | GPUs can talk to each other (clean skip when no GPU is present) | Week 2 |
| 6 | Certification report | JSON + Markdown report: per-check verdicts with the evidence behind each | Week 1–2 |

Weeks 3–4: a kubebuilder operator (a `ClusterReadiness` CRD whose controller runs these checks on a schedule and publishes status conditions), scheduler/DRA notes, and the write-ups.

## Learn it: the tutorial

**[GPU Cluster Readiness, explained with pictures and animations](tutorial/)** — a visual tutorial that teaches the whole idea: why capacity is not readiness, what each check proves, and how the certification report is built. Start there if you are new to GPU clusters.

## Run it

```bash
go run ./cmd/readiness
```

Day 1 state: prints the check plan. On a cluster with no GPUs (a laptop `kind` cluster counts), a correct inventory result is **zero GPUs, reported honestly**. Zeros are a pass at this stage; a tool that invents GPUs would fail the only test that matters.

## Design principles

- **Evidence over claims.** Every check returns machine-readable evidence, not a vibe. The report is the product.
- **Clean skips.** A check that cannot run (no GPU, no NCCL) says so and skips. It never fails loudly and never passes silently.
- **Boring code.** Straight-line Go, one goroutine per node, errors checked at every call. This tool gets read in interviews; it should read well.

## About

Built by **Bishow Pandey** — Senior Platform Engineer (AWS EKS/Kubernetes, Java/Spring Boot, application security) — as the flagship build of a 30-day NVIDIA prep, aimed at the DSX cluster-readiness lane. Part of the public [AI Engineer Journey](https://github.com/bishowpratap7/ai-engineer-journey).
