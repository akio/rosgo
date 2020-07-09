package libgengo

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func isRosPackage(dir string) bool {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, f := range files {
		if f.Name() == "package.xml" {
			return true
		}
	}
	return false
}

func findPackages(pkgType string, rosPkgPaths []string) (map[string]string, error) {
	pkgs := make(map[string]string)

	recurseDirForMsgs := func(path string, info os.FileInfo, err error) error {
		// Check whether we could open the path ok.
		if err != nil {
			// Just skip this item and try the next one.
			return nil
		}

		// Check whether this path is a directory.
		if !info.IsDir() {
			// It's not a directory, so we skip it; we only care about directories.
			return nil
		}

		// Check whether this directory is a ROS package.
		if isRosPackage(path) {
			// It's a ROS package.
			pkgName := filepath.Base(path)
			pkgPath := filepath.Join(path, pkgType)
			pkgPaths, err := filepath.Glob(pkgPath + fmt.Sprintf("/*.%s", pkgType))
			if err != nil {
				return nil
			}
			for _, p := range pkgPaths {
				basename := filepath.Base(p)
				rootname := basename[:len(basename)-4] // This is chopping off the file extension.  Horribly.
				fullname := pkgName + "/" + rootname
				pkgs[fullname] = p
			}

			// No point checking INSIDE this one, since it's already a package.
			return filepath.SkipDir
		}

		// Else just keep walking.
		return nil
	}

	// Iterate over the list of paths to search.
	for _, p := range rosPkgPaths {
		err := filepath.Walk(p, recurseDirForMsgs)
		if err != nil {
			// If someone complains, then we just skip searching the reset of this path.
			continue
		}
	}
	return pkgs, nil
}

func FindAllMessages(rosPkgPaths []string) (map[string]string, error) {
	return findPackages("msg", rosPkgPaths)
}

func FindAllServices(rosPkgPaths []string) (map[string]string, error) {
	return findPackages("srv", rosPkgPaths)
}

func FindAllActions(rosPkgPaths []string) (map[string]string, error) {
	return findPackages("action", rosPkgPaths)
}

type MsgContext struct {
	msgPathMap      map[string]string
	srvPathMap      map[string]string
	actionPathMap   map[string]string
	msgRegistry     map[string]*MsgSpec
	msgRegistryLock sync.RWMutex
}

func NewMsgContext(rosPkgPaths []string) (*MsgContext, error) {
	ctx := new(MsgContext)
	msgs, err := FindAllMessages(rosPkgPaths)
	if err != nil {
		return nil, err
	}
	ctx.msgPathMap = msgs

	srvs, err := FindAllServices(rosPkgPaths)
	if err != nil {
		return nil, err
	}
	ctx.srvPathMap = srvs

	acts, err := FindAllActions(rosPkgPaths)
	if err != nil {
		return nil, err
	}
	ctx.actionPathMap = acts

	ctx.msgRegistry = make(map[string]*MsgSpec)
	return ctx, nil
}

func (ctx *MsgContext) Register(fullname string, spec *MsgSpec) {
	ctx.msgRegistryLock.Lock()
	ctx.msgRegistry[fullname] = spec
	ctx.msgRegistryLock.Unlock()
}

func (ctx *MsgContext) LoadMsgFromString(text string, fullname string) (*MsgSpec, error) {
	packageName, shortName, e := packageResourceName(fullname)
	if e != nil {
		return nil, e
	}

	var fields []Field
	var constants []Constant
	for lineno, origLine := range strings.Split(text, "\n") {
		cleanLine := stripComment(origLine)
		if len(cleanLine) == 0 {
			// Skip empty line
			continue
		} else if strings.Contains(cleanLine, ConstChar) {
			constant, e := loadConstantLine(origLine)
			if e != nil {
				return nil, NewSyntaxError(fullname, lineno, e.Error())
			}
			constants = append(constants, *constant)
		} else {
			field, e := loadFieldLine(origLine, packageName)
			if e != nil {
				return nil, NewSyntaxError(fullname, lineno, e.Error())
			}
			fields = append(fields, *field)
		}
	}
	spec, _ := NewMsgSpec(fields, constants, text, fullname, OptionPackageName(packageName), OptionShortName(shortName))
	var err error
	md5sum, err := ctx.ComputeMsgMD5(spec)
	if err != nil {
		return nil, err
	}
	spec.MD5Sum = md5sum
	ctx.Register(fullname, spec)
	return spec, nil
}

