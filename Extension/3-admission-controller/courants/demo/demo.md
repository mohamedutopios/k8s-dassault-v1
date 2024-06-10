Je vais commenter directement les changements dans chaque fichier de déploiement. Voici vos fichiers de déploiement mis à jour avec les commentaires sur les modifications apportées pour intégrer les contrôleurs d'admission courants.

**db-deployment.yaml**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: db
  name: db
  namespace: vote  # Ajouté pour s'assurer qu'il se trouve dans le namespace "vote"
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
      containers:
        - image: postgres:15.3-alpine3.18
          name: postgres
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
          resources:
            requests:
              memory: 128Mi
              cpu: 100m
            limits:
              memory: 128Mi 
              cpu: 100m
          securityContext:
            allowPrivilegeEscalation: false  # Sécurité ajoutée
            readOnlyRootFilesystem: true     # Sécurité ajoutée
            seccompProfile:
              type: RuntimeDefault            # Sécurité ajoutée
          volumeMounts:
            - mountPath: /var/run/postgresql
              name: run
      volumes:
        - name: run
          emptyDir: {}
```

**redis-deployment.yaml**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: redis
  name: redis
  namespace: vote  # Ajouté pour s'assurer qu'il se trouve dans le namespace "vote"
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
      containers:
        - image: redis:7.0.7-alpine3.17
          name: redis
          livenessProbe:
            initialDelaySeconds: 10
            periodSeconds: 5
            successThreshold: 1
            failureThreshold: 3
            tcpSocket:
              port: 6379
          ports:
            - containerPort: 6379
          resources:
            requests:
              memory: 64Mi
              cpu: 50m
            limits:
              memory: 64Mi 
              cpu: 50m
          securityContext:
            allowPrivilegeEscalation: false  # Sécurité ajoutée
            readOnlyRootFilesystem: true     # Sécurité ajoutée
            seccompProfile:
              type: RuntimeDefault            # Sécurité ajoutée
          volumeMounts:
            - mountPath: /data
              name: redis-data
      volumes:
        - name: redis-data
          emptyDir: {}
```

**result-deployment.yaml**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: result
  name: result
  namespace: vote  # Ajouté pour s'assurer qu'il se trouve dans le namespace "vote"
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
      containers:
        - image: voting/result:v1.0.16
          name: result
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
          livenessProbe:
            initialDelaySeconds: 30
            periodSeconds: 10
            successThreshold: 1
            failureThreshold: 3
            timeoutSeconds: 5
            tcpSocket:
              port: 5000
          ports:
            - containerPort: 5000
          resources:
            requests:
              memory: 256Mi
              cpu: 250m
            limits:
              memory: 256Mi 
          securityContext:
            allowPrivilegeEscalation: false  # Sécurité ajoutée
            # readOnlyRootFilesystem: true   # Peut être ajouté pour plus de sécurité
            # runAsUser: 10000               # Peut être ajouté pour plus de sécurité
            # runAsNonRoot: true             # Peut être ajouté pour plus de sécurité
            seccompProfile:
              type: RuntimeDefault            # Sécurité ajoutée
            capabilities:
              drop:
              - ALL                           # Sécurité ajoutée
```

**result-ui-deployment.yaml**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: result-ui
  name: result-ui
  namespace: vote  # Ajouté pour s'assurer qu'il se trouve dans le namespace "vote"
spec:
  replicas: 1
  selector:
    matchLabels:
      app: result-ui
  template:
    metadata:
      labels:
        app: result-ui
    spec:
      containers:
        - image: voting/result-ui:v1.0.15
          name: result-ui
          imagePullPolicy: Always
          livenessProbe:
            initialDelaySeconds: 10
            periodSeconds: 5
            successThreshold: 1
            failureThreshold: 3
            tcpSocket:
              port: 80
          ports:
            - containerPort: 80
          resources:
            requests:
              memory: 64Mi
              cpu: 50m
            limits:
              memory: 64Mi 
              cpu: 50m
          securityContext:
            allowPrivilegeEscalation: false  # Sécurité ajoutée
            readOnlyRootFilesystem: true     # Sécurité ajoutée
            seccompProfile:
              type: RuntimeDefault            # Sécurité ajoutée
          volumeMounts:
          - name: run
            mountPath: /var/run
          - name: cache
            mountPath: /var/cache/nginx
      volumes:
      - name: run
        emptyDir: {}
      - name: cache
        emptyDir: {}
```

