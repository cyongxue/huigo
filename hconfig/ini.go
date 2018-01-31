package hconfig

/**
 和ini配置文件相关
 */

import (
	"io/ioutil"
	"sync"
	"path/filepath"
	"bytes"
	"bufio"
	"io"
	"os"
	"strings"
	"fmt"
	"os/user"
	"strconv"
	"errors"
)

var (
	defaultSection  = "default"              // 顶层缺省的section
	bNumComment     = []byte{'#'}            // 多行注释
	bSemComment     = []byte{';'}            // 直接的行注释
	bEmpty          = []byte{}
	bEqual          = []byte{'='}            // 单引号
	bDQuote         = []byte{'"'}            // 双引号
	sectionStart    = []byte{'['}            // 段开始，section start
	sectionEnd      = []byte{']'}            // 段截止，section end
	lineBreak       = "\n"                   // 行分割, line break

	includeSec      = "include"
	sepSecKey       = "::"
)


/************************************************************************
 定义一个空类用于囊括ini配置文件的解析处理操作
 todo: 该方式和interface的使用结合和紧密的
 */
type InitConfig struct {

}

/**
 文件解析
 */
func (c *InitConfig) parseFile(fileName string) (*IniConfigContainer, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return c.parseData(filepath.Dir(fileName), data)
}

/**
 解析data
 */
func (c *InitConfig) parseData(dir string, data []byte) (*IniConfigContainer, error) {
	cfgContainer := &IniConfigContainer{
		data:           make(map[string]map[string]string),
		sectionComment: make(map[string]string),
		keyComment:     make(map[string]string),
		RWMutex:        sync.RWMutex{},
	}
	cfgContainer.Lock()
	defer cfgContainer.Unlock()

	buf := bufio.NewReader(bytes.NewBuffer(data))
	// check the BOM
	head, err := buf.Peek(3)
	if err == nil && head[0] == 239 && head[1] == 187 && head[2] == 191 {
		for i := 1; i <= 3; i++ {
			buf.ReadByte()
		}
	}

	var comment bytes.Buffer                        // save comment here
	section := defaultSection
	for {
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		}
		// 任何路径相关的错误，均throw
		if _, ok := err.(*os.PathError); ok {
			return nil, err
		}

		line = bytes.TrimSpace(line)
		if bytes.Equal(line, bEmpty) {
			continue                    // if empty line, continue
		}

		// parse line, if comment, and write comment to comment buf
		var bComment []byte
		switch {
		case bytes.HasPrefix(line, bNumComment):
			bComment = bNumComment
		case bytes.HasPrefix(line, bSemComment):
			bComment = bSemComment
		}
		if bComment != nil {
			line = bytes.TrimLeft(line, string(bComment))
			if comment.Len() > 0 {
				comment.WriteByte('\n')             // one comment append '\n'
			}
			comment.Write(line)
			continue
		}

		// parse section, and create section map[string]string, maybe section---->comment
		if bytes.HasPrefix(line, sectionStart) && bytes.HasSuffix(line, sectionEnd) {
			section = strings.ToLower(string(line[1: len(line)-1]))
			if comment.Len() > 0 {
				cfgContainer.sectionComment[section] = comment.String()
				comment.Reset()
			}

			if _, ok := cfgContainer.data[section]; !ok {
				cfgContainer.data[section] = make(map[string]string)
			}
			continue
		}
		if _, ok := cfgContainer.data[section]; !ok {
			cfgContainer.data[section] = make(map[string]string)
		}
		keyValue := bytes.SplitN(line, bEqual, 2)
		key := strings.ToLower(string(bytes.TrimSpace(keyValue[0])))        // get key

		// maybe key is 'include', eg: include "other.conf"
		if len(keyValue) == 1 && strings.HasSuffix(key, includeSec) {
			includeFiles := strings.Fields(key)
			if includeFiles[0] == includeSec && len(includeFiles) == 2 {
				otherFile := strings.Trim(includeFiles[1], "\"")
				if !filepath.IsAbs(otherFile) {                         // 相对路径，需要join
					otherFile = filepath.Join(dir, otherFile)
				}

				innerC, err := c.parseFile(otherFile)
				if err != nil {
					return nil, err
				}

				// include file's section key value
				for sec, dt := range innerC.data {
					if _, ok := cfgContainer.data[sec]; !ok {
						cfgContainer.data[sec] = make(map[string]string)
					}
					for k, v := range dt {
						cfgContainer.data[sec][k] = v
					}
				}

				// include file's section comment
				for sec, comm := range innerC.sectionComment {
					cfgContainer.sectionComment[sec] = comm
				}

				// include file's key comment
				for k, comm := range innerC.keyComment {
					cfgContainer.keyComment[k] = comm
				}

				continue
			}
		}

		// not include deal continue
		if len(keyValue) != 2 {
			return nil, fmt.Errorf("read the content error: '%s', should key = val", string(line))
		}
		val := bytes.TrimSpace(keyValue[1])
		if bytes.HasPrefix(val, bDQuote) {
			val = bytes.Trim(val, `"`)
		}

		cfgContainer.data[section][key] = ExpandValueEnv(string(val))
		if comment.Len() > 0 {
			cfgContainer.keyComment[section + "." + key] = comment.String()
			comment.Reset()
		}
	}

	return cfgContainer, nil
}

