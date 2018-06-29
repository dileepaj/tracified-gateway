package stellarexecuter

import "fmt"
import "main/proofs/proofBuilder"

var currentId int

var todos proofBuilder.Todos

// Give us some seed data
func init() {
	RepoCreateTodo(proofBuilder.Todo{Name: "Write presentation"})
	RepoCreateTodo(proofBuilder.Todo{Name: "Host meetup"})
}

func RepoFindTodo(id int) proofBuilder.Todo {
	for _, t := range todos {
		if t.Id == id {
			return t
		}
	}
	// return empty Todo if not found
	return proofBuilder.Todo{}
}

//this is bad, I don't think it passes race condtions
func RepoCreateTodo(t proofBuilder.Todo) proofBuilder.Todo {
	currentId += 1
	t.Id = currentId
	todos = append(todos, t)
	return t
}

func RepoDestroyTodo(id int) error {
	for i, t := range todos {
		if t.Id == id {
			todos = append(todos[:i], todos[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Todo with id of %d to delete", id)
}
