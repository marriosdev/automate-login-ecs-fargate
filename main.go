package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
)

func main() {
	cmd := exec.Command("aws", "ecs", "list-clusters", "--query", "clusterArns[*]", "--output", "json")

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Erro ao executar comando:", err)
		return
	}

	var clusters []string
	if err := json.Unmarshal(output, &clusters); err != nil {
		fmt.Println("Erro ao fazer parse do JSON:", err)
		return
	}

	for i, cluster := range clusters {
		clusters[i] = strings.Split(cluster, "/")[1]
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

	fmt.Printf("Você escolheu %q\n", cluster)

	cmd = exec.Command("aws", "ecs", "list-services", "--cluster", cluster, "--query", "serviceArns[*]", "--output", "json")
	output, err = cmd.Output()

	var services []string
	if err := json.Unmarshal(output, &services); err != nil {
		fmt.Println("Erro ao fazer parse do JSON:", err)
		return
	}

	promptServices := promptui.Select{
		Label: "Qual serviço você quer?",
		Items: services,
		Size:  10,
	}
	_, service_selected, err := promptServices.Run()

	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	service_selected_split := strings.Split(service_selected, "/")
	service_selected = service_selected_split[len(service_selected_split)-1]

	fmt.Printf("Você escolheu %q\n", service_selected)

	cmd = exec.Command(
		"aws", "ecs", "list-tasks",
		"--cluster", cluster,
		"--service-name", service_selected,
		"--query", "taskArns[0]",
		"--output", "json",
	)

	task_id, err := cmd.Output()

	if err != nil {
		fmt.Println("Erro ao executar comando:", err)
		return
	}

	cmd = exec.Command(
		"aws", "ecs", "describe-services",
		"--cluster", cluster,
		"--services", service_selected,
		"--query", "services[0].loadBalancers[0].containerName",
		"--output", "json",
	)

	container, err := cmd.Output()
	fmt.Print(string(container))

	if err != nil {
		fmt.Println("Erro ao executar comando:", err)
		return
	}

	cmd = exec.Command(
		"aws", "ecs", "execute-command",
		"--cluster", cluster,
		"--task", string(task_id),
		"--container", string(container),
		"--command", "'/bin/bash'",
		"--interactive",
	)

	fmt.Println("Executando comando: ", cmd.String())
	container_exec, err := cmd.Output()

	if err != nil {
		fmt.Println("Erro ao executar comando:", err)
		return
	}

	fmt.Print(string(container_exec))
}
