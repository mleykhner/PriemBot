package faq

import (
	"PriemBot/config"
	"PriemBot/faq/models"
	"log"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

type Manager struct {
	path     string
	faqList  models.FAQList
	faqIndex map[int]models.FAQ
	mu       sync.RWMutex
}

func NewFAQManager(config *config.FAQConfig) (*Manager, error) {
	m := &Manager{path: config.FilePath}
	if err := m.reload(); err != nil {
		return nil, err
	}
	go m.watch()
	return m, nil
}

func (m *Manager) reload() error {
	data, err := os.ReadFile(m.path)
	if err != nil {
		return err
	}
	var list models.FAQList
	if err := yaml.Unmarshal(data, &list); err != nil {
		return err
	}
	idx := make(map[int]models.FAQ)
	for _, f := range list.FAQs {
		idx[f.ID] = f
	}
	m.mu.Lock()
	m.faqList = list
	m.faqIndex = idx
	m.mu.Unlock()
	return nil
}

func (m *Manager) watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("FAQ watcher init error:", err)
		return
	}
	defer watcher.Close()
	if err := watcher.Add(m.path); err != nil {
		log.Println("FAQ watcher file add error:", err)
		return
	}
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				log.Println("FAQ file updated, reloading...")
				if err := m.reload(); err != nil {
					log.Printf("FAQ reload failed: %v", err)
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("FAQ watcher error:", err)
		}
	}
}

func (m *Manager) List() []models.FAQ {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]models.FAQ(nil), m.faqList.FAQs...) // копия
}

func (m *Manager) Get(id int) (models.FAQ, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	f, ok := m.faqIndex[id]
	return f, ok
}

func (m *Manager) Info() string {
	return m.faqList.Info
}
