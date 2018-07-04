package doc2cvariable

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Converter struct {
	outFilePath string
	inFilePaths []string
}

func NewDoc2CVariable(out string, in []string) *Converter {
	c := &Converter{
		outFilePath: out,
		inFilePaths: in,
	}
	return c
}

func (c *Converter) getDefineName() string {
	name := strings.Replace(filepath.Base(c.outFilePath), ".", "_", -1)
	return strings.ToUpper(name) + "_"
}

func getVariableName(name string) string {
	return strings.Replace(filepath.Base(name), ".", "_", -1)
}

func replaceEOL(text string) string {
	reg := regexp.MustCompile(`\r\n|\r|\n`)
	return string(reg.Copy().ReplaceAll([]byte(text), []byte("\\n")))
}

func (c *Converter) createHeader() string {
	template := "#ifndef <+name+>\n#define <+name+>\n\n"
	return strings.Replace(template, "<+name+>", c.getDefineName(), -1)
}

func (c *Converter) createBody() (string, error) {
	template := "const char *<+name+> = \"<+doc+>\";\n"
	result := []string{}
	for _, name := range c.inFilePaths {
		text, err := ioutil.ReadFile(name)
		if err != nil {
			return "", err
		}

		replaced := strings.Replace(template, "<+name+>", getVariableName(name), -1)
		replaced = strings.Replace(replaced, "<+doc+>", replaceEOL(string(text)), -1)
		result = append(result, replaced)
	}
	return strings.Join(result[:], "\n"), nil
}

func (c *Converter) createFooter() string {
	return "\n#endif"
}

func (c *Converter) Convert() (string, error) {
	result := c.createHeader()
	body, err := c.createBody()
	if err != nil {
		return "", err
	}
	result += body
	result += c.createFooter()
	return result, nil
}

func (c *Converter) WriteFile() error {
	converted, err := c.Convert()
	if err != nil {
		return err
	}
	ioutil.WriteFile(c.outFilePath, []byte(converted), os.ModePerm)
	return nil
}
