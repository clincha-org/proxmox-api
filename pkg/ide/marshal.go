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

	slog.Debug("ide-unmarshal", "method", "Unmarshal", "data", data)

	storage.ID = id

	IDEFields := strings.Split(data, ",")

	if IDEFields[0] != "none" {
		storage.Storage = &strings.Split(IDEFields[0], ":")[0]
		storage.Path = &strings.Split(IDEFields[0], ":")[1]
	}

	for _, value := range IDEFields[1:] {
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

	slog.Debug("ide-marshal", "method", "Marshal", "storage", storage)

	if storage == nil {
		return "", fmt.Errorf("cannot marshal into nil InternalDataStorage object")
	}

	if storage.ID < 0 || storage.ID > 3 {
		return "", fmt.Errorf("invalid ID for IDE device: %v", storage.ID)
	}

	// Handle special syntax STORAGE_ID:SIZE_IN_GiB to allocate a new volume. See Proxmox API documentation.
	if storage.Path == nil && storage.Storage != nil && *storage.Storage != "" && storage.Size != nil && *storage.Size != "" {
		// Remove the trailing "G" from the size
		*storage.Size = strings.TrimSuffix(*storage.Size, "G")

		slog.Debug("ide-marshal", "method", "Marshal", "new volume", *storage.Storage+":"+*storage.Size)

		return *storage.Storage + ":" + *storage.Size, nil
	}

	// Handle empty storage media
	if storage.Path == nil && storage.Storage == nil && storage.Media != nil && *storage.Media != "" {
		slog.Debug("ide-marshal", "method", "Marshal", "empty storage media", *storage.Media)

		return "none,media=" + *storage.Media, nil
	}

	data := *storage.Storage + ":" + *storage.Path

	if storage.Media != nil && *storage.Media != "" {
		data += ",media=" + *storage.Media
	}
	if storage.Size != nil && *storage.Size != "" {
		data += ",size=" + *storage.Size
	}

	return data, nil
}
