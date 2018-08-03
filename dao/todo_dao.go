package dao

import (
	"log"

	. "github.com/kypseli/todo-api/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type TodosDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "todos"
)

// Establish a connection to database
func (m *TodosDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

// Find list of todos
func (m *TodosDAO) FindAll() ([]Todo, error) {
	var todos []Todo
	err := db.C(COLLECTION).Find(bson.M{}).All(&todos)
	return todos, err
}

// Find a todo by its id
func (m *TodosDAO) FindById(id string) (Todo, error) {
	var todo Todo
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&todo)
	return todo, err
}

// Insert a todo into database
func (m *TodosDAO) Insert(todo Todo) error {
	err := db.C(COLLECTION).Insert(&todo)
	return err
}

// Delete an existing todo
func (m *TodosDAO) Delete(todo Todo) error {
	err := db.C(COLLECTION).Remove(&todo)
	return err
}

// Update an existing todo
func (m *TodosDAO) Update(todo Todo) error {
	err := db.C(COLLECTION).UpdateId(todo.ID, &todo)
	return err
}