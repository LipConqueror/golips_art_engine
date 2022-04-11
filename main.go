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

	"github.com/LipConqueror/golips_art_engine/conf"
	"github.com/LipConqueror/golips_art_engine/models"
	"github.com/LipConqueror/golips_art_engine/utils"
)

const (
	inputDir          = "layers-lipc"
	outputDir         = "builds"
	outputImagesDir   = "images"
	outputMetadataDir = "json"
)

var (
	debug             bool = true
	rarityDelimiter        = "#"
	colorSetDelimiter      = "|"
	limitDelimiter         = "*"
	dnaDelimiter           = "-"
)

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

	// use cache to boost the render
	var (
		genCount = 0

		imgCache = make(map[string]image.Image, 0)
		imgMutex = sync.RWMutex{}

		processChan = make(chan bool, processCount)
	)

	for i := 0; i < processCount; i++ {
		processChan <- true
	}

	for batch, c := range config.LayerConfigurations {

		var (
			existDNAs = make(map[string]bool, 0)
			dnaMutex  = sync.RWMutex{}
		)

		if config.LogSettings.ShowGeneratingProgress {
			log.Println("Generating batch: ", batch)
		}

		for i := 1; i <= c.GrowEditionSizeTo; i++ {

			canProcess := <-processChan

			if !canProcess {
				break
			}

			// use these num in sync process instead of i
			num := i

			// if start id in config has been set, use it.
			// -1 is because the i start at 1
			if config.DnaSettings.StartId > 0 {
				num = config.DnaSettings.StartId + i - 1
			}

			if config.LogSettings.ShowGeneratingProgress {
				log.Println("Generating id: ", num)
			}

			go func() {
				dst := image.NewRGBA(image.Rect(0, 0, config.Format.Width, config.Format.Height))

				if config.Background.Generate {
					backColor := genColor(config.Background.BrightnessNum)

					draw.Draw(dst, dst.Bounds(), &image.Uniform{backColor}, image.ZP, draw.Src)

				}

				var (
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

						if dnaCheckTimes > 5 {
							panic("Too many duplicate times. Please make sure layers have enough amount.")
						}
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
					}

					imgMutex.RLock()
					img, exist := imgCache[e.Path]

					imgMutex.RUnlock()
					if exist {
						draw.Draw(dst, dst.Bounds(), img, image.ZP, draw.Over)
						continue
					}

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

					draw.Draw(dst, dst.Bounds(), img, image.ZP, draw.Over)
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

				saveMetadata(num, dna, config, attributesList)

				genCount += 1
				processChan <- true
			}()
		}
	}

	// make sure all the render works finish
	for {
		if len(processChan) == processCount {
			break
		}
	}

	log.Printf("NFT Generated: %d\nAll Done!\n", genCount)
}

func saveMetadata(id int, dna string, config *models.Config, attributes []models.MetaDataAttribute) {
	var metadata = models.MetadataErc721{}

	metadata.Name = fmt.Sprintf("%s #%d", config.NamePrefix, id)

	metadata.Description = config.Description

	metadata.Image = fmt.Sprintf("%s/%d.png", config.BaseUri, id)

	if config.MetadataSettings.SaveDnaInMetadata {
		metadata.Dna = utils.GetSha1Hash(dna)
	}

	metadata.Attributes = attributes
	metadata.Compiler = "GoLips Art Engine"

	newJson, err := os.Create(filepath.Join(".", outputDir, outputMetadataDir, fmt.Sprintf("%d.json", id)))

	if err != nil {
		if debug {
			log.Println("[CreateJson]", err)
		}
		panic(err)
	}

	je := json.NewEncoder(newJson)

	err = je.Encode(&metadata)

	if err != nil {
		if debug {
			log.Println("[JsonMarshal]", err)
		}
		panic(err)
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
				v.BelongLayerName = layer.Options.DisplayName
				v.HideInMetadata = layer.Options.HideInMetadata
				elementList = append(elementList, v)

				dnaKey := getLimitKey(layer.Options.DisplayName, v.Name)
				usedElements[dnaKey] = true

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
	for i, v := range layer.LayersOrder {

		if v.Options.DisplayName == "" {
			layer.LayersOrder[i].Options.DisplayName = v.Name
		}

		list, limits := getElementsFromDir(filepath.Join(".", inputDir, v.Name), v.Options.ColorSet != "", 0)

		layer.LayersOrder[i].Elements = list
		layer.LayersOrder[i].Limits = limits

	}
}

func getElementsFromDir(dir string, isColorSet bool, startId int) ([]models.LayerElement, map[string][]models.LayerElement) {
	fileArray, err := ioutil.ReadDir(dir)

	if err != nil {
		if debug {
			log.Println("[ReadFile]", err)
		}
		panic("Read File Failed")
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
