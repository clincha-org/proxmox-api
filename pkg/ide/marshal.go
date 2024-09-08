package ide

import (
	"fmt"
	"log/slog"
	"strings"
)

func Unmarshal(id int64, data string, storage *InternalDataStorage) error {
	if data == "" {
		return nil
	}
	commaSeparated := strings.Split(data, ",")

	slog.Debug("ide-unmarshal", "method", "Unmarshal", "data", data)

	storage.ID = id

	storage.Storage = strings.Split(commaSeparated[0], ":")[0]
	storage.Path = &strings.Split(commaSeparated[0], ":")[1]

	for _, value := range commaSeparated[1:] {
		keyValue := strings.Split(value, "=")
		switch keyValue[0] {
		case "media":
			storage.Media = &keyValue[1]
		case "size":
			storage.Size = &keyValue[1]
		}
	}
	return nil
}

func Marshal(storage *InternalDataStorage) (string, error) {

	if storage == nil {
		return "", fmt.Errorf("cannot marshal into nil InternalDataStorage object")
	}

	if storage.ID < 0 || storage.ID > 3 {
		return "", fmt.Errorf("invalid ID for IDE device: %v", storage.ID)
	}

	if storage.Storage == "" {
		return "", fmt.Errorf("storage is required for IDE device: %v", storage.ID)
	}

	var data string
	// Handle special syntax STORAGE_ID:SIZE_IN_GiB to allocate a new volume. See Proxmox API documentation.
	if storage.Path == nil && storage.Size != nil && *storage.Size != "" {
		return storage.Storage + ":" + *storage.Size, nil
	}

	data = storage.Storage + ":" + *storage.Path

	if storage.Media != nil && *storage.Media != "" {
		data += ",media=" + *storage.Media
	}
	if storage.Size != nil && *storage.Size != "" {
		data += ",size=" + *storage.Size
	}

	return data, nil
}
