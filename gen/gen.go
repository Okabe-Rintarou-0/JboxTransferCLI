package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gen"
	"gorm.io/gorm"
	"jtrans/db/models"
)

func main() {
	var (
		err error
		db  *gorm.DB
	)
	db, err = gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})

	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	g := gen.NewGenerator(gen.Config{
		OutPath:       "./db/query",
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	// Use the above `*gorm.DB` instance to initialize the generator,
	// which is required to generate structs from db when using `GenerateModel/GenerateModelAs`
	g.UseDB(db)

	// Generate default DAO interface for those specified structs
	g.ApplyBasic(models.FileSyncTask{})

	// Execute the generator
	g.Execute()
}
