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
- SOL metadata
- static background
- extra metadata
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
|processCount|how many threads used to generate at the same time|
|layersOrder.options.hideInMetadata|hide this layer from metadata|
|layersOrder.options.colorSet|which colorset should this layer belong|
|layersOrder.options.isColorBase|if set to 'true', colors will come from this layer for that colorset|

Hope you create some awesome artworks with this code💄