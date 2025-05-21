package models

// CloudProvider 云服务商模型
type CloudProvider struct {
	BaseModel
	Name        string `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Code        string `gorm:"column:code;type:varchar(50);not null;uniqueIndex:uk_code" json:"code"`
	Description string `gorm:"column:description;type:text" json:"description"`
	// 关联产品
	Products []CloudProduct `gorm:"foreignKey:CloudProviderID" json:"products,omitempty"`
}

// TableName 表名
func (CloudProvider) TableName() string {
	return "cloud_providers"
}