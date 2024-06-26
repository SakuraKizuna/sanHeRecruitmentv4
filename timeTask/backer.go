package timeTask

import (
	"log"
	"os"
	"path/filepath"
	"sanHeRecruitment/config"
	"sanHeRecruitment/module/backupModule"
	"strings"
	"time"
)

func backer() {
	backupModule.Backer()
}

func expireBackerRemove() {
	expireBackupRemover(config.BackUpConfig.SavePath)
}

func expireBackupRemover(backupPath string) {
	err := filepath.Walk(backupPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		//println(path)
		expireTime := time.Now().AddDate(0, -config.BackerExpireTime, 0)
		expireTimeEx := time.Now().AddDate(0, -config.BackerExpireTime-1, 0)
		if strings.Index(path, "backer_"+expireTime.Format("20060102")) != -1 ||
			strings.Index(path, "backer_"+expireTimeEx.Format("20060102")) != -1 {
			errR := os.Remove(path)
			if errR != nil {
				log.Println("timeTask expireBackupRemover failed,err", errR)
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("timeTask expireBackupRemover filepath.Walk() failed, returned err : %v\n", err)
	}
}
