package main_test

import (
	"testing"
)

// This file holds all of the functions for validating
// expected values on the Typha deployment.

func TestTypha(t *testing.T) {
	t.Run("typha enabled with kubernetes datastore", func(t *testing.T) {
		testYml := `datastore: kdd
typha:
  enabled: true`

		resources, err := render(testYml)
		if err != nil {
			t.Error(err, "Error while rendering resources with typha enabled and datastore set to a kubernetes datastore")
		}

		// Validate that the typha resources exist
		if resources.calicoTypha == nil {
			t.Error(err, "calico-typha deployment should exist when typha is enabled with a Kubernetes datastore")
		}

		if resources.calicoTyphaService == nil {
			t.Error(err, "calico-typha service should exist when typha is enabled with a Kubernetes datastore")
		}

		if resources.calicoTyphaPDB == nil {
			t.Error(err, "calico-typha pod disruption budget should exist when typha is enabled with a Kubernetes datastore")
		}
	})

	t.Run("typha disabled with kubernetes datastore", func(t *testing.T) {
		testYml := `datastore: kdd
typha:
  enabled: false`

		resources, err := render(testYml)
		if err != nil {
			t.Error(err, "Error while rendering resources with typha disabled and datastore set to a kubernetes datastore")
		}

		// Validate that the typha resources exist
		if resources.calicoTypha != nil {
			t.Error(err, "calico-typha deployment should not exist when typha is disabled with a Kubernetes datastore")
		}

		if resources.calicoTyphaService != nil {
			t.Error(err, "calico-typha service should not exist when typha is disabled with a Kubernetes datastore")
		}

		if resources.calicoTyphaPDB != nil {
			t.Error(err, "calico-typha pod disruption budget should not exist when typha is disabled with a Kubernetes datastore")
		}
	})

	t.Run("typha enabled with etcd datastore", func(t *testing.T) {
		testYml := `datastore: etcd
etcd:
  endpoints: foo
typha:
  enabled: true`

		resources, err := render(testYml)
		if err != nil {
			t.Error(err, "Error while rendering resources with typha enabled and datastore set to etcd")
		}

		// Validate that the typha resources exist
		if resources.calicoTypha != nil {
			t.Error(err, "calico-typha deployment should not exist when typha is enabled with an etcd datastore")
		}

		if resources.calicoTyphaService != nil {
			t.Error(err, "calico-typha service should not exist when typha is enabled with an etcd datastore")
		}

		if resources.calicoTyphaPDB != nil {
			t.Error(err, "calico-typha pod disruption budget should not exist when typha is enabled with an etcd datastore")
		}
	})

	t.Run("typha disabled with etcd datastore", func(t *testing.T) {
		testYml := `datastore: etcd
etcd:
  endpoints: foo
typha:
  enabled: false`

		resources, err := render(testYml)
		if err != nil {
			t.Error(err, "Error while rendering resources with typha disabled and datastore set to etcd")
		}

		// Validate that the typha resources exist
		if resources.calicoTypha != nil {
			t.Error(err, "calico-typha deployment should not exist when typha is disabled with an etcd datastore")
		}

		if resources.calicoTyphaService != nil {
			t.Error(err, "calico-typha service should not exist when typha is disabled with an etcd datastore")
		}

		if resources.calicoTyphaPDB != nil {
			t.Error(err, "calico-typha pod disruption budget should not exist when typha is disabled with an etcd datastore")
		}
	})
}
