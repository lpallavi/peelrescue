package buyer

import (
	"errors"
	"fmt"
	"projectGoLive/application/apiclient"
	"sync"
)

// Linked list node containing an item of type person, and a pointer to next node
type CartNode struct {
	Item apiclient.ItemsDetails
	Next *CartNode
}

// Linked list structure containing the head node, size of linked list, and mutex
type CartLinkedList struct {
	Head *CartNode
	Size int
	mu   sync.Mutex
}

// This is a method for linked list struct
// It is used to get an item at the specified index
// It returns item struct and any errors
func (c *CartLinkedList) Get(index int) (apiclient.ItemsDetails, error) {
	emptyItem := apiclient.ItemsDetails{}
	if c.Head == nil {
		return emptyItem, errors.New("Empty Linked list!")
	}
	if index > 0 && index <= c.Size {
		currentNode := c.Head
		for i := 1; i <= index-1; i++ {
			currentNode = currentNode.Next
		}
		item := currentNode.Item
		return item, nil

	}
	return emptyItem, errors.New("Invalid Index")
}

// This is a method for linked list struct
// It is used to add an item to the linked list
// It takes in the item to be added
// The node for the item is added at the end of the linked list
// Mutex lock is enabled during addition of node
func (c *CartLinkedList) AddNode(name apiclient.ItemsDetails) error {
	c.mu.Lock()
	{
		newNode := &CartNode{
			Item: name,
			Next: nil,
		}
		if c.Head == nil {
			c.Head = newNode
		} else {
			currentNode := c.Head
			for currentNode.Next != nil {
				currentNode = currentNode.Next
			}
			currentNode.Next = newNode
		}
		c.Size++

	}
	c.mu.Unlock()
	return nil
}

// This is a method for linked list struct
// It is used to add an item to the linked list at a certain index
// It traverses the linked list, and adds name of type person at location index
func (c *CartLinkedList) AddAtPos(index int, name apiclient.ItemsDetails) error {
	newNode := &CartNode{
		Item: name,
		Next: nil,
	}

	if index > 0 && index <= c.Size+1 {
		if index == 1 {
			newNode.Next = c.Head
			c.Head = newNode

		} else {

			currentNode := c.Head
			var prevNode *CartNode
			for i := 1; i <= index-1; i++ {
				prevNode = currentNode
				currentNode = currentNode.Next
			}
			newNode.Next = currentNode
			prevNode.Next = newNode

		}
		c.Size++
		return nil
	} else {
		return errors.New("Invalid Index")
	}
}

// This is a method for linked list struct
// It is used to remove an item from the linked list at a certain index
// It traverses the linked list, and removes the item at location index
// It returns the item that is removed, and any errors if present
func (c *CartLinkedList) Remove(index int) (apiclient.ItemsDetails, error) {
	var item apiclient.ItemsDetails
	emptyItem := apiclient.ItemsDetails{}

	if c.Head == nil {
		return emptyItem, errors.New("Empty Linked list!")
	}
	if index > 0 && index <= c.Size {
		if index == 1 {
			item = c.Head.Item
			c.Head = c.Head.Next
		} else {
			var currentNode *CartNode = c.Head
			var prevNode *CartNode
			for i := 1; i <= index-1; i++ {
				prevNode = currentNode
				currentNode = currentNode.Next

			}
			item = currentNode.Item
			prevNode.Next = currentNode.Next
		}
	}
	c.Size--
	return item, nil
}

// This is a method for linked list struct
// It is used to list all the items in the linked list
// It traverses the linked list, and returns all the items
// This is used by the buyer template to display list of all items in the cart
func (c *CartLinkedList) GetAllItems() (msg []string, allitems []apiclient.ItemsDetails) {
	var message []string
	var items []apiclient.ItemsDetails

	count := 1
	currentNode := c.Head
	if currentNode == nil {
		message = append(message, fmt.Sprintf("No users found."))
		return message, items
	}
	message = append(message, fmt.Sprintf("\nListing all Items in your shopping cart:"))
	items = append(items, currentNode.Item)

	count++
	for currentNode.Next != nil {
		currentNode = currentNode.Next
		items = append(items, currentNode.Item)
		count++
	}
	return message, items
}

