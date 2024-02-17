package goforever

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/com"
)

func TestGlob(t *testing.T) {
	fileName := `test-abc-` + time.Now().Format(`20060102T150405.000`) + `.log`
	testFile := `./` + fileName
	err := os.WriteFile(testFile, []byte{}, os.ModePerm)
	assert.NoError(t, err)
	files, err := filepath.Glob(`./test-abc-` + globPattern)
	assert.NoError(t, err)
	com.Dump(files)
	assert.True(t, com.InSlice(fileName, files))
}

func TestRotateLog(t *testing.T) {
	err := os.WriteFile(`./test-abc.log`, []byte(`test-1`), os.ModePerm)
	assert.NoError(t, err)
	err = RotateLog(`./test-abc.log`, 2)
	assert.NoError(t, err)
}
