// File: config/iotda.go
package config

import (
	"fmt"
	"log"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/region" // Pastikan import ini ditambahkan
	iotda "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"
)

func NewIoTDAClient() (*iotda.IoTDAClient, error) {

	credentials, err := basic.
		NewCredentialsBuilder().
		WithAk("HPUAVLKXCF9DLINH1GSI").
		WithSk("DnOrCwhzGBnkvu4nT7xXXsIpV5N3kdBdWjDUAlZ7").
		WithProjectId("019e15e534287864a97ed18e6e5bc19d"). // Project ID AP-Jakarta
		WithDerivedPredicate(auth.GetDefaultDerivedPredicate()).
		SafeBuild()

	if err != nil {
		return nil, err
	}

	// Buat Custom Region yang mengikat nama region "ap-southeast-4" dengan URL Instance Anda
	apiEndpoint := "https://d45409fb27.st1.iotda-app.ap-southeast-4.myhuaweicloud.com"
	// customRegion := region.NewRegion("ap-southeast-4", apiEndpoint)

	hcClient, err := iotda.
		IoTDAClientBuilder().
		WithRegion(region.NewRegion("ap-southeast-4", apiEndpoint)). // Gunakan WithRegion, bukan WithEndpoint
		WithCredential(credentials).
		SafeBuild()

	if err != nil {
		return nil, err
	}

	client := iotda.NewIoTDAClient(hcClient)

	log.Println("✅ IOTDA SDK CONNECTED")

	// Instantiate a request object.
	request := &model.ListDevicesRequest{}
	// Call the API for querying the device list.
	response, err := client.ListDevices(request)
	if err != nil {

		return nil, err
	} else {
		fmt.Println("✅ BERHASIL! Autentikasi API lolos.")
		if response.Page != nil && response.Page.Count != nil {
			fmt.Printf("Berhasil menemukan %d device di IoTDA Anda.\n", *response.Page.Count)
		}
	}

	return client, nil
}
