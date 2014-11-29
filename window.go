package main

import (
	"github.com/MovingtoMars/gotk3/glib"
	"github.com/MovingtoMars/gotk3/gtk"
	"github.com/MovingtoMars/gtkgain/library"

	"fmt"
	"log"
	"sync"
)

type window struct {
	win       *gtk.Window
	scroll    *gtk.ScrolledWindow
	headerBar *gtk.HeaderBar
	treeView  *gtk.TreeView
	listStore *gtk.ListStore
	columns   []*gtk.TreeViewColumn
	paths     map[string]*gtk.TreePath

	colSpinnerOn chan []*library.Song
	queueDraw    chan bool

	tagUntaggedButton, untagTaggedButton *gtk.Button
	spinner                              *gtk.Spinner
	finishLoadingChan, setSpinner        chan bool
	taggingDone                          chan bool
	fcButton                             *gtk.Button
	inTask                               bool
	setSongGainsChan                     chan *library.Song

	songQueue     []*library.Song
	songQueueLock sync.Mutex

	lib *library.Library
}

func (w *window) onFolderSelect(uri string) {
	w.spinner.Start()
	go w.lib.ImportFromDir(uri)
}

func (w *window) onFcButtonClick(b *gtk.Button) {
	w.showFolderChooser(w.onFolderSelect)
}

func (w *window) setupHeaderBar() {
	var err error
	w.headerBar, err = gtk.HeaderBarNew()
	crashIf("Unable to create headerbar", err)

	w.headerBar.SetShowCloseButton(true)
	w.headerBar.SetTitle("GtkGain")
	w.win.SetTitlebar(w.headerBar)

	// TODO handle errs
	w.fcButton, err = gtk.ButtonNewFromIconName("document-open-symbolic", gtk.ICON_SIZE_BUTTON)
	crashIf("Unable to create folder chooser button", err)
	w.headerBar.PackEnd(w.fcButton)
	w.fcButton.Connect("clicked", w.onFcButtonClick)

	w.tagUntaggedButton, err = gtk.ButtonNewWithLabel("Tag Untagged")
	w.untagTaggedButton, err = gtk.ButtonNewWithLabel("Untag Tagged")

	w.headerBar.PackStart(w.tagUntaggedButton)
	w.headerBar.PackStart(w.untagTaggedButton)

	w.tagUntaggedButton.Connect("clicked", w.onTagUntaggedClicked)
	w.untagTaggedButton.Connect("clicked", w.onUntagTaggedClicked)

	w.setTagButtonsSensitive(false)

	w.spinner, err = gtk.SpinnerNew()
	w.headerBar.PackStart(w.spinner)
}

func (w *window) setTagButtonsSensitive(s bool) {
	w.tagUntaggedButton.SetSensitive(s)
	w.untagTaggedButton.SetSensitive(s)
}

func (w *window) onSongUpdate(s *library.Song) {
	w.setSongGainsChan <- s

}

