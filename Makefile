.PHONY: test test-integration tidy

# Lancer tous les tests unitaires
test:
	go test ./app/functions/... -v

# Lancer les tests d'intégration avec la base de données de test
test-integration:
	@echo "Lancement des tests d'intégration..."
	go test ./tests/... -v

# Nettoyer les dépendances
tidy:
	go mod tidy
