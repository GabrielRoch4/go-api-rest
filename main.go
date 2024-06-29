package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Pessoa struct {
	Id   int    `json:"id"`
	Nome string `json:"nome"`
}

func pessoaHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		exibir(w, r)
	case http.MethodPost:
		salvar(w, r)
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

func pessoaHandlerDeleteAndGetById(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		remover(w, r)
	case http.MethodGet:
		exibirPorID(w, r)
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

func exibir(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query("SELECT id, nome FROM pessoas")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pessoas []Pessoa
	for rows.Next() {
		var p Pessoa
		if err := rows.Scan(&p.Id, &p.Nome); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pessoas = append(pessoas, p)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pessoas)
}

func exibirPorID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var p Pessoa
	err = db.QueryRow("SELECT id, nome FROM pessoas WHERE id = ?", id).Scan(&p.Id, &p.Nome)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "ID informado não está cadastrado", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func salvar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	nome := r.URL.Query().Get("nome")
	if nome == "" {
		http.Error(w, "Nome é obrigatório", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO pessoas (nome) VALUES (?)", nome)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Pessoa{Id: int(id), Nome: nome})
}

func remover(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("DELETE FROM pessoas WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "ID informado não está cadastrado", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ID: %d foi removido!", id)
}

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/go-api-rest"
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erro ao abrir a conexão com o banco de dados:", err)
	}
	defer db.Close()

	// Verificar a conexão
	err = db.Ping()
	if err != nil {
		log.Fatal("Erro ao verificar a conexão com o banco de dados:", err)
	}

	fmt.Println("Conexão com o banco de dados estabelecida com sucesso!")

	http.HandleFunc("/pessoa", pessoaHandler)
	http.HandleFunc("/pessoa/", pessoaHandlerDeleteAndGetById)

	log.Fatal(http.ListenAndServe(":3333", nil))
}
