package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func NewGCSClient() (*storage.Client, error) {
	ctx := context.Background()
	return storage.NewClient(ctx)
}

func ListGCSFileNames(client *storage.Client, c Configuration, e ExportConfig) ([]string, error) {
	delim := "/"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	it := client.Bucket(c.GcsBucket).Objects(ctx, &storage.Query{
		Prefix:    c.GcsPrefix + "/" + e.Table + "/",
		Delimiter: delim,
	})
	names := []string{}
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []string{}, fmt.Errorf("Bucket(%q).Objects(): %v", c.GcsBucket, err)
		}
		names = append(names, attrs.Name)
	}
	return names, nil
}

func GetStorageLastIdFromFileNames(fileNames []string) (int64, error) {
	lastId := int64(0)
	for _, fileName := range fileNames {
		names := strings.Split(fileName, "/")
		name := names[len(names)-1]
		id, err := strconv.ParseInt(name, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("Error parsing Storage fileName %s to int64 for id", fileName)
		}
		if lastId < id {
			lastId = id
		}
	}
	return lastId, nil
}