/**
 对外提供的解析接口
 */
func (c *InitConfig) Parse(fileName string) (ConfContainer, error) {
	return c.parseFile(fileName)
}

func (c *InitConfig) ParseData(data []byte) (ConfContainer, error) {
	dir := "ini"
	currentUser, err := user.Current()
	if err == nil {
		dir = dir + "-" + currentUser.Username
	}
	dir = filepath.Join(os.TempDir(), dir)
	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	return c.parseData(dir, data)
}


/***********************************************************************
 解析之后的容器
 */
type IniConfigContainer struct {
	data            map[string]map[string]string            // section=> key:val        最多就是两层
	sectionComment  map[string]string                       // section: comment, 段的注释信息
	keyComment      map[string]string                       // id
	sync.RWMutex                            // 🔐机制，匿名成员，继承的机制
}

/**
 section.key or key
 内部接口
*/
func (c *IniConfigContainer) getData(key string) string {
	if len(key) == 0 {
		return ""
	}

	c.RLock()
	defer c.RUnlock()

	var (
		section, k string
		sectionKey = strings.Split(strings.ToLower(key), sepSecKey)
	)
	if len(sectionKey) >= 2 {
		section = sectionKey[0]
		k = sectionKey[1]
	} else {
		section = defaultSection
		k = sectionKey[0]
	}

	if v, ok := c.data[section]; ok {
		if vv, ok := v[k]; ok {
			return vv
		}
	}
	return ""
}


// Bool returns the boolean value for a given key.
func (c *IniConfigContainer) Bool(key string) (bool, error) {
	return ParseBool(c.getData(key))
}

// DefaultBool returns the boolean value for a given key.
// if err != nil return defaltval
func (c *IniConfigContainer) DefaultBool(key string, defaultVal bool) bool {
	v, err := c.Bool(key)
	if err != nil {
		return defaultVal
	}
	return v
}

// Int returns the integer value for a given key.
func (c *IniConfigContainer) Int(key string) (int, error) {
	return strconv.Atoi(c.getData(key))
}

// DefaultInt returns the integer value for a given key.
// if err != nil return defaltval
func (c *IniConfigContainer) DefaultInt(key string, defaultVal int) int {
	v, err := c.Int(key)
	if err != nil {
		return defaultVal
	}
	return v
}

// Int64 returns the int64 value for a given key.
func (c *IniConfigContainer) Int64(key string) (int64, error) {
	return strconv.ParseInt(c.getData(key), 10, 64)
}

// DefaultInt64 returns the int64 value for a given key.
// if err != nil return defaltval
func (c *IniConfigContainer) DefaultInt64(key string, defaultval int64) int64 {
	v, err := c.Int64(key)
	if err != nil {
		return defaultval
	}
	return v
}

// Float returns the float value for a given key.
func (c *IniConfigContainer) Float(key string) (float64, error) {
	return strconv.ParseFloat(c.getData(key), 64)
}

// DefaultFloat returns the float64 value for a given key.
// if err != nil return defaltval
func (c *IniConfigContainer) DefaultFloat(key string, defaultval float64) float64 {
	v, err := c.Float(key)
	if err != nil {
		return defaultval
	}
	return v
}

