openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj "/CN=mutating-webhook"


docker build -t mohamed1780/mutating-webhook:latest .
docker push mohamed1780/mutating-webhook:latest


kubectl create secret tls mutating-webhook-certs --cert=cert.pem --key=key.pem -n vote


cat cert.pem | base64 | tr -d '\n'


kubectl apply -f mutating-webhook-configuration.yaml


kubectl apply -f test-pod.yaml


kubectl get pod test-pod -n vote -o yaml



