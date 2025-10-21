<template>
  <div id="app" style="padding: 20px; font-family: Arial, sans-serif;">
    <h1>🎥 Just Watch — Стримминговый сервис</h1>

    <!-- Форма загрузки -->
    <div style="border: 1px solid #ccc; padding: 16px; margin-bottom: 30px; max-width: 600px;">
      <h2>📤 Загрузить видео</h2>
      <form @submit.prevent="uploadVideo">
        <div style="margin-bottom: 12px;">
          <label>Название:</label><br />
          <input
            v-model="uploadForm.title"
            type="text"
            placeholder="Введите название видео"
            style="width: 100%; padding: 8px; margin-top: 4px;"
          />
        </div>
        <div style="margin-bottom: 12px;">
          <label>Файл:</label><br />
          <input
            type="file"
            @change="onFileSelected"
            accept="video/*"
            style="margin-top: 4px;"
          />
        </div>
        <button
          type="submit"
          :disabled="!uploadForm.file"
          style="padding: 8px 16px; background: #42b883; color: white; border: none; cursor: pointer;"
        >
          Загрузить
        </button>
      </form>

      <p v-if="uploadResult" style="margin-top: 12px; color: green;">
        ✅ Загружено! ID: <code>{{ uploadResult.video_id }}</code>
      </p>
    </div>

    <!-- Список видео (заглушка) -->
    <div>
      <h2>📽️ Доступные видео</h2>
      <p><em>В реальном проекте здесь будет список из Metadata Service.</em></p>

      <!-- Пример видео для теста стриминга -->
      <div style="margin-top: 20px; padding: 12px; border: 1px dashed #999; max-width: 600px;">
        <h3>Демо-видео (если вы загрузили файл, замените ID)</h3>
        <input
          v-model="demoVideoId"
          placeholder="Введите video_id"
          style="padding: 6px; width: 200px; margin-right: 10px;"
        />
        <video
          v-if="demoVideoId"
          :src="`/stream/${demoVideoId}`"
          controls
          style="width: 100%; margin-top: 10px;"
        >
          Ваш браузер не поддерживает видео.
        </video>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'App',
  data() {
    return {
      uploadForm: {
        title: '',
        file: null
      },
      uploadResult: null,
      demoVideoId: '' // для ручного ввода ID
    }
  },
  methods: {
    onFileSelected(event) {
      this.uploadForm.file = event.target.files[0]
    },
    async uploadVideo() {
      if (!this.uploadForm.file) return

      const formData = new FormData()
      formData.append('video', this.uploadForm.file)
      if (this.uploadForm.title) {
        formData.append('title', this.uploadForm.title)
      } else {
        formData.append('title', this.uploadForm.file.name)
      }

      try {
        const response = await fetch('/upload', {
          method: 'POST',
          body: formData
        })

        if (response.ok) {
          this.uploadResult = await response.json()
          alert('Видео успешно загружено!')
          // Сброс формы
          this.uploadForm = { title: '', file: null }
        } else {
          const error = await response.json()
          alert('Ошибка загрузки: ' + (error.error || 'Неизвестная ошибка'))
        }
      } catch (err) {
        console.error('Upload error:', err)
        alert('Ошибка сети: ' + err.message)
      }
    }
  }
}
</script>

<style>
body {
  margin: 0;
  background: #f5f5f5;
}
</style>