package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type RegistrationCertificate struct {
	AssetType       string `json:"assetType"`
	RCId            string `json:"rcId"`
	Make            string `json:"make"`
	Model           string `json:"model"`
	Color           string `json:"color"`
	OwnerName       string `json:"ownerName"`
	OwnerAadhar     string `json:"ownerAadhar"`
	EngineNumber    string `json:"engineNumber"`
	InsuranceCert   string `json:"insuranceCert"`
	PollutionCert   string `json:"pollutionCert"`
	Status          string `json:"status"`
	RegisteredState string `json:"registeredState"`
}

type RCData struct {
	AssetType       string `json:"assetType"`
	RCId            string `json:"RCId"`
	Make            string `json:"make"`
	Model           string `json:"model"`
	Color           string `json:"color"`
	IsNOCapproved   bool   `json:"isNOCapproved"`
	RegisteredState string `json:"registeredState"`
	Status          string `json:"status"`
}

type TransferRC struct {
	RCId      string `json:"RCId"`
	FromState string `json:"fromState"`
	ToState   string `json:"toState"`
}
type RCHistory struct {
	Record    *RCData `json:"record"`
	TxId      string  `json:"txId"`
	Timestamp string  `json:"timestamp"`
	IsDelete  bool    `json:"isDelete"`
}
type PaginatedQueryResult struct {
	Records             []*RegistrationCertificate `json:"records"`
	FetchedRecordsCount int32                      `json:"fetchedRecordsCount"`
	Bookmark            string                     `json:"bookmark"`
}
type TransferRCData struct {
	RCId            string `json:"rcId"`
	AssetType       string `json:"assetType"`
	Make            string `json:"make"`
	Model           string `json:"model"`
	Color           string `json:"color"`
	OwnerName       string `json:"ownerName"`
	OwnerAadhar     string `json:"ownerAadhar"`
	EngineNumber    string `json:"engineNumber"`
	InsuranceCert   string `json:"insuranceCert"`
	PollutionCert   string `json:"pollutionCert"`
	Status          string `json:"status"`
	RegisteredState string `json:"registeredState"`
	TransferredFrom string `json:"transferredFrom"`
	TransferredTo   string `json:"transferredTo"`
}

