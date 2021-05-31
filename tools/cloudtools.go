/*
===========================================================================
MIT License

Copyright (c) 2021 Manish Meganathan, Mariyam A.Ghani

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
===========================================================================
FyrMesh gopkg tools
===========================================================================
*/
package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

// A constructore function that generates and returns a Firestore.Client object
// that is configures with the Service Account Credentials from cloud config.
func NewFirestoreClient() (*firestore.Client, error) {
	ctx := context.Background()

	// Read the 'FYRMESHCONFIG' env var.
	filedir := os.Getenv("FYRMESHCONFIG")
	// Construct the path to the config file
	filelocation := filepath.Join(filedir, "cloudconfig.json")
	if filedir == "" {
		return nil, fmt.Errorf("environment variable 'FYRMESHCONFIG' has not been set")
	}

	// Open the config file
	configfile, err := os.Open(filelocation)
	if err != nil {
		return nil, err
	}

	// Defer the closing of the file
	defer configfile.Close()

	// Define a cloudConfig struct to read the cloudconfig.json file.
	type cloudConfig struct {
		ProjectID string `json:"project_id"`
	}

	// Read the cloudconfig file into a CloudConfig
	var cloudconfig cloudConfig
	byteValue, _ := ioutil.ReadAll(configfile)
	json.Unmarshal([]byte(byteValue), &cloudconfig)

	// Set up the Service Account Credentials
	serviceaccount := option.WithCredentialsFile(filelocation)
	// Generate the Firestore Client with the Project ID from the cloudconfig and the Service Account Credentials
	client, err := firestore.NewClient(ctx, cloudconfig.ProjectID, serviceaccount)
	if err != nil {
		return nil, err
	}

	// Return the client
	return client, nil
}

// A struct that represents the Cloud Interface of the Mesh Orchestrator
type CloudInterface struct {
	// A Firestore Client object
	FirestoreClient firestore.Client

	// A Document Reference to the MeshDocument
	MeshDoc firestore.DocumentRef

	// A Collection Reference to collection of PingDocuments
	PingCollection firestore.CollectionRef
}

// A constructor function that generates a CloudInterface
// object from a given mesh ID, which is taken from the
// DeviceID field of the Config struct.
func NewCloudInterface(meshid string) (*CloudInterface, error) {
	// Create an empty CloudInterface
	cloudinterface := CloudInterface{}

	// Generate a new Firestore client
	client, err := NewFirestoreClient()
	if err != nil {
		return nil, fmt.Errorf("could not construct firestore client - %v", err)
	}

	// Assign the Firestore client
	cloudinterface.FirestoreClient = *client
	// Assign the Document Reference to the mesh document
	cloudinterface.MeshDoc = *client.Collection("meshes").Doc(meshid)
	// Assign the Collection Reference to the ping collection
	cloudinterface.PingCollection = *cloudinterface.MeshDoc.Collection("pings")

	// Return the cloud interface
	return &cloudinterface, nil
}

// A struct that represents the Firestore Document
// that contains the values that make up a MeshPing
type PingDocument struct {
	PingID          string                        `firestore:"pingid"`
	Pingtime        string                        `firestore:"pingtime"`
	Nodelist        []int64                       `firestore:"nodelist"`
	Sensordata      map[string]map[string]float64 `firestore:"sensordata"`
	Probabilitydata map[string]float64            `firestore:"probability"`
	AvgProbability  float64                       `firestore:"avgprobability"`
}

// A constructor function that generates and returns a PingDocument object from a given MeshPing.
func NewPingDocument(meshping *MeshPing) *PingDocument {
	// Create an empty MeshPing
	pingdoc := PingDocument{}

	// Assign the PingID, Pingtime and Nodelist from the MeshPing
	pingdoc.PingID = meshping.PingID
	pingdoc.Pingtime = meshping.Pingtime
	pingdoc.Nodelist = meshping.Nodelist
	// Generate and assign the Sensordata, Probabilitydata and AvgProbability
	pingdoc.Sensordata = meshping.GenerateSensordatamap()
	pingdoc.Probabilitydata = meshping.GenerateProbabilitydatamap()
	pingdoc.AvgProbability = meshping.GenerateAvgProbability()

	// Return the PingDocument
	return &pingdoc
}

// A method of PingDocument that writes the Document to the Firestore database
// A document is created with the Pingtime as the ID and the values of the PingDocument are written.
func (pingdoc *PingDocument) Push(cloudinterface *CloudInterface) error {
	_, err := cloudinterface.PingCollection.Doc(pingdoc.Pingtime).Create(context.Background(), pingdoc)
	return err
}

// A struct that represents the Credentials required to login to the mesh dashboard.
type Credentials struct {
	Username string `firestore:"username"`
	Password string `firestore:"password"`
}

// A struct that represents the hardware coded configurations of the sensor mesh.
type MeshConfiguration struct {
	Meshssid string `firestore:"MESHSSID"`
	Meshpswd string `firestore:"MESHPSWD"`
	Meshport int    `firestore:"MESHPORT"`
}

// A struct that represents the Firestore Document
// that contains the values that make up a MeshOrchestrator
type MeshDocument struct {
	Credentials       Credentials           `firestore:"credentials"`
	MeshConfiguration MeshConfiguration     `firestore:"meshconfiguration"`
	ControllerID      string                `firestore:"controllerID"`
	ControlnodeID     string                `firestore:"controlnodeID"`
	ControlnodeConfig ControlNode           `firestore:"controlnode"`
	Nodelist          map[string]SensorNode `firestore:"nodes"`
	NodeIDlist        []int64               `firestore:"nodeids"`
}

// A constructor function that generates and returns a
// MeshDocument object from a given MeshOrchestrator
func NewMeshDocument(meshorchestrator *MeshOrchestrator) *MeshDocument {
	// Create an empty MeshDocument
	meshdoc := MeshDocument{}

	// Create and assign the credentials. Fixed password is a temporary implementation
	meshdoc.Credentials = Credentials{
		Username: meshorchestrator.ControllerID,
		Password: "123456",
	}

	// Create and assign the mesh configuration.
	meshdoc.MeshConfiguration = MeshConfiguration{
		Meshssid: meshorchestrator.Controlnode.MeshSSID,
		Meshpswd: meshorchestrator.Controlnode.MeshPSWD,
		Meshport: meshorchestrator.Controlnode.MeshPORT,
	}

	// Assign the Controller ID
	meshdoc.ControllerID = meshorchestrator.ControllerID
	// Assign the ControlnodeID after converting to a string
	meshdoc.ControlnodeID = strconv.FormatInt(meshorchestrator.Controlnode.NodeID, 10)
	// Assign the ControlnodeConfig to the ControlNode object
	meshdoc.ControlnodeConfig = meshorchestrator.Controlnode
	// Assign the NodeIDlist
	meshdoc.NodeIDlist = meshorchestrator.NodeIDlist

	// Generate and assign a Nodelist from the MeshOrchestrator. Needs to have string keys for Firestore.
	meshdoc.Nodelist = make(map[string]SensorNode)
	for nodeid, node := range meshorchestrator.Nodelist {
		strnodeid := strconv.FormatInt(nodeid, 10)
		meshdoc.Nodelist[strnodeid] = node
	}

	// Return the MeshDocument
	return &meshdoc
}

// A method of MeshDocument that writes the Document to the Firestore database
// The MeshDoc of the cloud interface is set with the values from the MeshDocument.
func (meshdoc *MeshDocument) Push(cloudinterface *CloudInterface) error {
	_, err := cloudinterface.MeshDoc.Set(context.Background(), meshdoc)
	return err
}
