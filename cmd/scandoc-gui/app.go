package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"proWeb/lib/appCmds"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const NameApp = "scandoc-GUI"

// App struct
type App struct {
	Name    string
	ctx     context.Context
	appCmds appCmds.AppCmds
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		Name:    NameApp,
		ctx:     context.Background(),
		appCmds: appCmds.NewApp(),
	}
}

func (a *App) SetupApp(ctx context.Context) error {
	if err := a.appCmds.StartApp(a.Name); err != nil {
		return err
	}
	a.ctx = ctx

	return nil
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	if err := a.SetupApp(ctx); err != nil {
		panic(err)
	}
}

// OpenLog - открывает папку с логами
func (a *App) OpenLog() {
	const op = ".open_log"

	if err := a.appCmds.OpenLogFolder(a.Name+op, false); err != nil {
		runtime.LogErrorf(a.ctx, err.Error())
	}
}

// OpenStorage - открывает локальное хранилище
func (a *App) OpenStorage() {
	const op = ".open_storage"

	if err := a.appCmds.OpenStorage(a.Name+op, false); err != nil {
		runtime.LogErrorf(a.ctx, err.Error())
	}
}

func (a *App) StartInit() {
	const op = ".init"

	if ok, _ := a.appCmds.CheckInit(a.Name + op); ok {
		runtime.EventsEmit(a.ctx, "init_status", "already-init")
		return
	}

	runtime.EventsEmit(a.ctx, "init_status", "process")

	runtime.EventsEmit(a.ctx, "processing_start", nil)
	if err := a.appCmds.InitApp(a.Name + op); err != nil {
		runtime.EventsEmit(a.ctx, "init_status", "error")
		runtime.EventsEmit(a.ctx, "processing_end", nil)
		return
	}
	runtime.EventsEmit(a.ctx, "processing_end", nil)
	runtime.EventsEmit(a.ctx, "init_status", "success")
}

// CheckInitStatus - проверяет, проинициализировано приложение или нет
func (a *App) CheckInitStatus() string {
	const op = ".check_init_status"

	ok, err := a.appCmds.CheckInit(a.Name + op)
	if err != nil {
		return "error"
	}
	if ok {
		return "already-init"
	}
	return "ready"
}

func (a *App) GetFilesFromStorage() interface{} {
	const op = ".get_files_from_storage"

	files, err := a.appCmds.GetFilesFromStorage(a.Name + op)
	if err != nil {
		runtime.LogErrorf(a.ctx, err.Error())
		return nil
	}

	return files
}

func (a *App) ReadFileFromStorage(fileName string) string {
	const op = ".read_file_from_storage"

	content, err := a.appCmds.ReadFileFromStorage(a.Name+op, fileName)
	if err != nil {
		runtime.LogErrorf(a.ctx, err.Error())
		return ""
	}
	var pretty bytes.Buffer
	if err = json.Indent(&pretty, []byte(content), "", "  "); err != nil {
		runtime.LogErrorf(a.ctx, "JSON format error: %v", err)
		return content
	}
	return pretty.String()
}

func (a *App) SaveFileToStorage(fileName string, content string) error {
	const op = ".save_file_to_storage"
	var data interface{}
	err := json.Unmarshal([]byte(content), &data)
	if err != nil {
		runtime.LogErrorf(a.ctx, err.Error())
		return fmt.Errorf("невалидный json, проверьте ввод")
	}

	if err := a.appCmds.SaveFileToStorage(a.Name+op, fileName, "", data); err != nil {
		runtime.LogErrorf(a.ctx, err.Error())
		return err
	}
	return nil
}
