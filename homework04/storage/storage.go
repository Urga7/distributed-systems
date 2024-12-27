package storage

import (
	"errors"
	"sync"
)

type Todo struct {
	Task      string `json:"task"`
	Value     string `json:"value"`
	Committed bool   `json:"committed"`
}

type TodoStorage struct {
	dict map[string]Todo
	lock sync.RWMutex
}

var ErrorNotFound = errors.New("not found")

func NewTodoStorage() *TodoStorage {
	dict := make(map[string]Todo)
	return &TodoStorage{
		dict: dict,
	}
}

func NewTodo(task string, value string) *Todo {
	return &Todo{
		Task:      task,
		Value:     value,
		Committed: false,
	}
}

func (tds *TodoStorage) Put(key, value string) {
	tds.lock.Lock()
	defer tds.lock.Unlock()
	tds.dict[key] = Todo{
		Task:      key,
		Value:     value,
		Committed: false,
	}
}

func (tds *TodoStorage) Get(key string) (string, bool) {
	tds.lock.RLock()
	defer tds.lock.RUnlock()
	if todo, ok := tds.dict[key]; ok && todo.Committed {
		return todo.Value, true
	}
	return "", false
}

func (tds *TodoStorage) Commit(key string) error {
	tds.lock.Lock()
	defer tds.lock.Unlock()
	if todo, ok := tds.dict[key]; ok {
		todo.Committed = true
		tds.dict[key] = todo
		return nil
	}
	return ErrorNotFound
}
