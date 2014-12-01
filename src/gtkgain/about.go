package main

import (
	"github.com/MovingtoMars/gotk3/gtk"
)

func (w *window) showAboutDialog() {
	aboutDialog, _ := gtk.AboutDialogNew()
	
	aboutDialog.SetWebsite("https://github.com/MovingtoMars/gtkgain")
	aboutDialog.SetVersion("v0.1")
	aboutDialog.SetProgramName("GtkGain")
	aboutDialog.SetLogoIconName("gtkgain")
	aboutDialog.SetLicenseType(gtk.LICENSE_GPL_3_0)
	aboutDialog.SetCopyright("Copyright © github.com/MovingtoMars")
	aboutDialog.SetComments("GtkGain is a tool for viewing, adding and removing Replaygain tags")
	
	aboutDialog.SetTransientFor(w.win)
	
	aboutDialog.Connect("response", func(d *gtk.AboutDialog, r int) {
		if r == int(gtk.RESPONSE_DELETE_EVENT) {
			aboutDialog.Destroy()
		}
	})
	
	aboutDialog.Show()
}
