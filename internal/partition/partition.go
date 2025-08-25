package partition

type Owners struct {
	Primary   string
	Secondary string
}

type Partitioner interface {
	OwnersFor(key string) Owners
	Members() []string
	Add(node string)
	Remove(node string)
}
