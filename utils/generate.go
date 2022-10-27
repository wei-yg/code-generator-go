package utils

import (
	"bytes"
	"generate/config"
	"generate/mysql"
	"gorm.io/gorm"
	"strings"
	"text/template"
)

type ModelFileType struct {
	Package         string
	IsImportTime    bool
	TableInfo       Table
	Name            string
	ModelFields     []Field
	ModelStructName string //模型结构体名字
	IdType          string // 模型 id的类型
}
type Table struct {
	Name            string `gorm:"column:Name"`
	Comment         string `gorm:"column:Comment"`
	IsFilterDelete  bool
	DeleteFieldName string
	DeleteData      int
}

type Field struct {
	Field      string `gorm:"column:Field"`
	Type       string `gorm:"column:Type"`
	Null       string `gorm:"column:Null"`
	Key        string `gorm:"column:Key"`
	Default    string `gorm:"column:Default"`
	Extra      string `gorm:"column:Extra"`
	Privileges string `gorm:"column:Privileges"`
	Comment    string `gorm:"column:Comment"`
}

// CreateFiles 生成文件
func CreateFiles(config config.Yaml) {
	tableNames := config.ModelFileConfig.TableName
	// 获取表单信息
	db := mysql.InitMysql()
	tables := getTables(db, tableNames, config.MysqlConfig.DBName) //生成所有表信息
	for _, table := range tables {
		fields := getFields(db, &table)
		// 生成model
		loadTemplate2GoModelFile(&config.ModelFileConfig, table, fields)

		// 生成vo文件
		loadTemplate2GoVoFile(&config.ModelFileConfig, table, fields)

		//生成ts 接口文件
		loadTemplate2TsTypeFile(&config.ModelFileConfig, table, fields)
	}

	//生成 model 工具文件
	loadTemplate2ModelToolFile(&config.ModelFileConfig)

	// 格式化go文件
	GoFileFormat()
}

// 读取模板引擎创建 Ts文件
func loadTemplate2TsTypeFile(modelFileConfig *config.ModelFile, table Table, fields []Field) {
	//向模板注入函数
	funcMap := template.FuncMap{
		"MysqlTableName2TsInterfaceName": MysqlTableName2TsInterfaceName, //表名转结构体名
		"MysqlFiled2TsType":              MysqlFiled2TsType,              //数据库类型转go类型
		"Hump2JsonField":                 Hump2JsonField,                 //下划线转驼峰
		"FirstLower":                     FirstLower,                     // 首字母小写
		"TsDefaultValue":                 TsDefaultValue,                 //ts根据类型设置默认值
	}
	modelTpl, err := template.New("type.template").Funcs(funcMap).ParseFiles("./tpl/ts/type.template")
	model := &ModelFileType{}
	model.Package = modelFileConfig.PackageName
	model.TableInfo = table
	model.ModelFields = fields
	//model.TsInterfaceName = utils.MysqlTableName2GolangStructName(model.TableInfo.Name)

	var fileBuffer bytes.Buffer
	err = modelTpl.Execute(&fileBuffer, model)
	if err != nil {
		panic(err)
	}
	fileName := FirstLower(UnderlineToHump(table.Name))
	GenerateFile(modelFileConfig.TypePath, fileName+"Type.ts", fileBuffer.String(), modelFileConfig.IsCover)
}

// 生成model格式化文件
func loadTemplate2ModelToolFile(modelFileConfig *config.ModelFile) {
	modelTpl, _ := template.New("modelTool.template").ParseFiles("./tpl/go/modelTool.template")
	data := make(map[string]interface{})
	data["packageName"] = config.YamlConfig.ModelFileConfig.PackageName
	var fileBuffer bytes.Buffer
	_ = modelTpl.Execute(&fileBuffer, data)
	GenerateFile(modelFileConfig.ModelPath, "modelTool.go", fileBuffer.String(), modelFileConfig.IsCover)
}

// 读取模板引擎创建 go文件
func loadTemplate2GoVoFile(modelFileConfig *config.ModelFile, table Table, fields []Field) {
	//向模板注入函数
	funcMap := template.FuncMap{
		"MysqlTableName2GolangStructName": MysqlTableName2GolangStructName, //表名转结构体名
		"MysqlFiled2GolangType":           MysqlFiled2GolangType,           //数据库类型转go类型
		"Hump2JsonField":                  Hump2JsonField,                  //下划线转驼峰
		"FirstLower":                      FirstLower,                      // 首字母小写
	}
	modelTpl, err := template.New("vo.template").Funcs(funcMap).ParseFiles("./tpl/go/vo.template")
	model := &ModelFileType{}
	model.Package = modelFileConfig.PackageName
	model.TableInfo = table
	model.ModelFields = fields
	model.ModelStructName = MysqlTableName2GolangStructName(model.TableInfo.Name)

	var fileBuffer bytes.Buffer
	err = modelTpl.Execute(&fileBuffer, model)
	if err != nil {
		panic(err)
	}
	fileName := FirstLower(UnderlineToHump(table.Name))
	GenerateFile(modelFileConfig.VoPath, fileName+"Vo.go", fileBuffer.String(), modelFileConfig.IsCover)
}

// 读取模板引擎创建 go model文件
func loadTemplate2GoModelFile(modelFileConfig *config.ModelFile, table Table, fields []Field) {
	//向模板注入函数
	funcMap := template.FuncMap{
		"MysqlTableName2GolangStructName": MysqlTableName2GolangStructName, //表名转结构体名
		"MysqlFiled2GolangType":           MysqlFiled2GolangType,           //数据库类型转go类型
		"Hump2JsonField":                  Hump2JsonField,                  //下划线转驼峰
		"FirstLower":                      FirstLower,                      // 首字母小写
	}
	modelTpl, err := template.New("model.template").Funcs(funcMap).ParseFiles("./tpl/go/model.template")
	model := &ModelFileType{}
	model.Package = modelFileConfig.PackageName
	for _, field := range fields {
		// 如果包含 datetime 则导入时间包
		if field.Type == "datetime" {
			model.IsImportTime = true
			break
		}
	}
	for _, field := range fields {
		if field.Field == "id" {
			model.IdType = field.Type
			break
		}
	}
	model.TableInfo = table
	model.ModelFields = fields
	model.ModelStructName = MysqlTableName2GolangStructName(model.TableInfo.Name)

	var fileBuffer bytes.Buffer
	err = modelTpl.Execute(&fileBuffer, model)
	if err != nil {
		panic(err)
	}
	fileName := FirstLower(UnderlineToHump(table.Name))
	GenerateFile(modelFileConfig.ModelPath, fileName+"Model.go", fileBuffer.String(), modelFileConfig.IsCover)
}

// 获取具体表单
func getTables(db *gorm.DB, tableNames []string, dbName string) []Table {

	// 字符串拼接生成表名范围
	tableNamesStr := "'" + strings.Join(tableNames, "','") + "'"

	// 获取指定表信息
	var tables []Table
	if tableNamesStr == "''" {
		db.Raw("SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES " +
			"WHERE table_schema='" + dbName + "';").Find(&tables)
	} else {
		db.Raw("SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES " +
			"WHERE TABLE_NAME IN (" + tableNamesStr + ") AND " +
			"table_schema='" + dbName + "';").Find(&tables)
	}
	return tables
}

// 获取字段的详情信息
func getFields(db *gorm.DB, table *Table) (fields []Field) {
	db.Raw("show FULL COLUMNS from " + table.Name + ";").Find(&fields)
	return
}
