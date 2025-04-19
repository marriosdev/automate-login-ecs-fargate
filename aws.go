package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type AWS struct {
}

func (a *AWS) ListClusters() ([]string, error) {
	clusters, error := a.command("ecs", "list-clusters", "--query", "clusterArns[*]", "--output", "json")
	if error != nil {
		return nil, fmt.Errorf("erro ao executar comando: %w", error)
	}
	var clustersList []string
	if err := json.Unmarshal(clusters, &clustersList); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse do JSON: %w", err)
	}
	for i, cluster := range clustersList {
		clustersList[i] = strings.Split(cluster, "/")[1]
	}
	return clustersList, nil
}

func (a *AWS) ListServices(cluster string) ([]string, error) {
	services, err := a.command("ecs", "list-services", "--cluster", cluster, "--query", "serviceArns[*]", "--output", "json")
	if err != nil {
		return nil, fmt.Errorf("erro ao executar comando: %w", err)
	}
	var servicesList []string
	if err := json.Unmarshal(services, &servicesList); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse do JSON: %w", err)
	}
	for i, service := range servicesList {
		serviceSplit := strings.Split(service, "/")
		servicesList[i] = serviceSplit[len(serviceSplit)-1]
	}
	return servicesList, nil
}

func (a *AWS) ListTasks(cluster, service string) ([]string, error) {
	tasks, err := a.command("ecs", "list-tasks", "--cluster", cluster, "--service-name", service, "--query", "taskArns[*]", "--output", "json")
	if err != nil {
		return nil, fmt.Errorf("erro ao executar comando: %w", err)
	}
	var tasksList []string
	if err := json.Unmarshal(tasks, &tasksList); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse do JSON: %w", err)
	}
	for i, task := range tasksList {
		taskSplit := strings.Split(task, "/")
		tasksList[i] = taskSplit[len(taskSplit)-1]
	}
	return tasksList, nil
}

func (a *AWS) GetContainer(cluster string, service string) (string, error) {
	container, err := a.command("ecs", "describe-services",
		"--cluster", cluster,
		"--services", service,
		"--query", "services[0].loadBalancers[0].containerName",
		"--output", "json",
	)
	if err != nil {
		return "", fmt.Errorf("erro ao executar comando: %w", err)
	}
	return string(container), nil
}

func (a *AWS) GetConnectionCommand(cluster string, task string, container string) string {
	return fmt.Sprintf("aws ecs execute-command --cluster %s --task %s --container %s --interactive --command '/bin/sh'", cluster, task, container)
}

func (a *AWS) command(args ...string) ([]byte, error) {
	cmd := exec.Command("aws", args...)
	output, err := cmd.Output()
	if err != nil {
		return []byte{}, err
	}
	return output, nil
}
