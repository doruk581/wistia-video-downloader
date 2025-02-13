package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

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
		if len(videos) == 0 {
			dialog.ShowInformation("Information", "List is empty, add videos first.", w)
			return
		}

		go func() {
			for _, v := range videos {
				videoId, err := GetVideoIDFromCourseURL(v.ID)

				if err != nil {
					dialog.ShowError(fmt.Errorf("Video ID: %s, error: %v", v.ID, err), w)
					continue
				}

				url, err := Get1080pURL(videoId)
				if err != nil {
					dialog.ShowError(fmt.Errorf("Video ID: %s, error: %v", v.ID, err), w)
					continue
				}

				fileName := v.Name + ".mp4"
				err = DownloadFile(url, fileName)
				if err != nil {
					dialog.ShowError(fmt.Errorf("Video: %s, download error: %v", v.ID, err), w)
					continue
				}

				dialog.ShowInformation("Done", fmt.Sprintf("%s downloaded (%s)", fileName, v.ID), w)
			}
		}()
	})

	inputContainer := container.NewGridWithColumns(3,
		entryID, entryName, addButton,
	)

	mainContainer := container.NewVBox(
		inputContainer,
		listContainer,
		downloadButton,
	)

	w.SetContent(mainContainer)
	w.ShowAndRun()
}
