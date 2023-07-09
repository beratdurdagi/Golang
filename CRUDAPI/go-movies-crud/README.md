# Movie API

Bu, Gorilla Mux yönlendiricisini kullanarak Golang ile uygulanan basit bir film API'sidir.

## Kurulum

1. Depoyu klonlayın:

   ```bash
   git clone https://github.com/karalakrepp/golang.git
2. Proje dizinine gidin: 
    ```bash 
   cd golang

3. Bağımlılıkları yükleyin: 
   ```bash 
  go mod download

## Kullanım

1. Sunucuyu başlatın:
   ```bash  
   go run main.go


Sunucu varsayılan olarak 8000 numaralı portta çalışacaktır. PORT ortam değişkenini ayarlayarak özel bir port belirleyebilirsiniz.

2. API uç noktalarıyla etkileşime geçin:
- GET /movies: Tüm filmleri alın.
- GET /movies/{id}: Belirli bir filmi ID'ye göre alın.
- POST /movies: Yeni bir film oluşturun.
- PUT /movies/{id}: Varolan bir filmi güncelleyin.
- DELETE /movies/{id}: Bir filmi silin.
- Bu uç noktalara HTTP istekleri yapmak için cURL veya Postman gibi araçları kullanabilirsiniz.