func main() {
	router := gin.Default()
	var wg sync.WaitGroup
	wg.Add(1)
	go ChaincodeEventListener("tnrto", "vaahanchannel", "autoledger", &wg)
	router.Static("/public", "./public")
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(ctx *gin.Context) {
		result := submitTxnFn("tnrto", "vaahanchannel", "autoledger", "RegistrationCertificateContract", "query", make(map[string][]byte), "GetAllRCs")

		var rcs []RCData

		if len(result) > 0 {
			// Unmarshal the JSON array string into the cars slice
			if err := json.Unmarshal([]byte(result), &rcs); err != nil {
				fmt.Println("Error:", err)
				return
			}
		}

		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"title": "AutoLedger App", "rcList": rcs,
		})
	})

	router.GET("/tnrto", func(ctx *gin.Context) {
		result := submitTxnFn("tnrto", "vaahanchannel", "autoledger", "RegistrationCertificateContract", "query", make(map[string][]byte), "GetAllRCs")

		var rcs []RCData

		if len(result) > 0 {
			// Unmarshal the JSON array string into the cars slice
			if err := json.Unmarshal([]byte(result), &rcs); err != nil {
				fmt.Println("Error:", err)
				return
			}
		}

		ctx.HTML(http.StatusOK, "tnrto.html", gin.H{
			"title": "TamilnaduRTO Dashboard", "RCList": rcs,
		})
	})

	router.GET("/klrto", func(ctx *gin.Context) {
		result := submitTxnFn("klrto", "vaahanchannel", "autoledger", "RegistrationCertificateContract", "query", make(map[string][]byte), "GetAssetsbyState", "klrto")

		var rcs []RCData

		if len(result) > 0 {
			// Unmarshal the JSON array string into the cars slice
			if err := json.Unmarshal([]byte(result), &rcs); err != nil {
				fmt.Println("Error:", err)
				return
			}
		}

		ctx.HTML(http.StatusOK, "klrto.html", gin.H{
			"title": "KeralaRTO Dashboard", "RCList": rcs,
		})
	})

	router.POST("/api/rc/create", func(ctx *gin.Context) {
		var req RegistrationCertificate
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
			return
		}
		fmt.Printf("Registration Certificate Details  %s", req)
		// Convert the request into a transient data map format
		rcData := RegistrationCertificate{
			AssetType:       req.AssetType,
			RCId:            req.RCId,
			Make:            req.Make,
			Model:           req.Model,
			Color:           req.Color,
			OwnerName:       req.OwnerName,
			OwnerAadhar:     req.OwnerAadhar,
			EngineNumber:    req.EngineNumber,
			InsuranceCert:   req.InsuranceCert,
			PollutionCert:   req.PollutionCert,
			RegisteredState: req.RegisteredState,
		}
		// Marshal rcData to JSON
		rcJSON, err := json.Marshal(rcData)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to marshal RC data"})
			return
		}

		// Store the JSON in privateData
		privateData := map[string][]byte{
			"rc_properties": rcJSON,
		}

		fmt.Printf("Registration Certificate Details: %s\n", privateData)

		submitTxnFn(req.RegisteredState, "vaahanchannel", "autoledger", "PrivateAssetDetailsContract", "private", privateData, "CreatePrivateRC", req.RegisteredState)
		submitTxnFn(req.RegisteredState, "vaahanchannel", "autoledger", "RegistrationCertificateContract", "invoke", make(map[string][]byte), "CreateRC", req.RCId, req.AssetType, req.Make, req.Model, req.Color, req.RegisteredState)

		ctx.JSON(http.StatusOK, req)
	})
	router.GET("/api/event", func(ctx *gin.Context) {
		result := getEvents()
		fmt.Println("result:", result)

		ctx.JSON(http.StatusOK, gin.H{"RCEvent": result})

	})
	router.POST("api/rc/initiateTransfer", func(ctx *gin.Context) {
		fmt.Printf("Inside Initiate transfer")
		var req TransferRC
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
			return
		}
		fmt.Printf("Initiate RC Transfer  %s", req)
		submitTxnFn(req.FromState, "vaahanchannel", "autoledger", "RegistrationCertificateContract", "invoke", make(map[string][]byte), "InitiateTransferRC", req.RCId, req.FromState, req.ToState)

		ctx.JSON(http.StatusOK, req)

	})
	router.POST("api/rc/approveTransfer", func(ctx *gin.Context) {
		var req TransferRC
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
			return
		}
		fmt.Printf("Approve RC Transfer  %s", req)
		submitTxnFn(req.ToState, "vaahanchannel", "autoledger", "RegistrationCertificateContract", "invoke", make(map[string][]byte), "ApproveTransferRC", req.RCId, req.FromState, req.ToState)

		ctx.JSON(http.StatusOK, req)

	})
	router.POST("api/rc/deleterc", func(ctx *gin.Context) {
		var req TransferRC
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
			return
		}
		fmt.Printf("Delete RC   %s", req)
		submitTxnFn(req.FromState, "vaahanchannel", "autoledger", "PrivateAssetDetailsContract", "invoke", make(map[string][]byte), "MovetoTransferRCCollection", req.RCId, req.FromState, req.ToState)
		submitTxnFn(req.FromState, "vaahanchannel", "autoledger", "RegistrationCertificateContract", "invoke", make(map[string][]byte), "DeleteRC", req.RCId, req.FromState, req.ToState)

		ctx.JSON(http.StatusOK, req)

	})
	router.POST("/api/rc/addrc", func(ctx *gin.Context) {
		var req RegistrationCertificate
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
			return
		}
		fmt.Printf("Registration Certificate Details  %s", req)
		// Convert the request into a transient data map format
		addrcData := RegistrationCertificate{
			AssetType:       req.AssetType,
			RCId:            req.RCId,
			Make:            req.Make,
			Model:           req.Model,
			Color:           req.Color,
			OwnerName:       req.OwnerName,
			OwnerAadhar:     req.OwnerAadhar,
			EngineNumber:    req.EngineNumber,
			InsuranceCert:   req.InsuranceCert,
			PollutionCert:   req.PollutionCert,
			RegisteredState: req.RegisteredState,
		}
		fmt.Printf("Add Transferred RC   %s", req)
		// Marshal rcData to JSON
		addrcJSON, err := json.Marshal(addrcData)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to marshal RC data"})
			return
		}

		// Store the JSON in privateData
		transferprivateData := map[string][]byte{
			"add_rc_properties": addrcJSON,
		}

		fmt.Printf("Registration Certificate Details: %s\n", transferprivateData)

		submitTxnFn(req.RegisteredState, "vaahanchannel", "autoledger", "PrivateAssetDetailsContract", "private", transferprivateData, "AddTransferredRC", "KL01CX2525", req.RegisteredState)
		submitTxnFn(req.RegisteredState, "vaahanchannel", "autoledger", "RegistrationCertificateContract", "invoke", make(map[string][]byte), "CreateRC", "KL01CX2525", req.AssetType, req.Make, req.Model, req.Color, req.RegisteredState)

		ctx.JSON(http.StatusOK, req)
	})
	router.GET("/api/tnrto/rc/all", func(ctx *gin.Context) {

		result := submitTxnFn("tnrto", "vaahanchannel", "autoledger", "PrivateAssetDetailsContract", "query", make(map[string][]byte), "GetAssetCollection", "tnrto")

		var rcs []RegistrationCertificate

		if len(result) > 0 {
			// Unmarshal the JSON array string into the orders slice
			if err := json.Unmarshal([]byte(result), &rcs); err != nil {
				fmt.Println("Error:", err)
				return
			}
		} else {
			fmt.Println("no result found")
		}

		ctx.HTML(http.StatusOK, "regcertificates.html", gin.H{
			"title": "All RCs registered with TNRTO", "RCList": rcs,
		})
	})

	router.GET("/api/klrto/rc/all", func(ctx *gin.Context) {

		result := submitTxnFn("klrto", "vaahanchannel", "autoledger", "PrivateAssetDetailsContract", "query", make(map[string][]byte), "GetAssetCollection", "klrto")

		var rcs []RegistrationCertificate

		if len(result) > 0 {
			// Unmarshal the JSON array string into the orders slice
			if err := json.Unmarshal([]byte(result), &rcs); err != nil {
				fmt.Println("Error:", err)
				return
			}
		} else {
			fmt.Println("no result found")
		}

		ctx.HTML(http.StatusOK, "regcertificates.html", gin.H{
			"title": "All RCs registered with KLRTO", "RCList": rcs,
		})
	})

	router.GET("/api/tnrto/transferred_rc/all", func(ctx *gin.Context) {

		result := submitTxnFn("tnrto", "vaahanchannel", "autoledger", "PrivateAssetDetailsContract", "query", make(map[string][]byte), "GetTransferredRCCollection", "tnrto")

		var trcs []TransferRCData

		if len(result) > 0 {
			// Unmarshal the JSON array string into the orders slice
			if err := json.Unmarshal([]byte(result), &trcs); err != nil {
				fmt.Println("Error:", err)
				return
			}
		} else {
			fmt.Println("no result found")
		}

		ctx.HTML(http.StatusOK, "transfercertificates.html", gin.H{
			"title": "All Transferred RCs", "RCList": trcs,
		})
	})

	router.GET("/api/klrto/transferred_rc/all", func(ctx *gin.Context) {

		result := submitTxnFn("klrto", "vaahanchannel", "autoledger", "PrivateAssetDetailsContract", "query", make(map[string][]byte), "GetTransferredRCCollection", "klrto")

		var trcs []TransferRCData

		if len(result) > 0 {
			// Unmarshal the JSON array string into the orders slice
			if err := json.Unmarshal([]byte(result), &trcs); err != nil {
				fmt.Println("Error:", err)
				return
			}
		} else {
			fmt.Println("no result found")
		}

		ctx.HTML(http.StatusOK, "transfercertificates.html", gin.H{
			"title": "All Transferred RCs", "RCList": trcs,
		})
	})

	router.GET("/api/tnrto/rc/:id", func(ctx *gin.Context) {
		rcId := ctx.Param("id")
		//organisation:=ctx.GetClientIdentity().GetMSPID()

		result := submitTxnFn("tnrto", "vaahanchannel", "autoledger", "RegistrationCertificateContract", "query", make(map[string][]byte), "GetAssetDetails", rcId)

		ctx.JSON(http.StatusOK, gin.H{"data": result})
	})

	router.GET("/api/klrto/rc/:id", func(ctx *gin.Context) {
		rcId := ctx.Param("id")
		//organisation:=ctx.GetClientIdentity().GetMSPID()

		result := submitTxnFn("klrto", "vaahanchannel", "autoledger", "RegistrationCertificateContract", "query", make(map[string][]byte), "GetAssetDetails", rcId)

		ctx.JSON(http.StatusOK, gin.H{"data": result})
	})

	router.GET("/api/rc/history", func(ctx *gin.Context) {
		rcID := ctx.Query("rcId")
		result := submitTxnFn("tnrto", "vaahanchannel", "autoledger", "RegistrationCertificateContract", "query", make(map[string][]byte), "GetAssetHistory", rcID)

		// fmt.Printf("result %s", result)

		var rcs []RCHistory

		if len(result) > 0 {
			// Unmarshal the JSON array string into the orders slice
			if err := json.Unmarshal([]byte(result), &rcs); err != nil {
				fmt.Println("Error:", err)
				return
			}
		}

		ctx.HTML(http.StatusOK, "history.html", gin.H{
			"title": "RC History", "itemList": rcs,
		})
	})

	router.GET("/api/rc/asset_by_range", func(ctx *gin.Context) {

		result := submitTxnFn("tnrto", "vaahanchannel", "autoledger", "RegistrationCertificateContract", "query", make(map[string][]byte), "GetAssetsByRange", "TN01CX2525", "TN04CX2525")

		fmt.Printf("result %s", result)

		var rcs []RCData

		if len(result) > 0 {
			// Unmarshal the JSON array string into the orders slice
			if err := json.Unmarshal([]byte(result), &rcs); err != nil {
				fmt.Println("Error:", err)
				return
			}
		}
		ctx.JSON(http.StatusOK, result)
	})

	router.GET("/api/rc/asset_with_pagination", func(ctx *gin.Context) {

		result := submitTxnFn("tnrto", "vaahanchannel", "autoledger", "RegistrationCertificateContract", "query", make(map[string][]byte), "GetAssetsWithPagination", "2", "")

		fmt.Printf("result %s", result)

		var rcs PaginatedQueryResult

		if len(result) > 0 {
			// Unmarshal the JSON array string into the orders slice
			if err := json.Unmarshal([]byte(result), &rcs); err != nil {
				fmt.Println("Error:", err)
				return
			}
		}
		ctx.JSON(http.StatusOK, result)
	})
	router.Run("localhost:8080")

}
