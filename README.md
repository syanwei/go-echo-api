### echo 项目模板
#### 技术栈
- golang
- echo 框架

#### 开发原因 
- 基于规范达成项目架构一致性
- 参数通过Validator来校验

#### 计划功能
- [ ] 架构规范
- [x] [参数验证](standard/validator/validator.md) 
- [ ] 参数命名规范
- [ ] SQLite加密和并发

####  运行服务
- 安装依赖
- 安装 swag   
    ```go get -u github.com/swaggo/swag/cmd/swag```
- 运行```swag init ```生成api文档
- 运行```go run main.go```  

#### 参考文档
- [echo官方文档](https://echo.labstack.com/)