func (ctx *MsgContext) LoadMsgFromFile(filePath string, fullname string) (*MsgSpec, error) {
	bytes, e := ioutil.ReadFile(filePath)
	if e != nil {
		return nil, e
	}
	text := string(bytes)
	return ctx.LoadMsgFromString(text, fullname)
}

func (ctx *MsgContext) LoadMsg(fullname string) (*MsgSpec, error) {
	ctx.msgRegistryLock.RLock()
	if spec, ok := ctx.msgRegistry[fullname]; ok {
		ctx.msgRegistryLock.RUnlock()
		return spec, nil
	} else {
		ctx.msgRegistryLock.RUnlock()
		if path, ok := ctx.msgPathMap[fullname]; ok {
			spec, err := ctx.LoadMsgFromFile(path, fullname)
			if err != nil {
				return nil, err
			} else {
				ctx.msgRegistryLock.Lock()
				ctx.msgRegistry[fullname] = spec
				ctx.msgRegistryLock.Unlock()
				return spec, nil
			}
		} else {
			return nil, fmt.Errorf("Message definition of `%s` is not found", fullname)
		}
	}
}

func (ctx *MsgContext) LoadSrvFromString(text string, fullname string) (*SrvSpec, error) {
	packageName, shortName, err := packageResourceName(fullname)
	if err != nil {
		return nil, err
	}

	components := strings.Split(text, "---")
	if len(components) != 2 {
		return nil, fmt.Errorf("Syntax error: missing '---'")
	}

	reqText := components[0]
	resText := components[1]

	reqSpec, err := ctx.LoadMsgFromString(reqText, fullname+"Request")
	if err != nil {
		return nil, err
	}
	resSpec, err := ctx.LoadMsgFromString(resText, fullname+"Response")
	if err != nil {
		return nil, err
	}

	spec := &SrvSpec{
		packageName, shortName, fullname, text, "", reqSpec, resSpec,
	}
	md5sum, err := ctx.ComputeSrvMD5(spec)
	if err != nil {
		return nil, err
	}
	spec.MD5Sum = md5sum

	return spec, nil
}

func (ctx *MsgContext) LoadSrvFromFile(filePath string, fullname string) (*SrvSpec, error) {
	bytes, e := ioutil.ReadFile(filePath)
	if e != nil {
		return nil, e
	}
	text := string(bytes)
	return ctx.LoadSrvFromString(text, fullname)
}

func (ctx *MsgContext) LoadSrv(fullname string) (*SrvSpec, error) {
	if path, ok := ctx.srvPathMap[fullname]; ok {
		spec, err := ctx.LoadSrvFromFile(path, fullname)
		if err != nil {
			return nil, err
		} else {
			return spec, nil
		}
	} else {
		return nil, fmt.Errorf("Service definition of `%s` is not found", fullname)
	}
}

func (ctx *MsgContext) LoadActionFromString(text string, fullname string) (*ActionSpec, error) {
	packageName, shortName, err := packageResourceName(fullname)
	if err != nil {
		return nil, err
	}

	components := strings.Split(text, "---")
	if len(components) != 3 {
		return nil, fmt.Errorf("Syntax error: missing '---'")
	}

	goalText := components[0]
	resultText := components[1]
	feedbackText := components[2]
	goalSpec, err := ctx.LoadMsgFromString(goalText, fullname+"Goal")
	if err != nil {
		return nil, err
	}
	actionGoalText := "Header header\nactionlib_msgs/GoalID goal_id\n" + fullname + "Goal goal\n"
	actionGoalSpec, err := ctx.LoadMsgFromString(actionGoalText, fullname+"ActionGoal")
	if err != nil {
		return nil, err
	}
	feedbackSpec, err := ctx.LoadMsgFromString(feedbackText, fullname+"Feedback")
	if err != nil {
		return nil, err
	}
	actionFeedbackText := "Header header\nactionlib_msgs/GoalStatus status\n" + fullname + "Feedback feedback"
	actionFeedbackSpec, err := ctx.LoadMsgFromString(actionFeedbackText, fullname+"ActionFeedback")
	if err != nil {
		return nil, err
	}
	resultSpec, err := ctx.LoadMsgFromString(resultText, fullname+"Result")
	if err != nil {
		return nil, err
	}
	actionResultText := "Header header\nactionlib_msgs/GoalStatus status\n" + fullname + "Result result"
	actionResultSpec, err := ctx.LoadMsgFromString(actionResultText, fullname+"ActionResult")
	if err != nil {
		return nil, err
	}

	spec := &ActionSpec{
		Package:        packageName,
		ShortName:      shortName,
		FullName:       fullname,
		Text:           text,
		Goal:           goalSpec,
		Feedback:       feedbackSpec,
		Result:         resultSpec,
		ActionGoal:     actionGoalSpec,
		ActionFeedback: actionFeedbackSpec,
		ActionResult:   actionResultSpec,
	}

	md5sum, err := ctx.ComputeActionMD5(spec)
	if err != nil {
		return nil, err
	}
	spec.MD5Sum = md5sum
	return spec, nil
}

