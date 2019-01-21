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
	var containerDefinitions []ContainerDefinitions
	var requiresCompatibilities []string
	requiresCompatibilities = append(requiresCompatibilities, "FARGATE")
	jsonFile, err := os.Open(os.Args[1])
	executionRoleArn := os.Args[3]
	image := strings.Replace(os.Args[4], "\"", "", -1)
	name := os.Args[5]
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}

	json.Unmarshal(byteValue, &result)

	defer jsonFile.Close()

	for k, v := range result {
		envs = append(envs, Environment{Name: k, Value: v.(string)})
	}
	var portmappings []Portmappings
	portmappings = append(portmappings, Portmappings{
		HostPort:      80,
		Protocol:      "tcp",
		ContainerPort: 80,
	})

	containerDefinitions = append(containerDefinitions, ContainerDefinitions{
		Name:  name,
		Image: image,
		//Image:        "<IMAGE1_NAME>",
		Essential:    true,
		PortMappings: portmappings,
		Environment:  envs,
	})

	taskdef := Taskdef{

		ExecutionRoleArn:        executionRoleArn,
		NetworkMode:             "awsvpc",
		CPU:                     "256",
		Memory:                  "512",
		Family:                  "ecs-demo",
		RequiresCompatibilities: requiresCompatibilities,
		ContainerDefinitions:    containerDefinitions,
	}
	json, _ := JSONMarshal(taskdef)
	//json, _ := json.Marshal(taskdef)
	err = ioutil.WriteFile(os.Args[2], json, 0644)
	if err != nil {
		panic(err)
	}

}

// ContainerDefinitions ...
type ContainerDefinitions struct {
	Name         string         `json:"name"`
	Image        string         `json:"image"`
	Essential    bool           `json:"essential"`
	PortMappings []Portmappings `json:"portMappings"`
	Environment  []Environment  `json:"environment"`
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
