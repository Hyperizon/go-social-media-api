# Go Social Media API Project

Bu proje, Go programlama dili kullanılarak oluşturulan bir API örneğidir. Bu projede, kullanıcı yönetimi, gönderi oluşturma, beğenme, yorum yapma gibi temel işlevlerin nasıl uygulandığına dair bir örnek bulacaksınız.

## Teknolojiler

- Go Programlama Dili
- Fiber Framework (HTTP Sunucu ve Router)
- GORM (Veritabanı İşlemleri)
- JWT (Kimlik Doğrulama ve Yetkilendirme)
- Validator (Validasyon İşlemleri)

## Proje Yapısı

- `main.go`: Uygulama giriş noktası ve HTTP sunucu başlatma
- `middlewares/`: Orta katman fonksiyonları ve doğrulama middleware'leri
- `controllers/`: API işlevlerini gerçekleştiren kontrolörler
- `models/`: Veritabanı modelleri
- `routes/`: API rotalarını yöneten dosyalar
- `utils/`: Yardımcı işlevler ve araçlar

## Bağımlılıklar

Uygulama aşağıdaki bağımlılıkları kullanmaktadır:

- github.com/go-playground/validator/v10 v10.15.0
-	github.com/gofiber/fiber/v2 v2.48.0
-	github.com/golang-jwt/jwt v3.2.2+incompatible
-	github.com/joho/godotenv v1.5.1
-	gorm.io/driver/postgres v1.5.2
-	gorm.io/gorm v1.25.2

## Kullanım

API dökümanına adresten ulaşabilirsiniz
https://documenter.getpostman.com/view/11298546/2s9XxzvZ5K
