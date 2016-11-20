/*
Copyright 2016 Michael Morello

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"time"

	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/pkg/api"
	"k8s.io/client-go/1.4/tools/clientcmd"
)

var (
	kubeconfig = flag.String("kubeconfig", "./config", "absolute path to the kubeconfig file")
)

func main() {
	flag.Parse()
	// uses the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		nodes, err := clientset.Core().Nodes().List(api.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		now := time.Now()
		tenSecondsAgo := now.Add(-10 * time.Second).Unix()
		for _, element := range nodes.Items {
			conditions := element.Status.Conditions
			for _, condition := range conditions {
				if condition.LastHeartbeatTime.Unix() < tenSecondsAgo {
					fmt.Printf("state=KO|node=%s|delay=%d|now=%v|lastHB=%v\n", element.Name, tenSecondsAgo-condition.LastHeartbeatTime.Unix(), now, condition.LastHeartbeatTime)
					break
				}
			}
		}
		time.Sleep(10 * time.Second)
	}
}
