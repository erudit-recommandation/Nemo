# Nemo

Site web "Nemo", un engin de recherche d'article académique à partir de mot clef et de segment d'article.

# Support
Linux, non tester sur Windows (par contre avec Docker, il n'y a pas de problème)

# Installation et Dépendance

Il suffit de télécharger les dépendances suivantes: 

- [GO version 1.18](https://go.dev/dl/)
- [Docker](https://docs.docker.com/get-docker/)
- Make (pas obligatoire, facilite l'exécution de commandes)

  
# Usage

Le site web dépend de deux services, une base de données ArangoDB et [https://github.com/erudit-recommandation/text-analysis-service](text-analysis-service). Ces deux services peuvent être connectés de trois façons:
- Via l'internet 
- Via un réseau Docker-compose
- Via l'hôte local (localhost)

Il faut exécuter [le script d'initialisation](https://github.com/erudit-recommandation/initialisation-service) afin de configurer les services pour une première fois.

## Internet
Ce mode se sert du fichier `.env`, ce fichier doit être créé par l'usager en suivant le canvas de `.env_dev`. Il suffit de lancer la commande `make run-prod`.

L'utilité de ce mode est de soit déployer l'application dans un serveur ou encore pour utiliser l'application localement en se servant des données en ligne à des fins de test.

## Docker-compose
Ce mode se sert du fichier `.env_dev_docker` afin de configurer les services. Pour lancer ce mode il suffit de lancer à [la racine du répertoire](https://github.com/erudit-recommandation/Nemo) `make run run-docker`.

L'utilité de ce mode est de faire fonctionner localement l'application.
## Localhost
Ce mode se sert du fichier `.env_dev` afin de configurer les services. Il faut au préalable lancer les services à [la racine du répertoire](https://github.com/erudit-recommandation/Nemo) en lançant `make run-docker-debug`, par la suite il faut dans ce répertoire courant lancer la commande `make run`.

L'utilité de ce mode est de plus facilement faire du développement, pour lancer l'application localement, il est plus simple de servir de Docker-compose.