# Comores Events 2027 — Maquette web

Cette maquette reprend la direction graphique de la capture fournie (vert profond, fond ivoire, titres éditoriaux, sections culture, billetterie, informations pratiques et FAQ).

## Ouvrir le site

Décompressez le ZIP et **gardez les fichiers ensemble**, puis double-cliquez sur `index.html`.

- `index.html` : structure du site
- `style.css` : couleurs, mise en page et responsive mobile
- `script.js` : interactions et simulation du choix des billets
- `assets/` : visuels issus de la référence graphique fournie

Pour utiliser un serveur local : `python -m http.server 8000`, puis ouvrez `http://localhost:8000`.

## Fonctionnalités testables

- Navigation vers les sections et menu mobile
- Choix d'un billet 1 jour (samedi/dimanche) ou d'un pass 2 jours
- Modification de la quantité et calcul instantané du total
- Fenêtre de récapitulatif sans transaction réelle
- Questions fréquentes dépliables

## Configuration et intégration FOODYFLY

Les tarifs du prototype sont configurables dans `EVENT_CONFIG` au début de `script.js`.

**Attention : ce prototype n'encaisse aucun paiement.** Le bouton « Poursuivre avec FOODYFLY » ouvre une fenêtre de démonstration. En production, il faudra appeler un backend sécurisé qui créera une session de paiement via l'API FOODYFLY, puis gérer la confirmation des commandes, l'émission des QR codes, la protection contre la réutilisation et les notifications.

Les images ont été recadrées depuis la capture fournie par l'utilisateur pour préserver l'ambiance graphique. Pour la production, utilisez des photographies haute définition disposant des autorisations nécessaires. Le chargement des polices Google nécessite une connexion internet ; des polices de secours sont prévues.
