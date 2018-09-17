# caire-openfaas

This is an [OpenFaaS](https://github.com/openfaas/faas) function for [Caire](https://github.com/esimov/caire) content aware image resize library. This function facilitates to run the library without the need to have it installed locally. 

### Usage
To run the function locally you have to make sure OpenFaaS is up and running. Follow the official documentation for more details. https://docs.openfaas.com/

Clone the repository:
```bash
$ git clone https://github.com/esimov/caire-openfaas
```

#### Build
```bash 
$ faas-cli build -f stack.yml --gateway=http://<GATEWAY-IP>
```

#### Deploy
```bash 
$ faas-cli deploy -f stack.yml --gateway=http://<GATEWAY-IP>
```

![sample-screen](https://user-images.githubusercontent.com/883386/45596764-ec5ec900-b9c9-11e8-8b01-8f84f78b327e.png)


Once the function has been deployed you can access the UI on the url defined in `--gateway` parameter.

You have to provide each option parameter as a `JSON` string, defined in the UI `Request Body` section. The `json` should have the structure of the following form:

```
{
	"input":"https://user-images.githubusercontent.com/883386/37569642-0c5f49e8-2aee-11e8-8ac1-d096c0387ca0.jpg", 
	"width":20,
	"height":0,
	"perc":"true",
	"scale":"false",
	"face":"true"
	"classifier":"./data/facefinder"
}
```
For more details about the supported options check the project page: https://github.com/esimov/caire. 

Below are the supported commands:

| Command | Variable Type | Description |
| --- | --- | --- |
| `input` | string | Input file |
| `width` | int | New width |
| `height` | int | New height |
| `perc` | bool | Reduce image by percentage |
| `square` | bool | Reduce image to square dimensions |
| `scale` | bool | Proportional scaling |
| `debug` | bool | Use debugger |
| `blur` | int | Blur radius |
| `sobel` | int | Sobel filter threshold |
| `face` | false | Use face detection |

***Note:** all the boolean type option should be defined as string.*

### Results

| Original | Shrunk |
| --- | --- |
| ![broadway_tower_edit](https://user-images.githubusercontent.com/883386/35498083-83d6015e-04d5-11e8-936a-883e17b76f9d.jpg) | ![broadway_tower_edit](https://user-images.githubusercontent.com/883386/35498110-a4a03328-04d5-11e8-9bf1-f526ef033d6a.jpg) |
| ![waterfall](https://user-images.githubusercontent.com/883386/35498250-2f31e202-04d6-11e8-8840-a78f40fc1a0c.png) | ![waterfall](https://user-images.githubusercontent.com/883386/35498209-0411b16a-04d6-11e8-9ce2-ec4bce34828a.jpg) |
| ![dubai](https://user-images.githubusercontent.com/883386/35498466-1375b88a-04d7-11e8-8f8e-9d202da6a6b3.jpg) | ![dubai](https://user-images.githubusercontent.com/883386/35498499-3c32fc38-04d7-11e8-9f0d-07f63a8bd420.jpg) |
| ![boat](https://user-images.githubusercontent.com/883386/35498465-1317a678-04d7-11e8-9185-ec92ea57f7c6.jpg) | ![boat](https://user-images.githubusercontent.com/883386/35498498-3c0f182c-04d7-11e8-9af8-695bc071e0f1.jpg) |

| Original | Extended |
| --- | --- |
| ![gasadalur](https://user-images.githubusercontent.com/883386/35498662-e11853c4-04d7-11e8-98d7-fcdb27207362.jpg) | ![gasadalur](https://user-images.githubusercontent.com/883386/35498559-87eb6426-04d7-11e8-825c-2dd2abdfc112.jpg) |
| ![dubai](https://user-images.githubusercontent.com/883386/35498466-1375b88a-04d7-11e8-8f8e-9d202da6a6b3.jpg) | ![dubai](https://user-images.githubusercontent.com/883386/35498827-8cee502c-04d8-11e8-8449-05805f196d60.jpg) |

## License

This project is under the MIT License. See the LICENSE file for the full license text.
