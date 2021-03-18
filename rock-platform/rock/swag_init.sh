/e/mygopath/src/github.com/swaggo/swag/cmd/swag/swag.exe init -g ./cmd/main.go - o ./docs

# 如果出现: cannot find type definition: v1.xxxxx 这样的报错，有两种可能:
#        1. 确实不存在这个定义的结构体.
#        2. 可能 v1 是当前文件中被导入(import v1 "k8s.io/api/core/v1")了同名的模块，导致模块名(v1)和目录名(v1/xxx.go)冲突.
