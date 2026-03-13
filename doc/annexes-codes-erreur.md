# Codes d'erreur de reference

[< Retour au sommaire](README.md)

---

| Code | Type | Exemples |
|------|------|----------|
| 200 | Succes | Attempt reussi, mise a jour OK |
| 201 | Cree | Nouveau joueur, donjon, run, listing |
| 400 | Validation | Champs manquants, format JSON invalide |
| 404 | Non trouve | Donjon inexistant, joueur inconnu |
| 409 | Conflit metier | `NOT_IN_RANGE`, `WRONG_STEP_ORDER`, `INSUFFICIENT_GOLD`, `INSUFFICIENT_ITEMS`, `LISTING_NOT_ACTIVE` |
| 422 | Non traitable | Publier un donjon sans boss step |
| 500 | Erreur serveur | Erreur MongoDB, panique interne |
