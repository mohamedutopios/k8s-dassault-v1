### Étape 1 : Lister les PodSecurityPolicies
D'abord, vous pouvez lister toutes les `PodSecurityPolicy` disponibles dans votre cluster pour trouver celle dont vous avez besoin.

```bash
kubectl get podsecuritypolicy
```

### Étape 2 : Décrire une PodSecurityPolicy spécifique
Une fois que vous avez identifié le nom de la PSP que vous souhaitez examiner, utilisez la commande `kubectl describe` pour obtenir ses détails.

```bash
kubectl describe podsecuritypolicy <nom-de-la-psp>
```

### Exemple
Si votre PSP s'appelle `restrictive-psp`, la commande serait :

```bash
kubectl describe podsecuritypolicy restrictive-psp
```

### Exemple de sortie
Voici un exemple de sortie de cette commande :

```plaintext
Name:  restrictive-psp
Namespace:  
Labels:  <none>
Annotations:  <none>
API Version:  policy/v1beta1
Kind:  PodSecurityPolicy
Metadata:
  Creation Timestamp:  2023-06-08T00:00:00Z
  Resource Version:  123456
  Self Link:  /apis/policy/v1beta1/podsecuritypolicies/restrictive-psp
  UID:  12345678-1234-1234-1234-123456789012
Spec:
  Allow Privilege Escalation:  false
  Allowed Capabilities:  <none>
  Allowed Flex Volumes:  <none>
  Allowed Host