package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/kubernetes/scheme"
	"log"

	"os"
	"os/exec"
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

	const chartPath= "../_includes/master/charts/calico"
	f, err := ioutil.TempFile("", "helmfv")
	if err != nil {
		return err
	}
	_, err = f.WriteString(`datastore: kubernetes
network: calico`)
	if err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	cmd := exec.Command("helm", "template", "-f", f.Name(), chartPath)

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

	os.Remove(f.Name())

	var (
		calicoEtcdSecrets v1.Secret
		calicoConfig v1.ConfigMap
		calicoNode v1beta1.DaemonSet
		calicoNodeServiceAccount v1.ServiceAccount
		calicoKubeControllers v1beta1.Deployment
		kubeControllersSA v1.ServiceAccount

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
		//	configureCanal v1.Job
	)

	// there's gotta be a more efficient way to do this but f it
	byteObjs := bytes.Split(stdout.Bytes(), []byte("---"))
	for _, byteObj := range byteObjs {

		obj, gvk, err := scheme.Codecs.UniversalDeserializer().Decode(byteObj, nil, nil)
		if err != nil {
			log.Print(err)
			continue
		}

		fmt.Println(gvk.Kind)
		switch o := obj.(type) {

		case *v1.Secret:
			calicoEtcdSecrets = *o
		case *v1.ConfigMap:
			calicoConfig = *o
		case *v1beta1.DaemonSet:
			calicoNode  = *o
		case *v1.ServiceAccount:
			if o.ObjectMeta.Name == "calico-kube-controllers" {
				kubeControllersSA = *o
			} else if o.ObjectMeta.Name == "calico-node" {
				calicoNodeServiceAccount = *o
			} else {
				return fmt.Errorf("unexpected service account: %s", o.ObjectMeta.Name)
			}
		case *v1beta1.Deployment:
			calicoKubeControllers = *o
		}
		//case "apiextensions.k8s.io/v1beta1/CustomResourceDefinition:bgpconfigurations.crd.projectcalico.org":
		//case "apiextensions.k8s.io/v1beta1/CustomResourceDefinition:ippools.crd.projectcalico.org":
		//case "apiextensions.k8s.io/v1beta1/CustomResourceDefinition:hostendpoints.crd.projectcalico.org":
		//case "apiextensions.k8s.io/v1beta1/CustomResourceDefinition:clusterinformations.crd.projectcalico.org":
		//case "apiextensions.k8s.io/v1beta1/CustomResourceDefinition:globalnetworkpolicies.crd.projectcalico.org":
		//case "apiextensions.k8s.io/v1beta1/CustomResourceDefinition:globalnetworksets.crd.projectcalico.org":
		//case "apiextensions.k8s.io/v1beta1/CustomResourceDefinition:networkpolicies.crd.projectcalico.org":
		//case "rbac.authorization.k8s.io/v1beta1/ClusterRole:calico-kube-controllers":
		//case "rbac.authorization.k8s.io/v1beta1/ClusterRoleBinding:calico-kube-controllers":
		//case "rbac.authorization.k8s.io/v1beta1/ClusterRole:calico-node":
		//case "rbac.authorization.k8s.io/v1beta1/ClusterRoleBinding:calico-node":
		//case "batch/v1/Job:configure-canal":
	}

	fmt.Println(calicoEtcdSecrets.ObjectMeta.Name)
	fmt.Println(calicoConfig.ObjectMeta.Name)
	fmt.Println(calicoNode.ObjectMeta.Name)
	fmt.Println(kubeControllersSA.ObjectMeta.Name)
	fmt.Println(calicoNodeServiceAccount.ObjectMeta.Name)
	fmt.Println(calicoKubeControllers.ObjectMeta.Name)

	return nil
}