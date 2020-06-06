ATTEMPTS=0
ROLLOUT_STATUS_CMD="kubectl --insecure-skip-tls-verify --kubeconfig=kubeconfig.yml rollout status deployment/lineblocs-internals"
until $ROLLOUT_STATUS_CMD || [ $ATTEMPTS -eq 60 ]; do
  $ROLLOUT_STATUS_CMD
  ATTEMPTS=$((attempts + 1))
  sleep 10
done
