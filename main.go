// main
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"golips_art_engine/conf"
	"golips_art_engine/models"
	"golips_art_engine/utils"
)

const (
	inputDir             = "layers"
	outputDir            = "builds"
	outputImagesDir      = "images"
	outputMetadataDir    = "json"
	outputSolMetadataDir = "json-sol"
)

var (
	debug             bool = true
	rarityDelimiter        = "#"
	colorSetDelimiter      = "$"
	limitDelimiter         = "^"
	dnaDelimiter           = "-"
)

func getMultiVersionFolderName(layerName string) string {
	return outputImagesDir + "-" + layerName
}

func main() {

	log.Println("Reading Config...")

	config, err := conf.GetConfig(debug)

	if err != nil {
		panic(err)
	}

	debug = config.LogSettings.Debug

	if config.RarityDelimiter != "" {
		rarityDelimiter = config.RarityDelimiter
	}

	if config.ColorSetDelimiter != "" {
		colorSetDelimiter = config.ColorSetDelimiter
	}

	if config.LimitDelimiter != "" {
		limitDelimiter = config.LimitDelimiter
	}

	if config.DnaDelimiter != "" {
		dnaDelimiter = config.DnaDelimiter
	}

	log.Println("Set Folders...")

	err = os.RemoveAll(filepath.Join(".", outputDir, "."))

	if err != nil {
		if debug {
			log.Println("[CleanFolder]", err)
		}
		panic(err)
	}

	err = os.MkdirAll(filepath.Join(".", outputDir, outputImagesDir), os.ModePerm)

	if err != nil {
		if debug {
			log.Println("[CreateFolder]", err)
		}
		panic(err)
	}

	err = os.MkdirAll(filepath.Join(".", outputDir, outputMetadataDir), os.ModePerm)

	if err != nil {
		if debug {
			log.Println("[CreateFolder]", err)
		}
		panic(err)
	}

	if config.MetadataSettings.OutputSOLFormat {
		err = os.MkdirAll(filepath.Join(".", outputDir, outputSolMetadataDir), os.ModePerm)

		if err != nil {
			if debug {
				log.Println("[CreateFolder]", err)
			}
			panic(err)
		}
	}

	if config.MultiVersionSettings.LayerName != "" {
		err = os.MkdirAll(filepath.Join(".", outputDir, getMultiVersionFolderName(config.MultiVersionSettings.LayerName)), os.ModePerm)

		if err != nil {
			if debug {
				log.Println("[CreateFolder]", err)
			}
			panic(err)
		}
	}

	for i, _ := range config.LayerConfigurations {
		layersSetup(&config.LayerConfigurations[i])
	}

	if config.Background.Generate {
		config.Background.BrightnessNum = getBrightnessNum(config.Background.Brightness)
	}

	processes, _ := config.ProcessCount.Int64()

	processCount := int(processes)

	if processCount < 1 {
		processCount = 1
	}

	log.Println("Begin Generating...")

	log.Println("Async Process Count: ", processCount)

	var (
		genCount = 0

		// use cache to boost the render
		imgCache = make(map[string]image.Image, 0)
		imgMutex = sync.RWMutex{}

		// save dna to check
		existDNAs = make(map[string]bool, 0)
		dnaMutex  = sync.RWMutex{}

		rarityMutex = sync.RWMutex{}

		processChan = make(chan bool, processCount)
	)

	for i := 0; i < processCount; i++ {
		processChan <- true
	}

	for batch, c := range config.LayerConfigurations {

		if config.LogSettings.ShowGeneratingProgress {
			log.Println("Generating batch: ", batch)
		}

		for i := 1; i <= c.GrowEditionSizeTo; i++ {

			canProcess := <-processChan

			if !canProcess {
				break
			}

			// use this num in async process instead of i
			num := i

			// if start id in config has been set, use it.
			// -1 is because the i start at 1
			if config.DnaSettings.StartId > 0 {
				num = config.DnaSettings.StartId + i - 1
			}

			if batch > 0 {
				num += config.LayerConfigurations[batch-1].GrowEditionSizeTo
			}

			if config.LogSettings.ShowGeneratingProgress {
				log.Println("Generating id: ", num)
			}

			go func() {

				var (
					dstMv = &image.RGBA{}

					hasMultiVersion = config.MultiVersionSettings.LayerName != ""
				)

				if hasMultiVersion {
					dstMv = image.NewRGBA(image.Rect(0, 0, config.Format.Width, config.Format.Height))
				}

				dst := image.NewRGBA(image.Rect(0, 0, config.Format.Width, config.Format.Height))

				// generate random background
				if config.Background.Generate {
					backColor := genColor(config.Background.BrightnessNum)

					draw.Draw(dst, dst.Bounds(), &image.Uniform{backColor}, image.ZP, draw.Src)

					if hasMultiVersion {
						draw.Draw(dstMv, dst.Bounds(), &image.Uniform{backColor}, image.ZP, draw.Src)
					}
				}

				var (
					// make sure the program won't last forever, break it if it can not create new dna.
					dnaCheckTimes = 0
					dna           string
					elements      []models.LayerElement
				)

				for {
					dna, elements = createDNA(&c)

					if debug {
						fmt.Println(fmt.Sprintf("DNA FOR %d: %s", num, dna))
					}

					dnaMutex.Lock()
					if existDNAs[dna] {
						dnaCheckTimes += 1

						if dnaCheckTimes > 20 {
							log.Printf("NFT Generated: %d\n", num-config.DnaSettings.StartId)
							log.Println("Too many duplicate times. Please make sure traits have enough amount.")
							os.Exit(2)
						}
						dnaMutex.Unlock()
						continue
					} else {
						existDNAs[dna] = true
					}
					dnaMutex.Unlock()
					break
				}

				var (
					attributesList = make([]models.MetaDataAttribute, 0)
					attribute      = models.MetaDataAttribute{}
				)

				for _, e := range elements {

					if !e.HideInMetadata {
						if config.MetadataSettings.ShowNoneInMetadata || e.Name != config.MetadataSettings.NoneAttributeName {
							attribute.TraitType = e.BelongLayerName
							attribute.Value = e.Name
							attributesList = append(attributesList, attribute)

						}

						rarityMutex.Lock()
						traits, ok := c.Traits[e.BelongLayerName]

						if ok {
							traits[e.Name] += 1
						}
						rarityMutex.Unlock()
					}

					imgMutex.RLock()
					img, exist := imgCache[e.Path]

					imgMutex.RUnlock()
					if !exist {
						imgFile, err := os.Open(e.Path)

						if err != nil {
							if debug {
								log.Println("[ReadImage]", err)
								log.Println("[ImagePath]", e.Path)
							}
							panic(err)
						}

						defer imgFile.Close()

						img, _, err = image.Decode(imgFile)

						if err != nil {
							if debug {
								log.Println("[ParseImage]", err)
								log.Println("[ImagePath]", e.Path)
							}
							panic(err)
						}

						imgMutex.Lock()
						imgCache[e.Path] = img

						imgMutex.Unlock()
					}

					if e.BelongLayerName != config.MultiVersionSettings.LayerName {
						draw.Draw(dst, dst.Bounds(), img, image.ZP, draw.Over)
					}

					if hasMultiVersion {
						draw.Draw(dstMv, dst.Bounds(), img, image.ZP, draw.Over)
					}
				}

				newImg, err := os.Create(filepath.Join(".", outputDir, outputImagesDir, fmt.Sprintf("%d.png", num)))

				if err != nil {
					if debug {
						log.Println("[CreateImage]", err)
					}
					panic(err)
				}

				defer newImg.Close()

				png.Encode(newImg, dst)

				if hasMultiVersion {
					mvImg, err := os.Create(filepath.Join(".", outputDir, getMultiVersionFolderName(config.MultiVersionSettings.LayerName), fmt.Sprintf("%d.png", num)))

					if err != nil {
						if debug {
							log.Println("[CreateImage]", err)
						}
						panic(err)
					}

					defer mvImg.Close()

					png.Encode(mvImg, dstMv)
				}

				if config.MetadataSettings.NumberAttributes != nil {
					for _, v := range config.MetadataSettings.NumberAttributes {
						if (v.MaxValue - v.MinValue) <= 0 {
							continue
						}
						max := v.MaxValue - v.MinValue
						attribute.TraitType = v.Name
						valueInit := rand.Intn(max)
						attribute.Value = valueInit + v.MinValue
						attribute.DisplayType = "number"
						attribute.MaxValue = v.MaxValue
						attribute.MinValue = v.MinValue
						attributesList = append(attributesList, attribute)
					}
				}

				if config.MetadataSettings.OutputEthFormat {
					saveMetadataErc721(num, dna, config, attributesList)
				}

				if config.MetadataSettings.OutputSOLFormat {
					saveMetadataSolana(num, dna, config, attributesList)
				}

				genCount += 1
				processChan <- true
			}()
		}

		for {
			if len(processChan) == processCount {
				break
			}
		}

		saveRarityFile(batch, c.Traits)
	}

	// make sure all the render works finish
	for {
		if len(processChan) == processCount {
			break
		}
	}

	log.Printf("NFT Generated: %d\nAll Done!\n", genCount)
}

