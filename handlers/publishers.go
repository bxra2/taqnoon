// handlers/publishers.go
package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/bxra2/taqnun/models"

	"github.com/bxra2/taqnun/db"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func GetPublishers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	col := db.GetCollection("mydb", "publishers")
	cursor, err := col.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var publishers []models.Publisher
	if err = cursor.All(ctx, &publishers); err != nil {
		http.Error(w, "Decode error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(publishers)
}
