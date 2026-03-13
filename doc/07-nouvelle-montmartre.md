# La Quete du Donjon de Montmartre

[< Retour au sommaire](README.md)

*Une courte histoire illustrant une session de jeu complete avec un Game Master et trois joueurs.*

---

## Chapitre 1 - Le Maitre du Jeu

Theo s'installa a la terrasse du cafe, son ordinateur pose entre deux tasses vides. Il avait passe la semaine a reperer les lieux, a mesurer les distances, a imaginer les creatures. Ce soir, le Donjon de Montmartre ouvrirait ses portes.

Il ouvrit son terminal.

```
POST /v1/players
{ "display_name": "MaitreSombre" }
-> { "id": "mj-theo" }
```

Trois items d'abord. Il fallait que les recompenses en vaillent la peine :

```
POST /v1/items
{ "name": "Bouclier du Funiculaire", "type": "artifact", "rarity": "rare",
  "stats": {"defense": 30}, "tradable": true, "baseValue": 350 }
-> { "id": "item-bouclier" }

POST /v1/items
{ "name": "Fiole d'Absinthe Magique", "type": "consumable", "rarity": "uncommon",
  "stats": {"heal": 50}, "tradable": true, "baseValue": 80 }
-> { "id": "item-absinthe" }

POST /v1/items
{ "name": "Lame de la Butte", "type": "weapon", "rarity": "legendary",
  "stats": {"attack": 75, "crit": 15}, "tradable": true, "baseValue": 1200 }
-> { "id": "item-lame" }
```

Puis le donjon lui-meme :

```
POST /v1/mj/dungeons
{ "title": "Le Donjon de Montmartre",
  "description": "Trois gardiens protegent la Butte. Seuls les braves les vaincront.",
  "createdBy": "mj-theo",
  "area": { "name": "Montmartre, Paris 18e" } }
-> { "id": "dng-montmartre", "status": "draft" }
```

Trois boss. Trois lieux. Trois defis.

```
POST /v1/mj/dungeons/dng-montmartre/steps
{ "name": "Le Minotaure du Metro Abbesses",
  "location": { "lat": 48.8843, "lon": 2.3385, "radiusMeters": 40 },
  "zoneDescription": "A la sortie du metro, sous la fresque en mosaique",
  "difficulty": 3, "goldReward": 100,
  "lootTable": [
    { "itemId": "item-absinthe", "dropRate": 0.7, "minQty": 1, "maxQty": 2 }
  ] }
-> { "id": "step-abbesses", "order": 1 }

POST /v1/mj/dungeons/dng-montmartre/steps
{ "name": "La Gargouille du Sacre-Coeur",
  "location": { "lat": 48.8867, "lon": 2.3431, "radiusMeters": 50 },
  "zoneDescription": "Au pied de l'escalier monumental",
  "difficulty": 6, "goldReward": 200,
  "lootTable": [
    { "itemId": "item-bouclier", "dropRate": 0.4, "minQty": 1, "maxQty": 1 },
    { "itemId": "item-absinthe", "dropRate": 0.5, "minQty": 1, "maxQty": 3 }
  ] }
-> { "id": "step-sacrecoeur", "order": 2 }

POST /v1/mj/dungeons/dng-montmartre/steps
{ "name": "La Sorciere de la Place du Tertre",
  "location": { "lat": 48.8863, "lon": 2.3409, "radiusMeters": 25 },
  "zoneDescription": "Cachee parmi les chevalets des portraitistes",
  "difficulty": 9, "goldReward": 500,
  "lootTable": [
    { "itemId": "item-lame", "dropRate": 0.08, "minQty": 1, "maxQty": 1 },
    { "itemId": "item-bouclier", "dropRate": 0.3, "minQty": 1, "maxQty": 1 }
  ] }
-> { "id": "step-tertre", "order": 3 }
```

Theo relu sa creation. Trois boss, difficulte croissante, une legendaire a 8% de drop sur le dernier. Parfait.

```
POST /v1/mj/dungeons/dng-montmartre/publish
-> 200 "dungeon published"
```

Il envoya le lien dans le groupe. **"Montmartre vous attend. Qui ose ?"**

---

## Chapitre 2 - Les Trois Aventuriers

Samedi, 14h. Trois amis se retrouverent a la sortie du metro Anvers.