// String returns the string value for a given key.
func (c *IniConfigContainer) String(key string) string {
	return c.getData(key)
}

// DefaultString returns the string value for a given key.
// if err != nil return defaltval
func (c *IniConfigContainer) DefaultString(key string, defaultval string) string {
	v := c.String(key)
	if v == "" {
		return defaultval
	}
	return v
}

// Strings returns the []string value for a given key.
// Return nil if config value does not exist or is empty.
func (c *IniConfigContainer) Strings(key string) []string {
	v := c.String(key)
	if v == "" {
		return nil
	}
	return strings.Split(v, ";")
}

// DefaultStrings returns the []string value for a given key.
// if err != nil return defaltval
func (c *IniConfigContainer) DefaultStrings(key string, defaultval []string) []string {
	v := c.Strings(key)
	if v == nil {
		return defaultval
	}
	return v
}

// GetSection returns map for the given section
func (c *IniConfigContainer) GetSection(section string) (map[string]string, error) {
	if v, ok := c.data[section]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("not exist section")
}

// SaveConfigFile save the config into file.
//
// BUG(env): The environment variable config item will be saved with real value in SaveConfigFile Funcation.
func (c *IniConfigContainer) SaveConfigFile(filename string) (err error) {
	// Write configuration file by filename.
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Get section or key comments. Fixed #1607
	getCommentStr := func(section, key string) string {
		var (
			comment string
			ok      bool
		)
		if len(key) == 0 {
			comment, ok = c.sectionComment[section]
		} else {
			comment, ok = c.keyComment[section+"."+key]
		}

		if ok {
			// Empty comment
			if len(comment) == 0 || len(strings.TrimSpace(comment)) == 0 {
				return string(bNumComment)
			}
			prefix := string(bNumComment)
			// Add the line head character "#"
			return prefix + strings.Replace(comment, lineBreak, lineBreak+prefix, -1)
		}
		return ""
	}

	buf := bytes.NewBuffer(nil)
	// Save default section at first place
	if dt, ok := c.data[defaultSection]; ok {
		for key, val := range dt {
			if key != " " {
				// Write key comments.
				if v := getCommentStr(defaultSection, key); len(v) > 0 {
					if _, err = buf.WriteString(v + lineBreak); err != nil {
						return err
					}
				}

				// Write key and value.
				if _, err = buf.WriteString(key + string(bEqual) + val + lineBreak); err != nil {
					return err
				}
			}
		}

		// Put a line between sections.
		if _, err = buf.WriteString(lineBreak); err != nil {
			return err
		}
	}
	// Save named sections
	for section, dt := range c.data {
		if section != defaultSection {
			// Write section comments.
			if v := getCommentStr(section, ""); len(v) > 0 {
				if _, err = buf.WriteString(v + lineBreak); err != nil {
					return err
				}
			}

			// Write section name.
			if _, err = buf.WriteString(string(sectionStart) + section + string(sectionEnd) + lineBreak); err != nil {
				return err
			}

			for key, val := range dt {
				if key != " " {
					// Write key comments.
					if v := getCommentStr(section, key); len(v) > 0 {
						if _, err = buf.WriteString(v + lineBreak); err != nil {
							return err
						}
					}

					// Write key and value.
					if _, err = buf.WriteString(key + string(bEqual) + val + lineBreak); err != nil {
						return err
					}
				}
			}

			// Put a line between sections.
			if _, err = buf.WriteString(lineBreak); err != nil {
				return err
			}
		}
	}
	_, err = buf.WriteTo(f)
	return err
}

// Set writes a new value for key.
// if write to one section, the key need be "section::key".
// if the section is not existed, it panics.
func (c *IniConfigContainer) Set(key, value string) error {
	c.Lock()
	defer c.Unlock()
	if len(key) == 0 {
		return errors.New("key is empty")
	}

	var (
		section, k string
		sectionKey = strings.Split(strings.ToLower(key), "::")
	)

	if len(sectionKey) >= 2 {
		section = sectionKey[0]
		k = sectionKey[1]
	} else {
		section = defaultSection
		k = sectionKey[0]
	}

	if _, ok := c.data[section]; !ok {
		c.data[section] = make(map[string]string)
	}
	c.data[section][k] = value
	return nil
}