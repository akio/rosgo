package libtest_gengo

import (
	"github.com/edwinhayes/rosgo/libgengo"
	"path/filepath"
	"os"
	"testing"
)

const allMsgDefContent string = `
byte FOO=1
byte BAR=2
string HOGE=hoge

Header h
byte b
int8 i8
int16 i16
int32 i32
int64 i64
uint8 u8
uint16 u16
uint32 u32
uint64 u64
float32 f32
float64 f64
time t
duration d
string s
std_msgs/ColorRGBA c
uint32[] dyn_ary
uint32[2] fix_ary`

func RTTest(t *testing.T) {

	// Test whether a bad path is handled ok.
	msgs, err := libgengo.FindAllMessages([]string{"", " ", ":"})
	if err != nil { t.Error("error in FindAllMessages: " + err.Error()) }
	if len(msgs) != 0 { t.Errorf("wrong number of entries (should be 0, got %v)", len(msgs)) }

	// Create a temporary directory structure to parse.

	tmp := os.TempDir()
	testDir := filepath.Join(tmp, "gengo_test")
		// ./a1/b1/c1/d1 - Normal.
	p1 := filepath.Join(testDir, "a1", "b1", "c1", "d1")
	if err := os.MkdirAll(filepath.Join(p1, "msg"), 0744); err != nil { t.Error("couldn't create temp files to use for test: " + err.Error()) }
	defer os.RemoveAll(testDir)
	if err := mkFile(filepath.Join(p1, "package.xml"), ""); err != nil { t.Error("couldn't create temp files to use for test: " + err.Error()) }
	if err := mkFile(filepath.Join(p1, "msg", "test.msg"), allMsgDefContent); err != nil { t.Error("couldn't create temp files to use for test: " + err.Error()) }
		// ./a2/b2/c2/d2 - Package but no message.
	p2 := filepath.Join(testDir, "a2", "b2", "c2", "d2")
	if err := os.MkdirAll(filepath.Join(p2, "msg"), 0744); err != nil { t.Error("couldn't create temp files to use for test: " + err.Error()) }
	if err := mkFile(filepath.Join(p2, "package.xml"), ""); err != nil { t.Error("couldn't create temp files to use for test: " + err.Error()) }
		// ./a3/b3/c3/d3/e3/f3 - Deeply nested package.
	p3 := filepath.Join(testDir, "a3", "b3", "c3", "d3", "e3", "f3")
	if err := os.MkdirAll(filepath.Join(p3, "msg"), 0744); err != nil { t.Error("couldn't create temp files to use for test: " + err.Error()) }
	if err := mkFile(filepath.Join(p3, "package.xml"), ""); err != nil { t.Error("couldn't create temp files to use for test: " + err.Error()) }
	if err := mkFile(filepath.Join(p3, "msg", "test.msg"), allMsgDefContent); err != nil { t.Error("couldn't create temp files to use for test: " + err.Error()) }
		// ./a4 - Shallowly nested package.
	p4 := filepath.Join(testDir, "a4")
	if err := os.MkdirAll(filepath.Join(p4, "msg"), 0744); err != nil { t.Error("couldn't create temp files to use for test: " + err.Error()) }
	if err := mkFile(filepath.Join(p4, "package.xml"), ""); err != nil { t.Error("couldn't create temp files to use for test: " + err.Error()) }
	if err := mkFile(filepath.Join(p4, "msg", "test.msg"), allMsgDefContent); err != nil { t.Error("couldn't create temp files to use for test: " + err.Error()) }
		// ./a5/b5/c5ab - Branching directories. 
	p5a := filepath.Join(testDir, "a5", "b5", "c5a")
	p5b := filepath.Join(testDir, "a5", "b5", "c5b")
	if err := os.MkdirAll(filepath.Join(p5a, "msg"), 0744); err != nil { t.Error("couldn't create temp files to use for test: " + err.Error()) }
	if err := mkFile(filepath.Join(p5a, "package.xml"), ""); err != nil { t.Error("couldn't create temp files to use for test: " + err.Error()) }
	if err := mkFile(filepath.Join(p5a, "msg", "test.msg"), allMsgDefContent); err != nil { t.Error("couldn't create temp files to use for test: " + err.Error()) }
	if err := os.MkdirAll(filepath.Join(p5b, "msg"), 0744); err != nil { t.Error("couldn't create temp files to use for test: " + err.Error()) }
	if err := mkFile(filepath.Join(p5b, "package.xml"), ""); err != nil { t.Error("couldn't create temp files to use for test: " + err.Error()) }
	if err := mkFile(filepath.Join(p5b, "msg", "test.msg"), allMsgDefContent); err != nil { t.Error("couldn't create temp files to use for test: " + err.Error()) }
	// ./a6/b6/c6a/b - One inside the other.
	p6a := filepath.Join(testDir, "a6", "b6", "c6")
	p6b := filepath.Join(testDir, "a6", "b6", "c6", "d6")
	if err := os.MkdirAll(filepath.Join(p6a, "msg"), 0744); err != nil { t.Error("couldn't create temp files to use for test: " + err.Error()) }
	if err := mkFile(filepath.Join(p6a, "package.xml"), ""); err != nil { t.Error("couldn't create temp files to use for test: " + err.Error()) }
	if err := mkFile(filepath.Join(p6a, "msg", "test.msg"), allMsgDefContent); err != nil { t.Error("couldn't create temp files to use for test: " + err.Error()) }
	if err := os.MkdirAll(filepath.Join(p6b, "msg"), 0744); err != nil { t.Error("couldn't create temp files to use for test: " + err.Error()) }
	if err := mkFile(filepath.Join(p6b, "package.xml"), ""); err != nil { t.Error("couldn't create temp files to use for test: " + err.Error()) }
	if err := mkFile(filepath.Join(p6b, "msg", "test.msg"), allMsgDefContent); err != nil { t.Error("couldn't create temp files to use for test: " + err.Error()) }

	// Try to parse the directory structure.
	packagePaths := make([]string, 0)
	packagePaths = append(packagePaths, testDir)
	msgs, err = libgengo.FindAllMessages(packagePaths)
	if err != nil { t.Error("error in FindAllMessages: " + err.Error()) }
	if _, ok := msgs["d1/test"]; !ok { t.Error("didn't create d1/test successfully") }
	if _, ok := msgs["d2/test"]; ok { t.Error("created d2/test which shouldn't be possible") }
	if _, ok := msgs["f3/test"]; !ok { t.Error("didn't create f3/test successfully") }
	if _, ok := msgs["a4/test"]; !ok { t.Error("didn't create a4/test successfully") }
	if _, ok := msgs["c5a/test"]; !ok { t.Error("didn't create c5a/test successfully") }
	if _, ok := msgs["c5b/test"]; !ok { t.Error("didn't create c5b/test successfully") }
	if _, ok := msgs["c6/test"]; !ok { t.Error("didn't create c6/test successfully") }
	if _, ok := msgs["d6/test"]; ok { t.Error("created d6/test which shouldn't be possible") }
}

func mkFile(path string, contents string) error {
	f, err := os.Create(path)
	if err != nil { return err }
	defer f.Close()
	_, err = f.Write([]byte(contents))
	return err
}