**Elise**, la strategiste. Elle avait deja etudie la carte.
**Karim**, le fonceur. Il n'avait meme pas lu la description.
**Mei**, la collectionneuse. Elle ne revait que de la Lame de la Butte.

```
POST /v1/players { "display_name": "Elise_la_Sage" }    -> { "id": "p-elise" }
POST /v1/players { "display_name": "KarimLeBrave" }      -> { "id": "p-karim" }
POST /v1/players { "display_name": "Mei_Chasseuse" }     -> { "id": "p-mei" }
```

Ils consulterent les donjons disponibles :

```
GET /v1/dungeons
-> [{ "id": "dng-montmartre", "title": "Le Donjon de Montmartre", ... }]

GET /v1/dungeons/dng-montmartre
-> { ..., "bossSteps": [
     { "name": "Le Minotaure du Metro Abbesses", "order": 1, "difficulty": 3 },
     { "name": "La Gargouille du Sacre-Coeur", "order": 2, "difficulty": 6 },
     { "name": "La Sorciere de la Place du Tertre", "order": 3, "difficulty": 9 }
   ] }
```

"Trois boss," murmura Elise. "On commence par Abbesses, c'est juste en haut."

Chacun lanca son run :

```
POST /v1/runs { "dungeonId": "dng-montmartre", "playerId": "p-elise" }
-> { "id": "run-elise", "state": "active", "currentStep": 1 }

POST /v1/runs { "dungeonId": "dng-montmartre", "playerId": "p-karim" }
-> { "id": "run-karim", "state": "active", "currentStep": 1 }

POST /v1/runs { "dungeonId": "dng-montmartre", "playerId": "p-mei" }
-> { "id": "run-mei", "state": "active", "currentStep": 1 }
```

---

## Chapitre 3 - Le Minotaure du Metro Abbesses

Ils remonterent la rue Lepic. A la sortie du metro Abbesses, sous la fresque en mosaique, Karim degaina son telephone le premier.

```
POST /v1/runs/run-karim/steps/step-abbesses/attempt
{ "lat": 48.8844, "lon": 2.3386 }
```

Le serveur calcula. Distance : 14 metres. Rayon : 40 metres.

```
-> 200 { "success": true, "rewards": { "gold": 100,
     "items": [{"itemId": "item-absinthe", "qty": 2}] }, "runCompleted": false }
```

"Cent pieces d'or et deux fioles !" cria Karim. Les deux filles tenterent a leur tour.

Elise, depuis le meme trottoir :
```
POST /v1/runs/run-elise/steps/step-abbesses/attempt
{ "lat": 48.8842, "lon": 2.3384 }
-> 200 { "rewards": { "gold": 100, "items": [{"itemId": "item-absinthe", "qty": 1}] } }
```

Mei, un peu plus loin, testant les limites :
```
POST /v1/runs/run-mei/steps/step-abbesses/attempt
{ "lat": 48.8840, "lon": 2.3383 }
-> 200 { "rewards": { "gold": 100, "items": [] } }
```

Pas de loot pour Mei cette fois. "Le RNG me deteste," soupira-t-elle.

---

## Chapitre 4 - L'Erreur de Karim

Euphorique, Karim voulut bruler les etapes. Sans attendre les autres, il courut vers la Place du Tertre et tenta directement le boss final.

```
POST /v1/runs/run-karim/steps/step-tertre/attempt
{ "lat": 48.8864, "lon": 2.3410 }
-> 409 { "message": "WRONG_STEP_ORDER" }
```

"Tu dois d'abord battre la Gargouille," lui expliqua Elise par message. "C'est l'etape 2, tu es a l'etape 2, mais ce n'est pas le bon step ID."

Karim redescendit vers le Sacre-Coeur en grommelant.

---

## Chapitre 5 - La Gargouille du Sacre-Coeur

Au pied de l'escalier monumental, les trois se regroupent. La vue sur Paris etait magnifique, mais la Gargouille les attendait.

Mei tenta sa chance depuis le haut des marches - trop loin :

```
POST /v1/runs/run-mei/steps/step-sacrecoeur/attempt
{ "lat": 48.8872, "lon": 2.3435 }
-> 409 { "message": "NOT_IN_RANGE" }
```

56 metres. Le rayon etait de 50.

Elle descendit quelques marches et reessaya :

```
POST /v1/runs/run-mei/steps/step-sacrecoeur/attempt
{ "lat": 48.8868, "lon": 2.3432 }
-> 200 { "rewards": { "gold": 200,
     "items": [{"itemId": "item-bouclier", "qty": 1}] }, "runCompleted": false }
```

