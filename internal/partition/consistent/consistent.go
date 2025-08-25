package consistent

import (
	"hash/fnv"
	"sort"
	"strconv"

	"github.com/Hasan-Al-Khazraji/CuckooCache/internal/partition"
)

type Ring struct {
	points      []uint32
	pointToNode map[uint32]string // virtual node hash -> physical node
	nodes       map[string]struct{}
	vnodes      int
}

func New(nodes []string, vnodes int) *Ring {
	if len(nodes) < 2 {
		panic("at least 2 nodes are required")
	}
	if vnodes <= 0 {
		vnodes = 128
	}
	r := &Ring{
		pointToNode: make(map[uint32]string),
		nodes:       make(map[string]struct{}),
		vnodes:      vnodes,
	}
	for _, n := range nodes {
		r.nodes[n] = struct{}{}
	}
	r.rebuild()
	return r
}

func (r *Ring) Members() []string {
	out := make([]string, 0, len(r.nodes))
	for n := range r.nodes {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}

func (r *Ring) Add(node string) {
	if _, ok := r.nodes[node]; ok {
		return
	}
	r.nodes[node] = struct{}{}
	r.rebuild()
}

func (r *Ring) Remove(node string) {
	if _, ok := r.nodes[node]; !ok {
		return
	}
	delete(r.nodes, node)
	r.rebuild()
}

func (r *Ring) OwnersFor(key string) partition.Owners {
	if len(r.points) == 0 {
		return partition.Owners{}
	}
	h := hash32(key)
	i := sort.Search(len(r.points), func(i int) bool { return r.points[i] >= h })
	if i == len(r.points) {
		i = 0
	}
	primary := r.pointToNode[r.points[i]]

	j := (i + 1) % len(r.points)
	secondary := primary

	for j != i {
		n := r.pointToNode[r.points[j]]
		if n != primary {
			secondary = n
			break
		}
		j = (j + 1) % len(r.points)
	}

	return partition.Owners{Primary: primary, Secondary: secondary}
}

// Helpers

func (r *Ring) rebuild() {
	r.points = r.points[:0]
	for n := range r.nodes {
		for v := 0; v < r.vnodes; v++ {
			point := hash32(n + "#" + strconv.Itoa(v))
			r.pointToNode[point] = n
			r.points = append(r.points, point)
		}
	}
	sort.Slice(r.points, func(i, j int) bool { return r.points[i] < r.points[j] })
}

func hash32(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return h.Sum32()
}
