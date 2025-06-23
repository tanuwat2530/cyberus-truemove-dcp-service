package services

import (
	"CyberusGolangShareLibrary/redis_db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Struct to map the expected JSON fields
type TransactionCallbackRequest struct {
	Code      string `json:"code"`
	Desc      string `json:"desc"`
	Msisdn    string `json:"msisdn"`
	Operator  string `json:"operator"`
	Shortcode string `json:"short-code"`
	TranRef   string `json:"tran-ref"`
	Timestamp int    `json:"timestamp"`
}

func TransactionCallbackProcessRequest(r *http.Request) map[string]string {

	// Generate a random UUID (UUID v4)
	//transaction_id := uuid.New().String()
	redisConnection := os.Getenv("BN_REDIS_URL")
	res := map[string]string{}

	redis_key := ""

	var payload map[string]interface{}
	errPayload := json.NewDecoder(r.Body).Decode(&payload)
	if errPayload != nil {
		// Example: print the values
		//fmt.Println("Error decode Json to map[string]interface{} :", errPayload.Error())

		res["code"] = "-1"
		res["desc"] = "JSON Decode error"
		res["tran-ref"] = "undefined"
		return res
	}

	jsonData, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		//fmt.Println("Error marshalling JSON:", err.Error())
		res["code"] = "-2"
		res["desc"] = "JSON Marshalling error"
		res["tran-ref"] = "undefined"
		return res
	}

	// // Unmarshal JSON into struct
	var requestData TransactionCallbackRequest
	err = json.Unmarshal(jsonData, &requestData)
	if err != nil {
		//fmt.Println("Error map Json to Struct :" + err.Error())
		//fmt.Println("Error marshalling JSON:", err.Error())
		res["code"] = "-3"
		res["desc"] = "JSON Struct error"
		res["tran-ref"] = "undefined"
		return res
	}

	// // Print the data to the console
	// fmt.Println("##### Received Subscription Callback Data #####")
	// fmt.Println("IDPartner : " + requestData.Msisdn)
	// fmt.Println("Shortcode : " + requestData.Shortcode)
	// fmt.Println("Operator  : " + requestData.Operator)
	// fmt.Println("Action  : " + requestData.Action)
	// fmt.Println("Code  : " + requestData.Code)
	// fmt.Println("Desc  : " + requestData.Desc)
	// fmt.Println("Timestamp  : " + strconv.FormatInt(requestData.Timestamp))
	// fmt.Println("TranRef  : " + requestData.TranRef)
	// fmt.Println("Action  : " + requestData.Action)
	// fmt.Println("RefId  : " + requestData.RefId)
	// fmt.Println("Media  : " + requestData.Media)
	// fmt.Println("Token  : " + requestData.Token)

	// // Add ClientIP to Payload
	payload["cyberus-return"] = "200"
	// Convert the struct to JSON string
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Failed to convert payload to JSON:  %+v\n ", http.StatusInternalServerError)
		fmt.Println("Error marshalling JSON:", err.Error())
		res["code"] = "-4"
		res["desc"] = "Additional value error"
		res["tran-ref"] = requestData.TranRef
		return res
	}

	payloadString := string(payloadBytes)
	if requestData.Operator == "TRUEMOVE" {
		redis_key = "tmvh-transaction-callback-api:" + requestData.TranRef
	}
	if requestData.Operator == "DTAC" {
		redis_key = "dtac-transaction-callback-api:" + requestData.TranRef
	}
	if requestData.Operator == "AIS" {
		redis_key = "ais-transaction-callback-api:" + requestData.TranRef
	}
	ttl := 24 * time.Hour // expires in 1 Hour

	redis_db.ConnectRedis(redisConnection, "", 0)
	// // Set key with TTL
	if err := redis_db.SetWithTTL(redis_key, payloadString, ttl); err != nil {
		//write to file if Redis problem or forward request to AIS
		log.Fatalf("SetWithTTL error: %v", err)
		fmt.Println("Error marshalling JSON:", err.Error())
		res["code"] = "-5"
		res["desc"] = "Save process error"
		res["tran-ref"] = requestData.TranRef
		return res
	}

	//fmt.Println("##### Process completed #####")

	res["code"] = "200"
	res["desc"] = "Success"
	res["tran-ref"] = requestData.TranRef

	return res
}
