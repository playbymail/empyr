// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cluster

type Cluster struct {
	Radius  float64
	Systems []string // id for every system in the cluster
	Stars   []string // id for every star in the cluster
}
