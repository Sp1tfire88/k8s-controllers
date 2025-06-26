#!/bin/bash

# Exit on error
set -e

# Default cluster name
CLUSTER_NAME="${1:-cluster-a}"
CLUSTER_ID=""

# Port mappings for different clusters
case "$CLUSTER_NAME" in
    cluster-a)
        CLUSTER_ID="a"
        ETCD_CLIENT_PORT=2379
        ETCD_PEER_PORT=2380
        API_SERVER_PORT=6443
        KUBELET_PORT=10250
        ;;
    cluster-b)
        CLUSTER_ID="b"
        ETCD_CLIENT_PORT=2479
        ETCD_PEER_PORT=2480
        API_SERVER_PORT=6543
        KUBELET_PORT=10350
        ;;
    cluster-c)
        CLUSTER_ID="c"
        ETCD_CLIENT_PORT=2579
        ETCD_PEER_PORT=2580
        API_SERVER_PORT=6643
        KUBELET_PORT=10450
        ;;
    *)
        echo "Supported cluster names: cluster-a, cluster-b, cluster-c"
        exit 1
        ;;
esac

echo "Setting up Kubernetes cluster: $CLUSTER_NAME (ports: API=$API_SERVER_PORT, ETCD=$ETCD_CLIENT_PORT)"

# Cluster-specific directories
CLUSTER_DIR="./clusters/$CLUSTER_NAME"
KUBEBUILDER_DIR="$CLUSTER_DIR/kubebuilder"
ETCD_DIR="$CLUSTER_DIR/etcd"
KUBELET_DIR="$CLUSTER_DIR/kubelet"
CONTAINERD_DIR="$CLUSTER_DIR/containerd"
KUBE_CONFIG_DIR="$CLUSTER_DIR/.kube"

# Function to check if a process is running for specific cluster
is_running() {
    pgrep -f "$1.*$CLUSTER_NAME" >/dev/null || pgrep -f "$1.*$ETCD_CLIENT_PORT" >/dev/null || pgrep -f "$1.*$API_SERVER_PORT" >/dev/null
}

# Function to check if all components are running
check_running() {
    is_running "etcd" && \
    is_running "kube-apiserver" && \
    is_running "kube-controller-manager" && \
    is_running "kube-scheduler" && \
    is_running "kubelet" && \
    is_running "containerd"
}

# Function to kill process if running
stop_process() {
    echo "Stopping $1 for $CLUSTER_NAME..."
    pkill -f "$1.*$CLUSTER_NAME" 2>/dev/null || true
    pkill -f "$1.*$ETCD_CLIENT_PORT" 2>/dev/null || true
    pkill -f "$1.*$API_SERVER_PORT" 2>/dev/null || true
    sleep 2
}

