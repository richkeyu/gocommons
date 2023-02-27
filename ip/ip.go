package ip

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/richkeyu/gocommons/plog"
	"github.com/richkeyu/gocommons/server"
	"github.com/richkeyu/gocommons/util"

	"github.com/gin-gonic/gin"
	"github.com/ip2location/ip2location-go/v9"
)

var db *ip2location.DB
var backDb *ip2location.DB // 初始连接
var mutex sync.RWMutex
var downloadUrl = "https://download.ip2location.com/lite/IP2LOCATION-LITE-DB1.IPV6.BIN.ZIP"
var FilePath = "./config/ip2location/"
var FileName = "IP2LOCATION-LITE-DB1.IPV6.BIN"

func ToLocation(ip string) (LocationRecord, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	var result LocationRecord
	results, err := db.Get_all(ip)
	if err != nil {
		return result, err
	}
	return ip2locationRecordConvert(results), nil
}

// Init 初始化 定期更新IP库
func Init() {
	curDb, err := ip2location.OpenDB(FilePath + FileName)
	if err != nil {
		panic(fmt.Sprintf("ip2location init fail: %s", err))
	}
	db = curDb
	backDb = db

	// 更新进程
	go func() {
		ctx := server.NewContext(context.Background(), &gin.Context{})
		// 立即更新一次
		err = UpdateDbFile(ctx)
		if err != nil {
			plog.Errorf(ctx, "UpdateDbFile fail: %s", err)
			useBackDb(ctx)
		}
		// 每天更新一次
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err = UpdateDbFile(ctx)
				if err != nil {
					plog.Errorf(ctx, "UpdateDbFile fail: %s", err)
					useBackDb(ctx)
				}
			}
		}
	}()
}

func useBackDb(ctx context.Context) {
	mutex.Lock()
	defer mutex.Unlock()
	db = backDb
	plog.Infof(ctx, "useBackDb success")
}

func UpdateDbFile(ctx context.Context) error {
	mutex.Lock()
	defer mutex.Unlock()
	// 关闭旧连接
	if db != backDb {
		db.Close()
	}

	zipFilePath := FilePath + FileName + ".zip"
	dbFilePath := FilePath + FileName + ".new"
	err := util.DownloadFile(downloadUrl, zipFilePath, func(length, downLen int64) {})
	if err != nil {
		return fmt.Errorf("DownloadFile: %w", err)
	}
	defer func() {
		os.Remove(zipFilePath)
	}()

	// 解压
	zipFile, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return fmt.Errorf("UpdateDbFile OpenReader: %w", err)
	}
	defer zipFile.Close()
	file, err := zipFile.Open(FileName)
	if err != nil {
		return fmt.Errorf("UpdateDbFile OpenZip: %w", err)
	}
	defer file.Close()

	// 删除旧文件
	_, err = os.Stat(dbFilePath)
	if !os.IsNotExist(err) {
		err = os.Remove(dbFilePath)
		if err != nil {
			return fmt.Errorf("UpdateDbFile RemoveOld: %w", err)
		}
	}

	// 新文件
	nf, err := os.Create(dbFilePath)
	if err != nil {
		return fmt.Errorf("UpdateDbFile Create: %w", err)
	}
	defer nf.Close()
	_, err = io.Copy(nf, file)
	if err != nil {
		return fmt.Errorf("UpdateDbFile Copy: %w", err)
	}

	// 下载完 更换DB
	curDb, err := ip2location.OpenDB(dbFilePath)
	if err != nil {
		return fmt.Errorf("UpdateDbFile OpenDB: %w", err)
	}
	db = curDb
	plog.Infof(ctx, "UpdateDbFile success: %s", dbFilePath)
	return nil
}

func ip2locationRecordConvert(locationRecord ip2location.IP2Locationrecord) LocationRecord {
	return LocationRecord{
		Country_short:      locationRecord.Country_short,
		Country_long:       locationRecord.Country_long,
		Region:             locationRecord.Region,
		City:               locationRecord.City,
		Isp:                locationRecord.Isp,
		Latitude:           locationRecord.Latitude,
		Longitude:          locationRecord.Longitude,
		Domain:             locationRecord.Domain,
		Zipcode:            locationRecord.Zipcode,
		Timezone:           locationRecord.Timezone,
		Netspeed:           locationRecord.Netspeed,
		Iddcode:            locationRecord.Iddcode,
		Areacode:           locationRecord.Areacode,
		Weatherstationcode: locationRecord.Weatherstationcode,
		Weatherstationname: locationRecord.Weatherstationname,
		Mcc:                locationRecord.Mcc,
		Mnc:                locationRecord.Mnc,
		Mobilebrand:        locationRecord.Mobilebrand,
		Elevation:          locationRecord.Elevation,
		Usagetype:          locationRecord.Usagetype,
		Addresstype:        locationRecord.Addresstype,
		Category:           locationRecord.Category,
	}
}
