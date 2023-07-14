# Go-React Todo Uygulaması

Bu proje, Go ve React kullanarak basit bir Todo uygulamasının bir parçası olan bir ara katman (middleware) sunar. Go dilinde yazılan ara katman, MongoDB veritabanına bağlanır ve HTTP isteklerini işleyerek görevlerin yönetilmesini sağlar. React tarafı, kullanıcı arayüzünü sunmak için kullanılabilir.

## Kurulum

1. Bu projeyi yerel makinenize klonlayın:

   ```shell
   git clone https://github.com/kullaniciadi/go-react-todo.git

2. MongoDB veritabanı oluşturun ve bir veritabanı adı ile bir koleksiyon adı belirleyin.

3. .env dosyasını oluşturun ve aşağıdaki değişkenleri belirleyin:
   MONGODB_URL=<MongoDB Veritabanı URL'si>
   DB_NAME=<Veritabanı Adı>
   DB_COLLECTION_NAME=<Koleksiyon Adı>

4. Bağımlılıkları yüklemek için aşağıdaki komutu çalıştırın:
    ```shell
   go run main.go

5. Go ara katmanını başlatmak için aşağıdaki komutu çalıştırın:
      ```shell
      go run main.go

6. React uygulamasını başlatmak için frontend dizinine gidin ve aşağıdaki komutu çalıştırın:
   ```shell
   npm install
   npm start


## Kullanım
#### Uygulama, Todo listesini yönetmek için aşağıdaki HTTP isteklerini sağlar:

`GET /tasks`: Tüm görevleri alır.
`POST /tasks`: Yeni bir görev oluşturur.
`PUT /tasks/:id/complete`: Bir görevi tamamlanmış olarak işaretler.
`PUT /tasks/:id/undo`: Bir görevi tamamlanmamış olarak işaretler.
`DELETE /tasks/:id`: Bir görevi siler.
`DELETE /tasks`: Tüm görevleri siler.
