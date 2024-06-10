L'erreur que vous rencontrez indique que la commande `curl` n'a pas pu résoudre l'URL fournie. Cela peut se produire en raison de plusieurs raisons, comme des problèmes de connectivité ou des erreurs de formatage dans la commande.

Voici comment vous pouvez corriger cela et télécharger kubectl correctement :

1. **Télécharger kubectl :**

   Utilisez la commande suivante pour télécharger directement kubectl. Assurez-vous d'utiliser des guillemets simples pour éviter l'interpolation de variables par le shell :

   ```bash
   curl -LO 'https://dl.k8s.io/release/$(curl -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl'
   ```

   Si cela ne fonctionne pas, vous pouvez essayer de décomposer la commande en étapes plus simples :

   ```bash
   version=$(curl -s https://dl.k8s.io/release/stable.txt)
   curl -LO https://dl.k8s.io/release/$version/bin/linux/amd64/kubectl
   ```

2. **Rendre kubectl exécutable :**

   ```bash
   chmod +x kubectl
   ```

3. **Déplacer kubectl dans un répertoire de votre PATH :**

   ```bash
   sudo mv kubectl /usr/local/bin/
   ```

4. **Vérifier l'installation de kubectl :**

   ```bash
   kubectl version --client
   ```

### Installation complète avec Docker et Minikube

Voici les étapes complètes pour installer Minikube en utilisant Docker comme driver, après avoir installé kubectl :

1. **Installer Docker :**

   Si Docker n'est pas déjà installé, vous pouvez l'installer en suivant les étapes ci-dessous :

   ```bash
   sudo apt-get update
   sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common
   curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
   sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
   sudo apt-get update
   sudo apt-get install -y docker-ce
   sudo systemctl start docker
   sudo systemctl enable docker
   sudo usermod -aG docker $USER
   newgrp docker
   ```

2. **Installer Minikube :**

   ```bash
   curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
   chmod +x minikube-linux-amd64
   sudo mv minikube-linux-amd64 /usr/local/bin/minikube
   ```

3. **Démarrer Minikube avec Docker comme driver :**

   ```bash
   minikube start --driver=docker
   ```

4. **Configurer kubectl pour utiliser le contexte Minikube :**

   ```bash
   kubectl config use-context minikube
   ```

5. **Vérifier que le cluster est en place :**

   ```bash
   kubectl get nodes
   ```

Avec ces étapes, vous devriez avoir un cluster Kubernetes fonctionnel sur Minikube en utilisant Docker comme driver sur Ubuntu. Si vous rencontrez d'autres problèmes, n'hésitez pas à me le faire savoir !