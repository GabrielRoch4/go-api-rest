package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var slicePessoa = []Pessoa{}
var contId = 0

type Pessoa struct {
	Id   int
	Nome string
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

	for _, value := range slicePessoa {
		// fmt.Fprintf(w, "ID: %d, Nome: %s\n", slicePessoa[index].Id, slicePessoa[index].Nome)
		pessoaJson, err := json.Marshal(value)

		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Fprintln(w, string(pessoaJson))
	}
}

func remover(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	var removido Pessoa

	novaSlicePessoa := []Pessoa{}
	for _, pessoa := range slicePessoa {
		if pessoa.Id != id {
			novaSlicePessoa = append(novaSlicePessoa, pessoa)
		} else {
			removido = pessoa
		}
	}

	if removido.Id == 0 && removido.Nome == "" {
		// fmt.Fprintf(w, "O ID: %d não está cadastrado!", id)
		http.Error(w, "ID informado não está cadastrado!", http.StatusNotFound)
		return
	} else {
		fmt.Fprintf(w, "ID: %d foi removido!", removido.Id)
	}

	slicePessoa = novaSlicePessoa
}

func exibirPorID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	var encontrado bool

	for _, value := range slicePessoa {
		if value.Id != id {
			encontrado = false
		} else {
			encontrado = true

			pessoaJson, err := json.Marshal(value)

			if err != nil {
				log.Fatal(err.Error())
			}

			fmt.Fprintln(w, string(pessoaJson))

			return
		}
	}

	if encontrado == false {
		http.Error(w, "ID informado não está cadastrado!", http.StatusNotFound)
		return
	}
}

func salvar(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	nome := r.URL.Query().Get("nome")

	contId++

	pessoa := Pessoa{
		Id:   contId,
		Nome: nome,
	}

	slicePessoa = append(slicePessoa, pessoa)

	http.Error(w, "Cadastrado com sucesso!", http.StatusCreated)

	return
	// fmt.Fprintf(w, "%s, Cadastrado com sucesso", pessoa.Nome)
}

func main() {
	http.HandleFunc("/pessoa", pessoaHandler)
	http.HandleFunc("/pessoa/", pessoaHandlerDeleteAndGetById)

	_ = http.ListenAndServe(":3333", nil)
}
