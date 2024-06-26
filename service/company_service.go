package service

import (
	"github.com/jinzhu/gorm"
	"sanHeRecruitment/config"
	"sanHeRecruitment/dao"
	"sanHeRecruitment/models/mysqlModel"
	"sanHeRecruitment/util/formatUtil"
	"sanHeRecruitment/util/sqlUtil"
	"sanHeRecruitment/util/timeUtil"
	"time"
)

type CompanyService struct {
}

// FuzzyQueryCompaniesPage 模糊查找
func (cs *CompanyService) FuzzyQueryCompaniesPage(fuzzyComName, companyLevel, host string, desStatus, pageNum int) []mysqlModel.CompanyBasicInfo {
	var CompanyBasicInfos []mysqlModel.CompanyBasicInfo
	sqlPage := sqlUtil.PageNumToSqlPage(pageNum, config.PageSize)
	queryQ := dao.DB.Table("companies").
		Where("LOCATE(?,companies.company_name) > 0", fuzzyComName)
	if companyLevel != "0" {
		queryQ = queryQ.Where("com_level = ?", companyLevel)
	}
	if desStatus != -1 {
		queryQ = queryQ.Where("com_status = ?", desStatus)
	}
	queryQ.Offset(sqlPage).Limit(config.PageSize).Find(&CompanyBasicInfos)
	for i, n := 0, len(CompanyBasicInfos); i < n; i++ {
		CompanyBasicInfos[i].PicUrl = formatUtil.GetPicHeaderBody(host, CompanyBasicInfos[i].PicUrl)
	}
	return CompanyBasicInfos
}

func (cs *CompanyService) TotalCount() (mysqlModel.CompaniesTotal, error) {
	var ct mysqlModel.CompaniesTotal
	err := dao.DB.Table("companies").Select("COUNT(*) AS total_nums,"+
		"COUNT( CASE WHEN com_level = 1 THEN 1 ELSE NULL END ) AS companies,"+
		"COUNT( CASE WHEN com_level = 2 THEN 1 ELSE NULL END ) AS servers").
		Where("com_status = ?", 1).
		Find(&ct).Error
	return ct, err
}

func (cs *CompanyService) QueryCompanyInfoByName(CompanyName string) (mysqlModel.Company, error) {
	var companyInfo mysqlModel.Company
	err := dao.DB.Where("company_name=?", CompanyName).Find(&companyInfo).Error
	if err != nil {
		return companyInfo, NoRecord
	}
	return companyInfo, err
}

func (cs *CompanyService) QueryCompanyBasicInfoByName(CompanyName, host string) (mysqlModel.CompanyBasicInfo, error) {
	var companyInfo mysqlModel.CompanyBasicInfo
	err := dao.DB.Table("companies").Where("company_name=?", CompanyName).Find(&companyInfo).Error
	if err != nil {
		return companyInfo, NoRecord
	}
	companyInfo.PicUrl = formatUtil.GetPicHeaderBody(host, companyInfo.PicUrl)
	return companyInfo, err
}

func (cs *CompanyService) QueryCompanyInfoById(CompanyId int, host string) (companyInfo mysqlModel.CompanyOut, err error) {
	err = dao.DB.Table("companies").
		Select("com_id,pic_url,company_name,description,companies.phone,scale_tag,person_scale,address,update_time,update_person,com_level,vip,name").
		Joins("INNER JOIN users on users.username = companies.update_person").
		Where("com_id =?", CompanyId).Find(&companyInfo).Error
	if err != nil {
		return
	}
	companyInfo.UpdateTimeOut = timeUtil.TimeFormatToStr(companyInfo.UpdateTime)
	companyInfo.PicUrl = formatUtil.GetPicHeaderBody(host, companyInfo.PicUrl)
	return
}

// AddCompanyInfo 上传公司信息
func (cs *CompanyService) AddCompanyInfo(username, comHeadPic, companyName, description, scaleTag, personScale, address, phone string,
	updateTime time.Time, TargetLevelInt int) error {
	scaleLevel, ok := sqlUtil.CompanyScLevels[personScale]
	if !ok {
		scaleLevel = 1
	}
	var companyInfo = mysqlModel.Company{
		PicUrl:       comHeadPic,
		CompanyName:  companyName,
		Description:  description,
		ScaleTag:     scaleTag,
		PersonScale:  personScale,
		UpdatePerson: username,
		UpdateTime:   updateTime,
		Address:      address,
		ComLevel:     TargetLevelInt,
		Phone:        phone,
		Vip:          0,
		ComStatus:    0,
		ScaleLevel:   scaleLevel,
		Applicant:    username,
	}
	err := dao.DB.Save(&companyInfo).Error
	if err != nil {
		return err
	}
	return err
}

