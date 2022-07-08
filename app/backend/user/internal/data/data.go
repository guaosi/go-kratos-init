package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"maniverse/app/backend/user/internal/conf"

	"reflect"
	"time"

	// init mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewDB, NewUserRepo)

const (
	TxKey = "tx"
)

// Data .
type Data struct {
	db  *gorm.DB
	log *log.Helper
}

func NewDB(confApp *conf.App, confData *conf.Data, logger log.Logger) *gorm.DB {
	log := log.NewHelper(log.With(logger, "module", "user-service/data/gorm"))
	l := NewLogger(log)
	if confApp.Debug {
		l.LogLevel = gormlogger.Silent
	} else {
		l.LogLevel = gormlogger.Info
	}
	db, err := gorm.Open(mysql.Open(confData.Database.Source), &gorm.Config{
		Logger:                                   l,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}

	return db
}

// NewData .
func NewData(db *gorm.DB, logger log.Logger) (*Data, func(), error) {
	log := log.NewHelper(log.With(logger, "module", "user-service/data"))

	d := &Data{
		db:  db,
		log: log,
	}
	// 一些初始化操作
	d.init()

	return d, func() {}, nil
}

func (d *Data) init() {
	// 初始化表
	d.dbInitTables()
}
func (d *Data) Close() {
	// 关闭数据库链接
	d.dbClose()
}

func (d *Data) dbClose() {
	db, err := d.db.DB()
	if err == nil {
		db.Close()
	}
}

func (d *Data) dbInitTables() {
	initTables(d.db)
}

func (d *Data) DBTxBegin() *gorm.DB {
	return d.db.Begin()
}
func (d *Data) GetDBTxContext() context.Context {
	return context.WithValue(context.TODO(), TxKey, d.db)
}
func (d *Data) getDB(ctx context.Context) *gorm.DB {
	ctxValue := ctx.Value(TxKey)
	if ctxValue != nil {
		tx, ok := ctxValue.(*gorm.DB)
		if !ok {
			log.Errorf("unexpect context value type: %s \n", reflect.TypeOf(tx))
			return d.db.WithContext(ctx)
		}
		return tx
	}
	return d.db.WithContext(ctx)
}

type logger struct {
	log                       *log.Helper
	SlowThreshold             time.Duration
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
	LogLevel                  gormlogger.LogLevel
	SourceField               string
}

func NewLogger(log *log.Helper) *logger {
	return &logger{
		log:                       log,
		SkipCallerLookup:          false,
		IgnoreRecordNotFoundError: false,
		LogLevel:                  gormlogger.Warn,
		SlowThreshold:             2 * time.Second,
		SourceField:               "file",
	}
}

func (l *logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return &logger{
		log:                       l.log,
		SlowThreshold:             l.SlowThreshold,
		LogLevel:                  level,
		SkipCallerLookup:          l.SkipCallerLookup,
		IgnoreRecordNotFoundError: l.IgnoreRecordNotFoundError,
		SourceField:               "file",
	}
}
func (l *logger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Info {
		return
	}
	l.log.Debugf(str, args...)
}

func (l *logger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Warn {
		return
	}
	l.log.Warnf(str, args...)
}

func (l *logger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Error {
		return
	}
	l.log.Errorf(str, args...)
}

func (l *logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sqlString, _ := fc()
	fields := map[string]string{}
	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.IgnoreRecordNotFoundError) {
		fields["error"] = err.Error()
		l.log.WithContext(ctx).Errorf(fmt.Sprintf("%+v,%s [%s]", fields, sqlString, elapsed))
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.log.WithContext(ctx).Warnf(fmt.Sprintf("%+v,%s [%s]", fields, sqlString, elapsed))
		return
	}
	l.log.WithContext(ctx).Debugf(fmt.Sprintf("%+v,%s [%s]", fields, sqlString, elapsed))
}