download_components() {
    # Create necessary directories if they don't exist
    sudo mkdir -p $KUBEBUILDER_DIR/bin
    sudo mkdir -p $CLUSTER_DIR/etc/cni/net.d
    sudo mkdir -p $KUBELET_DIR
    sudo mkdir -p $CLUSTER_DIR/etc/kubernetes/manifests
    sudo mkdir -p $CLUSTER_DIR/var/log/kubernetes
    sudo mkdir -p $CLUSTER_DIR/etc/containerd/
    sudo mkdir -p $CONTAINERD_DIR/run

    # Download kubebuilder tools if not present
    if [ ! -f "$KUBEBUILDER_DIR/bin/etcd" ]; then
        echo "Downloading kubebuilder tools for $CLUSTER_NAME..."
        curl -L https://storage.googleapis.com/kubebuilder-tools/kubebuilder-tools-1.30.0-linux-amd64.tar.gz -o /tmp/kubebuilder-tools-$CLUSTER_ID.tar.gz
        sudo tar -C $KUBEBUILDER_DIR --strip-components=1 -zxf /tmp/kubebuilder-tools-$CLUSTER_ID.tar.gz
        rm /tmp/kubebuilder-tools-$CLUSTER_ID.tar.gz
        sudo chmod -R 755 $KUBEBUILDER_DIR/bin
    fi

    if [ ! -f "$KUBEBUILDER_DIR/bin/kubelet" ]; then
        echo "Downloading kubelet for $CLUSTER_NAME..."
        sudo curl -L "https://dl.k8s.io/v1.30.0/bin/linux/amd64/kubelet" -o $KUBEBUILDER_DIR/bin/kubelet
        sudo chmod 755 $KUBEBUILDER_DIR/bin/kubelet
    fi

    # Install CNI components if not present (shared across clusters)
    if [ ! -d "/opt/cni" ]; then
        sudo mkdir -p /opt/cni
        
        echo "Installing containerd..."
        wget https://github.com/containerd/containerd/releases/download/v2.0.5/containerd-static-2.0.5-linux-amd64.tar.gz -O /tmp/containerd.tar.gz
        sudo tar zxf /tmp/containerd.tar.gz -C /opt/cni/
        rm /tmp/containerd.tar.gz

        echo "Installing runc..."
        sudo curl -L "https://github.com/opencontainers/runc/releases/download/v1.2.6/runc.amd64" -o /opt/cni/bin/runc

        echo "Installing CNI plugins..."
        wget https://github.com/containernetworking/plugins/releases/download/v1.6.2/cni-plugins-linux-amd64-v1.6.2.tgz -O /tmp/cni-plugins.tgz
        sudo tar zxf /tmp/cni-plugins.tgz -C /opt/cni/bin/
        rm /tmp/cni-plugins.tgz

        # Set permissions for all CNI components
        sudo chmod -R 755 /opt/cni
    fi

    if [ ! -f "$KUBEBUILDER_DIR/bin/kube-controller-manager" ]; then
        echo "Downloading additional components for $CLUSTER_NAME..."
        sudo curl -L "https://dl.k8s.io/v1.30.0/bin/linux/amd64/kube-controller-manager" -o $KUBEBUILDER_DIR/bin/kube-controller-manager
        sudo curl -L "https://dl.k8s.io/v1.30.0/bin/linux/amd64/kube-scheduler" -o $KUBEBUILDER_DIR/bin/kube-scheduler
        sudo curl -L "https://dl.k8s.io/v1.30.0/bin/linux/amd64/cloud-controller-manager" -o $KUBEBUILDER_DIR/bin/cloud-controller-manager
        sudo chmod 755 $KUBEBUILDER_DIR/bin/kube-controller-manager
        sudo chmod 755 $KUBEBUILDER_DIR/bin/kube-scheduler
        sudo chmod 755 $KUBEBUILDER_DIR/bin/cloud-controller-manager
    fi
}

