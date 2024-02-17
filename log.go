package goforever

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/webx-top/com"
)

const globPattern = `*????????T??????.???.log`

// NewLog Create a new file for logging
func NewLog(path string) *os.File {
	if path == "" {
		return nil
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		log.Fatalf("%s\n", err)
		return nil
	}
	return file
}

func RotateLog(logfile string, keepNum int) error {
	name := strings.TrimSuffix(logfile, `.log`)
	prefix := name + `-`
	files, err := filepath.Glob(prefix + globPattern)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	delNum := len(files) + 1 - keepNum
	if delNum > 0 {
		for _, file := range files {
			if delNum < 1 {
				break
			}
			err = com.Remove(file)
			if err != nil {
				log.Printf(`failed to delete %q: %s`, file, err.Error())
			}
			delNum--
		}
	}
	bakFile := prefix + time.Now().Format(`20060102T150405.000`) + `.log`
	err = com.Copy(logfile, bakFile)
	if err != nil {
		return err
	}
	err = os.WriteFile(logfile, []byte{}, os.ModePerm)
	return err
}

func (p *Process) Rotate(keepNum int) {
	if len(p.Logfile) > 0 {
		if err := RotateLog(p.Logfile, keepNum); err != nil {
			log.Println(err.Error())
		}
	}
	if len(p.Errfile) > 0 {
		if err := RotateLog(p.Errfile, keepNum); err != nil {
			log.Println(err.Error())
		}
	}
	if p.children != nil {
		p.children.Rotate(keepNum)
	}
}

func (p *Process) DailyRotate(keepNum int) {
	if p.ctx == nil {
		p.ctx, p.cancel = context.WithCancel(context.Background())
	}
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	nxt := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, time.Local)
	dur := nxt.Sub(now)
	t := time.NewTimer(dur)
	defer t.Stop()
	for {
		select {
		case <-p.ctx.Done():
			return
		case tm := <-t.C:
			p.Rotate(keepNum)
			if tm.Before(tomorrow) {
				time.Sleep(tomorrow.Sub(tm))
				tm = time.Now()
			}
			tomorrow = tm.AddDate(0, 0, 1)
			nxt = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, time.Local)
			dur = nxt.Sub(now)
			t.Reset(dur)
		}
	}
}
