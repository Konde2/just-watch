<template>
  <div id="app">
    <!-- Navbar -->
    <header class="nav">
      <div class="nav-left">
        <div class="logo">JUST WATCH</div>
        <nav class="nav-links">
          <a href="#">Главная</a>
          <a href="#">Сериалы</a>
          <a href="#">Фильмы</a>
          <a href="#">Новое и популярное</a>
          <a href="#">Список</a>
        </nav>
      </div>
      <div class="nav-right">
        <button class="upload-btn" @click="showUpload = true">Загрузить</button>
        <div class="avatar">JW</div>
      </div>
    </header>

    <!-- Hero -->
    <section class="hero" :style="heroStyle">
      <div class="hero-overlay"></div>
      <div class="hero-content">
        <h1 class="hero-title">{{ featured.title }}</h1>
        <p class="hero-desc">{{ featured.description }}</p>
        <div class="hero-actions">
          <button class="btn btn-primary" @click="playFeatured">▶ Смотреть</button>
          <button class="btn btn-secondary" @click="showInfo = true">ℹ Подробнее</button>
        </div>
      </div>
      <div class="hero-gradient"></div>
    </section>

    <!-- Rows -->
    <main class="rows">
      <section v-for="row in rows" :key="row.id" class="row">
        <h2 class="row-title">{{ row.title }}</h2>
        <div class="scroller" ref="scrollers">
          <div
            v-for="item in row.items"
            :key="item.id"
            class="card"
            :title="item.title"
            @click="setDemo(item)"
          >
            <img :src="item.poster" :alt="item.title" />
            <div class="card-overlay">
              <div class="card-actions">
                <button class="mini-btn" @click.stop="play(item)">▶</button>
                <button class="mini-btn" @click.stop="like(item)">❤</button>
              </div>
              <div class="card-title">{{ item.title }}</div>
            </div>
          </div>
        </div>
      </section>
    </main>

    <!-- Sticky demo player (для проверки стриминга) -->
    <section v-if="demoVideoId" class="player">
      <div class="player-bar">
        <div class="player-label">Демо-проигрыватель</div>
        <input v-model="demoVideoId" placeholder="video_id" class="player-input" />
        <button class="btn btn-secondary" @click="demoVideoId = ''">Закрыть</button>
      </div>
      <video :src="`/stream/${demoVideoId}`" controls class="video"></video>
    </section>

    <!-- Upload modal -->
    <div v-if="showUpload" class="modal-backdrop" @click.self="showUpload = false">
      <div class="modal">
        <h3>Загрузить видео</h3>
        <form @submit.prevent="uploadVideo" class="form">
          <label class="label">Название</label>
          <input v-model="uploadForm.title" type="text" placeholder="Введите название видео" class="input" />
          <label class="label">Файл</label>
          <input type="file" @change="onFileSelected" accept="video/*" class="input" />
          <div class="modal-actions">
            <button type="button" class="btn btn-secondary" @click="showUpload = false">Отмена</button>
            <button type="submit" class="btn btn-primary" :disabled="!uploadForm.file">Загрузить</button>
          </div>
          <p v-if="uploadResult" class="success">✅ Загружено! ID: <code>{{ uploadResult.video_id }}</code></p>
        </form>
      </div>
    </div>

    <!-- Info modal -->
    <div v-if="showInfo" class="modal-backdrop" @click.self="showInfo = false">
      <div class="modal">
        <h3>{{ featured.title }}</h3>
        <p class="muted">{{ featured.description }}</p>
        <div class="modal-actions">
          <button class="btn btn-primary" @click="playFeatured">▶ Смотреть</button>
          <button class="btn btn-secondary" @click="showInfo = false">Закрыть</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'App',
  data() {
    return {
      showUpload: false,
      showInfo: false,
      uploadForm: {
        title: '',
        file: null
      },
      uploadResult: null,
      demoVideoId: '',
      featured: {
        title: 'Специально для вас',
        description: 'Откройте новые фильмы и сериалы в Just Watch.',
        backdrop: 'https://images.unsplash.com/photo-1497032628192-86f99bcd76bc?q=80&w=1920&auto=format&fit=crop'
      },
      rows: [
        {
          id: 'trending',
          title: 'В тренде',
          items: []
        },
        {
          id: 'top',
          title: 'Топ сегодня',
          items: []
        },
        {
          id: 'new',
          title: 'Новинки',
          items: []
        }
      ]
    }
  },
  computed: {
    heroStyle() {
      return {
        backgroundImage: `url('${this.featured.backdrop}')`
      }
    }
  },
  mounted() {
    // Заглушечные постеры (можно заменить данными из Metadata Service)
    const sample = Array.from({ length: 18 }).map((_, i) => ({
      id: `demo-${i + 1}`,
      title: `Видео ${i + 1}`,
      poster: `https://picsum.photos/seed/jw-${i}/360/540`
    }))
    this.rows = this.rows.map((r, idx) => ({ ...r, items: sample.slice(idx * 6, idx * 6 + 12) }))
  },
  methods: {
    onFileSelected(event) {
      this.uploadForm.file = event.target.files[0]
    },
    async uploadVideo() {
      if (!this.uploadForm.file) return

      const formData = new FormData()
      formData.append('video', this.uploadForm.file)
      formData.append('title', this.uploadForm.title || this.uploadForm.file.name)

      try {
        const response = await fetch('/upload', { method: 'POST', body: formData })
        if (response.ok) {
          this.uploadResult = await response.json()
          this.demoVideoId = this.uploadResult.video_id
          this.uploadForm = { title: '', file: null }
          this.showUpload = false
          alert('Видео успешно загружено!')
        } else {
          const error = await response.json().catch(() => ({}))
          alert('Ошибка загрузки: ' + (error.error || 'Неизвестная ошибка'))
        }
      } catch (err) {
        console.error('Upload error:', err)
        alert('Ошибка сети: ' + err.message)
      }
    },
    playFeatured() {
      if (!this.demoVideoId) {
        // если нет загруженного id — просто подсказка
        alert('Загрузите видео или укажите ID в демо-плеере ниже.')
        return
      }
    },
    play(item) {
      // Используем демо-плеер: предполагаем, что item.id — это video_id (для реальных данных — заменить)
      this.demoVideoId = item.id
    },
    like(item) {
      console.log('LIKE', item.id)
    },
    setDemo(item) {
      this.demoVideoId = item.id
    }
  }
}
</script>

