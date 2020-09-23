## Go语言Echo框架：输入校验，自定义Validator
**问题说明**
- 参数自动绑定和校验是web框架很重要的功能，可以极大提升开发速度。
- Echo没有内置数据校验的能力，即没有默认的Validator实现。可以集成第三方的数据校验库。
- 通过实现以下接口，可以方便的讲任何第三方数据校验库集成到Echo中。[详见](https://github.com/go-playground/validator)
```go
Validator interface {
  Validate(i interface{}) error
```

### go-playgroud/validator
最流行的第三方校验库，具有以下独特功能：
- 通过使用验证标签（tag）或自定义验证程序进行跨字段和跨结构体验证；
- 切片，数组和 map，可以验证任何的多维字段或多层级；
- 能够深入（多维）了解 map 键和值以进行验证；
- 通过在验证之前确定其基础类型来处理接口类型；
- 处理自定义字段类型，例如 sql driver Valuer；
- 别名验证标签，允许将多个验证映射到单个标签，以便更轻松地定义结构上的验证；
- 提取自定义定义的字段名称，例如可以指定在验证时提取 JSON 名称，并将其用于结果 FieldError 中；
- 可自定义的 i18n 错误消息；
- gin Web 框架的默认验证器；
### demo
执行命令：
```shell script
go run validator.go -name mamba -age 7 -email songyanwei@uniontech.com
```
什么都没输出，表示一切正常。 提供一个非法的邮箱地址：
```shell script
go run validator.go -name manba -age 7 -email songyanwei@uniontech
```
输出如下错误：
```shell script
Key: 'User.Email' Error:Field validation for 'Email' failed on the 'email' tag
```
错误显示不友好。怎么更友好？更国际化？

### 国际化
介绍校验库错误消息国际化之前，有个概念要了解下，CLDR
#### CLDR Unicode Common Locale Data Repository
它是 [i18n](https://baike.baidu.com/item/I18N/6771940?fr=aladdin) 的一套核心规范（ Common Locale Data Respository），即通用的本地化数据存储库，什么意思呢？比如我们的手机，电脑都可以选择语言模式为 英语、汉语、日语、法语等等，这套操作背后的规范，就是 CLDR；CLDR 是以 Unicode 的编码标准作为前提，将多国的语言文字进行编码的。
[详情](http://cldr.unicode.org/)

CLDR Go 语言实现，本文讲解的校验库是 go-playground 这个组织创建的，它们还提供了其他的一些有用库，其中就包括了 CLDR 的 Go 语言实现，这就是 locales。
> 该库是从 CLDR 项目生成的一组语言环境，可以单独使用或在 i18n 软件包中使用；这些是专为 https://github.com/go-playground/universal-translator 构建的，但也可以单独他用。

这引出了该组织的另外一个库：universal-translator

universal-translator：一个使用 CLDR 的 Go i18n 转换器（翻译器）。

这个通用的翻译器包主要包含了两个核心数据结构：Translator 接口和 UniversalTranslator 结构体。（注意，该包的包名是 ut）

#### Validator 集成以上两个库提供i18n
Validator 库提供了相应的子库，对以上两个库进行了封装。[中文库](https://github.com/go-playground/validator/translations/zh) ，这些子库提供了一个 RegisterDefaultTranslations，为所有内置标签的验证器注册一组默认翻译。
```shell script
go run validatorZh.go
```
输出
```shell script
Name为必填字段
Age必须大于或等于1
Email为必填字段
```
#### Echo 集成Validator
```shell script
go run validatorEcho.go
```

[http测试](https://httpie.org/docs#examples) 
```shell script
http -v :8888 age:=0 email=song@163.com
```

**具体支持**
- [validator](https://github.com/go-playground/validator)

- [定制参数校验](https://github.com/go-playground/validator/blob/master/_examples/custom/main.go)
