Bien sûr ! Voici les déploiements corrigés afin de respecter les contraintes définies dans la `PodSecurityPolicy` très restrictive que nous avons créée.

### Corrected Deployments

#### db-deployment.yaml
**db-deployment.yaml**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: db
  name: db
  namespace: vote
spec:
  replicas: 1
  selector:
    matchLabels:
      app: db
  template:
    metadata:
      labels:
        app: db
    spec:
      serviceAccountName: vote-sa # Utilisation du ServiceAccount
      containers:
        - image: postgres:15.3-alpine3.18
          name: postgres
          securityContext:
            runAsUser: 1000              # Correct: Running as non-root user
            fsGroup: 2000                # Correct: Using required fsGroup
            capabilities:
              drop: ["ALL"]              # Correct: Dropping all capabilities
          env:
            - name: POSTGRES_USER
              valueFrom:
                secretKeyRef:
                  name: db
                  key: username
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: db
                  key: password
          ports:
            - containerPort: 5432
          volumeMounts:
            - mountPath: /var/run/postgresql
              name: run
              readOnly: true              # Correct: Read-only root filesystem
      volumes:
        - name: run
          emptyDir: {}
```

#### redis-deployment.yaml
**redis-deployment.yaml**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: redis
  name: redis
  namespace: vote
spec:
  replicas: 1
  selector:
    matchLabels:
      app: redis
  template:
    metadata:
      labels:
        app: redis
    spec:
      serviceAccountName: vote-sa # Utilisation du ServiceAccount
      containers:
        - image: redis:7.0.7-alpine3.17
          name: redis
          securityContext:
            runAsUser: 1000              # Correct: Running as non-root user
            runAsGroup: 1000             # Correct: Running with the required group
            allowPrivilegeEscalation: false # Correct: Not allowing privilege escalation
          ports:
            - containerPort: 6379
          volumeMounts:
            - mountPath: /data
              name: redis-data
              readOnly: true              # Correct: Read-only volume
      volumes:
        - name: redis-data
          emptyDir: {}                    # Correct: Using emptyDir volume
```

#### result-deployment.yaml
**result-deployment.yaml**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: result
  name: result
  namespace: vote
spec:
  replicas: 1
  selector:
    matchLabels:
      app: result
  template:
    metadata:
      labels:
        app: result
    spec:
      serviceAccountName: vote-sa # Utilisation du ServiceAccount
      containers:
        - image: voting/result:v1.0.16
          name: result
          securityContext:
            readOnlyRootFilesystem: true # Correct: Read-only root filesystem
            seLinuxOptions:               # Correct: Using the required SELinux context
              level: 's0:c123,c456'
          env:
            - name: POSTGRES_USER
              valueFrom:
                secretKeyRef:
                  name: db
                  key: username
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: db
                  key: password
          imagePullPolicy: Always
          ports:
            - containerPort: 5000
          volumeMounts:
            - mountPath: /var/run/postgresql
              name: run
              readOnly: true              # Correct: Read-only volume
      volumes:
        - name: run
          emptyDir: {}
```

#### worker-deployment.yaml
**worker-deployment.yaml**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: worker
  name: worker
  namespace: vote
spec:
  replicas: 1
  selector:
    matchLabels:
      app: worker
  template:
    metadata:
      labels:
        app: worker
    spec:
      serviceAccountName: vote-sa # Utilisation du ServiceAccount
      containers:
        - image: voting/worker:v1.0.15
          name: worker
          securityContext:
            runAsUser: 1000              # Correct: Running as non-root user
            fsGroup: 2000                # Correct: Using required fsGroup
            capabilities:
              drop: ["ALL"]              # Correct: Dropping all capabilities
            seLinuxOptions:               # Correct: Using the required SELinux context
              level: 's0:c123,c456'
          env:
            - name: POSTGRES_USER
              valueFrom:
                secretKeyRef:
                  name: db
                  key: username
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: db
                  key: password
          imagePullPolicy: Always
          ports:
            - containerPort: 5000
          volumeMounts:
            - mountPath: /var/run/postgresql
              name: run
              readOnly: true              # Correct: Read-only volume
      volumes:
        - name: run
          emptyDir: {}
```

### Application des Configurations

Appliquez les fichiers de configuration corrigés dans votre cluster Kubernetes :

```bash
kubectl apply -f pod-security-policy.yaml
kubectl apply -f db-deployment.yaml
kubectl apply -f redis-deployment.yaml
kubectl apply -f result-deployment.yaml
kubectl apply -f worker-deployment.yaml
```

Ces déploiements corrigés devraient maintenant respecter toutes les contraintes définies par la `PodSecurityPolicy` restrictive et réussir à se lancer dans le cluster Kubernetes.