<style>
/* Base */
:root {
  --bg: #141414;
  --muted: #b3b3b3;
  --text: #ffffff;
  --primary: #e50914;
  --card: #2f2f2f;
}

* { box-sizing: border-box; }

html, body, #app {
  height: 100%;
}

body {
  margin: 0;
  background: var(--bg);
  color: var(--text);
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Arial, 'Helvetica Neue', sans-serif;
}

/* Navbar */
.nav {
  position: fixed;
  z-index: 50;
  top: 0;
  left: 0;
  right: 0;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 28px;
  background: linear-gradient(180deg, rgba(0,0,0,0.8), rgba(0,0,0,0));
}

.nav-left {
  display: flex;
  align-items: center;
  gap: 24px;
}

.logo {
  font-weight: 800;
  letter-spacing: 1px;
  color: var(--primary);
}

.nav-links a {
  color: #e5e5e5;
  text-decoration: none;
  margin-right: 16px;
  font-size: 14px;
}

.nav-links a:hover { color: #ffffff; }

.nav-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.avatar {
  width: 32px;
  height: 32px;
  background: #333;
  border-radius: 4px;
  display: grid;
  place-items: center;
  font-size: 12px;
  color: #ccc;
  border: 1px solid #444;
}

.upload-btn {
  background: transparent;
  color: #e5e5e5;
  border: 1px solid #555;
  padding: 6px 12px;
  border-radius: 3px;
  cursor: pointer;
}
.upload-btn:hover { border-color: #777; color: #fff; }

/* Hero */
.hero {
  height: 68vh;
  background-position: center center;
  background-size: cover;
  position: relative;
}

.hero-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0,0,0,0.3);
}

