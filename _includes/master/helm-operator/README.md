kubectl create clusterrolebinding operator-admin \
  --clusterrole=tigera-manager-user --serviceaccount=calico-operator