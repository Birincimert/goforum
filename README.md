# GoForum ✨

Topluluk odaklı, hafif ve hızlı bir forum uygulaması. Onaylı yazılar, kategoriler, detay sayfaları, yorumlar (nested-div şeklinde yanıtlarla), beğeniler, kaydetme (bookmark) ve kullanıcı profili gibi modern özellikleri bir araya getirir. 💬🔖

## Özellikler 🚀
- 🏠 Anasayfa: Onaylı yazıları listeler, kategoriye göre filtreleme.
- 🗂️ Kategoriler: `/kategori/:slug` ile kategori bazlı listeleme ve sayım.
- 📝 Detay sayfası: Zengin içerik, toplam yorum/yanıt sayısı, yazar bilgisi.
- 💬 Yorumlar: İç içe yanıtlar, beğeni ve beğeni sayacı.
- 🔖 Kaydetme (Bookmark): Yazıyı kaydet/kaldır, profil sayfasında listele.
- 👤 Profil: Kullanıcı yazıları, yorumları, beğendikleri ve kaydettikleri.
- ➕ İçerik yönetimi: Yeni yazı ekleme, düzenleme ve silme, yazılara kapak görseli ekleme (opsiyonel).
- 🛡️ Oturum yönetimi: `gorilla/sessions` ile güvenli cookie store.
- 🧭 Router: `httprouter` ile hızlı ve temiz yönlendirme.

## Mimari Özeti 🧩
- Dil ve çekirdek: Go (net/http, html/template)
- Router: `github.com/julienschmidt/httprouter`
- Oturum: `github.com/gorilla/sessions`
- Ortam: `github.com/joho/godotenv` (.env yükleme)
- UI: Font Awesome 5 (JS), Bootstrap sınıfları
- Kaydetme/Beğeni: AJAX ile REST uç noktaları
- Veritabanı: `github.com/denisenkom/go-mssqldb` (SQL Server)
- ORM: Basit SQL sorguları, `database/sql`
- Sha256: `crypto/sha256` (parola şifrelemek için)
> Not: Proje varsayılan olarak SQL Server DSN kullanır (Windows geliştirici ortamı için uygundur). `site/models/Database.go` içindeki DSN’i kendi ortamınıza göre güncelleyebilirsiniz.

## Dizin Yapısı 🗂️
- `admin/` — Admin paneli (kontroller, modeller, görünümler, statikler)
- `site/` — Site tarafı (kontroller, modeller, görünümler, statikler)
- `config/Routes.go` — Tüm HTTP rotaları
- `main.go` — Sunucu başlangıç noktası, migrasyonlar ve servis
- `uploads/` — Yüklenen görseller (kapak vb.)

## Hızlı Başlangıç ⚡
1) Depoyu klonla
```bash
  git clone https://github.com/USERNAME/goforum.git
  cd goforum
```

2) Go bağımlılıklarını indir
```bash
  go mod download
```

3) .env oluştur ve oturum anahtarını ayarla
```bash
# .env dosyası
SESSION_KEY=buraya-uzun-ve-rastgele-bir-gizli-anahtar-yazin
```

4) (İsteğe bağlı) Ön yüz paketleri için
```bash
  npm install
```

5) Çalıştır
```bash
  go run ./...
  # veya
  go run main.go
```
Default Sunucu: http://localhost:8080

## Veritabanı ⚙️
- Varsayılan DSN (SQL Server, Windows): `site/models/Database.go`
```
server=localhost,52175;database=goforum;trusted_connection=yes
```
- Kendi sunucunuza göre `server`, `database`, `trusted_connection`/`user id`/`password` alanlarını düzenleyin.
- Uygulama başlangıcında gerekli tablolar için migrasyonlar otomatik tetiklenir (bkz. `main.go`).

## Önemli Uç Noktalar 🔗
- Genel
  - `GET /` — Anasayfa (liste)
  - `GET /post/:slug` — Yazı detay (ID fallback destekli)
  - `GET /yazilar/:slug` — Yazı detay (alternatif yol)
  - `GET /kategori/:slug` — Kategori listesi
  - `GET /about`, `GET /contact`
- Yorumlar
  - `POST /comment/add` — Yorum ekle
  - `POST /comment/upvote/:id` — Yorum beğen (toggle)
  - `GET /comment/likes/:id` — Yorum beğeni sayısı
  - `GET /comment/liked/:id` — Beğeni durumu
- Kaydetme (Bookmark)
  - `POST /post/save/:id` — Yazıyı kaydet/kaldır (toggle)
  - `GET /save/status/:id` — Kaydetme durumu
- Profil ve içerik
  - `GET /profile` — Profil (giriş gerekli)
  - `GET /profile/new-post` — Yeni yazı formu
  - `POST /profile/new-post` — Yeni yazı ekle
  - `GET /profile/post/edit/:id` — Yazı düzenleme formu
  - `POST /profile/post/edit/:id` — Yazı düzenle
  - `GET /profile/post/delete/:id` — Yazı sil

```
Admin Route'ları da `admin/` altında benzer şekilde yapılandırılmıştır.
```

## UI İpuçları 🎨
- Bookmark ikonu için iki ikon yaklaşımı kullanılır:
  - Boş: `far fa-bookmark`, Dolu: `fas fa-bookmark`
  - JS yalnızca butona `.saved` sınıfını ekler/çıkarır, CSS görünürlüğü yönetir.
- Font Awesome 5 (JS) `<i>` etiketlerini SVG’ye çevirdiği için boyutlandırmayı butondan yapın. Örneğin:
```css
.bookmark-btn { font-size: 3rem; }
.bookmark-btn.saved { color: #f4d03f; }
```

## Geliştirme 🛠️
- Windows ortamında test edilmiştir.
- Statik dosyalar: `/assets/*` ve `/uploads/*` rota servisleri.
- Şablon render’ı ResponseWriter’a yazılmadan önce buffer’da yapılır; böylece `superfluous WriteHeader` uyarısı engellenir.

## Sorun Giderme 🧯
- `wsasend: Kurulan bir bağlantı … iptal edildi` (Windows): Genelde kullanıcı sayfayı kapattığında/yenilediğinde oluşur; kritik değildir.
- `http: superfluous response.WriteHeader`: Birden fazla header yazımı; şablonlar buffer üzerinden yazıldığı için çözülmüştür.
- DSN bağlantı sorunları: `site/models/Database.go` içindeki DSN’i kontrol edin, sunucu/port erişimini doğrulayın.

## Dağıtım ☁️
- Tek binary derleyin:
```bash
  go build -o goforum
```
- Reverse proxy (Nginx/Apache) arkasında çalıştırın, `uploads` klasörü izinlerini doğru ayarlayın.
- Ortam değişkenlerini (.env) prod için güvenle yönetin.

## Katkı 🤝
- Issue açın, küçük PR’lar memnuniyetle kabul edilir.
- PR öncesi `go fmt`, `go vet` ve temel manuel kontrolleri çalıştırın.

## Lisans 📄
- Ayrıntılar için `LICENSE` dosyasına bakın.

## Dipnot 📝

```
- Bu proje eğitim amaçlıdır ve gerçek dünya uygulamalarında ek güvenlik/optimizasyon gerektirebilir.

- Proje henüz geliştirme aşamasındadır, güncellendikçe en son sürümü GitHub’da kontrol edin.
```