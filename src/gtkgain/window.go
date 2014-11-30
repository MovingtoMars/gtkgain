package main

import (
	"github.com/MovingtoMars/gotk3/glib"
	"github.com/MovingtoMars/gotk3/gtk"
	//"github.com/MovingtoMars/gotk3/gdk"
	"github.com/MovingtoMars/gtkgain/src/library"

	"log"
	"sync"
	"strconv"
	"fmt"
)

type window struct {
	win       *gtk.Window
	scroll    *gtk.ScrolledWindow
	headerBar *gtk.HeaderBar
	treeView  *gtk.TreeView
	listStore *gtk.ListStore
	vbox      *gtk.Box // vertical box
	bottomBox      *gtk.Box // bottom box
	pbar      *gtk.ProgressBar
	paths     map[string]*gtk.TreePath
	
	importHelperChan chan string
	
	tagUntaggedButton, untagTaggedButton *gtk.Button
	spinner                              *gtk.Spinner
	fcButton                             *gtk.Button
	inTask                               bool
	taskProgress, taskTotal int

	songQueue     []*library.Song
	songQueueLock sync.Mutex

	lib *library.Library
}

func (w *window) importHelper() {
	for {
		select {
		case path := <- w.importHelperChan:
			w.lib.ImportFromDir(path)
		}	
	}
}

func (w *window) onFolderSelect(path string) {
	w.setSpinner(true)
	w.importHelperChan <- path
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

	w.fcButton, err = gtk.ButtonNewFromIconName("document-open-symbolic", gtk.ICON_SIZE_BUTTON)
	crashIf("Unable to create folder chooser button", err)
	w.headerBar.PackEnd(w.fcButton)
	w.fcButton.Connect("clicked", w.onFcButtonClick)
	w.fcButton.SetTooltipText("Add a folder")

	w.tagUntaggedButton, err = gtk.ButtonNewWithLabel("Tag Untagged")
	crashIf("Unable to create taguntagged button", err)
	w.tagUntaggedButton.SetTooltipText("Calculate ReplayGain for all untagged files")
	crashIf("Unable to create untagtagged button", err)
	w.untagTaggedButton, err = gtk.ButtonNewWithLabel("Untag Tagged")
	w.untagTaggedButton.SetTooltipText("Remove ReplayGain from all tagged files")

	w.headerBar.PackStart(w.tagUntaggedButton)
	w.headerBar.PackStart(w.untagTaggedButton)

	w.tagUntaggedButton.Connect("clicked", w.onTagUntaggedClicked)
	w.untagTaggedButton.Connect("clicked", w.onUntagTaggedClicked)

	w.setTagButtonsSensitive(false)
	w.fcButton.SetSensitive(true)

	w.spinner, err = gtk.SpinnerNew()
	crashIf("Unable to create spinner", err)
	w.headerBar.PackStart(w.spinner)
}

func (w *window) setTagButtonsSensitive(s bool) {
	w.tagUntaggedButton.SetSensitive(s)
	w.untagTaggedButton.SetSensitive(s)
	w.fcButton.SetSensitive(s)
}

func (w *window) onSongUpdate(s *library.Song) {
	glib.IdleAdd(func() {w.setSongGains(s)})
}

func (w *window) setProgressBarFraction(cur, total int) {
	w.pbar.SetFraction(float64(cur) / float64(total))
	stotal := strconv.FormatInt(int64(total), 10)
	w.pbar.SetText(fmt.Sprintf("%*d / %s", len(stotal), cur, stotal))
}

