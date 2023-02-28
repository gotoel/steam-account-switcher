package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/andygrunwald/vdf"
	"os"
	"strconv"
	"strings"
)

func switchAccountLoginUsersVdf(settings *Settings, account *Account) error {
	loginUsersVdfPath := fmt.Sprintf("%s\\config\\loginusers.vdf",
		strings.TrimRight(strings.TrimRight(settings.SteamPath, "\\"), "/"))
	// Get login user map
	loginUsersData, err := os.ReadFile(loginUsersVdfPath)

	p := vdf.NewParser(bytes.NewReader(loginUsersData))
	loginUsers, err := p.Parse()
	if err != nil {
		return err
	}

	if err == nil {
		for _, usersMap := range loginUsers {
			if users, ok := usersMap.(map[string]interface{}); ok {
				for _, userMap := range users {
					if user, okUser := userMap.(map[string]interface{}); okUser {
						if user["AccountName"] == account.Username {
							user["AllowAutoLogin"] = "1"
							user["RememberPassword"] = "1"
							user["MostRecent"] = "1"
						} else {
							user["AllowAutoLogin"] = "0"
							user["RememberPassword"] = "0"
							user["MostRecent"] = "0"
						}
					}
				}
			}
		}

		return writeLoginUsersVdf(loginUsersVdfPath, loginUsers)
	} else {
		return fmt.Errorf("error retrieving login users vdf: %s", err)
	}

	return nil
}

func writeLoginUsersVdf(filePath string, loginUsers map[string]interface{}) error {
	var mapBuffer bytes.Buffer
	writer := bufio.NewWriter(&mapBuffer)
	buildVdfMap(writer, loginUsers, 0)
	err := writer.Flush()
	if err == nil {
		var file *os.File
		file, err = os.Create(filePath)
		if err != nil {
			return err
		} else {
			_, err = file.Write(mapBuffer.Bytes())
			return err
		}
	} else {
		return err
	}
	return nil
}

func buildVdfMap(writer *bufio.Writer, m map[string]interface{}, indent int) {
	for k, v := range m {
		// Write the key with the appropriate indentation
		for i := 0; i < indent; i++ {
			writer.WriteString("\t")
		}
		writer.WriteString(`"` + k + `"`)

		// Write the value
		switch value := v.(type) {
		case int:
			writer.WriteString("\t\t\"" + strconv.Itoa(value) + "\"\n")
		case string:
			writer.WriteString("\t\t\"" + value + "\"\n")
		case []int:
			writer.WriteString(" \n")
			for _, i := range value {
				for i := 0; i < indent+1; i++ {
					writer.WriteString("\t")
				}
				writer.WriteString(strconv.Itoa(i) + "\n")
			}
		case map[string]interface{}:
			writer.WriteString(" \n")
			for i := 0; i < indent; i++ {
				writer.WriteString("\t")
			}
			writer.WriteString("{")
			writer.WriteString(" \n")
			buildVdfMap(writer, value, indent+1)
			for i := 0; i < indent; i++ {
				writer.WriteString("\t")
			}
			writer.WriteString("}")
			writer.WriteString(" \n")
		}
	}
}