"Le Bouclier du Funiculaire !" Mei le serra precieusement dans son inventaire virtuel.

Karim et Elise vainquirent la Gargouille a leur tour. Karim obtint trois Fioles d'Absinthe, Elise un Bouclier.

---

## Chapitre 6 - La Sorciere de la Place du Tertre

La Place du Tertre fourmillait de touristes et de peintres. Quelque part entre les chevalets, la Sorciere attendait.

Le rayon etait serre : 25 metres seulement. Il fallait etre au coeur de la place.

Elise, pragmatique, se placa en plein centre :

```
POST /v1/runs/run-elise/steps/step-tertre/attempt
{ "lat": 48.8863, "lon": 2.3409 }
-> 200 { "rewards": { "gold": 500, "items": [] }, "runCompleted": true }
```

500 pieces d'or, mais pas de Lame. Run termine.

Karim, depuis le meme spot :
```
POST /v1/runs/run-karim/steps/step-tertre/attempt
{ "lat": 48.8862, "lon": 2.3408 }
-> 200 { "rewards": { "gold": 500,
     "items": [{"itemId": "item-bouclier", "qty": 1}] }, "runCompleted": true }
```

Un bouclier de plus. Pas de Lame non plus.

Mei, le coeur battant. 8% de chance.

```
POST /v1/runs/run-mei/steps/step-tertre/attempt
{ "lat": 48.8863, "lon": 2.3410 }
-> 200 { "rewards": { "gold": 500,
     "items": [{"itemId": "item-lame", "qty": 1}] }, "runCompleted": true }
```

**La Lame de la Butte.**

Mei poussa un cri que les portraitistes de la place n'oublierent pas de sitot.

---

## Chapitre 7 - Le Marche

De retour au cafe, les trois aventuriers comparerent leur butin :

```
GET /v1/inventory?playerId=p-elise
-> { "items": [{"itemId": "item-absinthe", "qty": 1}, {"itemId": "item-bouclier", "qty": 1}] }

GET /v1/inventory?playerId=p-karim
-> { "items": [{"itemId": "item-absinthe", "qty": 5}, {"itemId": "item-bouclier", "qty": 1}] }

GET /v1/inventory?playerId=p-mei
-> { "items": [{"itemId": "item-bouclier", "qty": 1}, {"itemId": "item-lame", "qty": 1}] }
```

Karim, assis sur ses 5 fioles, vit une opportunite :

```
POST /v1/auction/listings
{ "sellerId": "p-karim", "itemId": "item-absinthe", "qty": 3, "pricePerUnit": 60 }
-> 201 { "id": "list-001", "status": "active" }
```

Elise, qui n'avait qu'une fiole, acheta immediatement :

```
POST /v1/auction/listings/list-001/buy
{ "buyerId": "p-elise", "qty": 2 }
-> 200 "purchase successful"
```

Elise : -120 gold, +2 fioles. Karim : +120 gold.

Mei, elle, n'avait pas l'intention de vendre la Lame. Mais elle consulta le leaderboard :

```
GET /v1/leaderboard?type=completions&limit=5
-> [
  { "playerId": "p-elise", "displayName": "Elise_la_Sage", "score": 1 },
  { "playerId": "p-karim", "displayName": "KarimLeBrave", "score": 1 },
  { "playerId": "p-mei", "displayName": "Mei_Chasseuse", "score": 1 }
]

GET /v1/leaderboard?type=speed&dungeonId=dng-montmartre&limit=3
-> [
  { "playerId": "p-karim", "displayName": "KarimLeBrave", "score": 2847.5 },
  { "playerId": "p-elise", "displayName": "Elise_la_Sage", "score": 3012.1 },
  { "playerId": "p-mei", "displayName": "Mei_Chasseuse", "score": 3198.7 }
]
```

Karim, malgre son erreur de parcours, avait ete le plus rapide. Il leva sa biere.

---

## Epilogue

Theo recut les resultats sur son tableau de bord. Trois runs completes, des items distribues, le marche qui commencait a vivre. Il ouvrit son terminal une derniere fois.

Il avait deja une idee pour le prochain donjon. Les Catacombes de Paris, cette fois. Plus sombre. Plus profond. Plus dangereux.

Mais ca, c'est une autre histoire.

---

*Fin.*