func (w *window) incProgressBarFraction(delta int) {
	w.taskProgress += delta
	w.setProgressBarFraction(w.taskProgress, w.taskTotal)
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

func (w *window) tagAlbums(albums []*library.Album) {
	tasks := make(chan *library.Album, 0)
	
	w.taskTotal = 0
	for _, al := range albums {
		w.taskTotal += len(al.GetSongs())
	}
	
	w.taskProgress = 0
	w.setProgressBarFraction(0, w.taskTotal)
	
	go func() {
		for _, a := range albums {
			tasks <- a
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
					return
				}
				
				paths := make([]string, len(at.GetSongs()))
				for i, so := range at.GetSongs() {
					paths[i] = so.Path()
				}
				numSongs := len(paths)
				
				glib.IdleAdd(func() {
					for _, s := range paths {
						w.setSpinnerForSong(s, true)
					}
					w.treeView.QueueDraw()
				})
				
				err := at.TagGain(w.onSongUpdate)
				if err != nil {
					log.Println(err)
				}
				glib.IdleAdd(func() {
					w.incProgressBarFraction(numSongs)
					w.treeView.QueueDraw()
				})
			}
		}(tasks)
	}
	
	wg.Wait()
	glib.IdleAdd(func() {
		w.inTask = false
		w.setSpinner(false)
		w.setTagButtonsSensitive(true)
		//w.fcButton.SetSensitive(true)
		w.bottomBox.Set("visible", false)
	})
}

func (w *window) untagSongs(songs []*library.Song) {
	// TODO make this use multiple goroutines. Perhaps divide song list up into groups ie. of 20?
	err := library.SongsUntagGain(songs, w.onSongUpdate)
	if err != nil {
		log.Println("Error untagging songs:", err)
	}

	glib.IdleAdd(func() {
		w.inTask = false
		w.treeView.QueueDraw()
		w.setSpinner(false)
		w.setTagButtonsSensitive(true)
		//w.fcButton.SetSensitive(true)
	})
}

func (w *window) onTagUntaggedClicked() {
	a := w.lib.UntaggedAlbums()
	if len(a) == 0 {
		return
	}

	w.inTask = true
	//w.fcButton.SetSensitive(false)
	w.setTagButtonsSensitive(false)
	w.setSpinner(true)
	w.bottomBox.Set("visible", true)
	go w.tagAlbums(a)
}

func (w *window) onUntagTaggedClicked() {
	a := w.lib.TaggedSongs()
	if len(a) == 0 {
		return
	}

	w.inTask = true
	//w.fcButton.SetSensitive(false)
	w.setTagButtonsSensitive(false)
	w.setSpinner(true)
	go w.untagSongs(a)
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

func (w *window) onLoadFinish() {
	glib.IdleAdd(func() {
		w.setTagButtonsSensitive(true)
		w.setSpinner(false)
	})
}

func (w *window) setSpinner(b bool) {
	w.spinner.Set("active", b)
	w.spinner.Set("visible", b)
}

func (w *window) onSongImport(s *library.Song) {
	w.songQueueLock.Lock()
	w.songQueue = append(w.songQueue, s)
	w.songQueueLock.Unlock()
}

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
	w.paths = make(map[string]*gtk.TreePath)
	w.importHelperChan = make(chan string, 50)
	
	var err error

	w.win, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	crashIf("Unable to create window", err)

	w.win.Connect("destroy", w.onDestroy)

	w.setupHeaderBar()
	
	w.vbox, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	w.bottomBox, _ = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	w.vbox.PackEnd(w.bottomBox, false, false, 8)
	w.win.Add(w.vbox)
	
	w.pbar, _ = gtk.ProgressBarNew()
	w.pbar.SetFraction(0.0)
	w.pbar.SetShowText(true)
	w.bottomBox.PackStart(w.pbar, true, true, 16)
	
	w.setupTreeView()
	
	w.win.SetIconName("gtkgain")
	
	w.win.SetDefaultSize(1000, 800)

	w.win.ShowAll()
	w.spinner.Set("visible", false)
	w.bottomBox.Set("visible", false)

	w.songQueue = make([]*library.Song, 0)
	
	go w.importHelper()
	
	glib.TimeoutAdd(100, w.onTimer)

	return w
}