func saveRarityFile(batch int, traits map[string]map[string]int) {

	var (
		list      = make([]models.TraitLayer, 0)
		layer     = models.TraitLayer{}
		trait     = models.Trait{}
		traitList []models.Trait
	)

	for la, elements := range traits {

		layer.Name = la
		layer.Total = len(elements)

		var total = 0

		for _, count := range elements {
			total += count
		}

		traitList = make([]models.Trait, 0)

		for name, count := range elements {
			trait.Name = name
			trait.Total = count
			trait.Rate = fmt.Sprintf("%.1f%%", (float32(count)/float32(total))*100)

			traitList = append(traitList, trait)
		}

		layer.Elements = traitList

		list = append(list, layer)
	}

	newJson, err := os.Create(filepath.Join(".", outputDir, fmt.Sprintf("bathc-%d-rarity.json", batch+1)))

	if err != nil {
		if debug {
			log.Println("[CreateJson]", err)
		}
		panic(err)
	}

	defer newJson.Close()

	je := json.NewEncoder(newJson)

	err = je.Encode(&list)

	if err != nil {
		if debug {
			log.Println("[JsonMarshal]", err)
		}
		panic(err)
	}
}

func saveMetadataErc721(id int, dna string, config *models.Config, attributes []models.MetaDataAttribute) {
	var metadata = models.MetadataErc721{}

	metadata.Name = fmt.Sprintf("%s #%d", config.NamePrefix, id)

	metadata.Description = config.Description

	metadata.Image = fmt.Sprintf("%s/%d.png", config.BaseUri, id)

	if config.MetadataSettings.SaveDnaInMetadata {
		metadata.Dna = utils.GetSha1Hash(dna)
	}

	metadata.Attributes = attributes
	metadata.Compiler = "GoLips Art Engine"

	if config.MetadataSettings.ExtraMetadata != nil {
		metadata.ExtraMetadata = "has#@!"
	}

	if config.MetadataSettings.ShowEditionInMetadata {
		metadata.Edition = id
	}

	newJson, err := os.Create(filepath.Join(".", outputDir, outputMetadataDir, fmt.Sprintf("%d.json", id)))

	if err != nil {
		if debug {
			log.Println("[CreateJson]", err)
		}
		panic(err)
	}

	if config.MetadataSettings.ExtraMetadata != nil {
		data, err := json.Marshal(&metadata)

		if err != nil {
			if debug {
				log.Println("[JsonMarshal]", err)
			}
			panic(err)
		}

		extraData, err := json.Marshal(config.MetadataSettings.ExtraMetadata)

		if err != nil {
			if debug {
				log.Println("[JsonMarshal]", err)
			}
			panic(err)
		}

		extraStr := strings.Replace(string(extraData), "{", "", -1)

		extraStr = strings.Replace(extraStr, "}", "", -1)

		newString := strings.Replace(string(data), `"extra!@#":"has#@!"`, extraStr, -1)

		_, err = newJson.WriteString(newString)

		if err != nil {
			if debug {
				log.Println("[WriteFile]", err)
			}
			panic(err)
		}
	} else {
		je := json.NewEncoder(newJson)

		err = je.Encode(&metadata)

		if err != nil {
			if debug {
				log.Println("[JsonMarshal]", err)
			}
			panic(err)
		}
	}

}

