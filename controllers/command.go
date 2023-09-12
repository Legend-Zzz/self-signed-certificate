package controllers

import (
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type AllFormData struct {
	Selection        string `form:"selection"`
	Domain           string `form:"domain"`
	IP               string `form:"ip"`
	CASubject        string `form:"caSubject"`
	CAValidityDays   int    `form:"caValidityDays"`
	RootSubject      string `form:"rootSubject"`
	RootValidityDays int    `form:"rootValidityDays"`
	SerialNumber     int    `form:"serialNumber"`
}

type RootFormData struct {
	Selection            string `form:"selection"`
	RootOnlySubject      string `form:"rootOnlySubject"`
	RootOnlyValidityDays int    `form:"rootOnlyValidityDays"`
	RootOnlySerialNumber int    `form:"rootOnlySerialNumber"`
}

func generateCommand(data interface{}) (string, error) {
	var cmd string
	switch v := data.(type) {
	case *RootFormData:
		cmd = "./gen.root.sh"
		if v.RootOnlySubject != "" {
			cmd += " -s " + v.RootOnlySubject
		}
		if v.RootOnlyValidityDays != 0 {
			cmd += " -d " + strconv.Itoa(v.RootOnlyValidityDays)
		}
		if v.RootOnlySerialNumber != 0 {
			cmd += " -sn " + strconv.Itoa(v.RootOnlySerialNumber)
		}
	case *AllFormData:
		cmd = "./gen.cert.sh"
		if v.Domain != "" {
			cmd += " -d " + v.Domain
		}
		if v.IP != "" {
			cmd += " -i " + v.IP
		}
		if v.CASubject != "" {
			cmd += " -s " + v.CASubject
		}
		if v.CAValidityDays != 0 {
			cmd += " -D " + strconv.Itoa(v.CAValidityDays)
		}
		if v.RootSubject != "" {
			cmd += " -rs " + v.RootSubject
		}
		if v.RootValidityDays != 0 {
			cmd += " -rD " + strconv.Itoa(v.RootValidityDays)
		}
		if v.SerialNumber != 0 {
			cmd += " -sn " + strconv.Itoa(v.SerialNumber)
		}
	}

	return cmd, nil
}

func executeCommand(command string) (string, error) {
	cmdArgs := strings.Fields(command)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	output, err := cmd.CombinedOutput()
	result := strings.TrimSpace(string(output))

	return result, err
}

func ExecCommand(c *gin.Context, fileName string) {
	selection := c.PostForm("selection")
	var formData interface{}

	if selection == "rootOnly" {
		formData = &RootFormData{}
	} else if selection == "all" {
		formData = &AllFormData{}
	}

	if formData == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid formData"})
		return
	}

	if err := c.ShouldBind(formData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	command, err := generateCommand(formData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result, cmdErr := executeCommand(command)
	result = "------------------------------" + formattedTime + "------------------------------" + "\n" + command + "\n" + result + "\n"

	WriteFiles(c, result, fileName)

	if cmdErr != nil {
		c.Redirect(http.StatusSeeOther, "/result")
		return
	}

	c.Redirect(http.StatusSeeOther, "/files")
}
