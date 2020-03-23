sudo echo "y" | sudo kubeadm reset

sudo kubeadm init --pod-network-cidr=10.244.0.0/16 --image-repository=registry.cn-hangzhou.aliyuncs.com/google_containers

sudo rm -rf $HOME/.kube
sudo mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

sudo sysctl net.bridge.bridge-nf-call-iptables=1

sudo kubectl apply -f ./kube-flannel.yml

sudo rm -rf /root/.kube
sudo mkdir -p /root/.kube
sudo cp -i /etc/kubernetes/admin.conf /root/.kube/config
sudo chown $(id -u):$(id -g) /root/.kube/config

##########################################
# Install the kubernetes-dashboard
# reference: https://zhuanlan.zhihu.com/p/91731765
##########################################

# install the kubernetes-dashboard.v2.0.0.beta4
sudo kubectl apply -f kubernetes-dashboard/kubernetes-dashboard-v2.0.0.beta4.yaml

# install the openssl certificates
# NOTE: IP in certs.sh needs to change!!! `ifconfig`
sudo bash kubernetes-dashboard/key/certs.sh

# user-binding
sudo kubectl create -f kubernetes-dashboard/admin-user.yaml
sudo kubectl create -f kubernetes-dashboard/admin-user-role-binding.yaml

# get the token for access
alias token="kubectl -n kubernetes-dashboard describe secret $(kubectl -n kubernetes-dashboard get secret | grep admin-user | awk '{print $1}')"
token

# print the Master IP, get from `ifconfig`, which is the same as in certs.sh!
print ("\n192.168.206.165\n")  

# print the service port
kubectl get service kubernetes-dashboard -n kubernetes-dashboard


#########################################
# Install the metrics-server
#########################################

# allow to deply pods to master
kubectl taint nodes --all node-role.kubernetes.io/master-

# deploy all components about metrics-server
kubectl create -f metrics-server/deploy/kubernetes/





