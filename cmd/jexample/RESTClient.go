package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	as1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string

	//Home is the home directory. If you can get the value of home directory, it can be used as the default value
	if home := homedir.HomeDir(); home != "" {
		//If the kubeconfig parameter is entered, the value of this parameter is the absolute path of the kubeconfig file,
		//If the kubeconfig parameter is not entered, the default path ~ / kube/config
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		//If the home directory of the current user cannot be obtained, there is no way to set the default directory of kubeconfig. You can only get it from the input parameter
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()

	//The kubeconfig configuration file is loaded natively, so the first parameter is an empty string
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	//Kubeconfig failed to load and exited directly
	if err != nil {
		panic(err.Error())
	}

	//Reference path: / API / V1 / namespaces / {namespace} / pods
	config.APIPath = "api"
	//The group of pod is an empty string
	config.GroupVersion = &corev1.SchemeGroupVersion
	//Specify serialization tool
	config.NegotiatedSerializer = scheme.Codecs

	//Build the restclient instance according to the configuration information
	restClient, err := rest.RESTClientFor(config)

	if err != nil {
		panic(err.Error())
	}

	//Save the data structure instance of pod results
	result := &corev1.PodList{}

	//Specify namespace
	// namespace := "kube-system"
	namespace := "default"
	//Set the request parameters and then initiate the request
	//Get request
	err = restClient.Get().
		//Specify namespace，参考path : /api/v1/namespaces/{namespace}/pods
		Namespace(namespace).
		//Find multiple pods. Refer to path: / API / V1 / namespaces / {namespace} / pods
		Resource("pods").
		//Specify size limits and serialization tools
		VersionedParams(&metav1.ListOptions{Limit: 100}, scheme.ParameterCodec).
		//Request
		Do(context.TODO()).
		//Save the result into result
		Into(result)

	if err != nil {
		panic(err.Error())
	}

	//Header
	fmt.Printf("namespace\t status\t\t name\n")

	//Each pod prints namespace and status Phase and name fields
	for _, d := range result.Items {
		fmt.Printf("%v\t %v\t %v\n",
			d.Namespace,
			d.Status.Phase,
			d.Name)
	}

	// hpa test
	// https://blog.csdn.net/cuiwangfeng/article/details/119796915
	// TARGETS <unknown>/50%

	// /api/v1/namespaces/{namespace}/pods/{name}
	// /apis/autoscaling/v1/namespaces/{namespace}/horizontalpodautoscalers/{name}
	config.APIPath = "apis"
	//The group of pod is an empty string
	config.GroupVersion = &as1.SchemeGroupVersion
	//Specify serialization tool
	config.NegotiatedSerializer = scheme.Codecs

	//Build the restclient instance according to the configuration information
	restClient, err = rest.RESTClientFor(config)

	if err != nil {
		panic(err.Error())
	}
	hpa := &as1.HorizontalPodAutoscalerList{}

	err = restClient.Get().
		Namespace(namespace).
		Resource("horizontalpodautoscalers").
		//Specify size limits and serialization tools
		VersionedParams(&metav1.ListOptions{}, scheme.ParameterCodec).
		//Request
		Do(context.TODO()).
		//Save the result into result
		Into(hpa)
	if err != nil {
		fmt.Println(err)
	}

	for _, d := range hpa.Items {
		fmt.Printf("%+v\n", d)
		max := d.Spec.MaxReplicas
		min := d.Spec.MinReplicas
		percent := d.Spec.TargetCPUUtilizationPercentage
		ns := d.Namespace
		name := d.Name

		fmt.Println(ns, name, *min, max, *percent)
	}
}
