package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "quickmetrics/backend/handlers"
)

func main() {
    // Configuração do servidor Gin
    r := gin.Default()

    // Endpoint para calcular as métricas
    r.GET("/metrics", handlers.CalculateMetrics)

    // Executar o servidor na porta 8080
    r.Run(":8080")
}
