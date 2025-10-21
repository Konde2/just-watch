<template>
  <div>
    <h1>🎥 Just Watch</h1>

    <div>
      <h2>Upload Video</h2>
      <form @submit.prevent="upload">
        <input v-model="title" placeholder="Title" />
        <input type="file" @change="onFileChange" accept="video/*" />
        <button type="submit">Upload</button>
      </form>
    </div>

    <div v-if="uploadResult">
      <p>✅ Uploaded! ID: {{ uploadResult.video_id }}</p>
    </div>

    <div v-if="videos.length">
      <h2>Videos</h2>
      <ul>
        <li v-for="video in videos" :key="video.id">
          <strong>{{ video.title }}</strong> ({{ video.status }})
          <video v-if="video.status === 'ready'" :src="`/stream/${video.id}`" controls width="400"></video>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      title: '',
      file: null,
      uploadResult: null,
      videos: []
    }
  },
  async mounted() {
    // В реальном проекте: GET /api/videos → metadata-service
    // Пока заглушка
    this.videos = [
      { id: 'demo', title: 'Demo Video', status: 'ready' }
    ]
  },
  methods: {
    onFileChange(e) {
      this.file = e.target.files[0]
    },
    async upload() {
      if (!this.file) return
      const formData = new FormData()
      formData.append('video', this.file)
      formData.append('title', this.title)

      const res = await fetch('/upload', {
        method: 'POST',
        body: formData
      })
      this.uploadResult = await res.json()
      // В реальности: обновить список видео
    }
  }
}
</script>