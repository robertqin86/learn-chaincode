/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors" // standard go error format.
	"fmt" // contains Println for debugging/logging.
	"github.com/hyperledger/fabric/core/chaincode/shim" // contains the definition for the chaincode interface and the chaincode stub, which you will need to interact with the ledger.
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init is called when you first deploy your chaincode. As the name implies, this function 
// should be used to do any initialization your chaincode needs. 
// In our example, we use Init to configure the initial state of a single key/value pair on the ledger
// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

    // Stores the first element in the args argument to the key "hello_world".
    // This is done by using the stub function stub.PutState. 
    // The function interprets the first argument sent in the deployment request as the value to be stored under the key 'hello_world' 
    // in the ledger. All will be explained after we finish implementing the chaincode interface. 
	err := stub.PutState("hello_world", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke is called when you want to call chaincode functions to do real work. 
// Invocations will be captured as a transactions, which get grouped into blocks on the chain.
// Invoke is our entry point to invoke a chaincode function
// When you need to update the ledger, you will do so by invoking your chaincode. 
// The structure of Invoke is simple. It receives a function and an array of arguments. 
// Based on what function was passed in through the function parameter in the invoke request, 
// Invoke will either call a helper function or return an error.

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "write" {
		// Write uses two arguments, allowing you to pass in both the key and the value for the call to PutState. 
		// Basically, this function allows you to store any key/value pair you want into the blockchain ledger.
		return t.write(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is called whenever you query your chaincode's state. Queries do not result in blocks being added to the chain, 
// and you cannot use functions like PutState inside of Query or any helper functions it calls. 
// You will use Query to read the value of your chaincode state's key/value pairs.
// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" {											
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query: " + function)
}

func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string)([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key   = args[0]
	value = args[1]
	err   = stub.PutState(key, []byte(value))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	//  While PutState allows you to set a key/value pair, 
	//  GetState lets you read the value for a previously written key.
	valAsbytes, err := stub.GetState(key)

	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}
