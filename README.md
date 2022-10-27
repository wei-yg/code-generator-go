# 一个用于根据数据库表 生成go model(集成gorm操作) 生成器
###### 一枚前端 学习go ing... 希望看到的朋友 star一下就是给我最大的动力
###### 哪里做的不好也希望大佬指正

## 更新日志
#### 2022-10-27
###### 1.支持go model 代码生成
###### 2.支持前端TS 接口模型生成
###### 3.集成gorm 查询
###### 4.支持事务处理


## 目录结构
```
code-generator-go
├─main.go                           //主入口
├─utils                             //工具包 也是生成器核心代码
|   ├─file.go
|   ├─format.go
|   └generate.go
├─tpl                               // 模板文件目录   (可以根据自己项目情况修改模板文件)
|  ├─ts                                 // 生成前端的ts文件模板目录
|  | └type.template
|  ├─go                                 // 生成go的文件模板目录
|  | ├─model.template
|  | ├─modelTool.template
|  | └vo.template
├─mysql                             // mysql初始化 gorm
|   └init.go
├─generateFile                      // 文件生成的目录
|      ├─vo                             
|      ├─types                          // 存放给前端的ts文件
|      ├─model                          // 存放go的model文件  集成了gorm
├─config
|   ├─conf.yaml
|   └config.go
```

## run run run

### 虽然还不知道怎么用,但是 先试一把
```
go 1.18+ (当前使用1.19,可修改为1.18)
在model中使用了泛型 所以需要1.18+  可自行修改 tpl/go/model.template 修改删除泛型代码适配低版本
```
```
修改config目录下 conf.yaml文件 修改为你自己的数据库信息
```
```
执行 main.go
```
##

####  不出意外 generateFile 目录下即生成了对应文件
###### 此刻代码生成完结,可以自己参照代码 修改 template 模板实现自己想要的代码

#
#
# 生成的 model 文件使用说明

### 看到此文件了解生成的代码约定   ````````          重要! 重要! 重要!
#### 
#### 约定一 : 数据库如需要 创建时间 修改时间 软删除 需按照以下字段(也是gorm的默认字段)
```
DeletedAt gorm.DeletedAt `json:"deletedAt"`                       // 删除时间
CreatedAt time.Time      `json:"createdAt"`                       // 创建时间
UpdatedAt time.Time      `json:"updatedAt"`                       // 修改时间
```
###### 如需要改此字段也可以在 tpl/go/model.template 内自行修改
#### 约定二 : 命名方式
```
struct 大驼峰  AbcAbc
json   小驼峰  abcAbc
数据库  下划线  abc_abc
```
###### 如需要改此字段也可以在 tpl/go/model.template和 utils/ 实现文件 内自行修改

