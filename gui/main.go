package main

import (
	libSvm "github.com/ewalker544/libsvm-go"
	"github.com/gotk3/gotk3/gtk"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/log"
	"licence-plate-recognition/licencePlateReconition"
	"licence-plate-recognition/utilities"
)

var (
	modelPath string
	imagePath string
)

func main() {
	model := new(libSvm.Model)
	dataMap := utilities.DataMap()
	_ = model
	_ = dataMap

	gtk.Init(nil)

	builder, e := gtk.BuilderNewFromFile("main.glace")
	if e != nil {
		panic(e)
	}

	obj, e := builder.GetObject("mainWindow")
	if e != nil {
		panic(e)
	}

	win, ok := obj.(*gtk.Window)
	if !ok {
		panic("can't get mainWindow")
	}

	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	obj, e = builder.GetObject("btnModel")
	if e != nil {
		panic(e)
	}

	btnModel, ok := obj.(*gtk.Button)
	if !ok {
		panic("can't get btnModel")
	}

	_, e = btnModel.Connect("clicked", func() {
		modelPath, e = ChooseFile(win, "Choose Model", "*.svm")
		if e != nil {
			panic(e)
		}
		model = libSvm.NewModelFromFile(modelPath)
	})
	if e != nil {
		panic(e)
	}

	obj, e = builder.GetObject("img")
	if e != nil {
		panic(e)
	}

	imgDisplay, ok := obj.(*gtk.Image)
	if !ok {
		panic("cant get container image")
	}

	obj, e = builder.GetObject("btnImage")
	if e != nil {
		panic(e)
	}

	btnImage, ok := obj.(*gtk.Button)
	if !ok {
		panic("can't get btnImage")
	}

	_, e = btnImage.Connect("clicked", func() {
		imagePath, e = ChooseFile(win, "Choose Image", "*.png", "*.jpg")
		if e != nil {
			log.Error(e)
		}

		imgDisplay.SetSizeRequest(500, 500)
		imgDisplay.SetFromFile(imagePath)
	})
	if e != nil {
		panic(e)
	}

	obj, e = builder.GetObject("btnRecognize")
	if e != nil {
		panic(e)
	}

	btnRecognize, ok := obj.(*gtk.Button)
	if !ok {
		panic("can't get btnRecognize")
	}

	obj, e = builder.GetObject("txtResult")
	if e != nil {
		panic(e)
	}

	txtResult, ok := obj.(*gtk.TextView)
	if !ok {
		panic("can't get txtResult")
	}

	_, e = btnRecognize.Connect("clicked", func() {
		licencePlateDetected, e := licencePlateReconition.Detect(imagePath, true, 8)
		if e != nil {
			log.Info(e)
			return
		}

		result, e := licencePlateReconition.Recognize(licencePlateDetected, model, dataMap)
		if e != nil {
			log.Info(e)
			return
		}

		buffer, e := txtResult.GetBuffer()
		if e != nil {
			log.Info(e)
			return
		}

		buffer.SetText("")
		buffer.SetText(result)
	})
	if e != nil {
		panic(e)
	}

	win.ShowAll()
	gtk.Main()
}

func ChooseFile(win *gtk.Window, title string, patterns ...string) (string, error) {
	btn, e := gtk.FileChooserDialogNewWith1Button(title, win, gtk.FILE_CHOOSER_ACTION_OPEN, "Open", gtk.RESPONSE_ACCEPT)
	if e != nil {
		return "", e
	}

	filter, e := gtk.FileFilterNew()
	if e != nil {
		return "", e
	}

	for _, v := range patterns {
		filter.AddPattern(v)
	}

	btn.AddFilter(filter)
	btn.Run()
	defer btn.Destroy()

	return btn.GetFilename(), nil
}
