package main

import (
	"fmt"
	"math"
	"strings"
	"time"
	"unsafe"

	"github.com/ibm/chic-sched/demos"
	"github.com/ibm/chic-sched/pkg/builder"
	"github.com/ibm/chic-sched/pkg/placement"
	"github.com/ibm/chic-sched/pkg/topology"
	"github.com/ibm/chic-sched/pkg/util"
)

// Scaling demo
//
//	 evaluate the placement algorithm as a function of system and group sizes
//		- create a pTree and PEs given tree degrees at various levels
//		- generate random allocation with given distribution
//		- generate a group with level constraints
//		- place group, generate lTree
func main() {

	// run parameters
	isPrint := false

	// resource parameters
	capacity := []int{16}
	demand := []int{2}

	// loading parameters
	loadFactor := 0.4
	alpha := 0.250
	beta := 0.125
	cov := 1.5
	avg := demos.ComputeAverageLoad(loadFactor, alpha, beta)

	// constraint parameters
	// constraint at rack level
	level := 1
	affinity := util.Spread
	isHard := false
	lc1 := placement.NewLevelConstraint("lc-1", level, util.Affinity(affinity), isHard)

	//constraint at server level
	level = 0
	affinity = util.Pack
	isHard = false
	lc0 := placement.NewLevelConstraint("lc-0", level, util.Affinity(affinity), isHard)

	// misc parameters
	lineLength := 64

	// experiments parameters
	degreeBase := []int{1, 1, 2}
	groupSizeBase := 8
	minSizeFactor := 1
	maxSizeFactor := 6
	numExperiments := 10000

	// print configuration
	if !isPrint {
		fmt.Println(strings.Repeat("=", lineLength))
		fmt.Println("Configuration parameters: ")
		fmt.Printf("capacity=%v; demand=%v \n", capacity, demand)
		fmt.Printf("loadFactor=%f; alpha=%f; beta=%f \n", loadFactor, alpha, beta)
		fmt.Printf("avg=%f; cov=%f \n", avg, cov)
		fmt.Println(lc1)
		fmt.Println(lc0)
		fmt.Printf("degreeBase=%v; groupSizeBase=%d; minSizeFactor=%d; maxSizeFactor=%d; numExperiments=%d \n",
			degreeBase, groupSizeBase, minSizeFactor, maxSizeFactor, numExperiments)
		fmt.Println(strings.Repeat("=", lineLength))
	}

	degree := make([]int, len(degreeBase))
	groupSize := 1

	for sizefactor := minSizeFactor; sizefactor <= maxSizeFactor; sizefactor++ {
		size := int(math.Pow(2, float64(sizefactor)))
		nserv := 1
		degree[0] = 1
		for i := 1; i < len(degreeBase); i++ {
			degree[i] = size * degreeBase[i]
			nserv *= degree[i]
		}
		groupSize = size * groupSizeBase
		fmt.Printf("numServers=%v; degree=%v; groupSize=%v; ", nserv, degree, groupSize)

		bw, err := util.NewBoxWhisker(numExperiments)
		if err != nil {
			fmt.Printf("failure creating a BoxWhisker object")
			break
		}

		experiment := 0
		avgDuration := int64(0)
		isPlaced := true
		numFailures := 0
		for experiment < numExperiments {

			// create placement group
			groupDemand, _ := util.NewAllocationCopy(demand)
			pg := placement.NewPGroup("pg0", groupSize, groupDemand)
			pg.AddLevelConstraint(lc1)
			pg.AddLevelConstraint(lc0)

			// create physical tree and servers
			if isPrint {
				fmt.Println(strings.Repeat("=", lineLength))
				fmt.Print("Create physical tree topology: ")
				fmt.Printf("(degree=%v; capacity=%v)", degree, capacity)
				fmt.Println()
				fmt.Println(strings.Repeat("=", lineLength))
			}

			tg := builder.NewTreeGen()
			pTree := tg.CreateUniformTree(degree, capacity)
			//fmt.Print(pTree)
			pes := tg.GetPEs()
			numServers := len(pes)

			// place some random allocation on the servers
			if isPrint {
				fmt.Println(strings.Repeat("=", lineLength))
				fmt.Println("Place some load: ")
				fmt.Printf("(loadFactor=%f; alpha=%f; beta=%f; \n", loadFactor, alpha, beta)
				fmt.Printf("avg=%f; cov=%f)", avg, cov)
				fmt.Println()
				fmt.Println(strings.Repeat("=", lineLength))
			}

			demos.PlaceBackgroungLoad(pes, loadFactor, alpha, beta, cov)

			start := time.Now()
			pTree.PercolateResources()
			//fmt.Print(pTree)

			if isPrint {
				fmt.Println(strings.Repeat("=", lineLength))
				fmt.Println("Servers:")
				fmt.Println(strings.Repeat("=", lineLength))
				for i := 0; i < numServers; i++ {
					fmt.Println(pes[i])
				}
			}

			// print level constraint
			if isPrint {
				fmt.Println(strings.Repeat("=", lineLength))
				fmt.Println("Create level constraint:")
				fmt.Println(strings.Repeat("=", lineLength))
				fmt.Println(lc1)
				fmt.Println(lc0)
				fmt.Println()
			}

			// print placement group
			if isPrint {
				fmt.Println(strings.Repeat("=", lineLength))
				fmt.Println("Create placement group:")
				fmt.Println(strings.Repeat("=", lineLength))
				fmt.Println(pg)
				fmt.Println()
			}

			// place group
			p := placement.NewPlacer(pTree)
			ltree, err := p.PlaceGroup(pg)
			isPlaced = isPlaced && (err == nil)
			if isPlaced {
				lRoot := (*topology.LNode)(unsafe.Pointer(ltree.GetRoot()))
				if lRoot.GetCount() == groupSize {
					experiment++
					duration := time.Since(start).Microseconds()
					bw.AddSample(int(duration))
					avgDuration += duration
				} else {
					numFailures++
				}
			}

			if isPrint {
				fmt.Println(strings.Repeat("=", lineLength))
				fmt.Println("Place group result: (logical tree)")
				fmt.Println(strings.Repeat("=", lineLength))
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Print(pg)
			}

			pg.ClaimAll(pTree)

			if isPrint {
				fmt.Println(strings.Repeat("=", lineLength))
				fmt.Println("Physical tree after group allocation:")
				fmt.Println(strings.Repeat("=", lineLength))
				//fmt.Print(pTree)

				fmt.Println(strings.Repeat("=", lineLength))
				fmt.Println("Allocation on servers:")
				fmt.Println(strings.Repeat("=", lineLength))
				for i := 0; i < numServers; i++ {
					fmt.Println(pes[i])
				}
			}
		}

		fmt.Printf("numFailures=%d; ", numFailures)
		fmt.Printf("duration (microsec)=%d; ", avgDuration/int64(numExperiments))

		bw.Calculate()
		fmt.Println(bw)
	}
}
