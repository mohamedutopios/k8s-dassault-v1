En programmation d'applications web, un **controller d'admission** (ou admission controller) est une composante essentielle du processus de gestion des requêtes et des réponses. Il s'agit d'un logiciel ou d'un module qui intervient avant ou après que les requêtes soient traitées par le serveur d'application. Voici les types d'admission controllers les plus courants :

1. **Admission Controller Mutating** : Modifie les requêtes avant qu'elles ne soient traitées par le serveur. Par exemple, il peut ajouter des labels, des annotations, ou modifier des configurations de ressources.

2. **Admission Controller Validating** : Valide les requêtes pour s'assurer qu'elles respectent certaines règles ou politiques avant d'être acceptées par le serveur. Il peut rejeter les requêtes qui ne respectent pas les critères définis.

3. **Webhook Admission Controller** : Utilise des webhooks pour déléguer la logique de mutation ou de validation à des services externes. Ceci permet de centraliser la logique d'admission dans des services indépendants.

4. **Initializers** : Un type d'admission controller qui permet d'ajouter des initialisateurs à des objets avant leur création complète. Les initializers sont progressivement dépréciés au profit des webhooks mutating.

### Exemples de contrôleurs d'admission courants dans Kubernetes

1. **NamespaceLifecycle** : Gère la création et la suppression des namespaces, en empêchant la création de ressources dans des namespaces en cours de suppression.

2. **ResourceQuota** : Assure que les namespaces ne dépassent pas les quotas de ressources alloués (CPU, mémoire, etc.).

3. **LimitRanger** : Applique des limites et des demandes par défaut pour les ressources si elles ne sont pas spécifiées par l'utilisateur.

4. **ServiceAccount** : Assure que les pods ont un compte de service associé, permettant ainsi une gestion fine des permissions et de la sécurité.

5. **PodSecurityPolicy** : Gère la politique de sécurité des pods, contrôlant ce qu'un pod peut et ne peut pas faire, par exemple, l'utilisation de volumes hostPath ou de privilèges escaladés.

6. **NodeRestriction** : Restreint les modifications que les kubelets (agents de nœuds) peuvent apporter aux ressources des nœuds et des pods.

### Fonctionnement général

1. **Intercept Request** : Les contrôleurs d'admission interceptent les requêtes qui arrivent au serveur d'API avant qu'elles ne soient persistées dans etcd (la base de données clé-valeur de Kubernetes).

2. **Evaluate Request** : Ils évaluent la requête en fonction des règles et politiques définies. Cela peut inclure la vérification des quotas, la validation de la sécurité, et l'application de valeurs par défaut.

3. **Modify or Reject Request** : En fonction de l'évaluation, le contrôleur peut modifier la requête (mutation) ou la rejeter (validation).

4. **Pass to Next Stage** : Si la requête est acceptée, elle est passée aux étapes suivantes pour le traitement par le serveur d'API et, éventuellement, persistée dans etcd.

En résumé, les admission controllers jouent un rôle crucial dans la sécurisation et la gestion des ressources dans les environnements de conteneurs, en particulier dans Kubernetes. Ils assurent que les requêtes sont conformes aux politiques de l'organisation avant d'être acceptées et traitées par le système.