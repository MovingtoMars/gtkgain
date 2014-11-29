package main

import (
	"github.com/MovingtoMars/gotk3/gtk"
)

type folderChooser struct {
	wid *gtk.FileChooserWidget
}

func (w *window) showFolderChooser(callback func(uri string)) {
	chooser, _ := gtk.FileChooserWidgetNew(gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER)

	dialog, _ := gtk.DialogNew()
	c, _ := dialog.GetContentArea()
	c.PackStart(chooser, true, true, 0)

	hbar, _ := gtk.HeaderBarNew()
	hbar.SetTitle("Add Music Folder")
	dialog.SetTitlebar(hbar)

	dialog.SetDefaultSize(800, 600)
	dialog.SetTransientFor(w.win)

	accept, _ := gtk.ButtonNewWithLabel("Open")
	hbar.PackEnd(accept)
	accept.Connect("clicked", func() { callback(chooser.GetFilename()); dialog.Destroy() })

	cancel, _ := gtk.ButtonNewWithLabel("Cancel")
	hbar.PackStart(cancel)
	cancel.Connect("clicked", func() { dialog.Destroy() })

	dialog.ShowAll()
}
