package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func main() {
	AWS := AWS{}
	clusters, err := AWS.ListClusters()
	if err != nil {
		fmt.Println("Erro ao listar clusters:", err)
		return
	}
	prompt := promptui.Select{
		Label: "Qual cluster você quer?",
		Items: clusters,
		Size:  10,
	}
	_, cluster, err := prompt.Run()
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}
	services, err := AWS.ListServices(cluster)
	if err != nil {
		fmt.Println("Erro ao listar serviços:", err)
		return
	}
	promptServices := promptui.Select{
		Label: "Qual serviço você quer?",
		Items: services,
		Size:  10,
	}
	_, service, err := promptServices.Run()
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}
	tasks, err := AWS.ListTasks(cluster, service)
	if err != nil {
		fmt.Println("Erro ao listar tarefas:", err)
		return
	}
	task := tasks[0]
	container, err := AWS.GetContainer(cluster, service)
	if err != nil {
		fmt.Println("Erro ao obter container:", err)
		return
	}
	fmt.Println("Seu comando", AWS.GetConnectionCommand(cluster, task, container))
}
