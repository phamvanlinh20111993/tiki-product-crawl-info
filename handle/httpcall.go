package handle

import "fmt"

func getApi() {

	client := resty.New()
	defer client.Close()

	res, err := client.R().
		EnableTrace().
		Get("https://httpbin.org/get")
	fmt.Println(err, res)
	fmt.Println(res.Request.TraceInfo())

}
