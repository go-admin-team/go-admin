package models

type SysTables struct {
	TableId             int    `gorm:"primaryKey;autoIncrement" json:"tableId"`        //表编码
	TBName              string `gorm:"column:table_name;size:255;" json:"tableName"`   //表名称
	TableComment        string `gorm:"size:255;" json:"tableComment"`                  //表备注
	ClassName           string `gorm:"size:255;" json:"className"`                     //类名
	TplCategory         string `gorm:"size:255;" json:"tplCategory"`                   //
	PackageName         string `gorm:"size:255;" json:"packageName"`                   //包名
	ModuleName          string `gorm:"size:255;" json:"moduleName"`                    //go文件名
	ModuleFrontName     string `gorm:"size:255;comment:前端文件名;" json:"moduleFrontName"` //前端文件名
	BusinessName        string `gorm:"size:255;" json:"businessName"`                  //
	FunctionName        string `gorm:"size:255;" json:"functionName"`                  //功能名称
	FunctionAuthor      string `gorm:"size:255;" json:"functionAuthor"`                //功能作者
	PkColumn            string `gorm:"size:255;" json:"pkColumn"`
	PkGoField           string `gorm:"size:255;" json:"pkGoField"`
	PkJsonField         string `gorm:"size:255;" json:"pkJsonField"`
	Options             string `gorm:"size:255;" json:"options"`
	TreeCode            string `gorm:"size:255;" json:"treeCode"`
	TreeParentCode      string `gorm:"size:255;" json:"treeParentCode"`
	TreeName            string `gorm:"size:255;" json:"treeName"`
	Tree                bool   `gorm:"size:1;default:0;" json:"tree"`
	Crud                bool   `gorm:"size:1;default:1;" json:"crud"`
	Remark              string `gorm:"size:255;" json:"remark"`
	IsDataScope         int    `gorm:"size:1;" json:"isDataScope"`
	IsActions           int    `gorm:"size:1;" json:"isActions"`
	IsAuth              int    `gorm:"size:1;" json:"isAuth"`
	IsLogicalDelete     string `gorm:"size:1;" json:"isLogicalDelete"`
	LogicalDelete       bool   `gorm:"size:1;" json:"logicalDelete"`
	LogicalDeleteColumn string `gorm:"size:128;" json:"logicalDeleteColumn"`
	ModelTime
	ControlBy
}

func (SysTables) TableName() string {
	return "sys_tables"
}
