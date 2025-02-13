package main

import (
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func runOnMain(fn func()) {
	if driver, ok := fyne.CurrentApp().Driver().(interface {
		RunOnMain(func())
	}); ok {
		driver.RunOnMain(fn)
	} else {
		fn()
	}
}

type VideoInfo struct {
	ID   string
	Name string
}

func main() {
	a := app.New()
	w := a.NewWindow("Wistia 1080p Downloader")
	w.Resize(fyne.NewSize(600, 400))

	var videos []VideoInfo

	entryID := widget.NewEntry()
	entryID.SetPlaceHolder("Wistia Video ID...")

	entryName := widget.NewEntry()
	entryName.SetPlaceHolder("File name...")

	entryFolder := widget.NewEntry()
	entryFolder.SetPlaceHolder("Download folder...")
	browseButton := widget.NewButton("Browse", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if uri != nil {
				entryFolder.SetText(uri.Path())
			}
		}, w)
	})
	folderContainer := container.NewHBox(entryFolder, browseButton)

	list := widget.NewList(
		func() int {
			return len(videos)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			vi := videos[i]
			o.(*widget.Label).SetText(fmt.Sprintf("%s - %s", vi.ID, vi.Name))
		},
	)
	listContainer := container.NewGridWrap(fyne.NewSize(550, 200), list)

	addButton := widget.NewButton("Add to List", func() {
		idText := strings.TrimSpace(entryID.Text)
		nameText := strings.TrimSpace(entryName.Text)
		if idText == "" || nameText == "" {
			dialog.ShowInformation("Error", "Video ID and file name cannot be empty!", w)
			return
		}
		videos = append(videos, VideoInfo{ID: idText, Name: nameText})
		list.Refresh()
		entryID.SetText("")
		entryName.SetText("")
	})

	downloadButton := widget.NewButton("Download All", func() {
		folderPath := strings.TrimSpace(entryFolder.Text)
		if folderPath == "" {
			dialog.ShowInformation("Error", "Download folder cannot be empty!", w)
			return
		}
		if _, err := os.Stat(folderPath); os.IsNotExist(err) {
			err := os.MkdirAll(folderPath, os.ModePerm)
			if err != nil {
				dialog.ShowError(fmt.Errorf("failed to create download folder: %v", err), w)
				return
			}
		}
		if len(videos) == 0 {
			dialog.ShowInformation("Information", "List is empty, add videos first.", w)
			return
		}

		go func() {
			for i := len(videos) - 1; i >= 0; i-- {
				v := videos[i]
				videoId, err := GetVideoIDFromCourseURL(v.ID)
				if err != nil {
					runOnMain(func() {
						dialog.ShowError(fmt.Errorf("Video ID: %s, error: %v", v.ID, err), w)
					})
					continue
				}

				url, err := Get1080pURL(videoId)
				if err != nil {
					runOnMain(func() {
						dialog.ShowError(fmt.Errorf("Video ID: %s, error: %v", v.ID, err), w)
					})
					continue
				}

				fileName := fmt.Sprintf("%s/%s.mp4", folderPath, v.Name)
				err = DownloadFile(url, fileName)
				if err != nil {
					runOnMain(func() {
						dialog.ShowError(fmt.Errorf("Video: %s, download error: %v", v.ID, err), w)
					})
					continue
				}

				runOnMain(func() {
					dialog.ShowInformation("Done", fmt.Sprintf("%s downloaded (%s)", fileName, v.ID), w)
					videos = append(videos[:i], videos[i+1:]...)
					list.Refresh()
				})
			}
		}()
	})

	inputContainer := container.NewGridWithColumns(3, entryID, entryName, addButton)
	mainContainer := container.NewVBox(
		inputContainer,
		folderContainer,
		listContainer,
		downloadButton,
	)

	w.SetContent(mainContainer)
	w.ShowAndRun()
}
