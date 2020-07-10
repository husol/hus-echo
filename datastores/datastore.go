package datastores

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file" // golang-migrate
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func CreateConnection() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	name := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")

	uri := fmt.Sprintf(
		"%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, name,
	)
	db, err := gorm.Open("mysql", uri)
	if db != nil {
		envMaxIdleConns := os.Getenv("DB_MAX_IDLE_CONNS")
		envMaxOpenConns := os.Getenv("DB_MAX_OPEN_CONNS")
		maxIdleConns := 10
		maxOpenConns := 44
		if envMaxIdleConns != "" {
			numValue, err := strconv.Atoi(envMaxIdleConns)
			if err != nil {
				return db, err
			}
			maxIdleConns = numValue
		}
		if envMaxOpenConns != "" {
			numValue, err := strconv.Atoi(envMaxOpenConns)
			if err != nil {
				return db, err
			}
			maxOpenConns = numValue
		}
		db.DB().SetMaxIdleConns(maxIdleConns)
		db.DB().SetMaxOpenConns(maxOpenConns)
		db.DB().SetConnMaxLifetime(time.Minute * 5)
		db.LogMode(true)
	}
	return db, err
}

func Migrate(db *gorm.DB) {
	driver, err := mysql.WithInstance(db.DB(), &mysql.Config{
		MigrationsTable: "custom_schema_migrations",
	})

	if err != nil {
		log.Fatal(err)
	}

	var m *migrate.Migrate
	if os.Getenv("IS_DEBUG") == "true" {
		m, err = migrate.NewWithDatabaseInstance("file://./migrations", os.Getenv("DB_NAME"), driver)
	} else {
		m, err = migrate.NewWithDatabaseInstance("file://./migrations", os.Getenv("DB_NAME"), driver)
	}

	if err != nil {
		log.Fatal(err)
	}
	err = m.Up()

	if err != nil {
		if !strings.Contains(err.Error(), "no change") && err.Error() != "file does not exist" {
			log.Fatal(err)
		}
	}

	log.Println("Migrate done!")
}

type BatchDb struct {
	*gorm.DB
}

func (db *BatchDb) BatchInsert(objArr []interface{}) (int64, error) {
	// If there is no data, nothing to do.
	if len(objArr) == 0 {
		return 0, errors.New("insert a slice length of 0")
	}

	mainObj := objArr[0]
	mainScope := db.NewScope(mainObj)
	mainFields := mainScope.Fields()
	quoted := make([]string, 0, len(mainFields))
	for i := range mainFields {
		// If primary key has blank value (0 for int, "" for string, nil for interface ...), skip it.
		// If field is ignore field, skip it.
		if (mainFields[i].IsPrimaryKey && mainFields[i].IsBlank) || (mainFields[i].IsIgnored) {
			continue
		}
		quoted = append(quoted, mainScope.Quote(mainFields[i].DBName))
	}

	placeholdersArr := make([]string, 0, len(objArr))

	for _, obj := range objArr {
		scope := db.NewScope(obj)
		fields := scope.Fields()
		placeholders := make([]string, 0, len(fields))
		for i := range fields {
			if (fields[i].IsPrimaryKey && fields[i].IsBlank) || (fields[i].IsIgnored) {
				continue
			}
			var vars interface{}
			if (fields[i].Name == "CreatedAt" || fields[i].Name == "UpdatedAt") && fields[i].IsBlank {
				vars = gorm.NowFunc()
			} else {
				vars = fields[i].Field.Interface()
			}
			placeholders = append(placeholders, scope.AddToVars(vars))
		}
		placeholdersStr := "(" + strings.Join(placeholders, ", ") + ")"
		placeholdersArr = append(placeholdersArr, placeholdersStr)
		// add real variables for the replacement of placeholders' '?' letter later.
		mainScope.SQLVars = append(mainScope.SQLVars, scope.SQLVars...)
	}
	mainScope.Raw(fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
		mainScope.QuotedTableName(),
		strings.Join(quoted, ", "),
		strings.Join(placeholdersArr, ", "),
	))
	//Execute and Log
	if err := mainScope.Exec().DB().Error; err != nil {
		return 0, err
	}
	return mainScope.DB().RowsAffected, nil
}
