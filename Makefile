help:	## Message d'aide.
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

install-deps:	## Installe les dépendances.
	@echo "Installation des dépendances..."
	go install golang.org/x/tools/cmd/goimports@latest
	go install mvdan.cc/gofumpt@latest

format:	## Formate le code.
	@echo "Formatage..."
	gofmt -s -w .
	gofumpt -l -w .
	goimports -w .
	go mod tidy
