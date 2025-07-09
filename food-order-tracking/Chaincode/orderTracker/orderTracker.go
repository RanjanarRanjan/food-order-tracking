// File: food-order-tracking/chaincode/orderTracker/orderTracker.go
package orderTracker

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Menu struct {
	DocType  string `json:"docType"`
	MenuID   string `json:"menuID"`
	FoodName string `json:"foodName"`
	Price    int    `json:"price"`
}

type Order struct {
	DocType string `json:"docType"`
	OrderID string `json:"orderID"`
	Status  string `json:"status"`
	Time    string `json:"timestamp"`
}

type PrivateOrderDetails struct {
	OrderID  string `json:"orderID"`
	FoodName string `json:"foodName"`
	Quantity int    `json:"quantity"`
}

// ✅ Org1 only
func (s *SmartContract) CreateMenu(ctx contractapi.TransactionContextInterface, menuID, foodName string, price int) error {
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil || clientMSPID != "Org1MSP" {
		return fmt.Errorf("only Org1 can create menu items")
	}

	exists, err := ctx.GetStub().GetState(menuID)
	if err != nil {
		return err
	}
	if exists != nil {
		return fmt.Errorf("menu %s already exists", menuID)
	}

	menu := Menu{
		DocType:  "menu",
		MenuID:   menuID,
		FoodName: foodName,
		Price:    price,
	}
	menuBytes, _ := json.Marshal(menu)
	return ctx.GetStub().PutState(menuID, menuBytes)
}

// ✅ All orgs
func (s *SmartContract) SearchMenuByFoodName(ctx contractapi.TransactionContextInterface, foodName string) ([]*Menu, error) {
	query := fmt.Sprintf(`{"selector":{"docType":"menu","foodName":"%s"}}`, foodName)
	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var menus []*Menu
	for resultsIterator.HasNext() {
		kv, _ := resultsIterator.Next()
		var menu Menu
		_ = json.Unmarshal(kv.Value, &menu)
		menus = append(menus, &menu)
	}
	return menus, nil
}

// ✅ Org2 only
func (s *SmartContract) PlaceOrder(ctx contractapi.TransactionContextInterface, orderID, timestamp string) error {
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil || clientMSPID != "Org2MSP" {
		return fmt.Errorf("only Org2 can place orders")
	}

	exists, err := ctx.GetStub().GetState(orderID)
	if err != nil {
		return err
	}
	if exists != nil {
		return fmt.Errorf("order %s already exists", orderID)
	}

	order := Order{
		DocType: "order",
		OrderID: orderID,
		Status:  "Placed",
		Time:    timestamp,
	}
	orderBytes, _ := json.Marshal(order)
	return ctx.GetStub().PutState(orderID, orderBytes)
}

// ✅ Org2 only
func (s *SmartContract) AddPrivateOrderDetails(ctx contractapi.TransactionContextInterface) error {
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil || clientMSPID != "Org2MSP" {
		return fmt.Errorf("only Org2 can add private order details")
	}

	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("error getting transient: %v", err)
	}
	orderJSON, ok := transientMap["order"]
	if !ok {
		return fmt.Errorf("order key not found in transient map")
	}
	var privateDetails PrivateOrderDetails
	_ = json.Unmarshal(orderJSON, &privateDetails)
	orderBytes, _ := json.Marshal(privateDetails)
	return ctx.GetStub().PutPrivateData("privateOrderDetails", privateDetails.OrderID, orderBytes)
}

// ✅ Org1, Org2
func (s *SmartContract) GetPrivateOrderDetails(ctx contractapi.TransactionContextInterface, orderID string) (*PrivateOrderDetails, error) {
	orderBytes, err := ctx.GetStub().GetPrivateData("privateOrderDetails", orderID)
	if err != nil {
		return nil, err
	} else if orderBytes == nil {
		return nil, fmt.Errorf("order details not found")
	}
	var order PrivateOrderDetails
	_ = json.Unmarshal(orderBytes, &order)
	return &order, nil
}

// ✅ Org3 only
func (s *SmartContract) UpdateOrderStatus(ctx contractapi.TransactionContextInterface, orderID, newStatus string) error {
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get MSP ID: %v", err)
	}
	if clientMSPID != "Org3MSP" {
		return fmt.Errorf("only Org3 (Delivery Partner) can update order status")
	}

	orderBytes, err := ctx.GetStub().GetState(orderID)
	if err != nil || orderBytes == nil {
		return fmt.Errorf("order not found")
	}
	var order Order
	_ = json.Unmarshal(orderBytes, &order)
	order.Status = newStatus
	updatedBytes, _ := json.Marshal(order)
	return ctx.GetStub().PutState(orderID, updatedBytes)
}

// ✅ All orgs
func (s *SmartContract) GetAllMenus(ctx contractapi.TransactionContextInterface) ([]*Menu, error) {
	query := `{"selector":{"docType":"menu"}}`
	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var menus []*Menu
	for resultsIterator.HasNext() {
		kv, _ := resultsIterator.Next()
		var menu Menu
		_ = json.Unmarshal(kv.Value, &menu)
		menus = append(menus, &menu)
	}
	return menus, nil
}

// ✅ All orgs
func (s *SmartContract) GetMenuByID(ctx contractapi.TransactionContextInterface, menuID string) (*Menu, error) {
	menuBytes, err := ctx.GetStub().GetState(menuID)
	if err != nil {
		return nil, fmt.Errorf("failed to get menu: %v", err)
	}
	if menuBytes == nil {
		return nil, fmt.Errorf("menu %s does not exist", menuID)
	}
	var menu Menu
	err = json.Unmarshal(menuBytes, &menu)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}
	return &menu, nil
}

func (s *SmartContract) GetOrderHistory(ctx contractapi.TransactionContextInterface, orderID string) ([]map[string]interface{}, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get history for %s: %v", orderID, err)
	}
	defer resultsIterator.Close()

	var records []map[string]interface{}
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var record map[string]interface{}
		if response.Value != nil {
			err = json.Unmarshal(response.Value, &record)
			if err != nil {
				return nil, err
			}
		} else {
			record = make(map[string]interface{})
		}
		record["TxId"] = response.TxId
		record["Timestamp"] = response.Timestamp.String()
		record["IsDelete"] = response.IsDelete

		records = append(records, record)
	}
	return records, nil
}
