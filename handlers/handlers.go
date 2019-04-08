package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/capossele/zmq-backend/dao"
	"github.com/capossele/zmq-backend/models"

	"github.com/gorilla/mux"
)

var txs []models.Tx

// GetTxEndpoint gets a tx
func GetTxEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	payload := dao.GetTx(params["hash"])

	json.NewEncoder(w).Encode(payload)
	//return

	//json.NewEncoder(w).Encode("Tx not found")
}

// GetAllTxsEndpoint gets all txs
func GetAllTxsEndpoint(w http.ResponseWriter, r *http.Request) {
	payload := dao.GetAllTxs()
	json.NewEncoder(w).Encode(payload)
}

// CreateTxEndpoint creates a tx
func CreateTxEndpoint(w http.ResponseWriter, r *http.Request) {
	var tx models.Tx
	_ = json.NewDecoder(r.Body).Decode(&tx)
	dao.InsertOneValue(tx)
	json.NewEncoder(w).Encode(tx)
}

// // DeleteTxEndpoint deletes a tx
// func DeleteTxEndpoint(w http.ResponseWriter, r *http.Request) {
// 	var tx models.Tx
// 	_ = json.NewDecoder(r.Body).Decode(&tx)
// 	dao.DeleteTx(tx)
// }

// DeleteAllTxsEndpoints deletes a tx
func DeleteAllTxsEndpoints(w http.ResponseWriter, r *http.Request) {
	dao.DeleteAllTxs()
}
