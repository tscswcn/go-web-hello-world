1, ### Task 0: Install a ubuntu 16.04 server 64-bit

download http://releases.ubuntu.com/16.04/ubuntu-16.04.6-server-amd64.iso
install, meet one bug, Chineses lanage will make one issue, the disk issue

sudo mkdir /etc/ipstables;
vi /etc/iptables/rules.v4


/etc/iptables/rules.v4
*filter
:INPUT DROP [0:0]
:FORWARD ACCEPT [0:0]
:OUTPUT ACCEPT [0:0]
:syn-flood - [0:0]
-A INPUT -i lo -j ACCEPT
-A INPUT -m state --state RELATED,ESTABLISHED -j ACCEPT
-A INPUT -p tcp  -m tcp --dport 22 -j ACCEPT
-A INPUT -p tcp  -m tcp --dport 80 -j ACCEPT
-A INPUT -p tcp  -m tcp --dport 2222 -j ACCEPT
-A INPUT -p tcp  -m tcp --dport 8080 -j ACCEPT
-A INPUT -p tcp  -m tcp --dport 8081 -j ACCEPT
-A INPUT -p tcp  -m tcp --dport 8082 -j ACCEPT
-A INPUT -p tcp  -m tcp --dport 31080 -j ACCEPT
-A INPUT -p tcp  -m tcp --dport 31081 -j ACCEPT


-A INPUT -p icmp -m limit --limit 100/sec --limit-burst 100 -j ACCEPT
-A INPUT -p icmp -m limit --limit 1/s --limit-burst 10 -j ACCEPT
-A INPUT -p tcp -m tcp --tcp-flags FIN,SYN,RST,ACK SYN -j syn-flood
-A INPUT -j REJECT --reject-with icmp-host-prohibited
-A syn-flood -p tcp -m limit --limit 3/sec --limit-burst 6 -j RETURN
-A syn-flood -j REJECT --reject-with icmp-port-unreachable
COMMIT


:wq save to exit
sudo iptables-restore < /etc/iptables/rules.v4 

sudo iptables -t nat -A PREROUTING -p tcp --dport 22 -j REDIRECT --to-ports 2222
sudo iptables -t nat -A PREROUTING -p tcp --dport 80 -j REDIRECT --to-ports 8080
sudo iptables -A INPUT -p tcp -m tcp --dport 80 -j ACCEPT


### Task 1: Update system

sudo apt update
sudo apt upgrade

### Task 2: install gitlab-ce version in the host

sudo apt-get install -y curl openssh-server ca-certificates
sudo apt-get install -y postfix
#选择“Internet Site”并按Enter键，其他选择则默认

curl https://packages.gitlab.com/install/repositories/gitlab/gitlab-ce/script.deb.sh | sudo bash

sudo apt-get update
sudo EXTERNAL_URL="http://10.0.1.6" apt-get install gitlab-ce  
## 我的机器的IP 地址是 10.0.1.6


配置https访问
#将ssl证书放入/etc/gitlab/ssl
sudo vim /etc/gitlab/gitlab.rb
external_url 'https://域名'
nginx['enable'] = true
nginx['redirect_http_to_https'] =true
nginx['ssl_certificate'] = "/etc/gitlab/ssl/域名的ssl证书.crt"
nginx['ssl_certificate_key'] = "/etc/gitlab/ssl/域名的ssl证书.key"
#保存退出:wq

# 我在 10.0.1.6的这个机器上没做，在的本地机器上 设置了1个域名。

sudo gitlab-ctl reconfigure
sudo gitlab-ctl restart
service sshd start
service postfix start

visit ip地址http://106.12.168.234 设置密码，登陆gitlab 
# 我自己有个域名cloud4u.top,绑在地址上，因此访问cloud4u.top 也可以

### Task 3: create a demo group/project in gitlab

1,create group and project go-web-hello-world.
2,use golang wite one go program: hello.go 
package main

import (
"fmt"
"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Go Web Hello World!")
		//fmt.print("Go Web Hello World!")
		//fmt.Print("Go Web Hello World!")
	})

	http.ListenAndServe(":8081", nil)
}

