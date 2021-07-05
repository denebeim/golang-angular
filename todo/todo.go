package todo

import (
	"errors"
	"sync"

	"github.com/rs/xid"
)

var (
	list []Todo
	mtx  sync.RWMutex
	once sync.Once
)

func init() {
	once.Do(initialiseList)
}

func initialiseList() {
	list = []Todo{}
}

// Todo data structure for a task with a description of what to do
type Todo struct {
	ID       string `json:"id"`
	Message  string `json:"message"`
	Complete bool   `json:"complete"`
}

func Get() []Todo {
	return list
}

func newTodo(msg string) Todo {
	return Todo{
		ID:       xid.New().String(),
		Message:  msg,
		Complete: false,
	}
}

// Add will add a new todo based on a message
func Add(message string) string {
	t := newTodo(message)
	mtx.Lock()
	list = append(list, t)
	mtx.Unlock()
	return t.ID
}

// Delete will remove a Todo from the Todo list
func Delete(id string) error {
	location, err := findTodoLocation(id)
	if err != nil {
		return err
	}
	removeElementByLocation(location)
	return nil
}

func Complete(id string) error {
	location, err := findTodoLocation(id)
	if err != nil {
		return err
	}
	setTodoCompleteByLocation(location)
	return nil
}

func findTodoLocation(id string) (int, error) {
	mtx.RLock()
	defer mtx.Unlock()
	for i, t := range list {
		if isMatchingID(t.ID, id) {
			return i, nil
		}
	}
	return 0, errors.New("could not find todo based on id")
}

func removeElementByLocation(i int) {
	mtx.Lock()
	list = append(list[:i], list[i+1:]...)
	mtx.Unlock()
}

func setTodoCompleteByLocation(location int) {
	mtx.Lock()
	list[location].Complete = true
	mtx.Unlock()
}

func isMatchingID(a string, b string) bool {
	return a == b
}
