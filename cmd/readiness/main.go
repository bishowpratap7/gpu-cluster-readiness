// Command readiness is the gpu-cluster-readiness CLI.
//
// Day 1 state: prints the check plan. The node-inventory check
// (client-go, one goroutine per node) lands in Week 1.
package main

import "fmt"

type check struct {
	name string
	what string
}

var plan = []check{
	{"node-inventory", "nvidia.com/gpu capacity vs allocatable, per node"},
	{"device-plugin", "GPU Operator + device plugin DaemonSet health"},
	{"driver-cuda", "driver version present and consistent across nodes"},
	{"vllm-smoke", "a real model serves; tokens/sec recorded as evidence"},
	{"nccl-smoke", "GPU-to-GPU comms (clean skip when no GPU present)"},
	{"report", "JSON + Markdown certification report"},
}

func main() {
	fmt.Println("gpu-cluster-readiness — is this cluster actually ready?")
	fmt.Println()
	for i, c := range plan {
		fmt.Printf("  [%d] %-16s %s  (planned)\n", i+1, c.name, c.what)
	}
	fmt.Println()
	fmt.Println("Evidence over claims. Clean skips. Zeros are a pass on a GPU-less cluster.")
}
