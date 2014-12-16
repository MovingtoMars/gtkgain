package main

import (
	"github.com/MovingtoMars/gotk3/glib"
	"github.com/MovingtoMars/gotk3/gtk"
	"github.com/MovingtoMars/gtkgain/src/library"

	"log"
)

func (w *window) setSongGains(s *library.Song) {
	path := w.paths[s.Path()]
	if path == nil {
		log.Fatal("path is nil")
	}
	iter, err := w.listStore.GetIter(path)
	crashIf("Unable to convert path to iter", err)

	if val, err := w.listStore.GetValue(iter, COL_PATH); err == nil {
		if str, _ := val.GetString(); str == s.Path() {
			w.listStore.Set(iter, []int{COL_AGAIN, COL_TGAIN}, []interface{}{s.Gain(library.GAIN_ALBUM), s.Gain(library.GAIN_TRACK)})
		}
	}
	/*for iter, b := w.listStore.GetIterFirst(); b; b = w.listStore.IterNext(iter) {
		if val, err := w.listStore.GetValue(iter, COL_PATH); err == nil {
			if str, _ := val.GetString(); str == s.Path() {
				w.listStore.Set(iter, []int{COL_AGAIN, COL_TGAIN}, []interface{} {s.AlbumGain(), s.TrackGain()})
			}
		}
	}*/
}

const (
	COL_TRACK = iota
	COL_TITLE
	COL_ALBUM
	COL_TGAIN
	COL_AGAIN
	COL_PATH
)

var columnList = []int{COL_TRACK, COL_TITLE, COL_ALBUM, COL_TGAIN, COL_AGAIN, COL_PATH}
var columnNames = []string{"Track", "Title", "Album", "Track Gain", "Album Gain", "Path"}

func (w *window) setupTreeView() {
	var err error
	w.listStore, err = gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING)
	crashIf("Unable to create list store", err)

	w.treeView, err = gtk.TreeViewNewWithModel(w.listStore)
	crashIf("Unable to create tree view", err)

	w.scroll, err = gtk.ScrolledWindowNew(nil, nil)
	w.vbox.PackStart(w.scroll, true, true, 0)
	w.scroll.Add(w.treeView)

	for _, i := range columnList {
		w.addColumn(i, columnNames[i], true)
	}
	w.treeView.Set("search-column", COL_TITLE)
	
	/*tentry, _ := gtk.TargetEntryNew("text/uri-list", gtk.TARGET_OTHER_APP, 0)
	w.fcButton.DragDestSet(gtk.DEST_DEFAULT_HIGHLIGHT, []gtk.TargetEntry {*tentry}, gdk.ACTION_DEFAULT)
	w.fcButton.Connect("drag-motion", w.onDragMotion)*/
}

func (w *window) addColumn(id int, name string, resizable bool) *gtk.TreeViewColumn {
	cr, err := gtk.CellRendererTextNew()
	crashIf("Unable to create cell renderer", err)
	col, err := gtk.TreeViewColumnNewWithAttribute(name, cr, "text", id)
	crashIf("Unable to add column", err)
	w.treeView.AppendColumn(col)
	col.Set("resizable", resizable)
	return col
}

/*func (w *window) onDragMotion(widget *gtk.Button, dc *gdk.DragContext, x, y int, time uint) {
	fmt.Println("test", x, y, gdk.Atom(dc.ListTargets().Data).Name())
}*/

func (w *window) setSpinnerForSong(spath string, going bool) {
	path := w.paths[spath]
	if path == nil {
		log.Fatal("path is nil")
	}
	iter, err := w.listStore.GetIter(path)
	crashIf("Unable to convert path to iter", err)

	if val, err := w.listStore.GetValue(iter, COL_PATH); err == nil {
		if str, _ := val.GetString(); str == spath {
			//w.listStore.Set(iter, []int{COL_SPIN}, []interface{} {going})
			w.listStore.Set(iter, []int{COL_AGAIN, COL_TGAIN}, []interface{}{"···", "···"})
		}
	}
}