setup_configs() {
    # Generate certificates and tokens if they don't exist
    if [ ! -f "$CLUSTER_DIR/sa.key" ]; then
        openssl genrsa -out $CLUSTER_DIR/sa.key 2048
        openssl rsa -in $CLUSTER_DIR/sa.key -pubout -out $CLUSTER_DIR/sa.pub
    fi

    if [ ! -f "$CLUSTER_DIR/token.csv" ]; then
        TOKEN="1234567890-$CLUSTER_ID"
        echo "${TOKEN},admin,admin,system:masters" > $CLUSTER_DIR/token.csv
    fi

    # Always regenerate and copy CA certificate to ensure it exists
    echo "Generating CA certificate for $CLUSTER_NAME..."
    openssl genrsa -out $CLUSTER_DIR/ca.key 2048
    openssl req -x509 -new -nodes -key $CLUSTER_DIR/ca.key -subj "/CN=kubelet-ca-$CLUSTER_NAME" -days 365 -out $CLUSTER_DIR/ca.crt
    sudo mkdir -p $KUBELET_DIR/pki
    sudo cp $CLUSTER_DIR/ca.crt $KUBELET_DIR/ca.crt
    sudo cp $CLUSTER_DIR/ca.crt $KUBELET_DIR/pki/ca.crt

    # Set up kubeconfig for this cluster
    mkdir -p $KUBE_CONFIG_DIR
    export KUBECONFIG=$KUBE_CONFIG_DIR/config
    
    $KUBEBUILDER_DIR/bin/kubectl config set-credentials $CLUSTER_NAME-user --token=1234567890-$CLUSTER_ID
    $KUBEBUILDER_DIR/bin/kubectl config set-cluster $CLUSTER_NAME --server=https://127.0.0.1:$API_SERVER_PORT --insecure-skip-tls-verify
    $KUBEBUILDER_DIR/bin/kubectl config set-context $CLUSTER_NAME-context --cluster=$CLUSTER_NAME --user=$CLUSTER_NAME-user --namespace=default 
    $KUBEBUILDER_DIR/bin/kubectl config use-context $CLUSTER_NAME-context

    # Configure CNI with cluster-specific subnet
    SUBNET_BASE=$((200 + CLUSTER_ID))  # cluster-a=200, cluster-b=201, etc.
    cat <<EOF | sudo tee $CLUSTER_DIR/etc/cni/net.d/10-mynet.conf
{
    "cniVersion": "0.3.1",
    "name": "mynet-$CLUSTER_NAME",
    "type": "bridge",
    "bridge": "cni$CLUSTER_ID",
    "isGateway": true,
    "ipMasq": true,
    "ipam": {
        "type": "host-local",
        "subnet": "10.$SUBNET_BASE.0.0/16",
        "routes": [
            { "dst": "0.0.0.0/0" }
        ]
    }
}
EOF

    # Configure containerd with cluster-specific socket
    cat <<EOF | sudo tee $CLUSTER_DIR/etc/containerd/config.toml
version = 3

[grpc]
  address = "$CONTAINERD_DIR/run/containerd.sock"

[plugins.'io.containerd.cri.v1.runtime']
  enable_selinux = false
  enable_unprivileged_ports = true
  enable_unprivileged_icmp = true
  device_ownership_from_security_context = false

[plugins.'io.containerd.cri.v1.images']
  snapshotter = "native"
  disable_snapshot_annotations = true

[plugins.'io.containerd.cri.v1.runtime'.cni]
  bin_dir = "/opt/cni/bin"
  conf_dir = "$CLUSTER_DIR/etc/cni/net.d"

[plugins.'io.containerd.cri.v1.runtime'.containerd.runtimes.runc]
  runtime_type = "io.containerd.runc.v2"

[plugins.'io.containerd.cri.v1.runtime'.containerd.runtimes.runc.options]
  SystemdCgroup = false

[plugins.'io.containerd.grpc.v1.cri']
  root = "$CONTAINERD_DIR/lib"
  state = "$CONTAINERD_DIR/run"
EOF

    # Ensure containerd data directory exists with correct permissions
    sudo mkdir -p $CONTAINERD_DIR/lib
    sudo chmod 711 $CONTAINERD_DIR/lib

    # Configure kubelet
    cat << EOF | sudo tee $KUBELET_DIR/config.yaml
apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration
authentication:
  anonymous:
    enabled: true
  webhook:
    enabled: true
  x509:
    clientCAFile: "$KUBELET_DIR/ca.crt"
authorization:
  mode: AlwaysAllow
clusterDomain: "cluster.local"
clusterDNS:
  - "10.0.0.10"
resolvConf: "/etc/resolv.conf"
runtimeRequestTimeout: "15m"
failSwapOn: false
seccompDefault: true
serverTLSBootstrap: false
containerRuntimeEndpoint: "unix://$CONTAINERD_DIR/run/containerd.sock"
staticPodPath: "$CLUSTER_DIR/etc/kubernetes/manifests"
port: $KUBELET_PORT
readOnlyPort: 0
EOF

    # Create required directories with proper permissions
    sudo mkdir -p $KUBELET_DIR/pods
    sudo chmod 750 $KUBELET_DIR/pods
    sudo mkdir -p $KUBELET_DIR/plugins
    sudo chmod 750 $KUBELET_DIR/plugins
    sudo mkdir -p $KUBELET_DIR/plugins_registry
    sudo chmod 750 $KUBELET_DIR/plugins_registry

    # Ensure proper permissions
    sudo chmod 644 $KUBELET_DIR/ca.crt
    sudo chmod 644 $KUBELET_DIR/config.yaml

    # Generate self-signed kubelet serving certificate if not present
    if [ ! -f "$KUBELET_DIR/pki/kubelet.crt" ] || [ ! -f "$KUBELET_DIR/pki/kubelet.key" ]; then
        echo "Generating self-signed kubelet serving certificate for $CLUSTER_NAME..."
        sudo openssl req -x509 -newkey rsa:2048 -nodes \
            -keyout $KUBELET_DIR/pki/kubelet.key \
            -out $KUBELET_DIR/pki/kubelet.crt \
            -days 365 \
            -subj "/CN=$(hostname)-$CLUSTER_NAME"
        sudo chmod 600 $KUBELET_DIR/pki/kubelet.key
        sudo chmod 644 $KUBELET_DIR/pki/kubelet.crt
    fi

    # Create kubeconfig for kubelet
    cat << EOF | sudo tee $KUBELET_DIR/kubeconfig
apiVersion: v1
kind: Config
clusters:
- cluster:
    insecure-skip-tls-verify: true
    server: https://127.0.0.1:$API_SERVER_PORT
  name: $CLUSTER_NAME
contexts:
- context:
    cluster: $CLUSTER_NAME
    user: $CLUSTER_NAME-user
  name: $CLUSTER_NAME-context
current-context: $CLUSTER_NAME-context
users:
- name: $CLUSTER_NAME-user
  user:
    token: 1234567890-$CLUSTER_ID
EOF
}

