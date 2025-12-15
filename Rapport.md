# AI30 - Simulation de Black Friday

## Description du projet

Le projet que nous avons décidé de développer est une simulation de supermarché dans une période de Black Friday ou de solde. L’objectif principal de cette simulation est de reproduire le comportement des clients dans un contexte de forte affluence.

Pour représenter un supermarché nous avons adopté une vue du dessus. Sur notre interface nous avons représenté les éléments suivants : 
- Des rayons : qui contiennent des produits que les agents vont aller rechercher.
- Des caisses : pour que notre supermarché puisse enregistrer les produits achetés par nos agents.
- Des portes : où les agents apparaissent à intervalle régulier et où ils vont une fois qu'ils ont payé leurs produits.

Les agents sont des clients qui vont collecter jusqu'à x objets simultanément et avoir plusieurs comportements :
- Collaboratifs : ces types d'agents vont aller chercher les produits dont ils ont besoin.
- Égoïstes : agissent comme les agents collaboratifs, s'ils recherchent un produit qu'un agent à proximité possède, ils peuvent voler ses produits.
- Rancunier : s'il se fait voler, il va essayer de voler quelqu'un d'autre pour récupérer son produit.
- Agent de sécurité : ces types d'agents ont la particularité de se déplacer aléatoirement dans le supermarché et permettent de réduire le vol inter-agents autour d'eux.

Dans le but de mieux comprendre les actions qui se déroulent dans notre simulation, en cliquant sur des éléments sur la carte (agents, rayons, ...), des informations liées à celui-ci s'afficheront dans un HUD (Heads up display), de plus, lorsqu'un agent récupère un objet, il est en surbrillance en vert.

## Question à répondre

La question que nous cherchons à répondre à l'aide de notre simulation est : quel est le meilleur agencement de magasin pour réaliser le plus de profit, le plus rapidement possible?

## Indicateurs de réponse

Les indicateurs utilisés pour répondre à cette question sont les suivants :
- Le bénéfice du supermarché sur un agencement : combien le supermarché peut bénéficier dans son magasin à partir d'un stock donné.
- Le temps que prennent les agents à faire leurs achats : le meilleur agencement passe par le temps que prennent les agents pour trouver rapidement leurs produits et les payer rapidement avant de partir.
- Le nombre de collisions inter-agents : afin de considérer le bien-être des agents, le nombre de collisions est un nombre qui affecte négativement l'agencement du magasin (couloirs trop serrés, accès difficiles aux rayons ou aux caisses).

## Architecture employée

Concernant l'architecture que nous employons, nous utilisons la bibliothèque graphique Ebitengine pour notre fenêtre.

### Architecture des packages internes

**pkg/constants** : centralise toutes les constantes de notre application (taille des cellules, coefficient des forces sociales, taille de la vision des agents, ...).

**pkg/graphics** : s'occupe de l'interface utilisateur, dont les interactions avec l'interface, les animations des agents, gestion des images, ...

**pkg/hud** : gère l'affichage du HUD sur la fenêtre.

**pkg/logger** : permet d'afficher des informations pour debugger notre application.

**pkg/map** : s'occupe de la carte de la simulation.

**pkg/pathfinding** : gère la recherche de chemin des agents en utilisant l'algorithme A*.

**pkg/shopping** : génère et lit dans des fichiers, les listes d'achat pour nos agents

**pkg/simulation** : contient le noyau de notre simulation, comprenant la gestion des agents, des accès concurrents, collisions, visions...

**pkg/utils** : dossier utilitaire pour gérer des types partagés et des calculs mathématiques.

### Commandes

**cmd/blackfriday** :  lance notre simulation en lisant un fichier pour une carte à générer, une liste JSON pour les stocks et une liste JSON pour la liste d'achat des agents.
**cmd/generate_shopping_lists_from_stocks** : génère une liste d'achats pour tous les agents à partir des données du stock. Pouvant randre, si l'on le souhaite la liste d'achat déterministe au travers des différentes simulations
**cmd/map_generator** : génère aléatoirement des cartes comprenant une taille, un nombre de portes, caisses, rayons, murs donnés.

## Paramètres

Les paramètres que nous proposons de changer sont les suivants. Ceux-ci sont à renseigner dans un fichier `.env`. Un fichier `.env.example` est dans le dépôt avec les dits paramètres :
- NB_AGENTS (le nombre d'agents)
- AGENT_SPEED (la vitesse de déplacement des agents)
- AGENT_MAX_SHOPPING_LIST (le nombre de produits maximum qu'un agent peut vouloir chercher)

Si ces paramètres ne sont pas renseignés, la simulation attribue des valeurs arbitraires.

## Comment lancer l'application

Pour lancer l'application, il faut exécuter les commandes suivantes :

#### Installation des dépendances
```
sudo apt install -y libx11-dev libxrandr-dev libxinerama-dev libxcursor-dev libxi-dev
sudo apt install -y libgl1-mesa-dev xorg-dev
```

#### Configuration
```
cp .env.example .env
# Modifier .env selon vos besoins (nombre d'agents, vitesse, etc.)

go run cmd/generate_shopping_lists_from_stocks/main.go
# Pour générer les listes d'achats pour les agents, dans le but de comparer plusieurs simulations avec les mêmes données.
```

#### Exécution
```
go mod download
go run cmd/blackfriday/main.go
```

#### Analyse des ventes
```
./stats/plot.sh
# Lance un script qui va afficher le graphique à partir des données extraites des simulations.
```