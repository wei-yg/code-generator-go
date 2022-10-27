package utils

import (
	"fmt"
	"generate/config"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"os/exec"
	"strings"
)

// MysqlTableName2GolangStructName 数据库表名转结构体名
func MysqlTableName2GolangStructName(str string) string {
	return FirstUpper(UnderlineToHump(str))
}

// MysqlTableName2TsInterfaceName 数据库表面转ts 接口名
func MysqlTableName2TsInterfaceName(str string) string {
	return MysqlTableName2GolangStructName(str)
}

// TsDefaultValue ts默认值
func TsDefaultValue(str string) string {
	str = MysqlFiled2TsType(str)
	if str == "number" {
		return "0"
	}
	return "\"\""
}

// UnderlineToHump 下划线转驼峰
func UnderlineToHump(str string) string {
	newStr := strings.Replace(str, "_", " ", -1)
	newStr = cases.Title(language.English).String(newStr)
	newStr = strings.Replace(newStr, " ", "", -1)
	return newStr
}

// Hump2JsonField 驼峰转json字段 (首字母小写)
func Hump2JsonField(str string) string {
	return FirstLower(UnderlineToHump(str))
}

// FirstLower 字符串首字母小写
func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// FirstUpper 字符串首字母大写
func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// MysqlFiled2GolangType 数据库字段转go类型
func MysqlFiled2GolangType(fieldType string) string {
	typeArr := strings.Split(fieldType, "(")
	switch typeArr[0] {
	case "int":
		return "int"
	case "integer":
		return "int"
	case "mediumint":
		return "int"
	case "bit":
		return "int"
	case "year":
		return "int"
	case "smallint":
		return "int"
	case "tinyint":
		return "int"
	case "bigint":
		return "int64"
	case "decimal":
		return "float32"
	case "double":
		return "float32"
	case "float":
		return "float32"
	case "real":
		return "float32"
	case "numeric":
		return "float32"
	case "timestamp":
		return "time.Time"
	case "datetime":
		return "time.Time"
	case "time":
		return "time.Time"
	default:
		return "string"
	}
}

// MysqlFiled2TsType 数据库字段转Ts类型
func MysqlFiled2TsType(fieldType string) string {
	typeArr := strings.Split(fieldType, "(")
	switch typeArr[0] {
	case "int":
		return "number"
	case "integer":
		return "number"
	case "mediumint":
		return "number"
	case "bit":
		return "number"
	case "year":
		return "number"
	case "smallint":
		return "number"
	case "tinyint":
		return "number"
	case "bigint":
		return "number"
	case "decimal":
		return "number"
	case "double":
		return "number"
	case "float":
		return "number"
	case "real":
		return "number"
	case "numeric":
		return "number"
	case "timestamp":
		return "string"
	case "datetime":
		return "string"
	case "time":
		return "string"
	default:
		return "string"
	}
}

// GoFileFormat go文件格式化
func GoFileFormat() {
	cmd := exec.Command("/bin/bash", "-c", "gofmt -s -w "+config.YamlConfig.ModelFileConfig.ModelPath+"*.go")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() 格式化model文件失败 %s\n", err)
	}
	cmd = exec.Command("/bin/bash", "-c", "gofmt -s -w "+config.YamlConfig.ModelFileConfig.VoPath+"*.go")
	out, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() 格式化vo文件失败 %s\n", err)
	}
	fmt.Printf("格式化文件成功!\n%s\n", string(out))
}
