import difflib
import yaml
import StringIO

manifests = [
    "master/getting-started/kubernetes/installation/hosted/calico.yaml",
    "master/getting-started/kubernetes/installation/hosted/canal/canal-etcd.yaml",
    "master/getting-started/kubernetes/installation/hosted/canal/canal.yaml",
    "master/getting-started/kubernetes/installation/hosted/kubernetes-datastore/calico-networking/1.7/calico.yaml",
    "master/getting-started/kubernetes/installation/hosted/kubernetes-datastore/calico-networking/typha/calico.yaml",
    "master/getting-started/kubernetes/installation/hosted/kubernetes-datastore/policy-only/1.7/calico.yaml",
    "master/getting-started/kubernetes/installation/manifests/app-layer-policy/etcd/calico-networking/calico.yaml",
    "master/getting-started/kubernetes/installation/manifests/app-layer-policy/kubernetes-datastore/calico-networking/calico.yaml",
    "master/getting-started/kubernetes/installation/manifests/app-layer-policy/kubernetes-datastore/flannel/canal.yaml",
    "master/getting-started/kubernetes/installation/manifests/app-layer-policy/kubernetes-datastore/policy-only/calico.yaml"
]

def getmap(filename):
    fold = open(filename, 'r')
    sold = fold.read()
    objs = sold.split("---")

    m = {}

    for obj in objs:
        ffff = yaml.load(obj, Loader=yaml.Loader)
        if ffff != None:
            m[ffff["kind"] + "/" + ffff["metadata"]["name"]] = obj

    return m


for manifest in manifests:
    print("############################### " + manifest + " ###################################")
    old = "_old/" + manifest
    neww = "_site/" + manifest

    oldmap = getmap(old)
    newmap = getmap(neww)

    for key in oldmap.keys():
        if key not in newmap:
            print(key + " not found in helm templates for " + manifest)
            continue

        jekyll = oldmap[key].splitlines(1)
        helm = newmap[key].splitlines(1)

        diff = difflib.unified_diff(jekyll, helm)

            # import pdb; pdb.set_trace()
        print(key)
        print("".join(diff))
    print("############################### " + manifest + " ###################################")
    _ = raw_input("continue?")
# read
# master/getting-started/kubernetes/installation/hosted/calico.yaml