#!/bin/bash

# Multi-cluster manager script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MULTI_CLUSTER_SCRIPT="$SCRIPT_DIR/multi-cluster-setup.sh"

CLUSTERS=("cluster-a" "cluster-b" "cluster-c")

print_usage() {
    echo "Usage: $0 {start-all|stop-all|cleanup-all|status-all|start|stop|cleanup|status} [cluster-name]"
    echo ""
    echo "Commands:"
    echo "  start-all    - Start all clusters (cluster-a, cluster-b, cluster-c)"
    echo "  stop-all     - Stop all clusters"
    echo "  cleanup-all  - Cleanup all clusters"
    echo "  status-all   - Show status of all clusters"
    echo "  start        - Start specific cluster"
    echo "  stop         - Stop specific cluster"
    echo "  cleanup      - Cleanup specific cluster"
    echo "  status       - Show status of specific cluster"
    echo ""
    echo "Examples:"
    echo "  $0 start-all"
    echo "  $0 start cluster-a"
    echo "  $0 status-all"
    echo ""
    echo "Port mappings:"
    echo "  cluster-a: API=6443, ETCD=2379, Kubelet=10250"
    echo "  cluster-b: API=6543, ETCD=2479, Kubelet=10350"
    echo "  cluster-c: API=6643, ETCD=2579, Kubelet=10450"
}

start_all() {
    echo "Starting all clusters..."
    for cluster in "${CLUSTERS[@]}"; do
        echo "=================================="
        echo "Starting $cluster..."
        echo "=================================="
        chmod +x "$MULTI_CLUSTER_SCRIPT"
        "$MULTI_CLUSTER_SCRIPT" "$cluster" start
        echo ""
    done
    echo "All clusters started!"
    echo ""
    show_connection_info
}

stop_all() {
    echo "Stopping all clusters..."
    for cluster in "${CLUSTERS[@]}"; do
        echo "Stopping $cluster..."
        chmod +x "$MULTI_CLUSTER_SCRIPT"
        "$MULTI_CLUSTER_SCRIPT" "$cluster" stop
    done
    echo "All clusters stopped!"
}

cleanup_all() {
    echo "Cleaning up all clusters..."
    for cluster in "${CLUSTERS[@]}"; do
        echo "Cleaning up $cluster..."
        chmod +x "$MULTI_CLUSTER_SCRIPT"
        "$MULTI_CLUSTER_SCRIPT" "$cluster" cleanup
    done
    echo "All clusters cleaned up!"
}

status_all() {
    echo "Status of all clusters:"
    echo "======================"
    for cluster in "${CLUSTERS[@]}"; do
        chmod +x "$MULTI_CLUSTER_SCRIPT"
        "$MULTI_CLUSTER_SCRIPT" "$cluster" status
        echo "----------------------"
    done
}

show_connection_info() {
    echo "Connection Information:"
    echo "======================"
    echo ""
    for cluster in "${CLUSTERS[@]}"; do
        case "$cluster" in
            cluster-a)
                api_port=6443
                ;;
            cluster-b)
                api_port=6543
                ;;
            cluster-c)
                api_port=6643
                ;;
        esac
        
        kubeconfig_path="./clusters/$cluster/.kube/config"
        if [ -f "$kubeconfig_path" ]; then
            echo "Cluster: $cluster"
            echo "  KUBECONFIG: export KUBECONFIG=$kubeconfig_path"
            echo "  API Server: https://127.0.0.1:$api_port"
            echo "  Test connection: kubectl --kubeconfig=$kubeconfig_path get nodes"
            echo ""
        fi
    done
    
    echo "Multi-cluster kubectl usage examples:"
    echo "  # Switch between clusters"
    echo "  export KUBECONFIG=./clusters/cluster-a/.kube/config"
    echo "  kubectl get nodes"
    echo ""
    echo "  export KUBECONFIG=./clusters/cluster-b/.kube/config"  
    echo "  kubectl get nodes"
    echo ""
    echo "  # Use specific kubeconfig without switching"
    echo "  kubectl --kubeconfig=./clusters/cluster-a/.kube/config get pods -A"
    echo "  kubectl --kubeconfig=./clusters/cluster-b/.kube/config get pods -A"
}

# Check if multi-cluster-setup.sh exists
if [ ! -f "$MULTI_CLUSTER_SCRIPT" ]; then
    echo "Error: multi-cluster-setup.sh not found in the same directory"
    echo "Please ensure both scripts are in the same directory"
    exit 1
fi

case "${1:-}" in
    start-all)
        start_all
        ;;
    stop-all)
        stop_all
        ;;
    cleanup-all)
        cleanup_all
        ;;
    status-all)
        status_all
        ;;
    start|stop|cleanup|status)
        if [ -z "${2:-}" ]; then
            echo "Error: cluster name required for individual operations"
            echo "Available clusters: ${CLUSTERS[*]}"
            exit 1
        fi
        chmod +x "$MULTI_CLUSTER_SCRIPT"
        "$MULTI_CLUSTER_SCRIPT" "$2" "$1"
        ;;
    *)
        print_usage
        exit 1
        ;;
esac