start() {
    if check_running; then
        echo "Kubernetes components for $CLUSTER_NAME are already running"
        return 0
    fi

    HOST_IP=$(hostname -I | awk '{print $1}')
    
    # Download components if needed
    download_components
    
    # Setup configurations
    setup_configs

    # Start components if not running
    if ! is_running "etcd"; then
        echo "Starting etcd for $CLUSTER_NAME..."
        $KUBEBUILDER_DIR/bin/etcd \
            --name $CLUSTER_NAME \
            --advertise-client-urls http://$HOST_IP:$ETCD_CLIENT_PORT \
            --listen-client-urls http://0.0.0.0:$ETCD_CLIENT_PORT \
            --data-dir $ETCD_DIR \
            --listen-peer-urls http://0.0.0.0:$ETCD_PEER_PORT \
            --initial-cluster $CLUSTER_NAME=http://$HOST_IP:$ETCD_PEER_PORT \
            --initial-advertise-peer-urls http://$HOST_IP:$ETCD_PEER_PORT \
            --initial-cluster-state new \
            --initial-cluster-token $CLUSTER_NAME-token &
    fi

    sleep 3

    if ! is_running "kube-apiserver"; then
        echo "Starting kube-apiserver for $CLUSTER_NAME..."
        $KUBEBUILDER_DIR/bin/kube-apiserver \
            --etcd-servers=http://$HOST_IP:$ETCD_CLIENT_PORT \
            --service-cluster-ip-range=10.0.0.0/24 \
            --bind-address=0.0.0.0 \
            --secure-port=$API_SERVER_PORT \
            --advertise-address=$HOST_IP \
            --authorization-mode=AlwaysAllow \
            --token-auth-file=$CLUSTER_DIR/token.csv \
            --enable-priority-and-fairness=false \
            --allow-privileged=true \
            --profiling=false \
            --storage-backend=etcd3 \
            --storage-media-type=application/json \
            --v=0 \
            --service-account-issuer=https://kubernetes.default.svc.cluster.local \
            --service-account-key-file=$CLUSTER_DIR/sa.pub \
            --service-account-signing-key-file=$CLUSTER_DIR/sa.key &
    fi

    if ! is_running "containerd"; then
        echo "Starting containerd for $CLUSTER_NAME..."
        export PATH=$PATH:/opt/cni/bin:$KUBEBUILDER_DIR/bin
        PATH=$PATH:/opt/cni/bin:/usr/sbin /opt/cni/bin/containerd -c $CLUSTER_DIR/etc/containerd/config.toml &
    fi

    sleep 5

    # Set up environment for this cluster
    export KUBECONFIG=$KUBE_CONFIG_DIR/config

    if ! is_running "kube-scheduler"; then
        echo "Starting kube-scheduler for $CLUSTER_NAME..."
        $KUBEBUILDER_DIR/bin/kube-scheduler \
            --kubeconfig=$KUBECONFIG \
            --leader-elect=false \
            --v=2 \
            --bind-address=0.0.0.0 &
    fi

    if ! is_running "kubelet"; then
        echo "Starting kubelet for $CLUSTER_NAME..."
        PATH=$PATH:/opt/cni/bin:/usr/sbin $KUBEBUILDER_DIR/bin/kubelet \
            --kubeconfig=$KUBELET_DIR/kubeconfig \
            --config=$KUBELET_DIR/config.yaml \
            --root-dir=$KUBELET_DIR \
            --cert-dir=$KUBELET_DIR/pki \
            --tls-cert-file=$KUBELET_DIR/pki/kubelet.crt \
            --tls-private-key-file=$KUBELET_DIR/pki/kubelet.key \
            --hostname-override=$(hostname)-$CLUSTER_NAME \
            --pod-infra-container-image=registry.k8s.io/pause:3.10 \
            --node-ip=$HOST_IP \
            --cgroup-driver=cgroupfs \
            --max-pods=4 \
            --v=1 &
    fi

    sleep 3

    # Create service account and configmap if they don't exist
    $KUBEBUILDER_DIR/bin/kubectl create sa default 2>/dev/null || true
    $KUBEBUILDER_DIR/bin/kubectl create configmap kube-root-ca.crt --from-file=ca.crt=$CLUSTER_DIR/ca.crt -n default 2>/dev/null || true

    # Label the node so static pods with nodeSelector can be scheduled
    NODE_NAME="$(hostname)-$CLUSTER_NAME"
    $KUBEBUILDER_DIR/bin/kubectl label node "$NODE_NAME" node-role.kubernetes.io/master="" --overwrite || true
    $KUBEBUILDER_DIR/bin/kubectl label node "$NODE_NAME" cluster=$CLUSTER_NAME --overwrite || true

    if ! is_running "kube-controller-manager"; then
        echo "Starting kube-controller-manager for $CLUSTER_NAME..."
        PATH=$PATH:/opt/cni/bin:/usr/sbin $KUBEBUILDER_DIR/bin/kube-controller-manager \
            --kubeconfig=$KUBELET_DIR/kubeconfig \
            --leader-elect=false \
            --service-cluster-ip-range=10.0.0.0/24 \
            --cluster-name=$CLUSTER_NAME \
            --root-ca-file=$KUBELET_DIR/ca.crt \
            --service-account-private-key-file=$CLUSTER_DIR/sa.key \
            --use-service-account-credentials=true \
            --v=2 &
    fi

    echo "Waiting for $CLUSTER_NAME components to be ready..."
    sleep 10

    echo "Verifying $CLUSTER_NAME setup..."
    echo "KUBECONFIG=$KUBECONFIG"
    $KUBEBUILDER_DIR/bin/kubectl get nodes
    $KUBEBUILDER_DIR/bin/kubectl get all -A
    $KUBEBUILDER_DIR/bin/kubectl get --raw='/readyz?verbose' || true
    
    echo "Cluster $CLUSTER_NAME is ready!"
    echo "To use this cluster, set: export KUBECONFIG=$KUBE_CONFIG_DIR/config"
}

