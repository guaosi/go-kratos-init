package data

import (
	"database/sql/driver"
	"fmt"
	"gorm.io/gorm"
	"sort"
	"strconv"
	"strings"
	"time"
)

// 表名，根据字母排序
const (
	TableUser = "user"
)

// 基本模型的定义
type BaseModel struct {
	ID        uint64    `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `gorm:"not null; default: current_timestamp(3);comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null; default: current_timestamp(3) on update current_timestamp(3);comment:更新时间" json:"updated_at"`
}
type GormUintList []uint

func (g GormUintList) Len() int           { return len(g) }
func (g GormUintList) Less(i, j int) bool { return g[i] < g[j] }
func (g GormUintList) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }
func (g GormUintList) Value() (driver.Value, error) {
	// 排序
	sort.Sort(g)
	str := ""
	first := true
	for _, v := range g {
		if first {
			str = fmt.Sprintf("%d", v)
			first = false
		} else {
			str += fmt.Sprintf(",%d", v)
		}
	}
	return []byte(str), nil
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (g *GormUintList) Scan(value interface{}) error {
	str := string(value.([]byte))
	if str != "" {
		strSlice := strings.Split(string(value.([]byte)), ",")
		data := make([]uint, 0, len(strSlice))
		for _, v := range strSlice {
			tmp, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return err
			}
			data = append(data, uint(tmp))
		}
		*g = data
	}
	return nil
}

type User struct {
	BaseModel
	UUID     string `gorm:"type:varchar(36);unique;not null;default:''; comment:uuid" json:"uuid"`
	Phone    string `gorm:"type:varchar(11);unique;not null;default:''; comment:手机号" json:"phone"`
	Password string `gorm:"type:varchar(100);not null; comment:密码" json:"password"`
	NickName string `gorm:"type:varchar(20);not null; default:'';comment:昵称" json:"nick_name"`
}

func (User) TableName() string {
	return TableUser
}

func initTables(db *gorm.DB) {
	db.AutoMigrate(
		new(User),
	)
}