**vote-deployment.yaml**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: vote
  name: vote
  namespace: vote  # Ajouté pour s'assurer qu'il se trouve dans le namespace "vote"
spec:
  replicas: 1
  selector:
    matchLabels:
      app: vote
  template:
    metadata:
      labels:
        app: vote
    spec:
      containers:
        - image: voting/vote:v1.0.13
          name: vote
          imagePullPolicy: Always
          livenessProbe:
            initialDelaySeconds: 30
            periodSeconds: 10
            successThreshold: 1
            failureThreshold: 3
            tcpSocket:
              port: 5000
          ports:
            - containerPort: 5000
          resources:
            requests:
              memory: 256Mi
              cpu: 200m
            limits:
              memory: 256Mi 
          securityContext:
            allowPrivilegeEscalation: false  # Sécurité ajoutée
            readOnlyRootFilesystem: true     # Sécurité ajoutée
            runAsUser: 10000                 # Sécurité ajoutée
            runAsNonRoot: true               # Sécurité ajoutée
            seccompProfile:
              type: RuntimeDefault            # Sécurité ajoutée
            capabilities:
              drop:
              - ALL                           # Sécurité ajoutée
          volumeMounts:
          - name: temp
            mountPath: /tmp
      volumes:
      - name: temp
        emptyDir: {}
```

**vote-ui-deployment.yaml**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: vote-ui
  name: vote-ui
  namespace: vote  # Ajouté pour s'assurer qu'il se trouve dans le namespace "vote"
spec:
  replicas: 1
  selector:
    matchLabels:
      app: vote-ui
  template:
    metadata:
      labels:
        app: vote-ui
    spec:
      containers:
        - image: voting/vote-ui:v1.0.19
          name: vote-ui
          imagePullPolicy: Always
          livenessProbe:
            initialDelaySeconds: 10
            periodSeconds: 5
            successThreshold: 1
            failureThreshold: 3
            tcpSocket:
              port: 80
          ports:
            - containerPort: 80
          resources:
            requests:
              memory: 64Mi
              cpu: 50m
            limits:
              memory: 64Mi 
              cpu: 50m
          securityContext:
            allowPrivilegeEscalation: false  # Sécurité ajoutée
            readOnlyRootFilesystem: true     # Sécurité ajoutée
            seccompProfile:
              type: RuntimeDefault            # Sécurité ajoutée
          volumeMounts:
          - name: run
            mountPath: /var/run
          - name: cache
            mountPath: /var/cache/nginx
      volumes:
      - name: run
        emptyDir: {}
      - name: cache
        emptyDir: {}
```

**worker-deployment.yaml**
```

yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: worker
  name: worker
  namespace: vote  # Ajouté pour s'assurer qu'il se trouve dans le namespace "vote"
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
      containers:
        - image: voting/worker:v1.0.15
          name: worker
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
          resources:
            requests:
              memory: 64Mi
              cpu: 50m
            limits:
              memory: 64Mi 
              cpu: 50m
          securityContext:
            allowPrivilegeEscalation: false  # Sécurité ajoutée
            readOnlyRootFilesystem: true     # Sécurité ajoutée
            runAsUser: 10000                 # Sécurité ajoutée
            runAsNonRoot: true               # Sécurité ajoutée
            seccompProfile:
              type: RuntimeDefault            # Sécurité ajoutée
            capabilities:
              drop:
              - ALL                           # Sécurité ajoutée
```

### Application des Configurations

Appliquez les fichiers de déploiement mis à jour dans votre cluster Kubernetes :
```bash
kubectl apply -f db-deployment.yaml
kubectl apply -f redis-deployment.yaml
kubectl apply -f result-deployment.yaml
kubectl apply -f result-ui-deployment.yaml
kubectl apply -f vote-deployment.yaml
kubectl apply -f vote-ui-deployment.yaml
kubectl apply -f worker-deployment.yaml
```

Ces mises à jour incluent des configurations de sécurité et de ressources conformément aux contrôleurs d'admission courants, et chaque déploiement est maintenant placé dans le namespace `vote`. Assurez-vous que les `LimitRange`, `ResourceQuota`, et `PodSecurityPolicy` sont bien appliqués dans ce namespace pour garantir le respect des règles.