# pym.sh is used to install all yaml resources in deploy

# create CSI and storageclass resources
printf "STEP1 starting..."
cd deploy/csi-plugin/
kubectl create -f rbac.yaml
kubectl create -f csidriver.yaml
kubectl create -f csinodeinfo.yaml
printf "STEP1 finished.\n"

# create resources related to lvm
printf "\nSTEP2 starting..."
cd lvm/
kubectl create -f lvm-attacher.yaml
kubectl create -f lvm-plugin.yaml
kubectl create -f storageclass.yaml
printf "STEP2 finished.\n"

# change directory to deploy/
# define new crd
printf "\nSTEP3 starting..."
cd ../../
# shellcheck disable=SC2006
kubectl create -f crd.yaml
printf "\nSTEP3 finished\n"

# change directory to storage/
cd ../

