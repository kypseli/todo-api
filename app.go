package main

import (
	"encoding/json"
	"log"
	"fmt"
	"os"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/handlers"
    "github.com/gorilla/mux"
	"github.com/kypseli/todo-api/dao"
	"github.com/kypseli/todo-api/models"
)

var todosDao = dao.TodosDAO{}

//app health check
func HealthCheckEndpoint(w http.ResponseWriter, r *http.Request) {
    // A very simple health check.
    // In the future report back on the status of the DB
    respondWithJson(w, http.StatusOK,map[string]bool{"alive": true})
}

// GET list of todos
func AllTodosEndPoint(w http.ResponseWriter, r *http.Request) {
	todos, err := todosDao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if todos == nil {
		respondWithJson(w, http.StatusOK, map[string]string{})
	}
	respondWithJson(w, http.StatusOK, todos)
}

// GET a todo by its ID
func FindTodoEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	todo, err := todosDao.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Todo ID")
		return
    }
	respondWithJson(w, http.StatusOK, todo)
}

// POST a new todo
func CreateTodoEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	todo.ID = bson.NewObjectId()
	if err := todosDao.Insert(todo); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
    }
	respondWithJson(w, http.StatusCreated, todo)
}

// PUT update an existing todo
func UpdateTodoEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := todosDao.Update(todo); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
    }
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// DELETE an existing todo
func DeleteTodoEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := todosDao.Delete(todo); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
    }
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// establish a connection to DB
func init() {
	mongoHost := getEnv("MONGO_HOST", "localhost")
	mongoPort := getEnv("MONGO_PORT", "27017")
	todosDao.Server = fmt.Sprintf("%s:%s", mongoHost, mongoPort)
	todosDao.Database = "todo_db"
	todosDao.Connect()
}

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

// Define HTTP request routes
func main() {
    r := mux.NewRouter()
    r.HandleFunc("/health", HealthCheckEndpoint).Methods("GET")
	r.HandleFunc("/todos", AllTodosEndPoint).Methods("GET","OPTIONS")
	r.HandleFunc("/todos", CreateTodoEndPoint).Methods("POST","OPTIONS")
	r.HandleFunc("/todos", UpdateTodoEndPoint).Methods("PUT","OPTIONS")
	r.HandleFunc("/todos", DeleteTodoEndPoint).Methods("DELETE","OPTIONS")
	r.HandleFunc("/todos/{id}", FindTodoEndpoint).Methods("GET","OPTIONS")

    headersOk := handlers.AllowedHeaders([]string{"Content-Type", "X-Requested-With"})
    originsOk := handlers.AllowedOrigins([]string{"*"})
    methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	if err := http.ListenAndServe(":3000", handlers.CORS(headersOk, originsOk, methodsOk)(r)); err != nil {
		log.Fatal(err)
	}
}