.hero-content {
  position: absolute;
  bottom: 20%;
  left: 60px;
  max-width: 640px;
  z-index: 2;
}

.hero-title {
  font-size: 48px;
  margin: 0 0 12px 0;
  line-height: 1.1;
  text-shadow: 0 2px 12px rgba(0,0,0,0.6);
}

.hero-desc {
  color: #e5e5e5;
  font-size: 18px;
  margin: 0 0 18px 0;
  text-shadow: 0 2px 8px rgba(0,0,0,0.5);
}

.hero-actions { display: flex; gap: 10px; }

.hero-gradient {
  position: absolute;
  bottom: -1px;
  left: 0;
  right: 0;
  height: 160px;
  background: linear-gradient(180deg, rgba(20,20,20,0), rgba(20,20,20,1) 60%);
}

/* Buttons */
.btn {
  border: none;
  padding: 10px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 600;
}

.btn-primary {
  background: var(--primary);
  color: #fff;
}

.btn-secondary {
  background: rgba(255,255,255,0.2);
  color: #fff;
}

/* Rows */
.rows {
  margin-top:  -80px; /* накладываем блоки под градиент hero */
  padding: 0 0 60px 0;
}

.row {
  padding: 12px 0 0 0;
}

.row-title {
  margin: 0 0 8px 60px;
  font-size: 20px;
}

.scroller {
  padding-left: 60px;
  display: grid;
  grid-auto-flow: column;
  grid-auto-columns: 180px;
  gap: 10px;
  overflow-x: auto;
  scroll-snap-type: x mandatory;
}

.scroller::-webkit-scrollbar { height: 8px; }
.scroller::-webkit-scrollbar-track { background: #1b1b1b; }
.scroller::-webkit-scrollbar-thumb { background: #333; border-radius: 8px; }

.card {
  position: relative;
  height: 270px;
  border-radius: 4px;
  overflow: hidden;
  background: var(--card);
  scroll-snap-align: start;
  transition: transform 120ms ease;
  transform-origin: center left;
}

.card:hover {
  transform: scale(1.06);
  z-index: 2;
}

.card img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.card-overlay {
  position: absolute;
  inset: auto 0 0 0;
  height: 42%;
  background: linear-gradient(180deg, rgba(0,0,0,0), rgba(0,0,0,0.85));
  color: #fff;
  padding: 10px;
  display: grid;
  align-content: end;
  gap: 6px;
}

.card-actions {
  display: flex;
  gap: 8px;
}

.mini-btn {
  background: rgba(255,255,255,0.2);
  color: #fff;
  border: none;
  padding: 6px 8px;
  border-radius: 50%;
  cursor: pointer;
}

/* Player */
.player {
  position: sticky;
  bottom: 0;
  background: rgba(0,0,0,0.95);
  padding: 16px 20px;
  border-top: 1px solid #222;
}

.player-bar {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}

.player-label { color: var(--muted); }
.player-input {
  background: #111;
  border: 1px solid #333;
  color: #fff;
  padding: 6px 8px;
  border-radius: 4px;
}

/* Modals */
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.7);
  display: grid;
  place-items: center;
  z-index: 100;
}

.modal {
  width: 520px;
  max-width: calc(100vw - 32px);
  background: #181818;
  border: 1px solid #2a2a2a;
  border-radius: 8px;
  padding: 20px;
  color: #fff;
}

.form { display: grid; gap: 10px; }
.label { color: var(--muted); font-size: 14px; }
.input {
  background: #0f0f0f;
  border: 1px solid #2a2a2a;
  color: #fff;
  padding: 8px 10px;
  border-radius: 4px;
}
.modal-actions { display: flex; justify-content: flex-end; gap: 8px; margin-top: 10px; }
.success { color: #4caf50; }
</style>