sudo echo "y" | sudo kubeadm reset

sudo kubeadm init --pod-network-cidr=10.244.0.0/16 --image-repository=registry.cn-hangzhou.aliyuncs.com/google_containers

sudo rm -rf $HOME/.kube
sudo mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

sudo sysctl net.bridge.bridge-nf-call-iptables=1

sudo kubectl apply -f ./kube-flannel.yml

echo "123456" | su

rm -rf $HOME/.kube
mkdir -p $HOME/.kube
cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
chown $(id -u):$(id -g) $HOME/.kube/config

exit
