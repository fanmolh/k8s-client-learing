package main

import (
	"fmt"
	"runtime"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, _ := clientcmd.BuildConfigFromFlags("", "../kbconfig")
	client, _ := kubernetes.NewForConfig(config)
	stopCh := make(chan struct{})
	defer close(stopCh)
	informerFactory := informers.NewSharedInformerFactory(client, time.Second*30)
	podInformer := informerFactory.Core().V1().Pods().Informer()
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			fmt.Printf("Pod %s/%s added\n", pod.Namespace, pod.Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldpod := oldObj.(*v1.Pod)
			newpod := newObj.(*v1.Pod)
			fmt.Printf("Pod update ->:  %s->%s\n", oldpod.Name, newpod.Name)
			fmt.Printf("Pod Message ->:  %s", newpod.Status.Message)
		},
	})
	informerFactory.Start(stopCh)
	if !cache.WaitForCacheSync(stopCh, podInformer.HasSynced) {
		runtime.Goexit()
	}
	select {}

}