func (ctx *MsgContext) LoadActionFromFile(filePath string, fullname string) (*ActionSpec, error) {
	bytes, e := ioutil.ReadFile(filePath)
	if e != nil {
		return nil, e
	}
	text := string(bytes)
	return ctx.LoadActionFromString(text, fullname)
}

func (ctx *MsgContext) LoadAction(fullname string) (*ActionSpec, error) {
	if path, ok := ctx.actionPathMap[fullname]; ok {
		spec, err := ctx.LoadActionFromFile(path, fullname)
		if err != nil {
			return nil, err
		} else {
			return spec, nil
		}
	} else {
		return nil, fmt.Errorf("Action definition of `%s` is not found", fullname)
	}
}

func (ctx *MsgContext) ComputeMD5Text(spec *MsgSpec) (string, error) {
	var buf bytes.Buffer
	for _, c := range spec.Constants {
		buf.WriteString(fmt.Sprintf("%s %s=%s\n", c.Type, c.Name, c.ValueText))
	}
	for _, f := range spec.Fields {
		if f.Package == "" {
			buf.WriteString(fmt.Sprintf("%s\n", f.String()))
		} else {
			subspec, err := ctx.LoadMsg(f.Package + "/" + f.Type)
			if err != nil {
				return "", nil
			}
			submd5, err := ctx.ComputeMsgMD5(subspec)
			if err != nil {
				return "", nil
			}
			buf.WriteString(fmt.Sprintf("%s %s\n", submd5, f.Name))
		}
	}
	return strings.Trim(buf.String(), "\n"), nil
}

func (ctx *MsgContext) ComputeMsgMD5(spec *MsgSpec) (string, error) {
	md5text, err := ctx.ComputeMD5Text(spec)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	hash.Write([]byte(md5text))
	sum := hash.Sum(nil)
	md5sum := hex.EncodeToString(sum)
	return md5sum, nil
}

func (ctx *MsgContext) ComputeActionMD5(spec *ActionSpec) (string, error) {
	goalText, err := ctx.ComputeMD5Text(spec.ActionGoal)
	if err != nil {
		return "", err
	}
	feedbackText, err := ctx.ComputeMD5Text(spec.ActionFeedback)
	if err != nil {
		return "", err
	}
	resultText, err := ctx.ComputeMD5Text(spec.ActionResult)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	hash.Write([]byte(goalText))
	hash.Write([]byte(feedbackText))
	hash.Write([]byte(resultText))
	sum := hash.Sum(nil)
	md5sum := hex.EncodeToString(sum)
	return md5sum, nil
}

func (ctx *MsgContext) ComputeSrvMD5(spec *SrvSpec) (string, error) {
	reqText, err := ctx.ComputeMD5Text(spec.Request)
	if err != nil {
		return "", err
	}
	resText, err := ctx.ComputeMD5Text(spec.Response)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	hash.Write([]byte(reqText))
	hash.Write([]byte(resText))
	sum := hash.Sum(nil)
	md5sum := hex.EncodeToString(sum)
	return md5sum, nil
}

func (ctx *MsgContext) GetMsgs() map[string]*MsgSpec {
	ctx.msgRegistryLock.RLock()
	msgs := ctx.msgRegistry
	ctx.msgRegistryLock.RUnlock()
	return msgs
}
