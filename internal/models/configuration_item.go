package models

// ConfigurationItem 安全配置基线项模型
type ConfigurationItem struct {
	BaseModel
	CloudProviderID    uint          `gorm:"column:cloud_provider_id;not null;index:idx_provider_product,priority:1" json:"cloud_provider_id"`
	ProductID          uint          `gorm:"column:product_id;not null;index:idx_provider_product,priority:2" json:"product_id"`
	Name               string        `gorm:"column:name;type:varchar(200);not null" json:"name"`
	RecommendedValue   string        `gorm:"column:recommended_value;type:text;not null" json:"recommended_value"`
	RiskDescription    string        `gorm:"column:risk_description;type:text" json:"risk_description"`
	CheckMethod        string        `gorm:"column:check_method;type:text" json:"check_method"`
	ConfigurationMethod string        `gorm:"column:configuration_method;type:text" json:"configuration_method"`
	Reference          string        `gorm:"column:reference;type:text" json:"reference"`
	Provider           CloudProvider `gorm:"foreignKey:CloudProviderID" json:"provider,omitempty"`
	Product            CloudProduct  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// TableName 表名
func (ConfigurationItem) TableName() string {
	return "configuration_items"
}