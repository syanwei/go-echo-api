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
具有以下独特功能：
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
go run main.go -name mamba -age 7 -email songyanwei@uniontech.com
```
什么都没输出，表示一切正常。 提供一个非法的邮箱地址：
```shell script
go run main.go -name manba -age 7 -email songyanwei@uniontech
```
输出如下错误：
```shell script
Key: 'User.Email' Error:Field validation for 'Email' failed on the 'email' tag
```
错误显示不友好。怎么更友好？更国际化？
### 国际化
介绍校验库错误国际化，有个概念要了解下，CLDR
#### CLDR Unicode Common Locale Data Repository
它是 i18n 的一套核心规范（ Common Locale Data Respository），即通用的本地化数据存储库，什么意思呢？比如我们的手机，电脑都可以选择语言模式为 英语、汉语、日语、法语等等，这套操作背后的规范，就是 CLDR；CLDR 是以 Unicode 的编码标准作为前提，将多国的语言文字进行编码的。
[详情](http://cldr.unicode.org/)
需要进行国际化和本地化的主要包括：
- 用于格式化和解析的特定于语言环境的模式：日期，时间，时区，数字和货币值，度量单位，…
- 名称的翻译：语言，脚本，国家和地区，货币，时代，月份，工作日，白天，时区，城市和时间单位，表情符号字符和序列（和搜索关键字），…
- 语言和文字信息：使用的字符；复数情况；性别；大写；分类和搜索规则；写作方向；音译规则；拼写数字的规则；将文本分割成字符，单词和句子的规则；键盘布局…
- 国家/地区信息：语言使用情况，货币信息，日历首选项，星期惯例等…
- 有效性：Unicode 语言环境，语言，脚本，区域和扩展名的定义，别名和有效性信息，…
#### CLDR Go 语言实现
本文讲解的校验库是 go-playground 这个组织创建的，它们还提供了其他的一些有用库，其中就包括了 CLDR 的 Go 语言实现，这就是 locales。
> 该库是从 CLDR 项目生成的一组语言环境，可以单独使用或在 i18n 软件包中使用；这些是专为 https://github.com/go-playground/universal-translator 构建的，但也可以单独他用。

这引出了该组织的另外一个库：universal-translator

universal-translator：一个使用 CLDR 数据+复数规则（比如英语很多复数规则是加 s）的 Go i18n 转换器（翻译器）。该库是  locales 的薄包装，以便存储和翻译文本，供你在应用程序中使用。

#### universal-translator 简明教程
这个通用的翻译器包主要包含了两个核心数据结构：Translator 接口和 UniversalTranslator 结构体，其他的是错误类型。我们先看 Translator 接口。（注意，该包的包名是 ut）

**Translator 接口**
```go
type Translator interface {
    locales.Translator

    // adds a normal translation for a particular language/locale
    // {#} is the only replacement type accepted and are ad infinitum
    // eg. one: '{0} day left' other: '{0} days left'
    Add(key interface{}, text string, override bool) error

    // adds a cardinal plural translation for a particular language/locale
    // {0} is the only replacement type accepted and only one variable is accepted as
    // multiple cannot be used for a plural rule determination, unless it is a range;
    // see AddRange below.
    // eg. in locale 'en' one: '{0} day left' other: '{0} days left'
    AddCardinal(key interface{}, text string, rule locales.PluralRule, override bool) error

    // adds an ordinal plural translation for a particular language/locale
    // {0} is the only replacement type accepted and only one variable is accepted as
    // multiple cannot be used for a plural rule determination, unless it is a range;
    // see AddRange below.
    // eg. in locale 'en' one: '{0}st day of spring' other: '{0}nd day of spring'
    // - 1st, 2nd, 3rd...
    AddOrdinal(key interface{}, text string, rule locales.PluralRule, override bool) error

    // adds a range plural translation for a particular language/locale
    // {0} and {1} are the only replacement types accepted and only these are accepted.
    // eg. in locale 'nl' one: '{0}-{1} day left' other: '{0}-{1} days left'
    AddRange(key interface{}, text string, rule locales.PluralRule, override bool) error

    // creates the translation for the locale given the 'key' and params passed in
    T(key interface{}, params ...string) (string, error)

    // creates the cardinal translation for the locale given the 'key', 'num' and 'digit' arguments
    //  and param passed in
    C(key interface{}, num float64, digits uint64, param string) (string, error)

    // creates the ordinal translation for the locale given the 'key', 'num' and 'digit' arguments
    // and param passed in
    O(key interface{}, num float64, digits uint64, param string) (string, error)

    //  creates the range translation for the locale given the 'key', 'num1', 'digit1', 'num2' and
    //  'digit2' arguments and 'param1' and 'param2' passed in
    R(key interface{}, num1 float64, digits1 uint64, num2 float64, digits2 uint64, param1, param2 string) (string, error)

    // VerifyTranslations checks to ensures that no plural rules have been
    // missed within the translations.
    VerifyTranslations() error
}
```
关于该接口需要需要如下几点说明

- 内嵌了 locales.Translator 接口；
- 几类复数规则：cardinal plural（基数复数规则，即单数和复数两种）；ordinal plural（序数复数规则，如 1st, 2nd, 3rd…）；ordinal plural （范围复数规则，如 0-1）。对中文来说，这里大部分不需要。
- 几个 Add 方法，和上面几类规则对应；一个 key 和 一个带站位符的 text；
- 单字符的几个方法和 Add 几个方法的对应关系：T -> Add；C -> AddCardinal；O -> AddOrdinal；R -> AddRange ；表示用具体的值替换 key 表示的文本 text 中的占位符。
- 以上方法参数中，num 表示占位符处的值，但对于有复数形式的语言，这个值必须符合复数语言的规范，否则会报错；digits 表示 num 值的有效数字（或者说小数位数）；
- VerifyTranslations 确保翻译库中没有缺少对应的语言规则；

**UniversalTranslator 结构体**

它用于保存所有语言环境和翻译数据。该结构体方法不多，我们关注几个核心的。

```go
func New(fallback locales.Translator, supportedLocales ...locales.Translator) *UniversalTranslator
```

New 返回一个 UniversalTranslator 实例，该实例具有后备语言环境（fallback）和应支持的语言环境（supportedLocales）。可以看到，New 函数接收的参数是 locales.Translator 类型，因此我们肯定需要用到 locales 包。

得到 UniversalTranslator 实例后，需要获得 universal-translator 包中的 Translator 接口实例，这就用到了下面几个方法。

- GetTranslator
- GetFallback
- FindTranslator

#### Validator 集成以上两个库提供i18n
Validator 库提供了相应的子库，对以上两个库进行了封装。[中文库](https://github.com/go-playground/validator/translations/zh) ，这些子库提供了一个 RegisterDefaultTranslations，为所有内置标签的验证器注册一组默认翻译。

#### Echo 集成Validator
