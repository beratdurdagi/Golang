# Get Students

**Bu proje, Golang programlama dili kullanılarak öğrenci bilgilerini saklamak ve yönetmek için basit bir API sunucusu sağlamaktadır. API, öğrenci listesini görüntülemek, öğrenci eklemek, öğrenci bilgilerini güncellemek ve öğrenci hesaplarını pasif duruma getirmek gibi temel işlemleri gerçekleştirmek için kullanılabilir.**

## Kurulum

Projenin çalışması için aşağıdaki adımları izleyin:

1. Bu projeyi bilgisayarınıza kopyalayın veya indirin.

2. `.env` dosyasını projenin kök dizinine ekleyin ve aşağıdaki değişkenleri doldurun:
   * PORT="8080"
   * POSTGRES_URL="postgres://username:password@localhost:port/db_name?sslmode=disable"


3. Gerekli Go paketlerini indirmek için aşağıdaki komutu çalıştırın:

   ```bash
   go mod tidy


## API Endpoint'leri
* GET /api/students: Tüm öğrenci bilgilerini almak için kullanılır.
* GET /api/student/{number}: Belirli bir öğrenci numarasına göre öğrenci bilgilerini almak için kullanılır.
* POST /api/student: Yeni bir öğrenci eklemek için kullanılır. Öğrenci verileri JSON formatında gönderilmelidir.
* PUT /api/student/{number}: Belirli bir öğrenci numarasına göre öğrenci bilgilerini güncellemek için kullanılır. Güncellenen öğrenci  verileri JSON formatında gönderilmelidir.
* DELETE /api/student/{number}: Belirli bir öğrenci numarasına göre öğrenciyi pasif duruma getirmek için kullanılır.

## Kullanılan Teknolojiler
* Go (Golang)
* Gorilla Mux
* PostgreSQL