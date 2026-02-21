package controller

import (
	"autoscaling-hetzner/model"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type target struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

func GetTargets(g *gin.Context) {
	var targets []target
	rows, err := model.GetAllServers()
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Print(err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var ip net.IP
		var groupId int
		if err := rows.Scan(&name, &groupId, &ip); err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Print(err.Error())
		}
		targets = append(targets, target{
			Targets: []string{ip.String() + ":9100"},
			Labels:  map[string]string{"groupId": strconv.Itoa(groupId), "name": name},
		})
	}
	if err := rows.Err(); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, targets)
}