#### 开始使用 : 打开任意 model 文件  
###### 我这里以一个简单的 user 表为例
```
type User struct {
	Id        int            `json:"id" gorm:"column:id;primaryKey;not null"`
	Name      string         `json:"name" gorm:"column:name"`         // 姓名
	Password  string         `json:"password" gorm:"column:password"` // 密码
	Gender    int            `json:"gender" gorm:"column:gender;"`    // 性别 1 女  2 男
	Age       int            `json:"age" gorm:"column:age;"`          // 年龄
	DeletedAt gorm.DeletedAt `json:"deletedAt"`                       // 删除时间
	CreatedAt time.Time      `json:"createdAt"`                       // 创建时间
	UpdatedAt time.Time      `json:"updatedAt"`                       // 修改时间
}

type UserModel interface {
	Create(data *User, ops ...SetUserWhereOption) (id int, err error)                           // 创建
	UpdateById(id int, data User, ops ...SetUserWhereOption) (err error)                        // 根据id修改
	UpdateByCondition(data User, ops ...SetUserWhereOption) (err error)                         // 根据条件map 批量修改
	FindById(id int, ops ...SetUserWhereOption) (result User, err error)                        // 根据id 查询
	FindOneByCondition(ops ...SetUserWhereOption) (result User, err error)                      // 根据条件查询一个
	FindByCondition(ops ...SetUserWhereOption) (result []User, err error)                       // 根据条件查询多个
	FindCountByCondition(ops ...SetUserWhereOption) (count int64, err error)                    // 查询符合条件的个数
	FindByConditionWithPage(ops ...SetUserWhereOption) (result ResultPageData[User], err error) // 根据条件分页查询
}
type userModel struct {
	db    *gorm.DB
	table string
}

func NewUserModel(db *gorm.DB) UserModel {
	return &userModel{
		db:    db,
		table: "user",
	}
}

// SetUserWhereOption 设置查询条件
type SetUserWhereOption func(o *WhereOption[User])

(`````````````此处省略部分代码`````````)

// 分页查询
func (m *userModel) FindByConditionWithPage(ops ...SetUserWhereOption) (result ResultPageData[User], err error) {
	query := make(map[string]interface{})
	whereOption := &WhereOption[User]{
		PageNum:  1,
		PageSize: 10,
		QueryMap: query,
	}
	for _, o := range ops {
		o(whereOption)
	}
	offsetVal := (whereOption.PageNum - 1) * whereOption.PageSize
	tx, _ := userByConditionBase(m, ops...)
	// 获取where条件
	tx = getWhereStrByWhereOption[User](tx, whereOption)
	tx = getOrderByWhereOption[User](tx, whereOption)
	err = tx.Count(&result.Total).Offset(offsetVal).Limit(whereOption.PageSize).Find(&result.List).Error
	if int64(whereOption.PageNum*whereOption.PageSize) >= result.Total {
		result.NextPage = -1
	} else {
		result.NextPage = result.PageNum + 1
	}
	result.PageNum = whereOption.PageNum
	result.PageSize = whereOption.PageSize
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return result, nil
	}
	return result, err
}
```
#### 老规矩先试为敬  我们先调用这个 FindByConditionWithPage 根据条件分页查询 (如果不懂泛型 先去学1.18更新的泛型再往下看)
```
db := 假如你已经初始化gorm db
userModel := model.NewUserModel(db)
pageRes, err := userModel.FindByConditionWithPage()   //获取 第一页数据 每页10条(默认PageNum 1,PageSize 10)
```
##### pageRes 分页返回结构体
```
// 文件在  generateFile/model/modelTool.go (生成的文件)
type ResultPageData[T any] struct {
	List     []T   `json:"list"`            // 结构体切片 泛型
	Total    int64 `json:"total"`           // 总条数
	PageNum  int   `json:"pageNum"`         // 当前页码
	NextPage int   `json:"nextPage"`        // 下一页码数  -1 即为最后一页
	PageSize int   `json:"pageSize"`        // 每页大小
}
```
##
#### 加入查询条件试  依然是调用 FindByConditionWithPage  利用go的 选项模式 实现参数可传可不传(不懂选项模式的可以继续去学习)
```
// 如果你用的是goland 在括号内输入 fun是会自动提示的  
// 这个回调接收一个 o 类型为 *model.WhereOption[model.User]
// 任意查询实现即是对这个 o 进行赋值 model内会解析o内 内容作为查询条件   
pageRes, err := userModel.FindByConditionWithPage(func(o *model.WhereOption[model.User]) {
    o.PageNum = 1   // 设置查询第一页
    o.PageSize = 5  // 设置每页大小为5页
    
    // 这里 o.QueryEntry即为泛型传进来的 model.User
    o.QueryEntry.Gender = 1  // 查询 条件 女
})
```
###### 此刻即可满足定义查询
##
### 下面看看复杂查询
#### 假如我需要查询的是年龄小于30,id大于1的 且带分页的
```
// 如果你用的是goland 在括号内输入 fun是会自动提示的  
// 这个回调接收一个 o 类型为 *model.WhereOption[model.User]
// 任意查询实现即是对这个 o 进行赋值 model内会解析o内 内容作为查询条件   
pageRes, err := userModel.FindByConditionWithPage(func(o *model.WhereOption[model.User]) {
    o.PageNum = 1   // 设置查询第一页
    o.PageSize = 5  // 设置每页大小为5页
    
    // 这里 o.QueryMap 则是参考 gorm的map传参  
    // 可选项参考  文件 generateFile/model/modelTool.go 内 mapCondition2whereSql 方法 case 2
    o.QueryMap["age < ?"] = 30  // 查询年龄小于30
    o.QueryMap["id > ?"] = 1   // 查询id大于1的
})
```
#### 事务处理
```
// 利用 o.Tx 传入 同一个db操作即可 (参考 gorm Transaction)
db.Transaction(func(tx *gorm.DB) error {
    err := customerDao.UpdateByCondition(currentCustomerInfo, func(o *model.WhereOption[model.Customer]) {
        o.Tx = tx
    })
    if err != nil {
        return err
    }
    // 操作记录
    consumptionDao := l.svcCtx.ConsumptionLogDao
    consumptionLog := model.ConsumptionLog{
        ConsumerId: req.Id,
        OperatorId: tokenInfo.Id,
        Money:      req.Money,
    }
    _, err = consumptionDao.Create(&consumptionLog, func(o *model.WhereOption[model.ConsumptionLog]) {
        o.Tx = tx
    })
    if err != nil {
        return err
    }
    return nil
})
```
###### 还有模糊查询  排序等就不一一列举 自行开发吧!
###### model 内还有很大优化空间 继续加油...

