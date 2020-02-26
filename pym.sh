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
printf "STEP2 finished.\n"

# change directory to storage/
cd ../../..
# shellcheck disable=SC2006
echo "Change to storage/: " `pwd`


