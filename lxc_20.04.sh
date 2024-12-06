## DO NOT RUN IT AS A SCRIPT
# You were warned
lxc launch ubuntu:20.04 test-container
lxc config set test-container security.nesting true
lxc config set test-container security.privileged true
lxc config device add test-container myport8080 proxy listen=tcp:0.0.0.0:8080 connect=tcp:127.0.0.1:8080
lxc restart test-container
lxc exec test-container -- bash

apt update && apt upgrade -y

apt install -y docker.io make protobuf-compiler
# Go
wget -c https://golang.org/dl/go1.23.3.linux-amd64.tar.gz -O - | sudo tar -xz -C /tmp/
sudo mv /tmp/go /usr/local/
mkdir -p ~/go/bin/
export PATH=$PATH:/usr/local/go/bin
export GOPATH=~/go
export PATH=$PATH:~/go/bin

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

systemctl start docker
systemctl enable docker

mkdir -p ~/.docker/cli-plugins/
curl -SL https://github.com/docker/compose/releases/download/v2.27.0/docker-compose-linux-x86_64 -o ~/.docker/cli-plugins/docker-compose
chmod +x ~/.docker/cli-plugins/docker-compose

lxc config device add test-container project disk source=$(pwd) path=/root/project

# Push project somehow into container
# Either with `lxc file push project.tar.gz test-container/root/`
# Or anyway you want
cd path/to/project

make gen_proto

docker compose up

# Open localhost:8080 in your browser