func saveMetadataSolana(id int, dna string, config *models.Config, attributes []models.MetaDataAttribute) {
	var metadata = models.MetadataSolana{}

	metadata.Name = fmt.Sprintf("%s #%d", config.NamePrefix, id)
	metadata.Symbol = config.SolanaMetadata.Symbol
	metadata.SellerFeeBasisPoints = config.SolanaMetadata.SellerFeeBasisPoints

	metadata.Description = config.Description

	metadata.Image = fmt.Sprintf("%d.png", id)
	metadata.ExternalUrl = config.SolanaMetadata.ExternalUrl

	if config.MetadataSettings.SaveDnaInMetadata {
		metadata.Dna = utils.GetSha1Hash(dna)
	}

	metadata.Attributes = attributes
	metadata.Compiler = "GoLips Art Engine"

	if config.MetadataSettings.ExtraMetadata != nil {
		metadata.ExtraMetadata = "has#@!"
	}

	if config.MetadataSettings.ShowEditionInMetadata {
		metadata.Edition = id
	}

	propFile := models.SolanaPropertyFile{
		Uri:  fmt.Sprintf("%d.png", id),
		Type: "image/png",
	}

	prop := models.SolanaProperty{
		Category: "image",
		Creators: config.SolanaMetadata.Creators,
		Files:    []models.SolanaPropertyFile{propFile},
	}

	metadata.Properties = prop

	newJson, err := os.Create(filepath.Join(".", outputDir, outputSolMetadataDir, fmt.Sprintf("%d.json", id)))

	if err != nil {
		if debug {
			log.Println("[CreateJson]", err)
		}
		panic(err)
	}

	if config.MetadataSettings.ExtraMetadata != nil {
		data, err := json.Marshal(&metadata)

		if err != nil {
			if debug {
				log.Println("[JsonMarshal]", err)
			}
			panic(err)
		}

		extraData, err := json.Marshal(config.MetadataSettings.ExtraMetadata)

		if err != nil {
			if debug {
				log.Println("[JsonMarshal]", err)
			}
			panic(err)
		}

		extraStr := strings.Replace(string(extraData), "{", "", -1)

		extraStr = strings.Replace(extraStr, "}", "", -1)

		newString := strings.Replace(string(data), `"extra!@#":"has#@!"`, extraStr, -1)

		_, err = newJson.WriteString(newString)

		if err != nil {
			if debug {
				log.Println("[WriteFile]", err)
			}
			panic(err)
		}
	} else {
		je := json.NewEncoder(newJson)

		err = je.Encode(&metadata)

		if err != nil {
			if debug {
				log.Println("[JsonMarshal]", err)
			}
			panic(err)
		}
	}

}

