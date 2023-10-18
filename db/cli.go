package db

import (
	"context"
	"jtrans/db/models"
	"jtrans/db/query"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	FilePath = "./sqlite.db"
)

var (
	ctx = context.Background()
	db  *gorm.DB
)

func init() {
	var err error
	db, err = gorm.Open(sqlite.Open(FilePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}

	_ = db.AutoMigrate(&models.FileSyncTask{})
}

func GetMinOrder() int {
	f := query.Use(db).FileSyncTask
	file, _ := f.WithContext(ctx).Order(f.Order).First()
	if file != nil {
		return file.Order
	}
	return 0
}

func GetMaxOrder() int {
	f := query.Use(db).FileSyncTask
	file, _ := f.WithContext(ctx).Order(f.Order.Desc()).First()
	if file != nil {
		return file.Order
	}
	return 0
}

func GetAllOrderedByOrder(descend bool) []*models.FileSyncTask {
	var files []*models.FileSyncTask
	f := query.Use(db).FileSyncTask
	if descend {
		files, _ = f.WithContext(ctx).Order(f.Order.Desc()).Find()
	} else {
		files, _ = f.WithContext(ctx).Order(f.Order).Find()
	}
	return files
}

func Update(task *models.FileSyncTask) error {
	f := query.Use(db).FileSyncTask
	return f.WithContext(ctx).Save(task)
}

func Begin() *gorm.DB {
	return db.Begin()
}

func Transaction(f func(tx *query.Query) error) error {
	q := query.Use(db)
	return q.Transaction(f)
}

func Insert(task *models.FileSyncTask) error {
	f := query.Use(db).FileSyncTask
	return f.WithContext(ctx).Save(task)
}

func FindByPath(path string) *models.FileSyncTask {
	f := query.Use(db).FileSyncTask
	task, _ := f.WithContext(ctx).Where(f.FilePath.Eq(path)).First()
	return task
}

func FindIdleTasks() []*models.FileSyncTask {
	f := query.Use(db).FileSyncTask
	tasks, _ := f.WithContext(ctx).Where(f.State.Eq(models.Idle)).Order(f.Order).Find()
	return tasks
}

func FindIdleTasksWithSmallerOrder(order int) []*models.FileSyncTask {
	f := query.Use(db).FileSyncTask
	tasks, _ := f.WithContext(ctx).Where(f.State.Eq(models.Idle), f.Order.Lt(order)).Find()
	return tasks
}
