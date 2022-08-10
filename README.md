# Nemo

Site web "Nemo", un engin de recherche d'article académique à partir de mot clef et de segment d'article.

# Support
Linux, non tester sur Windows (par contre avec Docker, il n'y a pas de problème)

# Installation et Dépendance

Il suffit de télécharger:

- [Docker](https://docs.docker.com/get-docker/)
- Make (pas obligatoire, facilite l'exécution de commandes)

Il faut exécuter [le script d'initialisation](https://github.com/erudit-recommandation/initialisation-service) afin de configurer les services pour une première fois.


# Usage

Afin d'exécuter localement avec `docker-compose` il suffit d'exécuter la commande `make run-docker`. Il est également possible à des fins de développement d’exécuter seulement la base de données et le service d'analyse de texte en exécutant la commande `make run-docker-debug`, pour plus d'information voir le répertoire [nemo](https://github.com/erudit-recommandation/Nemo/blob/main/nemo/README.md).