package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}
func main() {

	var envs []Environment
	var requiresCompatibilities []string
	requiresCompatibilities = append(requiresCompatibilities, "FARGATE")
	jsonFile := os.Getenv("path")
	image := strings.Replace(os.Getenv("imageID"), "\"", "", -1)

	taskOldFile := os.Getenv("taskold")

	var oldTaskDefinition Conf

	json.Unmarshal([]byte(taskOldFile), &oldTaskDefinition)

	//byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}

	json.Unmarshal([]byte(jsonFile), &result)

	for k, v := range result {
		envs = append(envs, Environment{Name: k, Value: v.(string)})
	}

	oldTaskDefinition.TaskDefinition.ContainerDefinitions[0].Environment = envs
	oldTaskDefinition.TaskDefinition.ContainerDefinitions[0].Image = image
	json, _ := JSONMarshal(oldTaskDefinition.TaskDefinition)
	//json, _ := json.Marshal(taskdef)
	err := ioutil.WriteFile("teste.json", json, 0644)
	if err != nil {
		panic(err)
	}

}

// Conf ...
type Conf struct {
	TaskDefinition Taskdef `json:"taskDefinition"`
}

// ContainerDefinitions ...
type ContainerDefinitions struct {
	Name         string         `json:"name"`
	Image        string         `json:"image"`
	Essential    bool           `json:"essential"`
	PortMappings []Portmappings `json:"portMappings"`
	Environment  []Environment  `json:"environment"`
	CPU          int            `json:"cpu"`
}

// Portmappings ...
type Portmappings struct {
	HostPort      int    `json:"hostPort"`
	Protocol      string `json:"protocol"`
	ContainerPort int    `json:"containerPort"`
}

// Taskdef ...
type Taskdef struct {
	ExecutionRoleArn        string                 `json:"executionRoleArn"`
	ContainerDefinitions    []ContainerDefinitions `json:"containerDefinitions"`
	RequiresCompatibilities []string               `json:"requiresCompatibilities"`
	NetworkMode             string                 `json:"networkMode"`
	CPU                     string                 `json:"cpu"`
	Memory                  string                 `json:"memory"`
	Family                  string                 `json:"family"`
}

// Environment ...
type Environment struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
