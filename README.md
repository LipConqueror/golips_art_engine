# golips_art_engine

The inspiration for this tool is based on [HashLips Art Engine](https://github.com/HashLips/hashlips_art_engine).

## Our official links
[🐦  Twitter](https://twitter.com/LipConqueror)

[💄  Discord](https://discord.gg/ey2Ek2bXAD)

[💄  NFT](https://lipconq.com/)


## Why building a new tool?
JS is very general, but there are still some dependency problems (especially canvas-related libraries), we believe many people will face the same problems as we do.

Also, when making our own NFT project, we also encountered many new requirements, some of which can be achieved by modifying the existing code, but some are difficult to achieve.

So in the end, we decided to rewrite a tool in a language we were familiar with.

And our goal is: efficient and sustainable use (only use official libraries).

## Feature Differences

### Unsupported
- canvas blend mode
	- If canvas blend mode is necessary for your project, please keep using [HashLips Art Engine](https://github.com/HashLips/hashlips_art_engine)
	- It takes 50 seconds for a single thread to generate 100 images, but it can be shortened to 20 seconds after multi-threading is enabled. (Complete test on our computer)
- ~~SOL metadata~~ (Supported in v0.0.2)
- static background
- ~~extra metadata~~ (Supported in v0.0.2)
- pixel or gif output
### New
- Multi-threaded generation
	- Spawns 50% faster on our computer compare to single thread
- Colorsets for multi-component color combinations
- Limited component collocation
	- for example: this hair may only appears on that head
- Set start id
- Set custom 'none' property name
- More metadata output options
	- show 'none' property or not, show dna or not, etc.
- Save and load DNA history (coming soon)
- Numerical attributes generation (coming soon)

## Installation

### Use Release

Put the file of the your operating system to the folder to generate the image, and ensure that the directory of the file contains the following structure
```
- Your folder
  - golips_art_engine_XXX
  - conf/config.json
  - layers/your_layers_for_generating
```

Examples of `conf` and `layers` folders you can find in this project

After making sure that the config file and layer files have been configured, go into the folder and run `golips_art_engine_XXX`

### Code

Please make sure to install the official environment of the go language. (>= go 1.16)
Official Go: https://go.dev

Then run:

    git clone https://github.com/LipConqueror/golips_art_engine.git
Enter the folder you just cloned, then run:

    go run .
And DONE! 
Like mentioned above, this tool relay on no third-party library, so you don't need to do anything else.
You can find your NFT and metadata files in `builds`

## Config
You can find `config.json` in `golips_art_engine/conf/`, which decided how the NFT series will be generated.
And here are some descriptions about some fields in `config.json`

| fields | description |
|--|--|
| dnaSettings.startId | what the first nft's id will be |
|metadataSettings.saveDnaInMetadata|save dna in metadata or not|
|metadataSettings.showNoneInMetadata|save none attribute in metadata or not|
|metadataSettings.noneAttributeName|specify your own 'none' file name|
|processCount|how many threads used to generate at the same time, Recommended 2 ~ 3|
|layersOrder.options.hideInMetadata|hide this layer from metadata|
|layersOrder.options.colorSet|which colorset should this layer belong|
|layersOrder.options.isColorBase|if set to 'true', colors will come from this layer for that colorset|

Hope you create some awesome artworks with this code💄

# 中文
本工具的灵感与基础来源于 [HashLips Art Engine](https://github.com/HashLips/hashlips_art_engine).

## 我们的官方链接
[🐦  Twitter](https://twitter.com/LipConqueror)

[💄  Discord](https://discord.gg/ey2Ek2bXAD)

[💄  NFT](https://lipconq.com/)

## 为什么要打造一个新的工具?
JS 非常通用, 但仍然会有一些依赖方面的问题(尤其是涉及到canvas相关库的安装与使用时), 我们相信很多人面临着和我们一样的问题。

同时，在制作我们自己的NFT项目时，我们也遇到了许多的新的需求，其中一些可以通过修改已有的代码来解决，但也有一些很难实现。

所以最终，我们决定用我们熟悉的语言（Go）来重写整个工具。

并且我们的目标是：高效且持续可用（通过只使用官方库来实现）。

## 功能差异

### 不支持
- canvas blend mode
	- 如果canvas blend mode对你项目来说是必须的，那么请继续使用[HashLips Art Engine](https://github.com/HashLips/hashlips_art_engine)
- ~~SOL metadata~~ (v0.0.2已支持)
- 静态背景
- ~~额外的metadata~~ (v0.0.2已支持)
- 像素化和生成动图
### 新增
- 多线程生成
	- 生成速度在我们的电脑上提高了50%以上，对比单线程生成
	- 同样生成100张图片，单线程需要50秒，开启多线程后可以缩短到20秒（在我们的电脑上完成测试）
- 色彩集合-跨组件/层级的多组件/层级颜色统一方案
- 限定组件搭配
	- 举个🌰：只有这个头型能顶那个特殊发型
- 设定起始ID
- 设置自定义的‘空’组件名称
- 更多元数据自定义选项
	- 是否展示‘空’组件, 是否展示DNA, 等等
- 保存和读取DNA历史 (即将到来)
- 数值属性生成 (即将到来)

## 安装使用

### 使用可执行文件(Release)

将对应系统的可执行文件放到要生成图片的文件夹下，并确保文件的目录包含以下结构
```
- Your folder
  - golips_art_engine_XXX
  - conf/config.json
  - layers/your_layers_for_generating
```

关于`conf`和`layers`文件夹的示例，你都可以在本项目中找到

在确保配置文件和图层文件都已经配置完毕后，进入文件夹，运行`golips_art_engine_XXX`

### 代码

请先确保你已经安装了官方的go语言环境。(>= go 1.16)

官方Go: https://go.dev

然后运行:

    git clone https://github.com/LipConqueror/golips_art_engine.git
进入你刚刚clone的文件夹, 然后运行:

    go run .
然后就没了！

就像我们上面提到的那样，这个工具完全没有依赖任何第三方库，所以你不需要再做任何其他的事情。

你可以在`builds`文件夹中找到NFT和元数据。

## 配置文件
你可以在`golips_art_engine/conf/`文件夹下找到`config.json`，其中包含了所有生成NFT的相关配置。

下面是`config.json`中的一些属性及其解释

| 字段 | 解释 |
|--|--|
| dnaSettings.startId | 生成的NFT的起始ID |
|metadataSettings.saveDnaInMetadata|是否要在元数据中保存DNA|
|metadataSettings.showNoneInMetadata|是否要在元数据中保存属性为‘空’的图层|
|metadataSettings.noneAttributeName|设定你自己的‘空’属性名|
|processCount|同时进行生成的线程数，推荐是2~3|
|layersOrder.options.hideInMetadata|是否在元数据中隐藏这一图层|
|layersOrder.options.colorSet|这一图层属于哪个色彩集合|
|layersOrder.options.isColorBase|如果这个字段为'true', 那么这个色彩集合的颜色名，将出自此图层|

希望大家可以用我们的工具创造出更多优秀的作品！💄

也希望国产项目能更多得走向世界！