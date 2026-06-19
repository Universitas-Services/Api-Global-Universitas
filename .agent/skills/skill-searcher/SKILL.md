---
name: search_antigravity_skills
description: Busca en GitHub nuevas skills, herramientas o integraciones para el IDE Antigravity. Úsala cuando el usuario te pida capacidades que actualmente no tienes.
parameters:
  - name: query
    type: string
    description: La funcionalidad o tecnología que necesitas buscar (ej. "supabase", "docker", "pdf parser").
    required: true
---

# Ejecución

Para ejecutar esta skill, el agente debe usar el siguiente comando:

```bash
python .agents/skills/skill-searcher/search_github.py "{{query}}"
```