// pass layer config
func createDNA(layerConfig *models.LayerConfiguration) (string, []models.LayerElement) {
	var (
		elementList = make([]models.LayerElement, 0)
		colorSets   = make(map[string]string, 0)
		// to find out limit quickly, key is 'layer-element'
		usedElements = make(map[string]bool, 0)
		dnaKeys      = make([]string, 0)
		// key: elements that can not be used due to conflict
		conflictUsed = make(map[string]bool, 0)
	)

	for _, layer := range layerConfig.LayersOrder {
		var (
			totalWeight float64 = 0
			color               = ""
		)

		// check if this layer is in a color set
		if layer.Options.ColorSet != "" && !layer.Options.IsColorBase {
			color = colorSets[layer.Options.ColorSet]

			if color == "" {
				continue
			}
		}

		var tempElementList = make([]models.LayerElement, 0)

		for _, v := range layer.Elements {

			if color != "" {
				if v.Color != color {
					continue
				}
			}

			// check if this element has conflict
			if conflictUsed[v.Name] {
				continue
			}

			totalWeight += v.Weight
			tempElementList = append(tempElementList, v)
		}

		if layer.Limits != nil {
			for k, elementList := range layer.Limits {
				if usedElements[k] {
					for _, v := range elementList {

						if color != "" {
							if v.Color != color {
								continue
							}
						}

						// check if this element has conflict
						if conflictUsed[v.Name] {
							continue
						}

						totalWeight += v.Weight
						tempElementList = append(tempElementList, v)

					}
				}
			}
		}

		rand.Seed(time.Now().UnixNano())

		target := rand.Float64() * totalWeight

		for _, v := range tempElementList {
			target -= v.Weight

			if target < 0 {
				// save layer info in elements to simplify the logic
				v.BelongLayerName = layer.Options.DisplayName
				v.HideInMetadata = layer.Options.HideInMetadata

				elementList = append(elementList, v)

				dnaKey := getLimitKey(layer.Options.DisplayName, v.Name)
				usedElements[dnaKey] = true

				conflictNames, exist := layerConfig.ConflictElements[v.Name]

				if exist {
					AddNewConflicts(conflictUsed, conflictNames)
					if debug {
						fmt.Println("Conflict Added")
						fmt.Println(conflictUsed)
					}
				}

				if !layer.Options.BypassDNA {
					dnaKeys = append(dnaKeys, dnaKey)
				}

				// set base color for color set
				if layer.Options.IsColorBase {
					colorSets[layer.Options.ColorSet] = v.Name
				}
				break
			}
		}
	}

	return strings.Join(dnaKeys, dnaDelimiter), elementList
}

