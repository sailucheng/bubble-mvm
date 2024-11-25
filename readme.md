# Bubble-MVM

Bubble-MVM 是一个基于 Go 语言实现的基于bubble tea 库的简单扩展框架。

## 主要组件

### 1. Pipe (管道)

`pipe.go` 文件定义了 `Pipe` 结构体，用于管理中间件和控制器的执行流程。主要功能包括：

- **RegisterControllers**: 注册控制器。
- **RegisterMiddlewares**: 注册中间件。
- **CreatePipe**: 创建一个新的管道实例。
- **Execute**: 执行管道中的所有中间件。

### 2. Context (上下文)

用于在管道中传递和管理状态。主要字段包括：

- **Msg**: 当前处理的 `tea.Msg`。
- **TeaModel**: 内部使用的 bubble tea model.
- **Model**: 业务逻辑模型。
- **Viewer**: 负责渲染逻辑的接口。
- **Result**: 存储上一次管道或控制器逻辑的结果。

结果返回工具函数：

- **Cmd**: 创建并返回带有命令的 `Result`。
- **NoAction**: 返回一个不带命令或错误的 `Result`。
- **View**: 创建一个带有新 `Viewer` 的 `Result`。
- **WithError**: 创建一个带有错误的 `Result`。
- **Redirect**: 创建一个带有新 `Viewer` 和 `Model` 的 `Result`。
- **Quit**: 创建一个终止程序的 `Result`。
- **None**: 返回一个预定义的“无操作” `Result`。
- **Abort**: 标记上下文为已中止,退出管道。
- **IsAbort**: 检查上下文是否已中止。

### 3. Controller (控制器)
用于处理应用的业务逻辑,`update`方法的封装,主要方法包括：

- **Filter**: 判断是否需要处理当前上下文。
- **Handle**: 处理当前上下文并返回结果。

### 4. Viewer (视图)

用于处理应用的渲染逻辑。主要方法包括：

- **Init**: 初始化视图。
- **Render**: 渲染视图。


### 5. Model (模型)

框架用于管理应用的状态和逻辑，等同于bubble tea中自定义的Model，负责调用管道逻辑。
根据controller返回的Result进行相关视图更新、跳转、退出等.

主要方法包括：

- **CreateModel**: 负责mvm的模型实例。
- **WithPipe**: 设置管道实例。
- **WithInitView**: 设置初始视图。

## 使用示例
参考examples 代码。



# generate tools
基于mvm的代码生成工具。
构建后使用help进行命令查看。