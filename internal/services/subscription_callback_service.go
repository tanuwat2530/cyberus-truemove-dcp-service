package services

import (
	"CyberusGolangShareLibrary/redis_db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Struct to map the expected JSON fields
type SubscriptionCallbackRequest struct {
	Msisdn    string `json:"msisdn"`
	Shortcode string `json:"short-code"`
	Operator  string `json:"operator"`
	TranRef   string `json:"tran-ref"`
	Action    string `json:"action"`
	Code      string `json:"code"`
	Desc      string `json:"desc"`
	RefId     string `json:"ref-id"`
	Media     string `json:"media"`
	Token     string `json:"token"`
	Timestamp int    `json:"timestamp"`
}

func SubscriptionCallbackProcessRequest(r *http.Request) map[string]string {

	// Generate a random UUID (UUID v4)
	//transaction_id := uuid.New().String()

	res := map[string]string{}

	var payload map[string]interface{}
	errPayload := json.NewDecoder(r.Body).Decode(&payload)
	if errPayload != nil {
		// Example: print the values
		//fmt.Println("Error decode Json to map[string]interface{} :", errPayload.Error())

		res["code"] = "-1"
		res["desc"] = "JSON Decode error"
		res["ref-id"] = "undefined"
		res["tran-ref"] = "undefined"
		res["media"] = "undefined"
		return res
	}

	jsonData, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		//fmt.Println("Error marshalling JSON:", err.Error())
		res["code"] = "-2"
		res["desc"] = "JSON Marshalling error"
		res["ref-id"] = "undefined"
		res["tran-ref"] = "undefined"
		res["media"] = "undefined"
		return res
	}

	// // Unmarshal JSON into struct
	var requestData SubscriptionCallbackRequest
	err = json.Unmarshal(jsonData, &requestData)
	if err != nil {
		//fmt.Println("Error map Json to Struct :" + err.Error())
		//fmt.Println("Error marshalling JSON:", err.Error())
		res["code"] = "-3"
		res["desc"] = "JSON Struct error"
		res["ref-id"] = "undefined"
		res["tran-ref"] = "undefined"
		res["media"] = "undefined"
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
		res["ref-id"] = requestData.RefId
		res["tran-ref"] = requestData.TranRef
		res["media"] = requestData.Media
		return res
	}

	payloadString := string(payloadBytes)
	redis_key := "tmvh-subscription-callback-api:" + requestData.Media + ":" + requestData.RefId
	ttl := 24 * time.Hour // expires in 1 Hour

	redis_db.ConnectRedis()
	// // Set key with TTL
	if err := redis_db.SetWithTTL(redis_key, payloadString, ttl); err != nil {
		//write to file if Redis problem or forward request to AIS
		log.Fatalf("SetWithTTL error: %v", err)
		fmt.Println("Error marshalling JSON:", err.Error())
		res["code"] = "-5"
		res["desc"] = "Save process error"
		res["ref-id"] = requestData.RefId
		res["tran-ref"] = requestData.TranRef
		res["media"] = requestData.Media
		return res
	}

	//fmt.Println("##### Process completed #####")
	res["code"] = "200"
	res["desc"] = "Success"
	res["ref-id"] = requestData.RefId
	res["tran-ref"] = requestData.TranRef
	res["media"] = requestData.Media

	return res
}