func (cs *CompanyService) AddCompanyInfoTX(username, comHeadPic, companyName, description, scaleTag, personScale, address, phone, userPresident string,
	updateTime time.Time, TargetLevelInt int) (err error) {
	err = dao.DB.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作
		scaleLevel, ok := sqlUtil.CompanyScLevels[personScale]
		if !ok {
			scaleLevel = 1
		}
		var companyInfo = mysqlModel.Company{
			PicUrl:       comHeadPic,
			CompanyName:  companyName,
			Description:  description,
			ScaleTag:     scaleTag,
			PersonScale:  personScale,
			UpdatePerson: username,
			UpdateTime:   updateTime,
			Address:      address,
			ComLevel:     TargetLevelInt,
			Phone:        phone,
			Vip:          0,
			ComStatus:    0,
			ScaleLevel:   scaleLevel,
			Applicant:    username,
		}
		if err := tx.Table("companies").Save(&companyInfo).Error; err != nil {
			return err
		}

		var comNewInfo mysqlModel.Company
		err := tx.Table("companies").Where("company_name=?", companyName).Find(&comNewInfo).Error
		if err != nil {
			return err
		}

		var UpgradeInfo = mysqlModel.Upgrade{
			Qualification: 0,
			TargetLevel:   TargetLevelInt,
			FromUsername:  username,
			CompanyId:     comNewInfo.ComId,
			ApplyTime:     updateTime,
			CompanyExist:  0,
			Show:          0,
			TimeId:        time.Now().Unix(),
		}
		if err := tx.Table("upgrades").Save(&UpgradeInfo).Error; err != nil {
			return err
		}

		if err := tx.Table("users").Where("username = ?", username).
			UpdateColumns(map[string]interface{}{
				"president": userPresident,
			}).
			Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}
		return nil
	})
	return
}

// FuzzyQueryCompanies 模糊查找
func (cs *CompanyService) FuzzyQueryCompanies(fuzzyComName, companyLevel string, desStatus int) (companyName []mysqlModel.CompanyName) {
	queryQ := dao.DB.Table("companies").Select("com_id,company_name,com_level").
		Where("LOCATE(?,companies.company_name) > 0", fuzzyComName)
	if companyLevel != "0" {
		queryQ = queryQ.Where("com_level = ?", companyLevel)
	}
	if desStatus != -1 {
		queryQ = queryQ.Where("com_status = ?", desStatus)
	}
	queryQ.Find(&companyName)
	return
}

// FuzzyQueryAllCompanies 模糊查找
func (cs *CompanyService) FuzzyQueryAllCompanies(fuzzyComName string) (companyName []mysqlModel.CompanyName) {
	sql := "select com_id,company_name,com_level from `companies` WHERE  LOCATE(?,companies.company_name) > 0"
	dao.DB.Raw(sql, fuzzyComName).Scan(&companyName)
	return
}

// UpdateCompanyHeadPic 修改公司头像
func (cs *CompanyService) UpdateCompanyHeadPic(comId int, picUrl, updatePerson string) (err error) {
	var companyInfo mysqlModel.Company
	err = dao.DB.Where("com_id = ?", comId).Find(&companyInfo).Error
	if err != nil {
		return
	}
	companyInfo.PicUrl = picUrl
	companyInfo.UpdateTime = time.Now()
	companyInfo.UpdatePerson = updatePerson
	err = dao.DB.Save(&companyInfo).Error
	return
}

// UpdateCompanyInfo 修改公司信息
func (cs *CompanyService) UpdateCompanyInfo(comId int,
	scaleTag, personScale, address, updatePerson, description, phone string) (err error) {
	var companyInfo mysqlModel.Company
	err = dao.DB.Where("com_id = ?", comId).Find(&companyInfo).Error
	if err != nil {
		return
	}
	companyInfo.ScaleTag = scaleTag
	companyInfo.PersonScale = personScale
	companyInfo.Address = address
	companyInfo.UpdateTime = time.Now()
	companyInfo.UpdatePerson = updatePerson
	companyInfo.Description = description
	companyInfo.Phone = phone
	err = dao.DB.Save(&companyInfo).Error
	return
}

func (cs *CompanyService) QueryAllCompanies(comLevel, pageNum int, host string) ([]*mysqlModel.Company, error) {
	var comInfos []*mysqlModel.Company
	sqlPage := sqlUtil.PageNumToSqlPage(pageNum, webPageSize)
	err := dao.DB.Table("companies").
		Where("com_level = ?", comLevel).Where("com_status = ?", 1).Limit(webPageSize).Offset(sqlPage).
		Find(&comInfos).Error
	for i, m := 0, len(comInfos); i < m; i++ {
		comInfos[i].PicUrl = formatUtil.GetPicHeaderBody(host, comInfos[i].PicUrl)
	}
	return comInfos, err
}
