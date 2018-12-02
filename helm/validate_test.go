package main_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/kubernetes/scheme"
	"log"
	"testing"

	"os"
	"os/exec"
)

func TestIt(t *testing.T) {
	r, err := run(`datastore: kubernetes
network: calico`)
	if err != nil {
		t.Error(err)
	}

	if r.calicoNode.Name != "calico-node" {
		t.Fail()
	}
}

type AllResources struct{
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
}

func run(valuesyml string) (AllResources, error) {
	const chartPath = "../_includes/master/charts/calico"
	f, err := ioutil.TempFile("", "helmfv")
	if err != nil {
		return AllResources{}, err
	}
	_, err = f.WriteString(valuesyml)
	if err != nil {
		return AllResources{}, err
	}
	if err := f.Close(); err != nil {
		return AllResources{}, err
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
		return AllResources{}, err
	}

	os.Remove(f.Name())

	var r AllResources


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
			r.calicoEtcdSecrets = *o
		case *v1.ConfigMap:
			r.calicoConfig = *o
		case *v1beta1.DaemonSet:
			r.calicoNode  = *o
		case *v1.ServiceAccount:
			if o.ObjectMeta.Name == "calico-kube-controllers" {
				r.kubeControllersSA = *o
			} else if o.ObjectMeta.Name == "calico-node" {
				r.calicoNodeServiceAccount = *o
			} else {
				return r, fmt.Errorf("unexpected service account: %s", o.ObjectMeta.Name)
			}
		case *v1beta1.Deployment:
			r.calicoKubeControllers = *o
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

	return r, nil
}