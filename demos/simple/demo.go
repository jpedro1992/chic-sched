package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"unsafe"

	"k8s.io/klog/v2"

	"github.com/ibm/chic-sched/pkg/placement"
	"github.com/ibm/chic-sched/pkg/system"
	"github.com/ibm/chic-sched/pkg/topology"
	"github.com/ibm/chic-sched/pkg/util"
)

// Small system demo
//   - create a pTree and PEs given tree degrees at various levels
//   - generate random allocation using normal distribution
//   - generate a group with level constraints
//   - place group, generate lTree
//   - deallocate resources
func main() {
	klog.InitFlags(nil)
	flag.Set("v", "4")
	flag.Set("skip_headers", "true")
	klog.SetOutput(os.Stdout)
	flag.Parse()
	defer klog.Flush()

	// parameters
	numResources := 2
	maxLevel := 2
	numServers := 5
	serverCapacity := []int{16, 256}
	groupSize := 4
	demand := []int{4, 32}

	groupDemand, _ := util.NewAllocationCopy(demand)

	// create non-server nodes in the physical tree
	root := topology.NewPNode(topology.NewNode(&system.Entity{ID: "root"}), maxLevel, numResources)
	rack0 := topology.NewPNode(topology.NewNode(&system.Entity{ID: "rack-0"}), maxLevel-1, numResources)
	rack1 := topology.NewPNode(topology.NewNode(&system.Entity{ID: "rack-1"}), maxLevel-1, numResources)

	// create servers
	capacity, _ := util.NewAllocationCopy(serverCapacity)
	allocated, _ := util.NewAllocationCopy(demand)

	fmt.Println("Servers:")
	pes := make([]*system.PE, numServers)
	servers := make([]*topology.PNode, numServers)
	for i := 0; i < numServers; i++ {
		namei := "server" + strconv.FormatInt(int64(i), 10)
		pei := system.NewPE(namei, capacity)
		pes[i] = pei
		fmt.Println(pei)
		pei.SetAllocated(allocated.Clone())
		nodei := topology.NewNode((*system.Entity)(unsafe.Pointer(pei)))
		servers[i] = topology.NewPNode(nodei, 0, numResources)
	}
	fmt.Println()

	fmt.Println("Place some load:")
	pes[0].GetAllocated().Scale(3)
	pes[1].GetAllocated().Scale(2)
	pes[2].GetAllocated().Scale(1)
	pes[3].GetAllocated().Scale(2)
	pes[4].GetAllocated().Scale(3)
	for i := 0; i < numServers; i++ {
		fmt.Println(pes[i])
	}
	fmt.Println()

	// build topology
	fmt.Println("Build physical topology:")
	root.AddChild(&rack0.Node)
	root.AddChild(&rack1.Node)

	rack0.AddChild(&servers[0].Node)
	rack0.AddChild(&servers[1].Node)
	rack0.AddChild(&servers[2].Node)

	rack1.AddChild(&servers[3].Node)
	rack1.AddChild(&servers[4].Node)

	tree := topology.NewTree(&root.Node)
	pTree := topology.NewPTree(tree)
	pTree.PercolateResources()
	fmt.Print(pTree)

	// create level constraint
	fmt.Println("Create level constraints:")

	// constraint at rack level
	level := 1
	affinity := util.Pack
	isHard := false
	lc1 := placement.NewLevelConstraint("lc-1", level, util.Affinity(affinity), isHard)
	fmt.Println(lc1)

	//constraint at server level
	level = 0
	affinity = util.Spread
	isHard = false
	lc0 := placement.NewLevelConstraint("lc-0", level, util.Affinity(affinity), isHard)
	fmt.Println(lc0)
	fmt.Println()

	// create placement group
	fmt.Println("Create placement group:")
	pg := placement.NewPGroup("pg0", groupSize, groupDemand)
	pg.AddLevelConstraint(lc1)
	pg.AddLevelConstraint(lc0)

	// place group
	p := placement.NewPlacer(pTree)
	_, err := p.PlaceGroup(pg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(pg)

	// allocate resources
	fmt.Println("Allocate logical tree:")
	pg.ClaimAll(pTree)
	fmt.Print(pTree)
	fmt.Print(pg)

	fmt.Println("Servers:")
	for i := 0; i < numServers; i++ {
		fmt.Println(pes[i])
	}
	fmt.Println()

	// de-allocate resources
	fmt.Println("DeAllocate logical tree:")
	pg.UnClaimAll(pTree)
	fmt.Print(pTree)
	fmt.Print(pg)
}
