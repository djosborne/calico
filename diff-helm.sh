#!/bin/bash
rm -rf _old
git checkout 8eb2bdd016a769667f8f5ed662f061a41fc33cad
jekyll build  --config _config.yml,_config_dev.yml
mv _site _old

git checkout helm-convert-templates
jekyll build  --config _config.yml,_config_dev.yml

diff -y _site/master/getting-started/kubernetes/installation/hosted/calico.yaml _old/master/getting-started/kubernetes/installation/hosted/calico.yaml
diff _site/master/getting-started/kubernetes/installation/hosted/canal/canal-etcd.yaml _old/master/getting-started/kubernetes/installation/hosted/canal/canal-etcd.yaml
diff _site/master/getting-started/kubernetes/installation/hosted/canal/canal.yaml _old/master/getting-started/kubernetes/installation/hosted/canal/canal.yaml
diff _site/master/getting-started/kubernetes/installation/hosted/kubernetes-datastore/calico-networking/1.7/calico.yaml _old/master/getting-started/kubernetes/installation/hosted/kubernetes-datastore/calico-networking/1.7/calico.yaml
diff _site/master/getting-started/kubernetes/installation/hosted/kubernetes-datastore/calico-networking/typha/calico.yaml _old/master/getting-started/kubernetes/installation/hosted/kubernetes-datastore/calico-networking/typha/calico.yaml
diff _site/master/getting-started/kubernetes/installation/hosted/kubernetes-datastore/policy-only/1.7/calico.yaml _old/master/getting-started/kubernetes/installation/hosted/kubernetes-datastore/policy-only/1.7/calico.yaml
diff _site/master/getting-started/kubernetes/installation/manifests/app-layer-policy/etcd/calico-networking/calico.yaml _old/master/getting-started/kubernetes/installation/manifests/app-layer-policy/etcd/calico-networking/calico.yaml
diff _site/master/getting-started/kubernetes/installation/manifests/app-layer-policy/kubernetes-datastore/calico-networking/calico.yaml _old/master/getting-started/kubernetes/installation/manifests/app-layer-policy/kubernetes-datastore/calico-networking/calico.yaml
diff _site/master/getting-started/kubernetes/installation/manifests/app-layer-policy/kubernetes-datastore/flannel/canal.yaml _old/master/getting-started/kubernetes/installation/manifests/app-layer-policy/kubernetes-datastore/flannel/canal.yaml
diff _site/master/getting-started/kubernetes/installation/manifests/app-layer-policy/kubernetes-datastore/policy-only/calico.yaml _old/master/getting-started/kubernetes/installation/manifests/app-layer-policy/kubernetes-datastore/policy-only/calico.yaml
