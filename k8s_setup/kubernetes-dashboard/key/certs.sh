#生成证书
openssl genrsa -out dashboard.key 2048 

#我这里写的自己的master节点,ip 地址可以由ifconfig获得
openssl req -new -out dashboard.csr -key dashboard.key -subj '/CN=192.168.206.165'
openssl x509 -req -in dashboard.csr -signkey dashboard.key -out dashboard.crt 

#删除原有的证书secret
kubectl delete secret kubernetes-dashboard-certs -n kubernetes-dashboard

#创建新的证书secret
kubectl create secret generic kubernetes-dashboard-certs --from-file=dashboard.key --from-file=dashboard.crt -n kubernetes-dashboard

#重启kubernetes-dashboard命名空间下的所有pod
kubectl delete po --all -n kubernetes-dashboard
