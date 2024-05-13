package indexer

import (
	"context"
	"encoding/hex"
	"net/http"

	"github.com/keep-starknet-strange/art-peace/backend/core"
)

func processUsernameClaimedEvent(event IndexerEvent, w http.ResponseWriter) {
	address := event.Event.Keys[1][2:]     // Remove 0x prefix
	usernameHex := event.Event.Data[0][2:] // Remove 0x prefix

	// Parse username hex as bytes encoded in utf-8
	decodedUsername, err := hex.DecodeString(usernameHex)
	if err != nil {
		PrintIndexerError("processUsernameClaimedEvent", "Error decoding username hex", address, usernameHex)
		return
	}
	// Trim off 0s at the start
	trimmedUsername := []byte{}
	trimming := true
	for _, b := range decodedUsername {
		if b == 0 && trimming {
			continue
		}
		trimming = false
		trimmedUsername = append(trimmedUsername, b)
	}
	username := string(trimmedUsername)

	// Set username in postgres
	_, err = core.ArtPeaceBackend.Databases.Postgres.Exec(context.Background(), "INSERT INTO Users (address, name) VALUES ($1, $2)", address, username)
	if err != nil {
		PrintIndexerError("processUsernameClaimedEvent", "Error inserting username into postgres", address, username)
		return
	}
}

func processUsernameChangedEvent(event IndexerEvent, w http.ResponseWriter) {
	address := event.Event.Keys[1][2:]     // Remove 0x prefix
	usernameHex := event.Event.Data[1][2:] // Remove 0x prefix

	// Parse username hex as bytes encoded in utf-8
	decodedUsername, err := hex.DecodeString(usernameHex)
	if err != nil {
		PrintIndexerError("processUsernameChangedEvent", "Error decoding username hex", address, usernameHex)
		return
	}
	// Trim off 0s at the start
	trimmedUsername := []byte{}
	trimming := true
	for _, b := range decodedUsername {
		if b == 0 && trimming {
			continue
		}
		trimming = false
		trimmedUsername = append(trimmedUsername, b)
	}
	username := string(trimmedUsername)

	// Set username in postgres
	_, err = core.ArtPeaceBackend.Databases.Postgres.Exec(context.Background(), "UPDATE Users SET name = $1 WHERE address = $2", username, address)
	if err != nil {
		PrintIndexerError("processUsernameChangedEvent", "Error updating username in postgres", address, username)
		return
	}
}
