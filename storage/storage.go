package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

type Task struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsDone      bool   `json:"isDone"`
}

type Storage struct {
	Path string
}

type Service interface {
	Add(t Task) error
	List() ([]Task, error)
	Done(id int) error
	Delete(id int) error
}

func NewStorage(path string) *Storage {
	return &Storage{Path: path}
}

func (s *Storage) LoadTasks() ([]Task, error) {
	var tasks []Task
	taskFile, err := os.ReadFile(s.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return tasks, nil
		}
		return nil, err
	}
	err = json.Unmarshal(taskFile, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *Storage) SaveTasks(tasks []Task) error {
	tasksByte, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.Path, tasksByte, 0644)
}

func (s *Storage) Add(t Task) error {
	var maxId int
	tasks, err := s.LoadTasks()
	if err != nil {
		return err
	}

	for _, task := range tasks {
		if task.ID > maxId {
			maxId = task.ID
		}
	}

	t.ID = maxId + 1
	tasks = append(tasks, t)
	err = s.SaveTasks(tasks)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) List() ([]Task, error) {
	tasks, err := s.LoadTasks()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *Storage) Done(id int) error {
	tasks, err := s.LoadTasks()
	if err != nil {
		return fmt.Errorf("Error in loading tasks: %v", err)
	}

	if len(tasks) == 0 {
		return fmt.Errorf("No tasks found")
	}

	found := false
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].IsDone = true
			err = s.SaveTasks(tasks)
			if err != nil {
				return fmt.Errorf("Error saving tasks: %v", err)
			}
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("Task ID not found")
	}
	return nil
}

func (s *Storage) Delete(id int) error {
	tasks, err := s.LoadTasks()
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		return fmt.Errorf("No tasks found")
	}

	found := false
	for i, task := range tasks {
		if task.ID == id {
			found = true
			tasks = append(tasks[:i], tasks[i+1:]...)
			err = s.SaveTasks(tasks)
			if err != nil {
				return fmt.Errorf("Error saving tasks: %v", err)
			}
			break
		}
	}
	if !found {
		return fmt.Errorf("Task ID not found")
	}
	return nil
}
