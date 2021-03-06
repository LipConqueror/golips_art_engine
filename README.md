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
- Conflict elements (v0.0.3)
    - for example: this necklace will not match this dress
- Set start id
- Set custom 'none' property name
- More metadata output options
	- show 'none' property or not, show dna or not, etc.
- Save and load DNA history (coming soon)
- Numerical attributes generation (v0.0.5)

### Limited combination

For example, let's assume that `layer A` will be rendered before `layer B`. In `layer A`, there is an element called `hairstyle 1`, which is a cute ball head, and another element called `hairstyle 2`, is ordinary long hair, and `layer B` contains all kinds of accessories that may appear on the head, but this happens, we designed a `hairpin 1`, which is only suitable for ball head `hairstyle 1`, but not for `hairstyle 2`, then we definitely hope that only when `hairstyle 1` is randomly selected in `layer A`, it is possible to randomly get `hairpin 1`.

Then, at this time, it is our turn to use the 'limited combination' function. First, we need to find the `limitDelimiter` field in `config.json`, and then get the value configured by this field. Then, we can enter `layer B` folder, create a new folder inside and name it `layer A^hairstyle 1`, then move the file of `hairpin 1` into the folder we just created, and done! No additional configuration is required! The program will automatically match. Only in `layer A` when `hairstyle 1` is randomly selected, it will read the elements under the folder `layer A^hairstyle 1`, but there is one thing to note, elements in this folder share random probabilities with all other elements in `layer B`.

![](https://github.com/LipConqueror/golips_art_engine/blob/main/golips_example_limit_06.png)  ![](https://github.com/LipConqueror/golips_art_engine/blob/main/golips_example_limit_05.png)

![](https://github.com/LipConqueror/golips_art_engine/blob/main/golips_example_limit_01.png)


![](https://github.com/LipConqueror/golips_art_engine/blob/main/golips_example_limit_02.png)

### Conflict Elements

For example, let's assume `layer A` is rendered before `layer B`, and there is an element named `cloth1` in `layer A`, which is conflict with the element named `necklace1` in `layer B`, before this, maybe you need to manually check these conditions, but now, you can directly add the following settings to `config.json` to achieve this goal:

When the random algorithm selects `cloth1` in `layer A`, `necklace1` in `layer B` will be ignored directly, instead, it will select one of the other elements.

![](https://github.com/LipConqueror/golips_art_engine/blob/main/conflict_example_1.png)


![](https://github.com/LipConqueror/golips_art_engine/blob/main/conflict_example_2.png)


![](https://github.com/LipConqueror/golips_art_engine/blob/main/conflict_example_3.png)


Also, if `cloth1` is also conflict with `necklace2`, you can change the line to below style:
```
"cloth1": "necklace1,necklace2"
```
i.e.: you can use commas `,` to connect multiple upper-level elements that conflict with the same underlying element

Note: you should always use the name of the underlying element as the key

### Numeric properties

Add the following configuration to the `metadataSettings` field in the `config.json` file to generate a numerical attribute field

![](https://github.com/LipConqueror/golips_art_engine/blob/main/number_attr_config_en.jpg)
![](https://github.com/LipConqueror/golips_art_engine/blob/main/number_attr_output_en.jpg)

## Installation

### Use Release

Put the file of the your operating system to the folder to generate the image, and ensure that the directory of the file contains the following structure
```
- Your folder
  - golips_art_engine_XXX (The executable file corresponding to your computer system)
  - conf/config.json (Configuration file, create a 'conf' folder and put the 'config.json' file in it)
  - layers/your_layers_for_generating (Your layer files)
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
- 冲突元素(v0.0.3)
    - 举个🌰：当穿了这件衣服，就不能搭配那条项链了
- 设定起始ID
- 设置自定义的‘空’组件名称
- 更多元数据自定义选项
	- 是否展示‘空’组件, 是否展示DNA, 等等
- 保存和读取DNA历史 (即将到来)
- 数值属性生成 (v0.0.5)

### 限定组合

举个例子，我们假设`layer A`会在`layer B`之前渲染, 在`layer A`中呢，有一个元素叫`发型1`，是一个可爱的丸子头，另一个元素叫`发型2`，是普通长发，而`layer B`中包含的，是各种可能出现在头上的饰品，但就会出现这种情况，我们设计了一款`簪子1`，只适用于丸子头的`发型1`，而不适用于`发型2`，那么我们肯定希望只有`layer A`中随机到`发型1`的时候，才有可能随机到`簪子1`。

那么这种时候，就轮到我们的‘限定组合’功能登场了，我们先在`config.json`中找到`limitDelimiter`字段，然后获取这个字段配置的值，接着，我们就可以进入`layer B`的文件夹，在里面再创建一个新的文件夹，并将其命名为`layer A^发型1`，接着，将`簪子1`的文件，移入到我们刚刚创建的文件夹中，然后就结束了！不再需要其他的配置！程序会自动匹配，只有在`layer A`中，随机到`发型1`时，才会去读取`layer A^发型1`这个文件夹下的元素，但有一点需要注意，这个文件夹下的元素，会和`layer B`中的其他所有元素共享随机概率。

![](https://github.com/LipConqueror/golips_art_engine/blob/main/golips_example_limit_03.png)  ![](https://github.com/LipConqueror/golips_art_engine/blob/main/golips_example_limit_04.png)

![](https://github.com/LipConqueror/golips_art_engine/blob/main/golips_example_limit_01.png)


![](https://github.com/LipConqueror/golips_art_engine/blob/main/golips_example_limit_02.png)

### 冲突元素

再次举个例子，我们假设`layer A`会在`layer B`之前渲染，在`layer A`中呢，有一个元素叫`cloth1`，但这个元素与`layer B`中的元素`necklace1`冲突了，在这之前，也许你需要手动检查找出这种情况，但现在，你可以直接在`config.json`加入如下配置来实现这个效果：

当随机算法在`layer A`中选择了`cloth1`之后，算法将直接无视`layer B`中的`necklace1`元素，而从其他元素中选择一个。

![](https://github.com/LipConqueror/golips_art_engine/blob/main/conflict_example_1.png)


![](https://github.com/LipConqueror/golips_art_engine/blob/main/conflict_example_2.png)


![](https://github.com/LipConqueror/golips_art_engine/blob/main/conflict_example_3.png)


同时，如果`cloth1`同时还与`necklace2`冲突，那么你可以修改这一行代码为如下的形式：
```
"cloth1": "necklace1,necklace2"
```
即：你可以使用英文逗号`,`来连接与同一个底层元素冲突的多个上层元素

注意：应该永远用底层元素的名称作为key来使用

### 数值属性

在`config.json`文件中的`metadataSettings`字段中添加如下配置，即可生成数值化的属性字段

![](https://github.com/LipConqueror/golips_art_engine/blob/main/number_attr_config_zh.jpg)
![](https://github.com/LipConqueror/golips_art_engine/blob/main/number_attr_output_zh.jpg)

## 安装使用

### 使用可执行文件(Release)

将对应系统的可执行文件放到要生成图片的文件夹下，并确保文件的目录包含以下结构
```
- 你的文件夹
  - golips_art_engine_XXX（对应你电脑系统的可执行文件）
  - conf/config.json （配置文件，创建一个conf文件夹，再将config.json文件放入其中）
  - layers/your_layers_for_generating （你的图层文件们）
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