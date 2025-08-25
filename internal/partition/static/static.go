package static

import "hash/fnv"

type Owners struct {
	Primary   string
	Secondary string
}

type Partitioner struct {
	nodes []string
}

func New(nodes []string) *Partitioner {
	if len(nodes) < 2 {
		panic("at least 2 nodes are required")
	}
	return &Partitioner{nodes: append([]string(nil), nodes...)}
}

func (p *Partitioner) Members() []string { return append([]string(nil), p.nodes...) }

func (p *Partitioner) Add(node string) { p.nodes = append(p.nodes, node) }
func (p *Partitioner) Remove(node string) {
	out := p.nodes[:0]
	for _, n := range p.nodes {
		if n != node {
			out = append(out, n)
		}
	}
	p.nodes = out
}

// TODO: use consistent hashing
// Static modulo, simple but it moves many keys when membership changes
func (p *Partitioner) OwnersFor(key string) Owners {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	idx := int(h.Sum32()) % len(p.nodes)
	primary := p.nodes[idx]
	secondary := p.nodes[(idx+1)%len(p.nodes)]
	return Owners{Primary: primary, Secondary: secondary}
}
