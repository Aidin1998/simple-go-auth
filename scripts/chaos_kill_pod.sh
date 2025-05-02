#!/usr/bin/env bash
set -euo pipefail

NAMESPACE=${1:-staging}
DEPLOYMENT=auth-service

# pick a random pod
POD=$(kubectl get pods -n $NAMESPACE -l app=$DEPLOYMENT -o jsonpath='{.items[*].metadata.name}' \
  | tr ' ' '\n' | shuf -n1)

echo "Killing pod $POD in $NAMESPACE..."
kubectl delete pod $POD -n $NAMESPACE

echo "Waiting for rollout..."
kubectl rollout status deployment/$DEPLOYMENT -n $NAMESPACE
echo "Chaos drill complete."
