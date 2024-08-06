package functions

import (
	"context"
	"encoding/json"
	"time"

	"urbangrid.com/database"
)

type model struct {
	Id     string `json:"id,omitempty"`
	Count  int64  `json:"count,omitempty"`
	Status string `json:"status,omitempty"`
}

func SetSignalById(id string, count int64, status string) error {
	data, err := json.Marshal(map[string]interface{}{
		"count":  count,
		"status": status,
	})
	if err != nil {
		return err
	}

	redisData := database.RedisClient.Set(context.Background(), id, data, 0)
	if redisData.Err() != nil {
		return redisData.Err()
	}

	return nil
}

func GetSignalById(id string) (model, error) {
	var data model

	val, err := database.RedisClient.Get(context.Background(), id).Result()
	if err != nil {
		return data, err
	}

	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		return data, err
	}

	data.Id = id

	return data, nil
}

func UpdateSignalData(id string, count int64, status string) error {
	err := SetSignalById(id, count, status)

	if err != nil {
		return err
	}

	return nil
}

func UpdateCount(id string, count int64) error {
	data, err := GetSignalById(id)

	if err != nil && err.Error() != "redis: nil" {
		return err
	}

	data.Count = count

	err = SetSignalById(id, data.Count, data.Status)

	return err
}

func UpdateStatus(id string, status string) error {
	data, err := GetSignalById(id)

	if err != nil && err.Error() != "redis: nil" {
		return err
	}

	err = SetSignalById(id, data.Count, status)

	return err
}

func GetSignalData() ([]model, error) {
	data := []model{}

	keys, err := database.RedisClient.Keys(context.Background(), "*").Result()
	if err != nil {
		return data, err
	}

	for _, key := range keys {
		if key == "timer" {
			continue
		}
		val, err := database.RedisClient.Get(context.Background(), key).Result()
		if err != nil {
			return data, err
		}

		var temp model
		err = json.Unmarshal([]byte(val), &temp)
		if err != nil {
			return data, err
		}

		temp.Id = key
		data = append(data, temp)
	}

	return data, nil
}

func TurnOnSignal(signal model) error {
	err := UpdateStatus(signal.Id, "on")

	if err != nil {
		return err
	}

	return nil
}

func TurnOffSignal(signal model) error {
	err := UpdateStatus(signal.Id, "off")

	if err != nil {
		return err
	}

	return nil
}

const timeLayout = time.RFC3339 // Standard time layout for storing time in Redis

type timer struct {
	Timer            string `json:"timer,omitempty"`
	LastActiveSignal string `json:"lastActiveSignal,omitempty"`
}

// SetSignalTimer sets the signal timer in Redis
func SetSignalTimer(t time.Time, lastActiveSignal string) error {
	timeBytes, err := json.Marshal(timer{
		Timer:            t.Format(timeLayout),
		LastActiveSignal: lastActiveSignal,
	})

	if err != nil {
		return err
	}

	_, err = database.RedisClient.Set(context.Background(), "timer", timeBytes, 0).Result()

	if err != nil {
		return err
	}

	return nil
}

// GetSignalTimer retrieves the signal timer from Redis
func GetSignalTimer() (timer, error) {
	val, err := database.RedisClient.Get(context.Background(), "timer").Result()
	if err != nil {
		return timer{}, err
	}

	var t timer

	err = json.Unmarshal([]byte(val), &t)
	if err != nil {
		return timer{}, err
	}

	return t, nil
}