func (w *window) setSongGains(s *library.Song) {
	path := w.paths[s.Path()]
	if path == nil {
		log.Fatal("path is nil")
	}
	iter, err := w.listStore.GetIter(path)
	crashIf("Unable to convert path to iter", err)

	if val, err := w.listStore.GetValue(iter, COL_PATH); err == nil {
		if str, _ := val.GetString(); str == s.Path() {
			w.listStore.Set(iter, []int{COL_AGAIN, COL_TGAIN}, []interface{}{s.AlbumGain(), s.TrackGain()})
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

const NUM_HELPERS = 4

func (w *window) tagAlbums(albums []*library.Album, onDone chan bool) {
	tasks := make(chan *library.Album, 0)

	go func() {
		for _, a := range albums {
			tasks <- a
			fmt.Println("task")
		}
		close(tasks)
	}()

	var wg sync.WaitGroup
	for i := 0; i < NUM_HELPERS; i++ {
		wg.Add(1)
		go func(taskChan chan *library.Album) {
			defer wg.Done()
			for {

				at, ok := <-taskChan
				if !ok {
					fmt.Println("ret")
					return
				}
				w.colSpinnerOn <- at.GetSongs()
				err := at.TagGain(w.onSongUpdate)
				if err != nil {
					log.Println(err)
				}
				w.queueDraw <- true
			}
		}(tasks)
	}

	wg.Wait()
	onDone <- true
	w.setSpinner <- false
}

func (w *window) untagSongs(songs []*library.Song, onDone chan bool) {
	err := library.SongsUntagGain(songs, w.onSongUpdate)
	if err != nil {
		log.Println("Error untagging songs:", err)
	}

	w.queueDraw <- true
	onDone <- true
	w.setSpinner <- false
}

func (w *window) onTagUntaggedClicked() {
	a := w.lib.UntaggedAlbums()
	if len(a) == 0 {
		return
	}

	w.inTask = true
	w.setTagButtonsSensitive(false)
	w.spinner.Start()
	go w.tagAlbums(a, w.taggingDone)
}

func (w *window) onUntagTaggedClicked() {
	a := w.lib.TaggedSongs()
	if len(a) == 0 {
		return
	}

	w.inTask = true
	w.setTagButtonsSensitive(false)
	w.spinner.Start()
	go w.untagSongs(a, w.taggingDone)
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
	w.win.Add(w.scroll)
	w.scroll.Add(w.treeView)

	w.columns = make([]*gtk.TreeViewColumn, len(columnList))
	for _, i := range columnList {
		w.columns[i] = w.addColumn(i, columnNames[i], true)
	}
	w.treeView.Set("search-column", COL_TITLE)
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

func (w *window) onLoadFinish() {
	w.finishLoadingChan <- true
	w.setSpinner <- false
	//glib.IdleAdd(w.finishLoadingChan, true)
}

func (w *window) onSongImport(s *library.Song) {
	w.songQueueLock.Lock()
	w.songQueue = append(w.songQueue, s)
	w.songQueueLock.Unlock()
}

func (w *window) setSpinnerForSong(s *library.Song, going bool) {
	path := w.paths[s.Path()]
	if path == nil {
		log.Fatal("path is nil")
	}
	iter, err := w.listStore.GetIter(path)
	crashIf("Unable to convert path to iter", err)

	if val, err := w.listStore.GetValue(iter, COL_PATH); err == nil {
		if str, _ := val.GetString(); str == s.Path() {
			//w.listStore.Set(iter, []int{COL_SPIN}, []interface{} {going})
			w.listStore.Set(iter, []int{COL_AGAIN, COL_TGAIN}, []interface{}{"···", "···"})
		}
	}
}

func (w *window) onTimer() bool {
	w.songQueueLock.Lock()
	for _, s := range w.songQueue {
		w.appendSong(s)
	}
	w.songQueue = make([]*library.Song, 0)
	w.songQueueLock.Unlock()

	all := false

	for !all {
		select {
		case b := <-w.finishLoadingChan:
			w.setTagButtonsSensitive(b)
		case <-w.taggingDone:
			w.setTagButtonsSensitive(true)
			w.inTask = false
		case b := <-w.setSpinner:
			w.spinner.Set("active", b)
		case ss := <-w.colSpinnerOn:
			for _, s := range ss {
				w.setSpinnerForSong(s, true)
			}
			w.treeView.QueueDraw()
		case <-w.queueDraw:
			w.treeView.QueueDraw()
		case s := <-w.setSongGainsChan:
			w.setSongGains(s)
		default:
			all = true
		}
	}

	glib.TimeoutAdd(150, w.onTimer)

	return false
}

func (w *window) appendSong(song *library.Song) {
	if w.paths[song.Path()] != nil {
		return
	}

	iter := w.listStore.Append()
	err := w.listStore.Set(iter, []int{COL_TRACK, COL_TITLE, COL_ALBUM, COL_TGAIN, COL_AGAIN, COL_PATH},
		[]interface{}{song.Track(), song.Title(), song.AlbumName(), song.TrackGain(), song.AlbumGain(), song.Path()})
	crashIf("Unable to add song", err)

	path, err := w.listStore.GetPath(iter)
	crashIf("Unable to get path from iter", err)
	w.paths[song.Path()] = path
}

func (w *window) onDestroy() {
	gtk.MainQuit()
}

func createWindow(lib *library.Library) *window {
	w := &window{lib: lib}

	lib.SetSongLoadReceiver(w.onSongImport)
	lib.SetLoadFinishReceiver(w.onLoadFinish)
	w.finishLoadingChan = make(chan bool, 2)
	w.queueDraw = make(chan bool, 2)
	w.taggingDone = make(chan bool, 1)
	w.setSpinner = make(chan bool, 1)
	w.colSpinnerOn = make(chan []*library.Song, NUM_HELPERS)
	w.paths = make(map[string]*gtk.TreePath)
	w.setSongGainsChan = make(chan *library.Song, NUM_HELPERS*20)

	var err error

	w.win, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	crashIf("Unable to create window", err)

	w.win.Connect("destroy", w.onDestroy)

	w.setupHeaderBar()
	w.setupTreeView()

	w.win.SetDefaultSize(1000, 800)

	w.win.ShowAll()

	w.songQueue = make([]*library.Song, 0)

	glib.TimeoutAdd(100, w.onTimer)

	return w
}
