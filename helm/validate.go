package main

import (
	"bytes"
	"fmt"
	"io"
	"k8s.io/api/core/v1"

	//"k8s.io/api/core/v1"
	//"k8s.io/api/extensions/v1beta1"
	"log"
	"os"
	"os/exec"
	//yaml "gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
//	helmInput := `
//datastore: kubernetes
//ipPool: 192.168.0.0/16`

	const chartPath = "../_includes/master/charts/calico"
	cmd := exec.Command("helm", "template", "--set", "network=calico", "--set", "datastore=etcd", chartPath)

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		fmt.Println(stderr.String())
		return err
	}

	var (
	calicoEtcdSecrets v1.Secret
//	calicoConfig v1.ConfigMap
////bgpconfigurations.crd.projectcalico.org apiextensions.k8s.io/v1beta1/CustomResourceDefinition
////ippools.crd.projectcalico.org apiextensions.k8s.io/v1beta1/CustomResourceDefinition
////hostendpoints.crd.projectcalico.org apiextensions.k8s.io/v1beta1/CustomResourceDefinition
////clusterinformations.crd.projectcalico.org apiextensions.k8s.io/v1beta1/CustomResourceDefinition
////globalnetworkpolicies.crd.projectcalico.org apiextensions.k8s.io/v1beta1/CustomResourceDefinition
////globalnetworksets.crd.projectcalico.org apiextensions.k8s.io/v1beta1/CustomResourceDefinition
////networkpolicies.crd.projectcalico.org apiextensions.k8s.io/v1beta1/CustomResourceDefinition
////calico-kube-controllers rbac.authorization.k8s.io/v1beta1/ClusterRole
////calico-kube-controllers rbac.authorization.k8s.io/v1beta1/ClusterRoleBinding
////calico-node rbac.authorization.k8s.io/v1beta1/ClusterRole
////calico-node rbac.authorization.k8s.io/v1beta1/ClusterRoleBinding
//	calicoNode v1beta1.DaemonSet
//	calicoNodeServiceAccount v1.ServiceAccount
//	calicoKubeControllers v1beta1.Deployment
//	kubeControllersSA v1.ServiceAccount
//	configureCanal v1.Job
	)


	// parse output
	type generic map[string]interface{}

	dec := yaml.NewYAMLOrJSONDecoder(&stdout, 2048)
	//dec := yaml.NewDecoder(&stdout)
	for {
		var emptyI interface{}
		err := dec.Decode(&emptyI)
		if err != nil {
			if err == io.EOF {
				log.Print("finished parsing yaml")
				break
			} else {
				return fmt.Errorf("unexpected error while parsing yaml: %v", err)
			}
		}

		obj := emptyI.(generic)

		if obj["metadata"] == nil {
			log.Print("object has no 'metadata' field. Likely an empty yaml block. Skipping")
			continue
		}
		metaData := obj["metadata"].(map[string]interface{})
		identifier := fmt.Sprintf("%s/%s:%s", obj["apiVersion"], obj["kind"], metaData["name"])
		switch identifier {
		case "v1/Secret:calico-etcd-secrets":
			fmt.Println("found it")
			calicoEtcdSecrets = emptyI.(v1.Secret)
		case "v1/ConfigMap:calico-config":
		case "apiextensions.k8s.io/v1beta1/CustomResourceDefinition:bgpconfigurations.crd.projectcalico.org":
		case "apiextensions.k8s.io/v1beta1/CustomResourceDefinition:ippools.crd.projectcalico.org":
		case "apiextensions.k8s.io/v1beta1/CustomResourceDefinition:hostendpoints.crd.projectcalico.org":
		case "apiextensions.k8s.io/v1beta1/CustomResourceDefinition:clusterinformations.crd.projectcalico.org":
		case "apiextensions.k8s.io/v1beta1/CustomResourceDefinition:globalnetworkpolicies.crd.projectcalico.org":
		case "apiextensions.k8s.io/v1beta1/CustomResourceDefinition:globalnetworksets.crd.projectcalico.org":
		case "apiextensions.k8s.io/v1beta1/CustomResourceDefinition:networkpolicies.crd.projectcalico.org":
		case "rbac.authorization.k8s.io/v1beta1/ClusterRole:calico-kube-controllers":
		case "rbac.authorization.k8s.io/v1beta1/ClusterRoleBinding:calico-kube-controllers":
		case "rbac.authorization.k8s.io/v1beta1/ClusterRole:calico-node":
		case "rbac.authorization.k8s.io/v1beta1/ClusterRoleBinding:calico-node":
		case "extensions/v1beta1/DaemonSet:calico-node":
		case "v1/ServiceAccount:calico-node":
		case "extensions/v1beta1/Deployment:calico-kube-controllers":
		case "v1/ServiceAccount:calico-kube-controllers":
		case "batch/v1/Job:configure-canal":
		}
	}

	fmt.Println(calicoEtcdSecrets.ObjectMeta.Name)
	return nil
}