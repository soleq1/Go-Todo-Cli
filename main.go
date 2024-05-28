package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type todo struct {
	Name     string `json:"name"`
	Msg      string `json:"msg"`
	Complete bool   `json:"complete"`
}

func showMenu() {
	fmt.Println("Todo CLI")
	fmt.Println("1. Add todo")
	fmt.Println("2. List todos")
	fmt.Println("3. Mark todo as completed")
	fmt.Println("4. Delete todo")
	fmt.Println("5. Exit")
}

func handleMenu(todos []todo) {
	fmt.Println("1. Menu || 2. Exit")
	choice := getUserInput("Enter your choice: ")

	switch choice {
	case "1":
		return

	case "2":
		fmt.Println("...Exiting")
		os.Exit(0)
	}
}

func addTodo() todo {
	name := getUserInput("Enter todo name: ")
	msg := getUserInput("Enter todo message: ")
	return todo{Name: name, Msg: msg, Complete: false}
}

func saveTodoList(filename string, todos []todo) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(todos)
}

func loadTodoList(filename string) ([]todo, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var todos []todo
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&todos)
	return todos, err
}

func listTodos(todos []todo) {
	if len(todos) == 0 {
		fmt.Println("Empty List")
		return
	}
	fmt.Println("Todos:")
	for i, t := range todos {
		completeTodo := "☐"
		if t.Complete {
			completeTodo = "☑"
		}
		fmt.Printf("%d. %s Name: %s \n   Message: %s\n", i+1, completeTodo, t.Name, t.Msg)
	}
}

func getUserInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewScanner(os.Stdin)
	if reader.Scan() {
		return strings.TrimSpace(reader.Text())
	}
	return ""
}

func completeTodo(todos []todo) {
	listTodos(todos)
	check := getUserInput("Todo Index :")
	index, err := strconv.Atoi(check)
	if err != nil {
		fmt.Println("invaild Input")
		return
	}
	if index > 0 && index <= len(todos) {
		todos[index-1].Complete = !todos[index-1].Complete
	} else {
		fmt.Println("No  Entry")
	}
}
func deleteTodo(todos []todo) []todo {
	listTodos(todos)
	check := getUserInput("Todo Number")
	index, err := strconv.Atoi(check)
	if err != nil {

		fmt.Println("Invalid Input")
		return todos
	}
	if index > 0 && index <= len(todos) {
		index-- // Adjust for 0-based index
		todos = append(todos[:index], todos[index+1:]...)
	} else {
		fmt.Println("Index Out Of Range")
	}
	return todos
}

func clearMenu() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
func main() {
	// Corrected the file path to reflect the current directory structure
	todoFile := "todo.json"
	todos, err := loadTodoList(todoFile)
	if err != nil {
		fmt.Println("Error loading todo list:", err)
		// Initialize with an empty list if the file doesn't exist or can't be loaded
		todos = []todo{}
	}

	for {
		clearMenu()
		showMenu()
		choice := getUserInput("\nEnter your choice: ")

		switch choice {
		case "1":
			t := addTodo()
			todos = append(todos, t)
			err := saveTodoList(todoFile, todos)
			if err != nil {
				fmt.Println("Error saving todo list:", err)
			}
			handleMenu(todos)
		case "2":
			listTodos(todos)
			handleMenu(todos)
		case "3":
			completeTodo(todos)
			err := saveTodoList(todoFile, todos)
			if err != nil {
				fmt.Println("Error Checking Off", err)
			}
			handleMenu(todos)
			// Implement mark as completed functionality
		case "4":
			todos = deleteTodo(todos)
			err := saveTodoList(todoFile, todos)
			if err != nil {
				fmt.Println("Error saving", err)
			}
			handleMenu(todos)
			// Implement delete functionality
		case "5":
			fmt.Println("Exiting...")
			err := saveTodoList(todoFile, todos)
			if err != nil {
				fmt.Println("Error saving todo list:", err)
			}
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}
