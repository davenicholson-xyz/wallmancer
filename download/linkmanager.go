package download

import "sync"

type LinkManager struct {
	Links []string
	mu    sync.RWMutex
}

func NewLinkManager() *LinkManager {
	return &LinkManager{
		Links: make([]string, 0),
	}
}

func (lm *LinkManager) AddLinks(newLinks []string) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	lm.Links = append(lm.Links, newLinks...)
}
