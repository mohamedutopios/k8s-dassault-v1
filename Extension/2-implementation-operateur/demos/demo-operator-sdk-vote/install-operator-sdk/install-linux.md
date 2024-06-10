# Télécharger et extraire Go
wget https://golang.org/dl/go1.20.1.linux-amd64.tar.gz  # Remplacez par la version appropriée
sudo tar -C /usr/local -xzf go1.20.1.linux-amd64.tar.gz

# Configurer l'environnement
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
echo 'export GOPATH=$HOME/go' >> ~/.profile
echo 'export GOBIN=$GOPATH/bin' >> ~/.profile
echo 'export PATH=$PATH:$GOBIN' >> ~/.profile
source ~/.profile

# Vérifier l'installation de Go
go version

# Installer l'Operator SDK
go install github.com/operator-framework/operator-sdk/cmd/operator-sdk@latest

# Vérifier l'installation de l'Operator SDK
operator-sdk version
