package models

// CloudProduct 云产品模型
type CloudProduct struct {
	BaseModel
	CloudProviderID uint          `gorm:"column:cloud_provider_id;not null;index:idx_provider" json:"cloud_provider_id"`
	Name            string        `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Code            string        `gorm:"column:code;type:varchar(50);not null;uniqueIndex:uk_provider_code,priority:2" json:"code"`
	Description     string        `gorm:"column:description;type:text" json:"description"`
	Provider        CloudProvider `gorm:"foreignKey:CloudProviderID" json:"provider,omitempty"`
	// 关联配置项
	ConfigItems []ConfigurationItem `gorm:"foreignKey:ProductID" json:"config_items,omitempty"`
}

// TableName 表名
func (CloudProduct) TableName() string {
	return "cloud_products"
}