
---

### 📄 `docs/deployment.md` — Деплой

```markdown
# 🌐 Деплой

## Production-сборка

1. Обновите `.env` с production-настройками.
2. Соберите образы:
   ```bash
   docker-compose -f docker-compose.yml -f docker-compose.prod.yml build