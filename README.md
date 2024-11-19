# chic-sched

Topology-aware group placement algorithm

**chic-sched** is a generic scheduler which uses a heuristics-based, topology-aware placement algorithm for HPC Placement Groups with level constraints.

- [Chic-sched: a HPC Placement-Group Scheduler on Hierarchical Topologies with Constraints](https://ieeexplore.ieee.org/document/10177454)

- [Description of the algorithm](docs/PG-sched-algo.pdf)

- [Modeling and analysis of placement group heuristics](docs/heuristics-modeling.pdf)

- [HPC placement group policies: requirements, specifications, and comparisons](docs/hpc-placement-policies.pdf)

Background on placement group offerings by [AWS](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/placement-groups.html), [Azure](https://docs.microsoft.com/en-us/azure/virtual-machine-scale-sets/virtual-machine-scale-sets-placement-groups), and [Google Cloud](https://cloud.google.com/compute/docs/instances/define-instance-placement).

An [example](demos/treebuild/demo.go) is provided.
