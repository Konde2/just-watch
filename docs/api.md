cd video-service
go mod tidy
go build .
./video-service


---

### 📄 `docs/api.md` — API

```markdown
# 📡 API

Все запросы проходят через **API Gateway** (`http://localhost:8080`).

## Загрузка видео

```http
POST /upload
Content-Type: multipart/form-data

Form fields:
- video: файл (video/*)
- title: (опционально) название
- owner_id: (опционально) ID владельца

Ответ

---

### 📄 `docs/api.md` — API

```markdown
# 📡 API

Все запросы проходят через **API Gateway** (`http://localhost:8080`).

## Загрузка видео

```http
POST /upload
Content-Type: multipart/form-data

Form fields:
- video: файл (video/*)
- title: (опционально) название
- owner_id: (опционально) ID владельца


Ответ
{ "video_id": "uuid", "message": "Upload successful" }

Стриминг видео
GET /stream/{video_id}

Поддерживает Range-запросы для перемотки.