// This is a method for linked list struct
// It is used to search for a specific item in the linked list of any seller
// It traverses the linked list, and returns all the items of a given name
// NO: It also returns the index at which the item is found, and any errors if present
func (c *CartLinkedList) SearchItemName(itemname string) ([]apiclient.ItemsDetails, error) {
	itemslist := []apiclient.ItemsDetails{}
	//index := 1
	if c.Head == nil {
		return itemslist, errors.New("Empty Linked list!")
	}
	currentNode := c.Head
	for i := 1; i <= c.Size; i++ {
		if currentNode.Item.Item != itemname {
			currentNode = currentNode.Next
		} else {
			item := currentNode.Item
			itemslist = append(itemslist, item)
		}
	}
	if itemslist == nil {
		return itemslist, errors.New("Sellername not found in list")
	} else {
		// one or more items found
		return itemslist, nil
	}
}

// This is a method for linked list struct
// It is used to search for items in the linked list of a particular seller
// It traverses the linked list, and returns all the items for the seller name
// NO :It also returns the index at which the item is found, and any errors if present
func (c *CartLinkedList) SearchSellerName(sellername string) ([]apiclient.ItemsDetails, error) {
	itemslist := []apiclient.ItemsDetails{}
	//index := 1
	if c.Head == nil {
		return itemslist, errors.New("Empty Linked list!")
	}
	currentNode := c.Head
	for i := 1; i <= c.Size; i++ {
		if currentNode.Item.Username != sellername {
			currentNode = currentNode.Next
		} else {
			item := currentNode.Item
			itemslist = append(itemslist, item)
		}
	}
	if itemslist == nil {
		return itemslist, errors.New("Sellername not found in list")
	} else {
		// one or more items found
		return itemslist, nil
	}
}

// This is a method for linked list struct
// It is used to search for a specific item in the linked list of any seller
// It traverses the linked list, and returns all the items of a given name
// NO: It also returns the index at which the item is found, and any errors if present
func (c *CartLinkedList) SearchItemandSellerName(itemname, sellername string) (apiclient.ItemsDetails, int, error) {
	item := apiclient.ItemsDetails{}
	index := 1
	if c.Head == nil {
		return item, -1, errors.New("Empty Linked list!")
	}
	currentNode := c.Head
	for i := 1; i <= c.Size; i++ {
		if currentNode.Item.Item != itemname || currentNode.Item.Username != sellername {
			currentNode = currentNode.Next
			index++
		} else {
			item := currentNode.Item
			return item, index, nil
		}
	}
	return item, -1, errors.New("Item not found in list")
}

// This is a method for linked list struct
// It is used to write Item data at a specified index in the linked list
// It takes in an index of type int, and thisItem of type ItemDetails
// The item data at the index specified is overwritten by contents of thisItem
// It returns any error if present
func (c *CartLinkedList) WriteAtIndex(index int, thisItem apiclient.ItemsDetails) error {

	if c.Head == nil {
		return errors.New("Empty Linked list!")
	}
	if index > 0 && index <= c.Size {
		currentNode := c.Head
		for i := 1; i <= index-1; i++ {
			currentNode = currentNode.Next
		}
		currentNode.Item = thisItem
		return nil

	}
	return errors.New("Invalid Index")
}

// This is a method for linked list struct
// It is a wrapper used to write item data at a specified index in the linked list
// It takes in an index and an item
// The item at index specified is overwritten
// It calls writeAtIndex function to write the item data at index
func (c *CartLinkedList) WriteItemData(thisItem apiclient.ItemsDetails, itemIndex int) error {
	err := c.WriteAtIndex(itemIndex, thisItem)
	if err != nil {
		return err
	}
	return nil
}