func AddNewConflicts(origin map[string]bool, newC string) {

	conflictNames := strings.Split(newC, ",")

	for _, v := range conflictNames {
		origin[v] = true
	}
}

func getLimitKey(layerName string, elementName string) string {
	return fmt.Sprintf("%s%s%s", layerName, limitDelimiter, elementName)
}

func getBrightnessNum(brightness string) float64 {
	brightness = strings.Replace(brightness, "%", "", -1)

	brightNum, err := strconv.ParseFloat(brightness, 64)

	if err != nil {
		if debug {
			log.Println("[Background]get color brightness fail: ", err)
		}

		panic(err)
	}

	return brightNum / 100
}

func genColor(brightness float64) color.RGBA {

	rand.Seed(time.Now().UnixNano())

	hue := rand.Float64()

	return utils.HSLToRGB(hue, 1, brightness)
}

func layersSetup(layer *models.LayerConfiguration) {

	layer.Traits = make(map[string]map[string]int, 0)

	for i, v := range layer.LayersOrder {

		if v.Options.DisplayName == "" {
			layer.LayersOrder[i].Options.DisplayName = v.Name
		}

		list, limits := getElementsFromDir(filepath.Join(".", inputDir, v.Name), v.Options.ColorSet != "", 0)

		layer.LayersOrder[i].Elements = list
		layer.LayersOrder[i].Limits = limits

		// ignore traits if layer shoule be hide in metedata
		if v.Options.HideInMetadata {
			continue
		}

		traits := make(map[string]int, 0)

		for _, e := range list {
			traits[e.Name] = 0
		}

		for _, le := range limits {
			for _, e := range le {
				traits[e.Name] = 0
			}
		}

		layer.Traits[layer.LayersOrder[i].Options.DisplayName] = traits
	}
}

func getElementsFromDir(dir string, isColorSet bool, startId int) ([]models.LayerElement, map[string][]models.LayerElement) {
	fileArray, err := ioutil.ReadDir(dir)

	if err != nil {
		if debug {
			log.Println("[ReadFile]", err)
		}
		panic("Read File Failed:" + err.Error())
	}

	var (
		element = models.LayerElement{}
		list    = make([]models.LayerElement, 0)
		limits  = make(map[string][]models.LayerElement, 0)
	)

	for id, e := range fileArray {

		if e.IsDir() {

			limitList, _ := getElementsFromDir(filepath.Join(dir, e.Name()), isColorSet, len(fileArray)+len(limits))

			limits[e.Name()] = limitList

			continue
		}

		element.Id = startId + id

		name, rarity, color, err := cleanName(e.Name(), isColorSet)

		if err != nil {
			if debug {
				log.Println("[ReadFileName]", err)
				log.Println("[FileName]", dir+"/"+e.Name())
			}
			panic("Read File Name Failed")
		}

		// ignore empty name
		if name == "" {
			continue
		}
		element.Name = name
		element.Weight = rarity
		element.Color = color
		element.Path = dir + "/" + e.Name()

		list = append(list, element)
	}

	return list, limits
}

// get name , rarity , color
func cleanName(name string, isColorSet bool) (string, float64, string, error) {

	// filter system files, such as .DS_Store
	if strings.HasPrefix(name, ".") {
		return "", 0, "", nil
	}

	var color = ""

	if isColorSet {
		colorList := strings.Split(name, colorSetDelimiter)

		if len(colorList) == 2 {
			color = colorList[0]
			name = colorList[1]
		}
	}

	name = name[0 : len(name)-4]

	nameList := strings.Split(name, rarityDelimiter)

	length := len(nameList)

	if length > 2 {
		return "", -1, "", errors.New("Too many rarity delimiter for: " + name)
	}

	var (
		rarity float64 = 1
		err    error
	)

	if length == 2 {
		rarity, err = strconv.ParseFloat(nameList[1], 64)

		if err != nil {
			return "", -1, "", errors.New("Rarity parse error: " + err.Error())
		}
	}

	return nameList[0], rarity, color, nil
}
