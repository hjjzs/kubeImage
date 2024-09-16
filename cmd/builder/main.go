package main

import (
	"context"
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	clientset "test/pkg/generated/clientset/versioned"
)

func main() {
	cfg, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	client, err := clientset.NewForConfig(cfg)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	//factory := externalversions.NewSharedInformerFactory(client, time.Second*30)
	//informer := factory.Hjjzs().V1().DockerFiles()
	list, err := client.HjjzsV1().DockerFiles("default").List(context.Background(), v1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range list.Items {
		fmt.Println(item.Name)
	}

}
