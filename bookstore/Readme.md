# Bookstore Projesi

Bu proje, kitap bilgilerini yönetmek için basit bir RESTful API sağlayan bir Go (Golang) projesidir. Proje, Gorilla Mux yönlendirme kütüphanesi ve GORM ORM kullanılarak oluşturulmuştur.

## Kurulum ve Çalıştırma

1. Proje kaynak kodlarını indirin veya kopyalayın.
2. PostgreSQL veritabanınızın kurulu olduğundan ve çalıştığından emin olun.
3. `config/db.go` dosyasında yer alan veritabanı bağlantı bilgilerini düzenleyerek kendi veritabanı bilgilerinizle güncelleyin.
4. Terminali açın ve proje dizinine gidin.
5. Proje bağımlılıklarını yüklemek için aşağıdaki komutu çalıştırın:

   ```bash
   go mod download
* Projeyi derlemek ve çalıştırmak için aşağıdaki komutu çalıştırın:
   ```bash
   go run main.go

* Bu projeyi yapmamda videosu ile bana yardımcı olan Akhil Sharma ya teşekkür ederim


*Proje başarıyla çalıştığında, tarayıcınızdan veya API istemci aracından şu URL'leri ziyaret edebilirsiniz:* 
* http://localhost:8080/book/ (GET) - Tüm kitapları listeler.
* http://localhost:8080/book/{bookId} (GET) - Belirli bir kitabın detaylarını gösterir.
* http://localhost:8080/book/ (POST) - Yeni bir kitap ekler.
* http://localhost:8080/book/{bookId} (PUT) - Varolan bir kitabı günceller.
* http://localhost:8080/book/{bookId} (DELETE) - Varolan bir kitabı siler.


## Katkıda Bulunma
Eğer bu projeye katkıda bulunmak isterseniz, lütfen GitHub deposunu forklayın ve pull request gönderin. Her türlü katkı ve geri bildirim değerlidir.