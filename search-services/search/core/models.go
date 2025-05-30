package core

import (
	"slices"
	"sync"
)

type Comics struct {
	ID       int      `db:"id"`
	URL      string   `db:"url"`
	Keywords []string `db:"words"`
	Score    int
}

type Index struct {
	index map[string][]int
	lock  sync.RWMutex
}

func NewIndex() *Index {
	return &Index{
		index: make(map[string][]int),
	}
}

func (i *Index) Clear() {
	i.lock.Lock()
	i.index = make(map[string][]int)
	i.lock.Unlock()
}

func (i *Index) Put(id int, keywords []string) {
	i.lock.Lock()
	for _, keyword := range keywords {
		i.index[keyword] = append(i.index[keyword], id)
	}
	i.lock.Unlock()
}

func (i *Index) Get(keyword string) []int {
	i.lock.RLock()
	defer i.lock.RUnlock()
	return slices.Clone(i.index[keyword])
}