stop() {
    echo "Stopping Kubernetes components for $CLUSTER_NAME..."
    stop_process "kube-controller-manager"
    stop_process "kubelet"
    stop_process "kube-scheduler"
    stop_process "kube-apiserver"
    stop_process "containerd"
    stop_process "etcd"
    echo "All components for $CLUSTER_NAME stopped"
}

cleanup() {
    stop
    echo "Cleaning up $CLUSTER_NAME..."
    sudo rm -rf $CLUSTER_DIR
    echo "Cleanup complete for $CLUSTER_NAME"
}

status() {
    echo "Status for $CLUSTER_NAME:"
    echo "API Server: $(is_running "kube-apiserver" && echo "Running" || echo "Stopped") (Port: $API_SERVER_PORT)"
    echo "ETCD: $(is_running "etcd" && echo "Running" || echo "Stopped") (Port: $ETCD_CLIENT_PORT)"
    echo "Kubelet: $(is_running "kubelet" && echo "Running" || echo "Stopped") (Port: $KUBELET_PORT)"
    echo "Scheduler: $(is_running "kube-scheduler" && echo "Running" || echo "Stopped")"
    echo "Controller Manager: $(is_running "kube-controller-manager" && echo "Running" || echo "Stopped")"
    echo "Containerd: $(is_running "containerd" && echo "Running" || echo "Stopped")"
    echo ""
    if [ -f "$KUBE_CONFIG_DIR/config" ]; then
        echo "KUBECONFIG: $KUBE_CONFIG_DIR/config"
        echo "To connect: export KUBECONFIG=$KUBE_CONFIG_DIR/config"
    fi
}

case "${2:-start}" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    cleanup)
        cleanup
        ;;
    status)
        status
        ;;
    *)
        echo "Usage: $0 <cluster-name> {start|stop|cleanup|status}"
        echo "Supported cluster names: cluster-a, cluster-b, cluster-c"
        echo ""
        echo "Examples:"
        echo "  $0 cluster-a start"
        echo "  $0 cluster-b stop"
        echo "  $0 cluster-c status"
        exit 1
        ;;
esac