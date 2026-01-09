package api

import (
	"bytes"
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/gingray/swisstools/pkg/common"
	"net/http"
	"os/exec"
)

type Server struct {
	Config *common.Config
}

type commandReq struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Cwd     string   `json:"cwd"`
}

func CreateServer(config *common.Config) *Server {
	return &Server{Config: config}
}

func (s *Server) Health(c *gin.Context) {
	c.JSONP(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func (s *Server) ExecuteCMD(c *gin.Context) {
	var req commandReq
	err := c.BindJSON(&req)
	if err != nil {
		log.Error(err)
		c.JSONP(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var stdBuffer bytes.Buffer
	var stdErrBuffer bytes.Buffer
	cmd := exec.Command(req.Command, req.Args...)
	cmd.Dir = req.Cwd
	cmd.Stdout = &stdBuffer
	cmd.Stderr = &stdErrBuffer
	err = cmd.Run()
	if err != nil {
		log.Error(err)
		c.JSONP(http.StatusInternalServerError, gin.H{"execution error": err.Error()})
		return
	}
	log.Infof("Command %s executed successfully", req.Command)
	c.JSONP(http.StatusOK, gin.H{"message": "OK", "cmd": req.Command, "args": req.Args, "cwd": req.Cwd, "stdout": stdBuffer.String(), "stderr": stdErrBuffer.String()})
}
