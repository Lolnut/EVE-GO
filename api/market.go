package api

import (
	"EVE-GO/specifications"
	"EVE-GO/util"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

func (c *Client) Market() *MarketEndpoint {
	me := &MarketEndpoint{c:c}
	return me
}

type MarketEndpoint struct {
	c *Client
}

type MarketQuery struct {
	c           *Client
	Endpoint    string
	QueryParams map[string]string
	Spec        specifications.Specification
}

func (m *MarketQuery) Get(result interface{}) error {

	if m.Endpoint == "" || m.Spec == nil {
		return fmt.Errorf("%s", "The query endpoint and specification must be defined. Did you forget to select an endpoint?")
	}

	if !m.Spec.IsSatisfiedBy(m) {
		return fmt.Errorf("%s", "The required query parameters were not satisfied")
	}

	data, err := m.c.Get(context.Background(), m.Endpoint, m.QueryParams)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return err
	}
	return nil
}

//-- Do Not Require Authentication --

type AlwaysTrueSpecification struct {
	specifications.AbstractSpecification
}

func (s *AlwaysTrueSpecification) IsSatisfiedBy(object interface{}) bool {
	return true
}

type MarketPrice struct {
	AdjustedPrice float32 `json:"adjusted_price"`
	AveragePrice float32 `json:"average_price"`
	TypeId int `json:"type_id"`
}

func (m *MarketEndpoint) MarketPrices() *MarketQuery {
	mq := &MarketQuery{c:m.c, QueryParams:make(map[string]string)}
	mq.Endpoint = "/markets/prices/"
	mq.Spec = &AlwaysTrueSpecification{}
	return mq
}

type ItemHistorySpecification struct {
	specifications.AbstractSpecification
}

func (s *ItemHistorySpecification) IsSatisfiedBy(object interface{}) bool {
	s2 := HasTypeIdParameter{}
	return s2.IsSatisfiedBy(object)
}

type ItemHistoryDatapoint struct {
	Average float32 `json:"average"`
	Date string `json:"date"`
	Highest float32 `json:"highest"`
	Lowest float32 `json:"lowest"`
	OrderCount int `json:"order_count"`
	Volume int `json:"volume"`
}

func (m *MarketEndpoint) ItemRegionHistory(regionId int, params ...util.KVP) *MarketQuery {
	mq := &MarketQuery{c:m.c, QueryParams:make(map[string]string)}
	mq.Endpoint = fmt.Sprintf("/markets/%d/history/", regionId)
	mq.Spec = &ItemHistorySpecification{}

	for _, param := range params {
		mq.QueryParams[param.Key] = param.Value.(string)
	}

	return mq
}

type MarketOrder struct {
	Duration int `json:"duration"`
	IsBuyOrder bool `json:"is_buy_order"`
	Issued string `json:"issued"`
	LocationID int64 `json:"location_id"`
	MinVolume int `json:"min_volume"`
	OrderID int64 `json:"order_id"`
	Price float32 `json:"price"`
	Range string `json:"range"`
	SystemID int `json:"system_id"`
	TypeID int `json:"type_id"`
	VolumeRemain int `json:"volume_remain"`
	VolumeTotal int `json:"volume_total"`
}

type RegionOrdersSpecification struct {
	specifications.AbstractSpecification
}

func (s *RegionOrdersSpecification) IsSatisfiedBy(object interface{}) bool {
	s2 := HasOrderTypeParameter{}
	return s2.IsSatisfiedBy(object)
}

func (m *MarketEndpoint) RegionOrders(regionId int, params ...util.KVP) *MarketQuery {
	mq := &MarketQuery{c:m.c, QueryParams:make(map[string]string)}
	mq.Endpoint = fmt.Sprintf("/markets/%d/orders/", regionId)
	mq.Spec = &AlwaysTrueSpecification{}

	for _, param := range params {
		mq.QueryParams[param.Key] = param.Value.(string)
	}

	return mq
}

// -- Requires Authentication --

// -- Predefined Parameter Functions --

func TypeId(id int) util.KVP {
	return util.KVP{Key:"type_id", Value:strconv.Itoa(id)}
}

func OrderType(orderType string) util.KVP {
	return util.KVP{Key:"order_type", Value:orderType}
}

func Page(page int) util.KVP {
	return util.KVP{Key:"page", Value:strconv.Itoa(page)}
}

// -- Base Specifications --

type HasTypeIdParameter struct {
	specifications.AbstractSpecification
}

func (s *HasTypeIdParameter) IsSatisfiedBy(object interface{}) bool {
	if market, ok := object.(*MarketQuery); ok {
		_, ok := market.QueryParams["type_id"]
		return ok
	}
	return false
}

type HasOrderTypeParameter struct {
	specifications.AbstractSpecification
}

func (s *HasOrderTypeParameter) IsSatisfiedBy(object interface{}) bool {
	if market, ok := object.(*MarketQuery); ok {
		_, ok := market.QueryParams["order_type"]
		return ok
	}
	return false
}

