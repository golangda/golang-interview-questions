package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ResourceSpec struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

type FlowConfiguration struct {
	v1.TypeMeta   `json:",inline"`
	v1.ObjectMeta `json:"metadata,omitempty"`
	Spec          struct {
		Sources      []string     `json:"sources"`
		Destinations []string     `json:"destinations"`
		Resources    ResourceSpec `json:"resources"`
	} `json:"spec"`
}

var dynamicClient dynamic.Interface
var namespace = "default"

func initK8sClient() {
	config, err := rest.InClusterConfig()
	if err != nil {
		kubeconfig := os.Getenv("KUBECONFIG")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatalf("Error creating kubernetes client config: %v", err)
		}
	}

	dynamicClient, err = dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating dynamic client: %v", err)
	}
}

func createFlowConfiguration(w http.ResponseWriter, r *http.Request) {
	var fc FlowConfiguration
	err := json.NewDecoder(r.Body).Decode(&fc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	gvr := schema.GroupVersionResource{
		Group:    "example.com",
		Version:  "v1",
		Resource: "flowconfigurations",
	}

	unstructuredFC, err := dynamicClient.Resource(gvr).Namespace(namespace).Create(
		context.TODO(),
		&fc,
		v1.CreateOptions{},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(unstructuredFC)
}

func updateFlowConfiguration(w http.ResponseWriter, r *http.Request) {
	var fc FlowConfiguration
	err := json.NewDecoder(r.Body).Decode(&fc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// name := fc.ObjectMeta.Name
	gvr := schema.GroupVersionResource{
		Group:    "example.com",
		Version:  "v1",
		Resource: "flowconfigurations",
	}

	unstructuredFC, err := dynamicClient.Resource(gvr).Namespace(namespace).Update(
		context.TODO(),
		&fc,
		v1.UpdateOptions{},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(unstructuredFC)
}

func deleteFlowConfiguration(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	gvr := schema.GroupVersionResource{
		Group:    "example.com",
		Version:  "v1",
		Resource: "flowconfigurations",
	}

	err := dynamicClient.Resource(gvr).Namespace(namespace).Delete(
		context.TODO(),
		name,
		v1.DeleteOptions{},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	initK8sClient()

	http.HandleFunc("/create", createFlowConfiguration)
	http.HandleFunc("/update", updateFlowConfiguration)
	http.HandleFunc("/delete", deleteFlowConfiguration)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
