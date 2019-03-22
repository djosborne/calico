package main_test

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	appv1beta1 "k8s.io/api/apps/v1beta1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	polv1beta1 "k8s.io/api/policy/v1beta1"
	"k8s.io/client-go/kubernetes/scheme"
)

type AllResources struct {
	calicoEtcdSecrets        *v1.Secret
	calicoConfig             *v1.ConfigMap
	calicoNode               *v1beta1.DaemonSet
	calicoNodeServiceAccount *v1.ServiceAccount
	calicoKubeControllers    *v1beta1.Deployment
	calicoTypha              *appv1beta1.Deployment
	kubeControllersSA        *v1.ServiceAccount
	calicoTyphaService       *v1.Service
	calicoTyphaPDB           *polv1beta1.PodDisruptionBudget
}

var (
	chartPath string
)

func init() {
	flag.StringVar(&chartPath, "chart-path", "../_includes/master/charts/calico", "path to the chart")
}

// TODO: Add call to kubeval to verify helm resources are valid
func render(valuesyml string) (AllResources, error) {
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

	cmd := exec.Command("helm", "template", "-f", f.Name())
	cmd.Args = append(cmd.Args, chartPath)

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err = cmd.Run(); err != nil {
		fmt.Println(stderr.String())
		return AllResources{}, fmt.Errorf("error running helm: %s", stderr.String())
	}

	os.Remove(f.Name())

	var r AllResources

	// TODO: there's gotta be a more efficient way to do this
	byteObjs := bytes.Split(stdout.Bytes(), []byte("---"))
	for _, byteObj := range byteObjs {

		obj, gvk, err := scheme.Codecs.UniversalDeserializer().Decode(byteObj, nil, nil)
		if err != nil {
			// TODO: Silence the ugly errors for unrendered manifests
			log.Print(err)
			continue
		}

		fmt.Println(gvk.Kind)
		switch o := obj.(type) {

		case *v1.Secret:
			r.calicoEtcdSecrets = o
		case *v1.ConfigMap:
			r.calicoConfig = o
		case *v1beta1.DaemonSet:
			r.calicoNode = o
		case *v1.ServiceAccount:
			switch o.ObjectMeta.Name {
			case "calico-kube-controllers":
				r.kubeControllersSA = o
			case "calico-node":
				r.calicoNodeServiceAccount = o
			}
		case *v1beta1.Deployment:
			r.calicoKubeControllers = o
		case *appv1beta1.Deployment:
			r.calicoTypha = o
		case *v1.Service:
			r.calicoTyphaService = o
		case *polv1beta1.PodDisruptionBudget:
			r.calicoTyphaPDB = o
		}
	}

	return r, nil
}
