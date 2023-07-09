# Movie API

Bu, Gorilla Mux yönlendiricisini kullanarak Golang ile uygulanan basit bir film API'sidir.

## Kurulum

1. Depoyu klonlayın:

   ```bash
   git clone https://github.com/karalakrepp/golang.git
2. Proje dizinine gidin: cd golang

3. Bağımlılıkları yükleyin: go mod download

## Kullanım

1. Sunucuyu başlatın:go run main.go
   Sunucu varsayılan olarak 8000 numaralı portta çalışacaktır. PORT ortam değişkenini ayarlayarak özel bir port belirleyebilirsiniz.

2. API uç noktalarıyla etkileşime geçin:
GET /movies: Tüm filmleri alın.
GET /movies/{id}: Belirli bir filmi ID'ye göre alın.
POST /movies: Yeni bir film oluşturun.
PUT /movies/{id}: Varolan bir filmi güncelleyin.
DELETE /movies/{id}: Bir filmi silin.
Bu uç noktalara HTTP istekleri yapmak için cURL veya Postman gibi araçları kullanabilirsiniz.

## API Belgesi
# Film Nesnesi

Bir film nesnesi aşağıdaki özelliklere sahiptir:

. id (string): Filmin ID'si.
. isbn (string): Filmin ISBN numarası.
. title (string): Filmin başlığı.
. director (object): Filmin yönetmeni.

# Uç Noktalar
. GET /movies
. Tüm filmleri alın.

# Yanıt

. Durum Kodu: 200 (OK)
. Yanıt Gövdesi: Film nesnelerinin bir dizisi
. Örnek Yanıt Gövdesi:[
  {
    "id": "1",
    "isbn": "1234567890",
    "title": "Film 1",
    "director": {
      "firstname": "John",
      "lastname": "Doe"
    }
  },
  {
    "id": "2",
    "isbn": "0987654321",
    "title": "Film 2",
    "director": {
      "firstname": "Jane",
      "lastname": "Smith"
    }
  }
]

. GET /movies/{id}: Belirli bir filmi ID'ye göre alın.

# Yol Parametreleri

. id (string, gerekli): Film ID'si.
. Durum Kodu: 200 (OK)
. Yanıt Gövdesi: Film nesnesi
. Örnek Yanıt Gövdesi:{
  "id": "1",
  "isbn": "1234567890",
  "title": "Film 1",
  "director": {
    "firstname": "John",
    "lastname": "Doe"
  }
}

. POST /movies:Yeni bir film oluşturun.

# Request Body

. Movie Object
 Example Request Body:{
  "isbn": "9876543210",
  "title": "New Movie",
  "director": {
    "firstname": "Alice",
    "lastname": "Johnson"
  }
}

# Response
. Status Code: 200 (OK)
. Response Body: Created movie object
Example Response Body:{
  "id": "3",
  "isbn": "9876543210",
  "title": "New Movie",
  "director": {
    "firstname": "Alice",
    "lastname": "Johnson"
  }
}
...