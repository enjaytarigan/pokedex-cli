BIN_DIR := ./bin

.PHONY: pokedex/start
 pokedex/start:
	go build -o ${BIN_DIR} . && ${BIN_DIR}/pokedexcli