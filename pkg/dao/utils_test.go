package dao_test

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/in-rich/uservice-linkedin-data/config"
	"github.com/in-rich/uservice-linkedin-data/migrations"
	_ "github.com/in-rich/uservice-linkedin-data/migrations"
	"github.com/in-rich/uservice-linkedin-data/pkg/entities"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"google.golang.org/api/iterator"
	"net/http"
	"time"
)

func OpenDB() *bun.DB {
	dsn := "postgres://test:test@localhost:5432/test?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	err := db.Ping()
	for i := 0; i < 10 && err != nil; i++ {
		time.Sleep(1 * time.Second)
		err = db.Ping()
	}

	// Just in case something went wrong on latest run.
	ClearDB(db)

	if err := migrations.Migrate(db); err != nil {
		panic(err)
	}

	return db
}

func ClearDB(db *bun.DB) {
	if _, err := db.ExecContext(context.TODO(), "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"); err != nil {
		panic(err)
	}

	if _, err := db.ExecContext(context.TODO(), "GRANT ALL ON SCHEMA public TO public;"); err != nil {
		panic(err)
	}
	if _, err := db.ExecContext(context.TODO(), "GRANT ALL ON SCHEMA public TO test;"); err != nil {
		panic(err)
	}
}

func CloseDB(db *bun.DB) {
	ClearDB(db)

	if err := db.Close(); err != nil {
		panic(err)
	}
}

func BeginTX[T any](db bun.IDB, fixtures []T) bun.Tx {
	tx, err := db.BeginTx(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	for _, fixture := range fixtures {
		_, err := tx.NewInsert().Model(fixture).Exec(context.TODO())
		if err != nil {
			panic(err)
		}
	}

	return tx
}

func RollbackTX(tx bun.Tx) {
	_ = tx.Rollback()
}

func NewStorage(bucketName string, fixtures map[string][]byte) *storage.BucketHandle {
	bucket, err := config.StorageClient.Bucket(bucketName)
	if err != nil {
		panic(err)
	}

	for vanityName, profilePic := range fixtures {
		writer := bucket.Object(vanityName).NewWriter(context.TODO())
		if _, err := writer.Write(profilePic); err != nil {
			_ = writer.Close()
			panic(err)
		}
		_ = writer.Close()
	}

	return bucket
}

func StorageURI(id string) string {
	const baseURL = "http://127.0.0.1:1251/download/storage/v1/b/user-profile-pictures-test/o"
	return fmt.Sprintf("%s/%s?", baseURL, id)
}

func DownloadBase64(url string) string {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		strBody := new(bytes.Buffer)
		_, _ = strBody.ReadFrom(resp.Body)
		panic(fmt.Sprintf("unexpected status code %d: %s", resp.StatusCode, strBody))
	}

	body := new(bytes.Buffer)
	if _, err = body.ReadFrom(resp.Body); err != nil {
		panic(err)
	}

	return body.String()
}

func ClearStorage(bucket *storage.BucketHandle) {
	iter := bucket.Objects(context.TODO(), nil)

	for {
		attrs, err := iter.Next()

		if errors.Is(err, iterator.Done) {
			break
		}

		if err != nil {
			panic(err)
		}

		if err := bucket.Object(attrs.Name).Delete(context.TODO()); err != nil {
			panic(err)
		}
	}
}

func UsersCompare(a, b *entities.User) string {
	return cmp.Diff(a, b, cmpopts.IgnoreFields(entities.User{}, "UpdatedAt"))
}

func CompaniesCompare(a, b *entities.Company) string {
	return cmp.Diff(a, b, cmpopts.IgnoreFields(entities.Company{}, "UpdatedAt"))
}

func UsersCompareAll(a, b []*entities.User) string {
	return cmp.Diff(a, b, cmpopts.IgnoreFields(entities.User{}, "UpdatedAt"))
}

func CompaniesCompareAll(a, b []*entities.Company) string {
	return cmp.Diff(a, b, cmpopts.IgnoreFields(entities.Company{}, "UpdatedAt"))
}
