package services

import (
	"fmt"

	pb "github.com/rtmelsov/metrigger/proto"

	"errors"
)

func UpdateMetric(resp *pb.Metric) (*pb.Metric, error) {
	var aliasErr error

	fmt.Println("resp.MType", resp.MType)

	switch resp.MType {
	case "counter":
		aliasErr = MetricsCounterSet(resp)
	case "gauge":
		aliasErr = MetricsGaugeSet(resp)
	default:
		return nil, errors.New("can't find type")
	}

	if aliasErr != nil {
		return nil, aliasErr
	}

	if resp.MType == "counter" {
		obj, _, err := MetricsCounterGet(resp.ID)
		if err != nil {
			return nil, errors.New("failed to find element")
		}
		resp.Delta = obj.Value
		return resp, nil
	} else {
		_, obj, err := MetricsGaugeGet(resp.ID)
		if err != nil {
			return nil, errors.New("failed to find element")
		}
		resp.Value = obj.Value
		return resp, nil
	}
}
