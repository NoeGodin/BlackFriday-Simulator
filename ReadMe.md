# Simulateur de Black Friday

## Installation

**Linux (Ubuntu/Debian)**: Installer les dépendances système :
```bash
sudo apt install -y libx11-dev libxrandr-dev libxinerama-dev libxcursor-dev libxi-dev
sudo apt install -y libgl1-mesa-dev xorg-dev
```

**Exécution:**
```bash
go mod download
go run cmd/blackfriday/main.go
```

## Problématique

Point de vue du magasin : Comment faire le plus de ventes ?

## Fonctionnement

Sur l'interface graphique on peut mettre :
- des rayons (obstacles)
- des objets (objectif)
- une entrée / sortie

Les agents seraient les clients qui pourraient collecter jusqu'à x objets simultanément et avoir plusieurs comportements :
- collaboratifs : laisser l'objet à l'agent à une taille de distance sans objet autour de lui sauf s'il est aussi collaboratif
- compétitifs : va prendre le chemin optimal pour entrer et sortir avec 1 objet
- égoïstes : va voler jusqu'à 3 objets aux autres agents si pas d'autres agents avec objets va en chercher un directement
- rancunier : compétitif, s'il se fait voler, va voler quelqu'un d'autre que son voleur

On pourrait également avoir des agents qui réapprovisionneraient les rayons.
Voir comment on gère les quantités d'objets.