3,run the application after built, curl http://127.0.0.1:8081
get "Go Web Hello World!" messages 
4, write one Dockerfile to build one image  
#source
FROM golang:latest
#author
MAINTAINER luopeng "755200@qq.com"
#workdir
WORKDIR $GOPATH/
#add code
ADD hello.go $GOPATH/
#build
RUN go build hello.go
#expose
EXPOSE 8081
#entrypoint
ENTRYPOINT  ["./hello"]

docker build -t go-web-hello-world .

5,### Task 5: install docker
https://docs.docker.com/install/linux/docker-ce/ubuntu/
sudo apt-get remove docker docker-engine docker-ce docker.io 
sudo apt-get update
apt-get install apt-transport-https  ca-certificates curl 

curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository \ "deb [arch=amd64] https://download.docker.com/linux/ubuntu \ $(lsb_release -cs) \stable"
apt-get update
sudo apt-get install docker-engine
docker version

docker run -d -p 8082:8081 go-web-hello-world 

curl 127.0.0.1:8082 
Go Web Hello World!

### Task 7: push image to dockerhub
docker tag go-web-hello-world tscswcn/go-web-hello-world:v0.1 
docker push tscswcn/go-web-hello-world:v0.1 
#my dockerhub  的账号是 tscswcn



Expect output: https://hub.docker.com/repository/docker/your_dockerhub_id/go-web-hello-world

### Task 8: document the procedure in a MarkDown file

create a README.md file in the gitlab repo and add the technical procedure above (0-7) in this file

-----------------------------------

### Task 9: install a single node Kubernetes cluster using kubeadm
https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/

Check in the admin.conf file into the gitlab repo
 swapoff -a 
 wget https://mirrors.aliyun.com/kubernetes/apt/pool/kubectl_1.15.5-00_amd64_bc99b7c6736e0d254263f270a4fec7e303fd6cb77d5ee97209ea7b34e539e4bc.deb
 wget https://mirrors.aliyun.com/kubernetes/apt/pool/kubelet_1.15.5-00_amd64_feba4d4831a02a994a71708885f0fd043b983ae2787a6d2eb1f1ae80b0f199f0.deb
 wget https://mirrors.aliyun.com/kubernetes/apt/pool/kubeadm_1.15.5-00_amd64_cffe0070e6279c8cdca599202eabeab1774b3265d0c590933d5e1115e739668b.deb


apt-get update && apt-get install -y apt-transport-https
curl https://mirrors.aliyun.com/kubernetes/apt/doc/apt-key.gpg | apt-key add - 
cat <<EOF >/etc/apt/sources.list.d/kubernetes.list
deb https://mirrors.aliyun.com/kubernetes/apt/ kubernetes-xenial main
EOF
apt-get update
apt-get install -y kubelet kubeadm kubectl

整理好 需要k8s1.15.5 需要的镜像，从ali网站上下载，然后tag成gcr.k8s.io的镜像

kubeadm init --pod-network-cidr=192.168.0.0/16
让后  taintd 掉master 节点，就可以创建 deployedment 
遇到几个问题，解决了问题



### Task 10: deploy the hello world container

in the kubernetes above and expose the service to nodePort 31080

Expect output:
```
curl http://127.0.0.1:31080
Go Web Hello World!

kubectl run goweb1 --image=tscswcn/go-web-hello-world:v1 --port=8081 --hostport=31080
yaml file has check in gitlab repo


------------------------------------

### Task 11: install kubernetes dashboard

and expose the service to nodeport 31081

https://kubernetes.io/docs/tasks/access-application-cluster/web-ui-dashboard/

Expect output: https://127.0.0.1:31081 (asking for token)

docker pull registry.cn-hangzhou.aliyuncs.com/rsqlh/kubernetes-dashboard:v1.10.1
docker tag registry.cn-hangzhou.aliyuncs.com/rsqlh/kubernetes-dashboard:v1.10.1 k8s.gcr.io/kubernetes-dashboard-amd64:v1.10.1 



### Task 12: generate token for dashboard login in task 11

figure out how to generate token to login to the dashboard and publish the procedure to the gitlab.

--------------------------------------

### Task 13: publish your work

push all files/procedures in your local gitlab repo to remote github repo (e.g. https://github.com/your_github_id/go-web-hello-world)

if this is for an interview session, please send it to bo.cui@ericsson.com, no later than two calendar days after the interview.

has submit to github.com/tscswcn/






 













